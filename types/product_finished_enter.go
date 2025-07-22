package types

import (
	"jdy/enums"
	"time"

	"github.com/shopspring/decimal"
)

type ProductFinishedEnterCreateReq struct {
	StoreId string `json:"store_id" binding:"required"` // 门店ID
	Remark  string `json:"remark"`                      // 备注
}

type ProductFinishedEnterWhere struct {
	Id      string                   `json:"id" label:"单号" input:"text" type:"string" find:"true" sort:"1" required:"false"`                                       // ID
	StoreId string                   `json:"store_id" label:"门店" input:"text" type:"string" find:"false" create:"true" sort:"2" required:"true"`                   // 门店ID
	Status  enums.ProductEnterStatus `json:"status" label:"状态" input:"select" type:"number" find:"false" create:"true" sort:"3" required:"false" preset:"typeMap"` // 状态
	Remark  string                   `json:"remark" label:"备注" input:"text" type:"string" find:"false" create:"true" sort:"4" required:"false"`                    // 备注
	StartAt *time.Time               `json:"start_at" label:"开始时间" input:"date" type:"time" find:"true" sort:"5" required:"false"`                                 // 开始时间
	EndAt   *time.Time               `json:"end_at" label:"结束时间" input:"date" type:"time" find:"true" sort:"6" required:"false"`                                   // 结束时间
}

type ProductFinishedEnterListReq struct {
	PageReq
	Where ProductFinishedEnterWhere `json:"where"`
}

type ProductFinishedEnterInfoReq struct {
	Id string `json:"id" binding:"required"`
	PageReq
}

type ProductFinishedEnterAddProductReq struct {
	EnterId  string                           `json:"enter_id" binding:"required"` // 入库单ID
	Products []ProductFinishedEnterReqProduct `json:"products" binding:"required"` // 商品信息
}

type ProductFinishedEnterEditProductReq struct {
	EnterId   string                         `json:"enter_id" binding:"required"`   // 入库单ID
	ProductId string                         `json:"product_id" binding:"required"` // 商品ID
	Product   ProductFinishedEnterReqProduct `json:"product" binding:"-"`           // 商品信息
}

type ProductFinishedEnterDelProductReq struct {
	EnterId    string   `json:"enter_id" binding:"required"`    // 入库单ID
	ProductIds []string `json:"product_ids" binding:"required"` // 商品ID列表
}

type ProductFinishedEnterClearProductReq struct {
	EnterId string `json:"enter_id" binding:"required"` // 入库单ID
}

type ProductFinishedEnterFinishReq struct {
	EnterId string `json:"enter_id" binding:"required"` // 入库单ID
}

type ProductFinishedEnterCancelReq struct {
	EnterId string `json:"enter_id" binding:"required"` // 入库单ID
}

type ProductFinishedEnterReqProduct struct {
	Code string `json:"code" binding:"required"` // 条码
	Name string `json:"name" binding:"required"` // 名称

	AccessFee  *decimal.Decimal `json:"access_fee" binding:"required"`  // 入网费
	LabelPrice *decimal.Decimal `json:"label_price" binding:"required"` // 标签价
	LaborFee   *decimal.Decimal `json:"labor_fee" binding:"required"`   // 工费

	WeightTotal *decimal.Decimal        `json:"weight_total"`                    // 总重量
	WeightMetal *decimal.Decimal        `json:"weight_metal" binding:"required"` // 金重
	WeightGem   *decimal.Decimal        `json:"weight_gem"`                      // 主石重
	WeightOther *decimal.Decimal        `json:"weight_other"`                    // 杂料重
	NumGem      int                     `json:"num_gem"`                         // 主石数
	NumOther    int                     `json:"num_other"`                       // 杂料数
	ColorMetal  string                  `json:"color_metal"`                     // 贵金属颜色
	ColorGem    enums.ProductColor      `json:"color_gem"`                       // 颜色
	Clarity     enums.ProductClarity    `json:"clarity"`                         // 净度
	RetailType  enums.ProductRetailType `json:"retail_type" binding:"required"`  // 零售方式
	Supplier    enums.ProductSupplier   `json:"supplier" binding:"required"`     // 供应商
	Material    enums.ProductMaterial   `json:"material" binding:"required"`     // 材质
	Quality     enums.ProductQuality    `json:"quality" binding:"required"`      // 成色
	Gem         enums.ProductGem        `json:"gem" binding:"required"`          // 主石

	Category enums.ProductCategory `json:"category" binding:"required"` // 品类
	Brand    enums.ProductBrand    `json:"brand"`                       // 品牌
	Craft    enums.ProductCraft    `json:"craft"`                       // 工艺
	Style    string                `json:"style"`                       // 款式
	Size     string                `json:"size"`                        // 手寸
	Series   string                `json:"series"`                      // 系列
	Remark   string                `json:"remark"`                      // 备注

	IsSpecialOffer bool     `json:"is_special_offer"` // 是否特价
	Certificate    []string `json:"certificate"`      // 证书

	EnterTime *time.Time `json:"enter_time"` // 入库时间
}
