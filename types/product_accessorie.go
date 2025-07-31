package types

import "jdy/enums"

// 配件条件
type ProductAccessorieWhere struct {
	Id         string                            `json:"id" label:"编号" find:"true" create:"false" update:"false" sort:"1" type:"string" input:"text" required:"false"`                             // ID
	StoreId    string                            `json:"store_id" label:"门店" find:"true" create:"false" update:"false" sort:"2" type:"string" input:"text" required:"true"`                        // 门店
	Name       string                            `json:"name" label:"名称" find:"true" create:"true" update:"false" sort:"3" type:"string" input:"text" required:"true"`                             // 名称
	Type       enums.ProductAccessorieType       `json:"type" label:"类型" find:"true" create:"true" update:"false" sort:"4" type:"string" input:"select" required:"true" preset:"typeMap"`          // 类型
	RetailType enums.ProductAccessorieRetailType `json:"retail_type" label:"零售方式" find:"true" create:"true" update:"false" sort:"5" type:"string" input:"select" required:"true" preset:"typeMap"` // 零售类型
	Remark     string                            `json:"remark" label:"备注" find:"true" create:"true" update:"true" sort:"6" type:"string" input:"textarea" required:"false"`                       // 备注
	Stock      int64                             `json:"stock" label:"库存" find:"false" create:"true" update:"false" sort:"7" type:"int" input:"text" required:"true"`                              // 库存
	Status     enums.ProductAccessorieStatus     `json:"status" label:"状态" find:"true" create:"false" update:"false" sort:"8" type:"string" input:"select" required:"true" preset:"typeMap"`       // 状态

	EnterId string `json:"enter_id" label:"入库单" find:"true" create:"false" update:"false" sort:"9" type:"string" input:"text" required:"true"` // 入库单
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
