package types

import (
	"jdy/enums"
	"time"

	"github.com/shopspring/decimal"
)

type ProductAccessorieEnterCreateReq struct {
	StoreId string `json:"store_id" binding:"required"` // 门店ID
	Remark  string `json:"remark"`                      // 备注
}

type ProductAccessorieEnterWhere struct {
	StoreId string                   `json:"store_id" label:"门店" input:"text" type:"string" find:"false" create:"true" sort:"2" required:"true"`                   // 门店ID
	Status  enums.ProductEnterStatus `json:"status" label:"状态" input:"select" type:"number" find:"false" create:"true" sort:"3" required:"false" preset:"typeMap"` // 状态
	Remark  string                   `json:"remark" label:"备注" input:"text" type:"string" find:"false" create:"true" sort:"4" required:"false"`                    // 类型
	Code    string                   `json:"code" label:"条码" input:"text" type:"string" find:"true" sort:"5" required:"false"`                                     // 条码
	StartAt *time.Time               `json:"start_at" label:"开始时间" input:"date" type:"time" find:"true" sort:"6" required:"false"`                                 // 开始时间
	EndAt   *time.Time               `json:"end_at" label:"结束时间" input:"date" type:"time" find:"true" sort:"7" required:"false"`                                   // 结束时间
}

type ProductAccessorieEnterListReq struct {
	PageReq
	Where ProductAccessorieEnterWhere `json:"where"`
}

type ProductAccessorieEnterInfoReq struct {
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
	Code      string          `json:"code" label:"编号" input:"text" type:"string" find:"true" sort:"2" required:"true"` // 条码
	AccessFee decimal.Decimal `json:"access_fee" label:"入网费" input:"number" type:"number" find:"false" create:"true" sort:"3"`
	Stock     int64           `json:"stock" label:"库存" input:"number" type:"number" find:"false" create:"true" sort:"3" required:"true"` // 库存
}
