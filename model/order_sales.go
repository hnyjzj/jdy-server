package model

import (
	"jdy/enums"
	"jdy/types"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// 销售单
type OrderSales struct {
	SoftDelete

	StoreId string `json:"store_id" gorm:"type:varchar(255);not NULL;comment:门店ID;"`  // 门店ID
	Store   Store  `json:"store" gorm:"foreignKey:StoreId;references:Id;comment:门店;"` // 门店

	Status enums.OrderSalesStatus `json:"status" gorm:"type:tinyint(2);not NULL;comment:订单状态;"` // 订单状态

	CashierId string            `json:"cashier_id" gorm:"type:varchar(255);not NULL;comment:收银员ID;"`    // 收银员ID
	Cashier   Staff             `json:"cashier" gorm:"foreignKey:CashierId;references:Id;comment:收银员;"` // 收银员
	Clerks    []OrderSalesClerk `json:"clerks" gorm:"foreignKey:OrderId;references:Id;comment:导购员列表;"`  // 导购员列表

	Source enums.OrderSource `json:"source" gorm:"type:tinyint(2);not NULL;comment:订单来源;"` // 订单来源

	MemberId string `json:"member_id" gorm:"type:varchar(255);not NULL;comment:会员ID;"`   // 会员ID
	Member   Member `json:"member" gorm:"foreignKey:MemberId;references:Id;comment:会员;"` // 会员

	HasIntegral       bool            `json:"has_integral" gorm:"type:tinyint(1);not NULL;comment:是否积分;"`          // 是否积分
	DiscountRate      decimal.Decimal `json:"discount_rate" gorm:"type:decimal(10,2);not NULL;comment:整单折扣;"`      // 整单折扣
	IntegralDeduction decimal.Decimal `json:"integral_deduction" gorm:"type:decimal(10,2);not NULL;comment:积分抵扣;"` // 积分抵扣
	RoundOff          decimal.Decimal `json:"round_off" gorm:"type:decimal(10,2);not NULL;comment:抹零;"`            // 抹零

	ProductFinisheds       []OrderSalesProductFinished   `json:"product_finisheds" gorm:"foreignKey:OrderId;references:Id;comment:成品;"`     // 成品
	ProductFinishedPrice   decimal.Decimal               `json:"product_finished_price" gorm:"type:decimal(10,2);not NULL;comment:货品金额;"`   // 货品金额
	ProductOlds            []OrderSalesProductOld        `json:"product_olds" gorm:"foreignKey:OrderId;references:Id;comment:旧料;"`          // 旧料
	ProductOldPrice        decimal.Decimal               `json:"product_old_price" gorm:"type:decimal(10,2);not NULL;comment:旧料抵扣;"`        // 旧料抵扣
	ProductAccessories     []OrderSalesProductAccessorie `json:"product_accessories" gorm:"foreignKey:OrderId;references:Id;comment:配件;"`   // 配件
	ProductAccessoriePrice decimal.Decimal               `json:"product_accessorie_price" gorm:"type:decimal(10,2);not NULL;comment:配件礼品;"` // 配件礼品

	// 其他单
	OrderDeposits []OrderDeposit  `json:"order_deposits" gorm:"many2many:order_sales_deposits;"`
	PriceDeposit  decimal.Decimal `json:"price_deposit" gorm:"type:decimal(10,2);not NULL;comment:定金抵扣;"`

	PriceOriginal decimal.Decimal `json:"price_original" gorm:"type:decimal(10,2);not NULL;comment:原价;"`   // 原价
	Price         decimal.Decimal `json:"price" gorm:"type:decimal(10,2);not NULL;comment:应付金额;"`          // 应付金额
	PricePay      decimal.Decimal `json:"price_pay" gorm:"type:decimal(10,2);not NULL;comment:实付金额;"`      // 实付金额
	PriceDiscount decimal.Decimal `json:"price_discount" gorm:"type:decimal(10,2);not NULL;comment:优惠金额;"` // 优惠金额
	Integral      decimal.Decimal `json:"integral" gorm:"type:decimal(10,2);not NULL;comment:积分;"`         // 积分

	Remark   string         `json:"remark" gorm:"type:varchar(255);not NULL;comment:订单备注;"`         // 订单备注
	Payments []OrderPayment `json:"payments" gorm:"foreignKey:OrderId;references:Id;comment:支付信息;"` // 支付信息

	OperatorId string `json:"operator_id" gorm:"type:varchar(255);not NULL;comment:操作员ID;"`     // 操作员ID
	Operator   Staff  `json:"operator" gorm:"foreignKey:OperatorId;references:Id;comment:操作员;"` // 操作员
	IP         string `json:"ip" gorm:"type:varchar(255);not NULL;comment:IP地址;"`               // IP地址
}

func (OrderSales) WhereCondition(db *gorm.DB, req *types.OrderSalesWhere) *gorm.DB {
	if req.Id != "" {
		db = db.Where("id = ?", req.Id)
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

// 销售单成品
type OrderSalesProductFinished struct {
	SoftDelete

	Status enums.OrderSalesStatus `json:"status" gorm:"type:tinyint(1);not NULL;comment:状态;"` // 状态

	OrderId string     `json:"order_id" gorm:"type:varchar(255);not NULL;comment:销售单ID;"`            // 销售单ID
	Order   OrderSales `json:"order,omitempty" gorm:"foreignKey:OrderId;references:Id;comment:销售单;"` // 销售单

	StoreId string `json:"store_id" gorm:"type:varchar(255);not NULL;comment:门店ID;"`            // 门店ID
	Store   Store  `json:"store,omitempty" gorm:"foreignKey:StoreId;references:Id;comment:门店;"` // 门店

	ProductId string          `json:"product_id" gorm:"type:varchar(255);not NULL;comment:产品ID;"`              // 产品ID
	Product   ProductFinished `json:"product,omitempty" gorm:"foreignKey:ProductId;references:Id;comment:产品;"` // 产品

	PriceGold         decimal.Decimal `json:"price_gold" gorm:"type:decimal(10,2);not NULL;comment:金价;"`           // 金价
	LaborFee          decimal.Decimal `json:"labor_fee" gorm:"type:decimal(10,2);not NULL;comment:工费;"`            // 工费
	DiscountFixed     decimal.Decimal `json:"discount_fixed" gorm:"type:decimal(10,2);not NULL;comment:固定折扣;"`     // 固定折扣
	IntegralDeduction decimal.Decimal `json:"integral_deduction" gorm:"type:decimal(10,2);not NULL;comment:积分抵扣;"` // 积分抵扣
	DiscountMember    decimal.Decimal `json:"discount_member" gorm:"type:decimal(10,2);not NULL;comment:会员折扣;"`    // 会员折扣
	RoundOff          decimal.Decimal `json:"round_off" gorm:"type:decimal(10,2);not NULL;comment:抹零;"`            // 抹零
	Integral          decimal.Decimal `json:"integral" gorm:"type:decimal(10,2);not NULL;comment:积分;"`             // 积分

	PriceOriginal decimal.Decimal `json:"price_original" gorm:"type:decimal(10,2);not NULL;comment:原价;"` // 原价
	Price         decimal.Decimal `json:"price" gorm:"type:decimal(10,2);not NULL;comment:应付金额;"`        // 应付金额
	DiscountFinal decimal.Decimal `json:"discount_final" gorm:"type:decimal(10,2);not NULL;comment:折扣;"` // 折扣
}

// 销售单旧料
type OrderSalesProductOld struct {
	SoftDelete

	Status enums.OrderSalesStatus `json:"status" gorm:"type:tinyint(1);not NULL;comment:状态;"` // 状态

	OrderId string     `json:"order_id" gorm:"type:varchar(255);not NULL;comment:销售单ID;"`            // 销售单ID
	Order   OrderSales `json:"order,omitempty" gorm:"foreignKey:OrderId;references:Id;comment:销售单;"` // 销售单

	StoreId string `json:"store_id" gorm:"type:varchar(255);not NULL;comment:门店ID;"`            // 门店ID
	Store   Store  `json:"store,omitempty" gorm:"foreignKey:StoreId;references:Id;comment:门店;"` // 门店

	ProductId string     `json:"product_id" gorm:"type:varchar(255);not NULL;comment:产品ID;"`              // 产品ID
	Product   ProductOld `json:"product,omitempty" gorm:"foreignKey:ProductId;references:Id;comment:产品;"` // 产品

	WeightMetal             decimal.Decimal            `json:"weight_metal" gorm:"type:decimal(10,2);comment:金重;"`                          // 金重
	RecyclePriceGold        decimal.Decimal            `json:"recycle_price_gold" gorm:"type:decimal(10,2);comment:回收金价;"`                  // 回收金价
	RecyclePriceLabor       decimal.Decimal            `json:"recycle_price_labor" gorm:"type:decimal(10,2);comment:回收工费;"`                 // 回收工费
	RecyclePriceLaborMethod enums.ProductRecycleMethod `json:"recycle_price_labor_method,omitempty" gorm:"type:tinyint(2);comment:回收工费方式;"` // 回收工费方式
	QualityActual           decimal.Decimal            `json:"quality_actual" gorm:"type:decimal(3,2);comment:实际成色;"`                       // 实际成色
	RecyclePrice            decimal.Decimal            `json:"recycle_price" gorm:"type:decimal(10,2);comment:回收金额;"`                       // 回收金额

	Integral decimal.Decimal `json:"integral" gorm:"type:decimal(10,2);not NULL;comment:积分;"` // 积分
}

// 销售单配件
type OrderSalesProductAccessorie struct {
	SoftDelete

	Status enums.OrderSalesStatus `json:"status" gorm:"type:tinyint(1);not NULL;comment:状态;"` // 状态

	OrderId string     `json:"order_id" gorm:"type:varchar(255);not NULL;comment:销售单ID;"`            // 销售单ID
	Order   OrderSales `json:"order,omitempty" gorm:"foreignKey:OrderId;references:Id;comment:销售单;"` // 销售单

	StoreId string `json:"store_id" gorm:"type:varchar(255);not NULL;comment:门店ID;"`            // 门店ID
	Store   Store  `json:"store,omitempty" gorm:"foreignKey:StoreId;references:Id;comment:门店;"` // 门店

	ProductId string            `json:"product_id" gorm:"type:varchar(255);not NULL;comment:产品ID;"`              // 产品ID
	Product   ProductAccessorie `json:"product,omitempty" gorm:"foreignKey:ProductId;references:Id;comment:产品;"` // 产品

	Quantity int64           `json:"quantity" gorm:"type:int(11);not NULL;comment:数量;"`       // 数量
	Price    decimal.Decimal `json:"price" gorm:"type:decimal(10,2);not NULL;comment:应付金额;"`  // 应付金额
	Integral decimal.Decimal `json:"integral" gorm:"type:decimal(10,2);not NULL;comment:积分;"` // 积分
}

// 销售单导购员
type OrderSalesClerk struct {
	SoftDelete

	OrderId string     `json:"order_id" gorm:"type:varchar(255);not NULL;comment:销售单ID;"`            // 销售单ID
	Order   OrderSales `json:"order,omitempty" gorm:"foreignKey:OrderId;references:Id;comment:销售单;"` // 销售单

	SalesmanId string `json:"salesman_id" gorm:"type:varchar(255);not NULL;comment:导购员ID;"`               // 导购员ID
	Salesman   Staff  `json:"salesman,omitempty" gorm:"foreignKey:SalesmanId;references:Id;comment:导购员;"` // 导购员

	PerformanceAmount decimal.Decimal `json:"performance_amount" gorm:"type:decimal(10,2);not NULL;comment:业绩金额;"` // 业绩金额
	PerformanceRate   decimal.Decimal `json:"performance_rate" gorm:"type:decimal(5,2);not NULL;comment:业绩比例;"`    // 业绩比例

	IsMain bool `json:"is_main" gorm:"type:tinyint(1);not NULL;comment:是否主导购员;"` // 是否主导购员
}

func init() {
	// 注册模型
	RegisterModels(
		&OrderSales{},
		&OrderSalesProductFinished{},
		&OrderSalesProductOld{},
		&OrderSalesProductAccessorie{},
		&OrderSalesClerk{},
	)
	// 重置表
	RegisterRefreshModels(
	// &OrderSales{},
	// &OrderSalesProductFinished{},
	// &OrderSalesProductOld{},
	// &OrderSalesProductAccessorie{},
	// &OrderSalesClerk{},
	)
}
