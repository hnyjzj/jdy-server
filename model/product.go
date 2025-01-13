package model

import (
	"fmt"
	"jdy/enums"
	"jdy/types"

	"gorm.io/gorm"
)

// 产品
type Product struct {
	SoftDelete

	Code   string   `json:"code" gorm:"uniqueIndex;type:varchar(255);<-:create;not NULL;comment:条码;"` // 条码
	Name   string   `json:"name" gorm:"type:varchar(255);not NULL;comment:名称;"`                       // 名称
	Images []string `json:"images" gorm:"type:text;serializer:json;comment:图片;"`                      // 图片

	AccessFee float64 `json:"access_fee" gorm:"type:decimal(10,2);not NULL;comment:入网费;"` // 入网费
	Price     float64 `json:"price" gorm:"type:decimal(10,2);not NULL;comment:价格;"`       // 价格
	LaborFee  float64 `json:"labor_fee" gorm:"type:decimal(10,2);not NULL;comment:工费;"`   // 工费

	Weight      float64                 `json:"weight" gorm:"type:decimal(10,2);comment:总重量;"`             // 总重量
	WeightMetal float64                 `json:"weight_metal" gorm:"type:decimal(10,2);comment:金重;"`        // 金重
	WeightGem   float64                 `json:"weight_gem" gorm:"type:decimal(10,2);comment:主石重;"`         // 主石重
	WeightOther float64                 `json:"weight_other" gorm:"type:decimal(10,2);comment:杂料重;"`       // 杂料重
	NumGem      int                     `json:"num_gem" gorm:"type:tinyint(2);comment:主石数;"`               // 主石数
	NumOther    int                     `json:"num_other" gorm:"type:tinyint(2);comment:杂料数;"`             // 杂料数
	ColorMetal  enums.ProductColor      `json:"color_metal" gorm:"type:tinyint(2);comment:金颜色;"`           // 金颜色
	ColorGem    enums.ProductColor      `json:"color_gem" gorm:"type:tinyint(2);comment:主石色;"`             // 主石色
	Clarity     enums.ProductClarity    `json:"clarity" gorm:"type:tinyint(2);comment:主石净度;"`              // 净度
	RetailType  enums.ProductRetailType `json:"retail_type" gorm:"type:tinyint(2);not NULL;comment:零售方式;"` // 零售方式
	Class       enums.ProductClass      `json:"class" gorm:"type:tinyint(2);not NULL;comment:大类;"`         // 大类
	Supplier    enums.ProductSupplier   `json:"supplier" gorm:"type:tinyint(2);not NULL;comment:供应商;"`     // 供应商
	Material    enums.ProductMaterial   `json:"material" gorm:"type:tinyint(2);not NULL;comment:材质;"`      // 材质
	Quality     enums.ProductQuality    `json:"quality" gorm:"type:tinyint(2);not NULL;comment:成色;"`       // 成色
	Gem         enums.ProductGem        `json:"gem" gorm:"type:tinyint(2);not NULL;comment:宝石;"`           // 宝石
	Category    enums.ProductCategory   `json:"category" gorm:"type:tinyint(2);not NULL;comment:品类;"`      // 品类
	Brand       enums.ProductBrand      `json:"brand" gorm:"type:tinyint(2);comment:品牌;"`                  // 品牌
	Craft       enums.ProductCraft      `json:"craft" gorm:"type:tinyint(2);comment:工艺;"`                  // 工艺
	Style       string                  `json:"style" gorm:"type:varchar(255);comment:款式;"`                // 款式
	Size        string                  `json:"size" gorm:"type:varchar(255);comment:手寸;"`                 // 手寸

	IsSpecialOffer bool                `json:"is_special_offer" gorm:"comment:是否特价;"`                    // 是否特价
	Remark         string              `json:"remark" gorm:"type:text;comment:备注;"`                      // 备注
	Certificate    []string            `json:"certificate" gorm:"type:text;serializer:json;comment:证书;"` // 证书
	Status         enums.ProductStatus `json:"status" gorm:"type:tinyint(2);comment:状态;"`                // 状态
	Type           enums.ProductType   `json:"type" gorm:"type:tinyint(2);comment:类型;"`                  // 类型

	ProductEnterId string        `json:"product_enter_id" gorm:"type:varchar(255);not NULL;comment:产品入库单ID;"`         // 产品入库单ID
	ProductEnter   *ProductEnter `json:"product_enter" gorm:"foreignKey:ProductEnterId;references:Id;comment:产品入库单;"` // 产品入库单

	ProductDamages []ProductDamage `json:"product_damage" gorm:"foreignKey:ProductId;references:Id;comment:报损记录;"` // 报损记录

	StoreId string `json:"store_id" gorm:"type:varchar(255);comment:店铺ID;"`           // 店铺ID
	Store   *Store `json:"store" gorm:"foreignKey:StoreId;references:Id;comment:店铺;"` // 店铺
}

func (Product) WhereCondition(db *gorm.DB, query *types.ProductWhere) *gorm.DB {
	if query.Code != "" {
		db = db.Where("code = ?", fmt.Sprint(query.Code))
	}
	if query.Name != "" {
		db = db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", query.Name))
	}
	if query.AccessFee != 0 {
		db = db.Where("access_fee = ?", float64(query.AccessFee))
	}
	if query.Price != 0 {
		db = db.Where("price = ?", float64(query.Price))
	}
	if query.LaborFee != 0 {
		db = db.Where("labor_fee = ?", float64(query.LaborFee))
	}
	if query.Weight != 0 {
		db = db.Where("weight = ?", float64(query.Weight))
	}
	if query.WeightMetal != 0 {
		db = db.Where("weight_metal = ?", float64(query.WeightMetal))
	}
	if query.WeightGem != 0 {
		db = db.Where("weight_gem = ?", float64(query.WeightGem))
	}
	if query.WeightOther != 0 {
		db = db.Where("weight_other = ?", float64(query.WeightOther))
	}
	if query.NumGem != 0 {
		db = db.Where("num_gem = ?", int(query.NumGem))
	}
	if query.NumOther != 0 {
		db = db.Where("num_other = ?", int(query.NumOther))
	}
	if query.ColorMetal != 0 {
		db = db.Where("color_metal = ?", query.ColorMetal)
	}
	if query.ColorGem != 0 {
		db = db.Where("color_gem = ?", query.ColorGem)
	}
	if query.Clarity != 0 {
		db = db.Where("clarity = ?", query.Clarity)
	}
	if query.RetailType != 0 {
		db = db.Where("retail_type = ?", query.RetailType)
	}
	if query.Class != 0 {
		db = db.Where("class = ?", query.Class)
	}
	if query.Supplier != 0 {
		db = db.Where("supplier = ?", query.Supplier)
	}
	if query.Material != 0 {
		db = db.Where("material = ?", query.Material)
	}
	if query.Quality != 0 {
		db = db.Where("quality = ?", query.Quality)
	}
	if query.Gem != 0 {
		db = db.Where("gem = ?", query.Gem)
	}
	if query.Category != 0 {
		db = db.Where("category = ?", query.Category)
	}
	if query.Brand != 0 {
		db = db.Where("brand = ?", query.Brand)
	}
	if query.Craft != 0 {
		db = db.Where("craft = ?", query.Craft)
	}
	if query.Style != "" {
		db = db.Where("style LIKE ?", fmt.Sprintf("%%%s%%", query.Name))
	}
	if query.Size != "" {
		db = db.Where("size LIKE ?", fmt.Sprintf("%%%s%%", query.Name))
	}
	if query.IsSpecialOffer {
		db = db.Where("is_special_offer = ?", query.IsSpecialOffer)
	}
	if query.Certificate != nil {
		db = db.Where("certificate IN ?", query.Certificate)
	}
	if query.Status != 0 {
		db = db.Where("status = ?", query.Status)
	}
	if query.Type != 0 {
		db = db.Where("type = ?", query.Type)
	}
	if query.ProductEnterId != "" {
		db = db.Where("product_enter_id = ?", query.ProductEnterId)
	}

	return db
}

// 产品入库单
type ProductEnter struct {
	SoftDelete

	Products []Product `json:"products" gorm:"foreignKey:ProductEnterId;references:Id;comment:产品;"` // 产品

	OperatorId string `json:"operator_id" gorm:"type:varchar(255);not NULL;comment:操作人ID;"`     // 操作人ID
	Operator   *Staff `json:"operator" gorm:"foreignKey:OperatorId;references:Id;comment:操作人;"` // 操作人

	IP string `json:"ip" gorm:"type:varchar(255);not NULL;comment:IP;"` // IP
}

func (ProductEnter) WhereCondition(db *gorm.DB, query *types.ProductEnterWhere) *gorm.DB {
	if query.Id != "" {
		db = db.Where("id = ?", query.Id)
	}
	if query.StartTime != nil && query.EndTime != nil {
		db = db.Where("created_at BETWEEN ? AND ?", query.StartTime, query.EndTime)
	}
	return db
}

type ProductDamage struct {
	SoftDelete

	ProductId string   `json:"product_id" gorm:"type:varchar(255);not NULL;comment:产品ID;"`
	Product   *Product `json:"product" gorm:"foreignKey:ProductId;references:Id;comment:产品;"`

	OperatorId string `json:"-" gorm:"type:varchar(255);not NULL;comment:操作人ID;"`        // 操作人ID
	Operator   *Staff `json:"-" gorm:"foreignKey:OperatorId;references:Id;comment:操作人;"` // 操作人

	Reason string `json:"reason" gorm:"type:text;not NULL;comment:原因;"`    // 原因
	IP     string `json:"-" gorm:"type:varchar(255);not NULL;comment:IP;"` // IP
}

type ProductAllocate struct {
	SoftDelete

	Method enums.ProductAllocateMethod `json:"method" gorm:"type:tinyint(2);not NULL;comment:调拨方式;"` // 调拨方式
	Type   enums.ProductType           `json:"type" gorm:"type:tinyint(2);not NULL;comment:产品类型;"`   // 仓库类型
	Reason enums.ProductAllocateReason `json:"reason" gorm:"type:tinyint(2);not NULL;comment:调拨原因;"` // 调拨原因
	Remark string                      `json:"remark" gorm:"type:text;comment:备注;"`                  // 备注

	Status enums.ProductAllocateStatus `json:"status" gorm:"type:tinyint(2);comment:状态;"` // 状态

	StoreId string `json:"store_id" gorm:"type:tinyint(2);comment:调拨门店;"` // 调拨门店
	Store   *Store `json:"store" gorm:"foreignKey:StoreId;references:Id;comment:调拨门店;"`

	Products []Product `json:"product" gorm:"many2many:product_allocate_list;"` // 产品

	OperatorId string `json:"-" gorm:"type:varchar(255);not NULL;comment:操作人ID;"`        // 操作人ID
	Operator   *Staff `json:"-" gorm:"foreignKey:OperatorId;references:Id;comment:操作人;"` // 操作人
	IP         string `json:"-" gorm:"type:varchar(255);not NULL;comment:IP;"`           // IP
}

func (ProductAllocate) WhereCondition(db *gorm.DB, query *types.ProductAllocateWhere) *gorm.DB {
	if query.Method != 0 {
		db = db.Where("method = ?", query.Method)
	}
	if query.Type != 0 {
		db = db.Where("type = ?", query.Type)
	}
	if query.Reason != 0 {
		db = db.Where("reason = ?", query.Reason)
	}
	if query.StoreId != "" {
		db = db.Where("store_id = ?", query.StoreId)
	}
	if query.StartTime != nil && query.EndTime != nil {
		db = db.Where("created_at BETWEEN ? AND ?", query.StartTime, query.EndTime)
	}
	return db
}

func init() {
	// 注册模型
	RegisterModels(
		&Product{},
		&ProductEnter{},
		&ProductDamage{},
		&ProductAllocate{},
	)
	// 重置表
	RegisterRefreshModels(
	// &Product{},
	// &ProductEnter{},
	// &ProductDamage{},
	// &ProductAllocate{},
	)
}
