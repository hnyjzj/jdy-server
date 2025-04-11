package model

import (
	"jdy/enums"
	"jdy/types"

	"gorm.io/gorm"
)

// 产品历史记录
type ProductHistory struct {
	Model

	Type   enums.ProductType   `json:"type" gorm:"type:tinyint(2);comment:产品类型;"` // 产品类型
	Action enums.ProductAction `json:"action" gorm:"type:tinyint(2);comment:操作;"` // 操作

	NewValue any `json:"new_value" gorm:"type:text;serializer:json;comment:值;"`  // 值
	OldValue any `json:"old_value" gorm:"type:text;serializer:json;comment:旧值;"` // 旧值

	ProductId string `json:"product_id" gorm:"type:varchar(255);not NULL;comment:产品ID;"` // 产品ID

	StoreId string `json:"store_id" gorm:"type:varchar(255);not NULL;comment:门店ID;"`  // 门店ID
	Store   *Store `json:"store" gorm:"foreignKey:StoreId;references:Id;comment:门店;"` // 门店

	SourceId string `json:"source_id" gorm:"column:source_id;size:255;not NULL;comment:来源id;"` // 来源id

	Reason string `json:"reason" gorm:"type:varchar(255);comment:原因;"` // 原因

	OperatorId string `json:"operator_id" gorm:"type:varchar(255);not NULL;comment:操作人ID;"`     // 操作人ID
	Operator   *Staff `json:"operator" gorm:"foreignKey:OperatorId;references:Id;comment:操作人;"` // 操作人
	IP         string `json:"ip" gorm:"type:varchar(255);not NULL;comment:IP;"`                 // IP
}

func (ProductHistory) WhereCondition(db *gorm.DB, query *types.ProductHistoryWhereReq) *gorm.DB {
	if len(query.Type) > 0 {
		db = db.Where("type in (?)", query.Type)
	}
	if query.ProductId != "" {
		db = db.Where("product_id = ?", query.ProductId)
	}
	if query.StoreId != "" {
		db = db.Where("store_id = ?", query.StoreId)
	}
	if query.Action != 0 {
		db = db.Where("action = ?", query.Action)
	}

	return db
}

func init() {
	// 注册模型
	RegisterModels(
		&ProductHistory{},
	)
	// 重置表
	RegisterRefreshModels(
	// &ProductHistory{},
	)
}
