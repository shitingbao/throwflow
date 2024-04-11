package biz

import (
	douyinv1 "company/api/service/douyin/v1"
	v1 "company/api/service/weixin/v1"
	"company/internal/conf"
	"company/internal/domain"
	"company/internal/pkg/tool"
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	ctos "github.com/volcengine/ve-tos-golang-sdk/v2/tos"

	"github.com/go-kratos/kratos/v2/log"
)

var (
	CompanyTaskCreateError                 = errors.InternalServer("COMPANY_TASK_CREATE_ERROR", "领取任务失败")
	CompanyTaskCreateLevelError            = errors.InternalServer("COMPANY_TASK_CREATE_LEVEL_ERROR", "非会员，领取任务失败")
	CompanyTaskExists                      = errors.InternalServer("COMPANY_TASK_EXISTS", "已存在任务")
	CompanyTaskUpperLimit                  = errors.InternalServer("COMPANY_TASK_UPPER_LIMIT", "任务领取达到上限")
	CompanyTaskRecoverExpireTimeCountError = errors.InternalServer("COMPANY_TASK_RECOVER_EXPIRE_TIME_COUNT_ERROR", "任务恢复数量出错")
	CompanyTaskRelationListError           = errors.InternalServer("COMPANY_TASK_RELATION_LIST_ERROR", "达人任务列表数量出错")
)

type CompanyTaskAccountRelationRepo interface {
	GetById(context.Context, uint64) (*domain.CompanyTaskAccountRelation, error)
	GetUserOrganizationRelations(context.Context, uint64) (*v1.GetUserOrganizationRelationsReply, error)
	GetByProductOutIdAndUserId(context.Context, uint64, uint64) (*domain.CompanyTaskAccountRelation, error)
	List(context.Context, uint64, uint64, int, int, int, string, string, string) ([]*domain.CompanyTaskAccountRelation, error)
	ListOpenDouyinUsers(context.Context, uint64, uint64, uint64, string) (*v1.ListOpenDouyinUsersReply, error)
	ListVideoTokensOpenDouyinVideos(context.Context, uint64, time.Time, []*domain.CompanyTaskClientKeyAndOpenId) ([]*douyinv1.ListVideoTokensOpenDouyinVideosReply_OpenDouyinVideo, error)
	ListByUserIds(context.Context, uint64, []uint64) ([]*domain.CompanyTaskAccountRelation, error)
	ListSettle(context.Context, string) ([]*domain.CompanyTaskAccountRelation, error)
	Count(context.Context, uint64, uint64) (int64, error)
	CountAvailableByTaskId(context.Context, uint64) (int64, error)
	CountByCondition(context.Context, uint64, uint64, int, string, string, string) (int64, error)
	CountByUserIds(context.Context, uint64, []uint64) (int64, error)
	Update(context.Context, *domain.CompanyTaskAccountRelation) (*domain.CompanyTaskAccountRelation, error)
	UpdateStatusByIds(context.Context, int, []uint64) error
	Save(context.Context, *domain.CompanyTaskAccountRelation) (*domain.CompanyTaskAccountRelation, error)

	PutContent(context.Context, string, io.Reader) (*ctos.PutObjectV2Output, error)
	// UpdateCacheHash 执行 lua 扣减脚本
	UpdateCacheHash(context.Context, string) error
	SaveCache(context.Context, string, time.Duration) bool
	DelCache(context.Context, string) error
}

type CompanyTaskAccountRelationUsecase struct {
	repo    CompanyTaskAccountRelationRepo
	ctrepo  CompanyTaskRepo
	ctdrepo CompanyTaskDetailRepo
	cprepo  CompanyProductRepo
	dorepo  DoukeOrderRepo
	tm      Transaction
	conf    *conf.Data
	vconf   *conf.Volcengine
	log     *log.Helper
}

func NewCompanyTaskAccountRelationUsecase(repo CompanyTaskAccountRelationRepo, ctrepo CompanyTaskRepo, ctdrepo CompanyTaskDetailRepo, cprepo CompanyProductRepo, dorepo DoukeOrderRepo, tm Transaction, conf *conf.Data, vconf *conf.Volcengine, logger log.Logger) *CompanyTaskAccountRelationUsecase {
	return &CompanyTaskAccountRelationUsecase{repo: repo, ctrepo: ctrepo, ctdrepo: ctdrepo, cprepo: cprepo, dorepo: dorepo, tm: tm, conf: conf, vconf: vconf, log: log.NewHelper(logger)}
}

// ListCompanyTaskAccountRelation
// 反馈微信用户对应领取的任务关系，任务信息，商品信息，明细信息
func (c *CompanyTaskAccountRelationUsecase) ListCompanyTaskAccountRelation(ctx context.Context, status int32, companyTaskId, userId, pageNum, pageSize uint64, expireTime, expiredTime, productName string) (*domain.CompanyTaskAccountRelationList, error) {
	list, err := c.repo.List(ctx, companyTaskId, userId, int(pageNum), int(pageSize), int(status), expireTime, expiredTime, productName)

	if err != nil {
		return nil, CompanyTaskRelationListError
	}

	total, err := c.repo.CountByCondition(ctx, companyTaskId, userId, int(status), expireTime, expiredTime, productName)

	if err != nil {
		return nil, CompanyTaskRelationListError
	}

	keys := []*domain.UserOpenDouyin{}
	productOutIds := []uint64{}
	taskIds := []uint64{}
	products := make(map[uint64]*domain.CompanyProduct)
	companyTaskMap := make(map[uint64]*domain.CompanyTask)

	for _, t := range list {
		for _, ds := range t.CompanyTaskDetails {
			keys = append(keys, &domain.UserOpenDouyin{
				ClientKey: ds.ClientKey,
				OpenId:    ds.OpenId,
			})
		}

		productOutIds = append(productOutIds, t.ProductOutId)
		taskIds = append(taskIds, t.CompanyTaskId)
	}

	// 需要用到商品信息
	companyProducts, err := c.cprepo.ListByProductOutIds(ctx, "1", productOutIds)

	if err != nil {
		return nil, CompanyTaskRelationListError
	}

	for _, companyProduct := range companyProducts {
		products[companyProduct.ProductOutId] = companyProduct
	}

	// 需要用到任务中的价格
	companyTasks, err := c.ctrepo.ListByIds(ctx, taskIds)

	if err != nil {
		return nil, CompanyTaskRelationListError
	}

	for _, companyTask := range companyTasks {
		companyTaskMap[companyTask.ProductOutId] = companyTask
	}

	b, err := json.Marshal(keys)

	if err != nil {
		return nil, CompanyTaskGetDouyinUserError
	}

	users, err := c.ctdrepo.ListByClientKeyAndOpenIds(ctx, 0, 40, string(b), "")

	if err != nil {
		return nil, CompanyTaskGetDouyinUserError
	}

	for _, t := range list {
		for _, d := range t.CompanyTaskDetails {
			d.SetNicknameAndAvatar(ctx, users.Data.List)
		}

		if companyTaskMap[t.ProductOutId] != nil {
			t.CompanyTask = *(companyTaskMap[t.ProductOutId])
		}

		companyProduct := products[t.ProductOutId]

		if companyProduct != nil {
			companyProduct.SetProductImgs(ctx)

			if len(companyProduct.ProductImgs) > 0 {
				companyProduct.SetProductImg(ctx, companyProduct.ProductImgs[0])
			}

			t.CompanyProduct = companyProduct
		}

	}

	return &domain.CompanyTaskAccountRelationList{
		PageNum:  pageNum,
		PageSize: pageSize,
		Total:    uint64(total),
		List:     list,
	}, nil
}

// one people only one task
// when not find task key in redis, get this task then set
func (c *CompanyTaskAccountRelationUsecase) CreateCompanyTaskAccountRelation(ctx context.Context, companyTaskId, productOutId, userId uint64) (*domain.CompanyTaskAccountRelation, error) {
	keyword := strconv.FormatUint(companyTaskId, 10) + ":" + strconv.FormatUint(productOutId, 10) + ":" + strconv.FormatUint(userId, 10)

	if !c.repo.SaveCache(ctx, keyword, c.conf.Redis.ProductLockTimeout.AsDuration().Abs()) {
		return nil, CompanyTaskExists
	}

	defer c.repo.DelCache(ctx, keyword)

	count, err := c.repo.Count(ctx, companyTaskId, userId) // 验证同一个微信是否已经领取任务

	if err != nil {
		return nil, CompanyTaskCreateError
	}

	if count > 0 {
		return nil, CompanyTaskExists
	}

	res, err := c.repo.GetUserOrganizationRelations(ctx, userId)

	if err != nil {
		return nil, CompanyTaskCreateLevelError
	}

	if res.Data.Level <= 0 {
		return nil, CompanyTaskCreateLevelError
	}

	_, err = c.ctrepo.GetCacheHash(ctx, strconv.FormatUint(companyTaskId, 10))

	if err != nil {
		// 如果 key 丢失，等待定时任务恢复，恢复之前无法领取
		return nil, CompanyTaskUpperLimit
	}

	rel := &domain.CompanyTaskAccountRelation{}

	err = c.tm.InTx(ctx, func(ctx context.Context) error {
		task, err := c.ctrepo.GetById(ctx, companyTaskId)

		if err != nil {
			return CompanyTaskCreateError
		}

		availCount, err := c.repo.CountAvailableByTaskId(ctx, companyTaskId)

		if err != nil {
			return CompanyTaskCreateError
		}

		if availCount >= int64(task.Quota) {
			return CompanyTaskUpperLimit
		}

		tm := time.Now().AddDate(0, 0, int(task.ExpireTime))

		relation := domain.NewCompanyTaskAccountRelation(ctx, companyTaskId, userId, productOutId)
		relation.SetClaimTime(ctx)
		relation.SetCreateTime(ctx)
		relation.SetUpdateTime(ctx)
		relation.SetExpireTime(ctx, tm)

		rel, err = c.repo.Save(ctx, relation)

		if err != nil {
			return CompanyTaskCreateError
		}

		if err := c.repo.UpdateCacheHash(ctx, strconv.FormatUint(companyTaskId, 10)); err != nil {
			return CompanyTaskCreateError
		}

		return nil
	})

	return rel, err
}

func (c *CompanyTaskAccountRelationUsecase) UpdateScreenshotById(ctx context.Context, relationId uint64, screenshot string) (*domain.CompanyTaskAccountRelation, error) {
	tk, err := c.repo.GetById(ctx, relationId)

	if err != nil {
		return nil, CompanyTaskDetailUpdateError
	}

	screenshots := strings.Split(screenshot, ",")

	if len(screenshots) != 2 {
		return nil, CompanyTaskDetailUpdateError
	}

	if _, ok := Mime[screenshots[0][5:len(screenshots[0])-7]]; !ok {
		return nil, CompanyTaskDetailUpdateError
	}

	imagePath := c.vconf.Tos.Task.SubFolder + "/" + tool.GetRandCode(time.Now().String()) + Mime[screenshots[0][5:len(screenshots[0])-7]]
	imageContent, err := base64.StdEncoding.DecodeString(screenshots[1])

	if err != nil {
		return nil, CompanyTaskDetailUpdateError
	}

	if _, err = c.repo.PutContent(ctx, imagePath, strings.NewReader(string(imageContent))); err != nil {
		return nil, CompanyTaskDetailUpdateError
	}

	tk.SetIsScreenshotAvailable(ctx, domain.ScreenshotAvailable)
	tk.SetScreenshotAddress(ctx, c.vconf.Tos.Task.Url+"/"+imagePath)
	tk.SetUpdateTime(ctx)

	task, err := c.repo.Update(ctx, tk)

	if err != nil {
		return nil, CompanyTaskDetailUpdateError
	}

	return task, nil
}

func (c *CompanyTaskAccountRelationUsecase) UpdateScreenshotAvailableById(ctx context.Context, isScreenshotAvailable uint8, relationId uint64) (*domain.CompanyTaskAccountRelation, error) {
	relation, err := c.repo.GetById(ctx, relationId)

	if err != nil {
		return nil, CompanyTaskDetailUpdateError
	}

	relation.SetIsScreenshotAvailable(ctx, isScreenshotAvailable)
	relation.SetUpdateTime(ctx)

	task, err := c.ctrepo.GetById(ctx, relation.CompanyTaskId)

	if err != nil {
		return nil, CompanyTaskDetailUpdateError
	}

	if task.IsGoodReviews == 1 {
		// 需要好评情况下
		// 取消截图有效时，判断当前的任务是否时过期状态
		statusFlag := relation.ExpireTime.After(time.Now())
		var status uint8 = domain.GoingStatus

		if !statusFlag {
			status = domain.ExpireStatus
		}

		if relation.Status == 1 && isScreenshotAvailable == 0 {
			relation.SetStatus(ctx, status)
		}
	}

	newRelation, err := c.repo.Update(ctx, relation)

	if err != nil {
		return nil, CompanyTaskDetailUpdateError
	}

	return newRelation, nil
}

// 定时任务更新任务详情
// 1.获取需要更新的任务关系，达人领取后未过期且未完成
// 2.获取 购买，播放率，截图等，更新数据，完成标识
// 3.过期的任务数进行恢复，加入 redis，并在 mysql 中标识失败
func (c *CompanyTaskAccountRelationUsecase) SyncUpdateCompanyTaskDetail(ctx context.Context) error {
	// 根据任务循环
	// 获取任务对应关系
	// 根据关系去拉取·播放量·，视频id，视频封面，发布时间，成本购买，好评达标，视频发布，播放量达标，截图地址，任务结果
	// 更新任务结果，达人和任务的关系，达人视频的成功结果
	// 检查对应任务的 redis key 值
	// 将过期未完成的任务关系标记为失败（已过期），将任务数重新写回 redis 中
	tasks, err := c.ctrepo.List(ctx, 1, 40, 0, -1, []uint64{})

	if err != nil {
		return err
	}

	wg := &sync.WaitGroup{}

	for _, task := range tasks {
		wg.Add(1)

		go c.syncUpdateCompanyTaskDetailProcess(ctx, wg, task)
	}

	wg.Wait()

	return nil
}

func (c *CompanyTaskAccountRelationUsecase) syncUpdateCompanyTaskDetailProcess(ctx context.Context, wg *sync.WaitGroup, taskInfo *domain.CompanyTask) error {
	defer wg.Done()

	pageNum, pageSize := 1, 40

	for {
		product, err := c.cprepo.GetByProductOutId(ctx, taskInfo.ProductOutId, "", "")

		if err != nil {
			return err
		}

		if product != nil {
			taskInfo.SetCompanyProduct(ctx, *product)
		}

		// 分批次处理领取任务的达人关系，注意是微信信息和任务的关联
		expriedTime := tool.TimeToString("2006-01-02 15:04:05", time.Now())

		relations, err := c.repo.List(ctx, taskInfo.Id, 0, pageNum, pageSize, -1, expriedTime, "", "")

		if err != nil {
			return err
		}

		if len(relations) == 0 {
			// 说明已经没有领取的人了
			break
		}

		total, err := c.repo.CountByCondition(ctx, taskInfo.Id, 0, 0, expriedTime, "", "")

		if err != nil {
			return err
		}

		if err := c.companyTaskDetailRelationsProcess(ctx, taskInfo, relations); err != nil {
			return err
		}

		if pageSize*pageNum >= int(total) {
			break
		}

		pageNum++
	}

	if err := c.recoverCompanyTaskExpireTimeCount(ctx, taskInfo.Id); err != nil {
		return err
	}

	return nil
}

// 达人领取任务处理过程
// 获取微信号对应的抖音账号信息列表，有分页
// 删除不在最新微信号对应的抖音账号关系中的任务详情数据（因为微信对应抖音账号吧绑定关系会变动）
// 进行素材数据录入
func (c *CompanyTaskAccountRelationUsecase) companyTaskDetailRelationsProcess(ctx context.Context, taskInfo *domain.CompanyTask, relations []*domain.CompanyTaskAccountRelation) error {
	// 这里的关系就是每个微信的信息
	// 每次处理提交一次
	err := c.tm.InTx(ctx, func(ctx context.Context) error {
		successTaskIds := []uint64{}

		for _, re := range relations {
			// 先获取成本购信息
			var isCostBuy uint8 = 0
			tokens := []*domain.CompanyTaskClientKeyAndOpenId{}
			tokenMap := make(map[domain.CompanyTaskClientKeyAndOpenId]bool)
			videoIdMap := make(map[string]bool)
			deleteIds := []uint64{}

			claimTime := tool.TimeToString("2006-01-02 15:04:05", re.ClaimTime)
			// 获取一条购买成功的成本购
			doukeOrder, err := c.dorepo.Get(ctx, re.UserId, strconv.FormatUint(taskInfo.ProductOutId, 10), domain.DoukeOrderREFUND, claimTime)

			if err == nil {
				if doukeOrder.Data.FlowPoint != "" && doukeOrder.Data.FlowPoint != domain.DoukeOrderREFUND {
					isCostBuy = 1
				}

				re.SetIsCostBuy(ctx, isCostBuy)
				re.SetUpdateTime(ctx)

				// 如果已经完成，订单状态改变，取消完成状态，因为这里获取的是未过期的，所以不用判断时间
				if re.Status == 1 && isCostBuy == 0 {
					re.SetStatus(ctx, domain.GoingStatus)
				}

				c.repo.Update(ctx, re)
			}

			// 获取每个微信对应的抖音信息,这里需要拿出所有，因为有删除关系操作
			openDouyinUser, err := c.repo.ListOpenDouyinUsers(ctx, re.UserId, 0, 40, "")

			if err != nil {
				continue
			}

			for _, r := range openDouyinUser.Data.List {
				tokens = append(tokens, &domain.CompanyTaskClientKeyAndOpenId{
					ClientKey: r.ClientKey,
					OpenId:    r.OpenId,
				})

				tokenMap[domain.CompanyTaskClientKeyAndOpenId{
					ClientKey: r.ClientKey,
					OpenId:    r.OpenId,
				}] = true
			}

			oldDetails, err := c.ctdrepo.List(ctx, 0, 40, re.CompanyTaskId, []uint64{re.UserId}, []domain.CompanyTaskClientKeyAndOpenId{})

			if err != nil {
				continue
			}

			// 抖音信息对应的素材数据
			list, err := c.repo.ListVideoTokensOpenDouyinVideos(ctx, re.ProductOutId, re.ClaimTime, tokens)

			if err != nil {
				log.Error("ListVideoTokensOpenDouyinVideos:", err)
				continue
			}

			for _, v := range list {
				videoIdMap[v.VideoId] = true
			}

			for _, detail := range oldDetails {
				// 人员关系不存在，或者视频数据不存在，都删除明细数据
				if !tokenMap[domain.CompanyTaskClientKeyAndOpenId{
					ClientKey: detail.ClientKey,
					OpenId:    detail.OpenId,
				}] {
					deleteIds = append(deleteIds, detail.Id)
				}

				if !videoIdMap[detail.VideoId] {
					deleteIds = append(deleteIds, detail.Id)
				}
			}

			if len(deleteIds) > 0 {
				c.ctdrepo.DeleteOpenDouyinUsers(ctx, deleteIds)
			}

			isSuccess, err := c.createOrUpdateCompanyTaskDetail(ctx, isCostBuy == 1, re.IsScreenshotAvailable, re.Id, re.UserId, re.CompanyTaskId, taskInfo, list)

			if err != nil {
				continue
			}

			// 插入后查看该用户的该任务是否是已经完成任务后，视频数据不达标的（没有或者剩下的播放量不够）
			count, err := c.ctdrepo.CountIsPlauSuccess(ctx, re.CompanyTaskId, re.UserId)

			if err != nil {
				continue
			}

			if re.Status == domain.SuccessStatus && count == 0 {
				// 如果完成后没有符合的视频，清除完成状态
				c.repo.UpdateStatusByIds(ctx, domain.GoingStatus, []uint64{re.Id})
			}

			if isSuccess {
				successTaskIds = append(successTaskIds, re.Id)
			}
		}

		if len(successTaskIds) > 0 {
			c.repo.UpdateStatusByIds(ctx, domain.SuccessStatus, successTaskIds)
		}

		return nil
	})

	return err
}

// 获取已经存在的数据（clientKey和openId）更新，并插入
// 判断是否完成任务
// 如果完成，反馈 true
func (c *CompanyTaskAccountRelationUsecase) createOrUpdateCompanyTaskDetail(ctx context.Context, isCostBuySuccess bool, isScreenshotAvailable uint8, relationId, userId, companyTaskId uint64, taskInfo *domain.CompanyTask, list []*douyinv1.ListVideoTokensOpenDouyinVideosReply_OpenDouyinVideo) (bool, error) {
	isSuccess := false
	sourceDetails := make(map[domain.CompanyTaskClientKeyAndOpenId]*douyinv1.ListVideoTokensOpenDouyinVideosReply_OpenDouyinVideo)
	detailConditions := []domain.CompanyTaskClientKeyAndOpenId{}
	userIds := []uint64{userId}

	for _, v := range list {
		if isCostBuySuccess && uint64(v.Statistics.PlayCount) >= (taskInfo.PlayNum) && (taskInfo.IsGoodReviews == 0 || isScreenshotAvailable == 1) {
			// 判断是否完成任务
			isSuccess = true
		}

		detailCondition := domain.CompanyTaskClientKeyAndOpenId{
			ClientKey: v.ClientKey,
			OpenId:    v.OpenId,
		}

		detailConditions = append(detailConditions, detailCondition)
		sourceDetails[domain.CompanyTaskClientKeyAndOpenId{
			ClientKey: v.ClientKey,
			OpenId:    v.OpenId,
			VideoId:   v.VideoId,
		}] = v
	}

	createList := []*domain.CompanyTaskDetail{}
	// 先根据 clientKey 和 openId 查出本地有的数据，用于更新
	updateList, err := c.ctdrepo.List(ctx, 0, 0, companyTaskId, userIds, detailConditions)

	if err != nil {
		return false, err
	}

	existList := make(map[domain.CompanyTaskClientKeyAndOpenId]bool)

	for _, detail := range updateList {
		source := sourceDetails[domain.CompanyTaskClientKeyAndOpenId{
			ClientKey: detail.ClientKey,
			OpenId:    detail.OpenId,
			VideoId:   detail.VideoId,
		}]

		if source != nil {
			detail.SetPlayCount(ctx, uint64(source.Statistics.PlayCount))
		}

		// 更新要重置该值
		var isPlaySuccess uint8 = 0

		if detail.PlayCount >= taskInfo.PlayNum {
			isPlaySuccess = 1
		}

		detail.SetIsPlaySuccess(ctx, isPlaySuccess)
		detail.SetUpdateTime(ctx)

		existList[domain.CompanyTaskClientKeyAndOpenId{
			ClientKey: detail.ClientKey,
			OpenId:    detail.OpenId,
			VideoId:   detail.VideoId,
		}] = true
	}

	for k, v := range sourceDetails {
		if !existList[k] {
			releaseTime, err := tool.StringToTime("2006-01-02 15:04:05", v.CreateTime)

			if err != nil {
				continue
			}

			detail := domain.NewCompanyTaskDetail(ctx, companyTaskId, relationId, userId, uint64(v.Statistics.PlayCount), taskInfo.CompanyProduct.ProductName, v.ClientKey, v.OpenId, v.ItemId, v.Cover, v.Nickname, v.Avatar, releaseTime)

			var isPlaySuccess uint8 = 0

			if uint64(v.Statistics.PlayCount) >= (taskInfo.PlayNum) {
				isPlaySuccess = 1
			}

			detail.SetVideoId(ctx, v.VideoId)
			detail.SetIsReleaseVideo(ctx)
			detail.SetCreateTime(ctx)
			detail.SetUpdateTime(ctx)
			detail.SetIsPlaySuccess(ctx, isPlaySuccess)

			createList = append(createList, detail)
		}
	}

	if len(updateList) > 0 {
		if err := c.ctdrepo.UpdateOnDuplicateKey(ctx, updateList); err != nil {
			return false, err
		}
	}

	if len(createList) > 0 {
		if err := c.ctdrepo.SaveList(ctx, createList); err != nil {
			return false, err
		}
	}

	return isSuccess, nil
}

// 恢复过期的任务数量，并标记为过期状态
// 检查对应任务的 redis key 值
// 更新任务领取数量和完成数量
// redis 中恢复可用数量
func (c *CompanyTaskAccountRelationUsecase) recoverCompanyTaskExpireTimeCount(ctx context.Context, taskId uint64) error {
	err := c.tm.InTx(ctx, func(ctx context.Context) error {
		// 获取过期的数量
		list, err := c.repo.List(ctx, taskId, 0, 0, 0, 0, "", tool.TimeToString("2006-01-02 15:04:05", time.Now()), "")

		if err != nil {
			return CompanyTaskRecoverExpireTimeCountError
		}

		ids := []uint64{}

		for _, v := range list {
			ids = append(ids, v.Id)
		}

		if len(ids) > 0 {
			c.repo.UpdateStatusByIds(ctx, domain.ExpireStatus, ids)
		}

		successQuota, err := c.repo.CountByCondition(ctx, taskId, 0, domain.SuccessStatus, "", "", "")

		if err != nil {
			return CompanyTaskRecoverExpireTimeCountError
		}

		// 找出成功的数量和正在运行的数量，就是领取数量
		goingQuota, err := c.repo.CountByCondition(ctx, taskId, 0, domain.GoingStatus, "", "", "")

		if err != nil {
			return CompanyTaskRecoverExpireTimeCountError
		}

		tk, err := c.ctrepo.GetById(ctx, taskId)

		if err != nil {
			return CompanyTaskRecoverExpireTimeCountError
		}

		tk.SetClaimQuota(ctx, uint64(goingQuota+successQuota))
		tk.SetSuccessQuota(ctx, uint64(successQuota))
		tk.SetUpdateTime(ctx)

		_, err = c.ctrepo.Update(ctx, tk)

		if err != nil {
			return err
		}

		_, err = c.ctrepo.GetCacheHash(ctx, strconv.FormatUint(taskId, 10))

		if err != nil {
			// 如果丢失，并且任务没有被关闭，重新生成
			tk, taskErr := c.ctrepo.GetById(ctx, taskId)

			if taskErr != nil {
				return CompanyTaskRecoverExpireTimeCountError
			}

			ct, err := c.repo.CountAvailableByTaskId(ctx, taskId)

			if err != nil {
				return CompanyTaskRecoverExpireTimeCountError
			}

			if err := c.ctrepo.SaveCacheHash(ctx, strconv.FormatUint(taskId, 10), tk.Quota-uint64(ct)); err != nil {
				return CompanyTaskRecoverExpireTimeCountError
			}

			return nil
		} else {
			if err := c.ctrepo.UpdateCacheHash(ctx, strconv.FormatUint(taskId, 10), int64(len(list))); err != nil {
				return CompanyTaskRecoverExpireTimeCountError
			}
		}
		return nil
	})

	return err
}

// 结算已经完成任务的金额
func (ctar *CompanyTaskAccountRelationUsecase) settlePrice(ctx context.Context) {
	relations, err := ctar.repo.ListSettle(ctx, tool.TimeToString("2006-01-02 15:04:05", time.Now()))

	if err != nil {
		return
	}

	if len(relations) == 0 {
		return
	}

	settleIds := []uint64{}

	// uint64
	for _, relation := range relations {
		// 结算方法
		// func settle
		settleIds = append(settleIds, relation.Id)
	}

	if len(settleIds) > 0 {
		ctar.repo.UpdateStatusByIds(ctx, domain.SettledStatus, settleIds)
	}
}
