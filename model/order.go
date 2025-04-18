package model

import (
	"jdy/enums"
	"jdy/types"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// 订单
type Order struct {
	SoftDelete

	Type   enums.OrderType   `json:"type" gorm:"type:tinyint(2);not NULL;comment:订单类型;"`     // 订单类型
	Status enums.OrderStatus `json:"status" gorm:"type:tinyint(2);not NULL;comment:订单状态;"`   // 订单状态
	Source enums.OrderSource `json:"source" gorm:"type:tinyint(2);not NULL;comment:订单来源;"`   // 订单来源
	Remark string            `json:"remark" gorm:"type:varchar(255);not NULL;comment:订单备注;"` // 订单备注

	Amount         decimal.Decimal `json:"amount" gorm:"type:decimal(10,2);not NULL;comment:应付金额;"`        // 应付金额
	AmountOriginal decimal.Decimal `json:"amount_original" gorm:"type:decimal(10,2);not NULL;comment:原价;"` // 原价
	AmountPay      decimal.Decimal `json:"amount_pay" gorm:"type:decimal(10,2);not NULL;comment:实付金额;"`    // 实付金额

	DiscountRate      decimal.Decimal `json:"discount_rate" gorm:"type:decimal(5,2);not NULL;comment:整单折扣;"`          // 整单折扣
	DiscountAmount    decimal.Decimal `json:"discount_amount" gorm:"type:decimal(10,2);not NULL;comment:整单折扣金额;"`     // 整单折扣金额
	AmountReduce      decimal.Decimal `json:"amount_reduce" gorm:"type:decimal(10,2);not NULL;comment:抹零金额;"`         // 抹零金额
	AmountOldMaterial decimal.Decimal `json:"amount_old_material" gorm:"type:decimal(10,2);not NULL;comment:旧料抵扣金额;"` // 旧料抵扣金额

	IntegralPresent decimal.Decimal `json:"integral_present" gorm:"type:int(11);not NULL;comment:赠送积分;"` // 赠送积分
	IntegralUse     decimal.Decimal `json:"integral_use" gorm:"type:int(11);not NULL;comment:使用积分;"`     // 使用积分

	MemberId string `json:"member_id" gorm:"type:varchar(255);not NULL;comment:会员ID;"`   // 会员ID
	Member   Member `json:"member" gorm:"foreignKey:MemberId;references:Id;comment:会员;"` // 会员

	StoreId   string `json:"store_id" gorm:"type:varchar(255);not NULL;comment:门店ID;"`       // 门店ID
	Store     Store  `json:"store" gorm:"foreignKey:StoreId;references:Id;comment:门店;"`      // 门店
	CashierId string `json:"cashier_id" gorm:"type:varchar(255);not NULL;comment:收银员ID;"`    // 收银员ID
	Cashier   Staff  `json:"cashier" gorm:"foreignKey:CashierId;references:Id;comment:收银员;"` // 收银员

	Salesmans         []OrderSalesman `json:"salesmans" gorm:"foreignKey:OrderId;references:Id;comment:订单导购员;"`      // 订单导购员
	ProductFinished   []OrderProduct  `json:"product_finished" gorm:"foreignKey:OrderId;references:Id;comment:成品;"`  // 成品
	ProductOld        []OrderProduct  `json:"product_old" gorm:"foreignKey:OrderId;references:Id;comment:旧料;"`       // 旧料
	ProductAccessorie []OrderProduct  `json:"product_accessory" gorm:"foreignKey:OrderId;references:Id;comment:配件;"` // 配件

	OperatorId string `json:"operator_id" gorm:"type:varchar(255);not NULL;comment:操作员ID;"`     // 操作员ID
	Operator   Staff  `json:"operator" gorm:"foreignKey:OperatorId;references:Id;comment:操作员;"` // 操作员
	IP         string `json:"ip" gorm:"type:varchar(255);not NULL;comment:IP地址;"`               // IP地址
}

func (Order) WhereCondition(db *gorm.DB, req *types.OrderWhere) *gorm.DB {
	if req.Id != "" {
		db = db.Where("id = ?", req.Id)
	}
	if req.Type != 0 {
		db = db.Where("type = ?", req.Type)
	}
	if req.Status != 0 {
		db = db.Where("status = ?", req.Status)
	}
	if req.Source != 0 {
		db = db.Where("source = ?", req.Source)
	}
	if req.MemberId != "" {
		db = db.Where("member_id = ?", req.MemberId)
	}
	if req.StoreId != "" {
		db = db.Where("store_id = ?", req.StoreId)
	}
	if req.CashierId != "" {
		db = db.Where("cashier_id = ?", req.CashierId)
	}
	if req.SalesmanId != "" {
		db = db.Where("id IN (SELECT order_id FROM order_salesmans WHERE salesman_id = ?)", req.SalesmanId)
	}
	if req.ProductId != "" {
		db = db.Where("id IN (SELECT order_id FROM order_products WHERE product_id = ?)", req.ProductId)
	}
	if req.StartDate != nil && req.EndDate == nil {
		db = db.Where("created_at >= ?", req.StartDate)
	}
	if req.StartDate == nil && req.EndDate != nil {
		db = db.Where("created_at <= ?", req.EndDate)
	}
	if req.StartDate != nil && req.EndDate != nil {
		db = db.Where("created_at BETWEEN ? AND ?", req.StartDate, req.EndDate)
	}

	return db
}

// 订单导购员
type OrderSalesman struct {
	Model

	OrderId string `json:"order_id" gorm:"type:varchar(255);not NULL;comment:订单ID;"`            // 订单ID
	Order   Order  `json:"order,omitempty" gorm:"foreignKey:OrderId;references:Id;comment:订单;"` // 订单

	SalesmanId string `json:"salesman_id" gorm:"type:varchar(255);not NULL;comment:导购员ID;"`               // 导购员ID
	Salesman   Staff  `json:"salesman,omitempty" gorm:"foreignKey:SalesmanId;references:Id;comment:导购员;"` // 导购员

	PerformanceAmount decimal.Decimal `json:"performance_amount" gorm:"type:decimal(10,2);not NULL;comment:业绩金额;"` // 业绩金额
	PerformanceRate   decimal.Decimal `json:"performance_rate" gorm:"type:decimal(5,2);not NULL;comment:业绩比例;"`    // 业绩比例

	IsMain bool `json:"is_main" gorm:"type:tinyint(1);not NULL;comment:是否主导购员;"` // 是否主导购员
}

func (OrderSalesman) TableName() string {
	return "order_salesmans"
}

// 订单商品
type OrderProduct struct {
	Model

	Status enums.OrderStatus `json:"status" gorm:"type:tinyint(1);not NULL;comment:状态;"` // 状态

	OrderId string `json:"order_id" gorm:"type:varchar(255);not NULL;comment:订单ID;"`            // 订单ID
	Order   Order  `json:"order,omitempty" gorm:"foreignKey:OrderId;references:Id;comment:订单;"` // 订单

	ProductId string `json:"product_id" gorm:"type:varchar(255);not NULL;comment:产品ID;"` // 产品ID

	Quantity       int64           `json:"quantity" gorm:"type:int(11);not NULL;comment:数量;"`              // 数量
	Price          decimal.Decimal `json:"price" gorm:"type:decimal(10,2);not NULL;comment:单价;"`           // 单价
	Amount         decimal.Decimal `json:"amount" gorm:"type:decimal(10,2);not NULL;comment:应付金额;"`        // 应付金额
	AmountOriginal decimal.Decimal `json:"amount_original" gorm:"type:decimal(10,2);not NULL;comment:原价;"` // 原价

	Discount       decimal.Decimal `json:"discount" gorm:"type:decimal(10,2);not NULL;comment:折扣;"`          // 折扣
	DiscountAmount decimal.Decimal `json:"discount_amount" gorm:"type:decimal(10,2);not NULL;comment:折扣金额;"` // 折扣金额

	Integral decimal.Decimal `json:"integral" gorm:"type:decimal(10,2);not NULL;comment:增加积分;"` // 增加积分
}

func init() {
	// 注册模型
	RegisterModels(
		&Order{},
		&OrderSalesman{},
		&OrderProduct{},
	)
	// 重置表
	RegisterRefreshModels(
	// &Order{},
	// &OrderSalesman{},
	// &OrderProduct{},
	)
}
