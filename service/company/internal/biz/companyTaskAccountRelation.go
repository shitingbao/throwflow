package biz

import (
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
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	ctos "github.com/volcengine/ve-tos-golang-sdk/v2/tos"

	"github.com/go-kratos/kratos/v2/log"
)

var (
	CompanyTaskCreateError                 = errors.InternalServer("COMPANY_TASK_CREATE_ERROR", "领取任务失败")
	CompanyTaskExists                      = errors.InternalServer("COMPANY_TASK_EXISTS", "已存在任务")
	CompanyTaskUpperLimit                  = errors.InternalServer("COMPANY_TASK_UPPER_LIMIT", "任务领取达到上限")
	CompanyTaskRecoverExpireTimeCountError = errors.InternalServer("COMPANY_TASK_RECOVER_EXPIRE_TIME_COUNT_ERROR", "任务恢复数量出错")
	CompanyTaskRelationListError           = errors.InternalServer("COMPANY_TASK_RELATION_LIST_ERROR", "达人任务列表数量出错")
)

type CompanyTaskAccountRelationRepo interface {
	GetById(context.Context, uint64) (*domain.CompanyTaskAccountRelation, error)
	GetCompanyTaskUserOrderStatus(context.Context, uint64, string, string) (*domain.DoukeOrderInfo, error)
	Save(context.Context, *domain.CompanyTaskAccountRelation) (*domain.CompanyTaskAccountRelation, error)
	Count(context.Context, uint64, uint64) (int64, error)
	CountAvailableByTaskId(context.Context, uint64) (int64, error)
	CountByProductOutId(context.Context, uint64, uint64) (int64, error)
	Update(context.Context, *domain.CompanyTaskAccountRelation) (*domain.CompanyTaskAccountRelation, error)
	UpdateStatusByIds(context.Context, int, []int) error
	SoftDelete(context.Context, uint64) error
	List(context.Context, uint64, uint64, int, int, int, string, string, string) ([]*domain.CompanyTaskAccountRelation, error)
	CountByCondition(context.Context, uint64, uint64, int, string, string) (int64, error)
	ListOpenDouyinUsers(context.Context, uint64, uint64, uint64, string) (*v1.ListOpenDouyinUsersReply, error)
	PutContent(context.Context, string, io.Reader) (*ctos.PutObjectV2Output, error)
	// 通过 ClientKey 和 OpenId 查询对应素材数据
	ListVideoTokensOpenDouyinVideos(context.Context, uint64, time.Time, []*domain.CompanyTaskClientKeyAndOpenId) ([]*domain.OpenDouyinVideo, error)
	// UpdateCacheHash 执行 lua 扣减脚本
	UpdateCacheHash(context.Context, string) error
	SaveCache(context.Context, string, time.Duration) bool
	DelCache(context.Context, string) error
}

type CompanyTaskAccountRelationUsecase struct {
	ctrepo  CompanyTaskRepo
	ctdrepo CompanyTaskDetailRepo
	cprepo  CompanyProductRepo
	repo    CompanyTaskAccountRelationRepo
	tm      Transaction
	// qcrepo  QrCodeRepo
	conf  *conf.Data
	vconf *conf.Volcengine
	log   *log.Helper
}

func NewCompanyTaskAccountRelationUsecase(ctrepo CompanyTaskRepo, ctdrepo CompanyTaskDetailRepo, cprepo CompanyProductRepo, repo CompanyTaskAccountRelationRepo, tm Transaction, conf *conf.Data, vconf *conf.Volcengine, logger log.Logger) *CompanyTaskAccountRelationUsecase {
	return &CompanyTaskAccountRelationUsecase{ctrepo: ctrepo, ctdrepo: ctdrepo, cprepo: cprepo, repo: repo, tm: tm, conf: conf, vconf: vconf, log: log.NewHelper(logger)}
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

		product, err := c.cprepo.GetByProductOutId(ctx, productOutId, "", "")

		if err != nil {
			return CompanyTaskCreateError
		}

		tm := time.Now().AddDate(0, 0, int(task.ExpireTime))

		relation := domain.NewCompanyTaskAccountRelation(ctx, companyTaskId, userId, productOutId, product.ProductName)
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

func (c *CompanyTaskAccountRelationUsecase) ListCompanyTaskAccountRelation(ctx context.Context, status int32, companyTaskId, userId, pageNum, pageSize uint64, expireTime, expiredTime, productName string) (*domain.CompanyTaskAccountRelationList, error) {
	list, err := c.repo.List(ctx, companyTaskId, userId, int(pageNum), int(pageSize), int(status), expireTime, expiredTime, productName)

	if err != nil {
		return nil, CompanyTaskRelationListError
	}

	total, err := c.repo.CountByCondition(ctx, companyTaskId, userId, int(status), expireTime, expiredTime)

	if err != nil {
		return nil, CompanyTaskRelationListError
	}

	keys := []*domain.UserOpenDouyin{}
	productOutIds := []uint64{}
	products := make(map[uint64]*domain.CompanyProduct)

	for _, t := range list {
		for _, ds := range t.CompanyTaskDetails {
			keys = append(keys, &domain.UserOpenDouyin{
				ClientKey: ds.ClientKey,
				OpenId:    ds.OpenId,
			})
		}

		productOutIds = append(productOutIds, t.ProductOutId)
	}

	for _, v := range list {
		productOutIds = append(productOutIds, v.ProductOutId)
	}

	companyProducts, err := c.cprepo.ListByProductOutIds(ctx, "1", productOutIds)

	if err != nil {
		return nil, CompanyTaskRelationListError
	}

	for _, companyProduct := range companyProducts {
		products[companyProduct.ProductOutId] = companyProduct
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
	taskIds, err := c.ctrepo.ListIds(ctx)

	if err != nil {
		return err
	}

	// 便利所有任务
	for _, taskId := range taskIds {
		if err := c.syncUpdateCompanyTaskDetailProcess(ctx, taskId); err != nil {
			log.Info("syncUpdateCompanyTaskDetailProcess:", err)
		}
	}
	return nil
}

func (c *CompanyTaskAccountRelationUsecase) syncUpdateCompanyTaskDetailProcess(ctx context.Context, taskId int) error {
	pageNum, pageSize := 1, 40
	for {
		taskInfo, err := c.ctrepo.GetById(ctx, uint64(taskId))

		if err != nil {
			return err
		}

		// 分批次处理领取任务的达人关系，注意是微信信息和任务的关联
		tm := tool.TimeToString("2006-01-02 15:04:05", time.Now())
		relations, err := c.repo.List(ctx, uint64(taskId), 0, pageNum, pageSize, 0, tm, "", "")

		if err != nil {
			return err
		}

		if len(relations) == 0 {
			// 说明已经没有领取的人了
			break
		}

		total, err := c.repo.CountByCondition(ctx, uint64(taskId), 0, 0, tm, "")

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

	if err := c.recoverCompanyTaskExpireTimeCount(ctx, uint64(taskId)); err != nil {
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
		successTaskIds := []int{}

		for _, re := range relations {
			// 获取每个微信对应的抖音信息,这里需要拿出所有，因为有删除关系操作
			openDouyinUser, err := c.repo.ListOpenDouyinUsers(ctx, re.UserId, 0, 40, "")

			if err != nil {
				continue
			}

			clictKeys := []string{}
			openIds := []string{}
			tokens := []*domain.CompanyTaskClientKeyAndOpenId{}

			for _, r := range openDouyinUser.Data.List {
				clictKeys = append(clictKeys, r.ClientKey)
				openIds = append(openIds, r.OpenId)
				tokens = append(tokens, &domain.CompanyTaskClientKeyAndOpenId{
					ClientKey: r.ClientKey,
					OpenId:    r.OpenId,
				})
			}

			if err := c.ctdrepo.DeleteOpenDouyinUsers(ctx, re.UserId, clictKeys, openIds); err != nil {
				return err
			}

			// 抖音信息对应的素材数据
			list, err := c.repo.ListVideoTokensOpenDouyinVideos(ctx, re.ProductOutId, re.ClaimTime, tokens)

			if err != nil {
				return err
			}

			if len(list) == 0 {
				continue
			}

			isCostBuySuccess := re.IsCostBuy > 0

			var isCostBuy uint8 = 0

			order, err := c.repo.GetCompanyTaskUserOrderStatus(ctx, re.UserId, strconv.FormatUint(taskInfo.ProductOutId, 10), "")

			if err == nil {
				if order.FlowPoint != domain.DoukeOrderREFUND {
					isCostBuySuccess = true
					isCostBuy = 1
				} else {
					isCostBuySuccess = false
				}

				re.SetIsCostBuy(ctx, isCostBuy)
				// 如果已经完成，订单状态改变，取消完成状态，因为这里获取的是未过期的，所以不用判断时间
				if re.Status == 1 && isCostBuy == 0 {
					re.SetStatus(ctx, domain.GoingStatus)
				}

				c.repo.Update(ctx, re)
			}

			isSuccess, err := c.createOrUpdateCompanyTaskDetail(ctx, isCostBuySuccess, re.IsScreenshotAvailable, re.Id, re.UserId, re.CompanyTaskId, taskInfo, list)

			if err != nil {
				return err
			}

			if isSuccess {
				successTaskIds = append(successTaskIds, int(re.Id))
			}
		}

		if len(successTaskIds) > 0 {
			err := c.repo.UpdateStatusByIds(ctx, domain.SuccessStatus, successTaskIds)

			if err != nil {
				return err
			}
		}
		return nil
	})

	return err
}

// 获取已经存在的数据（clientKey和openId）更新，并插入
// 判断是否完成任务
// 如果完成，反馈 true
func (c *CompanyTaskAccountRelationUsecase) createOrUpdateCompanyTaskDetail(ctx context.Context, isCostBuySuccess bool, isScreenshotAvailable uint8, relationId, userId, companyTaskId uint64, taskInfo *domain.CompanyTask, list []*domain.OpenDouyinVideo) (bool, error) {
	isSuccess := false
	sourceDetails := make(map[domain.CompanyTaskClientKeyAndOpenId]*domain.OpenDouyinVideo)
	detailConditions := []domain.CompanyTaskClientKeyAndOpenId{}
	userIds := []uint64{userId}

	for _, v := range list {
		if isCostBuySuccess && v.Statistics.PlayCount >= int32(taskInfo.PlayNum) && (taskInfo.IsGoodReviews == 0 || isScreenshotAvailable == 1) {
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
		if detail.PlayCount >= taskInfo.PlayNum {
			detail.IsPlaySuccess = 1
		}

		source := sourceDetails[domain.CompanyTaskClientKeyAndOpenId{
			ClientKey: detail.ClientKey,
			OpenId:    detail.OpenId,
		}]

		if source != nil {
			detail.PlayCount = uint64(source.Statistics.PlayCount)
		}

		existList[domain.CompanyTaskClientKeyAndOpenId{
			ClientKey: detail.ClientKey,
			OpenId:    detail.OpenId,
		}] = true
	}

	for k, v := range sourceDetails {
		if !existList[k] {
			de := &domain.CompanyTaskDetail{
				CompanyTaskId:                companyTaskId,
				CompanyTaskAccountRelationId: relationId,
				ProductName:                  taskInfo.CompanyProduct.ProductName,
				UserId:                       userId,
				ClientKey:                    v.ClientKey,
				OpenId:                       v.OpenId,
				ItemId:                       v.ItemId,
				PlayCount:                    uint64(v.Statistics.PlayCount),
				Cover:                        v.Cover,
				ReleaseTime:                  time.Unix(int64(v.CreateTime), 0),
				Nickname:                     v.Nickname,
				Avatar:                       v.Avatar,
			}

			var isPlaySuccess uint8 = 0

			if v.Statistics.PlayCount >= int32(taskInfo.PlayNum) {
				isPlaySuccess = 1
			}

			de.SetIsReleaseVideo(ctx)
			de.SetCreateTime(ctx)
			de.SetUpdateTime(ctx)
			de.SetIsPlaySuccess(ctx, isPlaySuccess)
			createList = append(createList, de)
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

		ids := []int{}

		for _, v := range list {
			ids = append(ids, int(v.Id))
		}

		if err := c.repo.UpdateStatusByIds(ctx, domain.ExpireStatus, ids); err != nil {
			return CompanyTaskRecoverExpireTimeCountError
		}

		successQuota, err := c.repo.CountByCondition(ctx, taskId, 0, domain.SuccessStatus, "", "")

		if err != nil {
			return CompanyTaskRecoverExpireTimeCountError
		}

		// 找出成功的数量和正在运行的数量，就是领取数量
		goingQuota, err := c.repo.CountByCondition(ctx, taskId, 0, domain.GoingStatus, "", "")

		if err != nil {
			return CompanyTaskRecoverExpireTimeCountError
		}

		tk, err := c.ctrepo.GetById(ctx, taskId)

		if err != nil {
			return CompanyTaskRecoverExpireTimeCountError
		}

		tk.SetClaimQuota(ctx, uint64(goingQuota+successQuota))
		tk.SetSuccessQuota(ctx, uint64(successQuota))
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

	if relation.CompanyTask.IsGoodReviews == 1 {
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
