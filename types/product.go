package types

import (
	"errors"
	"jdy/enums"
	"time"
)

type ProductEnterReq struct {
	Products []ProductEnterReqProduct `json:"products" binding:"required"`
}

type ProductEnterReqProduct struct {
	Code string `json:"code" binding:"required"` // 条码
	Name string `json:"name" binding:"required"` // 名称

	AccessFee float64 `json:"access_fee" binding:"required"` // 入网费
	Price     float64 `json:"price" binding:"required"`      // 标签价
	LaborFee  float64 `json:"labor_fee" binding:"required"`  // 工费

	Weight      float64                 `json:"weight" binding:"-"`             // 总重量
	WeightMetal float64                 `json:"weight_metal" binding:"-"`       // 金重
	WeightGem   float64                 `json:"weight_gem" binding:"-"`         // 主石重
	WeightOther float64                 `json:"weight_other" binding:"-"`       // 杂料重
	NumGem      int                     `json:"num_gem" binding:"-"`            // 主石数
	NumOther    int                     `json:"num_other" binding:"-"`          // 杂料数
	ColorMetal  enums.ProductColor      `json:"color_metal" binding:"-"`        // 金颜色
	ColorGem    enums.ProductColor      `json:"color_gem" binding:"-"`          // 主石色
	Clarity     enums.ProductClarity    `json:"clarity" binding:"-"`            // 净度
	RetailType  enums.ProductRetailType `json:"retail_type" binding:"required"` // 零售方式
	Class       enums.ProductClass      `json:"class" binding:"required"`       // 大类
	Supplier    enums.ProductSupplier   `json:"supplier" binding:"required"`    // 供应商
	Material    enums.ProductMaterial   `json:"material" binding:"required"`    // 材质
	Quality     enums.ProductQuality    `json:"quality" binding:"required"`     // 成色
	Gem         enums.ProductGem        `json:"gem" binding:"required"`         // 宝石
	Category    enums.ProductCategory   `json:"category" binding:"required"`    // 品类
	Brand       enums.ProductBrand      `json:"brand" binding:"-"`              // 品牌
	Craft       enums.ProductCraft      `json:"craft" binding:"-"`              // 工艺
	Style       string                  `json:"style" binding:"-"`              // 款式
	Size        string                  `json:"size" binding:"-"`               // 手寸

	IsSpecialOffer bool     `json:"is_special_offer" binding:"-"` // 是否特价
	Remark         string   `json:"remark" binding:"-"`           // 备注
	Certificate    []string `json:"certificate" binding:"-"`      // 证书
}

type ProductWhere struct {
	Code string `json:"code" label:"条码" show:"true" sort:"1" type:"string" input:"text" required:"true"` // 条码
	Name string `json:"name" label:"名称" show:"true" sort:"2" type:"string" input:"text" required:"true"` // 名称

	AccessFee float64 `json:"access_fee" label:"入网费" show:"true" sort:"3" type:"float" input:"text" required:"true"` // 入网费
	Price     float64 `json:"price" label:"价格" show:"true" sort:"4" type:"float" input:"text" required:"true"`       // 价格
	LaborFee  float64 `json:"labor_fee" label:"工费" show:"true" sort:"5" type:"float" input:"text" required:"true"`   // 工费

	Weight      float64                 `json:"weight" label:"总重量" show:"true" sort:"6" type:"float" input:"number"`                          // 总重量
	WeightMetal float64                 `json:"weight_metal" label:"金重" show:"true" sort:"7" type:"float" input:"number"`                     // 金重
	WeightGem   float64                 `json:"weight_gem" label:"主石重" show:"true" sort:"8" type:"float" input:"number"`                      // 主石重
	WeightOther float64                 `json:"weight_other" label:"杂料重" show:"true" sort:"9" type:"float" input:"number"`                    // 杂料重
	NumGem      int                     `json:"num_gem" label:"主石数" show:"true" sort:"10" type:"number" input:"number"`                       // 主石数
	NumOther    int                     `json:"num_other" label:"杂料数" show:"true" sort:"11" type:"number" input:"number"`                     // 杂料数
	ColorMetal  enums.ProductColor      `json:"color_metal" label:"金颜色" show:"true" sort:"12" type:"number" input:"select" preset:"typeMap"`  // 金颜色
	ColorGem    enums.ProductColor      `json:"color_gem" label:"主石色" show:"true" sort:"13" type:"number" input:"select" preset:"typeMap"`    // 主石色
	Clarity     enums.ProductClarity    `json:"clarity" label:"净度" show:"true" sort:"14" type:"number" input:"select" preset:"typeMap"`       // 净度
	RetailType  enums.ProductRetailType `json:"retail_type" label:"零售方式" show:"true" sort:"15" type:"number" input:"select" preset:"typeMap"` // 零售方式
	Class       enums.ProductClass      `json:"class" label:"大类" show:"true" sort:"16" type:"number" input:"select" preset:"typeMap"`         // 大类
	Supplier    enums.ProductSupplier   `json:"supplier" label:"供应商" show:"true" sort:"17" type:"number" input:"select" preset:"typeMap"`     // 供应商
	Material    enums.ProductMaterial   `json:"material" label:"材质" show:"true" sort:"18" type:"number" input:"select" preset:"typeMap"`      // 材质
	Quality     enums.ProductQuality    `json:"quality" label:"成色" show:"true" sort:"19" type:"number" input:"select" preset:"typeMap"`       // 成色
	Gem         enums.ProductGem        `json:"gem" label:"宝石" show:"true" sort:"20" type:"number" input:"select" preset:"typeMap"`           // 宝石

	Category enums.ProductCategory `json:"category" label:"品类" show:"true" sort:"21" type:"number" input:"select" preset:"typeMap"` // 品类
	Brand    enums.ProductBrand    `json:"brand" label:"品牌" show:"true" sort:"22" type:"number" input:"select" preset:"typeMap"`    // 品牌
	Craft    enums.ProductCraft    `json:"craft" label:"工艺" show:"true" sort:"23" type:"number" input:"select" preset:"typeMap"`    // 工艺
	Style    string                `json:"style" label:"款式" show:"true" sort:"24" type:"string" input:"text"`                       // 款式
	Size     string                `json:"size" label:"手寸" show:"true" sort:"25" type:"string" input:"text"`                        // 手寸

	IsSpecialOffer bool                `json:"is_special_offer" label:"是否特价" show:"true" sort:"26" type:"bool" input:"switch"`        // 是否特价
	Remark         string              `json:"remark" label:"备注" show:"true" sort:"27" type:"string" input:"textarea"`                // 备注
	Certificate    []string            `json:"certificate" label:"证书" show:"true" sort:"28" type:"string[]" input:"textarea"`         // 证书
	Status         enums.ProductStatus `json:"status" label:"状态" show:"true" sort:"29" type:"number" input:"select" preset:"typeMap"` // 状态
	Type           enums.ProductType   `json:"type" label:"类型" show:"true" sort:"30" type:"number" input:"select" preset:"typeMap"`   // 类型

	ProductEnterId string `json:"product_enter_id" label:"入库单" show:"true" sort:"2" type:"string" input:"text"` // 产品入库单ID
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

	Weight      float64                 `json:"weight"`       // 总重量
	WeightMetal float64                 `json:"weight_metal"` // 金重
	WeightGem   float64                 `json:"weight_gem"`   // 主石重
	WeightOther float64                 `json:"weight_other"` // 杂料重
	NumGem      int                     `json:"num_gem"`      // 主石数
	NumOther    int                     `json:"num_other"`    // 杂料数
	ColorMetal  enums.ProductColor      `json:"color_metal"`  // 金颜色
	ColorGem    enums.ProductColor      `json:"color_gem"`    // 主石色
	Clarity     enums.ProductClarity    `json:"clarity"`      // 净度
	RetailType  enums.ProductRetailType `json:"retail_type"`  // 零售方式
	Class       enums.ProductClass      `json:"class"`        // 大类
	Supplier    enums.ProductSupplier   `json:"supplier"`     // 供应商
	Material    enums.ProductMaterial   `json:"material"`     // 材质
	Quality     enums.ProductQuality    `json:"quality"`      // 成色
	Gem         enums.ProductGem        `json:"gem"`          // 宝石
	Category    enums.ProductCategory   `json:"category"`     // 品类
	Brand       enums.ProductBrand      `json:"brand"`        // 品牌
	Craft       enums.ProductCraft      `json:"craft"`        // 工艺
	Style       string                  `json:"style"`        // 款式
	Size        string                  `json:"size"`         // 手寸

	IsSpecialOffer bool     `json:"is_special_offer"` // 是否特价
	Remark         string   `json:"remark"`           // 备注
	Certificate    []string `json:"certificate"`      // 证书
}

type ProductDamageReq struct {
	Code   string `json:"code" binding:"required"`   // 条码
	Reason string `json:"reason" binding:"required"` // 损坏原因
}

type ProductAllocateCreateReq struct {
	Method  enums.ProductAllocateMethod `json:"method" binding:"required"` // 调拨方式
	Type    enums.ProductType           `json:"type" binding:"required"`   // 仓库类型
	Reason  enums.ProductAllocateReason `json:"reason" binding:"required"` // 调拨原因
	Remark  string                      `json:"remark" binding:"-"`        // 备注
	StoreId string                      `json:"store_id" binding:"-"`      // 调拨门店
}

func (req *ProductAllocateCreateReq) Validate() error {
	if req.Method == enums.ProductAllocateMethodStore && req.StoreId == "" {
		return errors.New("调拨门店不能为空")
	}

	return nil
}

type ProductAllocateWhere struct {
	Method  enums.ProductAllocateMethod `json:"method" label:"调拨方式" input:"select" type:"number" show:"true" sort:"1" required:"true" preset:"typeMap"` // 调拨方式
	Type    enums.ProductType           `json:"type" label:"仓库类型" input:"select" type:"number" show:"true" sort:"2" required:"true" preset:"typeMap"`   // 仓库类型
	Reason  enums.ProductAllocateReason `json:"reason" label:"调拨原因" input:"select" type:"number" show:"true" sort:"3" required:"true" preset:"typeMap"` // 调拨原因
	StoreId string                      `json:"store_id" label:"调拨门店" input:"search" type:"string" show:"true" sort:"4" required:"true"`                // 调拨门店

	StartTime *time.Time `json:"start_time" label:"开始时间" input:"date" type:"date" show:"true" sort:"5" required:"false"` // 开始时间
	EndTime   *time.Time `json:"end_time" label:"结束时间" input:"date" type:"date" show:"true" sort:"6" required:"false"`   // 结束时间
}

func (req *ProductAllocateWhere) Validate() error {
	if req.Method == enums.ProductAllocateMethodStore && req.StoreId == "" {
		return errors.New("调拨门店不能为空")
	}

	return nil
}

type ProductAllocateListReq struct {
	PageReq
	Where ProductAllocateWhere `json:"where"`
}

type ProductAllocateInfoReq struct {
	Id string `json:"id" binding:"required"`
}

type ProductAllocateAddReq struct {
	Id   string `json:"id" binding:"required"`   // 调拨单ID
	Code string `json:"code" binding:"required"` // 产品条码
}

type ProductEnterWhere struct {
	Id        string     `json:"id" label:"ID" input:"text" type:"string" show:"true" sort:"1" required:"false"`         // ID
	StartTime *time.Time `json:"start_time" label:"开始时间" input:"date" type:"date" show:"true" sort:"2" required:"false"` // 开始时间
	EndTime   *time.Time `json:"end_time" label:"结束时间" input:"date" type:"date" show:"true" sort:"3" required:"false"`   // 结束时间
}

type ProductEnterListReq struct {
	PageReq
	Where ProductEnterWhere `json:"where"`
}

type ProductEnterInfoReq struct {
	Id string `json:"id" binding:"required"`
}
