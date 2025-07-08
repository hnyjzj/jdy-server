package types

import (
	"jdy/enums"

	"github.com/shopspring/decimal"
)

type ProductFinishedWhere struct {
	Code string `json:"code" label:"条码" find:"true" create:"true" update:"true" sort:"1" type:"string" input:"text" required:"true"` // 条码
	Name string `json:"name" label:"名称" find:"true" create:"true" update:"true" sort:"2" type:"string" input:"text" required:"true"` // 名称

	AccessFee  *decimal.Decimal `json:"access_fee" label:"入网费" find:"true" create:"true" update:"true" sort:"3" type:"float" input:"number" required:"true"`  // 入网费
	LabelPrice *decimal.Decimal `json:"label_price" label:"标签价" find:"true" create:"true" update:"true" sort:"4" type:"float" input:"number" required:"true"` // 标签价
	LaborFee   *decimal.Decimal `json:"labor_fee" label:"工费" find:"true" create:"true" update:"true" sort:"5" type:"float" input:"number" required:"true"`    // 工费

	WeightTotal *decimal.Decimal           `json:"weight_total" label:"总重量" find:"true" create:"true" update:"true" sort:"6" type:"float" input:"number"`                    // 总重量
	WeightMetal *decimal.Decimal           `json:"weight_metal" label:"金重" find:"true" create:"true" update:"true" sort:"7" type:"float" input:"number"`                     // 金重
	WeightGem   *decimal.Decimal           `json:"weight_gem" label:"主石重" find:"true" create:"true" update:"true" sort:"8" type:"float" input:"number"`                      // 主石重
	WeightOther *decimal.Decimal           `json:"weight_other" label:"杂料重" find:"true" create:"true" update:"true" sort:"9" type:"float" input:"number"`                    // 杂料重
	NumGem      int                        `json:"num_gem" label:"主石数" find:"true" create:"true" update:"true" sort:"10" type:"number" input:"number"`                       // 主石数
	NumOther    int                        `json:"num_other" label:"杂料数" find:"true" create:"true" update:"true" sort:"11" type:"number" input:"number"`                     // 杂料数
	ColorMetal  string                     `json:"color_metal" label:"贵金属颜色" find:"true" create:"true" update:"true" sort:"12" type:"string" input:"text"`                   // 贵金属颜色
	ColorGem    enums.ProductColor         `json:"color_gem" label:"颜色" find:"true" create:"true" update:"true" sort:"13" type:"number" input:"select" preset:"typeMap"`     // 颜色
	Clarity     enums.ProductClarity       `json:"clarity" label:"净度" find:"true" create:"true" update:"true" sort:"14" type:"number" input:"select" preset:"typeMap"`       // 净度
	RetailType  enums.ProductRetailType    `json:"retail_type" label:"零售方式" find:"true" create:"true" update:"true" sort:"15" type:"number" input:"select" preset:"typeMap"` // 零售方式
	Class       enums.ProductClassFinished `json:"class" label:"大类" find:"true" create:"false" update:"false" sort:"16" type:"number" input:"select" preset:"typeMap"`       // 大类
	Supplier    enums.ProductSupplier      `json:"supplier" label:"供应商" find:"true" create:"true" update:"true" sort:"17" type:"number" input:"select" preset:"typeMap"`     // 供应商
	Material    enums.ProductMaterial      `json:"material" label:"材质" find:"true" create:"true" update:"true" sort:"18" type:"number" input:"select" preset:"typeMap"`      // 材质
	Quality     enums.ProductQuality       `json:"quality" label:"成色" find:"true" create:"true" update:"true" sort:"19" type:"number" input:"select" preset:"typeMap"`       // 成色
	Gem         enums.ProductGem           `json:"gem" label:"宝石" find:"true" create:"true" update:"true" sort:"20" type:"number" input:"select" preset:"typeMap"`           // 宝石

	Category enums.ProductCategory `json:"category" label:"品类" find:"true" create:"true" update:"true" sort:"21" type:"number" input:"select" preset:"typeMap"` // 品类
	Brand    enums.ProductBrand    `json:"brand" label:"品牌" find:"true" create:"true" update:"true" sort:"22" type:"number" input:"select" preset:"typeMap"`    // 品牌
	Craft    enums.ProductCraft    `json:"craft" label:"工艺" find:"true" create:"true" update:"true" sort:"23" type:"number" input:"select" preset:"typeMap"`    // 工艺
	Style    string                `json:"style" label:"款式" find:"true" create:"true" update:"true" sort:"24" type:"string" input:"text"`                       // 款式
	Size     string                `json:"size" label:"手寸" find:"true" create:"true" update:"true" sort:"25" type:"string" input:"text"`                        // 手寸

	IsSpecialOffer *bool               `json:"is_special_offer" label:"是否特价" find:"true" create:"true" update:"true" sort:"26" type:"boolean" input:"switch"`       // 是否特价
	Series         string              `json:"series" label:"系列" find:"true" create:"true" update:"true" sort:"26" type:"string" input:"text"`                      // 系列
	Remark         string              `json:"remark" label:"备注" find:"true" create:"true" update:"true" sort:"27" type:"string" input:"textarea"`                  // 备注
	Status         enums.ProductStatus `json:"status" label:"状态" find:"true" create:"false" update:"false" sort:"28" type:"number" input:"select" preset:"typeMap"` // 状态

	StoreId     string   `json:"store_id" label:"门店" find:"false" create:"true" update:"false" sort:"31" type:"string" input:"text" binding:"required"` // 门店
	Certificate []string `json:"certificate" label:"证书" find:"false" create:"true" update:"true" sort:"32" type:"string[]" input:"list"`                // 证书
	Images      []string `json:"images" label:"图片" find:"false" create:"false" update:"true" sort:"33" type:"string[]" input:"list"`                    // 图片

	EnterId string `json:"enter_id" label:"入库单" find:"true" sort:"2" type:"string" input:"text"` // 产品入库单ID

	All bool `json:"all" label:"全部" find:"true" sort:"1" type:"boolean" input:"switch"` // 全部
}

type ProductFinishedListReq struct {
	PageReq
	Where ProductFinishedWhere `json:"where" binding:"required"`
}

type ProductFinishedInfoReq struct {
	Code string `json:"code" binding:"required"` // 条码
}

type ProductFinishedUpdateReq struct {
	Id string `json:"id" binding:"required"` // ID
	ProductFinishedWhere
}

type ProductFinishedUploadReq struct {
	Id     string   `json:"id" binding:"required"`     // ID
	Images []string `json:"images" binding:"required"` // 图片
}
