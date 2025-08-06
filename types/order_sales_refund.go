package types

import "time"

type OrderSalesRefundWhere struct {
	OrderId string `json:"order_id" label:"订单号" find:"true" sort:"2" type:"string" input:"search" required:"true"` // 订单号
	StoreId string `json:"store_id" label:"门店" find:"false" sort:"3" type:"string" input:"search" required:"true"` // 门店
	Phone   string `json:"phone" label:"会员(手机号)" find:"true" create:"true" sort:"3" type:"string" input:"text"`    // 会员

	Code string `json:"code" label:"货品条码" find:"true" sort:"5" type:"string" input:"text"` // 货品条码
	Name string `json:"name" label:"货品名称" find:"true" sort:"6" type:"string" input:"text"` // 货品名称

	StartDate *time.Time `json:"start_date" label:"开始日期" find:"true" sort:"9" type:"string" input:"date"` // 开始日期
	EndDate   *time.Time `json:"end_date" label:"结束日期" find:"true" sort:"10" type:"string" input:"date"`  // 结束日期
}

type OrderSalesRefundListReq struct {
	PageReq
	Where OrderSalesRefundWhere `json:"where"`
}
