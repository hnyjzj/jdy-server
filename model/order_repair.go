package model

import (
	"jdy/enums"

	"github.com/shopspring/decimal"
)

// 维修单
type OrderRepair struct {
	SoftDelete

	StoreId string `json:"store_id" gorm:"type:varchar(255);not NULL;comment:门店ID;"`  // 门店ID
	Store   Store  `json:"store" gorm:"foreignKey:StoreId;references:Id;comment:门店;"` // 门店

	Status enums.OrderStatus `json:"status" gorm:"type:tinyint(2);not NULL;comment:订单状态;"` // 订单状态

	ReceptionistId string `json:"receptionist_id" gorm:"type:varchar(255);not NULL;comment:接待人ID;"`         // 接待人ID
	Receptionist   Staff  `json:"receptionist" gorm:"foreignKey:ReceptionistId;references:Id;comment:接待人;"` // 接待人

	MemberId string `json:"member_id" gorm:"type:varchar(255);not NULL;comment:会员ID;"`   // 会员ID
	Member   Member `json:"member" gorm:"foreignKey:MemberId;references:Id;comment:会员;"` // 会员

	Name string `json:"name" gorm:"type:varchar(255);not NULL;comment:维修项目;"` // 维修项目
	Desc string `json:"desc" gorm:"type:text;not NULL;comment:描述;"`           // 问题描述

	DeliveryMethod enums.DeliveryMethod `json:"delivery_method" gorm:"type:tinyint(2);not NULL;comment:取货方式;"` // 取货方式
	Province       string               `json:"province" gorm:"type:varchar(255);not NULL;comment:省;"`         // 省
	City           string               `json:"city" gorm:"type:varchar(255);not NULL;comment:市;"`             // 市
	District       string               `json:"district" gorm:"type:varchar(255);not NULL;comment:区;"`         // 区
	Address        string               `json:"address" gorm:"type:varchar(255);not NULL;comment:地址;"`         // 地址

	Products []OrderRepairProduct `json:"products" gorm:"foreignKey:OrderId;references:Id;comment:维修单商品;"` // 维修单商品

	Images []string `json:"images" gorm:"type:text;not NULL;serializer:json;comment:图片;"` // 图片

	Payments []OrderPayment `json:"payments" gorm:"foreignKey:OrderId;references:Id;comment:支付信息;"` // 支付信息

	OperatorId string `json:"operator_id" gorm:"type:varchar(255);not NULL;comment:操作员ID;"`     // 操作员ID
	Operator   Staff  `json:"operator" gorm:"foreignKey:OperatorId;references:Id;comment:操作员;"` // 操作员
	IP         string `json:"ip" gorm:"type:varchar(255);not NULL;comment:IP地址;"`               // IP地址
}

// func (OrderRepair) WhereCondition(db *gorm.DB, req *types.OrderRepairWhere) *gorm.DB {

// 	return db
// }

// 维修单商品
type OrderRepairProduct struct {
	SoftDelete

	Status enums.OrderStatus `json:"status" gorm:"type:tinyint(1);not NULL;comment:状态;"` // 状态

	OrderId string      `json:"order_id" gorm:"type:varchar(255);not NULL;comment:订单ID;"`  // 订单ID
	Order   OrderRepair `json:"order" gorm:"foreignKey:OrderId;references:Id;comment:订单;"` // 订单

	IsOur     bool            `json:"is_our" gorm:"type:tinyint(1);not NULL;comment:是否本店商品;"`          // 是否本店商品
	ProductId string          `json:"product_id" gorm:"type:varchar(255);not NULL;comment:商品ID;"`      // 商品ID
	Product   ProductFinished `json:"product" gorm:"type:text;not NULL;serializer:json;comment:货品信息;"` // 货品信息

	Expense decimal.Decimal `json:"expense" gorm:"type:decimal(10,2);not NULL;comment:维修费;"` // 维修费
	Cost    decimal.Decimal `json:"cost" gorm:"type:decimal(10,2);not NULL;comment:维修成本;"`   // 维修成本
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
