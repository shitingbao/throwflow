package domain

import (
	"context"
	"time"
)

type UserSampleOrder struct {
	Id               uint64
	UserId           uint64
	OpenDouyinUserId uint64
	ProductOutId     uint64
	ProductName      string
	ProductImg       string
	OrderSn          string
	Name             string
	Phone            string
	ProvinceAreaName string
	CityAreaName     string
	AreaAreaName     string
	AddressInfo      string
	Note             string
	IsCancel         uint8
	CancelNote       string
	KuaidiCompany    string
	KuaidiCode       string
	KuaidiNum        string
	KuaidiStateName  string
	CreateTime       time.Time
	UpdateTime       time.Time
}

type ExternalUserSampleOrderList struct {
	PageNum  uint64
	PageSize uint64
	Total    uint64
	List     []*ExternalUserOpenDouyin
}

type CompanyProduct struct {
	Key   string
	Value string
}

type SelectUserSampleOrders struct {
	CompanyProduct []*CompanyProduct
}

type StatisticsUserSampleOrder struct {
	Key   string
	Value string
}

type StatisticsUserSampleOrders struct {
	Statistics []*StatisticsUserSampleOrder
}

func NewSelectUserSampleOrders() *SelectUserSampleOrders {
	return &SelectUserSampleOrders{}
}

func NewUserSampleOrder(ctx context.Context, userId, openDouyinUserId, productOutId uint64, isCancel uint8, productName, productImg, orderSn, name, phone, provinceAreaName, cityAreaName, areaAreaName, addressInfo, note, cancelNote, kuaidiCompany, kuaidiCode, kuaidiNum, kuaidiStateName string) *UserSampleOrder {
	return &UserSampleOrder{
		UserId:           userId,
		OpenDouyinUserId: openDouyinUserId,
		ProductOutId:     productOutId,
		ProductName:      productName,
		ProductImg:       productImg,
		OrderSn:          orderSn,
		Name:             name,
		Phone:            phone,
		ProvinceAreaName: provinceAreaName,
		CityAreaName:     cityAreaName,
		AreaAreaName:     areaAreaName,
		AddressInfo:      addressInfo,
		Note:             note,
		IsCancel:         isCancel,
		CancelNote:       cancelNote,
		KuaidiCompany:    kuaidiCompany,
		KuaidiCode:       kuaidiCode,
		KuaidiNum:        kuaidiNum,
		KuaidiStateName:  kuaidiStateName,
	}
}

func (uso *UserSampleOrder) SetUserId(ctx context.Context, userId uint64) {
	uso.UserId = userId
}

func (uso *UserSampleOrder) SetOpenDouyinUserId(ctx context.Context, openDouyinUserId uint64) {
	uso.OpenDouyinUserId = openDouyinUserId
}

func (uso *UserSampleOrder) SetProductOutId(ctx context.Context, productOutId uint64) {
	uso.ProductOutId = productOutId
}

func (uso *UserSampleOrder) SetProductName(ctx context.Context, productName string) {
	uso.ProductName = productName
}

func (uso *UserSampleOrder) SetProductImg(ctx context.Context, productImg string) {
	uso.ProductImg = productImg
}

func (uso *UserSampleOrder) SetOrderSn(ctx context.Context, orderSn string) {
	uso.OrderSn = orderSn
}

func (uso *UserSampleOrder) SetName(ctx context.Context, name string) {
	uso.Name = name
}

func (uso *UserSampleOrder) SetPhone(ctx context.Context, phone string) {
	uso.Phone = phone
}

func (uso *UserSampleOrder) SetProvinceAreaName(ctx context.Context, provinceAreaName string) {
	uso.ProvinceAreaName = provinceAreaName
}

func (uso *UserSampleOrder) SetCityAreaName(ctx context.Context, cityAreaName string) {
	uso.CityAreaName = cityAreaName
}

func (uso *UserSampleOrder) SetAreaAreaName(ctx context.Context, areaAreaName string) {
	uso.AreaAreaName = areaAreaName
}

func (uso *UserSampleOrder) SetAddressInfo(ctx context.Context, addressInfo string) {
	uso.AddressInfo = addressInfo
}

func (uso *UserSampleOrder) SetNote(ctx context.Context, note string) {
	uso.Note = note
}

func (uso *UserSampleOrder) SetIsCancel(ctx context.Context, isCancel uint8) {
	uso.IsCancel = isCancel
}

func (uso *UserSampleOrder) SetCancelNote(ctx context.Context, cancelNote string) {
	uso.CancelNote = cancelNote
}

func (uso *UserSampleOrder) SetKuaidiCompany(ctx context.Context, kuaidiCompany string) {
	uso.KuaidiCompany = kuaidiCompany
}

func (uso *UserSampleOrder) SetKuaidiCode(ctx context.Context, kuaidiCode string) {
	uso.KuaidiCode = kuaidiCode
}

func (uso *UserSampleOrder) SetKuaidiNum(ctx context.Context, kuaidiNum string) {
	uso.KuaidiNum = kuaidiNum
}

func (uso *UserSampleOrder) SetKuaidiStateName(ctx context.Context, kuaidiStateName string) {
	uso.KuaidiStateName = kuaidiStateName
}

func (uso *UserSampleOrder) SetUpdateTime(ctx context.Context) {
	uso.UpdateTime = time.Now()
}

func (uso *UserSampleOrder) SetCreateTime(ctx context.Context) {
	uso.CreateTime = time.Now()
}

func (suso *SelectUserSampleOrders) SetCompanyProduct(ctx context.Context, companyProduct []*CompanyProduct) {
	suso.CompanyProduct = companyProduct
}
