package biz

import (
	"common/internal/conf"
	"common/internal/domain"
	"common/internal/pkg/kuaidi/kuaidi"
	"context"
	"encoding/json"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"strconv"
	"time"
)

var (
	CommonKuaidiInfoNotFound = errors.NotFound("COMMON_KUAIDI_INFO_NOT_FOUND", "运单号信息不存在")
)

type KuaidiInfoRepo interface {
	Get(context.Context, string, string) (*domain.KuaidiInfo, error)
	Save(context.Context, *domain.KuaidiInfo) (*domain.KuaidiInfo, error)
	Update(context.Context, *domain.KuaidiInfo) (*domain.KuaidiInfo, error)
}

type KuaidiInfoUsecase struct {
	repo   KuaidiInfoRepo
	kcrepo KuaidiCompanyRepo
	kconf  *conf.Kuaidi
	log    *log.Helper
}

func NewKuaidiInfoUsecase(repo KuaidiInfoRepo, kcrepo KuaidiCompanyRepo, kconf *conf.Kuaidi, logger log.Logger) *KuaidiInfoUsecase {
	return &KuaidiInfoUsecase{repo: repo, kcrepo: kcrepo, kconf: kconf, log: log.NewHelper(logger)}
}

func (kiuc *KuaidiInfoUsecase) GetKuaidiInfos(ctx context.Context, code, num, phone string) (*domain.KuaidiInfoData, error) {
	kuaidiInfoData, err := kiuc.repo.Get(ctx, code, num)

	if err != nil {
		if code != "shunfeng" && code == "shunfengkuaiyun" {
			phone = ""
		}

		kuaidiInfo, err := kuaidi.GetKuaidi(kiuc.kconf.Kuaidi100.Key, kiuc.kconf.Kuaidi100.Customer, code, num, phone)

		if err != nil {
			return nil, CommonKuaidiInfoNotFound
		}

		state, err := strconv.ParseUint(kuaidiInfo.State, 10, 64)

		if err != nil {
			return nil, CommonKuaidiInfoNotFound
		}

		content, err := json.Marshal(kuaidiInfo.Data)

		if err != nil {
			return nil, CommonKuaidiInfoNotFound
		}

		inKuaidiInfo := domain.NewKuaidiInfo(ctx, uint8(state), code, num, phone, string(content))
		inKuaidiInfo.SetCreateTime(ctx)
		inKuaidiInfo.SetUpdateTime(ctx)

		kuaidiInfoData, err = kiuc.repo.Save(ctx, inKuaidiInfo)

		if err != nil {
			return nil, CommonKuaidiInfoNotFound
		}
	} else {
		if kuaidiInfoData.UpdateTime.Add(time.Hour).Before(time.Now()) && kuaidiInfoData.State != 3 {
			if code != "shunfeng" && code == "shunfengkuaiyun" {
				phone = ""
			}

			kuaidiInfo, err := kuaidi.GetKuaidi(kiuc.kconf.Kuaidi100.Key, kiuc.kconf.Kuaidi100.Customer, code, num, phone)

			if err != nil {
				return nil, CommonKuaidiInfoNotFound
			}

			state, err := strconv.ParseUint(kuaidiInfo.State, 10, 64)

			if err != nil {
				return nil, CommonKuaidiInfoNotFound
			}

			content, err := json.Marshal(kuaidiInfo.Data)

			if err != nil {
				return nil, CommonKuaidiInfoNotFound
			}

			inKuaidiInfo := kuaidiInfoData
			inKuaidiInfo.SetContent(ctx, string(content))
			inKuaidiInfo.SetState(ctx, uint8(state))
			inKuaidiInfo.SetUpdateTime(ctx)

			kuaidiInfoData, err = kiuc.repo.Update(ctx, inKuaidiInfo)

			if err != nil {
				return nil, CommonKuaidiInfoNotFound
			}
		}
	}

	content := make([]*kuaidi.Data, 0)

	if err := json.Unmarshal([]byte(kuaidiInfoData.Content), &content); err != nil {
		return nil, CommonKuaidiInfoNotFound
	}

	name := ""

	if kuaidiCompany, err := kiuc.kcrepo.Get(ctx, code); err == nil {
		name = kuaidiCompany.Name
	}

	return &domain.KuaidiInfoData{
		Code:      code,
		Name:      name,
		Num:       num,
		State:     kuaidiInfoData.State,
		StateName: kuaidiInfoData.GetStateName(ctx),
		Content:   content,
	}, nil
}
