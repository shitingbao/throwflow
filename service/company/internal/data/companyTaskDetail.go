package data

import (
	v1 "company/api/service/weixin/v1"
	"company/internal/biz"
	"company/internal/domain"
	"company/internal/pkg/tool"
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

// 种草任务明细表
// 一个任务关系，对应该表多个视频明细
type CompanyTaskDetail struct {
	Id                           uint64                     `gorm:"column:id;primarykey;type:bigint(20) UNSIGNED;autoIncrement;not null;comment:自增ID"`
	CompanyTaskId                uint64                     `gorm:"column:company_task_id;type:bigint(20) UNSIGNED;not null;comment:任务ID"`
	UserId                       uint64                     `gorm:"column:user_id;type:bigint(20) UNSIGNED;not null;index:idx_user_id;comment:微信小程序用户ID"`
	ClientKey                    string                     `gorm:"column:client_key;type:varchar(50);not null;index:idx_client_key_open_id;comment:抖音开放平台应用Client Key"`
	OpenId                       string                     `gorm:"column:open_id;type:varchar(100);not null;index:idx_client_key_open_id;comment:抖音开放平台授权用户唯一标识"`
	CompanyTaskAccountRelationId uint64                     `gorm:"column:company_task_account_relation_id;type:bigint(20) UNSIGNED;not null;comment:任务关系ID"`
	ProductName                  string                     `gorm:"column:product_name;type:varchar(250);not null;comment:商品名称"`
	ItemId                       string                     `gorm:"column:item_id;type:varchar(250);not null;comment:视频id"` // 视频id // video_id
	PlayCount                    uint64                     `gorm:"column:play_count;type:int(10);not null;comment:播放数"`    // 播放数
	Cover                        string                     `gorm:"column:cover;type:text;not null;comment:视频封面"`           // 视频封面
	ReleaseTime                  time.Time                  `gorm:"column:release_time;type:datetime;not null;comment:发布时间"`
	IsReleaseVideo               uint8                      `gorm:"column:is_release_video;type:tinyint(3) UNSIGNED;not null;default:0;comment:视频发布,1:是,0:否"`
	IsPlaySuccess                uint8                      `gorm:"column:is_play_success;type:tinyint(3) UNSIGNED;not null;default:0;comment:播放量达标,1:是,0:否"`
	CreateTime                   time.Time                  `gorm:"column:create_time;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:新增时间"`
	UpdateTime                   time.Time                  `gorm:"column:update_time;type:datetime;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:修改时间"`
	CompanyTaskAccountRelation   CompanyTaskAccountRelation `gorm:"foreignKey:CompanyTaskAccountRelationId"`
}

func (CompanyTaskDetail) TableName() string {
	return "company_task_detail"
}

func (c *CompanyTaskDetail) ToDomain(ctx context.Context) *domain.CompanyTaskDetail {
	task := &domain.CompanyTaskDetail{
		Id:                           c.Id,
		CompanyTaskId:                c.CompanyTaskId,
		CompanyTaskAccountRelationId: c.CompanyTaskAccountRelationId,
		ProductName:                  c.ProductName,
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
		CompanyTaskAccountRelation: domain.CompanyTaskAccountRelation{
			Id:                    c.CompanyTaskAccountRelation.Id,
			CompanyTaskId:         c.CompanyTaskAccountRelation.CompanyTaskId,
			ProductOutId:          c.CompanyTaskAccountRelation.ProductOutId,
			ProductName:           c.CompanyTaskAccountRelation.ProductName,
			UserId:                c.CompanyTaskAccountRelation.UserId,
			ClaimTime:             c.CompanyTaskAccountRelation.ClaimTime,
			ExpireTime:            c.CompanyTaskAccountRelation.ExpireTime,
			Status:                c.CompanyTaskAccountRelation.Status,
			IsDel:                 c.CompanyTaskAccountRelation.IsDel,
			CreateTime:            c.CompanyTaskAccountRelation.CreateTime,
			UpdateTime:            c.CompanyTaskAccountRelation.UpdateTime,
			IsCostBuy:             c.CompanyTaskAccountRelation.IsCostBuy,
			ScreenshotAddress:     c.CompanyTaskAccountRelation.ScreenshotAddress,
			IsScreenshotAvailable: c.CompanyTaskAccountRelation.IsScreenshotAvailable,
		},
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

	db := ctr.data.db.WithContext(ctx).Table("company_task_detail").Preload("CompanyTaskAccountRelation")

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

func (ctr *companyTaskDetailRepo) Count(ctx context.Context, taskId uint64, clientKeyAndOpenIds []domain.CompanyTaskClientKeyAndOpenId) (int64, error) {
	db := ctr.data.db.WithContext(ctx).Table("company_task_detail")

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
		CompanyTaskAccountRelationId: in.CompanyTaskAccountRelationId,
		ProductName:                  in.ProductName,
		ItemId:                       in.ItemId,
		PlayCount:                    in.PlayCount,
		Cover:                        in.Cover,
		ReleaseTime:                  in.ReleaseTime,
		IsPlaySuccess:                in.IsPlaySuccess,
		// ScreenshotAddress:            in.ScreenshotAddress,
		// Status:                       in.Status,
		CreateTime: in.CreateTime,
		UpdateTime: in.UpdateTime,
	}

	if err := ctr.data.DB(ctx).Table("company_task_detail").Create(detail).Error; err != nil {
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
			CompanyTaskAccountRelationId: in.CompanyTaskAccountRelationId,
			ProductName:                  in.ProductName,
			ItemId:                       in.ItemId,
			PlayCount:                    in.PlayCount,
			Cover:                        in.Cover,
			ReleaseTime:                  in.ReleaseTime,
			IsPlaySuccess:                in.IsPlaySuccess,
			// ScreenshotAddress:            in.ScreenshotAddress,
			CreateTime: in.CreateTime,
			UpdateTime: in.UpdateTime,
		}
		details = append(details, detail)
	}

	return ctr.data.DB(ctx).Table("company_task_detail").Create(details).Error
}

func (ctr *companyTaskDetailRepo) Update(ctx context.Context, in *domain.CompanyTaskDetail) (*domain.CompanyTaskDetail, error) {
	task := &CompanyTaskDetail{
		Id:                           in.Id,
		CompanyTaskId:                in.CompanyTaskId,
		CompanyTaskAccountRelationId: in.CompanyTaskAccountRelationId,
		ProductName:                  in.ProductName,
		ItemId:                       in.ItemId,
		PlayCount:                    in.PlayCount,
		Cover:                        in.Cover,
		ReleaseTime:                  in.ReleaseTime,
		IsPlaySuccess:                in.IsPlaySuccess,
		// ScreenshotAddress:            in.ScreenshotAddress,
		CreateTime: in.CreateTime,
		UpdateTime: in.UpdateTime,
		UserId:     in.UserId,
		ClientKey:  in.ClientKey,
		OpenId:     in.OpenId,
		// IsCostBuy:             in.IsCostBuy,
		// IsScreenshotAvailable: in.IsScreenshotAvailable,
		// Status:                       in.Status,
	}

	if err := ctr.data.DB(ctx).Save(task).Error; err != nil {
		return nil, err
	}

	return task.ToDomain(ctx), nil
}

func (ctr *companyTaskDetailRepo) UpdateOnDuplicateKey(ctx context.Context, in []*domain.CompanyTaskDetail) error {
	values := []string{}

	for _, v := range in {
		val := "("
		val += strconv.FormatUint(v.Id, 10) + ","
		val += strconv.FormatUint(v.CompanyTaskId, 10) + ","
		val += strconv.FormatUint(v.UserId, 10) + ","
		val += "'" + v.ClientKey + "',"
		val += "'" + v.OpenId + "',"
		val += strconv.FormatUint(v.CompanyTaskAccountRelationId, 10) + ","
		val += "'" + v.ProductName + "',"
		val += "'" + v.ItemId + "',"
		val += strconv.FormatUint(v.PlayCount, 10) + ","
		val += "'" + v.Cover + "',"
		val += "'" + tool.TimeToString("2006-01-02 15:04", v.ReleaseTime) + "',"
		// val += strconv.Itoa(int(v.IsCostBuy)) + ","
		val += strconv.Itoa(int(v.IsReleaseVideo)) + ","
		val += strconv.Itoa(int(v.IsPlaySuccess)) + ","
		// val += "'" + v.ScreenshotAddress + "',"
		// val += strconv.Itoa(int(v.IsScreenshotAvailable)) + ","
		// val += strconv.Itoa(int(v.Status)) + ","
		val += "'" + tool.TimeToString("2006-01-02 15:04", v.CreateTime) + "',"
		val += "'" + tool.TimeToString("2006-01-02 15:04", v.UpdateTime) + "'"
		val += ")"
		values = append(values, val)
	}

	sql := `
	INSERT INTO company_task_detail (
		id, 
		company_task_id, 
		user_id, 
		client_key, 
		open_id, 
		company_task_account_relation_id, 
		product_name,
		item_id,
		play_count, 
		cover, 
		release_time, 
		is_release_video, 
		is_play_success, 
		create_time, 
		update_time
	) VALUES 
		 %s
	ON DUPLICATE KEY UPDATE
	play_count = VALUES(play_count),
	is_play_success = VALUES(is_play_success);
		`
	exec := fmt.Sprintf(sql, strings.Join(values, ","))

	return ctr.data.DB(ctx).Exec(exec).Error
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

// DeleteOpenDouyinUsers 删除没有对应关系的抖音账号
func (ctr *companyTaskDetailRepo) DeleteOpenDouyinUsers(ctx context.Context, userId uint64, clientKeys, openIds []string) error {
	sql := `delete from company_task_detail where user_id = ? and (client_key not in (?) or open_id not in (?))`

	if err := ctr.data.db.WithContext(ctx).Exec(sql, userId, clientKeys, openIds).Error; err != nil {
		return err
	}

	return nil
}
