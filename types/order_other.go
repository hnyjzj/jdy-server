package types

import (
	"errors"
	"jdy/enums"
	"time"

	"github.com/shopspring/decimal"
)

type OrderOtherWhere struct {
	Id        string                   `json:"id" label:"订单编号" find:"true" create:"false" sort:"1" type:"string" input:"text"`                                       // 订单编号
	StoreId   string                   `json:"store_id" label:"门店" find:"false" create:"true" sort:"2" type:"string" input:"search" required:"true"`                 // 门店
	Type      enums.FinanceType        `json:"type" label:"收支类型" find:"true" create:"true" sort:"3" type:"number" input:"select" required:"true" preset:"typeMap"`   // 收支类型
	Content   string                   `json:"content" label:"收支内容" find:"true" create:"true" sort:"4" type:"string" input:"text" required:"true"`                   // 收支内容
	Source    enums.FinanceSourceOther `json:"source" label:"收支来源" find:"true" create:"true" sort:"5" type:"number" input:"select" required:"true" preset:"typeMap"` // 收支来源
	ClerkId   string                   `json:"clerk_id" label:"导购员" find:"false" create:"true" sort:"6" type:"string" input:"search" required:"true"`                // 导购员
	MemberId  string                   `json:"member_id" label:"会员" find:"false" create:"true" sort:"7" type:"string" input:"search" required:"true"`                // 会员
	OrderId   string                   `json:"order_id" label:"销售单" find:"false" create:"true" sort:"8" type:"string" input:"search" required:"false"`               // 销售单
	Amount    decimal.Decimal          `json:"amount" label:"收支金额" find:"false" create:"true" sort:"9" type:"float" input:"number" required:"true"`
	StartDate *time.Time               `json:"start_date" label:"开始日期" find:"true" create:"true" sort:"10" type:"string" input:"date"` // 开始日期
	EndDate   *time.Time               `json:"end_date" label:"结束日期" find:"true" create:"true" sort:"11" type:"string" input:"date"`   // 结束日期
}

type OrderOtherCreateReq struct {
	StoreId  string                   `json:"store_id" binding:"required"`  // 门店ID
	Type     enums.FinanceType        `json:"type" binding:"required"`      // 收支类型
	Content  string                   `json:"content" binding:"required"`   // 收支内容
	Source   enums.FinanceSourceOther `json:"source" binding:"required"`    // 收支来源
	ClerkId  string                   `json:"clerk_id" binding:"required"`  // 导购员ID
	MemberId string                   `json:"member_id" binding:"required"` // 会员ID
	OrderId  string                   `json:"order_id"`                     // 销售单ID
	Amount   decimal.Decimal          `json:"amount" binding:"required"`    // 收支金额
	Payments []OrderPaymentMethods    `json:"payments" binding:"required"`  // 支付方式
}

func (req *OrderOtherCreateReq) Validate() error {
	// 检查支付方式
	if len(req.Payments) == 0 {
		return errors.New("支付方式不能为空")
	}

	return nil
}

type OrderOtherListReq struct {
	PageReq
	Where OrderOtherWhere `json:"where"`
}

type OrderOtherInfoReq struct {
	Id string `json:"id" binding:"required"`
}

type OrderOtherUpdateReq struct {
	Id string `json:"id" binding:"required"`

	StoreId  string                   `json:"store_id"`  // 门店ID
	Type     enums.FinanceType        `json:"type"`      // 收支类型
	Content  string                   `json:"content"`   // 收支内容
	Source   enums.FinanceSourceOther `json:"source"`    // 收支来源
	ClerkId  string                   `json:"clerk_id"`  // 导购员ID
	MemberId string                   `json:"member_id"` // 会员ID
	OrderId  string                   `json:"order_id"`  // 销售单ID
	Amount   decimal.Decimal          `json:"amount"`    // 收支金额
	Payments []OrderPaymentMethods    `json:"payments"`  // 支付方式
}

type OrderOtherDeleteReq struct {
	Id string `json:"id" binding:"required"`
}
