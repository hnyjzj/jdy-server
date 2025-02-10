package types

import (
	"errors"
	"jdy/enums"
	"time"

	"github.com/shopspring/decimal"
)

type OrderWhere struct {
	Id       string `json:"id" label:"订单编号" show:"true" sort:"1" type:"string" input:"text"`                       // 订单编号
	StoreId  string `json:"store_id" label:"门店ID" show:"true" sort:"2" type:"string" input:"text" required:"true"` // 门店ID
	MemberId string `json:"member_id" label:"会员ID" show:"true" sort:"3" type:"string" input:"text"`                // 会员ID

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

	DiscountRate decimal.Decimal `json:"discount_rate"` // 整单折扣率
	AmountReduce decimal.Decimal `json:"amount_reduce"` // 抹零
	IntegralUse  decimal.Decimal `json:"integral_use"`  // 使用积分

	MemberId  string `json:"member_id" required:"true"`  // 会员ID
	StoreId   string `json:"store_id" required:"true"`   // 门店ID
	CashierId string `json:"cashier_id" required:"true"` // 收银员ID

	Salesmens []*OrderCreateReqSalesmens `json:"salesmens" required:"true"` // 业务员
	Products  []*OrderCreateReqProduct   `json:"products" required:"true"`  // 商品

	Remark string `json:"remark"` // 备注
}

func (req *OrderCreateReq) Validate() error {
	if len(req.Products) == 0 {
		return errors.New("商品不能为空")
	}

	if len(req.Salesmens) == 0 {
		return errors.New("业务员不能为空")
	}

	if !req.DiscountRate.IsZero() {
		if req.DiscountRate.LessThan(decimal.NewFromFloat(0)) || req.DiscountRate.GreaterThan(decimal.NewFromFloat(10)) {
			return errors.New("整单折扣错误")
		}
	} else {
		req.DiscountRate = decimal.NewFromFloat(10)
	}

	for _, salesmen := range req.Salesmens {
		if !salesmen.CommissionRate.IsZero() {
			if salesmen.CommissionRate.LessThan(decimal.NewFromFloat(0)) || salesmen.CommissionRate.GreaterThan(decimal.NewFromFloat(10)) {
				return errors.New("佣金比例错误")
			}
		} else {
			salesmen.CommissionRate = decimal.NewFromFloat(10)
		}
	}

	for _, product := range req.Products {
		if product.Quantity <= 0 {
			return errors.New("商品数量错误")
		}

		if !product.Discount.IsZero() {
			if product.Discount.LessThan(decimal.NewFromFloat(0)) || product.Discount.GreaterThan(decimal.NewFromFloat(10)) {
				return errors.New("商品折扣错误")
			}
		} else {
			product.Discount = decimal.NewFromFloat(10)
		}
	}

	return nil
}

type OrderCreateReqSalesmens struct {
	SalesmenId     string          `json:"salesmen_id" required:"true"`     // 业务员ID
	CommissionRate decimal.Decimal `json:"commission_rate" required:"true"` // 佣金比例
	IsMain         bool            `json:"is_main" required:"true"`         // 是否主业务员
}

type OrderCreateReqProduct struct {
	ProductId string          `json:"product_id" required:"true"` // 商品ID
	Quantity  int64           `json:"quantity" required:"true"`   // 数量
	Discount  decimal.Decimal `json:"discount"`                   // 折扣
}

type OrderListReq struct {
	PageReq
	Where OrderWhere `json:"where"`
}

type OrderInfoReq struct {
	Id string `json:"id" required:"true"`
}
