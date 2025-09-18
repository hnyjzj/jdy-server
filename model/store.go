package model

import (
	"fmt"
	"jdy/enums"
	"jdy/types"
	"strings"

	"gorm.io/gorm"
)

type Store struct {
	SoftDelete

	IdWx  string `json:"id_wx" gorm:"index;size:255;comment:微信ID"` // 微信ID
	Name  string `json:"name" gorm:"index;size:255;comment:名称"`    // 名称
	Alias string `json:"alias" gorm:"index;size:255;comment:别名"`   // 别名
	Order int    `json:"order" gorm:"index;comment:排序"`            // 排序

	RegionId string `json:"region_id" gorm:"index;size:255;comment:区域ID"`               // 区域ID
	Region   Region `json:"region" gorm:"foreignKey:RegionId;references:Id;comment:区域"` // 区域

	Staffs    []Staff `json:"staffs" gorm:"many2many:store_staffs;"`       // 员工
	Superiors []Staff `json:"superiors" gorm:"many2many:store_superiors;"` // 负责人
	Admins    []Staff `json:"admins" gorm:"many2many:store_admins;"`       // 管理员
}

func (Store) WhereCondition(db *gorm.DB, query *types.StoreWhere) *gorm.DB {
	if query.Name != "" {
		db = db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", query.Name))
	}
	if query.Alias != "" {
		db = db.Where("alias = ?", query.Alias)
	}
	if query.RegionId != "" {
		db = db.Where("region_id = ?", query.RegionId)
	}

	return db
}

func (Store) Preloads(db *gorm.DB) *gorm.DB {
	db = db.Preload("Region")
	db = db.Preload("Staffs")
	db = db.Preload("Superiors")
	db = db.Preload("Admins")

	return db
}

func (Store) Default(identity enums.Identity) *Store {
	if identity < enums.IdentityAdmin {
		return nil
	}
	def := &Store{
		Name:  "全部",
		Alias: "全部",
	}

	return def
}

const StorePrefix = "店"
const RegionPrefix = "区域"
const HeaderquartersPrefix = "总部"

func (store *Store) IsHeadquarters() bool {
	if store == nil {
		return false
	}

	return strings.HasSuffix(store.Name, HeaderquartersPrefix)
}

func (store *Store) InStore(staff_id string) bool {
	for _, staff := range store.Staffs {
		if staff.Id == staff_id {
			return true
		}
	}
	for _, staff := range store.Superiors {
		if staff.Id == staff_id {
			return true
		}
	}
	for _, staff := range store.Admins {
		if staff.Id == staff_id {
			return true
		}
	}

	return false
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
