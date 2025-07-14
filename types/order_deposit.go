package types

import (
	"errors"
	"jdy/enums"
	"time"

	"github.com/shopspring/decimal"
)

type OrderDepositWhere struct {
	Id       string `json:"id" label:"订单编号" find:"true" sort:"1" type:"string" input:"text"`                                           // 订单编号
	StoreId  string `json:"store_id" label:"门店" find:"false" sort:"2" type:"string" input:"search" required:"true" binding:"required"` // 门店
	MemberId string `json:"member_id" label:"会员" find:"true" create:"true" sort:"3" type:"string" input:"search"`                      // 会员

	Status        enums.OrderDepositStatus `json:"status" label:"订单状态" find:"true" sort:"4" type:"number" input:"select" preset:"typeMap"`                                      // 订单状态
	PaymentMethod enums.OrderPaymentMethod `json:"payment_method" label:"支付方式" find:"false" create:"true" update:"true" sort:"6" type:"number" input:"select" preset:"typeMap"` // 支付方式

	CashierId string `json:"cashier_id" label:"收银员" find:"true" sort:"7" type:"string" input:"search"` // 收银员
	ClerkId   string `json:"clerk_id" label:"导购员" find:"true" sort:"8" type:"string" input:"search"`   // 导购员

	StartDate *time.Time `json:"start_date" label:"开始日期" find:"true" sort:"10" type:"string" input:"date"` // 开始日期
	EndDate   *time.Time `json:"end_date" label:"结束日期" find:"true" sort:"11" type:"string" input:"date"`   // 结束日期
}

type OrderDepositCreateReq struct {
	StoreId string `json:"store_id" binding:"required"` // 门店ID

	CashierId string `json:"cashier_id" binding:"required"` // 收银员ID
	ClerkId   string `json:"clerk_id" binding:"required"`   // 导购员ID
	MemberId  string `json:"member_id" binding:"required"`  // 会员ID

	Products []OrderDepositCreateReqProductFinished `json:"products" binding:"required"` // 货品
	Payments []OrderPaymentMethods                  `json:"payments" binding:"required"` // 支付方式
	Remarks  []string                               `json:"remarks"`                     // 备注
}

type OrderDepositCreateReqProductFinished struct {
	IsOur     bool   `json:"is_our" binding:"required"` // 是否本司货品
	ProductId string `json:"product_id"`                // 商品ID

	Name        string                  `json:"name" binding:"required"` // 名称
	LabelPrice  decimal.Decimal         `json:"label_price"`             // 标签价
	WeightMetal decimal.Decimal         `json:"weight_metal"`            // 金重
	LaborFee    decimal.Decimal         `json:"labor_fee"`               // 工费
	RetailType  enums.ProductRetailType `json:"retail_type"`             // 零售方式
	WeightGem   decimal.Decimal         `json:"weight_gem"`              // 主石重
	ColorGem    enums.ProductColor      `json:"color_gem"`               // 主石颜色
	ClarityGem  enums.ProductClarity    `json:"clarity_gem"`             // 主石净度

	PriceGold decimal.Decimal `json:"price_gold"` // 金价
	Price     decimal.Decimal `json:"price"`      // 定金金额
}

func (req *OrderDepositCreateReq) Validate() error {
	// 检查商品
	for _, p := range req.Products {
		if p.IsOur && p.ProductId == "" {
			return errors.New("商品ID不能为空")
		}
	}

	// 检查支付方式
	if len(req.Payments) == 0 {
		return errors.New("支付方式不能为空")
	}
	var total decimal.Decimal
	for _, payment := range req.Payments {
		if payment.Amount.LessThan(decimal.NewFromFloat(0)) {
			return errors.New("支付金额错误")
		}

		total = total.Add(payment.Amount) // 累加支付金额
	}

	return nil
}

type OrderDepositListReq struct {
	PageReq
	Where OrderDepositWhere `json:"where"`
}

type OrderDepositInfoReq struct {
	Id string `json:"id" required:"true"`
}

type OrderDepositRevokedReq struct {
	Id string `json:"id" required:"true"`
}

type OrderDepositPayReq struct {
	Id string `json:"id" required:"true"`
}

type OrderDepositRefundReq struct {
	Id        string `json:"id" required:"true"`         // 订单ID
	ProductId string `json:"product_id" required:"true"` // 商品ID
	Remark    string `json:"remark" required:"true"`     // 备注

	Payments []OrderPaymentMethods `json:"payments" binding:"required"` // 支付方式
}
