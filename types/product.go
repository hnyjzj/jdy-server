package types

type ProductEnterReq struct {
	Products []ProductEnterReqProduct `json:"products" binding:"required"`
}

type ProductEnterReqProduct struct {
	Code string `json:"code" binding:"required"` // 条码
	Name string `json:"name" binding:"required"` // 名称

	AccessFee float64 `json:"access_fee" binding:"required"` // 入网费
	Price     float64 `json:"price" binding:"required"`      // 标签价
	LaborFee  float64 `json:"labor_fee" binding:"required"`  // 工费

	Weight      float64           `json:"weight" binding:"-"`             // 总重量
	WeightMetal float64           `json:"weight_metal" binding:"-"`       // 金重
	WeightGem   float64           `json:"weight_gem" binding:"-"`         // 主石重
	WeightOther float64           `json:"weight_other" binding:"-"`       // 杂料重
	NumGem      int               `json:"num_gem" binding:"-"`            // 主石数
	NumOther    int               `json:"num_other" binding:"-"`          // 杂料数
	ColorMetal  ProductColor      `json:"color_metal" binding:"-"`        // 金颜色
	ColorGem    ProductColor      `json:"color_gem" binding:"-"`          // 主石色
	Clarity     ProductClarity    `json:"clarity" binding:"-"`            // 净度
	RetailType  ProductRetailType `json:"retail_type" binding:"required"` // 零售方式
	Class       ProductClass      `json:"class" binding:"required"`       // 大类
	Supplier    ProductSupplier   `json:"supplier" binding:"required"`    // 供应商
	Material    ProductMaterial   `json:"material" binding:"required"`    // 材质
	Quality     ProductQuality    `json:"quality" binding:"required"`     // 成色
	Gem         ProductGem        `json:"gem" binding:"required"`         // 宝石
	Category    ProductCategory   `json:"category" binding:"required"`    // 品类
	Brand       ProductBrand      `json:"brand" binding:"-"`              // 品牌
	Craft       ProductCraft      `json:"craft" binding:"-"`              // 工艺
	Style       string            `json:"style" binding:"-"`              // 款式
	Size        string            `json:"size" binding:"-"`               // 手寸

	IsSpecialOffer bool     `json:"is_special_offer" binding:"-"` // 是否特价
	Remark         string   `json:"remark" binding:"-"`           // 备注
	Certificate    []string `json:"certificate" binding:"-"`      // 证书
}

type ProductWhere struct {
	Code string `json:"code" where:"label:条码;type:text;required;"` // 条码
	Name string `json:"name" where:"label:名称;type:text;required;"` // 名称

	AccessFee float64 `json:"access_fee" where:"label:入网费;type:text;required;"` // 入网费
	Price     float64 `json:"price"  where:"label:价格;type:text;required;"`      // 价格
	LaborFee  float64 `json:"labor_fee" where:"label:工费;type:text;required;"`   // 工费

	Weight      float64           `json:"weight" where:"label:总重量;type:number;"`                             // 总重量
	WeightMetal float64           `json:"weight_metal" where:"label:金重;type:number;"`                        // 金重
	WeightGem   float64           `json:"weight_gem" where:"label:主石重;type:number;"`                         // 主石重
	WeightOther float64           `json:"weight_other" where:"label:杂料重;type:number;"`                       // 杂料重
	NumGem      int               `json:"num_gem" where:"label:主石数;type:number;"`                            // 主石数
	NumOther    int               `json:"num_other" where:"label:杂料数;type:number;"`                          // 杂料数
	ColorMetal  ProductColor      `json:"color_metal" where:"label:金颜色;type:select;preset:{{.ColorMetal}}"`  // 金颜色
	ColorGem    ProductColor      `json:"color_gem" where:"label:主石色;type:select;preset:{{.ColorGem}}"`      // 主石色
	Clarity     ProductClarity    `json:"clarity" where:"label:净度;type:select;preset:{{.Clarity}}"`          // 净度
	RetailType  ProductRetailType `json:"retail_type" where:"label:零售方式;type:select;preset:{{.RetailType}}"` // 零售方式
	Class       ProductClass      `json:"class" where:"label:大类;type:select;preset:{{.Class}}"`              // 大类
	Supplier    ProductSupplier   `json:"supplier" where:"label:供应商;type:select;preset:{{.Supplier}}"`       // 供应商
	Material    ProductMaterial   `json:"material" where:"label:材质;type:select;preset:{{.Material}}"`        // 材质
	Quality     ProductQuality    `json:"quality" where:"label:成色;type:select;preset:{{.Quality}}"`          // 成色
	Gem         ProductGem        `json:"gem" where:"label:宝石;type:select;preset:{{.Gem}}"`                  // 宝石
	Category    ProductCategory   `json:"category" where:"label:品类;type:select;preset:{{.Category}}"`        // 品类
	Brand       ProductBrand      `json:"brand" where:"label:品牌;type:select;preset:{{.Brand}}"`              // 品牌
	Craft       ProductCraft      `json:"craft" where:"label:工艺;type:select;preset:{{.Craft}}"`              // 工艺
	Style       string            `json:"style" where:"label:款式;type:text;"`                                 // 款式
	Size        string            `json:"size" where:"label:手寸;type:text;"`                                  // 手寸

	IsSpecialOffer bool          `json:"is_special_offer" where:"label:是否特价;type:select;preset:{{.IsSpecialOffer}}"` // 是否特价
	Remark         string        `json:"remark" where:"label:备注;type:text;"`                                         // 备注
	Certificate    []string      `json:"certificate" where:"label:证书;type:text;"`                                    // 证书
	Status         ProductStatus `json:"status" where:"label:状态;type:select;preset:{{.Status}}"`                     // 状态

	ProductEnterId string `json:"product_enter_id" where:"label:入库单;type:text;"` // 产品入库单ID
}

type ProductListReq struct {
	PageReq
	Where ProductWhere `json:"where" binding:"required"`
}

type ProductInfoReq struct {
	Code string `json:"code" binding:"required"` // 条码
}

type ProductUpdateReq struct {
	Id string `json:"id" binding:"required"` // 产品ID

	Name   string   `json:"name"`   // 名称
	Images []string `json:"images"` // 图片

	AccessFee float64 `json:"access_fee"` // 入网费
	Price     float64 `json:"price"`      // 标签价
	LaborFee  float64 `json:"labor_fee"`  // 工费

	Weight      float64           `json:"weight"`       // 总重量
	WeightMetal float64           `json:"weight_metal"` // 金重
	WeightGem   float64           `json:"weight_gem"`   // 主石重
	WeightOther float64           `json:"weight_other"` // 杂料重
	NumGem      int               `json:"num_gem"`      // 主石数
	NumOther    int               `json:"num_other"`    // 杂料数
	ColorMetal  ProductColor      `json:"color_metal"`  // 金颜色
	ColorGem    ProductColor      `json:"color_gem"`    // 主石色
	Clarity     ProductClarity    `json:"clarity"`      // 净度
	RetailType  ProductRetailType `json:"retail_type"`  // 零售方式
	Class       ProductClass      `json:"class"`        // 大类
	Supplier    ProductSupplier   `json:"supplier"`     // 供应商
	Material    ProductMaterial   `json:"material"`     // 材质
	Quality     ProductQuality    `json:"quality"`      // 成色
	Gem         ProductGem        `json:"gem"`          // 宝石
	Category    ProductCategory   `json:"category"`     // 品类
	Brand       ProductBrand      `json:"brand"`        // 品牌
	Craft       ProductCraft      `json:"craft"`        // 工艺
	Style       string            `json:"style"`        // 款式
	Size        string            `json:"size"`         // 手寸

	IsSpecialOffer bool     `json:"is_special_offer"` // 是否特价
	Remark         string   `json:"remark"`           // 备注
	Certificate    []string `json:"certificate"`      // 证书
}

type ProductDamageReq struct {
	Code   string `json:"code" binding:"required"`   // 条码
	Reason string `json:"reason" binding:"required"` // 损坏原因
}
