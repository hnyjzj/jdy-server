package model

import (
	"fmt"
	"jdy/enums"
	"jdy/types"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// 产品调拨单
type ProductAccessorieAllocate struct {
	SoftDelete

	Method enums.ProductAccessorieAllocateMethod `json:"method" gorm:"type:int(11);not NULL;comment:调拨类型;"` // 调拨类型
	Status enums.ProductAllocateStatus           `json:"status" gorm:"type:int(11);comment:状态;"`            // 状态
	Remark string                                `json:"remark" gorm:"type:text;comment:备注;"`               // 备注

	FromStoreId string  `json:"from_store_id" gorm:"type:varchar(255);comment:调出门店;"` // 调出门店
	FromStore   *Store  `json:"from_store" gorm:"foreignKey:FromStoreId;references:Id;comment:调出门店;"`
	ToStoreId   string  `json:"to_store_id" gorm:"type:varchar(255);comment:调入门店;"` // 调入门店
	ToStore     *Store  `json:"to_store" gorm:"foreignKey:ToStoreId;references:Id;comment:调入门店;"`
	ToRegionId  string  `json:"to_region_id" gorm:"type:varchar(255);comment:调入区域;"` // 调入区域
	ToRegion    *Region `json:"to_region" gorm:"foreignKey:ToRegionId;references:Id;comment:调入区域;"`

	Products     []ProductAccessorieAllocateProduct `json:"products" gorm:"foreignKey:AllocateId;references:Id;comment:产品;"` // 产品
	ProductCount int64                              `json:"product_count" gorm:"type:int(11);not NULL;comment:入库种类数;"`       // 入库种类数
	ProductTotal int64                              `json:"product_total" gorm:"type:int(11);not NULL;comment:入库总件数;"`       // 入库总件数

	OperatorId string `json:"operator_id" gorm:"type:varchar(255);NULL;comment:操作人ID;"`         // 操作人ID
	Operator   *Staff `json:"operator" gorm:"foreignKey:OperatorId;references:Id;comment:操作人;"` // 操作人
	IP         string `json:"-" gorm:"type:varchar(255);NULL;comment:IP;"`                      // IP

	InitiatorId string `json:"initiator_id" gorm:"type:varchar(255);NULL;comment:发起人ID;"`          // 发起人ID
	InitiatorIP string `json:"initiator_ip" gorm:"type:varchar(255);NULL;comment:发起人IP;"`          // 发起人IP
	Initiator   *Staff `json:"initiator" gorm:"foreignKey:InitiatorId;references:Id;comment:发起人;"` // 发起人

	ReceiverId string `json:"receiver_id" gorm:"type:varchar(255);NULL;comment:接收人ID;"`         // 接收人ID
	ReceiverIP string `json:"receiver_ip" gorm:"type:varchar(255);NULL;comment:接收人IP;"`         // 接收人IP
	Receiver   *Staff `json:"receiver" gorm:"foreignKey:ReceiverId;references:Id;comment:接收人;"` // 接收人
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
	if query.ToRegionId != "" {
		db = db.Where("to_region_id = ?", query.ToRegionId)
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

func (ProductAccessorieAllocate) Preloads(db *gorm.DB, pages *types.PageReq) *gorm.DB {
	db = db.Preload("FromStore").Preload("ToStore").Preload("ToRegion")
	if pages != nil {
		db = db.Preload("Products", func(tx *gorm.DB) *gorm.DB {
			pdb := tx
			if pages != nil {
				pdb = PageCondition(pdb, &types.PageReq{Page: pages.Page, Limit: pages.Limit})
			}

			return pdb
		})
	}
	db = db.Preload("Operator")
	db = db.Preload("Initiator")
	db = db.Preload("Receiver")

	return db
}

type ProductAccessorieAllocateProduct struct {
	SoftDelete

	Name       string                            `json:"name" gorm:"type:varchar(255);uniqueIndex:idx_allocate_name;comment:名称;"` // 名称
	Type       enums.ProductAccessorieType       `json:"type" gorm:"type:int(11);comment:配件类型;"`                                  // 配件类型
	RetailType enums.ProductAccessorieRetailType `json:"retail_type" gorm:"type:int(11);not NULL;comment:零售方式;"`                  // 零售方式
	Price      decimal.Decimal                   `json:"price" gorm:"type:decimal(10,2);comment:单价;"`                             // 单价
	Remark     string                            `json:"remark" gorm:"type:text;comment:备注;"`                                     // 备注
	Stock      int64                             `json:"stock" gorm:"type:int(9);default:1;comment:库存;"`                          // 库存

	Status enums.ProductAccessorieStatus `json:"status" gorm:"type:int(11);default:1;comment:状态;"` // 状态

	AllocateId string                     `json:"allocate_id" gorm:"type:varchar(255);uniqueIndex:idx_allocate_name;not NULL;comment:调拨单ID;"` // 调拨单ID
	Allocate   *ProductAccessorieAllocate `json:"allocate" gorm:"foreignKey:AllocateId;references:Id;comment:调拨单;"`                           // 调拨单
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
