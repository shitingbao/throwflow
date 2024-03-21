package biz

import (
	v1 "company/api/service/weixin/v1"
	"company/internal/conf"
	"company/internal/domain"
	"context"
	"encoding/json"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	CompanyTaskGetDouyinUserError = errors.InternalServer("COMPANY_TASK_GET_DOUYIN_USER_ERROR", "获取抖音账号信息出错")
	CompanyTaskDetailListError    = errors.InternalServer("COMPANY_TASK_DETAIL_LIST_ERROR", "获取任务细节列表出错")
	CompanyTaskDetailCreateError  = errors.InternalServer("COMPANY_TASK_DETAIL_CREATE_ERROR", "新建任务细节出错")
	CompanyTaskDetailUpdateError  = errors.InternalServer("COMPANY_TASK_DETAIL_UPDATE_ERROR", "任务细节更新出错")
)

type CompanyTaskDetailRepo interface {
	GetById(context.Context, uint64) (*domain.CompanyTaskDetail, error)
	List(context.Context, int, int, uint64, []uint64, []domain.CompanyTaskClientKeyAndOpenId) ([]*domain.CompanyTaskDetail, error)
	Save(context.Context, *domain.CompanyTaskDetail) (*domain.CompanyTaskDetail, error)
	SaveList(context.Context, []*domain.CompanyTaskDetail) error
	Update(context.Context, *domain.CompanyTaskDetail) (*domain.CompanyTaskDetail, error)
	UpdateOnDuplicateKey(context.Context, []*domain.CompanyTaskDetail) error
	Count(context.Context, uint64, []domain.CompanyTaskClientKeyAndOpenId) (int64, error)
	ListByClientKeyAndOpenIds(ctx context.Context, pageNum, pageSize uint64, clientKeyAndOpenIds, keyword string) (*v1.ListByClientKeyAndOpenIdsReply, error)
	DeleteOpenDouyinUsers(context.Context, uint64, []string, []string) error
}

type CompanyTaskDetailUsecase struct {
	repo CompanyTaskDetailRepo
	conf *conf.Data
	log  *log.Helper
}

func NewCompanyTaskDetailUsecase(repo CompanyTaskDetailRepo, conf *conf.Data, logger log.Logger) *CompanyTaskDetailUsecase {
	return &CompanyTaskDetailUsecase{repo: repo, conf: conf, log: log.NewHelper(logger)}
}

// first,use nickname to get opid in weixin
// if have nickname condition,to weixin get key
func (tuc *CompanyTaskDetailUsecase) ListCompanyTaskDetail(ctx context.Context,
	pageNum, pageSize uint64, taskId uint64, nickname string) (*domain.CompanyTaskDetailList, error) {

	keys := []domain.CompanyTaskClientKeyAndOpenId{}
	names := []*v1.ListByClientKeyAndOpenIdsReply_OpenDouyinUser{}
	userIds := []uint64{}

	if len(nickname) > 0 {
		users, err := tuc.repo.ListByClientKeyAndOpenIds(ctx, 0, 40, "", nickname)
		if err != nil {
			return nil, CompanyTaskGetDouyinUserError
		}

		names = users.Data.List

		for _, v := range users.Data.List {
			key := domain.CompanyTaskClientKeyAndOpenId{
				ClientKey: v.ClientKey,
				OpenId:    v.OpenId,
			}
			userIds = append(userIds, v.OpenDouyinUserId)
			keys = append(keys, key)
		}
	}

	total, err := tuc.repo.Count(ctx, taskId, keys)

	if err != nil {
		return nil, CompanyTaskDetailListError
	}

	tasks, err := tuc.repo.List(ctx, int(pageNum), int(pageSize), taskId, userIds, keys)

	if err != nil {
		return nil, CompanyTaskDetailListError
	}

	// 都需要反馈达人名称和头像
	// 如果是达人名称模糊查询，先查询对应的 client key 和 openid，结果中有名称和头像，然后将反馈在结果中
	// 如果是普通列表查询，需要去根据 client key 和 openid 列表查一下再反馈
	if len(nickname) == 0 {
		keys := []*domain.UserOpenDouyin{}

		for _, t := range tasks {
			keys = append(keys, &domain.UserOpenDouyin{
				ClientKey: t.ClientKey,
				OpenId:    t.OpenId,
			})
		}

		b, err := json.Marshal(keys)

		if err != nil {
			return nil, CompanyTaskGetDouyinUserError
		}

		users, err := tuc.repo.ListByClientKeyAndOpenIds(ctx, 0, 40, string(b), "")

		if err != nil {
			return nil, CompanyTaskGetDouyinUserError
		}

		for _, t := range tasks {
			t.SetNicknameAndAvatar(ctx, users.Data.List)
		}
	} else {
		for _, t := range tasks {
			t.SetNicknameAndAvatarByCompanyIds(ctx, names)
		}
	}

	return &domain.CompanyTaskDetailList{
		PageNum:  pageNum,
		PageSize: pageSize,
		Total:    uint64(total),
		List:     tasks,
	}, nil
}
