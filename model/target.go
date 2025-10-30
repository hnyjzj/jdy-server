package model

import (
	"fmt"
	"jdy/enums"
	"jdy/types"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Target struct {
	SoftDelete

	StoreId string `json:"store_id" gorm:"index:idx_pf_store_status_time,priority:1;type:varchar(255);comment:门店ID;"` // 门店ID
	Store   Store  `json:"store,omitzero" gorm:"foreignKey:StoreId;references:Id;comment:门店;"`                        // 门店

	Name      string    `json:"name" gorm:"type:varchar(255);comment:名称;"`                    // 名称
	IsDefault bool      `json:"is_default" gorm:"type:tinyint(1);comment:是否默认;"`              // 是否默认
	StartTime time.Time `json:"start_time" gorm:"type:datetime;not null;index;comment:开始时间;"` // 开始时间
	EndTime   time.Time `json:"end_time" gorm:"type:datetime;not null;index;comment:结束时间;"`   // 结束时间

	Scope  enums.TargetScope  `json:"scope" gorm:"type:tinyint(1);comment:统计范围;"`  // 统计范围
	Object enums.TargetObject `json:"object" gorm:"type:tinyint(1);comment:统计对象;"` // 统计对象
	Method enums.TargetMethod `json:"method" gorm:"type:tinyint(1);comment:统计方式;"` // 统计方式

	Class    []enums.ProductClassFinished `json:"class" gorm:"column:class;type:json;serializer:json;comment:产品大类;"`       // 产品大类
	Material []enums.ProductMaterial      `json:"material" gorm:"column:material;type:json;serializer:json;comment:产品材质;"` // 产品材质
	Quality  []enums.ProductQuality       `json:"quality" gorm:"column:quality;type:json;serializer:json;comment:产品成色;"`   // 产品成色
	Category []enums.ProductCategory      `json:"category" gorm:"column:category;type:json;serializer:json;comment:产品品类;"` // 产品品类
	Gem      []enums.ProductGem           `json:"gem" gorm:"column:gem;type:json;serializer:json;comment:产品主石;"`           // 产品主石
	Craft    []enums.ProductCraft         `json:"craft" gorm:"column:craft;type:json;serializer:json;comment:产品工艺;"`       // 产品工艺

	Groups    []TargetGroup    `json:"groups" gorm:"foreignKey:TargetId;references:Id;comment:分组;"`    // 分组
	Personals []TargetPersonal `json:"personals" gorm:"foreignKey:TargetId;references:Id;comment:个人;"` // 个人
}

func (Target) WhereCondition(db *gorm.DB, query *types.TargetWhere) *gorm.DB {
	if query.Id != "" {
		db = db.Where("id LIKE ?", fmt.Sprintf("%%%s%%", query.Id))
	}
	if query.Name != "" {
		db = db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", query.Name))
	}
	if query.IsDefault != nil {
		db = db.Where("is_default = ?", query.IsDefault)
	}
	if query.StartTime != nil {
		db = db.Where("start_time >= ?", query.StartTime)
	}
	if query.EndTime != nil {
		db = db.Where("end_time <= ?", query.EndTime)
	}
	if query.Method != 0 {
		db = db.Where("method = ?", query.Method)
	}
	if query.Scope != 0 {
		db = db.Where("scope = ?", query.Scope)
	}
	if query.Object != 0 {
		db = db.Where("object = ?", query.Object)
	}
	if query.StoreId != "" {
		db = db.Where("store_id = ?", query.StoreId)
	}
	if query.Class != 0 {
		db = db.Where("JSON_CONTAINS(class, ?)", fmt.Sprintf("[%d]", query.Class))
	}
	if query.Material != 0 {
		db = db.Where("JSON_CONTAINS(material, ?)", fmt.Sprintf("[%d]", query.Material))
	}
	if query.Quality != 0 {
		db = db.Where("JSON_CONTAINS(quality, ?)", fmt.Sprintf("[%d]", query.Quality))
	}
	if query.Category != 0 {
		db = db.Where("JSON_CONTAINS(category, ?)", fmt.Sprintf("[%d]", query.Category))
	}
	if query.Gem != 0 {
		db = db.Where("JSON_CONTAINS(gem, ?)", fmt.Sprintf("[%d]", query.Gem))
	}
	if query.Craft != 0 {
		db = db.Where("JSON_CONTAINS(craft, ?)", fmt.Sprintf("[%d]", query.Craft))
	}

	return db
}

func (Target) Preloads(db *gorm.DB) *gorm.DB {
	db = db.Preload("Store")
	db = db.Preload("Groups")
	db = db.Preload("Personals", func(tx *gorm.DB) *gorm.DB {
		tx = tx.Preload("Staff")
		tx = tx.Preload("Group")

		tx = tx.Order("Purpose desc")
		return tx
	})

	return db
}

type TargetGroup struct {
	SoftDelete

	TargetId string `json:"target_id" gorm:"index;type:varchar(255);comment:目标ID;"`               // 目标ID
	Target   Target `json:"target,omitzero" gorm:"foreignKey:TargetId;references:Id;comment:目标;"` // 目标

	Name string `json:"name" gorm:"type:varchar(255);comment:名称;"` // 名称

	Personals []TargetPersonal `json:"personals" gorm:"foreignKey:GroupId;references:Id;comment:个人;"` // 个人
}

func (TargetGroup) Preloads(db *gorm.DB) *gorm.DB {
	db = db.Preload("Target")
	db = db.Preload("Personals", func(tx *gorm.DB) *gorm.DB {
		tx = tx.Preload("Staff")
		tx = tx.Preload("Group")
		return tx
	})

	return db
}

type TargetPersonal struct {
	SoftDelete

	TargetId string `json:"target_id" gorm:"index;type:varchar(255);comment:目标ID;"`               // 目标ID
	Target   Target `json:"target,omitzero" gorm:"foreignKey:TargetId;references:Id;comment:目标;"` // 目标

	StaffId string `json:"staff_id" gorm:"index;type:varchar(255);comment:员工ID;"`              // 员工ID
	Staff   Staff  `json:"staff,omitzero" gorm:"foreignKey:StaffId;references:Id;comment:员工;"` // 员工

	GroupId string      `json:"group_id" gorm:"index;type:varchar(255);comment:分组ID;"`              // 分组ID
	Group   TargetGroup `json:"group,omitzero" gorm:"foreignKey:GroupId;references:Id;comment:分组;"` // 分组

	IsLeader bool            `json:"is_leader" gorm:"type:tinyint(1);comment:是否组长;"` // 是否组长
	Purpose  decimal.Decimal `json:"purpose" gorm:"type:decimal(10,2);comment:目标量;"` // 目标量

	Achieve decimal.Decimal `json:"achieved" gorm:"-"` // 达成量
}

func (TargetPersonal) Preloads(db *gorm.DB) *gorm.DB {
	db = db.Preload("Target")
	db = db.Preload("Staff")
	db = db.Preload("Group")

	return db
}

func init() {
	// 注册模型
	RegisterModels(
		&Target{},
		&TargetGroup{},
		&TargetPersonal{},
	)
	// 重置表
	RegisterRefreshModels(
	// &Target{},
	// &TargetGroup{},
	// &TargetPersonal{},
	)
}
