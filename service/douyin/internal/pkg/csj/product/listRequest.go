package product

import (
	"encoding/json"
)

type ListDataRequest struct {
	Page        uint64   `json:"page"`          // 商品第⼏⻚，从1开始
	PageSize    uint64   `json:"page_size"`     // 商品每⻚数量，最⼤20，最 ⼩1
	FirstCids   []uint64 `json:"first_cids"`    // 筛选商品⼀级类⽬，从商品类⽬接⼝可获得⼀级类⽬
	SecondCids  []uint64 `json:"second_cids"`   // 筛选商品⼆级类⽬， 从商品类⽬接⼝可获得⼆级类⽬
	ThirdCids   []uint64 `json:"third_cids"`    // 筛选商品三级类⽬，从商品类⽬接⼝可获得三级类⽬
	SearchType  int      `json:"search_type"`   // 排序类型：0 默认排序；1历史销量排序；2价格排序；3佣金排序；4佣金比例排序。不填默认为0。
	OrderType   int      `json:"order_type"`    // 0 升序，1 降序。不填默认为0。若search_type为0，则此值无意义
	CosRatioMin uint64   `json:"cos_ratio_min"` // 分佣⽐例百分⽐乘以100： 1.1%，传1.1*100 = 110
}

func (ldr ListDataRequest) String() string {
	data, _ := json.Marshal(ldr)

	return string(data)
}
