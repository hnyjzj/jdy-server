package model

import (
	"jdy/enums"

	"github.com/shopspring/decimal"
)

// 定金单
type OrderDeposit struct {
	SoftDelete

	StoreId string `json:"store_id" gorm:"type:varchar(255);not NULL;comment:门店ID;"`  // 门店ID
	Store   Store  `json:"store" gorm:"foreignKey:StoreId;references:Id;comment:门店;"` // 门店

	Status enums.OrderStatus `json:"status" gorm:"type:tinyint(2);not NULL;comment:订单状态;"` // 订单状态

	MemberId string `json:"member_id" gorm:"type:varchar(255);not NULL;comment:会员ID;"`   // 会员ID
	Member   Member `json:"member" gorm:"foreignKey:MemberId;references:Id;comment:会员;"` // 会员

	CashierId string `json:"cashier_id" gorm:"type:varchar(255);not NULL;comment:收银员ID;"`    // 收银员ID
	Cashier   Staff  `json:"cashier" gorm:"foreignKey:CashierId;references:Id;comment:收银员;"` // 收银员
	ClerkId   string `json:"clerk_id" gorm:"type:varchar(255);not NULL;comment:导购员ID;"`      // 导购员ID
	Clerk     Staff  `json:"clerk" gorm:"foreignKey:ClerkId;references:Id;comment:导购员;"`     // 导购员

	Products []OrderDepositProduct `json:"products" gorm:"foreignKey:OrderId;references:Id;comment:成品;"` // 成品

	Price    decimal.Decimal `json:"price" gorm:"type:decimal(10,2);not NULL;comment:应付金额;"`     // 应付金额
	PricePay decimal.Decimal `json:"price_pay" gorm:"type:decimal(10,2);not NULL;comment:实付金额;"` // 实付金额

	Remark   string         `json:"remark" gorm:"type:varchar(255);not NULL;comment:订单备注;"`         // 订单备注
	Payments []OrderPayment `json:"payments" gorm:"foreignKey:OrderId;references:Id;comment:支付信息;"` // 支付信息

	OperatorId string `json:"operator_id" gorm:"type:varchar(255);not NULL;comment:操作员ID;"`     // 操作员ID
	Operator   Staff  `json:"operator" gorm:"foreignKey:OperatorId;references:Id;comment:操作员;"` // 操作员
	IP         string `json:"ip" gorm:"type:varchar(255);not NULL;comment:IP地址;"`               // IP地址
}

// func (OrderDeposit) WhereCondition(db *gorm.DB, req *types.OrderDepositWhere) *gorm.DB {
// 	return db
// }

// 定金单商品
type OrderDepositProduct struct {
	Model

	Status enums.OrderStatus `json:"status" gorm:"type:tinyint(1);not NULL;comment:状态;"` // 状态

	OrderId string       `json:"order_id" gorm:"type:varchar(255);not NULL;comment:定金单ID;"`            // 定金单ID
	Order   OrderDeposit `json:"order,omitempty" gorm:"foreignKey:OrderId;references:Id;comment:定金单;"` // 定金单

	ProductId       string          `json:"product_id" gorm:"type:varchar(255);not NULL;comment:成品ID;"`                       // 成品ID
	Product         ProductFinished `json:"product,omitempty" gorm:"foreignKey:ProductId;references:Id;comment:手动添加;"`        // 手动添加
	ProductFinished ProductFinished `json:"product_finished,omitempty" gorm:"foreignKey:ProductId;references:Id;comment:成品;"` // 成品

	PriceGold    decimal.Decimal `json:"price_gold" gorm:"type:decimal(10,2);not NULL;comment:金价;"`      // 金价
	PricePayable decimal.Decimal `json:"price_payable" gorm:"type:decimal(10,2);not NULL;comment:应付金额;"` // 应付金额
}

func init() {
	// 注册模型
	RegisterModels(
	// &OrderDeposit{},
	// &OrderDepositClerk{},
	// &OrderDepositProduct{},
	)
	// 重置表
	RegisterRefreshModels(
	// &OrderDeposit{},
	// &OrderDepositClerk{},
	// &OrderDepositProduct{},
	)
}
