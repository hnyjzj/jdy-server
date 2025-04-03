package types

import (
	"jdy/enums"

	"github.com/shopspring/decimal"
)

type ProductOldWhere struct {
	Code        string                `json:"code" label:"条码" find:"true" create:"true" update:"false" sort:"4" type:"string" input:"text" required:"false"`
	Name        string                `json:"name" label:"名称" find:"true" create:"true" update:"false" sort:"5" type:"string" input:"text" required:"false"`
	Status      enums.ProductStatus   `json:"status" label:"状态" find:"true" create:"true" update:"false" sort:"13" type:"number" input:"select" required:"false" preset:"typeMap"`
	LabelPrice  *decimal.Decimal      `json:"label_price" label:"标签价" find:"true" create:"true" update:"false" sort:"12" type:"number" input:"number" required:"false"`
	Brand       enums.ProductBrand    `json:"brand" label:"品牌" find:"true" create:"true" update:"false" sort:"16" type:"number" input:"select" required:"false" preset:"typeMap"`
	Material    enums.ProductMaterial `json:"material" label:"材质" find:"true" create:"true" update:"false" sort:"6" type:"number" input:"select" required:"true" preset:"typeMap"`
	Quality     enums.ProductQuality  `json:"quality" label:"成色" find:"true" create:"true" update:"false" sort:"7" type:"number" input:"select" required:"true" preset:"typeMap"`
	Gem         enums.ProductGem      `json:"gem" label:"主石" find:"true" create:"true" update:"false" sort:"8" type:"number" input:"select" required:"true" preset:"typeMap"`
	Category    enums.ProductCategory `json:"category" label:"品类" find:"true" create:"true" update:"false" sort:"9" type:"number" input:"select" required:"false" preset:"typeMap"`
	Craft       enums.ProductCraft    `json:"craft" label:"工艺" find:"true" create:"true" update:"false" sort:"10" type:"number" input:"select" required:"false" preset:"typeMap"`
	WeightMetal *decimal.Decimal      `json:"weight_metal" label:"金重" find:"true" create:"true" update:"false" sort:"11" type:"number" input:"number" required:"true"`
	WeightTotal *decimal.Decimal      `json:"weight_total" label:"总重" find:"true" create:"true" update:"false" sort:"24" type:"number" input:"number" required:"false"`
	ColorGem    enums.ProductColor    `json:"color_gem" label:"颜色" find:"true" create:"true" update:"false" sort:"18" type:"number" input:"select" required:"false" preset:"typeMap"`
	WeightGem   *decimal.Decimal      `json:"weight_gem" label:"主石重" find:"true" create:"true" update:"false" sort:"17" type:"number" input:"number" required:"false"`
	NumGem      int                   `json:"num_gem" label:"主石数" find:"true" create:"true" update:"false" sort:"21" type:"number" input:"number" required:"false"`
	Clarity     enums.ProductClarity  `json:"clarity" label:"净度" find:"true" create:"true" update:"false" sort:"19" type:"number" input:"select" required:"false" preset:"typeMap"`
	Cut         enums.ProductCut      `json:"cut" label:"切工" find:"true" create:"true" update:"false" sort:"20" type:"number" input:"select" required:"false" preset:"typeMap"`
	WeightOther *decimal.Decimal      `json:"weight_other" label:"副石重" find:"true" create:"true" update:"false" sort:"22" type:"number" input:"number" required:"false"`
	NumOther    int                   `json:"num_other" label:"副石数" find:"true" create:"true" update:"false" sort:"23" type:"number" input:"number" required:"false"`
	Remark      string                `json:"remark" label:"备注" find:"true" create:"true" update:"false" sort:"26" type:"string" input:"textarea" required:"false"`

	StoreId string `json:"store_id" label:"所属店铺" find:"false" create:"true" update:"false" sort:"28" type:"string" input:"text" required:"false"`

	IsOur                   *bool                      `json:"is_our" label:"是否自有" find:"true" create:"true" update:"false" sort:"1" type:"boolean" input:"switch" required:"true"`
	RecycleMethod           enums.ProductRecycleMethod `json:"recycle_method" label:"回收方式" find:"true" create:"true" update:"false" sort:"2" type:"number" input:"select" required:"true" preset:"typeMap"`
	RecycleType             enums.ProductRecycleType   `json:"recycle_type" label:"回收类型" find:"true" create:"true" update:"false" sort:"3" type:"number" input:"select" preset:"typeMap"`
	RecyclePrice            *decimal.Decimal           `json:"recycle_price" label:"回收金额" find:"true" create:"true" update:"false" sort:"27" type:"number" input:"number" required:"false"`
	RecyclePriceGold        *decimal.Decimal           `json:"recycle_price_gold" label:"回收金价" find:"true" create:"true" update:"false" sort:"13" type:"number" input:"number" required:"false"`
	RecyclePriceLabor       *decimal.Decimal           `json:"recycle_price_labor" label:"回收工费" find:"true" create:"true" update:"false" sort:"15" type:"number" input:"number" required:"false"`
	RecyclePriceLaborMethod enums.ProductRecycleMethod `json:"recycle_price_labor_method" label:"回收工费方式" find:"true" create:"true" update:"false" sort:"14" type:"number" input:"select" required:"false" preset:"typeMap"`
	QualityActual           *decimal.Decimal           `json:"quality_actual" label:"实际成色" find:"true" create:"true" update:"false" sort:"25" type:"number" input:"number" required:"false"`
	RecycleSource           enums.ProductRecycleSource `json:"recycle_source" label:"回收来源" find:"true" create:"true" update:"false" sort:"4" type:"number" input:"select" required:"true" preset:"typeMap"`
	RecycleSourceId         string                     `json:"recycle_source_id" label:"回收来源" find:"true" create:"true" update:"false" sort:"5" type:"string" input:"text" required:"false"`
	RecycleStoreId          string                     `json:"recycle_store_id" label:"回收店铺" find:"true" create:"true" update:"false" sort:"6" type:"string" input:"text" required:"false"`
}

type ProductOldListReq struct {
	PageReq
	Where ProductOldWhere `json:"where" binding:"required"`
}

type ProductOldInfoReq struct {
	Id string `json:"id" binding:"required"` // 产品ID
}

type ProductOldUpdateReq struct {
	Id string `json:"id" binding:"required"` // 产品ID
	ProductOldWhere
}
