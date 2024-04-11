package data

import (
	v1 "company/api/service/weixin/v1"
	"company/internal/biz"
	"company/internal/domain"
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

// 种草任务明细表
// 一个任务关系，对应该表多个视频明细
type CompanyTaskDetail struct {
	Id                           uint64    `gorm:"column:id;primarykey;type:bigint(20) UNSIGNED;autoIncrement;not null;comment:自增ID"`
	CompanyTaskId                uint64    `gorm:"column:company_task_id;type:bigint(20) UNSIGNED;not null;comment:任务ID"`
	UserId                       uint64    `gorm:"column:user_id;type:bigint(20) UNSIGNED;not null;index:idx_user_id;comment:微信小程序用户ID"`
	ClientKey                    string    `gorm:"column:client_key;type:varchar(50);not null;comment:抖音开放平台应用Client Key"`
	OpenId                       string    `gorm:"column:open_id;type:varchar(100);not null;comment:抖音开放平台授权用户唯一标识"`
	VideoId                      string    `gorm:"column:video_id;type:varchar(100);not null;comment:视频id"`
	CompanyTaskAccountRelationId uint64    `gorm:"column:company_task_account_relation_id;type:bigint(20) UNSIGNED;not null;comment:任务关系ID"`
	ItemId                       string    `gorm:"column:item_id;type:varchar(250);not null;comment:视频id"` // 视频id // video_id
	PlayCount                    uint64    `gorm:"column:play_count;type:int(10);not null;comment:播放数"`    // 播放数
	Cover                        string    `gorm:"column:cover;type:text;not null;comment:视频封面"`           // 视频封面
	ReleaseTime                  time.Time `gorm:"column:release_time;type:datetime;not null;comment:发布时间"`
	IsReleaseVideo               uint8     `gorm:"column:is_release_video;type:tinyint(3) UNSIGNED;not null;default:0;comment:视频发布,1:是,0:否"`
	IsPlaySuccess                uint8     `gorm:"column:is_play_success;type:tinyint(3) UNSIGNED;not null;default:0;comment:播放量达标,1:是,0:否"`
	CreateTime                   time.Time `gorm:"column:create_time;type:datetime;not null;comment:新增时间"`
	UpdateTime                   time.Time `gorm:"column:update_time;type:datetime;not null;comment:修改时间"`
}

func (CompanyTaskDetail) TableName() string {
	return "company_task_detail"
}

func (c *CompanyTaskDetail) ToDomain(ctx context.Context) *domain.CompanyTaskDetail {
	task := &domain.CompanyTaskDetail{
		Id:                           c.Id,
		CompanyTaskId:                c.CompanyTaskId,
		CompanyTaskAccountRelationId: c.CompanyTaskAccountRelationId,
		ItemId:                       c.ItemId,
		PlayCount:                    c.PlayCount,
		Cover:                        c.Cover,
		ReleaseTime:                  c.ReleaseTime,
		IsPlaySuccess:                c.IsPlaySuccess,
		CreateTime:                   c.CreateTime,
		UpdateTime:                   c.UpdateTime,
		UserId:                       c.UserId,
		ClientKey:                    c.ClientKey,
		OpenId:                       c.OpenId,
		IsReleaseVideo:               c.IsReleaseVideo,
		VideoId:                      c.VideoId,
	}
	return task
}

type companyTaskDetailRepo struct {
	data *Data
	log  *log.Helper
}

func NewCompanyTaskDetailRepo(data *Data, logger log.Logger) biz.CompanyTaskDetailRepo {
	return &companyTaskDetailRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (ctr *companyTaskDetailRepo) GetById(ctx context.Context, id uint64) (*domain.CompanyTaskDetail, error) {
	task := &CompanyTaskDetail{}

	if err := ctr.data.db.WithContext(ctx).First(task, id).Error; err != nil {
		return nil, err
	}

	return task.ToDomain(ctx), nil
}

func (ctr *companyTaskDetailRepo) List(ctx context.Context, pageNum, pageSize int, taskId uint64, userIds []uint64, clientKeyAndOpenIds []domain.CompanyTaskClientKeyAndOpenId) ([]*domain.CompanyTaskDetail, error) {
	list := []*domain.CompanyTaskDetail{}
	task := []CompanyTaskDetail{}

	db := ctr.data.db.WithContext(ctx).Model(&CompanyTaskDetail{})

	if taskId > 0 {
		db = db.Where("company_task_id = ?", taskId)
	}

	if len(clientKeyAndOpenIds) > 0 {
		vals := []interface{}{}
		orConditions := ""

		for i, v := range clientKeyAndOpenIds {
			if i == 0 {
				orConditions = "(client_key = ? AND open_id = ?)"
			} else {
				orConditions += " OR (client_key = ? AND open_id = ?)"
			}

			vals = append(vals, v.ClientKey, v.OpenId)
		}

		db = db.Where(orConditions, vals...)
	}

	if len(userIds) > 0 {
		db = db.Where("user_id in (?)", userIds)
	}

	if pageNum > 0 {
		db = db.
			Limit(pageSize).Offset((pageNum - 1) * pageSize)
	}

	if err := db.Find(&task).Error; err != nil {
		return nil, err
	}

	for _, t := range task {
		list = append(list, t.ToDomain(ctx))
	}

	return list, nil
}

func (ctr *companyTaskDetailRepo) Count(ctx context.Context, taskId uint64, userIds []uint64, clientKeyAndOpenIds []domain.CompanyTaskClientKeyAndOpenId) (int64, error) {
	db := ctr.data.db.WithContext(ctx).Model(&CompanyTaskDetail{})

	if taskId > 0 {
		db = db.Where("company_task_id = ?", taskId)
	}

	if len(clientKeyAndOpenIds) > 0 {
		vals := []interface{}{}
		orConditions := ""

		for i, v := range clientKeyAndOpenIds {
			if i == 0 {
				orConditions = "(client_key = ? AND open_id = ?)"
			} else {
				orConditions += " OR (client_key = ? AND open_id = ?)"
			}
			vals = append(vals, v.ClientKey, v.OpenId)
		}
		db = db.Where(orConditions, vals...)
	}

	if len(userIds) > 0 {
		db = db.Where("user_id in (?)", userIds)
	}

	var count int64

	if err := db.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (ctr *companyTaskDetailRepo) CountIsPlauSuccess(ctx context.Context, taskId, userId uint64) (int64, error) {
	db := ctr.data.DB(ctx).Model(&CompanyTaskDetail{}).
		Where("company_task_id = ? and user_id = ? and is_play_success = 1", taskId, userId)

	var count int64

	if err := db.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (ctr *companyTaskDetailRepo) Save(ctx context.Context, in *domain.CompanyTaskDetail) (*domain.CompanyTaskDetail, error) {
	detail := &CompanyTaskDetail{
		CompanyTaskId:                in.CompanyTaskId,
		UserId:                       in.UserId,
		ClientKey:                    in.ClientKey,
		OpenId:                       in.OpenId,
		VideoId:                      in.VideoId,
		CompanyTaskAccountRelationId: in.CompanyTaskAccountRelationId,
		ItemId:                       in.ItemId,
		PlayCount:                    in.PlayCount,
		Cover:                        in.Cover,
		ReleaseTime:                  in.ReleaseTime,
		IsPlaySuccess:                in.IsPlaySuccess,
		IsReleaseVideo:               1,
		CreateTime:                   in.CreateTime,
		UpdateTime:                   in.UpdateTime,
	}

	if err := ctr.data.DB(ctx).Model(&CompanyTaskDetail{}).Create(detail).Error; err != nil {
		return nil, err
	}

	return detail.ToDomain(ctx), nil
}

func (ctr *companyTaskDetailRepo) SaveList(ctx context.Context, ins []*domain.CompanyTaskDetail) error {
	details := []*CompanyTaskDetail{}

	for _, in := range ins {
		detail := &CompanyTaskDetail{
			CompanyTaskId:                in.CompanyTaskId,
			UserId:                       in.UserId,
			ClientKey:                    in.ClientKey,
			OpenId:                       in.OpenId,
			VideoId:                      in.VideoId,
			CompanyTaskAccountRelationId: in.CompanyTaskAccountRelationId,
			ItemId:                       in.ItemId,
			PlayCount:                    in.PlayCount,
			Cover:                        in.Cover,
			ReleaseTime:                  in.ReleaseTime,
			IsPlaySuccess:                in.IsPlaySuccess,
			IsReleaseVideo:               1,
			CreateTime:                   in.CreateTime,
			UpdateTime:                   in.UpdateTime,
		}
		details = append(details, detail)
	}

	return ctr.data.DB(ctx).Model(&CompanyTaskDetail{}).Create(details).Error
}

func (ctr *companyTaskDetailRepo) Update(ctx context.Context, in *domain.CompanyTaskDetail) (*domain.CompanyTaskDetail, error) {
	task := &CompanyTaskDetail{
		Id:                           in.Id,
		CompanyTaskId:                in.CompanyTaskId,
		CompanyTaskAccountRelationId: in.CompanyTaskAccountRelationId,
		ItemId:                       in.ItemId,
		PlayCount:                    in.PlayCount,
		Cover:                        in.Cover,
		ReleaseTime:                  in.ReleaseTime,
		IsPlaySuccess:                in.IsPlaySuccess,
		CreateTime:                   in.CreateTime,
		UpdateTime:                   in.UpdateTime,
		UserId:                       in.UserId,
		ClientKey:                    in.ClientKey,
		OpenId:                       in.OpenId,
	}

	if err := ctr.data.DB(ctx).Save(task).Error; err != nil {
		return nil, err
	}

	return task.ToDomain(ctx), nil
}

func (ctr *companyTaskDetailRepo) UpdateOnDuplicateKey(ctx context.Context, ins []*domain.CompanyTaskDetail) error {
	tasks := []*CompanyTaskDetail{}

	for _, in := range ins {
		task := &CompanyTaskDetail{
			Id:                           in.Id,
			CompanyTaskId:                in.CompanyTaskId,
			UserId:                       in.UserId,
			ClientKey:                    in.ClientKey,
			OpenId:                       in.OpenId,
			VideoId:                      in.VideoId,
			CompanyTaskAccountRelationId: in.CompanyTaskAccountRelationId,
			ItemId:                       in.ItemId,
			PlayCount:                    in.PlayCount,
			Cover:                        in.Cover,
			ReleaseTime:                  in.ReleaseTime,
			IsReleaseVideo:               in.IsReleaseVideo,
			IsPlaySuccess:                in.IsPlaySuccess,
			CreateTime:                   in.CreateTime,
			UpdateTime:                   in.UpdateTime,
		}

		tasks = append(tasks, task)
	}

	return ctr.data.DB(ctx).Save(&tasks).Error
}

func (ctr *companyTaskDetailRepo) GetByIdUsers(ctx context.Context, userId uint64) (*v1.GetByIdUsersReply, error) {
	return ctr.data.weixinuc.GetByIdUsers(ctx, &v1.GetByIdUsersRequest{
		UserId: userId,
	})
}

func (ctr *companyTaskDetailRepo) ListByClientKeyAndOpenIds(ctx context.Context, pageNum, pageSize uint64, clientKeyAndOpenIds, keyword string) (*v1.ListByClientKeyAndOpenIdsReply, error) {
	return ctr.data.weixinuc.ListByClientKeyAndOpenIds(ctx, &v1.ListByClientKeyAndOpenIdsRequest{
		PageNum:             pageNum,
		PageSize:            pageSize,
		ClientKeyAndOpenIds: clientKeyAndOpenIds,
		Keyword:             keyword,
	})
}

// DeleteOpenDouyinUsers
// 同时更新视频信息，可能视频已经删除或者状态更新

func (ctr *companyTaskDetailRepo) DeleteOpenDouyinUsers(ctx context.Context, userIds []uint64) error {
	db := ctr.data.db.WithContext(ctx).
		Where("id in (?)", userIds)

	if err := db.Delete(&CompanyTaskDetail{}).Error; err != nil {
		return err
	}

	return nil
}
