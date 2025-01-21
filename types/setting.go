package types

import "jdy/enums"

type GoldPriceCreateReq struct {
	Price float64 `json:"price" required:"true"` // 金价
}

type GoldPriceUpdateReq struct {
	Id     int64                 `json:"id" required:"true"`     // 金价ID
	Status enums.GoldPriceStatus `json:"status" required:"true"` // 审批状态
}
