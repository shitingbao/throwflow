package domain

import (
	"context"
	"fmt"
	"material/internal/pkg/tool"
	"strconv"
	"time"
)

type Material struct {
	Id                 uint64
	VideoId            uint64
	VideoName          string
	VideoUrl           string
	VideoCover         string
	VideoLike          uint64
	VideoLikeShowA     string
	VideoLikeShowB     string
	IndustryId         uint64
	IndustryName       string
	CategoryId         uint64
	CategoryName       string
	Source             string
	AwemeId            uint64
	AwemeName          string
	AwemeAccount       string
	AwemeFollowers     string
	AwemeFollowersShow string
	AwemeImg           string
	AwemeLandingPage   string
	ProductId          uint64
	ProductName        string
	ProductImg         string
	ProductLandingPage string
	ProductPrice       string
	ShopId             string
	ShopName           string
	ShopLogo           string
	ShopScore          string
	IsHot              uint8
	TotalItemNum       uint64
	Platform           string
	PlatformName       string
	IsCollect          uint8
	UpdateDay          time.Time
	CreateTime         time.Time
	UpdateTime         time.Time
}

type MaterialList struct {
	PageNum  uint64
	PageSize uint64
	Total    uint64
	List     []*Material
}

type Product struct {
	ProductId          uint64
	ProductName        string
	ProductImg         string
	ProductLandingPage string
	ProductPrice       string
	IsHot              uint8
	VideoLike          uint64
	VideoLikeShowA     string
	VideoLikeShowB     string
	Awemes             uint64
	Videos             uint64
	Platform           string
	PlatformName       string
}

type ProductList struct {
	PageNum  uint64
	PageSize uint64
	Total    uint64
	List     []*Product
}

type ImageUrls struct {
	ImageUrl string `json:"image_url"`
}

type VideoUrlBodyData struct {
	VideoId   uint64       `json:"video_id"`
	VideoUrl  string       `json:"video_url"`
	ImageUrls []*ImageUrls `json:"image_urls"`
}

type VideoUrlBody struct {
	Code uint64           `json:"code"`
	Msg  string           `json:"msg"`
	Data VideoUrlBodyData `json:"data"`
}

type StatisticsMaterial struct {
	Key   string
	Value string
}

type Promotion struct {
	PromotionId            uint64
	PromotionName          string
	PromotionType          string
	PromotionAccount       string
	PromotionImg           string
	PromotionLandingPage   string
	PromotionFollowers     string
	PromotionFollowersShow string
	PromotionPrice         string
	PromotionPlatformName  string
	IndustryName           string
	Industry               []*MaterialAwemeIndustry
	ShopName               string
	ShopLogo               string
	PageNum                uint64
	PageSize               uint64
	Total                  uint64
	TotalPage              uint64
	List                   []*Material
}

type StatisticsMaterials struct {
	Statistics []*StatisticsMaterial
}

type ChildListCategory struct {
	Key   string
	Value string
}

type Category struct {
	Key       string
	Value     string
	ChildList []*ChildListCategory
}

type Msort struct {
	Key   string
	Value string
}

type Mplatform struct {
	Key   string
	Value string
}

type Search struct {
	Key   string
	Value string
}

type CompanyProductCategory struct {
	IndustryId uint64   `json:"industryId"`
	CategoryId []uint64 `json:"categoryId"`
}

type SelectMaterials struct {
	Category  []*Category
	Msort     []*Msort
	Mplatform []*Mplatform
	Search    []*Search
}

type MaterialAd struct {
	Type    string `json:"type"`
	Message struct {
		Name              string `json:"name"`
		CompanyId         uint64 `json:"companyId"`
		VideoId           uint64 `json:"videoId"`
		CompanyMaterialId uint64 `json:"companyMaterialId"`
		Content           string `json:"content"`
		SendTime          string `json:"sendTime"`
	} `json:"message"`
}

type MaterialAwemeIndustry struct {
	IndustryId    uint64
	IndustryName  string
	IndustryRatio string
	TotalItemNum  uint64
}

func NewSelectMaterials() *SelectMaterials {
	msort := make([]*Msort, 0)

	msort = append(msort, &Msort{Key: "time", Value: "更新时间"})
	msort = append(msort, &Msort{Key: "like", Value: "点赞热度"})
	msort = append(msort, &Msort{Key: "isHot", Value: "正在爆单"})

	mplatform := make([]*Mplatform, 0)

	mplatform = append(mplatform, &Mplatform{Key: "dy", Value: "抖音"})
	mplatform = append(mplatform, &Mplatform{Key: "ks", Value: "快手"})

	search := make([]*Search, 0)

	search = append(search, &Search{Key: "name", Value: "素材"})
	search = append(search, &Search{Key: "product", Value: "商品"})
	search = append(search, &Search{Key: "aweme", Value: "达人"})

	return &SelectMaterials{
		Msort:     msort,
		Mplatform: mplatform,
		Search:    search,
	}
}

func (sm *SelectMaterials) SetCategory(ctx context.Context, category []*Category) {
	sm.Category = category
}

func (m *Material) SetVideoUrl(ctx context.Context, videoUrl string) {
	m.VideoUrl = videoUrl
}

func (m *Material) SetPlatformName(ctx context.Context) {
	if m.Platform == "dy" {
		m.PlatformName = "抖音"
	} else if m.Platform == "ks" {
		m.PlatformName = "快手"
	}
}

func (m *Material) SetAwemeFollowersShow(ctx context.Context) {
	iawemeFollowers, err := strconv.ParseUint(m.AwemeFollowers, 10, 64)

	if err != nil {
		m.AwemeFollowersShow = ""
	}

	if iawemeFollowers > 10000 {
		m.AwemeFollowersShow = fmt.Sprintf("%.1f", float64(iawemeFollowers)/float64(10000)) + "w"
	} else {
		m.AwemeFollowersShow = m.AwemeFollowers
	}
}

func (m *Material) SetVideoLikeShowA(ctx context.Context) {
	if m.VideoLike > 10000 {
		m.VideoLikeShowA = fmt.Sprintf("%.1f", float64(m.VideoLike)/float64(10000)) + "w"
	} else {
		m.VideoLikeShowA = strconv.FormatUint(m.VideoLike, 10)
	}
}

func (m *Material) SetVideoLikeShowB(ctx context.Context) {
	svideoLike := strconv.FormatUint(m.VideoLike, 10)

	slen := len(svideoLike)

	if slen == 0 {
		m.VideoLikeShowB = ""
	} else {
		m.VideoLikeShowB = string(svideoLike[0])

		for slen > 1 {
			m.VideoLikeShowB += "0"

			slen -= 1
		}

		ivideoLikeShowB, _ := strconv.ParseUint(m.VideoLikeShowB, 10, 64)

		if ivideoLikeShowB >= 10000 {
			m.VideoLikeShowB = fmt.Sprintf("%d", ivideoLikeShowB/10000) + "万"
		}
	}
}

func (p *Product) SetVideoLikeShowA(ctx context.Context) {
	if p.VideoLike > 10000 {
		p.VideoLikeShowA = fmt.Sprintf("%.1f", float64(p.VideoLike)/float64(10000)) + "w"
	} else {
		p.VideoLikeShowA = strconv.FormatUint(p.VideoLike, 10)
	}
}

func (p *Product) SetVideoLikeShowB(ctx context.Context) {
	svideoLike := strconv.FormatUint(p.VideoLike, 10)

	slen := len(svideoLike)

	if slen == 0 {
		p.VideoLikeShowB = ""
	} else {
		p.VideoLikeShowB = string(svideoLike[0])

		for slen > 1 {
			p.VideoLikeShowB += "0"

			slen -= 1
		}

		ivideoLikeShowB, _ := strconv.ParseUint(p.VideoLikeShowB, 10, 64)

		if ivideoLikeShowB >= 10000 {
			p.VideoLikeShowB = fmt.Sprintf("%d", ivideoLikeShowB/10000) + "万"
		}
	}
}

func (mai *MaterialAwemeIndustry) SetIndustryRatio(ctx context.Context, totalItemNum uint64) {
	var industryRatio float64

	if totalItemNum > 0 {
		industryRatio = float64(mai.TotalItemNum) / float64(totalItemNum)
	}

	mai.IndustryRatio = strconv.FormatFloat(tool.Decimal(industryRatio*100, 2), 'f', 2, 64) + "%"
}
