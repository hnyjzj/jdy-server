package types

import "jdy/enums"

type GoldPriceGetRes struct {
	Price float64 `json:"price"` // 金价
}

type GoldPriceListReq struct {
	PageReq
}

type GoldPriceCreateReq struct {
	Price float64 `json:"price" required:"true"` // 金价
}

type GoldPriceUpdateReq struct {
	Id     string                `json:"id" required:"true"`     // 金价ID
	Status enums.GoldPriceStatus `json:"status" required:"true"` // 审批状态
}
