package model

import (
	"jdy/enums"
	"jdy/types"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// 支付信息
type OrderPayment struct {
	SoftDelete

	StoreId string `json:"store_id" gorm:"type:varchar(255);not NULL;comment:店铺ID;"`  // 店铺ID
	Store   Store  `json:"store" gorm:"foreignKey:StoreId;references:Id;comment:店铺;"` // 店铺

	Type   enums.FinanceType   `json:"type" gorm:"type:int(11);not NULL;comment:支付类型;"`   // 支付类型
	Source enums.FinanceSource `json:"source" gorm:"type:int(11);not NULL;comment:支付来源;"` // 支付来源

	OrderId      string          `json:"order_id" gorm:"type:varchar(255);not NULL;comment:订单ID;"`                     // 订单ID
	OrderType    enums.OrderType `json:"order_type" gorm:"type:int(11);not NULL;comment:订单类型;"`                        // 订单类型
	OrderSales   OrderSales      `json:"order,omitempty" gorm:"foreignKey:OrderId;references:Id;comment:销售单;"`         // 销售单
	OrderRepair  OrderRepair     `json:"order_repair,omitempty" gorm:"foreignKey:OrderId;references:Id;comment:维修单;"`  // 维修单
	OrderDeposit OrderDeposit    `json:"order_deposit,omitempty" gorm:"foreignKey:OrderId;references:Id;comment:定金单;"` // 定金单
	OrderOther   OrderOther      `json:"order_other,omitempty" gorm:"foreignKey:OrderId;references:Id;comment:其他单;"`   // 其他单
	OrderRefund  OrderRefund     `json:"order_refund,omitempty" gorm:"foreignKey:OrderId;references:Id;comment:退单;"`   // 退单

	PaymentMethod enums.OrderPaymentMethod `json:"payment_method" gorm:"type:int(11);not NULL;comment:支付方式;"` // 支付方式
	Amount        decimal.Decimal          `json:"amount" gorm:"type:decimal(10,2);not NULL;comment:金额;"`     // 金额
}

// 退单
type OrderRefund struct {
	SoftDelete

	StoreId string `json:"store_id" gorm:"type:varchar(255);not NULL;comment:店铺ID;"`  // 店铺ID
	Store   Store  `json:"store" gorm:"foreignKey:StoreId;references:Id;comment:店铺;"` // 店铺

	OrderId   string          `json:"order_id" gorm:"type:varchar(255);not NULL;comment:订单ID;"` // 订单ID
	OrderType enums.OrderType `json:"order_type" gorm:"type:int(11);not NULL;comment:订单类型;"`    // 订单类型

	MemberId string `json:"member_id" gorm:"type:varchar(255);not NULL;comment:会员ID;"`             // 会员ID
	Member   Member `json:"member,omitempty" gorm:"foreignKey:MemberId;references:Id;comment:会员;"` // 会员

	Type    enums.ProductType `json:"type" gorm:"type:int(11);not NULL;comment:产品类型;"`    // 产品类型
	Code    string            `json:"code" gorm:"type:varchar(255);not NULL;comment:条码;"` // 条码
	Product any               `json:"product,omitempty" gorm:"-"`
	Name    string            `json:"name" gorm:"type:varchar(255);not NULL;comment:名称;"`   // 名称
	Remark  string            `json:"remark" gorm:"type:varchar(255);not NULL;comment:备注;"` // 备注

	Quantity      int64           `json:"quantity" gorm:"type:int(11);not NULL;comment:数量;"`              // 数量
	Price         decimal.Decimal `json:"price" gorm:"type:decimal(10,2);not NULL;comment:实退金额;"`         // 实退金额
	PriceOriginal decimal.Decimal `json:"price_original" gorm:"type:decimal(10,2);not NULL;comment:原金额;"` // 原金额

	OperatorId string `json:"operator_id" gorm:"type:varchar(255);not NULL;comment:操作员ID;"`               // 操作员ID
	Operator   Staff  `json:"operator,omitempty" gorm:"foreignKey:OperatorId;references:Id;comment:操作员;"` // 操作员
	IP         string `json:"ip" gorm:"type:varchar(255);not NULL;comment:IP地址;"`                         // IP地址
}

func (OrderRefund) WhereCondition(db *gorm.DB, req *types.OrderSalesRefundWhere) *gorm.DB {
	if req.OrderId != "" {
		db = db.Where("order_id = ?", req.OrderId)
	}
	if req.StoreId != "" {
		db = db.Where("store_id = ?", req.StoreId)
	}
	if req.Phone != "" {
		db = db.Where("member_id = (SELECT id FROM members WHERE phone = ?)", req.Phone)
	}
	if req.Code != "" {
		db = db.Where("code = ?", req.Code)
	}
	if req.Name != "" {
		db = db.Where("name = ?", req.Name)
	}
	if req.StartDate != nil {
		db = db.Where("created_at >= ?", req.StartDate)
	}
	if req.EndDate != nil {
		db = db.Where("created_at <= ?", req.EndDate)
	}

	return db
}

func (OrderRefund) Preloads(db *gorm.DB) *gorm.DB {
	db = db.Preload("Store")
	db = db.Preload("Member")

	return db
}

func init() {
	// 注册模型
	RegisterModels(
		&OrderPayment{},
		&OrderRefund{},
	)
	// 重置表
	RegisterRefreshModels(
	// &OrderPayment{},
	// &OrderRefund{},
	)
}
