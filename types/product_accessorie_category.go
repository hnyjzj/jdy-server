package types

import (
	"jdy/enums"

	"github.com/shopspring/decimal"
)

// 配件条目条件
type ProductAccessorieCategoryWhere struct {
	Id            string                            `json:"id" label:"配件编号" input:"text" type:"string" find:"true" create:"false" sort:"1" required:"false"`
	TypePart      enums.ProductTypePart             `json:"type_part" label:"配件类型" input:"select" type:"number" find:"true" create:"true" sort:"1" required:"true" preset:"typeMap"`
	Name          string                            `json:"name" label:"配件名称" input:"text" type:"string" find:"true" create:"true" sort:"2" required:"true"`
	Code          string                            `json:"code" label:"配件条码" input:"text" type:"string" find:"true" create:"true" sort:"3" required:"false"`
	RetailType    enums.ProductAccessorieRetailType `json:"retail_type" label:"零售方式" input:"select" type:"number" find:"true" create:"true" sort:"4" required:"true" preset:"typeMap"`
	Weight        *decimal.Decimal                  `json:"weight" label:"重量" input:"number" type:"number" find:"true" create:"true" sort:"5" required:"false"`
	AccessFee     *decimal.Decimal                  `json:"access_fee" label:"入网费" input:"number" type:"number" find:"true" create:"true" sort:"6" required:"false"`
	LabelPrice    *decimal.Decimal                  `json:"label_price" label:"标签价" input:"number" type:"number" find:"true" create:"true" sort:"7" required:"false"`
	Material      enums.ProductAccessorieMaterial   `json:"material" label:"材质" input:"select" type:"number" find:"true" create:"true" sort:"8" required:"false" preset:"typeMap"`
	Quality       enums.ProductQuality              `json:"quality" label:"成色" input:"select" type:"number" find:"true" create:"true" sort:"9" required:"false" preset:"typeMap"`
	Gem           enums.ProductGem                  `json:"gem" label:"主石" input:"select" type:"number" find:"true" create:"true" sort:"10" required:"false" preset:"typeMap"`
	Category      enums.ProductCategory             `json:"category" label:"品类" input:"select" type:"number" find:"true" create:"true" sort:"11" required:"false" preset:"typeMap"`
	Specification string                            `json:"specification" label:"规格" input:"text" type:"string" find:"true" create:"true" sort:"12" required:"false"`
	Color         string                            `json:"color" label:"颜色" input:"text" type:"string" find:"true" create:"true" sort:"13" required:"false"`
	Series        string                            `json:"series" label:"系列" input:"text" type:"string" find:"true" create:"true" sort:"14" required:"false"`
	Supplier      string                            `json:"supplier" label:"供应商" input:"text" type:"string" find:"true" create:"true" sort:"15" required:"false"`
	Remark        string                            `json:"remark" label:"备注" input:"text" type:"string" find:"true" create:"true" sort:"16" required:"false"`
}

// 配件条目列表
type ProductAccessorieCategoryListReq struct {
	PageReq
	Where ProductAccessorieCategoryWhere `json:"where" binding:"required"`
}

// 配件条目详情
type ProductAccessorieCategoryInfoReq struct {
	Id string `json:"id" binding:"required"` // 条目ID
}

type ProductAccessorieCategoryCreateWhere struct {
	TypePart      enums.ProductTypePart             `json:"type_part" binding:"required"`
	Name          string                            `json:"name" binding:"required"`
	Code          string                            `json:"code" `
	RetailType    enums.ProductAccessorieRetailType `json:"retail_type" binding:"required"`
	Weight        *decimal.Decimal                  `json:"weight" `
	AccessFee     *decimal.Decimal                  `json:"access_fee" `
	LabelPrice    *decimal.Decimal                  `json:"label_price" `
	Material      enums.ProductAccessorieMaterial   `json:"material" `
	Quality       enums.ProductQuality              `json:"quality" `
	Gem           enums.ProductGem                  `json:"gem" `
	Category      enums.ProductCategory             `json:"category" `
	Specification string                            `json:"specification" `
	Color         string                            `json:"color" `
	Series        string                            `json:"series" `
	Supplier      string                            `json:"supplier" `
	Remark        string                            `json:"remark" `
}

type ProductAccessorieCategoryCreateReq struct {
	List []ProductAccessorieCategoryCreateWhere `json:"list" binding:"required"`
}

// 配件条目更新
type ProductAccessorieCategoryUpdateReq struct {
	Id string `json:"id" binding:"required"` // 条目ID
	ProductAccessorieCategoryWhere
}

// 配件条目删除
type ProductAccessorieCategoryDeleteReq struct {
	Id string `json:"id" binding:"required"` // 条目ID
}
