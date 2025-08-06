package model

import (
	"jdy/enums"
	"jdy/types"

	"gorm.io/gorm"
)

// 产品历史记录
type ProductHistory struct {
	Model

	Type   enums.ProductType   `json:"type" gorm:"type:int(11);comment:产品类型;"` // 产品类型
	Action enums.ProductAction `json:"action" gorm:"type:int(11);comment:操作;"` // 操作

	OldValue any `json:"old_value" gorm:"type:text;serializer:json;comment:旧值;"` // 旧值
	NewValue any `json:"new_value" gorm:"type:text;serializer:json;comment:新值;"` // 新值

	ProductId string `json:"product_id" gorm:"type:varchar(255);not NULL;comment:产品ID;"` // 产品ID

	StoreId string `json:"store_id" gorm:"type:varchar(255);not NULL;comment:门店ID;"`  // 门店ID
	Store   *Store `json:"store" gorm:"foreignKey:StoreId;references:Id;comment:门店;"` // 门店

	SourceId string `json:"source_id" gorm:"column:source_id;size:255;not NULL;comment:来源id;"` // 来源id

	Reason string `json:"reason" gorm:"type:varchar(255);comment:原因;"` // 原因

	OperatorId string `json:"operator_id" gorm:"type:varchar(255);not NULL;comment:操作人ID;"`     // 操作人ID
	Operator   *Staff `json:"operator" gorm:"foreignKey:OperatorId;references:Id;comment:操作人;"` // 操作人
	IP         string `json:"ip" gorm:"type:varchar(255);not NULL;comment:IP;"`                 // IP
}

func (ProductHistory) WhereCondition(db *gorm.DB, query *types.ProductHistoryWhere) *gorm.DB {
	db = db.Where("type <> ?", enums.ProductTypeAccessorie)

	if query.Type != 0 {
		db = db.Where("type = ?", query.Type)
	}
	if query.Code != "" {
		// 使用 JSON 函数更安全和高效
		db = db.Where(
			"JSON_EXTRACT(new_value, '$.code') = ? OR JSON_EXTRACT(old_value, '$.code') = ?",
			query.Code,
			query.Code,
		).Debug()
	}
	if query.StoreId != "" {
		db = db.Where("store_id = ?", query.StoreId)
	}
	if query.Action != 0 {
		db = db.Where("action = ?", query.Action)
	}

	if query.StartTime != nil {
		db = db.Where("created_at >= ?", query.StartTime)
	}
	if query.EndTime != nil {
		db = db.Where("created_at <= ?", query.EndTime)
	}

	return db
}

func (ProductHistory) WhereAccessorieCondition(db *gorm.DB, query *types.ProductAccessorieHistoryWhere) *gorm.DB {
	db = db.Where("type = ?", enums.ProductTypeAccessorie)

	if query.Name != "" {
		db = db.Where("new_value LIKE ? OR old_value LIKE ?", "%\"name\":\""+query.Name+"\"%", "%\"name\":\""+query.Name+"\"%")
	}
	if query.Code != "" {
		db = db.Where("new_value LIKE ? OR old_value LIKE ?", "%\"code\":\""+query.Code+"\"%", "%\"code\":\""+query.Code+"\"%")
	}
	if query.StoreId != "" {
		db = db.Where("store_id = ?", query.StoreId)
	}
	if query.Action != 0 {
		db = db.Where("action = ?", query.Action)
	}

	if query.StartTime != nil {
		db = db.Where("created_at >= ?", query.StartTime)
	}
	if query.EndTime != nil {
		db = db.Where("created_at <= ?", query.EndTime)
	}

	return db
}

func (ProductHistory) Preloads(db *gorm.DB) *gorm.DB {
	db = db.Preload("Store")
	db = db.Preload("Operator")

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
