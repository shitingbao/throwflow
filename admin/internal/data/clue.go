package data

import (
	v1 "admin/api/service/company/v1"
	"admin/internal/biz"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
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

func (cr *clueRepo) List(ctx context.Context, pageNum, industryId uint64, keyword string, status uint32) (*v1.ListCluesReply, error) {
	list, err := cr.data.companyuc.ListClues(ctx, &v1.ListCluesRequest{
		PageNum:    pageNum,
		PageSize:   uint64(cr.data.conf.Database.PageSize),
		Keyword:    keyword,
		IndustryId: industryId,
		Status:     status,
	})

	if err != nil {
		return nil, err
	}

	return list, err
}

func (cr *clueRepo) ListSelect(ctx context.Context) (*v1.ListSelectCluesReply, error) {
	list, err := cr.data.companyuc.ListSelectClues(ctx, &emptypb.Empty{})

	if err != nil {
		return nil, err
	}

	return list, err
}

func (cr *clueRepo) Statistics(ctx context.Context) (*v1.StatisticsCluesReply, error) {
	list, err := cr.data.companyuc.StatisticsClues(ctx, &emptypb.Empty{})

	if err != nil {
		return nil, err
	}

	return list, err
}

func (cr *clueRepo) Save(ctx context.Context, companyName, contactInformation, source, seller, facilitator, address, industryId string, userId, areaCode uint64, companyType, qianchuanUse, status uint32) (*v1.CreateCluesReply, error) {
	clue, err := cr.data.companyuc.CreateClues(ctx, &v1.CreateCluesRequest{
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
		Address:            address,
		AreaCode:           areaCode,
	})

	if err != nil {
		return nil, err
	}

	return clue, err
}

func (cr *clueRepo) Update(ctx context.Context, id, userId uint64, companyName, contactInformation, seller, facilitator, address, industryId string, areaCode uint64, companyType, qianchuanUse, status uint32) (*v1.UpdateCluesReply, error) {
	clue, err := cr.data.companyuc.UpdateClues(ctx, &v1.UpdateCluesRequest{
		Id:                 id,
		CompanyName:        companyName,
		IndustryId:         industryId,
		UserId:             userId,
		ContactInformation: contactInformation,
		CompanyType:        companyType,
		QianchuanUse:       qianchuanUse,
		Seller:             seller,
		Facilitator:        facilitator,
		Status:             status,
		Address:            address,
		AreaCode:           areaCode,
	})

	if err != nil {
		return nil, err
	}

	return clue, err
}

func (cr *clueRepo) UpdateOperationLog(ctx context.Context, id, userId uint64, content, operationTime string) (*v1.UpdateOperationLogCluesReply, error) {
	clue, err := cr.data.companyuc.UpdateOperationLogClues(ctx, &v1.UpdateOperationLogCluesRequest{
		Id:            id,
		UserId:        userId,
		Content:       content,
		OperationTime: operationTime,
	})

	if err != nil {
		return nil, err
	}

	return clue, err
}

func (cr *clueRepo) Delete(ctx context.Context, id uint64) (*v1.DeleteCluesReply, error) {
	clue, err := cr.data.companyuc.DeleteClues(ctx, &v1.DeleteCluesRequest{
		Id: id,
	})

	if err != nil {
		return nil, err
	}

	return clue, err
}
