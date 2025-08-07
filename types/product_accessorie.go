package types

import (
	"jdy/enums"

	"github.com/shopspring/decimal"
)

// 配件条件
type ProductAccessorieWhere struct {
	Id         string                            `json:"id" label:"编号" find:"true" create:"false" update:"false" info:"true" sort:"1" type:"string" input:"text" required:"false"`                             // ID
	StoreId    string                            `json:"store_id" label:"门店" find:"false" create:"false" update:"false" sort:"2" type:"string" input:"text" required:"true"`                                   // 门店
	Store      string                            `json:"store" label:"门店名称" info:"true" sort:"3" type:"string" input:"text" required:"true"`                                                                   // 门店名称
	Name       string                            `json:"name" label:"名称" find:"true" create:"true" update:"false" info:"true" sort:"4" type:"string" input:"text" required:"true"`                             // 名称
	Stock      int64                             `json:"stock" label:"库存" find:"false" create:"true" update:"false" info:"true" sort:"5" type:"number" input:"number" required:"true"`                         // 库存
	Type       enums.ProductAccessorieType       `json:"type" label:"类型" find:"true" create:"true" update:"false" info:"true" sort:"6" type:"number" input:"select" required:"true" preset:"typeMap"`          // 类型
	RetailType enums.ProductAccessorieRetailType `json:"retail_type" label:"零售方式" find:"true" create:"true" update:"false" info:"true" sort:"7" type:"number" input:"select" required:"true" preset:"typeMap"` // 零售类型
	Remark     string                            `json:"remark" label:"备注" find:"true" create:"true" update:"true" info:"true" sort:"8" type:"string" input:"textarea" required:"false"`                       // 备注
	Price      decimal.Decimal                   `json:"price" label:"单价" find:"false" create:"true" update:"true" info:"true" sort:"9" type:"decimal" input:"number" required:"true"`                         // 单价
	Status     enums.ProductAccessorieStatus     `json:"status" label:"状态" find:"true" create:"false" update:"false" info:"true" sort:"10" type:"number" input:"select" required:"true" preset:"typeMap"`      // 状态

	EnterId string `json:"enter_id" label:"入库单" find:"false" create:"false" update:"false" sort:"11" type:"string" input:"text" required:"true"` // 入库单

	CreatedAt string `json:"created_at" label:"创建时间"  info:"true" sort:"12" type:"date"` // 创建时间
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
