package types

import (
	"jdy/enums"

	"github.com/shopspring/decimal"
)

type OrderPaymentMethods struct {
	PaymentMethod enums.OrderPaymentMethod `json:"payment_method" required:"true"` // 支付方式
	Amount        decimal.Decimal          `json:"amount" required:"true"`         // 支付金额
}
