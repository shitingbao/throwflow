package biz

import (
	"context"
	"encoding/json"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	v1 "interface/api/service/company/v1"
	vd1 "interface/api/service/douyin/v1"
	"interface/internal/conf"
	"interface/internal/pkg/tool"
	"strconv"
	"strings"
	"unicode/utf8"
)

type CompanyUserRepo interface {
	GetCompanyUser(context.Context, string) (*v1.GetCompanyUserReply, error)
	ChangeCompanyUserCompany(context.Context, string, uint64) (*v1.ChangeCompanyUserCompanyReply, error)
	ListCompanyUserMenu(context.Context, uint64) (*v1.ListCompanyUserMenuReply, error)
	ListQianchuanAdvertisers(context.Context, uint64, uint64) (*v1.ListQianchuanAdvertisersCompanyUsersReply, error)
	List(context.Context, uint64, uint64, uint64) (*v1.ListCompanyUsersReply, error)
	ListByPhone(context.Context, string) (*v1.ListCompanyUsersByPhoneReply, error)
	Statistics(context.Context, uint64) (*v1.StatisticsCompanyUsersReply, error)
	Get(context.Context, uint64, uint64) (*v1.GetCompanyUsersReply, error)
	Save(context.Context, uint64, string, string, string, uint32) (*v1.CreateCompanyUsersReply, error)
	Update(context.Context, uint64, uint64, string, string, string, uint32) (*v1.UpdateCompanyUsersReply, error)
	UpdateStatus(context.Context, uint64, uint64, uint32) (*v1.UpdateStatusCompanyUsersReply, error)
	UpdateRole(context.Context, uint64, uint64, string) (*v1.UpdateRoleCompanyUsersReply, error)
	Delete(context.Context, uint64, uint64) (*v1.DeleteCompanyUsersReply, error)
}

type CompanyUserUsecase struct {
	repo   CompanyUserRepo
	qarepo QianchuanAdvertiserRepo
	ocrepo OceanengineConfigRepo
	conf   *conf.Data
	log    *log.Helper
}

func NewCompanyUserUsecase(repo CompanyUserRepo, qarepo QianchuanAdvertiserRepo, ocrepo OceanengineConfigRepo, conf *conf.Data, logger log.Logger) *CompanyUserUsecase {
	return &CompanyUserUsecase{repo: repo, qarepo: qarepo, ocrepo: ocrepo, conf: conf, log: log.NewHelper(logger)}
}

func (cuuc *CompanyUserUsecase) GetCompanyUser(ctx context.Context, verifyLoginStatus bool) (*v1.GetCompanyUserReply, error) {
	token := ctx.Value("token")

	companyUser, err := cuuc.repo.GetCompanyUser(ctx, token.(string))

	if err != nil {
		if tool.GetGRPCErrorCode(err) == "Unauthenticated" {
			return nil, errors.Unauthorized("INTERFACE_LOGIN_FAILED", tool.GetGRPCErrorInfo(err))
		} else {
			return nil, errors.InternalServer("INTERFACE_LOGIN_FAILED", tool.GetGRPCErrorInfo(err))
		}
	}

	if verifyLoginStatus {
		if l := utf8.RuneCountInString(companyUser.Data.Reason); l > 0 {
			return nil, errors.InternalServer("INTERFACE_lOGIN_ABNORMAL", companyUser.Data.Reason)
		}
	}

	return companyUser, nil
}

func (cuuc *CompanyUserUsecase) ChangeCompanyUserCompany(ctx context.Context, companyId uint64) error {
	token := ctx.Value("token")

	if _, err := cuuc.repo.ChangeCompanyUserCompany(ctx, token.(string), companyId); err != nil {
		return errors.InternalServer("INTERFACE_CHANGE_COMPANY_USER_COMPANY_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return nil
}

func (cuuc *CompanyUserUsecase) ListCompanyUserMenu(ctx context.Context, companyId uint64) (*v1.ListCompanyUserMenuReply, error) {
	menus, err := cuuc.repo.ListCompanyUserMenu(ctx, companyId)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_LIST_COMPANY_USER_MENU_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return menus, nil
}

func (cuuc *CompanyUserUsecase) ListCompanyUserQianchuanAdvertisers(ctx context.Context, pageNum, pageSize, companyId uint64, keyword string) (*vd1.ListQianchuanAdvertisersReply, error) {
	if pageSize == 0 {
		pageSize = uint64(cuuc.conf.Database.PageSize)
	}

	list, err := cuuc.qarepo.List(ctx, companyId, pageNum, pageSize, keyword)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_LIST_COMPANY_USER_ADVERTISER_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return list, nil
}

func (cuuc *CompanyUserUsecase) StatisticsQianchuanAdvertisers(ctx context.Context, companyId uint64) (*vd1.StatisticsQianchuanAdvertisersReply, error) {
	list, err := cuuc.qarepo.Statistics(ctx, companyId)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_STATISTICS_COMPANY_USER_ADVERTISER_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return list, nil
}

func (cuuc *CompanyUserUsecase) GetUrlCompanyUserOceanengineAccounts(ctx context.Context, oceanengineType uint32, companyId uint64) (*vd1.RandOceanengineConfigsReply, error) {
	oceanengineConfig, err := cuuc.ocrepo.Rand(ctx, oceanengineType)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_GET_URL_COMPANY_USER_OCEANENGINE_ACCOUNTS_FAILED", tool.GetGRPCErrorInfo(err))
	}

	urlParams := make(map[string]string)
	urlParams["appId"] = oceanengineConfig.Data.AppId
	urlParams["companyId"] = strconv.FormatUint(companyId, 10)

	surlParams, _ := json.Marshal(urlParams)

	oceanengineConfig.Data.RedirectUrl = strings.Replace(oceanengineConfig.Data.RedirectUrl, "your_custom_params", string(surlParams), -1)

	return oceanengineConfig, nil
}

func (cuuc *CompanyUserUsecase) UpdateStatusCompanyUserQianchuanAdvertisers(ctx context.Context, companyId, advertiserId uint64, status uint32) (*vd1.UpdateStatusQianchuanAdvertisersReply, error) {
	qianchuanAdvertiser, err := cuuc.qarepo.Update(ctx, companyId, advertiserId, status)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_UPDATE_STATUS_COMPANY_USER_ADVERTISER_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return qianchuanAdvertiser, nil
}

func (cuuc *CompanyUserUsecase) ListCompanyUsers(ctx context.Context, pageNum, pageSize, companyId uint64) (*v1.ListCompanyUsersReply, error) {
	if pageSize == 0 {
		pageSize = uint64(cuuc.conf.Database.PageSize)
	}

	companyUsers, err := cuuc.repo.List(ctx, companyId, pageNum, pageSize)

	if err != nil {
		return nil, errors.InternalServer("Interface_DATA_ERROR", tool.GetGRPCErrorInfo(err))
	}

	return companyUsers, nil
}

func (cuuc *CompanyUserUsecase) StatisticsCompanyUsers(ctx context.Context, companyId uint64) (*v1.StatisticsCompanyUsersReply, error) {
	list, err := cuuc.repo.Statistics(ctx, companyId)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_COMPANY_USER_STATISTICS_ERROR", tool.GetGRPCErrorInfo(err))
	}

	return list, nil
}

func (cuuc *CompanyUserUsecase) CreateCompanyUsers(ctx context.Context, companyId uint64, username, job, phone string, role uint32) (*v1.CreateCompanyUsersReply, error) {
	if role == 1 {
		return nil, errors.InternalServer("INTERFACE_COMPANY_USER_CREATE_ERROR", "不能新增主管理员")
	}

	inCompanyUser, err := cuuc.repo.Save(ctx, companyId, username, job, phone, role)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_COMPANY_USER_CREATE_ERROR", tool.GetGRPCErrorInfo(err))
	}

	return inCompanyUser, nil
}

func (cuuc *CompanyUserUsecase) UpdateCompanyUsers(ctx context.Context, id, companyId uint64, username, job, phone string, role uint32) (*v1.UpdateCompanyUsersReply, error) {
	dCompanyUser, err := cuuc.repo.Get(ctx, id, companyId)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_COMPANY_USER_UPDATE_ERROR", tool.GetGRPCErrorInfo(err))
	}

	if dCompanyUser.Data.Role == 1 {
		return nil, errors.InternalServer("INTERFACE_COMPANY_USER_UPDATE_ERROR", "不能对主管理员进行操作")
	}

	if role == 1 {
		return nil, errors.InternalServer("INTERFACE_COMPANY_USER_UPDATE_ERROR", "不能修改为主管理员")
	}

	inCompanyUser, err := cuuc.repo.Update(ctx, id, companyId, username, job, phone, role)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_COMPANY_USER_UPDATE_ERROR", tool.GetGRPCErrorInfo(err))
	}

	return inCompanyUser, nil
}

func (cuuc *CompanyUserUsecase) UpdateStatusCompanyUsers(ctx context.Context, id, companyId uint64, status uint32) (*v1.UpdateStatusCompanyUsersReply, error) {
	dCompanyUser, err := cuuc.repo.Get(ctx, id, companyId)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_COMPANY_USER_UPDATE_ERROR", tool.GetGRPCErrorInfo(err))
	}

	if dCompanyUser.Data.Role == 1 {
		return nil, errors.InternalServer("INTERFACE_COMPANY_USER_UPDATE_ERROR", "不能对主管理员进行操作")
	}

	inCompanyUser, err := cuuc.repo.UpdateStatus(ctx, id, companyId, status)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_COMPANY_USER_UPDATE_ERROR", tool.GetGRPCErrorInfo(err))
	}

	return inCompanyUser, nil
}

func (cuuc *CompanyUserUsecase) ListQianchuanAdvertisersCompanyUsers(ctx context.Context, id, companyId uint64) (*v1.ListQianchuanAdvertisersCompanyUsersReply, error) {
	dCompanyUser, err := cuuc.repo.Get(ctx, id, companyId)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_COMPANY_USER_LIST_QIANCHUAN_ADVERTISERS_ERROR", tool.GetGRPCErrorInfo(err))
	}

	if dCompanyUser.Data.Role == 1 {
		return nil, errors.InternalServer("INTERFACE_COMPANY_USER_LIST_QIANCHUAN_ADVERTISERS_ERROR", "主管理员不需要分配广告权限")
	}

	list, err := cuuc.repo.ListQianchuanAdvertisers(ctx, id, companyId)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_COMPANY_USER_LIST_QIANCHUAN_ADVERTISERS_ERROR", tool.GetGRPCErrorInfo(err))
	}

	return list, nil
}

func (cuuc *CompanyUserUsecase) UpdateRoleCompanyUsers(ctx context.Context, id, companyId uint64, roleIds string) (*v1.UpdateRoleCompanyUsersReply, error) {
	dCompanyUser, err := cuuc.repo.Get(ctx, id, companyId)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_COMPANY_USER_UPDATE_ERROR", tool.GetGRPCErrorInfo(err))
	}

	if dCompanyUser.Data.Role == 1 {
		return nil, errors.InternalServer("INTERFACE_COMPANY_USER_UPDATE_ERROR", "不能对主管理员进行操作")
	}

	if len(roleIds) == 0 {
		return nil, errors.InternalServer("INTERFACE_COMPANY_USER_UPDATE_ERROR", "至少选择一个广告权限")
	}

	sroleIds := strings.Split(roleIds, ",")
	droleIds := make([]string, 0)

	for _, sroleId := range sroleIds {
		isExit := true

		for _, droleId := range droleIds {
			if droleId == sroleId {
				isExit = false
				break
			}
		}

		if isExit {
			droleIds = append(droleIds, sroleId)
		}
	}

	inCompanyUser, err := cuuc.repo.UpdateRole(ctx, id, companyId, strings.Join(droleIds, ","))

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_COMPANY_USER_UPDATE_ERROR", tool.GetGRPCErrorInfo(err))
	}

	return inCompanyUser, nil
}

func (cuuc *CompanyUserUsecase) DeleteCompanyUsers(ctx context.Context, id, cid, companyId uint64) error {
	dCompanyUser, err := cuuc.repo.Get(ctx, id, companyId)

	if err != nil {
		return errors.InternalServer("INTERFACE_COMPANY_USER_DELETE_ERROR", tool.GetGRPCErrorInfo(err))
	}

	if dCompanyUser.Data.Role == 1 {
		return errors.InternalServer("INTERFACE_COMPANY_USER_DELETE_ERROR", "不能对主管理员进行操作")
	}

	if dCompanyUser.Data.Id == cid {
		return errors.InternalServer("INTERFACE_COMPANY_USER_DELETE_ERROR", "不能对自己进行删除操作")
	}

	if _, err := cuuc.repo.Delete(ctx, id, companyId); err != nil {
		return errors.InternalServer("INTERFACE_COMPANY_USER_DELETE_ERROR", tool.GetGRPCErrorInfo(err))
	}

	return nil
}
