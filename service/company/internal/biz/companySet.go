package biz

import (
	"company/internal/conf"
	"company/internal/domain"
	"company/internal/pkg/tool"
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"strconv"
	"time"
)

var (
	CompanyCompanySetNotFound    = errors.NotFound("COMPANY_COMPANY_SET_NOT_FOUND", "公司设置不存在")
	CompanyCompanySetCreateError = errors.InternalServer("COMPANY_COMPANY_SET_CREATE_ERROR", "公司设置创建失败")
	CompanyCompanySetUpdateError = errors.InternalServer("COMPANY_COMPANY_SET_UPDATE_ERROR", "公司设置失败")
)

type CompanySetRepo interface {
	GetByCompanyIdAndDayAndSetKey(context.Context, uint64, uint32, string) (*domain.CompanySet, error)
	Save(context.Context, *domain.CompanySet) (*domain.CompanySet, error)
	Update(context.Context, *domain.CompanySet) error
}

type CompanySetUsecase struct {
	repo  CompanySetRepo
	crepo CompanyRepo
	conf  *conf.Data
	log   *log.Helper
}

func NewCompanySetUsecase(repo CompanySetRepo, crepo CompanyRepo, conf *conf.Data, logger log.Logger) *CompanySetUsecase {
	return &CompanySetUsecase{repo: repo, crepo: crepo, conf: conf, log: log.NewHelper(logger)}
}

func (csuc *CompanySetUsecase) GetCompanySets(ctx context.Context, companyId uint64, day, setKey string) (*domain.CompanySet, error) {
	dayTime, _ := tool.StringToTime("2006-01-02", day)
	uday, _ := strconv.ParseUint(tool.TimeToString("20060102", dayTime), 10, 64)

	if _, err := csuc.crepo.GetById(ctx, companyId); err != nil {
		return nil, CompanyCompanyNotFound
	}

	companySet, err := csuc.repo.GetByCompanyIdAndDayAndSetKey(ctx, companyId, uint32(uday), setKey)

	if err != nil {
		if setKey == "sampleThreshold" {
			return &domain.CompanySet{
				CompanyId: companyId,
				Day:       uint32(uday),
				SetKey:    "sampleThreshold",
				SetValue:  "{\"type\":2,\"value\":10000}",
			}, nil
		}
	}

	return companySet, nil
}

func (csuc *CompanySetUsecase) UpdateCompanySets(ctx context.Context, companyId uint64, setKey, setValue string) error {
	day, _ := strconv.ParseUint(tool.TimeToString("20060102", time.Now()), 10, 64)

	if _, err := csuc.crepo.GetById(ctx, companyId); err != nil {
		return CompanyCompanyNotFound
	}

	if inCompanySet, err := csuc.repo.GetByCompanyIdAndDayAndSetKey(ctx, companyId, uint32(day), setKey); err != nil {
		inCompanySet = domain.NewCompanySet(ctx, companyId, uint32(day), setKey, setValue)
		inCompanySet.SetCreateTime(ctx)
		inCompanySet.SetUpdateTime(ctx)

		if ok := inCompanySet.VerifySetValue(ctx); !ok {
			return CompanyCompanySetUpdateError
		}

		if _, err := csuc.repo.Save(ctx, inCompanySet); err != nil {
			return CompanyCompanySetUpdateError
		}
	} else {
		if inCompanySet.Day == uint32(day) {
			inCompanySet.SetSetValue(ctx, setValue)
			inCompanySet.SetUpdateTime(ctx)

			if ok := inCompanySet.VerifySetValue(ctx); !ok {
				return CompanyCompanySetUpdateError
			}

			if err := csuc.repo.Update(ctx, inCompanySet); err != nil {
				return CompanyCompanySetUpdateError
			}
		} else {
			inCompanySet = domain.NewCompanySet(ctx, companyId, uint32(day), setKey, setValue)
			inCompanySet.SetCreateTime(ctx)
			inCompanySet.SetUpdateTime(ctx)

			if ok := inCompanySet.VerifySetValue(ctx); !ok {
				return CompanyCompanySetUpdateError
			}

			if _, err := csuc.repo.Save(ctx, inCompanySet); err != nil {
				return CompanyCompanySetUpdateError
			}
		}
	}

	return nil
}
