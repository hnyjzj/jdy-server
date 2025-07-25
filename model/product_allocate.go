package model

import (
	"jdy/enums"
	"jdy/types"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// 产品调拨单
type ProductAllocate struct {
	SoftDelete

	Method enums.ProductAllocateMethod `json:"method" gorm:"type:int(11);not NULL;comment:调拨类型;"` // 调拨类型
	Type   enums.ProductType           `json:"type" gorm:"type:int(11);not NULL;comment:产品类型;"`   // 仓库类型
	Reason enums.ProductAllocateReason `json:"reason" gorm:"type:int(11);not NULL;comment:调拨原因;"` // 调拨原因
	Status enums.ProductAllocateStatus `json:"status" gorm:"type:int(11);comment:状态;"`            // 状态
	Remark string                      `json:"remark" gorm:"type:text;comment:备注;"`               // 备注

	FromStoreId string `json:"from_store_id" gorm:"type:varchar(255);comment:调出门店;"` // 调出门店
	FromStore   *Store `json:"from_store" gorm:"foreignKey:FromStoreId;references:Id;comment:调出门店;"`
	ToStoreId   string `json:"to_store_id" gorm:"type:varchar(255);comment:调入门店;"` // 调入门店
	ToStore     *Store `json:"to_store" gorm:"foreignKey:ToStoreId;references:Id;comment:调入门店;"`

	ProductFinisheds []ProductFinished `json:"product_finisheds" gorm:"many2many:product_allocate_finished_products;comment:成品;"` // 成品
	ProductOlds      []ProductOld      `json:"product_olds" gorm:"many2many:product_allocate_old_products;comment:旧料;"`           // 旧料

	ProductCount            int64           `json:"product_count" gorm:"type:int(11);not NULL;comment:数量;"`                     // 数量
	ProductTotalWeightMetal decimal.Decimal `json:"product_total_weight_metal" gorm:"type:decimal(10,2);not NULL;comment:总重;"`  // 总重
	ProductTotalLabelPrice  decimal.Decimal `json:"product_total_label_price" gorm:"type:decimal(10,2);not NULL;comment:总标签价;"` // 总标签价
	ProductTotalAccessFee   decimal.Decimal `json:"product_total_access_fee" gorm:"type:decimal(10,2);not NULL;comment:总加工费;"`  // 总加工费

	OperatorId string `json:"operator_id" gorm:"type:varchar(255);not NULL;comment:操作人ID;"`     // 操作人ID
	Operator   *Staff `json:"operator" gorm:"foreignKey:OperatorId;references:Id;comment:操作人;"` // 操作人
	IP         string `json:"-" gorm:"type:varchar(255);not NULL;comment:IP;"`                  // IP
}

func (ProductAllocate) WhereCondition(db *gorm.DB, query *types.ProductAllocateWhere) *gorm.DB {
	if query.Id != "" {
		db = db.Where("id = ?", query.Id)
	}
	if query.Method != 0 {
		db = db.Where("method = ?", query.Method)
	}
	if query.Type != 0 {
		db = db.Where("type = ?", query.Type)
	}
	if query.Reason != 0 {
		db = db.Where("reason = ?", query.Reason)
	}
	if query.FromStoreId != "" {
		db = db.Where("from_store_id = ?", query.FromStoreId)
	}
	if query.ToStoreId != "" {
		db = db.Where("to_store_id = ?", query.ToStoreId)
	}
	if query.Status != 0 {
		db = db.Where("status = ?", query.Status)
	} else {
		db = db.Where("status IN (?)", []enums.ProductAllocateStatus{
			enums.ProductAllocateStatusDraft,
			enums.ProductAllocateStatusOnTheWay,
		})
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

func (ProductAllocate) Preloads(db *gorm.DB) *gorm.DB {
	db = db.Preload("FromStore")
	db = db.Preload("ToStore")
	db = db.Preload("Operator")

	return db
}

func init() {
	// 注册模型
	RegisterModels(
		&ProductAllocate{},
	)
	// 重置表
	RegisterRefreshModels(
	// &ProductAllocate{},
	)
}
