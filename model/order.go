package model

import (
	"jdy/enums"

	"github.com/shopspring/decimal"
)

// 支付信息
type OrderPayment struct {
	SoftDelete

	StoreId string `json:"store_id" gorm:"type:varchar(255);not NULL;comment:店铺ID;"`  // 店铺ID
	Store   Store  `json:"store" gorm:"foreignKey:StoreId;references:Id;comment:店铺;"` // 店铺

	Type   enums.FinanceType   `json:"type" gorm:"type:tinyint(1);not NULL;comment:支付类型;"`   // 支付类型
	Source enums.FinanceSource `json:"source" gorm:"type:tinyint(1);not NULL;comment:支付来源;"` // 支付来源

	OrderId      string       `json:"order_id" gorm:"type:varchar(255);not NULL;comment:销售单ID;"`                    // 销售单ID
	OrderSales   OrderSales   `json:"order,omitempty" gorm:"foreignKey:OrderId;references:Id;comment:销售单;"`         // 销售单
	OrderRepair  OrderRepair  `json:"order_repair,omitempty" gorm:"foreignKey:OrderId;references:Id;comment:维修单;"`  // 维修单
	OrderDeposit OrderDeposit `json:"order_deposit,omitempty" gorm:"foreignKey:OrderId;references:Id;comment:定金单;"` // 定金单
	OrderOther   OrderOther   `json:"order_other,omitempty" gorm:"foreignKey:OrderId;references:Id;comment:其他单;"`   // 其他单

	PaymentMethod enums.OrderPaymentMethod `json:"payment_method" gorm:"type:tinyint(1);not NULL;comment:支付方式;"` // 支付方式
	Amount        decimal.Decimal          `json:"amount" gorm:"type:decimal(10,2);not NULL;comment:金额;"`        // 金额
}

func init() {
	// 注册模型
	RegisterModels(
		&OrderPayment{},
	)
	// 重置表
	RegisterRefreshModels(
	// &OrderPayment{},
	)
}
