package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"strconv"
	v1 "weixin/api/service/common/v1"
	"weixin/internal/conf"
	"weixin/internal/domain"
)

var (
	WeixinUserSampleOrderNotFound             = errors.NotFound("WEIXIN_USER_SAMPLE_ORDER_NOT_FOUND", "微信用户寄样申请不存在")
	WeixinUserSampleOrderListProductError     = errors.InternalServer("WEIXIN_USER_SAMPLE_ORDER_LIST_PRODUCT_ERROR", "微信用户寄样申请商品列表获取失败")
	WeixinUserSampleOrderListError            = errors.InternalServer("WEIXIN_USER_SAMPLE_ORDER_LIST_ERROR", "微信用户寄样申请列表获取失败")
	WeixinUserSampleOrderCreateError          = errors.InternalServer("WEIXIN_USER_SAMPLE_ORDER_CREATE_ERROR", "微信用户寄样申请创建失败")
	WeixinUserSampleOrderUpdateError          = errors.InternalServer("WEIXIN_USER_SAMPLE_ORDER_UPDATE_ERROR", "微信用户寄样申请更新失败")
	WeixinUserSampleOrderNotMeetRequirement   = errors.InternalServer("WEIXIN_USER_SAMPLE_ORDER_NOT_MEET_REQUIREMENT", "不满足免费寄样申请要求")
	WeixinUserSampleOrderNotAcceptRequirement = errors.InternalServer("WEIXIN_USER_SAMPLE_ORDER_NOT_MEET_REQUIREMENT", "不接受免费寄样")
)

type UserSampleOrderRepo interface {
	NextId(context.Context) (uint64, error)
	Get(context.Context, uint64) (*domain.UserSampleOrder, error)
	List(context.Context, int, int, uint64) ([]*domain.UserSampleOrder, error)
	ListProduct(context.Context, string) ([]*domain.UserSampleOrder, error)
	ListOpenDouyinUser(context.Context, int, int, string, string, string) ([]*domain.UserOpenDouyin, error)
	ListByOpenDouyinUserIds(context.Context, []uint64, string, string, string) ([]*domain.UserSampleOrder, error)
	Count(context.Context, uint64) (int64, error)
	CountOpenDouyinUser(context.Context, string, string, string) (int64, error)
	Statistics(context.Context, string, string, string) (int64, error)
	Save(context.Context, *domain.UserSampleOrder) (*domain.UserSampleOrder, error)
	Update(context.Context, *domain.UserSampleOrder) (*domain.UserSampleOrder, error)
}

type UserSampleOrderUsecase struct {
	repo    UserSampleOrderRepo
	uodrepo UserOpenDouyinRepo
	uarepo  UserAddressRepo
	crepo   CompanyRepo
	cprepo  CompanyProductRepo
	jorepo  JinritemaiOrderRepo
	kir     KuaidiInfoRepo
	arepo   AreaRepo
	tm      Transaction
	conf    *conf.Data
	log     *log.Helper
}

func NewUserSampleOrderUsecase(repo UserSampleOrderRepo, uodrepo UserOpenDouyinRepo, uarepo UserAddressRepo, crepo CompanyRepo, cprepo CompanyProductRepo, jorepo JinritemaiOrderRepo, kir KuaidiInfoRepo, arepo AreaRepo, tm Transaction, conf *conf.Data, logger log.Logger) *UserSampleOrderUsecase {
	return &UserSampleOrderUsecase{repo: repo, uodrepo: uodrepo, uarepo: uarepo, crepo: crepo, cprepo: cprepo, jorepo: jorepo, kir: kir, arepo: arepo, tm: tm, conf: conf, log: log.NewHelper(logger)}
}

func (usouc *UserSampleOrderUsecase) GetKuaidiInfoUserSampleOrders(ctx context.Context, userSampleOrderId uint64) (*v1.GetKuaidiInfosReply, error) {
	inUserSampleOrder, err := usouc.repo.Get(ctx, userSampleOrderId)

	if err != nil {
		return nil, WeixinUserSampleOrderNotFound
	}

	if len(inUserSampleOrder.KuaidiNum) == 0 || len(inUserSampleOrder.KuaidiCode) == 0 {
		return nil, WeixinKuaidiInfoGetError
	}

	kuaidiInfo, err := usouc.kir.Get(ctx, inUserSampleOrder.KuaidiCode, inUserSampleOrder.KuaidiNum, inUserSampleOrder.Phone)

	if err != nil {
		return nil, WeixinKuaidiInfoGetError
	}

	inUserSampleOrder.SetKuaidiStateName(ctx, kuaidiInfo.Data.StateName)
	inUserSampleOrder.SetUpdateTime(ctx)

	if _, err := usouc.repo.Update(ctx, inUserSampleOrder); err != nil {
		return nil, WeixinKuaidiInfoGetError
	}

	return kuaidiInfo, nil
}

func (usouc *UserSampleOrderUsecase) ListUserSampleOrders(ctx context.Context, pageNum, pageSize uint64, day, keyword, searchType string) (*domain.ExternalUserSampleOrderList, error) {
	openDouyinUsers, err := usouc.repo.ListOpenDouyinUser(ctx, int(pageNum), int(pageSize), day, keyword, searchType)

	if err != nil {
		return nil, WeixinUserSampleOrderListError
	}

	total, err := usouc.repo.CountOpenDouyinUser(ctx, day, keyword, searchType)

	if err != nil {
		return nil, WeixinUserSampleOrderListError
	}

	openDouyinUserIds := make([]uint64, 0)

	for _, openDouyinUser := range openDouyinUsers {
		openDouyinUserIds = append(openDouyinUserIds, openDouyinUser.Id)
	}

	userSampleOrders, err := usouc.repo.ListByOpenDouyinUserIds(ctx, openDouyinUserIds, day, keyword, searchType)

	if err != nil {
		return nil, WeixinUserSampleOrderListError
	}

	list := make([]*domain.ExternalUserOpenDouyin, 0)

	for _, openDouyinUser := range openDouyinUsers {
		openDouyinUser.SetFansShow(ctx)

		list = append(list, &domain.ExternalUserOpenDouyin{
			Id:           openDouyinUser.Id,
			Nickname:     openDouyinUser.Nickname,
			AccountId:    openDouyinUser.AccountId,
			Avatar:       openDouyinUser.Avatar,
			AvatarLarger: openDouyinUser.AvatarLarger,
			Fans:         openDouyinUser.Fans,
			FansShow:     openDouyinUser.FansShow,
		})
	}

	for _, userSampleOrder := range userSampleOrders {
		for _, l := range list {
			if userSampleOrder.OpenDouyinUserId == l.Id {
				l.UserSampleOrders = append(l.UserSampleOrders, userSampleOrder)

				break
			}
		}
	}

	return &domain.ExternalUserSampleOrderList{
		PageNum:  pageNum,
		PageSize: pageSize,
		Total:    uint64(total),
		List:     list,
	}, nil
}

func (usouc *UserSampleOrderUsecase) StatisticsUserSampleOrders(ctx context.Context, day, keyword, searchType string) (*domain.StatisticsUserSampleOrders, error) {
	statistics := make([]*domain.StatisticsUserSampleOrder, 0)

	count, _ := usouc.repo.Statistics(ctx, day, keyword, searchType)

	statistics = append(statistics, &domain.StatisticsUserSampleOrder{
		Key:   "所有",
		Value: strconv.FormatInt(count, 10),
	})

	return &domain.StatisticsUserSampleOrders{
		Statistics: statistics,
	}, nil
}

func (usouc *UserSampleOrderUsecase) VerifyUserSampleOrders(ctx context.Context, userId, openDouyinUserId, productId uint64) error {
	/*openDouyinUsers, err := usouc.uodrepo.List(ctx, 0, 40, userId, "")

	if err != nil {
		return WeixinUserOpenDouyinListError
	}

	isNotExist := true

	clientKey := ""
	openId := ""

	for _, openDouyinUser := range openDouyinUsers {
		if openDouyinUser.Id == openDouyinUserId {
			isNotExist = false

			clientKey = openDouyinUser.ClientKey
			openId = openDouyinUser.OpenId

			break
		}
	}

	if isNotExist {
		return WeixinUserOpenDouyinNotFound
	}

	companyProduct, err := usouc.cprepo.GetExternal(ctx, productId)

	if err != nil {
		return WeixinCompanyProductNotFound
	}

	startDay := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	endDay := time.Now().Format("2006-01-02")

	statisticses, err := usouc.jorepo.StatisticsByClientKeyAndOpenId(ctx, clientKey, openId, startDay, endDay)

	if err != nil {
		return WeixinJinritemaiOrderGetError
	}

	if companyProduct.Data.SampleThresholdType == 1 {
		for _, statistics := range statisticses.Data.Statistics {
			if statistics.Key == "销量" {
				itemNum, err := strconv.ParseUint(statistics.Value, 10, 64)

				if err != nil {
					return WeixinUserSampleOrderNotMeetRequirement
				}

				if itemNum < companyProduct.Data.SampleThresholdValue {
					return WeixinUserSampleOrderNotMeetRequirement
				}
			}
		}
	} else if companyProduct.Data.SampleThresholdType == 2 {
		for _, statistics := range statisticses.Data.Statistics {
			if statistics.Key == "销额" {
				totalPayAmount, err := strconv.ParseFloat(statistics.Value, 10)

				if err != nil {
					return WeixinUserSampleOrderNotMeetRequirement
				}

				if totalPayAmount < float64(companyProduct.Data.SampleThresholdValue) {
					return WeixinUserSampleOrderNotMeetRequirement
				}
			}
		}
	} else if companyProduct.Data.SampleThresholdType == 3 {
		return WeixinUserSampleOrderNotAcceptRequirement
	}
	*/
	return nil
}

func (usouc *UserSampleOrderUsecase) CancelUserSampleOrders(ctx context.Context, userSampleOrderId uint64, cancelNote string) error {
	inUserSampleOrder, err := usouc.repo.Get(ctx, userSampleOrderId)

	if err != nil {
		return WeixinUserSampleOrderNotFound
	}

	inUserSampleOrder.SetIsCancel(ctx, 1)
	inUserSampleOrder.SetCancelNote(ctx, cancelNote)
	inUserSampleOrder.SetUpdateTime(ctx)

	if _, err := usouc.repo.Update(ctx, inUserSampleOrder); err != nil {
		return WeixinUserSampleOrderUpdateError
	}

	return nil
}

func (usouc *UserSampleOrderUsecase) CreateUserSampleOrders(ctx context.Context, userId, openDouyinUserId, productId, userAddressId uint64, note string) (*domain.UserSampleOrder, error) {
	/*openDouyinUsers, err := usouc.uodrepo.List(ctx, 0, 40, userId, "")

	if err != nil {
		return nil, WeixinUserOpenDouyinListError
	}

	isNotExist := true

	clientKey := ""
	openId := ""

	for _, openDouyinUser := range openDouyinUsers {
		if openDouyinUser.Id == openDouyinUserId {
			isNotExist = false

			clientKey = openDouyinUser.ClientKey
			openId = openDouyinUser.OpenId

			break
		}
	}

	if isNotExist {
		return nil, WeixinUserOpenDouyinNotFound
	}

	userAddress, err := usouc.uarepo.GetById(ctx, userId, userAddressId)

	if err != nil {
		return nil, WeixinUserAddressNotFound
	}

	companyProduct, err := usouc.cprepo.GetExternal(ctx, productId)

	if err != nil {
		return nil, WeixinCompanyProductNotFound
	}

	startDay := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	endDay := time.Now().Format("2006-01-02")

	statisticses, err := usouc.jorepo.StatisticsByClientKeyAndOpenId(ctx, clientKey, openId, startDay, endDay)

	if err != nil {
		return nil, WeixinJinritemaiOrderGetError
	}

	if companyProduct.Data.SampleThresholdType == 1 {
		for _, statistics := range statisticses.Data.Statistics {
			if statistics.Key == "销量" {
				itemNum, err := strconv.ParseUint(statistics.Value, 10, 64)

				if err != nil {
					return nil, WeixinUserSampleOrderNotMeetRequirement
				}

				if itemNum < companyProduct.Data.SampleThresholdValue {
					return nil, WeixinUserSampleOrderNotMeetRequirement
				}
			}
		}
	} else if companyProduct.Data.SampleThresholdType == 2 {
		for _, statistics := range statisticses.Data.Statistics {
			if statistics.Key == "销额" {
				totalPayAmount, err := strconv.ParseFloat(statistics.Value, 10)

				if err != nil {
					return nil, WeixinUserSampleOrderNotMeetRequirement
				}

				if totalPayAmount < float64(companyProduct.Data.SampleThresholdValue) {
					return nil, WeixinUserSampleOrderNotMeetRequirement
				}
			}
		}
	} else if companyProduct.Data.SampleThresholdType == 3 {
		return nil, WeixinUserSampleOrderNotAcceptRequirement
	}

	province, err := usouc.arepo.GetByAreaCode(ctx, userAddress.ProvinceAreaCode)

	if err != nil {
		return nil, WeixinAreaNotFound
	}

	city, err := usouc.arepo.GetByAreaCode(ctx, userAddress.CityAreaCode)

	if err != nil {
		return nil, WeixinAreaNotFound
	}

	if province.Data.AreaCode != city.Data.ParentAreaCode {
		return nil, WeixinAreaNotFound
	}

	area, err := usouc.arepo.GetByAreaCode(ctx, userAddress.AreaAreaCode)

	if err != nil {
		return nil, WeixinAreaNotFound
	}

	orderSn, err := usouc.repo.NextId(ctx)

	if err != nil {
		return nil, WeixinUserSampleOrderCreateError
	}

	productImg := ""

	if len(companyProduct.Data.ProductImgs) > 0 {
		productImg = companyProduct.Data.ProductImgs[0].ProductImg
	}

	inUserSampleOrder := domain.NewUserSampleOrder(ctx, userId, openDouyinUserId, companyProduct.Data.ProductId, 0, companyProduct.Data.ProductName, productImg, strconv.FormatUint(orderSn, 10), userAddress.Name, userAddress.Phone, province.Data.AreaName, city.Data.AreaName, area.Data.AreaName, userAddress.AddressInfo, note, "", "", "", "", "")
	inUserSampleOrder.SetCreateTime(ctx)
	inUserSampleOrder.SetUpdateTime(ctx)

	userSampleOrder, err := usouc.repo.Save(ctx, inUserSampleOrder)

	if err != nil {
		return nil, WeixinUserSampleOrderCreateError
	}*/

	return nil, nil
}

func (usouc *UserSampleOrderUsecase) UpdateKuaidiUserSampleOrders(ctx context.Context, userSampleOrderId uint64, kuaidiCode, kuaidiNum string) error {
	inUserSampleOrder, err := usouc.repo.Get(ctx, userSampleOrderId)

	if err != nil {
		return WeixinUserSampleOrderNotFound
	}

	if kuaidiInfo, err := usouc.kir.Get(ctx, kuaidiCode, kuaidiNum, inUserSampleOrder.Phone); err == nil {
		inUserSampleOrder.SetKuaidiCompany(ctx, kuaidiInfo.Data.Name)
		inUserSampleOrder.SetKuaidiStateName(ctx, kuaidiInfo.Data.StateName)
	}

	inUserSampleOrder.SetKuaidiCode(ctx, kuaidiCode)
	inUserSampleOrder.SetKuaidiNum(ctx, kuaidiNum)
	inUserSampleOrder.SetUpdateTime(ctx)

	if _, err := usouc.repo.Update(ctx, inUserSampleOrder); err != nil {
		return WeixinUserSampleOrderUpdateError
	}

	return nil
}
