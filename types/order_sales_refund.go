package types

import "time"

type OrderSalesRefundWhere struct {
	Id string `json:"id" label:"编号" find:"true" sort:"1" type:"string" input:"text"` // 编号

	OrderId  string `json:"order_id" label:"订单" find:"true" sort:"2" type:"string" input:"search" required:"true"`  // 订单
	StoreId  string `json:"store_id" label:"门店" find:"false" sort:"3" type:"string" input:"search" required:"true"` // 门店
	MemberId string `json:"member_id" label:"会员" find:"true" create:"true" sort:"4" type:"string" input:"search"`   // 会员

	Code string `json:"code" label:"产品编号" find:"true" sort:"5" type:"string" input:"text"` // 产品编号
	Name string `json:"name" label:"产品名称" find:"true" sort:"6" type:"string" input:"text"` // 产品名称

	StartDate *time.Time `json:"start_date" label:"开始日期" find:"true" sort:"9" type:"string" input:"date"` // 开始日期
	EndDate   *time.Time `json:"end_date" label:"结束日期" find:"true" sort:"10" type:"string" input:"date"`  // 结束日期
}

type OrderSalesRefundListReq struct {
	PageReq
	Where OrderSalesRefundWhere `json:"where"`
}
