package video

import (
	"github.com/google/go-querystring/query"
)

type ListVideoRequest struct {
	OpenId string `url:"open_id"`
	Cursor int64  `url:"cursor"`
	Count  int64  `url:"count"`
}

func (lvr ListVideoRequest) String() string {
	v, _ := query.Values(lvr)

	return v.Encode()
}
