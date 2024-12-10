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
	ColorMetal  ProductColor      `json:"color" binding:"-"`              // 金颜色
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
