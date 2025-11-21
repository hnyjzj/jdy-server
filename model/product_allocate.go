package model

import (
	"fmt"
	"jdy/enums"
	"jdy/types"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// 产品调拨单
type ProductAllocate struct {
	SoftDelete

	Method enums.ProductAllocateMethod `json:"method" gorm:"index;type:int(11);not NULL;comment:调拨类型;"` // 调拨类型
	Type   enums.ProductType           `json:"type" gorm:"index;type:int(11);not NULL;comment:产品类型;"`   // 仓库类型
	Reason enums.ProductAllocateReason `json:"reason" gorm:"index;type:int(11);not NULL;comment:调拨原因;"` // 调拨原因
	Status enums.ProductAllocateStatus `json:"status" gorm:"index;type:int(11);comment:状态;"`            // 状态
	Remark string                      `json:"remark" gorm:"type:text;comment:备注;"`                     // 备注

	FromStoreId string `json:"from_store_id" gorm:"index;type:varchar(255);comment:调出门店;"` // 调出门店
	FromStore   *Store `json:"from_store" gorm:"foreignKey:FromStoreId;references:Id;comment:调出门店;"`
	ToStoreId   string `json:"to_store_id" gorm:"index;type:varchar(255);comment:调入门店;"` // 调入门店
	ToStore     *Store `json:"to_store" gorm:"foreignKey:ToStoreId;references:Id;comment:调入门店;"`

	ProductFinisheds []ProductFinished `json:"product_finisheds" gorm:"many2many:product_allocate_finished_products;comment:成品;"` // 成品
	ProductOlds      []ProductOld      `json:"product_olds" gorm:"many2many:product_allocate_old_products;comment:旧料;"`           // 旧料
	Product          any               `json:"product" gorm:"-"`                                                                  // 产品

	ProductCount             int64           `json:"product_count" gorm:"type:int(11);not NULL;default:0;comment:数量;"`                        // 数量
	ProductTotalWeightMetal  decimal.Decimal `json:"product_total_weight_metal" gorm:"type:decimal(15,4);not NULL;default:0;comment:总金重;"`    // 总金重
	ProductTotalLabelPrice   decimal.Decimal `json:"product_total_label_price" gorm:"type:decimal(15,4);not NULL;default:0;comment:总标签价;"`    // 总标签价
	ProductTotalAccessFee    decimal.Decimal `json:"product_total_access_fee" gorm:"type:decimal(15,4);not NULL;default:0;comment:总入网费;"`     // 总入网费
	ProductTotalRecyclePrice decimal.Decimal `json:"product_total_recycle_price" gorm:"type:decimal(15,4);not NULL;default:0;comment:总回收金额;"` // 总回收金额

	OperatorId string `json:"operator_id" gorm:"type:varchar(255);NULL;comment:操作人ID;"`         // 操作人ID
	Operator   *Staff `json:"operator" gorm:"foreignKey:OperatorId;references:Id;comment:操作人;"` // 操作人
	IP         string `json:"-" gorm:"type:varchar(255);NULL;comment:IP;"`                      // IP

	InitiatorId string `json:"initiator_id" gorm:"index;type:varchar(255);NULL;comment:发起人ID;"`    // 发起人ID
	InitiatorIP string `json:"-" gorm:"type:varchar(255);NULL;comment:发起人IP;"`                     // 发起人IP
	Initiator   *Staff `json:"initiator" gorm:"foreignKey:InitiatorId;references:Id;comment:发起人;"` // 发起人

	ReceiverId string `json:"receiver_id" gorm:"index;type:varchar(255);NULL;comment:接收人ID;"`   // 接收人ID
	ReceiverIP string `json:"-" gorm:"type:varchar(255);NULL;comment:接收人IP;"`                   // 接收人IP
	Receiver   *Staff `json:"receiver" gorm:"foreignKey:ReceiverId;references:Id;comment:接收人;"` // 接收人
}

func (ProductAllocate) WhereCondition(db *gorm.DB, req *types.ProductAllocateWhere, staff *Staff) *gorm.DB {
	if req.Id != "" {
		db = db.Where("id LIKE ?", fmt.Sprintf("%%%s%%", req.Id))
	}
	if req.Method != 0 {
		db = db.Where("method = ?", req.Method)
	}
	if req.Type != 0 {
		db = db.Where("type = ?", req.Type)
	}
	if req.Reason != 0 {
		db = db.Where("reason = ?", req.Reason)
	}
	switch {
	case req.FromStoreId == "" && req.ToStoreId == "" && req.StoreId != "":
		{
			db = db.Where("(from_store_id = ? OR to_store_id = ?)", req.StoreId, req.StoreId)
		}
	case req.FromStoreId != "" && req.ToStoreId != "":
		{
			db = db.Where("(from_store_id = ? AND to_store_id = ?)", req.FromStoreId, req.ToStoreId)
		}
	case req.FromStoreId != "" && req.ToStoreId == "":
		{
			db = db.Where("from_store_id = ?", req.FromStoreId)
		}
	case req.FromStoreId == "" && req.ToStoreId != "":
		{
			db = db.Where("to_store_id = ?", req.ToStoreId)
		}
	default:
		{
			db = db.Where("(from_store_id IN (?) OR to_store_id IN (?))", staff.StoreIds, staff.StoreIds)
		}
	}
	if req.Status != 0 {
		db = db.Where("status = ?", req.Status)
	}
	if req.StartTime != nil {
		db = db.Where("created_at >= ?", req.StartTime)
	}
	if req.EndTime != nil {
		db = db.Where("created_at <= ?", req.EndTime)
	}
	if req.InitiatorId != "" {
		db = db.Where("initiator_id = ?", req.InitiatorId)
	}
	if req.ReceiverId != "" {
		db = db.Where("receiver_id = ?", req.ReceiverId)
	}

	return db
}

func (ProductAllocate) Preloads(db *gorm.DB) *gorm.DB {
	db = db.Preload("FromStore")
	db = db.Preload("ToStore")
	db = db.Preload("Operator")
	db = db.Preload("Initiator")
	db = db.Preload("Receiver")

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
