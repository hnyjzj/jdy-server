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

	Method enums.TargetMethod `json:"method" gorm:"type:tinyint(1);comment:统计方式;"` // 统计方式
	Scope  enums.TargetScope  `json:"scope" gorm:"type:tinyint(1);comment:统计范围;"`  // 统计范围
	Object enums.TargetObject `json:"object" gorm:"type:tinyint(1);comment:统计对象;"` // 统计对象

	Class    []enums.ProductClassFinished `json:"class" gorm:"column:class;type:text;serializer:json;comment:产品大类;"`       // 产品大类
	Material []enums.ProductMaterial      `json:"material" gorm:"column:material;type:text;serializer:json;comment:产品材质;"` // 产品材质
	Quality  []enums.ProductQuality       `json:"quality" gorm:"column:quality;type:text;serializer:json;comment:产品成色;"`   // 产品成色
	Category []enums.ProductCategory      `json:"category" gorm:"column:category;type:text;serializer:json;comment:产品品类;"` // 产品品类
	Gem      []enums.ProductGem           `json:"gem" gorm:"column:gem;type:text;serializer:json;comment:产品主石;"`           // 产品主石
	Craft    []enums.ProductCraft         `json:"craft" gorm:"column:craft;type:text;serializer:json;comment:产品工艺;"`       // 产品工艺

	Groups    []TargetGroup    `json:"groups" gorm:"foreignKey:TargetId;references:Id;comment:分组;"`    // 分组
	Personals []TargetPersonal `json:"personals" gorm:"foreignKey:TargetId;references:Id;comment:个人;"` // 个人
	Achieves  []TargetAchieve  `json:"achieves" gorm:"foreignKey:TargetId;references:Id;comment:达成;"`  // 达成
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
		db = db.Where("FIND_IN_SET(?, class)", query.Class)
	}
	if query.Material != 0 {
		db = db.Where("FIND_IN_SET(?, material)", query.Material)
	}
	if query.Quality != 0 {
		db = db.Where("FIND_IN_SET(?, quality)", query.Quality)
	}
	if query.Category != 0 {
		db = db.Where("FIND_IN_SET(?, category)", query.Category)
	}
	if query.Gem != 0 {
		db = db.Where("FIND_IN_SET(?, gem)", query.Gem)
	}
	if query.Craft != 0 {
		db = db.Where("FIND_IN_SET(?, craft)", query.Craft)
	}

	return db
}

func (Target) Preloads(db *gorm.DB) *gorm.DB {
	db = db.Preload("Store")
	db = db.Preload("Groups")
	db = db.Preload("Personals", func(tx *gorm.DB) *gorm.DB {
		tx = tx.Order("Achieve desc")
		tx = tx.Order("Purpose desc")

		tx = tx.Preload("Staff")
		tx = tx.Preload("Group")
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

	IsLeader bool            `json:"is_leader" gorm:"type:tinyint(1);comment:是否组长;"`  // 是否组长
	Purpose  decimal.Decimal `json:"purpose" gorm:"type:decimal(10,2);comment:目标量;"`  // 目标量
	Achieve  decimal.Decimal `json:"achieved" gorm:"type:decimal(10,2);comment:达成量;"` // 达成量
}

func (TargetPersonal) Preloads(db *gorm.DB) *gorm.DB {
	db = db.Preload("Target")
	db = db.Preload("Staff")
	db = db.Preload("Group")

	return db
}

type TargetAchieve struct {
	SoftDelete

	TargetId string `json:"target_id" gorm:"index;type:varchar(255);comment:目标ID;"`               // 目标ID
	Target   Target `json:"target,omitzero" gorm:"foreignKey:TargetId;references:Id;comment:目标;"` // 目标

	OrderId string     `json:"order_id" gorm:"index;type:varchar(255);comment:订单ID;"`              // 订单ID
	Order   OrderSales `json:"order,omitzero" gorm:"foreignKey:OrderId;references:Id;comment:订单;"` // 订单

	StaffId string `json:"staff_id" gorm:"index;type:varchar(255);comment:员工ID;"`              // 员工ID
	Staff   Staff  `json:"staff,omitzero" gorm:"foreignKey:StaffId;references:Id;comment:员工;"` // 员工

	Achieve decimal.Decimal `json:"achieved" gorm:"type:decimal(10,2);comment:达成量;"` // 达成量
}

func TargetAddAchieve(tx *gorm.DB, order_id, store_id, staff_id string, amount decimal.Decimal, quantity int64) error {
	var (
		personals []TargetPersonal
	)
	tdb := tx.Model(&Target{})
	tdb = tdb.Where("store_id =?", store_id)
	tdb = tdb.Where("start_time <= ? AND end_time >= ?", time.Now(), time.Now())

	tpdb := tx.Model(&TargetPersonal{})
	tpdb = tpdb.Where("staff_id =?", staff_id)
	tpdb = tpdb.Where("target_id IN (?)", tdb.Select("id"))
	tpdb = tpdb.Preload("Target")
	if err := tpdb.Find(&personals).Error; err != nil {
		return err
	}

	for _, personal := range personals {
		achieve := TargetAchieve{
			TargetId: personal.TargetId,
			OrderId:  order_id,
			StaffId:  personal.StaffId,
		}

		switch personal.Target.Method {
		case enums.TargetMethodAmount:
			{
				achieve.Achieve = amount
			}
		case enums.TargetMethodQuantity:
			{
				achieve.Achieve = decimal.NewFromInt(quantity)
			}
		}

		if achieve.Achieve.IsZero() {
			continue
		}

		if err := tx.Model(&TargetPersonal{}).Where("id = ?", personal.Id).UpdateColumn("achieve", gorm.Expr("achieve + ?", achieve.Achieve)).Error; err != nil {
			return err
		}
		if err := tx.Create(&achieve).Error; err != nil {
			return err
		}
	}

	return nil
}

func init() {
	// 注册模型
	RegisterModels(
		&Target{},
		&TargetGroup{},
		&TargetPersonal{},
		&TargetAchieve{},
	)
	// 重置表
	RegisterRefreshModels(
	// &Target{},
	// &TargetGroup{},
	// &TargetPersonal{},
	// &TargetAchieve{},
	)
}
