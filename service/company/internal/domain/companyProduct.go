package domain

import (
	"company/internal/pkg/tool"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Commission struct {
	CommissionRatio  float32 `json:"commissionRatio"`
	ServiceRatio     float32 `json:"serviceRatio"`
	CommissionOutUrl string  `json:"commissionOutUrl"`
}

type CompanyProduct struct {
	Id                    uint64
	ProductOutId          uint64
	ProductType           uint8
	ProductStatus         uint8
	ProductName           string
	ProductImg            string
	ProductImgs           []string
	ProductDetailImg      string
	ProductDetailImgs     []string
	ProductPrice          string
	IndustryId            uint64
	IndustryName          string
	CategoryId            uint64
	CategoryName          string
	SubCategoryId         uint64
	SubCategoryName       string
	ProductUrl            string
	ShopName              string
	ShopLogo              string
	ShopScore             float64
	IsTop                 uint8
	IsHot                 uint8
	IsExist               uint8
	TotalSale             uint64
	CommissionRatio       float32
	SampleThresholdType   uint8
	SampleThresholdValue  uint64
	MaterialOutUrl        string
	MaterialOutUrls       []string
	Commission            string
	Commissions           []*Commission
	InvestmentRatio       float32
	ForbidReason          string
	IsTask                uint8
	PureCommission        string
	PureServiceCommission string
	Awemes                []*Aweme
	CreateTime            time.Time
	UpdateTime            time.Time
}

type CompanyProductList struct {
	PageNum  uint64
	PageSize uint64
	Total    uint64
	List     []*CompanyProduct
}

type StatisticsCompanyProduct struct {
	Key   string
	Value string
}

type CompanyProductCreateJinritemaiStoreMessage struct {
	ProductName string
	AwemeName   string
	Content     string
}

type CompanyProductCreateJinritemaiStore struct {
	Content  string
	Messages []*CompanyProductCreateJinritemaiStoreMessage
}

type Aweme struct {
	Nickname     string
	AccountId    string
	Avatar       string
	AvatarLarger string
}

type ProductAd struct {
	Type    string `json:"type"`
	Message struct {
		Name      string `json:"name"`
		ProductId uint64 `json:"productId"`
		Content   string `json:"content"`
		SendTime  string `json:"sendTime"`
	} `json:"message"`
}

func NewCompanyProduct(ctx context.Context, productOutId, industryId, categoryId, subCategoryId uint64, productType, productStatus, isTop, isExist uint8, investmentRatio float32, productName, productImg, productPrice, materialOutUrl, commission, forbidReason string) *CompanyProduct {
	return &CompanyProduct{
		ProductOutId:    productOutId,
		ProductType:     productType,
		ProductStatus:   productStatus,
		ProductName:     productName,
		ProductImg:      productImg,
		ProductPrice:    productPrice,
		IndustryId:      industryId,
		CategoryId:      categoryId,
		SubCategoryId:   subCategoryId,
		IsTop:           isTop,
		IsExist:         isExist,
		MaterialOutUrl:  materialOutUrl,
		Commission:      commission,
		InvestmentRatio: investmentRatio,
		ForbidReason:    forbidReason,
	}
}

func (cp *CompanyProduct) SetProductOutId(ctx context.Context, productOutId uint64) {
	cp.ProductOutId = productOutId
}

func (cp *CompanyProduct) SetProductType(ctx context.Context, productType uint8) {
	cp.ProductType = productType
}

func (cp *CompanyProduct) SetProductStatus(ctx context.Context, productStatus uint8) {
	cp.ProductStatus = productStatus
}

func (cp *CompanyProduct) SetProductName(ctx context.Context, productName string) {
	cp.ProductName = productName
}

func (cp *CompanyProduct) SetProductImg(ctx context.Context, productImg string) {
	cp.ProductImg = productImg
}

func (cp *CompanyProduct) SetProductImgs(ctx context.Context) {
	cp.ProductImgs = tool.RemoveEmptyString(strings.Split(cp.ProductImg, ","))
}

func (cp *CompanyProduct) SetProductDetailImg(ctx context.Context, productDetailImg string) {
	cp.ProductDetailImg = productDetailImg
}

func (cp *CompanyProduct) SetProductDetailImgs(ctx context.Context) {
	cp.ProductDetailImgs = tool.RemoveEmptyString(strings.Split(cp.ProductDetailImg, ","))
}

func (cp *CompanyProduct) SetProductPrice(ctx context.Context, productPrice string) {
	cp.ProductPrice = productPrice
}

func (cp *CompanyProduct) SetIndustryId(ctx context.Context, industryId uint64) {
	cp.IndustryId = industryId
}

func (cp *CompanyProduct) SetIndustryName(ctx context.Context, industryName string) {
	cp.IndustryName = industryName
}

func (cp *CompanyProduct) SetCategoryId(ctx context.Context, categoryId uint64) {
	cp.CategoryId = categoryId
}

func (cp *CompanyProduct) SetCategoryName(ctx context.Context, categoryName string) {
	cp.CategoryName = categoryName
}

func (cp *CompanyProduct) SetSubCategoryId(ctx context.Context, subCategoryId uint64) {
	cp.SubCategoryId = subCategoryId
}

func (cp *CompanyProduct) SetSubCategoryName(ctx context.Context, subCategoryName string) {
	cp.SubCategoryName = subCategoryName
}

func (cp *CompanyProduct) SetProductUrl(ctx context.Context) {
	cp.ProductUrl = fmt.Sprintf("https://haohuo.jinritemai.com/ecommerce/trade/detail/index.html?origin_type=old_h5&id=%d", cp.ProductOutId)
}

func (cp *CompanyProduct) SetShopName(ctx context.Context, shopName string) {
	cp.ShopName = shopName
}

func (cp *CompanyProduct) SetShopLogo(ctx context.Context, shopLogo string) {
	cp.ShopLogo = shopLogo
}

func (cp *CompanyProduct) SetShopScore(ctx context.Context, shopScore float64) {
	cp.ShopScore = shopScore
}

func (cp *CompanyProduct) SetIsTop(ctx context.Context, isTop uint8) {
	cp.IsTop = isTop
}

func (cp *CompanyProduct) SetIsHot(ctx context.Context, isHot uint8) {
	cp.IsHot = isHot
}

func (cp *CompanyProduct) SetIsExist(ctx context.Context, isExist uint8) {
	cp.IsExist = isExist
}

func (cp *CompanyProduct) SetTotalSale(ctx context.Context, totalSale uint64) {
	cp.TotalSale = totalSale
}

func (cp *CompanyProduct) SetCommissionRatio(ctx context.Context, commissionRatio float32) {
	cp.CommissionRatio = commissionRatio
}

func (cp *CompanyProduct) SetSampleThresholdType(ctx context.Context, sampleThresholdType uint8) {
	cp.SampleThresholdType = sampleThresholdType
}

func (cp *CompanyProduct) SetSampleThresholdValue(ctx context.Context, sampleThresholdValue uint64) {
	cp.SampleThresholdValue = sampleThresholdValue
}

func (cp *CompanyProduct) SetMaterialOutUrl(ctx context.Context, materialOutUrl string) {
	cp.MaterialOutUrl = materialOutUrl
}

func (cp *CompanyProduct) SetMaterialOutUrls(ctx context.Context) {
	materialOutUrls := strings.Split(cp.MaterialOutUrl, ",")

	cp.MaterialOutUrls = tool.RemoveDuplicateString(materialOutUrls)
	cp.MaterialOutUrl = strings.Join(cp.MaterialOutUrls, ",")
}

func (cp *CompanyProduct) GetMaterialOutUrls(ctx context.Context) []string {
	materialOutUrls := make([]string, 0)

	if len(cp.MaterialOutUrl) > 0 {
		materialOutUrls = strings.Split(cp.MaterialOutUrl, ",")
	}

	return tool.RemoveDuplicateString(materialOutUrls)
}

func (cp *CompanyProduct) SetCommission(ctx context.Context, commission string) {
	cp.Commission = commission
}

func (cp *CompanyProduct) SetCommissions(ctx context.Context) {
	var commissions []*Commission

	if err := json.Unmarshal([]byte(cp.Commission), &commissions); err == nil {
		for _, commission := range commissions {
			cp.Commissions = append(cp.Commissions, commission)
		}
	}
}

func (cp *CompanyProduct) GetCommission(ctx context.Context) (pureCommission, pureServiceCommission, pureCommissionUrl string) {
	var fpureCommission float32 = 100.00
	var fpureServiceCommission float32 = 0.00

	isExistPureCommission := false

	for _, companyProductCommission := range cp.Commissions {
		isExistPureCommission = true

		if fpureCommission > companyProductCommission.CommissionRatio {
			fpureCommission = companyProductCommission.CommissionRatio
			fpureServiceCommission = companyProductCommission.ServiceRatio
			pureCommissionUrl = companyProductCommission.CommissionOutUrl
		}
	}

	if isExistPureCommission {
		pureCommission = fmt.Sprintf("%.f", tool.Decimal(float64(fpureCommission), 0))
		pureServiceCommission = fmt.Sprintf("%.f", tool.Decimal(float64(fpureServiceCommission), 0))
	}

	return
}

func (cp *CompanyProduct) SetInvestmentRatio(ctx context.Context, investmentRatio float32) {
	cp.InvestmentRatio = investmentRatio
}

func (cp *CompanyProduct) SetForbidReason(ctx context.Context, forbidReason string) {
	cp.ForbidReason = forbidReason
}

func (cp *CompanyProduct) SetPureCommission(ctx context.Context, pureCommission string) {
	cp.PureCommission = pureCommission
}

func (cp *CompanyProduct) SetPureServiceCommission(ctx context.Context, pureServiceCommission string) {
	cp.PureServiceCommission = pureServiceCommission
}

func (cp *CompanyProduct) SetAwemes(ctx context.Context, awemes []*Aweme) {
	cp.Awemes = awemes
}

func (cp *CompanyProduct) SetUpdateTime(ctx context.Context, time time.Time) {
	cp.UpdateTime = time
}

func (cp *CompanyProduct) SetCreateTime(ctx context.Context) {
	cp.CreateTime = time.Now()
}

func (cp *CompanyProduct) SetIsTask(ctx context.Context, isTask uint8) {
	cp.IsTask = isTask
}

func (cp *CompanyProduct) VerifyCommission(ctx context.Context) bool {
	commissionPrefixs := [1]string{"^https://haohuo.jinritemai.com"}

	var commissions []*Commission

	if err := json.Unmarshal([]byte(cp.Commission), &commissions); err != nil {
		return false
	}

	commissionPrefix := ""
	var productOutId uint64
	originType := "mengma"

	for _, commission := range commissions {
		if commission.CommissionRatio <= 0 {
			return false
		}

		if commission.CommissionRatio > 100 {
			return false
		}

		if commission.ServiceRatio < 0 {
			return false
		}

		if commission.ServiceRatio > 100 {
			return false
		}

		if len(commissionPrefix) == 0 {
			for _, lcommissionPrefix := range commissionPrefixs {
				if regexp.MustCompile(lcommissionPrefix).MatchString(commission.CommissionOutUrl) {
					commissionPrefix = lcommissionPrefix
				}
			}
		}

		if len(commissionPrefix) == 0 {
			return false
		}

		if !regexp.MustCompile(commissionPrefix).MatchString(commission.CommissionOutUrl) {
			return false
		}

		commissionOutUrl, err := url.Parse(commission.CommissionOutUrl)

		if err != nil {
			return false
		}

		sproductOutId := commissionOutUrl.Query().Get("id")

		iproductOutId, err := strconv.ParseUint(sproductOutId, 10, 64)

		if err != nil {
			return false
		}

		if productOutId == 0 {
			productOutId = iproductOutId
		} else {
			if iproductOutId != productOutId {
				return false
			}
		}

		if originType == "mengma" {
			originType = commissionOutUrl.Query().Get("origin_type")
		} else {
			if originType != commissionOutUrl.Query().Get("origin_type") {
				return false
			}
		}

		cp.Commissions = append(cp.Commissions, commission)
	}

	return true
}

func (cp *CompanyProduct) VerifyMaterialOutUrl(ctx context.Context) bool {
	materialOutUrlPrefixs := [2]string{"^https://pcp.kuaizi.co", "https://v.douyin.com"}

	for _, materialOutUrl := range cp.MaterialOutUrls {
		isNotExist := true

		for _, materialOutUrlPrefix := range materialOutUrlPrefixs {
			if regexp.MustCompile(materialOutUrlPrefix).MatchString(materialOutUrl) {
				isNotExist = false

				break
			}
		}

		if isNotExist {
			return false
		}
	}

	return true
}
