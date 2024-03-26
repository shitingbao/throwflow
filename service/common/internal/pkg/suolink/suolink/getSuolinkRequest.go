package suolink

import (
	"github.com/google/go-querystring/query"
)

type GetSuoLinkRequest struct {
	Url        string `url:"url"`
	Format     string `url:"format"`
	Key        string `url:"key"`
	ExpireDate string `url:"expireDate"`
	Domain     string `url:"domain"`
	Protocol   string `url:"protocol"`
}

func (gslr GetSuoLinkRequest) String() string {
	v, _ := query.Values(gslr)

	return v.Encode()
}
