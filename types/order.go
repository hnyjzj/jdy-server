package types

import (
	"errors"
	"jdy/enums"
	"time"
)

type OrderWhere struct {
	Id       string `json:"id" label:"订单编号" show:"true" sort:"1" type:"string" input:"text"`        // 订单编号
	StoreId  string `json:"store_id" label:"门店ID" show:"true" sort:"2" type:"string" input:"text"`  // 门店ID
	MemberId string `json:"member_id" label:"会员ID" show:"true" sort:"3" type:"string" input:"text"` // 会员ID

	Status enums.OrderStatus `json:"status" label:"订单状态" show:"true" sort:"4" type:"string" input:"select" preset:"typeMap"` // 订单状态
	Type   enums.OrderType   `json:"type" label:"订单类型" show:"true" sort:"5" type:"string" input:"select" preset:"typeMap"`   // 订单类型
	Source enums.OrderSource `json:"source" label:"订单来源" show:"true" sort:"6" type:"string" input:"select" preset:"typeMap"` // 订单来源

	CashierId   string `json:"cashier_id" label:"收银员ID" show:"true" sort:"7" type:"string" input:"search"`  // 收银员ID
	SalesmensId string `json:"salesmen_id" label:"业务员ID" show:"true" sort:"8" type:"string" input:"search"` // 业务员ID

	StartDate *time.Time `json:"start_date" label:"开始日期" show:"true" sort:"9" type:"string" input:"date"` // 开始日期
	EndDate   *time.Time `json:"end_date" label:"结束日期" show:"true" sort:"10" type:"string" input:"date"`  // 结束日期

}

type OrderCreateReq struct {
	Type   enums.OrderType   `json:"type" required:"true"`   // 订单类型
	Source enums.OrderSource `json:"source" required:"true"` // 订单来源

	DiscountRate *float64 `json:"discount_rate"` // 整单折扣率
	AmountReduce float64  `json:"amount_reduce"` // 抹零
	IntegralUse  float64  `json:"integral_use"`  // 使用积分

	MemberId  string `json:"member_id" required:"true"`  // 会员ID
	StoreId   string `json:"store_id" required:"true"`   // 门店ID
	CashierId string `json:"cashier_id" required:"true"` // 收银员ID

	Salesmens []OrderCreateReqSalesmens `json:"salesmens" required:"true"` // 业务员
	Products  []OrderCreateReqProduct   `json:"products" required:"true"`  // 商品

	Remark string `json:"remark"` // 备注
}

func (req *OrderCreateReq) Validate() error {
	if len(req.Products) == 0 {
		return errors.New("商品不能为空")
	}
	if len(req.Salesmens) == 0 {
		return errors.New("业务员不能为空")
	}

	return nil
}

type OrderCreateReqSalesmens struct {
	SalesmenId     string  `json:"salesmen_id" required:"true"`     // 业务员ID
	CommissionRate float64 `json:"commission_rate" required:"true"` // 佣金比例
	IsMain         bool    `json:"is_main" required:"true"`         // 是否主业务员
}

type OrderCreateReqProduct struct {
	ProductId string   `json:"product_id" required:"true"` // 商品ID
	Quantity  int      `json:"quantity" required:"true"`   // 数量
	Discount  *float64 `json:"discount"`                   // 折扣
}
