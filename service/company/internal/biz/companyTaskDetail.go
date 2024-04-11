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
	CompanyTaskDetailJsonError    = errors.InternalServer("COMPANY_TASK_DETAIL_JSON_ERROR", "任务细节 JSON 出错")
)

type CompanyTaskDetailRepo interface {
	GetById(context.Context, uint64) (*domain.CompanyTaskDetail, error)
	List(context.Context, int, int, uint64, []uint64, []domain.CompanyTaskClientKeyAndOpenId) ([]*domain.CompanyTaskDetail, error)
	Save(context.Context, *domain.CompanyTaskDetail) (*domain.CompanyTaskDetail, error)
	SaveList(context.Context, []*domain.CompanyTaskDetail) error
	Update(context.Context, *domain.CompanyTaskDetail) (*domain.CompanyTaskDetail, error)
	UpdateOnDuplicateKey(context.Context, []*domain.CompanyTaskDetail) error
	Count(context.Context, uint64, []uint64, []domain.CompanyTaskClientKeyAndOpenId) (int64, error)
	CountIsPlauSuccess(context.Context, uint64, uint64) (int64, error)
	ListByClientKeyAndOpenIds(context.Context, uint64, uint64, string, string) (*v1.ListByClientKeyAndOpenIdsReply, error)
	DeleteOpenDouyinUsers(context.Context, []uint64) error
}

type CompanyTaskDetailUsecase struct {
	repo    CompanyTaskDetailRepo
	ctarepo CompanyTaskAccountRelationRepo
	wurepo  WeixinUserRepo
	conf    *conf.Data
	log     *log.Helper
}

func NewCompanyTaskDetailUsecase(repo CompanyTaskDetailRepo, ctarepo CompanyTaskAccountRelationRepo, wurepo WeixinUserRepo, conf *conf.Data, logger log.Logger) *CompanyTaskDetailUsecase {
	return &CompanyTaskDetailUsecase{repo: repo, ctarepo: ctarepo, wurepo: wurepo, conf: conf, log: log.NewHelper(logger)}
}

// first,use nickname to get opid in weixin
// if have nickname condition,to weixin get key
func (ctduc *CompanyTaskDetailUsecase) ListCompanyTaskDetail(ctx context.Context, pageNum, pageSize uint64, taskId uint64, keyword string) (*domain.CompanyTaskAccountRelationList, error) {
	var err error
	var weixinUserList *v1.ListByIdsReply
	weixinUserMap := make(map[uint64]*v1.ListByIdsReply_User)
	keys := []*domain.UserOpenDouyin{}
	detailMap := make(map[uint64][]*domain.CompanyTaskAccountRelation)
	userIds := []uint64{}
	relations := []*domain.CompanyTaskAccountRelation{}
	var total int64

	if len(keyword) > 0 {
		// 有查询就在条件中获取
		weixinUserList, err = ctduc.wurepo.ListByIds(ctx, "", keyword, "")

		if err != nil {
			return nil, CompanyTaskDetailListError
		}

		for _, v := range weixinUserList.Data.List {
			userIds = append(userIds, v.Id)
		}

		relations, err = ctduc.ctarepo.ListByUserIds(ctx, taskId, userIds)

		if err != nil {
			return nil, CompanyTaskDetailListError
		}

		count, err := ctduc.ctarepo.CountByUserIds(ctx, taskId, userIds)

		if err != nil {
			return nil, CompanyTaskDetailListError
		}

		total = count
	} else {
		relations, err = ctduc.ctarepo.List(ctx, taskId, 0, int(pageNum), int(pageSize), -1, "", "", "")

		if err != nil {
			return nil, CompanyTaskDetailListError
		}

		count, err := ctduc.ctarepo.Count(ctx, taskId, 0)

		if err != nil {
			return nil, CompanyTaskDetailListError
		}

		total = count

		for _, v := range relations {
			userIds = append(userIds, v.UserId)
		}

		b, err := json.Marshal(userIds)

		if err != nil {
			return nil, CompanyTaskDetailJsonError
		}

		weixinUserList, err = ctduc.wurepo.ListByIds(ctx, "", "", string(b))

		if err != nil {
			return nil, CompanyTaskDetailListError
		}
	}

	for _, v := range weixinUserList.Data.List {
		weixinUserMap[v.Id] = v
	}

	for _, relation := range relations {
		user := weixinUserMap[relation.UserId]

		if user != nil {
			relation.SetNickName(ctx, user.NickName)
			relation.SetAvatarUrl(ctx, user.AvatarUrl)
		}

		for _, detail := range relation.CompanyTaskDetails {
			keys = append(keys, &domain.UserOpenDouyin{
				ClientKey: detail.ClientKey,
				OpenId:    detail.OpenId,
			})

			if detail.IsPlaySuccess == 1 {
				relation.SetIsPlaySuccess(ctx, detail.IsPlaySuccess)
			}
		}
	}

	if len(keys) > 0 {
		b, err := json.Marshal(keys)

		if err != nil {
			return nil, CompanyTaskDetailJsonError
		}

		users, err := ctduc.repo.ListByClientKeyAndOpenIds(ctx, 0, 40, string(b), "")

		if err != nil {
			return nil, CompanyTaskGetDouyinUserError
		}

		for _, relation := range relations {
			for _, t := range relation.CompanyTaskDetails {
				t.SetNicknameAndAvatar(ctx, users.Data.List)
			}

			detailMap[relation.UserId] = append(detailMap[relation.UserId], relation)
		}
	}

	return &domain.CompanyTaskAccountRelationList{
		PageNum:  pageNum,
		PageSize: pageSize,
		Total:    uint64(total),
		List:     relations,
	}, nil
}
