package biz

import (
	"context"
	"encoding/json"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	v1 "interface/api/service/company/v1"
	"interface/internal/pkg/tool"
)

type ClueRepo interface {
	ListSelectClues(context.Context) (*v1.ListSelectCluesReply, error)
	UpdateCompanyName(context.Context, uint64, string) (*v1.UpdateCompanyNameCluesReply, error)
	Save(context.Context, string, string, string, uint32, uint32, uint64) (*v1.CreateCluesReply, error)
}

type ContactInformation struct {
	ContactUsername string `json:"contactUsername"`
	ContactPosition string `json:"contactPosition"`
	ContactPhone    string `json:"contactPhone"`
}

type ContactInformations []ContactInformation

type ClueUsecase struct {
	repo ClueRepo
	log  *log.Helper
}

func NewClueUsecase(repo ClueRepo, logger log.Logger) *ClueUsecase {
	return &ClueUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (cuc *ClueUsecase) ListSelectClues(ctx context.Context) (*v1.ListSelectCluesReply, error) {
	list, err := cuc.repo.ListSelectClues(ctx)

	if err != nil {
		return nil, InterfaceDataError
	}

	return list, nil
}

func (cuc *ClueUsecase) SaveClues(ctx context.Context, companyName, contactPhone string, companyType uint32, areaCode uint64) (*v1.CreateCluesReply, error) {
	var status uint32 = 1
	source := "官网"

	contactInformations := make([]ContactInformation, 0)
	contactInformations = append(contactInformations, ContactInformation{
		ContactUsername: companyName,
		ContactPosition: "",
		ContactPhone:    contactPhone,
	})

	contactInformationsByte, err := json.Marshal(contactInformations)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_APPLY_FORM_ERROR", "申请表单提交失败")
	}

	applyForm, err := cuc.repo.Save(ctx, companyName, string(contactInformationsByte), source, companyType, status, areaCode)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_APPLY_FORM_CREATE_ERROR", tool.GetGRPCErrorInfo(err))
	}

	return applyForm, nil
}

func (cuc *ClueUsecase) UpdateClues(ctx context.Context, companyId uint64, companyName string) (*v1.UpdateCompanyNameCluesReply, error) {
	clue, err := cuc.repo.UpdateCompanyName(ctx, companyId, companyName)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_UPDATE_CLUES_ERROR", tool.GetGRPCErrorInfo(err))
	}

	return clue, nil
}
