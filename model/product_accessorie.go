package model

import (
	"fmt"
	"jdy/enums"
	"jdy/types"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// 配件
type ProductAccessorie struct {
	SoftDelete

	StoreId string `json:"store_id" gorm:"type:varchar(255);uniqueIndex:idx_store_name;comment:门店ID;"` // 门店ID
	Store   Store  `json:"store,omitempty" gorm:"foreignKey:StoreId;references:Id;comment:门店;"`        // 门店

	Name       string                            `json:"name" gorm:"type:varchar(255);uniqueIndex:idx_store_name;comment:名称;"` // 名称
	Type       enums.ProductAccessorieType       `json:"type" gorm:"type:int(11);comment:配件类型;"`                               // 配件类型
	RetailType enums.ProductAccessorieRetailType `json:"retail_type" gorm:"type:int(11);not NULL;comment:零售方式;"`               // 零售方式
	Price      decimal.Decimal                   `json:"price" gorm:"type:decimal(10,2);comment:单价;"`                          // 单价
	Remark     string                            `json:"remark" gorm:"type:text;comment:备注;"`                                  // 备注
	Stock      int64                             `json:"stock" gorm:"index;type:int(9);default:1;comment:库存;"`                 // 库存

	Status enums.ProductAccessorieStatus `json:"status" gorm:"type:int(11);default:1;comment:状态;"` // 状态
}

func (ProductAccessorie) WhereCondition(db *gorm.DB, req *types.ProductAccessorieWhere) *gorm.DB {
	if req.Id != "" {
		db = db.Where("id LIKE ?", fmt.Sprintf("%%%s%%", req.Id))
	}
	if req.StoreId != "" {
		db = db.Where("store_id = ?", req.StoreId)
	}
	if req.Name != "" {
		db = db.Where("name LIKE (?)", fmt.Sprintf("%%%s%%", req.Name))
	}
	if req.Type != 0 {
		db = db.Where("type = ?", req.Type)
	}
	if req.RetailType != 0 {
		db = db.Where("retail_type = ?", req.RetailType)
	}
	if req.Remark != "" {
		db = db.Where("remark LIKE (?)", fmt.Sprintf("%%%s%%", req.Remark))
	}
	if req.Stock != 0 {
		db = db.Where("stock >= ?", req.Stock)
	}
	if req.Status != 0 {
		db = db.Where("status = ?", req.Status)
	}

	return db
}

func (ProductAccessorie) Preloads(db *gorm.DB) *gorm.DB {
	db = db.Preload("Store")

	return db
}
func init() {
	// 注册模型
	RegisterModels(
		&ProductAccessorie{},
	)
	// 重置表
	RegisterRefreshModels(
	// &ProductAccessorie{},
	)
}
