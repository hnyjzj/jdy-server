package model

import (
	"jdy/enums"
	"jdy/types"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// 其他单
type OrderOther struct {
	SoftDelete

	StoreId string `json:"store_id" gorm:"type:varchar(255);not NULL;comment:门店ID;"`  // 门店ID
	Store   Store  `json:"store" gorm:"foreignKey:StoreId;references:Id;comment:门店;"` // 门店

	Type    enums.FinanceType        `json:"type" gorm:"type:int(11);not NULL;comment:订单类型;"`         // 订单类型
	Content string                   `json:"content" gorm:"type:varchar(500);not NULL;comment:订单内容;"` // 订单内容
	Source  enums.FinanceSourceOther `json:"source" gorm:"type:int(11);not NULL;comment:订单来源;"`       // 订单来源

	ClerkId string `json:"clerk_id" gorm:"type:varchar(255);not NULL;comment:导购员ID;"`  // 导购员ID
	Clerk   Staff  `json:"clerk" gorm:"foreignKey:ClerkId;references:Id;comment:导购员;"` // 导购员

	MemberId string `json:"member_id" gorm:"type:varchar(255);not NULL;comment:会员ID;"`   // 会员ID
	Member   Member `json:"member" gorm:"foreignKey:MemberId;references:Id;comment:会员;"` // 会员

	OrderId string     `json:"order_id" gorm:"type:varchar(255);comment:销售单ID;"`             // 销售单ID
	Order   OrderSales `json:"order" gorm:"foreignKey:OrderId;references:Id;comment:关联销售单;"` // 关联销售单

	Amount   decimal.Decimal `json:"amount" gorm:"type:decimal(10,2);not NULL;comment:金额;"`          // 金额
	Payments []OrderPayment  `json:"payments" gorm:"foreignKey:OrderId;references:Id;comment:支付信息;"` // 支付信息

	OperatorId string `json:"operator_id" gorm:"type:varchar(255);not NULL;comment:操作员ID;"`     // 操作员ID
	Operator   Staff  `json:"operator" gorm:"foreignKey:OperatorId;references:Id;comment:操作员;"` // 操作员
	IP         string `json:"ip" gorm:"type:varchar(255);not NULL;comment:IP地址;"`               // IP地址
}

func (OrderOther) WhereCondition(db *gorm.DB, req *types.OrderOtherWhere) *gorm.DB {
	if req.Id != "" {
		db = db.Where("id = ?", req.Id)
	}
	if req.StoreId != "" {
		db = db.Where("store_id = ?", req.StoreId)
	}
	if req.Type != 0 {
		db = db.Where("type = ?", req.Type)
	}
	if req.Content != "" {
		db = db.Where("content LIKE ?", "%"+req.Content+"%")
	}
	if req.Source != 0 {
		db = db.Where("source = ?", req.Source)
	}
	if req.ClerkId != "" {
		db = db.Where("clerk_id = ?", req.ClerkId)
	}
	if req.MemberId != "" {
		db = db.Where("member_id = ?", req.MemberId)
	}
	if req.OrderId != "" {
		db = db.Where("order_id = ?", req.OrderId)
	}
	if req.StartDate != nil {
		db = db.Where("created_at >= ?", req.StartDate)
	}
	if req.EndDate != nil {
		db = db.Where("created_at <= ?", req.EndDate)
	}

	return db
}

func (OrderOther) Preloads(db *gorm.DB) *gorm.DB {
	db = db.Preload("Store")
	db = db.Preload("Clerk")
	db = db.Preload("Member")
	db = db.Preload("Order")
	db = db.Preload("Payments")

	return db
}

func init() {
	// 注册模型
	RegisterModels(
		&OrderOther{},
	)
	// 重置表
	RegisterRefreshModels(
	// &OrderOther{},
	)
}
