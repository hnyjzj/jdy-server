package model

import (
	"jdy/types"

	"gorm.io/gorm"
)

type Remark struct {
	Model

	Content string `json:"content" gorm:"type:text;not null;comment:评论内容;"` // 评论内容

	StoreId string `json:"store_id" gorm:"index;type:varchar(255);not NULL;comment:门店ID;"` // 门店ID
	Store   Store  `json:"store" gorm:"foreignKey:StoreId;references:Id;comment:门店;"`      // 门店

	OperatorId string `json:"operator_id" gorm:"type:varchar(255);not NULL;comment:操作员ID;"`     // 操作员ID
	Operator   Staff  `json:"operator" gorm:"foreignKey:OperatorId;references:Id;comment:操作员;"` // 操作员
	IP         string `json:"ip" gorm:"type:varchar(255);not NULL;comment:IP地址;"`               // IP地址
}

func (Remark) WhereCondition(db *gorm.DB, query *types.RemarkWhere) *gorm.DB {
	if query.Id != "" {
		db = db.Where("id = ?", query.Id)
	}
	if query.StoreId != "" {
		db = db.Where("store_id = ?", query.StoreId)
	}
	if query.Content != "" {
		db = db.Where("content LIKE ?", "%"+query.Content+"%")
	}

	return db
}

func (Remark) Preloads(db *gorm.DB) *gorm.DB {
	db = db.Preload("Store")
	db = db.Preload("Operator")

	return db
}

func init() {
	// 注册模型
	RegisterModels(
		&Remark{},
	)
	// 重置表
	RegisterRefreshModels(
	// &Remark{},
	)
}
