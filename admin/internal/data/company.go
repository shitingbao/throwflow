package data

import (
	v1 "admin/api/service/company/v1"
	"admin/internal/biz"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
)

type companyRepo struct {
	data *Data
	log  *log.Helper
}

func NewCompanyRepo(data *Data, logger log.Logger) biz.CompanyRepo {
	return &companyRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (cr *companyRepo) List(ctx context.Context, pageNum, industryId uint64, keyword, status string, companyType uint32) (*v1.ListCompanysReply, error) {
	list, err := cr.data.companyuc.ListCompanys(ctx, &v1.ListCompanysRequest{
		PageNum:     pageNum,
		PageSize:    uint64(cr.data.conf.Database.PageSize),
		Keyword:     keyword,
		IndustryId:  industryId,
		Status:      status,
		CompanyType: companyType,
	})

	if err != nil {
		return nil, err
	}

	return list, err
}

func (cr *companyRepo) ListSelect(ctx context.Context) (*v1.ListSelectCompanysReply, error) {
	list, err := cr.data.companyuc.ListSelectCompanys(ctx, &emptypb.Empty{})

	if err != nil {
		return nil, err
	}

	return list, err
}

func (cr *companyRepo) Statistics(ctx context.Context) (*v1.StatisticsCompanysReply, error) {
	list, err := cr.data.companyuc.StatisticsCompanys(ctx, &emptypb.Empty{})

	if err != nil {
		return nil, err
	}

	return list, err
}

func (cr *companyRepo) Save(ctx context.Context, companyName, contactInformation, source, seller, facilitator, adminName, adminPhone, address, industryId string, userId, clueId, areaCode uint64, companyType, qianchuanUse, status uint32) (*v1.CreateCompanysReply, error) {
	company, err := cr.data.companyuc.CreateCompanys(ctx, &v1.CreateCompanysRequest{
		CompanyName:        companyName,
		IndustryId:         industryId,
		UserId:             userId,
		ContactInformation: contactInformation,
		CompanyType:        companyType,
		QianchuanUse:       qianchuanUse,
		Source:             source,
		Seller:             seller,
		Facilitator:        facilitator,
		Status:             status,
		ClueId:             clueId,
		AdminName:          adminName,
		AdminPhone:         adminPhone,
		Address:            address,
		AreaCode:           areaCode,
	})

	if err != nil {
		return nil, err
	}

	return company, err
}

func (cr *companyRepo) Update(ctx context.Context, id uint64, companyName, contactInformation, seller, facilitator, adminName, adminPhone, address, industryId string, clueId, areaCode uint64, companyType, qianchuanUse, status uint32) (*v1.UpdateCompanysReply, error) {
	company, err := cr.data.companyuc.UpdateCompanys(ctx, &v1.UpdateCompanysRequest{
		Id:                 id,
		CompanyName:        companyName,
		IndustryId:         industryId,
		ContactInformation: contactInformation,
		CompanyType:        companyType,
		QianchuanUse:       qianchuanUse,
		Seller:             seller,
		Facilitator:        facilitator,
		Status:             status,
		ClueId:             clueId,
		AdminName:          adminName,
		AdminPhone:         adminPhone,
		Address:            address,
		AreaCode:           areaCode,
	})

	if err != nil {
		return nil, err
	}

	return company, err
}

func (cr *companyRepo) UpdateStatus(ctx context.Context, id uint64, status uint32) (*v1.UpdateStatusCompanysReply, error) {
	company, err := cr.data.companyuc.UpdateStatusCompanys(ctx, &v1.UpdateStatusCompanysRequest{
		Id:     id,
		Status: status,
	})

	if err != nil {
		return nil, err
	}

	return company, err
}

func (cr *companyRepo) UpdateRole(ctx context.Context, id, userId uint64, menuIds, startTime, endTime string, accounts, qianchuanAdvertisers, companyType, isTermwork uint32) (*v1.UpdateRoleCompanysReply, error) {
	company, err := cr.data.companyuc.UpdateRoleCompanys(ctx, &v1.UpdateRoleCompanysRequest{
		Id:                   id,
		UserId:               userId,
		CompanyType:          companyType,
		MenuIds:              menuIds,
		Accounts:             accounts,
		QianchuanAdvertisers: qianchuanAdvertisers,
		IsTermwork:           isTermwork,
		StartTime:            startTime,
		EndTime:              endTime,
	})

	if err != nil {
		return nil, err
	}

	return company, err
}

func (cr *companyRepo) Delete(ctx context.Context, id uint64) (*v1.DeleteCompanysReply, error) {
	company, err := cr.data.companyuc.DeleteCompanys(ctx, &v1.DeleteCompanysRequest{
		Id: id,
	})

	if err != nil {
		return nil, err
	}

	return company, err
}
