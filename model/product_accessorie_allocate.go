package model

import (
	"fmt"
	"jdy/enums"
	"jdy/types"

	"gorm.io/gorm"
)

// 产品调拨单
type ProductAccessorieAllocate struct {
	SoftDelete

	Method enums.ProductAllocateMethod `json:"method" gorm:"type:int(11);not NULL;comment:调拨类型;"` // 调拨类型
	Status enums.ProductAllocateStatus `json:"status" gorm:"type:int(11);comment:状态;"`            // 状态
	Remark string                      `json:"remark" gorm:"type:text;comment:备注;"`               // 备注

	FromStoreId string `json:"from_store_id" gorm:"type:varchar(255);comment:调出门店;"` // 调出门店
	FromStore   *Store `json:"from_store" gorm:"foreignKey:FromStoreId;references:Id;comment:调出门店;"`
	ToStoreId   string `json:"to_store_id" gorm:"type:varchar(255);comment:调入门店;"` // 调入门店
	ToStore     *Store `json:"to_store" gorm:"foreignKey:ToStoreId;references:Id;comment:调入门店;"`

	Products     []ProductAccessorieAllocateProduct `json:"products" gorm:"foreignKey:AllocateId;references:Id;comment:产品;"` // 产品
	ProductCount int64                              `json:"product_count" gorm:"type:int(11);not NULL;comment:入库种类数;"`       // 入库种类数
	ProductTotal int64                              `json:"product_total" gorm:"type:int(11);not NULL;comment:入库总件数;"`       // 入库总件数

	OperatorId string `json:"operator_id" gorm:"type:varchar(255);not NULL;comment:操作人ID;"`     // 操作人ID
	Operator   *Staff `json:"operator" gorm:"foreignKey:OperatorId;references:Id;comment:操作人;"` // 操作人
	IP         string `json:"-" gorm:"type:varchar(255);not NULL;comment:IP;"`                  // IP
}

func (ProductAccessorieAllocate) WhereCondition(db *gorm.DB, query *types.ProductAccessorieAllocateWhere) *gorm.DB {
	if query.Id != "" {
		db = db.Where("id = ?", query.Id)
	}
	if query.Method != 0 {
		db = db.Where("method = ?", query.Method)
	}
	if query.Status != 0 {
		db = db.Where("status = ?", query.Status)
	}
	if query.Remark != "" {
		db = db.Where("remark LIKE ?", fmt.Sprintf("%%%s%%", query.Remark))
	}
	if query.FromStoreId != "" {
		db = db.Where("from_store_id = ?", query.FromStoreId)
	}
	if query.ToStoreId != "" {
		db = db.Where("to_store_id = ?", query.ToStoreId)
	}
	if query.StartTime != nil {
		db = db.Where("created_at >= ?", query.StartTime)
	}
	if query.EndTime != nil {
		db = db.Where("created_at <= ?", query.EndTime)
	}
	if query.StoreId != "" {
		db = db.Where("from_store_id = ? OR to_store_id = ?", query.StoreId, query.StoreId)
	}

	return db
}

type ProductAccessorieAllocateProduct struct {
	ProductAccessorie

	AllocateId string                     `json:"allocate_id" gorm:"type:varchar(255);not NULL;comment:调拨单ID;"`     // 调拨单ID
	Allocate   *ProductAccessorieAllocate `json:"allocate" gorm:"foreignKey:AllocateId;references:Id;comment:调拨单;"` // 调拨单
}

func init() {
	// 注册模型
	RegisterModels(
		&ProductAccessorieAllocate{},
		&ProductAccessorieAllocateProduct{},
	)
	// 重置表
	RegisterRefreshModels(
	// &ProductAccessorieAllocate{},
	// &ProductAccessorieAllocateProduct{},
	)
}
