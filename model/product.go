package model

import "jdy/types"

// 产品
type Product struct {
	SoftDelete

	Code   string   `json:"code" gorm:"uniqueIndex;type:varchar(255);<-:create;not NULL;comment:条码;" where:"label:条码;type:text;required;"` // 条码
	Name   string   `json:"name" gorm:"type:varchar(255);not NULL;comment:名称;"    where:"label:名称;type:text;required;"`                    // 名称
	Images []string `json:"images" gorm:"type:text;serializer:json;comment:图片;"`                                                           // 图片

	AccessFee float64 `json:"access_fee" gorm:"type:decimal(10,2);not NULL;comment:入网费;" where:"label:入网费;type:text;required;"` // 入网费
	Price     float64 `json:"price" gorm:"type:decimal(10,2);not NULL;comment:价格;" where:"label:价格;type:text;required;"`        // 价格
	LaborFee  float64 `json:"labor_fee" gorm:"type:decimal(10,2);not NULL;comment:工费;" where:"label:工费;type:text;required;"`    // 工费

	Weight      float64                 `json:"weight" gorm:"type:decimal(10,2);comment:总重量;" where:"label:总重量;type:number;"`                                    // 总重量
	WeightMetal float64                 `json:"weight_metal" gorm:"type:decimal(10,2);comment:金重;" where:"label:金重;type:number;"`                                // 金重
	WeightGem   float64                 `json:"weight_gem" gorm:"type:decimal(10,2);comment:主石重;" where:"label:主石重;type:number;"`                                // 主石重
	WeightOther float64                 `json:"weight_other" gorm:"type:decimal(10,2);comment:杂料重;" where:"label:主石重;type:number;"`                              // 杂料重
	NumGem      int                     `json:"num_gem" gorm:"type:tinyint(2);comment:主石数;" where:"label:主石重;type:number;"`                                      // 主石数
	NumOther    int                     `json:"num_other" gorm:"type:tinyint(2);comment:杂料数;" where:"label:主石重;type:number;"`                                    // 杂料数
	ColorMetal  types.ProductColor      `json:"color" gorm:"type:tinyint(2);comment:金颜色;" where:"label:金颜色;type:select;preset:{{.ColorMetal}}"`                  // 金颜色
	ColorGem    types.ProductColor      `json:"color_gem" gorm:"type:tinyint(2);comment:主石色;" where:"label:主石色;type:select;preset:{{.ColorGem}}"`                // 主石色
	Clarity     types.ProductClarity    `json:"clarity" gorm:"type:tinyint(2);comment:主石净度;" where:"label:净度;type:select;preset:{{.Clarity}}"`                   // 净度
	RetailType  types.ProductRetailType `json:"retail_type" gorm:"type:tinyint(2);not NULL;comment:零售方式;" where:"label:零售方式;type:select;preset:{{.RetailType}}"` // 零售方式
	Class       types.ProductClass      `json:"class" gorm:"type:tinyint(2);not NULL;comment:大类;" where:"label:大类;type:select;preset:{{.Class}}"`                // 大类
	Supplier    types.ProductSupplier   `json:"supplier" gorm:"type:tinyint(2);not NULL;comment:供应商;" where:"label:供应商;type:select;preset:{{.Supplier}}"`        // 供应商
	Material    types.ProductMaterial   `json:"material" gorm:"type:tinyint(2);not NULL;comment:材质;" where:"label:材质;type:select;preset:{{.Material}}"`          // 材质
	Quality     types.ProductQuality    `json:"quality" gorm:"type:tinyint(2);not NULL;comment:成色;" where:"label:成色;type:select;preset:{{.Quality}}"`            // 成色
	Gem         types.ProductGem        `json:"gem" gorm:"type:tinyint(2);not NULL;comment:宝石;" where:"label:宝石;type:select;preset:{{.Gem}}"`                    // 宝石
	Category    types.ProductCategory   `json:"category" gorm:"type:tinyint(2);not NULL;comment:品类;" where:"label:品类;type:select;preset:{{.Category}}"`          // 品类
	Brand       types.ProductBrand      `json:"brand" gorm:"type:tinyint(2);comment:品牌;" where:"label:品牌;type:select;preset:{{.Brand}}"`                         // 品牌
	Craft       types.ProductCraft      `json:"craft" gorm:"type:tinyint(2);comment:工艺;" where:"label:工艺;type:select;preset:{{.Craft}}"`                         // 工艺
	Style       string                  `json:"style" gorm:"type:varchar(255);comment:款式;" where:"label:款式;type:text;"`                                          // 款式
	Size        string                  `json:"size" gorm:"type:varchar(255);comment:手寸;" where:"label:手寸;type:text;"`                                           // 手寸

	IsSpecialOffer bool                `json:"is_special_offer" gorm:"comment:是否特价;" where:"label:是否特价;type:select;preset:{{.IsSpecialOffer}}"` // 是否特价
	Remark         string              `json:"remark" gorm:"type:text;comment:备注;" where:"label:备注;type:text;"`                                 // 备注
	Certificate    []string            `json:"certificate" gorm:"type:text;serializer:json;comment:证书;" where:"label:证书;type:text;"`            // 证书
	Status         types.ProductStatus `json:"status" gorm:"type:tinyint(2);comment:状态;" where:"label:状态;type:select;preset:{{.Status}}"`       // 状态

	ProductEnterId string        `json:"product_enter_id" gorm:"type:varchar(255);not NULL;comment:产品入库单ID;" where:"label:入库单;type:text;"` // 产品入库单ID
	ProductEnter   *ProductEnter `json:"product_enter" gorm:"foreignKey:ProductEnterId;references:Id;comment:产品入库单;"`                      // 产品入库单
}

// 产品入库单
type ProductEnter struct {
	SoftDelete

	Products []Product `json:"products" gorm:"foreignKey:ProductEnterId;references:Id;comment:产品;"`

	OperatorId string `json:"operator_id" gorm:"type:varchar(255);not NULL;comment:操作人ID;"`
	Operator   *Staff `json:"operator" gorm:"foreignKey:OperatorId;references:Id;comment:操作人;"`
}

func init() {
	// 注册模型
	RegisterModels(
		&Product{},
		&ProductEnter{},
	)
	// 重置表
	RegisterRefreshModels(
	// &Product{},
	// &ProductEnter{},
	)
}
