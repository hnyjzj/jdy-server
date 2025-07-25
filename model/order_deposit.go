package model

import (
	"jdy/enums"
	"jdy/types"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// 定金单
type OrderDeposit struct {
	SoftDelete

	StoreId string `json:"store_id" gorm:"type:varchar(255);not NULL;comment:门店ID;"`  // 门店ID
	Store   Store  `json:"store" gorm:"foreignKey:StoreId;references:Id;comment:门店;"` // 门店

	Status enums.OrderDepositStatus `json:"status" gorm:"type:int(11);not NULL;comment:订单状态;"` // 订单状态

	MemberId string `json:"member_id" gorm:"type:varchar(255);not NULL;comment:会员ID;"`   // 会员ID
	Member   Member `json:"member" gorm:"foreignKey:MemberId;references:Id;comment:会员;"` // 会员

	CashierId string `json:"cashier_id" gorm:"type:varchar(255);not NULL;comment:收银员ID;"`    // 收银员ID
	Cashier   Staff  `json:"cashier" gorm:"foreignKey:CashierId;references:Id;comment:收银员;"` // 收银员
	ClerkId   string `json:"clerk_id" gorm:"type:varchar(255);not NULL;comment:导购员ID;"`      // 导购员ID
	Clerk     Staff  `json:"clerk" gorm:"foreignKey:ClerkId;references:Id;comment:导购员;"`     // 导购员

	Products []OrderDepositProduct `json:"products" gorm:"foreignKey:OrderId;references:Id;comment:成品;"` // 成品

	Price    decimal.Decimal `json:"price" gorm:"type:decimal(10,2);not NULL;comment:应付金额;"`     // 应付金额
	PricePay decimal.Decimal `json:"price_pay" gorm:"type:decimal(10,2);not NULL;comment:实付金额;"` // 实付金额

	Remarks  []string       `json:"remarks" gorm:"type:text;not NULL;serializer:json;comment:订单备注;"` // 订单备注
	Payments []OrderPayment `json:"payments" gorm:"foreignKey:OrderId;references:Id;comment:支付信息;"`  // 支付信息

	OperatorId string `json:"operator_id" gorm:"type:varchar(255);not NULL;comment:操作员ID;"`     // 操作员ID
	Operator   Staff  `json:"operator" gorm:"foreignKey:OperatorId;references:Id;comment:操作员;"` // 操作员
	IP         string `json:"ip" gorm:"type:varchar(255);not NULL;comment:IP地址;"`               // IP地址

	OrderSales []OrderSales `json:"order_sales" gorm:"many2many:order_sales_deposits;"`
}

func (OrderDeposit) WhereCondition(db *gorm.DB, req *types.OrderDepositWhere) *gorm.DB {
	// 订单编号
	if req.Id != "" {
		db = db.Where("id = ?", req.Id)
	}
	if req.Status != 0 {
		db = db.Where("status = ?", req.Status)
	} else {
		db = db.Where("status IN (?)", []enums.OrderDepositStatus{
			enums.OrderDepositStatusWaitPay,
			enums.OrderDepositStatusBooking,
			enums.OrderDepositStatusRefund,
		})
	}
	// 门店
	if req.StoreId != "" {
		db = db.Where("store_id = ?", req.StoreId)
	}
	// 会员
	if req.MemberId != "" {
		db = db.Where("member_id = ?", req.MemberId)
	}
	// 收银员
	if req.CashierId != "" {
		db = db.Where("cashier_id = ?", req.CashierId)
	}
	// 导购员
	if req.ClerkId != "" {
		db = db.Where("clerk_id = ?", req.ClerkId)
	}
	if req.StartDate != nil {
		db = db.Where("created_at >= ?", req.StartDate)
	}
	if req.EndDate != nil {
		db = db.Where("created_at <= ?", req.EndDate)
	}

	return db
}

func (OrderDeposit) Preloads(db *gorm.DB) *gorm.DB {
	db = db.Preload("Member")
	db = db.Preload("Store")
	db = db.Preload("Cashier")
	db = db.Preload("Clerk")
	db = db.Preload("Products", func(tx *gorm.DB) *gorm.DB {
		return tx.Preload("ProductFinished")
	})
	db = db.Preload("OrderSales")
	db = db.Preload("Payments")

	return db
}

// 定金单商品
type OrderDepositProduct struct {
	Model

	Status enums.OrderDepositStatus `json:"status" gorm:"type:int(11);not NULL;comment:状态;"` // 状态

	OrderId string       `json:"order_id" gorm:"type:varchar(255);not NULL;comment:定金单ID;"`            // 定金单ID
	Order   OrderDeposit `json:"order,omitempty" gorm:"foreignKey:OrderId;references:Id;comment:定金单;"` // 定金单

	IsOur           bool            `json:"is_our" gorm:"type:int(11);not NULL;default:0;comment:是否我司;"`                        // 是否我司
	ProductId       string          `json:"product_id" gorm:"type:varchar(255);not NULL;comment:成品ID;"`                         // 成品ID
	ProductFinished ProductFinished `json:"product_finished,omitempty" gorm:"foreignKey:ProductId;references:Id;comment:系统货品;"` // 系统货品
	ProductDemand   ProductFinished `json:"product_demand,omitempty" gorm:"type:text;serializer:json;comment:手动添加;"`            // 手动添加

	PriceGold decimal.Decimal `json:"price_gold" gorm:"type:decimal(10,2);not NULL;comment:金价;"` // 金价
	Price     decimal.Decimal `json:"price" gorm:"type:decimal(10,2);not NULL;comment:应付金额;"`    // 应付金额
}

func init() {
	// 注册模型
	RegisterModels(
		&OrderDeposit{},
		&OrderDepositProduct{},
	)
	// 重置表
	RegisterRefreshModels(
	// &OrderDeposit{},
	// &OrderDepositProduct{},
	)
}
