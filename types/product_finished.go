package types

import (
	"errors"
	"jdy/enums"
	"time"

	"github.com/shopspring/decimal"
)

type ProductFinishedWhere struct {
	Code string `json:"code" label:"条码" find:"true" info:"true" create:"true" update:"true" sort:"1" type:"string" input:"text" required:"true"` // 条码
	Name string `json:"name" label:"名称" find:"true" info:"true" create:"true" update:"true" sort:"2" type:"string" input:"text" required:"true"` // 名称

	AccessFee  *decimal.Decimal `json:"access_fee" label:"入网费" find:"true" info:"true" create:"true" update:"true" sort:"3" type:"float" input:"text" required:"true"`  // 入网费
	LabelPrice *decimal.Decimal `json:"label_price" label:"标签价" find:"true" info:"true" create:"true" update:"true" sort:"4" type:"float" input:"text" required:"true"` // 标签价
	LaborFee   *decimal.Decimal `json:"labor_fee" label:"工费" find:"true" info:"true" create:"true" update:"true" sort:"5" type:"float" input:"text" required:"true"`    // 工费

	WeightTotal *decimal.Decimal           `json:"weight_total" label:"总重量" find:"true" info:"true" create:"true" update:"true" sort:"6" type:"float" input:"text"`                      // 总重量
	WeightMetal *decimal.Decimal           `json:"weight_metal" label:"金重" find:"true" info:"true" create:"true" update:"true" sort:"7" type:"float" input:"text"`                       // 金重
	WeightGem   *decimal.Decimal           `json:"weight_gem" label:"主石重" find:"true" info:"true" create:"true" update:"true" sort:"8" type:"float" input:"text"`                        // 主石重
	WeightOther *decimal.Decimal           `json:"weight_other" label:"杂料重" find:"true" info:"true" create:"true" update:"true" sort:"9" type:"float" input:"text"`                      // 杂料重
	NumGem      int                        `json:"num_gem" label:"主石数" find:"true" info:"true" create:"true" update:"true" sort:"10" type:"number" input:"number"`                       // 主石数
	NumOther    int                        `json:"num_other" label:"杂料数" find:"true" info:"true" create:"true" update:"true" sort:"11" type:"number" input:"number"`                     // 杂料数
	ColorMetal  string                     `json:"color_metal" label:"贵金属颜色" find:"true" info:"true" create:"true" update:"true" sort:"12" type:"string" input:"text"`                   // 贵金属颜色
	ColorGem    enums.ProductColor         `json:"color_gem" label:"颜色" find:"true" info:"true" create:"true" update:"true" sort:"13" type:"number" input:"select" preset:"typeMap"`     // 颜色
	Clarity     enums.ProductClarity       `json:"clarity" label:"净度" find:"true" info:"true" create:"true" update:"true" sort:"14" type:"number" input:"select" preset:"typeMap"`       // 净度
	RetailType  enums.ProductRetailType    `json:"retail_type" label:"零售方式" find:"true" info:"true" create:"true" update:"true" sort:"15" type:"number" input:"select" preset:"typeMap"` // 零售方式
	Class       enums.ProductClassFinished `json:"class" label:"大类" find:"true" info:"true" create:"false" update:"false" sort:"16" type:"number" input:"select" preset:"typeMap"`       // 大类
	Supplier    enums.ProductSupplier      `json:"supplier" label:"供应商" find:"true" info:"true" create:"true" update:"true" sort:"17" type:"number" input:"select" preset:"typeMap"`     // 供应商
	Material    enums.ProductMaterial      `json:"material" label:"材质" find:"true" info:"true" create:"true" update:"true" sort:"18" type:"number" input:"select" preset:"typeMap"`      // 材质
	Quality     enums.ProductQuality       `json:"quality" label:"成色" find:"true" info:"true" create:"true" update:"true" sort:"19" type:"number" input:"select" preset:"typeMap"`       // 成色
	Gem         enums.ProductGem           `json:"gem" label:"宝石" find:"true" info:"true" create:"true" update:"true" sort:"20" type:"number" input:"select" preset:"typeMap"`           // 宝石

	Category enums.ProductCategory `json:"category" label:"品类" find:"true" info:"true" create:"true" update:"true" sort:"21" type:"number" input:"select" preset:"typeMap"` // 品类
	Brand    enums.ProductBrand    `json:"brand" label:"品牌" find:"true" info:"true" create:"true" update:"true" sort:"22" type:"number" input:"select" preset:"typeMap"`    // 品牌
	Craft    enums.ProductCraft    `json:"craft" label:"工艺" find:"true" info:"true" create:"true" update:"true" sort:"23" type:"number" input:"select" preset:"typeMap"`    // 工艺
	Style    string                `json:"style" label:"款式" find:"true" info:"true" create:"true" update:"true" sort:"24" type:"string" input:"text"`                       // 款式
	Size     string                `json:"size" label:"手寸" find:"true" info:"true" create:"true" update:"true" sort:"25" type:"string" input:"text"`                        // 手寸

	IsSpecialOffer *bool               `json:"is_special_offer" label:"是否特价" find:"true" info:"true" create:"true" update:"true" sort:"26" type:"boolean" input:"switch"`       // 是否特价
	Series         string              `json:"series" label:"系列" find:"true" info:"true" create:"true" update:"true" sort:"27" type:"string" input:"text"`                      // 系列
	Remark         string              `json:"remark" label:"备注" find:"true" info:"true" create:"true" update:"true" sort:"28" type:"string" input:"textarea"`                  // 备注
	Status         enums.ProductStatus `json:"status" label:"状态" find:"true" info:"true" create:"false" update:"false" sort:"29" type:"number" input:"select" preset:"typeMap"` // 状态

	StoreId     string   `json:"store_id" label:"门店" find:"false" create:"true" update:"false" sort:"30" type:"string" input:"text"`                 // 门店
	Certificate []string `json:"certificate" label:"证书" find:"false" create:"true" update:"true" info:"true" sort:"31" type:"string[]" input:"list"` // 证书
	Images      []string `json:"images" label:"图片" find:"false" create:"false" update:"true" sort:"32" type:"string[]" input:"list"`                 // 图片

	EnterId   string    `json:"enter_id" label:"入库单" find:"true" info:"true" sort:"33" type:"string" input:"text"` // 产品入库单ID
	EnterTime time.Time `json:"enter_time" label:"入库时间" find:"false" sort:"34" type:"date" input:"date"`           // 产品入库时间

	StartTime *time.Time `json:"start_time" label:"入库时间（开始）" find:"true" sort:"35" type:"date" input:"date" required:"false"` // 入库时间（开始）
	EndTime   *time.Time `json:"end_time" label:"入库时间（结束）" find:"true" sort:"36" type:"date" input:"date" required:"false"`   // 入库时间（结束）
}

type ProductFinishedListReq struct {
	PageReq
	Where ProductFinishedWhere `json:"where" binding:"required"`
}

type ProductFinishedListRes[T any] struct {
	PageRes[T]

	AccessFee   decimal.Decimal `json:"access_fee"`   // 入网费
	LabelPrice  decimal.Decimal `json:"label_price"`  // 标签价
	WeightMetal decimal.Decimal `json:"weight_metal"` // 金重
}

type ProductFinishedEmptyImageReq struct {
	StoreId string `json:"store_id" binding:"required"` // 门店
}

type ProductFinishedInfoReq struct {
	Code string `json:"code" binding:"required"` // 条码
}

type ProductFinishedRetrievalReq struct {
	Code    string `json:"code" binding:"required"` // 条码
	StoreId string `json:"store_id"`                // 门店
}

type ProductFinishedUpdateReq struct {
	Id string `json:"id" binding:"required"` // ID
	ProductFinishedWhere
}

type ProductFinishedUpdatesReq struct {
	Data []ProductFinishedWhere `json:"data" binding:"required"`
}

func (r *ProductFinishedUpdatesReq) Validate() error {
	for _, v := range r.Data {
		if v.Code == "" {
			return errors.New("条码不能为空")
		}
	}

	return nil
}

type ProductFinishedUpdateCodeReq struct {
	Data []ProductFinishedUpdateCode `json:"data" binding:"required"`
}

type ProductFinishedUpdateCode struct {
	Code    string `json:"code" binding:"required"`     // 条码
	NewCode string `json:"new_code" binding:"required"` // 新条码
}

type ProductFinishedUploadReq struct {
	Id     string   `json:"id" binding:"required"`     // ID
	Images []string `json:"images" binding:"required"` // 图片
}

type ProductFinishedFindCodeReq struct {
	Codes []string `json:"codes" binding:"required"` // 条码
}
