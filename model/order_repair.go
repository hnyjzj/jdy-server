package model

import (
	"jdy/enums"
	"jdy/types"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// 维修单
type OrderRepair struct {
	SoftDelete

	StoreId string `json:"store_id" gorm:"type:varchar(255);not NULL;comment:门店ID;"`  // 门店ID
	Store   Store  `json:"store" gorm:"foreignKey:StoreId;references:Id;comment:门店;"` // 门店

	Status enums.OrderRepairStatus `json:"status" gorm:"type:int(11);not NULL;comment:订单状态;"` // 订单状态

	ReceptionistId string `json:"receptionist_id" gorm:"type:varchar(255);not NULL;comment:接待人ID;"`         // 接待人ID
	Receptionist   Staff  `json:"receptionist" gorm:"foreignKey:ReceptionistId;references:Id;comment:接待人;"` // 接待人
	CashierId      string `json:"cashier_id" gorm:"type:varchar(255);not NULL;comment:收银员ID;"`              // 收银员ID
	Cashier        Staff  `json:"cashier" gorm:"foreignKey:CashierId;references:Id;comment:收银员;"`           // 收银员

	MemberId string `json:"member_id" gorm:"type:varchar(255);not NULL;comment:会员ID;"`   // 会员ID
	Member   Member `json:"member" gorm:"foreignKey:MemberId;references:Id;comment:会员;"` // 会员

	Name string `json:"name" gorm:"type:varchar(255);not NULL;comment:维修项目;"` // 维修项目
	Desc string `json:"desc" gorm:"type:text;not NULL;comment:描述;"`           // 问题描述

	DeliveryMethod enums.DeliveryMethod `json:"delivery_method" gorm:"type:int(11);not NULL;comment:取货方式;"` // 取货方式
	Province       string               `json:"province" gorm:"type:varchar(255);not NULL;comment:省;"`      // 省
	City           string               `json:"city" gorm:"type:varchar(255);not NULL;comment:市;"`          // 市
	Area           string               `json:"area" gorm:"type:varchar(255);not NULL;comment:区;"`          // 区
	Address        string               `json:"address" gorm:"type:varchar(255);not NULL;comment:地址;"`      // 地址

	Products []OrderRepairProduct `json:"products" gorm:"foreignKey:OrderId;references:Id;comment:维修单商品;"` // 维修单商品

	Images []string `json:"images" gorm:"type:text;not NULL;serializer:json;comment:图片;"` // 图片

	Expense  decimal.Decimal `json:"expense" gorm:"type:decimal(10,2);not NULL;comment:维修费;"`        // 维修费
	Cost     decimal.Decimal `json:"cost" gorm:"type:decimal(10,2);not NULL;comment:维修成本;"`          // 维修成本
	Payments []OrderPayment  `json:"payments" gorm:"foreignKey:OrderId;references:Id;comment:支付信息;"` // 支付信息

	OperatorId string `json:"operator_id" gorm:"type:varchar(255);not NULL;comment:操作员ID;"`     // 操作员ID
	Operator   Staff  `json:"operator" gorm:"foreignKey:OperatorId;references:Id;comment:操作员;"` // 操作员
	IP         string `json:"ip" gorm:"type:varchar(255);not NULL;comment:IP地址;"`               // IP地址
}

func (OrderRepair) WhereCondition(db *gorm.DB, req *types.OrderRepairWhere) *gorm.DB {
	if req.Id != "" {
		db = db.Where("id = ?", req.Id)
	}
	if req.Status != 0 {
		db = db.Where("status = ?", req.Status)
	} else {
		db = db.Where("status IN (?)", []enums.OrderRepairStatus{
			enums.OrderRepairStatusWaitPay,
			enums.OrderRepairStatusStoreReceived,
			enums.OrderRepairStatusSendOut,
			enums.OrderRepairStatusRepairing,
			enums.OrderRepairStatusSendBack,
			enums.OrderRepairStatusWaitPickUp,
		})
	}
	if req.ReceptionistId != "" {
		db = db.Where("receptionist_id = ?", req.ReceptionistId)
	}
	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.StoreId != "" {
		db = db.Where("store_id = ?", req.StoreId)
	}
	if req.MemberId != "" {
		db = db.Where("member_id = ?", req.MemberId)
	}
	if req.DeliveryMethod != 0 {
		db = db.Where("delivery_method = ?", req.DeliveryMethod)
	}
	if req.Province != "" {
		db = db.Where("province = ?", req.Province)
	}
	if req.City != "" {
		db = db.Where("city = ?", req.City)
	}
	if req.Area != "" {
		db = db.Where("area = ?", req.Area)
	}
	if req.Address != "" {
		db = db.Where("address LIKE ?", "%"+req.Address+"%")
	}
	if req.StartDate != nil {
		db = db.Where("created_at >= ?", req.StartDate)
	}
	if req.EndDate != nil {
		db = db.Where("created_at <= ?", req.EndDate)
	}

	return db
}

func (OrderRepair) Preloads(db *gorm.DB) *gorm.DB {
	db = db.Preload("Store.Superiors")
	db = db.Preload("Receptionist")
	db = db.Preload("Member")
	db = db.Preload("Products")
	db = db.Preload("Products.Product")
	db = db.Preload("Payments")
	db = db.Preload("Operator")

	return db
}

// 维修单商品
type OrderRepairProduct struct {
	SoftDelete

	Status enums.OrderRepairStatus `json:"status" gorm:"type:int(11);not NULL;comment:状态;"` // 状态

	OrderId string      `json:"order_id" gorm:"type:varchar(255);not NULL;comment:订单ID;"`  // 订单ID
	Order   OrderRepair `json:"order" gorm:"foreignKey:OrderId;references:Id;comment:订单;"` // 订单

	IsOur     bool            `json:"is_our" gorm:"type:int(11);not NULL;comment:是否本店商品;"`           // 是否本店商品
	ProductId string          `json:"product_id" gorm:"type:varchar(255);not NULL;comment:商品ID;"`    // 商品ID
	Product   ProductFinished `json:"product" gorm:"foreignKey:ProductId;references:Id;comment:商品;"` // 商品

	Code        string                `json:"code" gorm:"type:varchar(255);comment:条码;"`                   // 条码
	Name        string                `json:"name" gorm:"type:varchar(255);comment:名称;"`                   // 名称
	LabelPrice  decimal.Decimal       `json:"label_price" gorm:"type:decimal(10,2);not NULL;comment:标签价;"` // 标签价
	Brand       enums.ProductBrand    `json:"brand" gorm:"type:int(11);comment:品牌;"`                       // 品牌
	Material    enums.ProductMaterial `json:"material" gorm:"type:int(11);not NULL;comment:材质;"`           // 材质
	Quality     enums.ProductQuality  `json:"quality" gorm:"type:int(11);not NULL;comment:成色;"`            // 成色
	Gem         enums.ProductGem      `json:"gem" gorm:"type:int(11);not NULL;comment:主石;"`                // 主石
	Category    enums.ProductCategory `json:"category" gorm:"type:int(11);not NULL;comment:品类;"`           // 品类
	Craft       enums.ProductCraft    `json:"craft" gorm:"type:int(11);comment:工艺;"`                       // 工艺
	WeightMetal decimal.Decimal       `json:"weight_metal" gorm:"type:decimal(10,2);comment:金重;"`          // 金重
	WeightTotal decimal.Decimal       `json:"weight_total" gorm:"type:decimal(10,2);comment:总重;"`          // 总重
	ColorGem    enums.ProductColor    `json:"color_gem" gorm:"type:int(11);comment:颜色;"`                   // 颜色
	WeightGem   decimal.Decimal       `json:"weight_gem" gorm:"type:decimal(10,2);comment:主石重;"`           // 主石重
	Clarity     enums.ProductClarity  `json:"clarity" gorm:"type:int(11);comment:主石净度;"`                   // 主石净度
	Cut         enums.ProductCut      `json:"cut" gorm:"type:int(11);comment:主石切工;"`                       // 主石切工
	Remark      string                `json:"remark" gorm:"type:text;comment:备注;"`                         // 备注
}

func init() {
	// 注册模型
	RegisterModels(
		&OrderRepair{},
		&OrderRepairProduct{},
	)
	// 重置表
	RegisterRefreshModels(
	// &OrderRepair{},
	// &OrderRepairProduct{},
	)
}
