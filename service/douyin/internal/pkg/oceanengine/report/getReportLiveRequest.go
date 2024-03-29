package report

import (
	"github.com/google/go-querystring/query"
)

type GetReportLiveRequest struct {
	AdvertiserId uint64 `url:"advertiser_id"`
	AwemeId      uint64 `url:"aweme_id"`
	StartTime    string `url:"start_time"`
	EndTime      string `url:"end_time"`
	Fields       string `url:"fields"`
}

func (grlr GetReportLiveRequest) String() string {
	v, _ := query.Values(grlr)

	return v.Encode()
}
