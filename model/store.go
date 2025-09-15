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

	IdWx  string `json:"id_wx" gorm:"size:255;comment:微信ID"` // 微信ID
	Name  string `json:"name" gorm:"size:255;comment:名称"`    // 名称
	Order int    `json:"order" gorm:"comment:排序"`            // 排序

	Logo     string `json:"logo" gorm:"size:255;comment:logo"`    // logo
	Contact  string `json:"contact" gorm:"size:255;comment:联系方式"` // 联系方式
	Province string `json:"province" gorm:"size:255;comment:省份"`  // 省份
	City     string `json:"city" gorm:"size:255;comment:城市"`      // 城市
	District string `json:"district" gorm:"size:255;comment:区域"`  // 区域
	Address  string `json:"address" gorm:"size:500;comment:地址"`   // 地址

	Staffs    []Staff `json:"staffs" gorm:"many2many:store_staffs;"`       // 员工
	Superiors []Staff `json:"superiors" gorm:"many2many:store_superiors;"` // 负责人
}

func (Store) WhereCondition(db *gorm.DB, query *types.StoreWhere) *gorm.DB {
	if query.Name != "" {
		db = db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", query.Name))
	}
	if query.Address != "" {
		db = db.Where("address LIKE ?", fmt.Sprintf("%%%s%%", query.Address))
	}
	if query.Contact != "" {
		db = db.Where("contact LIKE ?", fmt.Sprintf("%%%s%%", query.Contact))
	}
	if query.Field.Province != nil {
		db = db.Where("province LIKE ?", fmt.Sprintf("%%%s%%", *query.Field.Province))
	}
	if query.Field.City != nil {
		db = db.Where("city LIKE ?", fmt.Sprintf("%%%s%%", *query.Field.City))
	}
	if query.Field.District != nil {
		db = db.Where("district LIKE ?", fmt.Sprintf("%%%s%%", *query.Field.District))
	}
	if query.RegionId != "" {
		db = db.Where("id IN (SELECT store_id FROM region_stores WHERE region_id = ?)", query.RegionId)
	}

	return db
}

func (Store) Preloads(db *gorm.DB) *gorm.DB {
	db = db.Preload("Staffs")
	db = db.Preload("Superiors")

	return db
}

func (Store) Default(identity enums.Identity) *Store {
	if identity < enums.IdentityAdmin {
		return nil
	}
	def := &Store{
		Name: "全部",
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

	return strings.Contains(store.Name, HeaderquartersPrefix)
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
