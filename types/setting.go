package types

import (
	"jdy/enums"

	"github.com/shopspring/decimal"
)

type GoldPriceGetRes struct {
	Price decimal.Decimal `json:"price"` // 金价
}

type GoldPriceListReq struct {
	PageReq
}

type GoldPriceCreateReq struct {
	Price decimal.Decimal `json:"price" required:"true"` // 金价
}

type GoldPriceUpdateReq struct {
	Id     string                `json:"id" required:"true"`     // 金价ID
	Status enums.GoldPriceStatus `json:"status" required:"true"` // 审批状态
}
