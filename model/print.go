package model

import (
	"fmt"
	"jdy/enums"
	"jdy/types"

	"gorm.io/gorm"
)

type Print struct {
	SoftDelete

	StoreId string `json:"store_id" gorm:"index;column:store_id;type:varchar(255);not null;comment:店铺id"` // 店铺id
	Store   Store  `json:"store" gorm:"foreignKey:StoreId;references:Id"`

	Name      string            `json:"name" gorm:"column:name;type:varchar(255);not null;comment:模板名称"`          // 模板名称
	Type      enums.PrintType   `json:"type" gorm:"column:type;type:int(11);not null;comment:模板类型"`               // 模板类型
	Config    types.PrintConfig `json:"config" gorm:"column:config;type:json;serializer:json;comment:模板配置"`       // 模板配置
	IsDefault bool              `json:"is_default" gorm:"column:is_default;type:int(11);not null;comment:是否默认模板"` // 是否默认模板
}

func (Print) Default(t enums.PrintType) Print {
	return Print{
		IsDefault: true,
		Name:      "系统默认",
		Type:      t,
		Config:    types.PrintConfig{}.Default(),
	}
}

func (Print) WhereCondition(db *gorm.DB, req *types.PrintWhere) *gorm.DB {
	if req.Id != "" {
		db = db.Where("id LIKE ?", fmt.Sprintf("%%%s%%", req.Id))
	}
	if req.StoreId != "" {
		db = db.Where("store_id = ?", req.StoreId)
	}
	if req.Type != 0 {
		db = db.Where("type = ?", req.Type)
	}
	if req.Name != "" {
		db = db.Where("name like ?", fmt.Sprintf("%%%s%%", req.Name))
	}

	return db
}

func init() {
	// 注册模型
	RegisterModels(
		&Print{},
	)
	// 重置表
	RegisterRefreshModels(
	// &Print{},
	)
}
