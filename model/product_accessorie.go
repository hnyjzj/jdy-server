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

	StoreId string `json:"store_id" gorm:"type:varchar(255);comment:门店ID;"`                     // 门店ID
	Store   Store  `json:"store,omitempty" gorm:"foreignKey:StoreId;references:Id;comment:门店;"` // 门店

	Code     string                     `json:"code" gorm:"type:varchar(255);not NULL;comment:配件条目ID;"`                // 配件条目ID
	Category *ProductAccessorieCategory `json:"category,omitempty" gorm:"foreignKey:Id;references:Code;comment:配件条目;"` // 配件条目

	Stock     int64           `json:"stock" gorm:"type:int(9);default:1;comment:库存;"`             // 库存
	AccessFee decimal.Decimal `json:"access_fee" gorm:"type:decimal(10,2);not NULL;comment:入网费;"` // 入网费

	Status enums.ProductStatus `json:"status" gorm:"type:tinyint(2);not NULL;comment:状态;"` // 状态

	EnterId string                  `json:"enter_id,omitempty" gorm:"type:varchar(255);comment:产品入库单ID;"`           // 产品入库单ID
	Enter   *ProductAccessorieEnter `json:"enter,omitempty" gorm:"foreignKey:EnterId;references:Id;comment:产品入库单;"` // 产品入库单
}

func (ProductAccessorie) WhereCondition(db *gorm.DB, query *types.ProductAccessorieWhere) *gorm.DB {
	if !query.AccessFee.IsZero() {
		db = db.Where("access_fee = ?", query.AccessFee)
	}

	if query.StoreId != "" {
		db = db.Where("store_id = ?", query.StoreId)
	}

	return db
}

// 配件条目
type ProductAccessorieCategory struct {
	SoftDelete

	TypePart enums.ProductTypePart `json:"type_part,omitempty" gorm:"type:tinyint(2);comment:配件类型;"` // 配件类型

	Name          string                            `json:"name" gorm:"type:varchar(255);uniqueIndex;comment:名称;"`       // 名称
	Code          string                            `json:"code" gorm:"type:varchar(255);comment:条码;"`                   // 条码
	RetailType    enums.ProductAccessorieRetailType `json:"retail_type" gorm:"type:tinyint(2);not NULL;comment:零售方式;"`   // 零售方式
	Weight        decimal.Decimal                   `json:"weight" gorm:"type:decimal(10,2);comment:重量;"`                // 重量
	AccessFee     decimal.Decimal                   `json:"access_fee" gorm:"type:decimal(10,2);not NULL;comment:入网费;"`  // 入网费
	LabelPrice    decimal.Decimal                   `json:"label_price" gorm:"type:decimal(10,2);not NULL;comment:标签价;"` // 标签价
	Material      enums.ProductAccessorieMaterial   `json:"material" gorm:"type:tinyint(2);not NULL;comment:材质;"`        // 材质
	Quality       enums.ProductQuality              `json:"quality" gorm:"type:tinyint(2);not NULL;comment:成色;"`         // 成色
	Gem           enums.ProductGem                  `json:"gem" gorm:"type:tinyint(2);not NULL;comment:主石;"`             // 主石
	Category      enums.ProductCategory             `json:"category" gorm:"type:tinyint(2);not NULL;comment:品类;"`        // 品类
	Specification string                            `json:"specification" gorm:"type:varchar(255);comment:规格;"`          // 规格
	Color         string                            `json:"color" gorm:"type:varchar(255);comment:颜色;"`                  // 颜色
	Series        string                            `json:"series" gorm:"type:varchar(255);comment:系列;"`                 // 系列
	Supplier      string                            `json:"supplier" gorm:"type:varchar(255);comment:供应商;"`              // 供应商
	Remark        string                            `json:"remark" gorm:"type:text;comment:备注;"`                         // 备注

	Product ProductAccessorie `json:"product,omitempty" gorm:"foreignKey:Id;references:Code;comment:产品;"` // 产品

	OperatorId string `json:"operator_id" gorm:"type:varchar(255);not NULL;comment:操作员ID;"`     // 操作员ID
	Operator   Staff  `json:"operator" gorm:"foreignKey:OperatorId;references:Id;comment:操作员;"` // 操作员
	IP         string `json:"ip" gorm:"type:varchar(255);not NULL;comment:IP地址;"`               // IP地址
}

func (ProductAccessorieCategory) WhereCondition(db *gorm.DB, query *types.ProductAccessorieCategoryWhere) *gorm.DB {
	if query.Id != "" {
		db = db.Where("id = ?", query.Id)
	}

	if query.TypePart != 0 {
		db = db.Where("type_part = ?", query.TypePart)
	}

	if query.Name != "" {
		db = db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", query.Name))
	}

	if query.Code != "" {
		db = db.Where("code = ?", query.Code)
	}

	if query.RetailType != 0 {
		db = db.Where("retail_type = ?", query.RetailType)
	}

	if query.Material != 0 {
		db = db.Where("material = ?", query.Material)
	}

	if query.Gem != 0 {
		db = db.Where("gem = ?", query.Gem)
	}

	if query.Category != 0 {
		db = db.Where("category = ?", query.Category)
	}

	if query.Specification != "" {
		db = db.Where("specification = ?", query.Specification)
	}

	if query.Color != "" {
		db = db.Where("color = ?", query.Color)
	}

	if query.Series != "" {
		db = db.Where("series = ?", query.Series)
	}

	if query.Supplier != "" {
		db = db.Where("supplier = ?", query.Supplier)
	}

	if query.Remark != "" {
		db = db.Where("remark LIKE ?", fmt.Sprintf("%%%s%%", query.Remark))
	}

	if query.StoreId != "" {
		db = db.Preload("Product", func(*gorm.DB) *gorm.DB {
			return db.Where("store_id = ?", query.StoreId)
		})
	}

	return db
}

// 配件入库单
type ProductAccessorieEnter struct {
	SoftDelete

	StoreId string `json:"store_id" gorm:"type:varchar(255);not NULL;comment:门店ID;"`  // 门店ID
	Store   *Store `json:"store" gorm:"foreignKey:StoreId;references:Id;comment:门店;"` // 门店

	Remark string                   `json:"remark" gorm:"type:text;comment:备注;"`                // 备注
	Status enums.ProductEnterStatus `json:"status" gorm:"type:tinyint(2);not NULL;comment:状态;"` // 状态

	Products []ProductAccessorie `json:"products" gorm:"foreignKey:EnterId;references:Id;comment:产品;"` // 产品

	OperatorId string `json:"operator_id" gorm:"type:varchar(255);not NULL;comment:操作人ID;"`     // 操作人ID
	Operator   *Staff `json:"operator" gorm:"foreignKey:OperatorId;references:Id;comment:操作人;"` // 操作人
	IP         string `json:"ip" gorm:"type:varchar(255);not NULL;comment:IP;"`                 // IP
}

func (ProductAccessorieEnter) WhereCondition(db *gorm.DB, query *types.ProductAccessorieEnterWhere) *gorm.DB {
	if query.StoreId != "" {
		db = db.Where("store_id = ?", query.StoreId)
	}
	if query.Status != 0 {
		db = db.Where("status = ?", query.Status)
	}
	if query.Remark != "" {
		db = db.Where("remark = ?", query.Remark)
	}
	if query.Code != "" {
		db = db.Where("id IN (SELECT enter_id FROM product_accessories WHERE code = ?)", query.Code)
	}
	if query.StartAt != nil && query.EndAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", query.StartAt, query.EndAt)
	}
	return db
}

// 产品调拨单
type ProductAccessorieAllocate struct {
	SoftDelete

	Method enums.ProductAllocateMethod `json:"method" gorm:"type:tinyint(2);not NULL;comment:调拨类型;"` // 调拨类型
	Status enums.ProductAllocateStatus `json:"status" gorm:"type:tinyint(2);comment:状态;"`            // 状态
	Remark string                      `json:"remark" gorm:"type:text;comment:备注;"`                  // 备注

	FromStoreId string `json:"from_store_id" gorm:"type:varchar(255);comment:调出门店;"` // 调出门店
	FromStore   *Store `json:"from_store" gorm:"foreignKey:FromStoreId;references:Id;comment:调出门店;"`
	ToStoreId   string `json:"to_store_id" gorm:"type:varchar(255);comment:调入门店;"` // 调入门店
	ToStore     *Store `json:"to_store" gorm:"foreignKey:ToStoreId;references:Id;comment:调入门店;"`

	Products []ProductAccessorieAllocateProduct `json:"products" gorm:"foreignKey:AllocateId;references:Id;comment:产品;"`

	OperatorId string `json:"operator_id" gorm:"type:varchar(255);not NULL;comment:操作人ID;"`     // 操作人ID
	Operator   *Staff `json:"operator" gorm:"foreignKey:OperatorId;references:Id;comment:操作人;"` // 操作人
	IP         string `json:"-" gorm:"type:varchar(255);not NULL;comment:IP;"`                  // IP
}

func (ProductAccessorieAllocate) WhereCondition(db *gorm.DB, query *types.ProductAccessorieAllocateWhere) *gorm.DB {
	if query.Method != 0 {
		db = db.Where("method = ?", query.Method)
	}
	if query.Remark != "" {
		db = db.Where("remark LIKE ?", fmt.Sprintf("%%%s%%", query.Remark))
	}
	if query.FromStoreId != "" {
		db = db.Where("from_store_id = ?", query.FromStoreId)
	}
	if query.ToStoreId != "" {
		db = db.Where("to_store_id = ?", query.ToStoreId)
	}
	if query.StartTime != nil && query.EndTime != nil {
		db = db.Where("created_at BETWEEN ? AND ?", query.StartTime, query.EndTime)
	}
	return db
}

type ProductAccessorieAllocateProduct struct {
	ProductId string             `json:"product_id" gorm:"type:varchar(255);not NULL;comment:产品ID;"`    // 产品ID
	Product   *ProductAccessorie `json:"product" gorm:"foreignKey:ProductId;references:Id;comment:产品;"` // 产品

	AllocateId string                     `json:"allocate_id" gorm:"type:varchar(255);not NULL;comment:调拨单ID;"`     // 调拨单ID
	Allocate   *ProductAccessorieAllocate `json:"allocate" gorm:"foreignKey:AllocateId;references:Id;comment:调拨单;"` // 调拨单
	Quantity   int64                      `json:"quantity" gorm:"type:int(8);not NULL;comment:数量;"`                 // 数量
}

func init() {
	// 注册模型
	RegisterModels(
		&ProductAccessorie{},
		&ProductAccessorieCategory{},
		&ProductAccessorieEnter{},
		&ProductAccessorieAllocate{},
		&ProductAccessorieAllocateProduct{},
	)
	// 重置表
	RegisterRefreshModels(
	// &ProductAccessorie{},
	// &ProductAccessorieCategory{},
	// &ProductAccessorieEnter{},
	// &ProductAccessorieAllocate{},
	// &ProductAccessorieAllocateProduct{},
	)
}
