package types

import (
	"jdy/enums"
	"time"
)

type OrderSalesDetailWhere struct {
	Id string `json:"id" label:"编号" find:"true" sort:"1" type:"string" input:"text"` // 编号

	OrderId  string `json:"order_id" label:"订单" find:"true" sort:"2" type:"string" input:"search" required:"true"`  // 订单
	StoreId  string `json:"store_id" label:"门店" find:"false" sort:"3" type:"string" input:"search" required:"true"` // 门店
	MemberId string `json:"member_id" label:"会员" find:"true" create:"true" sort:"4" type:"string" input:"search"`   // 会员

	Type   enums.ProductType      `json:"type" label:"产品类型" find:"true" sort:"5" type:"number" input:"select" preset:"typeMap"`   // 产品类型
	Status enums.OrderSalesStatus `json:"status" label:"产品状态" find:"true" sort:"6" type:"number" input:"select" preset:"typeMap"` // 产品状态

	Code string `json:"code" label:"产品编号" find:"true" sort:"7" type:"string" input:"text"` // 产品编号

	StartDate *time.Time `json:"start_date" label:"开始日期" find:"true" sort:"8" type:"string" input:"date"` // 开始日期
	EndDate   *time.Time `json:"end_date" label:"结束日期" find:"true" sort:"9" type:"string" input:"date"`   // 结束日期
}

type OrderSalesDetailListReq struct {
	PageReq
	Where OrderSalesDetailWhere `json:"where"`
}

type OrderSalesDetailInfoReq struct {
	Id string `json:"id" required:"true"`
}
