package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
	v1 "interface/api/service/company/v1"
	"interface/internal/biz"
)

type clueRepo struct {
	data *Data
	log  *log.Helper
}

func NewClueRepo(data *Data, logger log.Logger) biz.ClueRepo {
	return &clueRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (cr *clueRepo) ListSelectClues(ctx context.Context) (*v1.ListSelectCluesReply, error) {
	list, err := cr.data.companyuc.ListSelectClues(ctx, &emptypb.Empty{})

	if err != nil {
		return nil, err
	}

	return list, err
}

func (cr *clueRepo) Save(ctx context.Context, companyName, contactInformation, source string, companyType, status uint32, areaCode uint64) (*v1.CreateCluesReply, error) {
	applyForm, err := cr.data.companyuc.CreateClues(ctx, &v1.CreateCluesRequest{
		CompanyName:        companyName,
		IndustryId:         "",
		ContactInformation: contactInformation,
		CompanyType:        companyType,
		AreaCode:           areaCode,
		Source:             source,
		Status:             status,
	})

	if err != nil {
		return nil, err
	}

	return applyForm, err
}

func (cr *clueRepo) UpdateCompanyName(ctx context.Context, companyId uint64, companyName string) (*v1.UpdateCompanyNameCluesReply, error) {
	list, err := cr.data.companyuc.UpdateCompanyNameClues(ctx, &v1.UpdateCompanyNameCluesRequest{
		CompanyId:   companyId,
		CompanyName: companyName,
	})

	if err != nil {
		return nil, err
	}

	return list, err
}
