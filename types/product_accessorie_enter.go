package types

import (
	"jdy/enums"
	"time"
)

type ProductAccessorieEnterCreateReq struct {
	StoreId string `json:"store_id" binding:"required"` // 门店ID
	Remark  string `json:"remark"`                      // 备注
}

type ProductAccessorieEnterWhere struct {
	Id      string                   `json:"id" label:"入库单号" input:"text" type:"string" find:"true" create:"false" sort:"1" required:"false"`                      // ID
	StoreId string                   `json:"store_id" label:"门店" input:"text" type:"string" find:"false" create:"true" sort:"2" required:"true"`                   // 门店ID
	Status  enums.ProductEnterStatus `json:"status" label:"状态" input:"select" type:"number" find:"true" create:"false" sort:"3" required:"false" preset:"typeMap"` // 状态
	Remark  string                   `json:"remark" label:"备注" input:"text" type:"string" find:"true" create:"true" sort:"4" required:"false"`                     // 类型
	Name    string                   `json:"name" label:"名称" input:"text" type:"string" find:"true" sort:"5" required:"false"`                                     // 名称

	StartTime *time.Time `json:"start_time" label:"开始时间" input:"date" type:"date" find:"true" sort:"6" required:"false"` // 开始时间
	EndTime   *time.Time `json:"end_time" label:"结束时间" input:"date" type:"date" find:"true" sort:"6" required:"false"`   // 结束时间
}

type ProductAccessorieEnterListReq struct {
	PageReq
	Where ProductAccessorieEnterWhere `json:"where"`
}

type ProductAccessorieEnterInfoReq struct {
	PageReq
	Id string `json:"id" binding:"required"`
}

type ProductAccessorieEnterAddProductReq struct {
	EnterId  string                             `json:"enter_id" binding:"required"` // 入库单ID
	Products []ProductAccessorieEnterReqProduct `json:"products" binding:"required"` // 商品信息
}

type ProductAccessorieEnterEditProductReq struct {
	EnterId   string                           `json:"enter_id" binding:"required"`   // 入库单ID
	ProductId string                           `json:"product_id" binding:"required"` // 商品ID
	Product   ProductAccessorieEnterReqProduct `json:"product" binding:"-"`           // 商品信息
}

type ProductAccessorieEnterDelProductReq struct {
	EnterId    string   `json:"enter_id" binding:"required"`    // 入库单ID
	ProductIds []string `json:"product_ids" binding:"required"` // 商品ID列表
}

type ProductAccessorieEnterFinishReq struct {
	EnterId string `json:"enter_id" binding:"required"` // 入库单ID
}

type ProductAccessorieEnterCancelReq struct {
	EnterId string `json:"enter_id" binding:"required"` // 入库单ID
}

type ProductAccessorieEnterReqProduct struct {
	Name       string                            `json:"name" label:"名称" find:"true" create:"true" update:"false" sort:"3" type:"string" input:"text" required:"true" binding:"required"`                             // 名称
	Type       enums.ProductAccessorieType       `json:"type" label:"类型" find:"true" create:"true" update:"false" sort:"4" type:"string" input:"select" required:"true" preset:"typeMap" binding:"required"`          // 类型
	RetailType enums.ProductAccessorieRetailType `json:"retail_type" label:"零售方式" find:"true" create:"true" update:"false" sort:"5" type:"string" input:"select" required:"true" preset:"typeMap" binding:"required"` // 零售类型
	Remark     string                            `json:"remark" label:"备注" find:"true" create:"true" update:"true" sort:"6" type:"string" input:"textarea" required:"false"`                                          // 备注
	Stock      int64                             `json:"stock" label:"库存" find:"false" create:"true" update:"true" sort:"7" type:"int" input:"number" required:"true" binding:"required"`                             // 库存
}
