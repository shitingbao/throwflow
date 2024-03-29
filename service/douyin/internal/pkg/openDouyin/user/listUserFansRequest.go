package user

import (
	"github.com/google/go-querystring/query"
)

type ListUserFansRequest struct {
	OpenId   string `url:"open_id"`
	DateType int64  `url:"date_type"`
}

func (lufr ListUserFansRequest) String() string {
	v, _ := query.Values(lufr)

	return v.Encode()
}
