package data

import (
	"company/internal/biz"
	"company/internal/domain"
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

var (
	// RedisTaskPre + taskId
	RedisCompanyTaskPre = "company:task:"
)

// 种草任务表
// product_out_id 对应使用的商品唯一标识
type CompanyTask struct {
	Id            uint64    `gorm:"column:id;primarykey;type:bigint(20) UNSIGNED;autoIncrement;not null;comment:自增ID"`
	ProductOutId  uint64    `gorm:"column:product_out_id;type:bigint(20) UNSIGNED;not null;comment:商品外ID"`
	ExpireTime    uint64    `gorm:"column:expire_time;type:int(10);not null;comment:过期时间(天)"`
	PlayNum       uint64    `gorm:"column:play_num;type:int(10) UNSIGNED;not null;default:0;comment:播放数量"`
	Price         float64   `gorm:"column:price;type:decimal(10, 2) UNSIGNED;not null;default:0;comment:每条价格"`
	Quota         uint64    `gorm:"column:quota;type:int(3) UNSIGNED;not null;default:0;comment:数量"`
	ClaimQuota    uint64    `gorm:"column:claim_quota;type:int(3) UNSIGNED;not null;default:0;comment:领取数量"`
	SuccessQuota  uint64    `gorm:"column:success_quota;type:int(3) UNSIGNED;not null;default:0;comment:成功数量"`
	IsTop         uint8     `gorm:"column:is_top;type:tinyint(3) UNSIGNED;not null;default:0;comment:置顶:1:是,0:否"`
	IsDel         uint8     `gorm:"column:is_del;type:tinyint(3) UNSIGNED;not null;default:0;comment:是否移除,1:已移除,0:未移除"`
	IsGoodReviews uint8     `gorm:"column:is_good_reviews;type:tinyint(3) UNSIGNED;not null;default:0;comment:是否需要好评,1:需要,0:不需要"`
	CreateTime    time.Time `gorm:"column:create_time;type:datetime;not null;comment:新增时间"`
	UpdateTime    time.Time `gorm:"column:update_time;type:datetime;not null;comment:修改时间"`
}

func (CompanyTask) TableName() string {
	return "company_task"
}

func (c *CompanyTask) ToDomain(ctx context.Context) *domain.CompanyTask {
	task := &domain.CompanyTask{
		Id:            c.Id,
		ProductOutId:  c.ProductOutId,
		ExpireTime:    c.ExpireTime,
		PlayNum:       c.PlayNum,
		Price:         c.Price,
		Quota:         c.Quota,
		ClaimQuota:    c.ClaimQuota,
		SuccessQuota:  c.SuccessQuota,
		IsTop:         c.IsTop,
		IsDel:         c.IsDel,
		IsGoodReviews: c.IsGoodReviews,
		CreateTime:    c.CreateTime,
		UpdateTime:    c.UpdateTime,
	}
	return task
}

type companyTaskRepo struct {
	data *Data
	log  *log.Helper
}

func NewCompanyTaskRepo(data *Data, logger log.Logger) biz.CompanyTaskRepo {
	return &companyTaskRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (ctr *companyTaskRepo) GetById(ctx context.Context, id uint64) (*domain.CompanyTask, error) {
	task := &CompanyTask{}

	if err := ctr.data.db.WithContext(ctx).
		Where("is_del = ?", 0).First(task, id).Error; err != nil {
		return nil, err
	}

	return task.ToDomain(ctx), nil
}

func (ctr *companyTaskRepo) GetByProductOutId(ctx context.Context, productOutId uint64, isDel uint32) (*domain.CompanyTask, error) {
	task := &CompanyTask{}

	if err := ctr.data.db.WithContext(ctx).Model(&CompanyTask{}).
		Where("product_out_id = ? and is_del = ?", productOutId, isDel).First(task).Error; err != nil {
		return nil, err
	}

	return task.ToDomain(ctx), nil
}

// 排序根据置顶和操作时间
// isDel >= -1
func (ctr *companyTaskRepo) List(ctx context.Context, pageNum, pageSize, isTop, isDel int, productOutIds []uint64) ([]*domain.CompanyTask, error) {
	list := []*domain.CompanyTask{}
	tasks := []CompanyTask{}

	db := ctr.data.db.WithContext(ctx).Model(&CompanyTask{})

	if isDel >= 0 {
		db = db.Where("is_del = ?", isDel)
	}

	if len(productOutIds) > 0 {
		db = db.Where("product_out_id in (?)", productOutIds)
	}

	if isTop > 0 {
		db = db.Order("is_top DESC")
	}

	if err := db.Order("create_time DESC").
		Limit(pageSize).Offset((pageNum - 1) * pageSize).
		Find(&tasks).Error; err != nil {
		return nil, err
	}

	for _, t := range tasks {
		list = append(list, t.ToDomain(ctx))
	}

	return list, nil
}

func (ctr *companyTaskRepo) ListByProductOutId(ctx context.Context, productOutIds []string) ([]*domain.CompanyTask, error) {
	var companyTasks []CompanyTask
	list := make([]*domain.CompanyTask, 0)

	if result := ctr.data.db.WithContext(ctx).
		Where("product_out_id in ?", productOutIds).
		Where("is_del = 0").
		Order("create_time DESC").
		Find(&companyTasks); result.Error != nil {
		return nil, result.Error
	}

	for _, companyTask := range companyTasks {
		list = append(list, companyTask.ToDomain(ctx))
	}

	return list, nil
}

func (ctr *companyTaskRepo) ListByIds(ctx context.Context, ids []uint64) ([]*domain.CompanyTask, error) {
	var companyTasks []CompanyTask
	list := make([]*domain.CompanyTask, 0)

	if result := ctr.data.db.WithContext(ctx).
		Where("id in (?)", ids).
		Find(&companyTasks); result.Error != nil {
		return nil, result.Error
	}

	for _, companyTask := range companyTasks {
		list = append(list, companyTask.ToDomain(ctx))
	}

	return list, nil
}

func (ctr *companyTaskRepo) Count(ctx context.Context, isDel int, productOutIds []uint64) (int64, error) {
	db := ctr.data.db.WithContext(ctx).Model(&CompanyTask{})

	if isDel >= 0 {
		db = db.Where(" is_del = ?", 0)
	}

	if len(productOutIds) > 0 {
		db = db.Where("product_out_id in (?)", productOutIds)
	}

	var count int64

	if err := db.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

// Create a new task and add the number of available tasks to Redis
// And no timeout,when close the task to delete it
func (ctr *companyTaskRepo) Save(ctx context.Context, in *domain.CompanyTask) (*domain.CompanyTask, error) {
	task := &CompanyTask{
		ProductOutId:  in.ProductOutId,
		ExpireTime:    in.ExpireTime,
		PlayNum:       in.PlayNum,
		Price:         in.Price,
		Quota:         in.Quota,
		IsTop:         in.IsTop,
		IsDel:         in.IsDel,
		IsGoodReviews: in.IsGoodReviews,
		CreateTime:    in.CreateTime,
		UpdateTime:    in.UpdateTime,
	}

	if err := ctr.data.DB(ctx).Model(&CompanyTask{}).Create(task).Error; err != nil {
		return nil, err
	}

	return task.ToDomain(ctx), nil
}

func (ctr *companyTaskRepo) Update(ctx context.Context, in *domain.CompanyTask) (*domain.CompanyTask, error) {
	task := &CompanyTask{
		Id:            in.Id,
		ProductOutId:  in.ProductOutId,
		ExpireTime:    in.ExpireTime,
		PlayNum:       in.PlayNum,
		Price:         in.Price,
		Quota:         in.Quota,
		ClaimQuota:    in.ClaimQuota,
		SuccessQuota:  in.SuccessQuota,
		IsTop:         in.IsTop,
		IsDel:         in.IsDel,
		IsGoodReviews: in.IsGoodReviews,
		CreateTime:    in.CreateTime,
		UpdateTime:    in.UpdateTime,
	}

	if err := ctr.data.DB(ctx).Save(task).Error; err != nil {
		return nil, err
	}

	return task.ToDomain(ctx), nil
}

func (ctr *companyTaskRepo) UpdateCompanyTaskIsDel(ctx context.Context, id uint64) error {
	return ctr.data.DB(ctx).Model(&CompanyTask{}).Where("id = ?", id).Update("is_del", 1).Error
}

func (ctr *companyTaskRepo) GetCacheHash(ctx context.Context, taskId string) (string, error) {
	return ctr.data.rdb.Get(ctx, RedisCompanyTaskPre+taskId).Result()
}

func (ctr *companyTaskRepo) SaveCacheHash(ctx context.Context, taskId string, useCount uint64) error {
	return ctr.data.rdb.Set(ctx, RedisCompanyTaskPre+taskId, useCount, 0).Err()
}

func (ctr *companyTaskRepo) DeleteCacheHash(ctx context.Context, taskId string) error {
	return ctr.data.rdb.Del(ctx, RedisCompanyTaskPre+taskId).Err()
}

func (ctr *companyTaskRepo) UpdateCacheHash(ctx context.Context, taskId string, count int64) error {
	return ctr.data.rdb.IncrBy(ctx, RedisCompanyTaskPre+taskId, count).Err()
}
