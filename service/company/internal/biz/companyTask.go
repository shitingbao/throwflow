package biz

import (
	"company/internal/conf"
	"company/internal/domain"
	"company/internal/pkg/tool"
	"context"
	"strconv"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

var (
	CompanyTaskInfoError          = errors.NotFound("COMPANY_TASK_INFO_ERROR", "商品任务详情出错")
	CompanyTaskNotFound           = errors.NotFound("COMPANY_TASK_NOT_FOUND", "商品任务不存在")
	CompanyTaskListError          = errors.NotFound("COMPANY_TASK_LIST_ERROR", "商品任务列表有误")
	CompanyTaskProductCreate      = errors.InternalServer("COMPANY_TASK_PRODUCT_CREATE", "商品任务创建失败")
	CompanyTaskProductExists      = errors.InternalServer("COMPANY_TASK_PRODUCT_EXISTS", "该商品已存在任务")
	CompanyTaskProductDelete      = errors.InternalServer("COMPANY_TASK_PRODUCT_DELETE", "删除失败")
	CompanyTaskProductUpdateTop   = errors.InternalServer("COMPANY_TASK_PRODUCT_UPDATE_TOP", "更新置顶失败")
	CompanyTaskProductUpdateQuota = errors.InternalServer("COMPANY_TASK_PRODUCT_UPDATE_QUOTA", "更新任务数量失败")
	CompanyTaskReleaseTimeError   = errors.InternalServer("COMPANY_TASK_RELEASE_TIME_ERROR", "商品任务发布时间出错")
)

type CompanyTaskRepo interface {
	GetById(context.Context, uint64) (*domain.CompanyTask, error)
	GetByProductOutId(context.Context, uint64, uint32, string) (*domain.CompanyTask, error)
	List(context.Context, int, int, int, int, []uint64, string) ([]*domain.CompanyTask, error)
	ListByProductOutId(context.Context, []string) ([]*domain.CompanyTask, error)
	ListByIds(context.Context, []uint64) ([]*domain.CompanyTask, error)
	Count(context.Context, int, []uint64, string) (int64, error)
	Save(context.Context, *domain.CompanyTask) (*domain.CompanyTask, error)
	Update(context.Context, *domain.CompanyTask) (*domain.CompanyTask, error)
	UpdateCompanyTaskIsDel(context.Context, uint64) error

	GetCacheHash(context.Context, string) (string, error)
	SaveCacheHash(context.Context, string, uint64) error
	DeleteCacheHash(context.Context, string) error
	UpdateCacheHash(context.Context, string, int64) error
}

type CompanyTaskUsecase struct {
	repo   CompanyTaskRepo
	ctarpo CompanyTaskAccountRelationRepo
	cprepo CompanyProductRepo
	tm     Transaction
	conf   *conf.Data
	log    *log.Helper
}

func NewCompanyTaskUsecase(repo CompanyTaskRepo, ctarpo CompanyTaskAccountRelationRepo, cprepo CompanyProductRepo, tm Transaction, conf *conf.Data, logger log.Logger) *CompanyTaskUsecase {
	return &CompanyTaskUsecase{repo: repo, ctarpo: ctarpo, cprepo: cprepo, tm: tm, conf: conf, log: log.NewHelper(logger)}
}

func (c *CompanyTaskUsecase) GetByProductOutId(ctx context.Context, productOutId, userId uint64) (*domain.CompanyTask, error) {
	task, err := c.repo.GetByProductOutId(ctx, productOutId, 0, tool.TimeToString("2006-01-02 15:04:05", time.Now()))

	if err != nil {
		return nil, CompanyTaskInfoError
	}

	claimQuota, err := c.repo.GetCacheHash(ctx, strconv.FormatUint(task.Id, 10))

	if err == nil {
		num, err := strconv.ParseUint(claimQuota, 10, 64)

		if err == nil {
			task.SetClaimQuota(ctx, task.Quota-num)
		}
	}

	_, err = c.ctarpo.GetByProductOutIdAndUserId(ctx, productOutId, userId, tool.TimeToString("2006-01-02 15:04:05", time.Now()))

	if err == nil {
		task.SetIsUserExist(ctx)
	}

	return task, nil
}

func (c *CompanyTaskUsecase) ListCompanyTask(ctx context.Context, isDel int, isTop uint32, pageNum, pageSize uint64, keyword, releaseTime string) (*domain.CompanyTaskList, error) {
	productIds := []uint64{}

	if len(keyword) > 0 {
		products, err := c.cprepo.List(ctx, 0, 40, 0, 0, 0, "", "1", keyword)

		if err != nil {
			return nil, CompanyTaskListError
		}

		for _, v := range products {
			productIds = append(productIds, v.ProductOutId)
		}

		if len(products) == 0 {
			return &domain.CompanyTaskList{
				PageNum:  pageNum,
				PageSize: pageSize,
				Total:    0,
				List:     []*domain.CompanyTask{},
			}, nil
		}
	}

	tasks, err := c.repo.List(ctx, int(pageNum), int(pageSize), int(isTop), isDel, productIds, releaseTime)

	if err != nil {
		return nil, CompanyTaskListError
	}

	productMap, err := c.getExistProductsMapByTasks(ctx, tasks)

	if err != nil {
		return nil, CompanyTaskListError
	}

	total, err := c.repo.Count(ctx, isDel, productIds, releaseTime)

	if err != nil {
		return nil, CompanyTaskListError
	}

	for _, t := range tasks {
		claimQuota, err := c.repo.GetCacheHash(ctx, strconv.FormatUint(t.Id, 10))

		if err == nil {
			num, err := strconv.ParseUint(claimQuota, 10, 64)

			if err == nil {
				t.SetClaimQuota(ctx, t.Quota-num)
			}
		}

		product := productMap[t.ProductOutId]

		if product != nil {
			product.SetProductImgs(ctx)

			t.CompanyProduct = *product
		}
	}

	return &domain.CompanyTaskList{
		PageNum:  pageNum,
		PageSize: pageSize,
		Total:    uint64(total),
		List:     tasks,
	}, nil
}

// CreateCompanyTask
// create task and set use‘s number in redis
func (c *CompanyTaskUsecase) CreateCompanyTask(ctx context.Context, productOutId, expireTime, playNum, quota uint64, isGoodReviews uint8, price float64, releaseTime string) (*domain.CompanyTask, error) {
	existTask, err := c.repo.GetByProductOutId(ctx, productOutId, 0, "")

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, CompanyTaskProductCreate
	}

	if existTask != nil {
		return nil, CompanyTaskProductExists
	}

	tk := domain.NewCompanyTask(ctx, productOutId, expireTime, playNum, quota, isGoodReviews, price)

	reTime := time.Now()

	if len(releaseTime) > 0 {
		tm, err := tool.StringToTime("2006-01-02 15:04:05", releaseTime)

		if err != nil {
			return nil, CompanyTaskReleaseTimeError
		}

		if tm.Before(time.Now()) {
			return nil, CompanyTaskReleaseTimeError
		}

		reTime = tm
	}

	tk.SetReleaseTime(ctx, reTime)
	tk.SetCreateTime(ctx)
	tk.SetUpdateTime(ctx)

	task, err := c.repo.Save(ctx, tk)

	if err != nil {
		return nil, CompanyTaskProductCreate
	}

	if err := c.repo.SaveCacheHash(ctx, strconv.FormatUint(task.Id, 10), quota); err != nil {
		return nil, CompanyTaskProductCreate
	}

	return task, nil
}

// UpdateCompanyTaskQuota 更新任务名额，需要把 redis 中可用数量更新
// 注意先更新 mysql
func (c *CompanyTaskUsecase) UpdateCompanyTaskQuota(ctx context.Context, taskId, quota uint64) (*domain.CompanyTask, error) {
	tk, err := c.repo.GetById(ctx, taskId)

	if err != nil {
		return nil, CompanyTaskProductUpdateQuota
	}

	tk.SetQuota(ctx, quota)
	tk.SetUpdateTime(ctx)

	task, err := c.repo.Update(ctx, tk)

	if err != nil {
		return nil, CompanyTaskProductUpdateQuota
	}

	available, err := c.ctarpo.CountAvailableByTaskId(ctx, taskId)

	if err != nil {
		return nil, CompanyTaskProductUpdateQuota
	}

	err = c.repo.SaveCacheHash(ctx, strconv.FormatUint(taskId, 10), (quota - uint64(available)))

	if err != nil {
		return nil, CompanyTaskProductUpdateQuota
	}

	return task, nil
}

func (c *CompanyTaskUsecase) UpdateCompanyTaskIsTop(ctx context.Context, taskId uint64, isTop uint8) (*domain.CompanyTask, error) {
	tk, err := c.repo.GetById(ctx, taskId)

	if err != nil {
		return nil, CompanyTaskProductUpdateTop
	}

	tk.SetIsTop(ctx, isTop)
	tk.SetUpdateTime(ctx)

	task, err := c.repo.Update(ctx, tk)

	if err != nil {
		return nil, CompanyTaskProductUpdateTop
	}

	return task, nil
}

// UpdateCompanyTaskIsDel 删除任务后
func (c *CompanyTaskUsecase) UpdateCompanyTaskIsDel(ctx context.Context, taskId uint64) error {
	err := c.tm.InTx(ctx, func(ctx context.Context) error {

		if err := c.repo.DeleteCacheHash(ctx, strconv.FormatUint(taskId, 10)); err != nil {
			return CompanyTaskProductDelete
		}

		task, err := c.repo.GetById(ctx, taskId)

		if err != nil {
			return CompanyTaskProductDelete
		}

		task.SetIsTop(ctx, 0)
		task.SetIsDel(ctx)

		if _, err := c.repo.Update(ctx, task); err != nil {
			return CompanyTaskProductDelete
		}

		return nil
	})

	return err
}

func (c *CompanyTaskUsecase) getExistProductsMapByTasks(ctx context.Context, tasks []*domain.CompanyTask) (map[uint64]*domain.CompanyProduct, error) {
	productOutIds := []uint64{}
	products := make(map[uint64]*domain.CompanyProduct)

	for _, v := range tasks {
		productOutIds = append(productOutIds, v.ProductOutId)
	}

	companyProducts, err := c.cprepo.ListByProductOutIds(ctx, "1", productOutIds)

	if err != nil {
		return nil, CompanyTaskListError
	}

	for _, companyProduct := range companyProducts {
		companyProduct.SetCommissions(ctx)

		pureCommission, pureServiceCommission, _ := companyProduct.GetCommission(ctx)

		companyProduct.SetPureCommission(ctx, pureCommission)
		companyProduct.SetPureServiceCommission(ctx, pureServiceCommission)

		products[companyProduct.ProductOutId] = companyProduct
	}

	return products, nil
}
