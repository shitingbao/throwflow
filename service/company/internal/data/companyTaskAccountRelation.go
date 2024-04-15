package data

import (
	douyinv1 "company/api/service/douyin/v1"
	v1 "company/api/service/weixin/v1"
	"company/internal/biz"
	"company/internal/domain"
	"company/internal/pkg/tool"
	"context"
	"encoding/json"
	"errors"
	"io"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	ctos "github.com/volcengine/ve-tos-golang-sdk/v2/tos"
)

const taskTosSubFolder = "task"

var (
	// 对应的可使用种草任务数量减1
	CompanyTaskDeductionLua = `
	local test = tonumber(redis.call('GET', KEYS[1]))
	if test and test > 0 then
		redis.call('DECR', KEYS[1])
		return 1
	else
		return 0
	end`
)

// 种草任务达人关系表
// 一个达人认领后，发布多个视频，对应多个任务明细，先保存关系，再根据发布的视频来获取明细
type CompanyTaskAccountRelation struct {
	Id                    uint64               `gorm:"column:id;primarykey;type:bigint(20) UNSIGNED;autoIncrement;not null;comment:自增ID"`
	CompanyTaskId         uint64               `gorm:"column:company_task_id;type:bigint(20) UNSIGNED;not null;comment:任务ID"`
	ProductOutId          uint64               `gorm:"column:product_out_id;type:bigint(20) UNSIGNED;not null;comment:商品ID"`
	UserId                uint64               `gorm:"column:user_id;type:bigint(20) UNSIGNED;not null;index:idx_user_id;comment:微信小程序用户ID"`
	ClaimTime             time.Time            `gorm:"column:claim_time;type:datetime;not null;comment:认领任务的时间"`
	ExpireTime            time.Time            `gorm:"column:expire_time;type:datetime;not null;comment:过期时间"`
	Status                uint8                `gorm:"column:status;type:tinyint(3) UNSIGNED;not null;default:0;comment:1:任务完成,0:未完成,2:已过期,3:已结算"`
	IsDel                 uint8                `gorm:"column:is_del;type:tinyint(3) UNSIGNED;not null;default:0;comment:是否移除,1:已移除,0:未移除"`
	IsCostBuy             uint8                `gorm:"column:is_cost_buy;type:tinyint(3) UNSIGNED;not null;default:0;comment:成本购买,1:是,0:否"`
	ScreenshotAddress     string               `gorm:"column:screenshot_address;type:varchar(250);not null;comment:截图地址"`
	IsScreenshotAvailable uint8                `gorm:"column:is_screenshot_available;type:tinyint(3) UNSIGNED;not null;default:0;comment:截图是否有效,1:是,0:否"`
	CreateTime            time.Time            `gorm:"column:create_time;type:datetime;not null;comment:新增时间"`
	UpdateTime            time.Time            `gorm:"column:update_time;type:datetime;not null;comment:修改时间"`
	CompanyTaskDetails    []*CompanyTaskDetail `gorm:"foreignKey:CompanyTaskAccountRelationId;references:Id"`
}

func (CompanyTaskAccountRelation) TableName() string {
	return "company_task_account_relation"
}

type companyTaskAccountRelationRepo struct {
	data *Data
	log  *log.Helper
}

func NewCompanyTaskAccountRelationRepo(data *Data, logger log.Logger) biz.CompanyTaskAccountRelationRepo {
	return &companyTaskAccountRelationRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (c *CompanyTaskAccountRelation) ToDomain(ctx context.Context) *domain.CompanyTaskAccountRelation {
	list := []*domain.CompanyTaskDetail{}

	for _, v := range c.CompanyTaskDetails {
		list = append(list, &domain.CompanyTaskDetail{
			Id:                           v.Id,
			CompanyTaskId:                v.CompanyTaskId,
			CompanyTaskAccountRelationId: v.CompanyTaskAccountRelationId,
			UserId:                       v.UserId,
			ClientKey:                    v.ClientKey,
			OpenId:                       v.OpenId,
			ItemId:                       v.ItemId,
			PlayCount:                    v.PlayCount,
			Cover:                        v.Cover,
			ReleaseTime:                  v.ReleaseTime,
			IsPlaySuccess:                v.IsPlaySuccess,
			CreateTime:                   v.CreateTime,
			UpdateTime:                   v.UpdateTime,
			IsReleaseVideo:               v.IsReleaseVideo,
		})
	}

	task := &domain.CompanyTaskAccountRelation{
		Id:                    c.Id,
		CompanyTaskId:         c.CompanyTaskId,
		ProductOutId:          c.ProductOutId,
		UserId:                c.UserId,
		ClaimTime:             c.ClaimTime,
		ExpireTime:            c.ExpireTime,
		Status:                c.Status,
		IsDel:                 c.IsDel,
		CreateTime:            c.CreateTime,
		UpdateTime:            c.UpdateTime,
		IsCostBuy:             c.IsCostBuy,
		ScreenshotAddress:     c.ScreenshotAddress,
		IsScreenshotAvailable: c.IsScreenshotAvailable,
		CompanyTaskDetails:    list,
	}

	return task
}

func (ctr *companyTaskAccountRelationRepo) GetById(ctx context.Context, id uint64) (*domain.CompanyTaskAccountRelation, error) {
	task := &CompanyTaskAccountRelation{}

	if err := ctr.data.db.WithContext(ctx).First(task, id).Error; err != nil {
		return nil, err
	}

	return task.ToDomain(ctx), nil
}

func (ctr *companyTaskAccountRelationRepo) GetByProductOutIdAndUserId(ctx context.Context, productOutId, userId uint64) (*domain.CompanyTaskAccountRelation, error) {
	task := &CompanyTaskAccountRelation{}

	db := ctr.data.db.WithContext(ctx)

	if productOutId > 0 {
		db = db.Where("product_out_id = ?", productOutId)
	}

	if userId > 0 {
		db = db.Where("user_id = ?", userId)
	}

	if err := db.First(task).Error; err != nil {
		return nil, err
	}

	return task.ToDomain(ctx), nil
}

func (ctr *companyTaskAccountRelationRepo) GetUserOrganizationRelations(ctx context.Context, userId uint64) (*v1.GetUserOrganizationRelationsReply, error) {
	return ctr.data.weixinuc.GetUserOrganizationRelations(ctx, &v1.GetUserOrganizationRelationsRequest{
		UserId: userId,
	})
}

// List 反馈基本列表
// status -1 表示没有
func (ctr *companyTaskAccountRelationRepo) List(ctx context.Context, companyTaskId, userId uint64, pageNum, pageSize, status int, expireTime, expiredTime, keyword string) ([]*domain.CompanyTaskAccountRelation, error) {
	list := []*domain.CompanyTaskAccountRelation{}
	taskDetails := []CompanyTaskAccountRelation{}

	db := ctr.data.db.WithContext(ctx).Model(&CompanyTaskAccountRelation{}).
		Where("is_del = 0")

	db = db.Preload("CompanyTaskDetails")

	if companyTaskId > 0 {
		db = db.Where("company_task_id = ?", companyTaskId)
	}

	if status >= 0 {
		if status == 1 {
			db = db.Where("status = 1 or status = 3")
		} else {
			db = db.Where("status = ?", status)
		}
	}

	if userId > 0 {
		db = db.Where("user_id = ?", userId)
	}

	if len(expireTime) > 0 {
		db = db.Where("expire_time > ?", expireTime)
	}

	if len(expiredTime) > 0 {
		// 已经过期的条件
		db = db.Where("expire_time <= ?", expiredTime)
	}

	if pageNum > 0 {
		db = db.Limit(pageSize).Offset((pageNum - 1) * pageSize)
	}

	if err := db.Find(&taskDetails).Error; err != nil {
		return nil, err
	}

	for _, v := range taskDetails {
		list = append(list, v.ToDomain(ctx))
	}

	return list, nil
}

func (ctr *companyTaskAccountRelationRepo) ListByUserIds(ctx context.Context, taskId uint64, userIds []uint64) ([]*domain.CompanyTaskAccountRelation, error) {
	list := []*domain.CompanyTaskAccountRelation{}
	taskDetails := []CompanyTaskAccountRelation{}

	db := ctr.data.db.WithContext(ctx).Model(&CompanyTaskAccountRelation{}).
		Preload("CompanyTaskDetails").
		Where("is_del = 0 and company_task_id = ? and user_id in (?)", taskId, userIds)

	if err := db.Find(&taskDetails).Error; err != nil {
		return nil, err
	}

	for _, v := range taskDetails {
		list = append(list, v.ToDomain(ctx))
	}

	return list, nil
}

func (ctr *companyTaskAccountRelationRepo) ListOpenDouyinUsers(ctx context.Context, userId, pageNum, pageSize uint64, keyword string) (*v1.ListOpenDouyinUsersReply, error) {
	return ctr.data.weixinuc.ListOpenDouyinUsers(ctx, &v1.ListOpenDouyinUsersRequest{
		PageNum:  pageNum,
		PageSize: pageSize,
		UserId:   userId,
		Keyword:  keyword,
	})
}

func (ctr *companyTaskAccountRelationRepo) ListVideoTokensOpenDouyinVideos(ctx context.Context, productOutId uint64, claimTime time.Time, tokens []*domain.CompanyTaskClientKeyAndOpenId) ([]*douyinv1.ListVideoTokensOpenDouyinVideosReply_OpenDouyinVideo, error) {
	tks := []*douyinv1.ListVideoTokensOpenDouyinVideosRequestToken{}

	for _, v := range tokens {
		tks = append(tks, &douyinv1.ListVideoTokensOpenDouyinVideosRequestToken{
			ClientKey: v.ClientKey,
			OpenId:    v.OpenId,
		})
	}

	res, err := ctr.data.douyinuc.ListVideoTokensOpenDouyinVideos(ctx, &douyinv1.ListVideoTokensOpenDouyinVideosRequest{
		ProductOutId: productOutId,
		Tokens:       tks,
		ClaimTime:    tool.TimeToString("2006-01-02 15:04:05", claimTime),
	})

	if err != nil {
		return nil, err
	}

	return res.Data.List, nil
}

// 结算状态为到过期时间，并且状态为完成（status = 1）
func (ctr *companyTaskAccountRelationRepo) ListSettle(ctx context.Context, expiredTime string) ([]*domain.CompanyTaskAccountRelation, error) {
	list := []*domain.CompanyTaskAccountRelation{}
	taskDetails := []CompanyTaskAccountRelation{}

	db := ctr.data.db.WithContext(ctx).Model(&CompanyTaskAccountRelation{}).
		Where("status = 1")

	if len(expiredTime) > 0 {
		// 已经过期的条件
		db = db.Where("expire_time <= ?", expiredTime)
	}

	if err := db.Find(&taskDetails).Error; err != nil {
		return nil, err
	}

	for _, v := range taskDetails {
		list = append(list, v.ToDomain(ctx))
	}

	return list, nil
}

func (ctr *companyTaskAccountRelationRepo) Count(ctx context.Context, companyTaskId, userId uint64) (int64, error) {
	var count int64
	db := ctr.data.db.WithContext(ctx).Model(&CompanyTaskAccountRelation{})

	if companyTaskId > 0 {
		db = db.Where("company_task_id = ?", companyTaskId)
	}

	if userId > 0 {
		db = db.Where("user_id = ?", userId)
	}

	if err := db.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (ctr *companyTaskAccountRelationRepo) CountAvailableByTaskId(ctx context.Context, companyTaskId uint64) (int64, error) {
	var count int64

	err := ctr.data.db.WithContext(ctx).Model(&CompanyTaskAccountRelation{}).
		Where("company_task_id = ? and (expire_time > ? or status = 1 or status = 3)", companyTaskId, time.Now()).
		Count(&count).Error

	if err != nil {
		return 0, err
	}
	return count, nil
}

func (ctr *companyTaskAccountRelationRepo) CountByCondition(ctx context.Context, companyTaskId, userId uint64, status int, expireTime, expiredTime, keyword string) (int64, error) {
	var count int64

	db := ctr.data.db.WithContext(ctx).Model(&CompanyTaskAccountRelation{}).
		Where("is_del = 0")

	if companyTaskId > 0 {
		db = db.Where("company_task_id = ?", companyTaskId)
	}

	if status >= 0 {
		if status == 1 {
			db = db.Where("status = 1 or status = 3")
		} else {
			db = db.Where("status = ?", status)
		}
	}

	if userId > 0 {
		db = db.Where("user_id = ?", userId)
	}

	if len(expireTime) > 0 {
		db = db.Where("expire_time > ?", expireTime)
	}

	if len(expiredTime) > 0 {
		// 已经过期的条件
		db = db.Where("expire_time <= ?", expireTime)
	}

	if err := db.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (ctr *companyTaskAccountRelationRepo) CountByUserIds(ctx context.Context, companyTaskId uint64, userIds []uint64) (int64, error) {
	var count int64

	db := ctr.data.db.WithContext(ctx).Model(&CompanyTaskAccountRelation{}).
		Where("is_del = 0 and company_task_id = ? and user_id in (?)", companyTaskId, userIds)

	if err := db.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (ctr *companyTaskAccountRelationRepo) UpdateStatusByIds(ctx context.Context, status int, ids []uint64) error {
	return ctr.data.DB(ctx).Model(&CompanyTaskAccountRelation{}).Where("id in (?)", ids).Update("status", status).Error
}

func (ctr *companyTaskAccountRelationRepo) Update(ctx context.Context, in *domain.CompanyTaskAccountRelation) (*domain.CompanyTaskAccountRelation, error) {
	task := &CompanyTaskAccountRelation{
		Id:                    in.Id,
		CompanyTaskId:         in.CompanyTaskId,
		ProductOutId:          in.ProductOutId,
		UserId:                in.UserId,
		ClaimTime:             in.ClaimTime,
		ExpireTime:            in.ExpireTime,
		Status:                in.Status,
		IsDel:                 in.IsDel,
		IsCostBuy:             in.IsCostBuy,
		ScreenshotAddress:     in.ScreenshotAddress,
		IsScreenshotAvailable: in.IsScreenshotAvailable,
		CreateTime:            in.CreateTime,
		UpdateTime:            in.UpdateTime,
	}

	if err := ctr.data.DB(ctx).Save(task).Error; err != nil {
		return nil, err
	}

	return task.ToDomain(ctx), nil
}

// UpdateCacheHash 执行 lua 扣减脚本
func (ctr *companyTaskAccountRelationRepo) UpdateCacheHash(ctx context.Context, taskId string) error {
	res, err := ctr.data.rdb.Eval(ctx, CompanyTaskDeductionLua, []string{RedisCompanyTaskPre + taskId}).Result()

	if err != nil {
		return err
	}

	b, err := json.Marshal(res)

	if err != nil {
		return err
	}

	if string(b) == "0" {
		return errors.New("decr task err")
	}

	return nil
}

// 注意扣减的过程，lua 控制 redis 中的可使用数量
// 并查询当前已经被领取到任务数量
func (ctr *companyTaskAccountRelationRepo) Save(ctx context.Context, in *domain.CompanyTaskAccountRelation) (*domain.CompanyTaskAccountRelation, error) {
	relation := &CompanyTaskAccountRelation{
		CompanyTaskId: in.CompanyTaskId,
		ProductOutId:  in.ProductOutId,
		UserId:        in.UserId,
		ClaimTime:     in.ClaimTime,
		ExpireTime:    in.ExpireTime,
		Status:        in.Status,
		CreateTime:    in.CreateTime,
		UpdateTime:    in.UpdateTime,
	}

	if err := ctr.data.DB(ctx).Model(&CompanyTaskAccountRelation{}).Create(relation).Error; err != nil {
		return nil, err
	}

	return relation.ToDomain(ctx), nil
}

func (cor *companyTaskAccountRelationRepo) PutContent(ctx context.Context, fileName string, content io.Reader) (*ctos.PutObjectV2Output, error) {
	for _, ltos := range cor.data.toses {
		if ltos.name == taskTosSubFolder {
			output, err := ltos.tos.PutContent(ctx, fileName, content)

			if err != nil {
				return nil, err
			}

			return output, nil
		}
	}

	return nil, errors.New("tos is not exist")
}

func (ctr *companyTaskAccountRelationRepo) SaveCacheHash(ctx context.Context, keyword string, timeout time.Duration) bool {
	ok, err := ctr.data.rdb.SetNX(ctx, keyword, 1, timeout).Result()

	if err != nil {
		return false
	}

	return ok
}

func (ctr *companyTaskAccountRelationRepo) DelCacheHash(ctx context.Context, keyword string) error {
	ctr.data.rdb.Del(ctx, keyword)

	return nil
}
