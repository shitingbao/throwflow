package data

import (
	"context"
	"material/internal/biz"
	"material/internal/domain"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

// 素材文案表
// 每个微信用户对应商品的一条素材唯一
type MaterialContent struct {
	Id         uint64    `gorm:"column:id;primarykey;type:bigint(20) UNSIGNED;autoIncrement;not null;comment:自增ID"`
	ProductId  uint64    `gorm:"column:product_id;type:bigint(20) UNSIGNED;not null;comment:商品ID"`
	UserId     uint64    `gorm:"column:user_id;type:bigint(20) UNSIGNED;not null;uniqueIndex:idx_unique_product_id_user_id;;comment:微信小程序用户ID"`
	VideoId    uint64    `gorm:"column:video_id;type:bigint(20) UNSIGNED;not null;uniqueIndex:idx_unique_product_id_user_id;comment:素材视频ID"`
	Content    string    `gorm:"column:content;type:varchar(250);not null;comment:ai文案"`
	CreateTime time.Time `gorm:"column:create_time;type:datetime;not null;comment:新增时间"`
	UpdateTime time.Time `gorm:"column:update_time;type:datetime;not null;comment:修改时间"`
}

func (MaterialContent) TableName() string {
	return "material_content"
}

type materialContentRepo struct {
	data *Data
	log  *log.Helper
}

func NewMaterialContentRepo(data *Data, logger log.Logger) biz.MaterialContentRepo {
	return &materialContentRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (c *MaterialContent) ToDomain(ctx context.Context) *domain.MaterialContent {
	task := &domain.MaterialContent{
		Id:         c.Id,
		ProductId:  c.ProductId,
		UserId:     c.UserId,
		Content:    c.Content,
		VideoId:    c.VideoId,
		CreateTime: c.CreateTime,
		UpdateTime: c.UpdateTime,
	}

	return task
}

func (mc *materialContentRepo) Get(ctx context.Context, userId, videoId uint64) (*domain.MaterialContent, error) {
	content := &MaterialContent{}

	if err := mc.data.db.WithContext(ctx).
		Where("user_id = ? and video_id = ?", userId, videoId).
		First(content).Error; err != nil {
		return nil, err
	}

	return content.ToDomain(ctx), nil
}

func (mc *materialContentRepo) List(ctx context.Context, pageNum, pageSize int, userId uint64) ([]*domain.MaterialContent, error) {
	list := []*MaterialContent{}
	contentList := []*domain.MaterialContent{}

	db := mc.data.db.WithContext(ctx).Where("user_id = ?", userId)

	if err := db.Order("create_time desc").
		Limit(pageSize).Offset((pageNum - 1) * pageSize).
		Find(list).Error; err != nil {
		return nil, err
	}

	for _, content := range list {
		contentList = append(contentList, content.ToDomain(ctx))
	}

	return contentList, nil
}

func (mc *materialContentRepo) Count(ctx context.Context, userId uint64) (int64, error) {
	var count int64

	if err := mc.data.db.WithContext(ctx).
		Where("user_id = ?", userId).
		Count(&count).Error; err != nil {
		return count, err
	}

	return count, nil
}

func (ctr *materialContentRepo) Save(ctx context.Context, in *domain.MaterialContent) (*domain.MaterialContent, error) {
	detail := &MaterialContent{
		Id:         in.Id,
		ProductId:  in.ProductId,
		UserId:     in.UserId,
		Content:    in.Content,
		CreateTime: in.CreateTime,
		UpdateTime: in.UpdateTime,
	}

	if err := ctr.data.db.WithContext(ctx).Model(&MaterialContent{}).Create(detail).Error; err != nil {
		return nil, err
	}

	return detail.ToDomain(ctx), nil
}

func (ctr *materialContentRepo) Update(ctx context.Context, in *domain.MaterialContent) (*domain.MaterialContent, error) {
	task := &MaterialContent{
		Id:         in.Id,
		ProductId:  in.ProductId,
		UserId:     in.UserId,
		Content:    in.Content,
		CreateTime: in.CreateTime,
		UpdateTime: in.UpdateTime,
	}

	if err := ctr.data.db.WithContext(ctx).Save(task).Error; err != nil {
		return nil, err
	}

	return task.ToDomain(ctx), nil
}
