package model

import (
	"fmt"
	"jdy/types"

	"gorm.io/gorm"
)

type Store struct {
	SoftDelete

	Name     string `json:"name" gorm:"size:255;comment:名称"`        // 名称
	ParentId string `json:"parent_id" gorm:"size:255;comment:父级ID"` // 父级ID
	Order    int    `json:"order" gorm:"comment:排序"`                // 排序

	Logo     string `json:"logo" gorm:"size:255;comment:logo"`    // logo
	Contact  string `json:"contact" gorm:"size:255;comment:联系方式"` // 联系方式
	Province string `json:"province" gorm:"size:255;comment:省份"`  // 省份
	City     string `json:"city" gorm:"size:255;comment:城市"`      // 城市
	District string `json:"district" gorm:"size:255;comment:区域"`  // 区域
	Address  string `json:"address" gorm:"size:500;comment:地址"`   // 地址

	Staffs    []Staff `json:"staffs" gorm:"many2many:stores_staffs;"`       // 员工
	Superiors []Staff `json:"superiors" gorm:"many2many:stores_superiors;"` // 负责人
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
