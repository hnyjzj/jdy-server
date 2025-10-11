package model

import (
	"fmt"
	"jdy/enums"
	"jdy/types"

	"gorm.io/gorm"
)

type Region struct {
	SoftDelete

	IdWx  string `json:"id_wx" gorm:"index;size:255;comment:微信ID"`    // 微信ID
	Name  string `json:"name" gorm:"uniqueIndex;size:255;comment:名称"` // 名称
	Alias string `json:"alias" gorm:"index;size:255;comment:别名"`      // 别名
	Order int    `json:"order" gorm:"comment:排序"`                     // 排序

	Stores    []Store `json:"stores" gorm:"foreignKey:RegionId;references:Id;comment:门店"` // 门店
	Staffs    []Staff `json:"staffs" gorm:"many2many:region_staffs;"`                     // 员工
	Superiors []Staff `json:"superiors" gorm:"many2many:region_superiors;"`               // 负责人
	Admins    []Staff `json:"admins" gorm:"many2many:region_admins;"`                     // 管理员
}

func (Region) WhereCondition(db *gorm.DB, query *types.RegionWhere) *gorm.DB {
	if query.Name != "" {
		db = db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", query.Name))
	}
	if query.Alias != "" {
		db = db.Where("alias = ?", query.Alias)
	}

	return db
}

func (Region) Preloads(db *gorm.DB) *gorm.DB {
	db = db.Preload("Stores")
	db = db.Preload("Staffs")
	db = db.Preload("Superiors")
	db = db.Preload("Admins")

	return db
}

func (region *Region) InRegion(staff_id string) bool {
	for _, staff := range region.Staffs {
		if staff.Id == staff_id {
			return true
		}
	}
	for _, staff := range region.Superiors {
		if staff.Id == staff_id {
			return true
		}
	}
	for _, staff := range region.Admins {
		if staff.Id == staff_id {
			return true
		}
	}

	return false
}

func (Region) Default(identity enums.Identity) *Region {
	if identity < enums.IdentityAdmin {
		return nil
	}
	def := &Region{
		Name:  "全部",
		Alias: "全部",
	}

	return def
}

func init() {
	// 注册模型
	RegisterModels(
		&Region{},
	)
	// 重置表
	RegisterRefreshModels(
	// &Region{},
	)
}
