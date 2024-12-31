package model

import (
	"fmt"
	"jdy/types"

	"gorm.io/gorm"
)

type Store struct {
	SoftDelete

	Name     string `json:"name" gorm:"size:255;comment:名称"`
	Address  string `json:"address" gorm:"size:500;comment:地址"`
	Contact  string `json:"contact" gorm:"size:255;comment:联系方式"`
	Logo     string `json:"logo" gorm:"size:255;comment:logo"`
	Sort     int    `json:"sort" gorm:"size:10;comment:排序"`
	Province string `json:"province" gorm:"size:255;comment:省份"`
	City     string `json:"city" gorm:"size:255;comment:城市"`
	District string `json:"district" gorm:"size:255;comment:区域"`

	Children []*Store `json:"children,omitempty" gorm:"-"`

	Staffs []Staff `json:"staffs" gorm:"many2many:stores_staffs;"`
}

func (Store) WhereCondition(db *gorm.DB, query *types.StoreWhere) *gorm.DB {
	if query.Name != nil {
		db = db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", *query.Name))
	}
	if query.Address != "" {
		db = db.Where("address LIKE ?", fmt.Sprintf("%%%s%%", query.Address))
	}
	if query.Contact != "" {
		db = db.Where("contact LIKE ?", fmt.Sprintf("%%%s%%", query.Contact))
	}
	if query.Region.Province != nil {
		db = db.Where("province LIKE ?", fmt.Sprintf("%%%s%%", *query.Region.Province))
	}
	if query.Region.City != nil {
		db = db.Where("city LIKE ?", fmt.Sprintf("%%%s%%", *query.Region.City))
	}
	if query.Region.District != nil {
		db = db.Where("district LIKE ?", fmt.Sprintf("%%%s%%", *query.Region.District))
	}

	return db
}

func init() {
	// 注册模型
	RegisterModels(
		&Store{},
	)
	// 重置表
	RegisterRefreshModels(
	// &Store{},
	)
}
