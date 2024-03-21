package biz

import (
	"company/internal/conf"
	"company/internal/domain"
	"context"
	"strconv"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
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
)

type CompanyTaskRepo interface {
	GetById(context.Context, uint64) (*domain.CompanyTask, error)
	GetByProductOutId(context.Context, uint64) (*domain.CompanyTask, error)
	List(context.Context, int, int, int, int, []uint64) ([]*domain.CompanyTask, error)
	ListIds(context.Context) ([]int, error)
	Save(context.Context, *domain.CompanyTask) (*domain.CompanyTask, error)
	Update(context.Context, *domain.CompanyTask) (*domain.CompanyTask, error)
	UpdateCompanyTaskIsDel(context.Context, uint64) error
	Count(context.Context, int, []uint64) (int64, error)
	CountByProjectOutId(context.Context, uint64) (int64, error)
	GetCacheHash(context.Context, string) (string, error)
	SaveCacheHash(context.Context, string, uint64) error
	DeleteCache(context.Context, string) error
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

func (c *CompanyTaskUsecase) GetByProductOutId(ctx context.Context, productOutId, userId uint64) (int64, *domain.CompanyTask, error) {
	task, err := c.repo.GetByProductOutId(ctx, productOutId)

	if err != nil {
		return 0, nil, CompanyTaskInfoError
	}

	count, err := c.ctarpo.CountByProductOutId(ctx, productOutId, userId)

	if err != nil {
		return 0, nil, CompanyTaskInfoError
	}

	return count, task, nil
}

func (c *CompanyTaskUsecase) ListCompanyTask(ctx context.Context, isDel int, isTop uint32, pageNum, pageSize uint64, keyword string) (*domain.CompanyTaskList, error) {
	productIds := []uint64{}

	if len(keyword) > 0 {
		products, err := c.cprepo.List(ctx, 0, 40, 0, 0, 0, "", "", keyword)

		if err != nil {
			return nil, CompanyTaskListError
		}

		for _, v := range products {
			productIds = append(productIds, v.ProductOutId)
		}
	}

	tasks, err := c.repo.List(ctx, int(pageNum), int(pageSize), int(isTop), isDel, productIds)

	if err != nil {
		return nil, CompanyTaskListError
	}

	productMap, err := c.getExistProductsMapByTasks(ctx, tasks)

	if err != nil {
		return nil, CompanyTaskListError
	}

	total, err := c.repo.Count(ctx, isDel, productIds)

	if err != nil {
		return nil, CompanyTaskListError
	}

	for _, t := range tasks {
		claimQuota, err := c.repo.GetCacheHash(ctx, strconv.FormatUint(t.Id, 10))

		if err == nil {
			num, err := strconv.ParseUint(claimQuota, 10, 64)

			if err == nil {
				t.SetClaimQuota(ctx, num)
			}
		}

		product := productMap[t.ProductOutId]

		if product != nil {
			product.SetProductImgs(ctx)

			if len(product.ProductImgs) > 0 {
				product.SetProductDetailImg(ctx, product.ProductImgs[0])
			}

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
func (c *CompanyTaskUsecase) CreateCompanyTask(ctx context.Context, productOutId, expireTime, playNum, price, quota uint64, isGoodReviews uint8) (*domain.CompanyTask, error) {
	count, err := c.repo.CountByProjectOutId(ctx, productOutId)

	if err != nil {
		return nil, CompanyTaskProductCreate
	}

	if count > 0 {
		return nil, CompanyTaskProductExists
	}

	tk := domain.NewCompanyTask(ctx, productOutId, expireTime, playNum, price, quota, isGoodReviews)
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

		if err := c.repo.DeleteCache(ctx, strconv.FormatUint(taskId, 10)); err != nil {
			return CompanyTaskProductDelete
		}

		if err := c.repo.UpdateCompanyTaskIsDel(ctx, taskId); err != nil {
			return CompanyTaskProductDelete
		}

		// if err := c.ctarpo.SoftDelete(ctx, taskId); err != nil {
		// // 关闭任务后，关闭对应达人关系（已领取的任务，待定）
		// 	return CompanyTaskProductDelete
		// }
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
