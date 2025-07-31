package model

import (
	"jdy/enums"
	"jdy/types"

	"gorm.io/gorm"
)

// 配件入库单
type ProductAccessorieEnter struct {
	SoftDelete

	StoreId string `json:"store_id" gorm:"type:varchar(255);not NULL;comment:门店ID;"`  // 门店ID
	Store   *Store `json:"store" gorm:"foreignKey:StoreId;references:Id;comment:门店;"` // 门店

	Status enums.ProductEnterStatus `json:"status" gorm:"type:int(11);not NULL;comment:状态;"` // 状态
	Remark string                   `json:"remark" gorm:"type:text;comment:备注;"`             // 备注

	Products     []ProductAccessorieEnterProduct `json:"products" gorm:"foreignKey:EnterId;references:Id;comment:产品;"` // 产品
	ProductCount int64                           `json:"product_count" gorm:"type:int(11);not NULL;comment:入库种类;"`     // 入库种类
	ProductTotal int64                           `json:"product_total" gorm:"type:int(11);not NULL;comment:入库总件数;"`    // 入库总件数

	OperatorId string `json:"operator_id" gorm:"type:varchar(255);not NULL;comment:操作人ID;"`     // 操作人ID
	Operator   *Staff `json:"operator" gorm:"foreignKey:OperatorId;references:Id;comment:操作人;"` // 操作人
	IP         string `json:"ip" gorm:"type:varchar(255);not NULL;comment:IP;"`                 // IP
}

func (ProductAccessorieEnter) WhereCondition(db *gorm.DB, query *types.ProductAccessorieEnterWhere) *gorm.DB {
	if query.Id != "" {
		db = db.Where("id = ?", query.Id)
	}
	if query.StoreId != "" {
		db = db.Where("store_id = ?", query.StoreId)
	}
	if query.Status != 0 {
		db = db.Where("status = ?", query.Status)
	} else {
		db = db.Where("status IN (?)", []enums.ProductEnterStatus{
			enums.ProductEnterStatusDraft,
		})
	}
	if query.Remark != "" {
		db = db.Where("remark LIKE ?", "%"+query.Remark+"%")
	}
	if query.Name != "" {
		db = db.Where("id IN (SELECT enter_id FROM product_accessorie_enter_products WHERE name LIKE ?)", "%"+query.Name+"%")
	}
	if query.StartTime != nil {
		db = db.Where("created_at >= ?", query.StartTime)
	}
	if query.EndTime != nil {
		db = db.Where("created_at <= ?", query.EndTime)
	}

	return db
}

func (ProductAccessorieEnter) Preloads(db *gorm.DB, pages *types.PageReq) *gorm.DB {
	db = db.Preload("Store")
	db = db.Preload("Products", func(tx *gorm.DB) *gorm.DB {
		pdb := tx
		if pages != nil {
			pdb = PageCondition(pdb, &types.PageReq{Page: pages.Page, Limit: pages.Limit})
		}

		return pdb
	})
	db = db.Preload("Operator")

	return db
}

type ProductAccessorieEnterProduct struct {
	ProductAccessorie

	EnterId string                  `json:"enter_id,omitempty" gorm:"type:varchar(255);comment:产品入库单ID;"`           // 产品入库单ID
	Enter   *ProductAccessorieEnter `json:"enter,omitempty" gorm:"foreignKey:EnterId;references:Id;comment:产品入库单;"` // 产品入库单
}

func init() {
	// 注册模型
	RegisterModels(
		&ProductAccessorieEnter{},
		&ProductAccessorieEnterProduct{},
	)
	// 重置表
	RegisterRefreshModels(
	// &ProductAccessorieEnter{},
	// &ProductAccessorieEnterProduct{},
	)
}
