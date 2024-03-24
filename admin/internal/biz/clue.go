package biz

import (
	v1 "admin/api/service/company/v1"
	"admin/internal/pkg/tool"
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"time"
	"unicode/utf8"
)

type ClueRepo interface {
	List(context.Context, uint64, uint64, string, uint32) (*v1.ListCluesReply, error)
	ListSelect(context.Context) (*v1.ListSelectCluesReply, error)
	Statistics(context.Context) (*v1.StatisticsCluesReply, error)
	Save(context.Context, string, string, string, string, string, string, string, uint64, uint64, uint32, uint32, uint32) (*v1.CreateCluesReply, error)
	Update(context.Context, uint64, uint64, string, string, string, string, string, string, uint64, uint32, uint32, uint32) (*v1.UpdateCluesReply, error)
	UpdateOperationLog(context.Context, uint64, uint64, string, string) (*v1.UpdateOperationLogCluesReply, error)
	Delete(context.Context, uint64) (*v1.DeleteCluesReply, error)
}

type ClueUsecase struct {
	repo  ClueRepo
	urepo UserRepo
	log   *log.Helper
}

func NewClueUsecase(repo ClueRepo, urepo UserRepo, logger log.Logger) *ClueUsecase {
	return &ClueUsecase{repo: repo, urepo: urepo, log: log.NewHelper(logger)}
}

func (cuc *ClueUsecase) ListClues(ctx context.Context, pageNum, industryId uint64, keyword string, status uint32) (*v1.ListCluesReply, error) {
	list, err := cuc.repo.List(ctx, pageNum, industryId, keyword, status)

	if err != nil {
		return nil, AdminDataError
	}

	for _, clue := range list.Data.List {
		for _, operationLog := range clue.OperationLogs {
			if operationLog.UserId > 0 {
				if user, err := cuc.urepo.GetById(ctx, operationLog.UserId); err == nil {
					if len(user.Nickname) > 0 {
						operationLog.UserName = user.Nickname
					} else {
						operationLog.UserName = user.Username
					}
				}
			}
		}
	}

	return list, nil
}

func (cuc *ClueUsecase) ListSelectClues(ctx context.Context) (*v1.ListSelectCluesReply, error) {
	list, err := cuc.repo.ListSelect(ctx)

	if err != nil {
		return nil, AdminDataError
	}

	return list, nil
}

func (cuc *ClueUsecase) StatisticsClues(ctx context.Context) (*v1.StatisticsCluesReply, error) {
	list, err := cuc.repo.Statistics(ctx)

	if err != nil {
		return nil, AdminDataError
	}

	return list, nil
}

func (cuc *ClueUsecase) CreateClues(ctx context.Context, companyName, contactInformation, seller, facilitator, address, industryId string, companyType, qianchuanUse uint32, userId, areaCode uint64) (*v1.CreateCluesReply, error) {
	var status uint32 = 1
	source := "录入"

	if l := utf8.RuneCountInString(seller); l > 1 {
		status = 2
	}

	clue, err := cuc.repo.Save(ctx, companyName, contactInformation, source, seller, facilitator, address, industryId, userId, areaCode, companyType, qianchuanUse, status)

	if err != nil {
		return nil, errors.InternalServer("ADMIN_CLUE_CREATE_ERROR", tool.GetGRPCErrorInfo(err))
	}

	for _, operationLog := range clue.Data.OperationLogs {
		if operationLog.UserId > 0 {
			if user, err := cuc.urepo.GetById(ctx, operationLog.UserId); err == nil {
				if len(user.Nickname) > 0 {
					operationLog.UserName = user.Nickname
				} else {
					operationLog.UserName = user.Username
				}
			}
		}
	}

	return clue, nil
}

func (cuc *ClueUsecase) UpdateClues(ctx context.Context, id, userId uint64, companyName, contactInformation, seller, facilitator, address, industryId string, companyType, qianchuanUse uint32, areaCode uint64) (*v1.UpdateCluesReply, error) {
	var status uint32 = 1

	if l := utf8.RuneCountInString(seller); l > 1 {
		status = 2
	}

	clue, err := cuc.repo.Update(ctx, id, userId, companyName, contactInformation, seller, facilitator, address, industryId, areaCode, companyType, qianchuanUse, status)

	if err != nil {
		return nil, errors.InternalServer("ADMIN_CLUE_UPDATE_ERROR", tool.GetGRPCErrorInfo(err))
	}

	for _, operationLog := range clue.Data.OperationLogs {
		if operationLog.UserId > 0 {
			if user, err := cuc.urepo.GetById(ctx, operationLog.UserId); err == nil {
				if len(user.Nickname) > 0 {
					operationLog.UserName = user.Nickname
				} else {
					operationLog.UserName = user.Username
				}
			}
		}
	}

	return clue, nil
}

func (cuc *ClueUsecase) UpdateOperationLogClues(ctx context.Context, id, userId uint64, content string, operationTime time.Time) (*v1.UpdateOperationLogCluesReply, error) {
	clue, err := cuc.repo.UpdateOperationLog(ctx, id, userId, content, tool.TimeToString("2006-01-02 15:04", operationTime))

	if err != nil {
		return nil, errors.InternalServer("ADMIN_CLUE_UPDATE_OPERATION_LOG_ERROR", tool.GetGRPCErrorInfo(err))
	}

	for _, operationLog := range clue.Data.OperationLogs {
		if operationLog.UserId > 0 {
			if user, err := cuc.urepo.GetById(ctx, operationLog.UserId); err == nil {
				if len(user.Nickname) > 0 {
					operationLog.UserName = user.Nickname
				} else {
					operationLog.UserName = user.Username
				}
			}
		}
	}

	return clue, nil
}

func (cuc *ClueUsecase) DeleteClues(ctx context.Context, id uint64) (*v1.DeleteCluesReply, error) {
	clue, err := cuc.repo.Delete(ctx, id)

	if err != nil {
		return nil, errors.InternalServer("ADMIN_CLUE_DELETE_ERROR", tool.GetGRPCErrorInfo(err))
	}

	return clue, nil
}
