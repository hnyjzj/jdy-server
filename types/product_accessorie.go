package types

import (
	"github.com/shopspring/decimal"
)

// 配件条件
type ProductAccessorieWhere struct {
	Id        string          `json:"id" label:"ID" find:"true" create:"false" update:"false" sort:"1" type:"string" input:"text" required:"false"`
	StoreId   string          `json:"store_id" label:"门店" find:"true" create:"true" update:"false" sort:"1" type:"string" input:"text" required:"true"`      // 门店
	EntryId   string          `json:"entry_id" label:"入库单" find:"true" create:"true" update:"false" sort:"2" type:"string" input:"text" required:"true"`     // 入库单
	Code      string          `json:"code" label:"类目" find:"true" create:"true" update:"false" sort:"3" type:"string" input:"text" required:"true"`          // 类目
	Name      string          `json:"name" label:"名称" find:"true" create:"true" update:"false" sort:"4" type:"string" input:"text" required:"true"`          // 名称
	Stock     int64           `json:"stock" label:"库存" find:"false" create:"true" update:"false" sort:"5" type:"number" input:"number" required:"true"`      // 库存
	AccessFee decimal.Decimal `json:"access_fee" label:"入网费" find:"false" create:"true" update:"false" sort:"6" type:"float" input:"number" required:"true"` // 入网费
}

// 配件列表
type ProductAccessorieListReq struct {
	PageReq
	Where ProductAccessorieWhere `json:"where" binding:"required"`
}

// 配件详情
type ProductAccessorieInfoReq struct {
	Id string `json:"id" binding:"required"` // 产品ID
}

// 配件更新条件
type ProductAccessorieUpdateReq struct {
	Id string `json:"id" binding:"required"` // 产品ID
	ProductAccessorieWhere
}
