package product

import (
	"encoding/json"
)

type DetailDataRequest struct {
	ProductIds []uint64 `json:"product_ids"`
}

func (ddr DetailDataRequest) String() string {
	data, _ := json.Marshal(ddr)

	return string(data)
}
