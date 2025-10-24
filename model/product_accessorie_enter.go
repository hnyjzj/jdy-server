package model

import (
	"fmt"
	"jdy/enums"
	"jdy/types"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// 配件入库单
type ProductAccessorieEnter struct {
	SoftDelete

	StoreId string `json:"store_id" gorm:"index;type:varchar(255);not NULL;comment:门店ID;"` // 门店ID
	Store   *Store `json:"store" gorm:"foreignKey:StoreId;references:Id;comment:门店;"`      // 门店

	Status enums.ProductEnterStatus `json:"status" gorm:"index;type:int(11);not NULL;comment:状态;"` // 状态
	Remark string                   `json:"remark" gorm:"type:text;comment:备注;"`                   // 备注

	Products     []ProductAccessorieEnterProduct `json:"products" gorm:"foreignKey:EnterId;references:Id;comment:产品;"` // 产品
	ProductCount int64                           `json:"product_count" gorm:"type:int(11);not NULL;comment:入库种类;"`     // 入库种类
	ProductTotal int64                           `json:"product_total" gorm:"type:int(11);not NULL;comment:入库总件数;"`    // 入库总件数

	OperatorId string `json:"operator_id" gorm:"type:varchar(255);not NULL;comment:操作人ID;"`     // 操作人ID
	Operator   *Staff `json:"operator" gorm:"foreignKey:OperatorId;references:Id;comment:操作人;"` // 操作人
	IP         string `json:"ip" gorm:"type:varchar(255);not NULL;comment:IP;"`                 // IP
}

func (ProductAccessorieEnter) WhereCondition(db *gorm.DB, req *types.ProductAccessorieEnterWhere) *gorm.DB {
	if req.Id != "" {
		db = db.Where("id LIKE ?", fmt.Sprintf("%%%s%%", req.Id))
	}
	if req.StoreId != "" {
		db = db.Where("store_id = ?", req.StoreId)
	}
	if req.Status != 0 {
		db = db.Where("status = ?", req.Status)
	}
	if req.Remark != "" {
		db = db.Where("remark LIKE ?", fmt.Sprintf("%%%s%%", req.Remark))
	}
	if req.Name != "" {
		db = db.Where("id IN (SELECT enter_id FROM product_accessorie_enter_products WHERE name LIKE ?)", fmt.Sprintf("%%%s%%", req.Name))
	}
	if req.StartTime != nil {
		db = db.Where("created_at >= ?", req.StartTime)
	}
	if req.EndTime != nil {
		db = db.Where("created_at <= ?", req.EndTime)
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
	SoftDelete

	Name       string                            `json:"name" gorm:"type:varchar(255);uniqueIndex:idx_enter_name;comment:名称;"` // 名称
	Type       enums.ProductAccessorieType       `json:"type" gorm:"type:int(11);comment:配件类型;"`                               // 配件类型
	RetailType enums.ProductAccessorieRetailType `json:"retail_type" gorm:"type:int(11);not NULL;comment:零售方式;"`               // 零售方式
	Price      decimal.Decimal                   `json:"price" gorm:"type:decimal(10,2);comment:单价;"`                          // 单价
	Remark     string                            `json:"remark" gorm:"type:text;comment:备注;"`                                  // 备注
	Stock      int64                             `json:"stock" gorm:"type:int(9);default:1;comment:库存;"`                       // 库存

	Status enums.ProductAccessorieStatus `json:"status" gorm:"type:int(11);default:1;comment:状态;"` // 状态

	EnterId string                  `json:"enter_id,omitempty" gorm:"type:varchar(255);uniqueIndex:idx_enter_name;comment:产品入库单ID;"` // 产品入库单ID
	Enter   *ProductAccessorieEnter `json:"enter,omitempty" gorm:"foreignKey:EnterId;references:Id;comment:产品入库单;"`                  // 产品入库单
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
