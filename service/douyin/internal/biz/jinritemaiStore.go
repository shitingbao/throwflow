package biz

import (
	"context"
	"douyin/internal/conf"
	"douyin/internal/domain"
	"douyin/internal/pkg/event/event"
	"douyin/internal/pkg/event/kafka"
	"douyin/internal/pkg/jinritemai"
	"douyin/internal/pkg/jinritemai/store"
	"douyin/internal/pkg/tool"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"golang.org/x/time/rate"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	DouyinJinritemaiStoreGetCommissionRatioError = errors.InternalServer("DOUYIN_JINRITEMAI_STORE_GET_COMMISSION_RATIO_ERROR", "达人橱窗佣金获取失败")
	DouyinJinritemaiStoreListError               = errors.InternalServer("DOUYIN_JINRITEMAI_STORE_LIST_ERROR", "达人橱窗列表获取失败")
	DouyinJinritemaiStoreCreateError             = errors.InternalServer("DOUYIN_JINRITEMAI_STORE_CREATE_ERROR", "达人橱窗创建失败")
	DouyinJinritemaiStoreReplaceError            = errors.InternalServer("DOUYIN_JINRITEMAI_STORE_REPLACE_ERROR", "达人橱窗替换高佣链接失败")
	DouyinJinritemaiStoreDeleteError             = errors.InternalServer("DOUYIN_JINRITEMAI_STORE_DELETE_ERROR", "达人橱窗删除失败")
)

type JinritemaiStoreRepo interface {
	SaveIndex(context.Context, string)
	Upsert(context.Context, string, *domain.JinritemaiStore) error
	DeleteByDayAndClientKeyAndOpenId(context.Context, string, string, string) error

	Send(context.Context, event.Event) error
}

type JinritemaiStoreUsecase struct {
	repo     JinritemaiStoreRepo
	jsirepo  JinritemaiStoreInfoRepo
	wurepo   WeixinUserRepo
	wuodrepo WeixinUserOpenDouyinRepo
	odtrepo  OpenDouyinTokenRepo
	cprepo   CompanyProductRepo
	tlrepo   TaskLogRepo
	jalrepo  JinritemaiApiLogRepo
	conf     *conf.Data
	cconf    *conf.Company
	dconf    *conf.Developer
	log      *log.Helper
}

func NewJinritemaiStoreUsecase(repo JinritemaiStoreRepo, jsirepo JinritemaiStoreInfoRepo, wurepo WeixinUserRepo, wuodrepo WeixinUserOpenDouyinRepo, odtrepo OpenDouyinTokenRepo, cprepo CompanyProductRepo, tlrepo TaskLogRepo, jalrepo JinritemaiApiLogRepo, conf *conf.Data, cconf *conf.Company, dconf *conf.Developer, logger log.Logger) *JinritemaiStoreUsecase {
	return &JinritemaiStoreUsecase{repo: repo, jsirepo: jsirepo, wurepo: wurepo, wuodrepo: wuodrepo, odtrepo: odtrepo, cprepo: cprepo, tlrepo: tlrepo, jalrepo: jalrepo, conf: conf, cconf: cconf, dconf: dconf, log: log.NewHelper(logger)}
}

func (jsuc *JinritemaiStoreUsecase) ListJinritemaiStores(ctx context.Context, pageNum, pageSize, userId uint64) (*domain.JinritemaiStoreList, error) {
	weixinUser, err := jsuc.wurepo.GetById(ctx, userId)

	if err != nil {
		return nil, DouyinWeixinUserNotFound
	}

	weixinUserOpenDouyins, err := jsuc.wuodrepo.List(ctx, weixinUser.Data.UserId)

	if err != nil {
		return nil, DouyinWeixinUserOpenDouyinListError
	}

	var total int64
	list := make([]*domain.JinritemaiStoreInfo, 0)

	if len(weixinUserOpenDouyins.Data.List) > 0 {
		openDouyinTokens := make([]*domain.OpenDouyinToken, 0)

		for _, weixinUserOpenDouyin := range weixinUserOpenDouyins.Data.List {
			openDouyinTokens = append(openDouyinTokens, &domain.OpenDouyinToken{
				ClientKey: weixinUserOpenDouyin.ClientKey,
				OpenId:    weixinUserOpenDouyin.OpenId,
			})
		}

		jinritemaiStores, err := jsuc.jsirepo.List(ctx, int(pageNum), int(pageSize), openDouyinTokens)

		if err != nil {
			return nil, DouyinJinritemaiStoreListError
		}

		total, err = jsuc.jsirepo.Count(ctx, openDouyinTokens)

		if err != nil {
			return nil, DouyinJinritemaiStoreListError
		}

		for _, jinritemaiStore := range jinritemaiStores {
			jinritemaiStore.SetCommissionTypeName(ctx)
			jinritemaiStore.SetPromotionTypeName(ctx)

			list = append(list, &domain.JinritemaiStoreInfo{
				Id:                 jinritemaiStore.Id,
				ClientKey:          jinritemaiStore.ClientKey,
				OpenId:             jinritemaiStore.OpenId,
				ProductId:          jinritemaiStore.ProductId,
				ProductName:        jinritemaiStore.ProductName,
				ProductImg:         jinritemaiStore.ProductImg,
				ProductPrice:       jinritemaiStore.ProductPrice,
				CommissionType:     jinritemaiStore.CommissionType,
				CommissionTypeName: jinritemaiStore.CommissionTypeName,
				CommissionRatio:    jinritemaiStore.CommissionRatio,
				PromotionType:      jinritemaiStore.PromotionType,
				PromotionTypeName:  jinritemaiStore.PromotionTypeName,
				ColonelActivityId:  jinritemaiStore.ColonelActivityId,
				CreateTime:         jinritemaiStore.CreateTime,
				UpdateTime:         jinritemaiStore.UpdateTime,
			})
		}
	}

	return &domain.JinritemaiStoreList{
		PageNum:  pageNum,
		PageSize: pageSize,
		Total:    uint64(total),
		List:     list,
	}, nil
}

func (jsuc *JinritemaiStoreUsecase) ListProductIdJinritemaiStores(ctx context.Context, pageNum, pageSize uint64) (*domain.JinritemaiStoreList, error) {
	list := make([]*domain.JinritemaiStoreInfo, 0)

	jinritemaiStores, err := jsuc.jsirepo.ListProductId(ctx, int(pageNum), int(pageSize))

	if err != nil {
		return nil, DouyinJinritemaiStoreListError
	}

	total, err := jsuc.jsirepo.CountProductId(ctx)

	if err != nil {
		return nil, DouyinJinritemaiStoreListError
	}

	for _, jinritemaiStore := range jinritemaiStores {
		jinritemaiStore.SetCommissionTypeName(ctx)
		jinritemaiStore.SetPromotionTypeName(ctx)

		list = append(list, &domain.JinritemaiStoreInfo{
			Id:                 jinritemaiStore.Id,
			ClientKey:          jinritemaiStore.ClientKey,
			OpenId:             jinritemaiStore.OpenId,
			ProductId:          jinritemaiStore.ProductId,
			ProductName:        jinritemaiStore.ProductName,
			ProductImg:         jinritemaiStore.ProductImg,
			ProductPrice:       jinritemaiStore.ProductPrice,
			CommissionType:     jinritemaiStore.CommissionType,
			CommissionTypeName: jinritemaiStore.CommissionTypeName,
			CommissionRatio:    jinritemaiStore.CommissionRatio,
			PromotionType:      jinritemaiStore.PromotionType,
			PromotionTypeName:  jinritemaiStore.PromotionTypeName,
			ColonelActivityId:  jinritemaiStore.ColonelActivityId,
			CreateTime:         jinritemaiStore.CreateTime,
			UpdateTime:         jinritemaiStore.UpdateTime,
		})
	}

	return &domain.JinritemaiStoreList{
		PageNum:  pageNum,
		PageSize: pageSize,
		Total:    uint64(total),
		List:     list,
	}, nil
}

func (jsuc *JinritemaiStoreUsecase) CreateJinritemaiStores(ctx context.Context, userId, companyId, productId uint64, openDouyinUserIds, activityUrl string) ([]*domain.AddStoreMessage, error) {
	weixinUserOpenDouyins, err := jsuc.wuodrepo.List(ctx, userId)

	if err != nil {
		return nil, DouyinWeixinUserOpenDouyinListError
	}

	sopenDouyinUserIds := strings.Split(openDouyinUserIds, ",")

	isNotExist := true

	openDouyinUserInfos := make([]*domain.OpenDouyinUserInfo, 0)

	for _, weixinUserOpenDouyin := range weixinUserOpenDouyins.Data.List {
		for _, sopenDouyinUserId := range sopenDouyinUserIds {
			if iopenDouyinUserId, err := strconv.ParseUint(sopenDouyinUserId, 10, 64); err == nil {
				if weixinUserOpenDouyin.OpenDouyinUserId == iopenDouyinUserId {
					openDouyinUserInfos = append(openDouyinUserInfos, &domain.OpenDouyinUserInfo{
						ClientKey: weixinUserOpenDouyin.ClientKey,
						OpenId:    weixinUserOpenDouyin.OpenId,
						Nickname:  weixinUserOpenDouyin.Nickname,
					})

					isNotExist = false

					break
				}
			}
		}
	}

	if isNotExist {
		return nil, DouyinWeixinUserOpenDouyinNotFound
	}

	list := make([]*domain.AddStoreMessage, 0)

	products := make([]*store.AddStoreBodyDataProductsRequest, 0)
	products = append(products, &store.AddStoreBodyDataProductsRequest{
		ProductId:   int64(productId),
		ActivityUrl: activityUrl,
	})

	companyProduct, cperr := jsuc.cprepo.GetByProductOutId(ctx, productId, "")

	for _, openDouyinUserInfo := range openDouyinUserInfos {
		if openDouyinToken, err := jsuc.odtrepo.GetByClientKeyAndOpenId(ctx, openDouyinUserInfo.ClientKey, openDouyinUserInfo.OpenId); err == nil {
			if stores, err := store.AddStore(openDouyinToken.AccessToken, openDouyinToken.OpenId, strconv.FormatUint(companyId, 10), products); err != nil {
				addStoreMessage := &domain.AddStoreMessage{
					AwemeName: "达人名称：" + openDouyinUserInfo.Nickname,
					Content:   "达人橱窗创建失败",
				}

				if cperr != nil {
					addStoreMessage.ProductName = "商品ID：" + strconv.FormatUint(productId, 10)
				} else {
					addStoreMessage.ProductName = "商品名称：" + companyProduct.Data.ProductName
				}

				list = append(list, addStoreMessage)
			} else {
				for _, store := range stores.Data.Results {
					if store.ErrorCode != 0 {
						addStoreMessage := &domain.AddStoreMessage{
							AwemeName: "达人名称：" + openDouyinUserInfo.Nickname,
							Content:   store.ErrorMsg,
						}

						if cperr != nil {
							addStoreMessage.ProductName = "商品ID：" + strconv.FormatUint(productId, 10)
						} else {
							addStoreMessage.ProductName = "商品名称：" + companyProduct.Data.ProductName
						}

						list = append(list, addStoreMessage)
					}
				}
			}
		} else {
			addStoreMessage := &domain.AddStoreMessage{
				AwemeName: "达人名称：" + openDouyinUserInfo.Nickname,
				Content:   "授权账户不存在",
			}

			list = append(list, addStoreMessage)
		}
	}

	return list, nil
}

func (jsuc *JinritemaiStoreUsecase) UpdateJinritemaiStores(ctx context.Context, userId uint64, stores string) (string, error) {
	weixinUser, err := jsuc.wurepo.GetById(ctx, userId)

	if err != nil {
		return "", DouyinWeixinUserNotFound
	}

	weixinUserOpenDouyins, err := jsuc.wuodrepo.List(ctx, weixinUser.Data.UserId)

	if err != nil {
		return "", DouyinWeixinUserOpenDouyinListError
	}

	inJinritemaiStore := domain.JinritemaiStore{}
	inJinritemaiStore.SetStore(ctx, stores)

	if ok := inJinritemaiStore.VerifyStores(ctx); !ok {
		return "", DouyinJinritemaiStoreReplaceError
	}

	var replaceProducts uint64 = 0

	type storeReplaceProduct struct {
		clientKey   string
		openId      string
		accessToken string
		products    []*store.AddStoreBodyDataProductsRequest
	}

	storeReplaceProducts := make([]*storeReplaceProduct, 0)

	if len(weixinUserOpenDouyins.Data.List) > 0 {
		wopenDouyinTokens := make([]*domain.OpenDouyinToken, 0)

		for _, store := range inJinritemaiStore.Stores {
			for _, weixinUserOpenDouyin := range weixinUserOpenDouyins.Data.List {
				if store.ClientKey == weixinUserOpenDouyin.ClientKey && store.OpenId == weixinUserOpenDouyin.OpenId {
					wopenDouyinTokens = append(wopenDouyinTokens, &domain.OpenDouyinToken{
						ClientKey: weixinUserOpenDouyin.ClientKey,
						OpenId:    weixinUserOpenDouyin.OpenId,
					})

					break
				}
			}
		}

		openDouyinTokens, err := jsuc.odtrepo.ListByClientKeyAndOpenId(ctx, wopenDouyinTokens)

		if err != nil {
			return "", DouyinOpenDouyinTokenListError
		}

		for _, lstore := range inJinritemaiStore.Stores {
			productIds := make([]string, 0)

			sproductIds := strings.Split(lstore.ProductId, ",")

			for _, sproductId := range sproductIds {
				productIds = append(productIds, sproductId)
			}

			if len(productIds) > 0 {
				if jinritemaiStores, err := jsuc.jsirepo.ListByClientKeyAndOpenIdAndProductIds(ctx, lstore.ClientKey, lstore.OpenId, productIds); err == nil {
					for _, jinritemaiStore := range jinritemaiStores {
						if productId, err := strconv.ParseUint(jinritemaiStore.ProductId, 10, 64); err == nil {
							if companyProduct, err := jsuc.cprepo.GetByProductOutId(ctx, productId, "1"); err == nil {
								var commissionRatio float64
								var differenceRatio float64
								var activityUrl string

								for _, commission := range companyProduct.Data.Commissions {
									if commissionRatio == 0 {
										commissionRatio = commission.CommissionRatio
										activityUrl = commission.CommissionOutUrl
									} else if commissionRatio > 0 {
										if commissionRatio > commission.CommissionRatio {
											commissionRatio = commission.CommissionRatio
											activityUrl = commission.CommissionOutUrl
										}
									}
								}

								if commissionRatio > 0 {
									differenceRatio = commissionRatio - float64(jinritemaiStore.CommissionRatio)

									if differenceRatio >= 0 && differenceRatio < 15 {
										accessToken := ""

										for _, openDouyinToken := range openDouyinTokens {
											if openDouyinToken.OpenId == jinritemaiStore.OpenId && openDouyinToken.ClientKey == jinritemaiStore.ClientKey {
												accessToken = openDouyinToken.AccessToken

												break
											}
										}

										products := make([]*store.AddStoreBodyDataProductsRequest, 0)
										products = append(products, &store.AddStoreBodyDataProductsRequest{
											ProductId:   int64(companyProduct.Data.ProductOutId),
											ActivityUrl: activityUrl,
										})

										storeReplaceProducts = append(storeReplaceProducts, &storeReplaceProduct{
											clientKey:   jinritemaiStore.ClientKey,
											openId:      jinritemaiStore.OpenId,
											accessToken: accessToken,
											products:    products,
										})
									}
								}
							}
						}
					}
				}
			}
		}

		if len(storeReplaceProducts) > 0 {
			for _, lstoreReplaceProduct := range storeReplaceProducts {
				if results, err := store.AddStore(lstoreReplaceProduct.accessToken, lstoreReplaceProduct.openId, strconv.FormatUint(jsuc.cconf.DefaultCompanyId, 10), lstoreReplaceProduct.products); err == nil {
					for _, result := range results.Data.Results {
						if result.ErrorCode == 0 {
							replaceProducts += 1
						}
					}
				}
			}
		}
	}

	return "替换了" + strconv.FormatUint(replaceProducts, 10) + "高佣链接", nil
}

func (jsuc *JinritemaiStoreUsecase) DeleteJinritemaiStores(ctx context.Context, userId uint64, stores string) (string, error) {
	weixinUser, err := jsuc.wurepo.GetById(ctx, userId)

	if err != nil {
		return "", DouyinWeixinUserNotFound
	}

	weixinUserOpenDouyins, err := jsuc.wuodrepo.List(ctx, weixinUser.Data.UserId)

	if err != nil {
		return "", DouyinWeixinUserOpenDouyinListError
	}

	inJinritemaiStore := domain.JinritemaiStore{}
	inJinritemaiStore.SetStore(ctx, stores)

	if ok := inJinritemaiStore.VerifyStores(ctx); !ok {
		return "", DouyinJinritemaiStoreDeleteError
	}

	var delProducts uint64 = 0

	if len(weixinUserOpenDouyins.Data.List) > 0 {
		wopenDouyinTokens := make([]*domain.OpenDouyinToken, 0)

		for _, store := range inJinritemaiStore.Stores {
			for _, weixinUserOpenDouyin := range weixinUserOpenDouyins.Data.List {
				if store.ClientKey == weixinUserOpenDouyin.ClientKey && store.OpenId == weixinUserOpenDouyin.OpenId {
					wopenDouyinTokens = append(wopenDouyinTokens, &domain.OpenDouyinToken{
						ClientKey: weixinUserOpenDouyin.ClientKey,
						OpenId:    weixinUserOpenDouyin.OpenId,
					})

					break
				}
			}
		}

		openDouyinTokens, err := jsuc.odtrepo.ListByClientKeyAndOpenId(ctx, wopenDouyinTokens)

		if err != nil {
			return "", DouyinOpenDouyinTokenListError
		}

		type storeRemoveProduct struct {
			clientKey   string
			openId      string
			accessToken string
			products    []*store.DelStoreBodyDataProductsRequest
		}

		storeRemoveProducts := make([]*storeRemoveProduct, 0)

		for _, lstore := range inJinritemaiStore.Stores {
			var accessToken string

			for _, openDouyinToken := range openDouyinTokens {
				if openDouyinToken.OpenId == lstore.OpenId && openDouyinToken.ClientKey == lstore.ClientKey {
					accessToken = openDouyinToken.AccessToken

					break
				}
			}

			products := make([]*store.DelStoreBodyDataProductsRequest, 0)

			productIds := make([]string, 0)

			sproductIds := strings.Split(lstore.ProductId, ",")

			for _, sproductId := range sproductIds {
				productIds = append(productIds, sproductId)
			}

			if len(productIds) > 0 {
				if jinritemaiStores, err := jsuc.jsirepo.ListByClientKeyAndOpenIdAndProductIds(ctx, lstore.ClientKey, lstore.OpenId, productIds); err == nil {
					for _, jinritemaiStore := range jinritemaiStores {
						if iproductId, err := strconv.ParseUint(jinritemaiStore.ProductId, 10, 64); err == nil {
							products = append(products, &store.DelStoreBodyDataProductsRequest{
								ProductId:   int64(iproductId),
								PromotionId: int64(jinritemaiStore.PromotionId),
							})
						}
					}
				}
			}

			if len(products) > 0 {
				storeRemoveProducts = append(storeRemoveProducts, &storeRemoveProduct{
					clientKey:   lstore.ClientKey,
					openId:      lstore.OpenId,
					accessToken: accessToken,
					products:    products,
				})
			}

			if len(storeRemoveProducts) > 0 {
				for _, lstoreRemoveProduct := range storeRemoveProducts {
					if _, err := store.DelStore(lstoreRemoveProduct.accessToken, lstoreRemoveProduct.openId, lstoreRemoveProduct.products); err == nil {
						delProducts += uint64(len(lstoreRemoveProduct.products))

						productIds := make([]string, 0)

						for _, product := range lstoreRemoveProduct.products {
							productIds = append(productIds, strconv.FormatInt(product.ProductId, 10))
						}

						jsuc.jsirepo.DeleteByDayAndClientKeyAndOpenIdAndProductIds(ctx, lstoreRemoveProduct.clientKey, lstoreRemoveProduct.openId, productIds)
					}
				}
			}
		}
	}

	return "移除了" + strconv.FormatUint(delProducts, 10) + "橱窗", nil
}

func (jsuc *JinritemaiStoreUsecase) SyncJinritemaiStores(ctx context.Context) error {
	openDouyinTokens, err := jsuc.odtrepo.List(ctx)

	if err != nil {
		inTaskLog := domain.NewTaskLog(ctx, "SyncJinritemaiStores", fmt.Sprintf("[SyncJinritemaiStoresError] Description=%s", "获取抖音开放平台token列表失败"))
		inTaskLog.SetCreateTime(ctx)

		jsuc.tlrepo.Save(ctx, inTaskLog)

		return err
	}

	day := tool.TimeToString("2006-01-02", time.Now())

	jsuc.repo.SaveIndex(ctx, day)

	var wg sync.WaitGroup

	limiter := rate.NewLimiter(0, 60)
	limiter.SetLimit(rate.Limit(60))

	for _, openDouyinToken := range openDouyinTokens {
		wg.Add(1)

		go jsuc.SyncJinritemaiStore(ctx, &wg, limiter, day, openDouyinToken)
	}

	wg.Wait()

	messageAd := domain.MessageAd{
		Type: "jinritemail_store_data_ready",
	}

	messageAd.Message.Name = "douyin"
	messageAd.Message.SyncDate = day
	messageAd.Message.SendTime = tool.TimeToString("2006-01-02 15:04:05", time.Now())

	bmessageAd, _ := json.Marshal(messageAd)

	jsuc.repo.Send(ctx, kafka.NewMessage(strconv.FormatUint(1, 10), bmessageAd))

	return nil
}

func (jsuc *JinritemaiStoreUsecase) SyncJinritemaiStore(ctx context.Context, wg *sync.WaitGroup, limiter *rate.Limiter, day string, openDouyinToken *domain.OpenDouyinToken) {
	defer wg.Done()

	var sjswg sync.WaitGroup

	sjswg.Add(1)

	if err := jsuc.repo.DeleteByDayAndClientKeyAndOpenId(ctx, day, openDouyinToken.ClientKey, openDouyinToken.OpenId); err != nil {
		inTaskLog := domain.NewTaskLog(ctx, "SyncJinritemaiStores", fmt.Sprintf("[SyncJinritemaiStoresError SyncJinritemaiStoreError] ClientKey=%d, OpenId=%d, Description=%s", openDouyinToken.ClientKey, openDouyinToken.OpenId, "同步精选联盟达人橱窗数据，删除数据库失败"))
		inTaskLog.SetCreateTime(ctx)

		jsuc.tlrepo.Save(ctx, inTaskLog)
	} else {
		stores, err := jsuc.listStores(ctx, &sjswg, limiter, day, openDouyinToken, 1)

		if err == nil {
			totalPage := int64(math.Ceil(float64(stores.Data.Total) / float64(jinritemai.PageSize20)))

			if totalPage > 1 {
				var page int64

				for page = 2; page <= totalPage; page++ {
					sjswg.Add(1)

					go jsuc.listStores(ctx, &sjswg, limiter, day, openDouyinToken, page)
				}
			}
		}

		sjswg.Wait()
	}
}

func (jsuc *JinritemaiStoreUsecase) listStores(ctx context.Context, sjswg *sync.WaitGroup, limiter *rate.Limiter, day string, openDouyinToken *domain.OpenDouyinToken, page int64) (*store.ListStoreResponse, error) {
	defer sjswg.Done()

	var stores *store.ListStoreResponse
	var err error

	for retryNum := 0; retryNum < 3; retryNum++ {
		limiter.Wait(ctx)

		stores, err = store.ListStore(openDouyinToken.AccessToken, openDouyinToken.OpenId, page)

		if err != nil {
			if retryNum == 2 {
				inJinritemaiApiLog := domain.NewJinritemaiApiLog(ctx, openDouyinToken.ClientKey, openDouyinToken.OpenId, openDouyinToken.AccessToken, err.Error())
				inJinritemaiApiLog.SetCreateTime(ctx)

				jsuc.jalrepo.Save(ctx, inJinritemaiApiLog)
			}
		} else {
			for _, store := range stores.Data.Results {
				inJinritemaiStore := domain.NewJinritemaiStore(ctx, store.ProductId, store.PromotionId, store.PromotionType, store.Price, store.CosType, store.ColonelActivityId, store.CosRatio, store.HideStatus, openDouyinToken.ClientKey, openDouyinToken.OpenId, store.Title, store.Cover)

				if err := jsuc.repo.Upsert(ctx, day, inJinritemaiStore); err != nil {
					sinJinritemaiStore, _ := json.Marshal(inJinritemaiStore)

					inTaskLog := domain.NewTaskLog(ctx, "SyncJinritemaiStores", fmt.Sprintf("[SyncJinritemaiStoresError SyncJinritemaiStoreError] ClientKey=%d, OpenId=%d, Data=%s, Description=%s", openDouyinToken.ClientKey, openDouyinToken.OpenId, sinJinritemaiStore, "同步精选联盟达人橱窗数据，插入数据库失败"))
					inTaskLog.SetCreateTime(ctx)

					jsuc.tlrepo.Save(ctx, inTaskLog)
				}
			}

			break
		}
	}

	return stores, err
}
