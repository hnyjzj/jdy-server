package types

import (
	"errors"
	"jdy/enums"
	"time"

	"github.com/shopspring/decimal"
)

type ProductEnterReq struct {
	StoreId  string                   `json:"store_id" binding:"required"` // 门店ID
	Products []ProductEnterReqProduct `json:"products" binding:"required"`
}

func (p *ProductEnterReq) Validate() error {
	if len(p.Products) == 0 {
		return errors.New("products is required")
	}
	for _, product := range p.Products {
		if product.Code == "" {
			return errors.New("code is required")
		}
	}

	return nil
}

type ProductEnterReqProduct struct {
	Code string `json:"code" binding:"required"` // 条码
	Name string `json:"name" binding:"required"` // 名称

	AccessFee decimal.Decimal `json:"access_fee" binding:"required"` // 入网费
	Price     decimal.Decimal `json:"price" binding:"required"`      // 标签价
	LaborFee  decimal.Decimal `json:"labor_fee" binding:"required"`  // 工费

	Weight      decimal.Decimal         `json:"weight"`                         // 总重量
	WeightMetal decimal.Decimal         `json:"weight_metal"`                   // 金重
	WeightGem   decimal.Decimal         `json:"weight_gem"`                     // 主石重
	WeightOther decimal.Decimal         `json:"weight_other"`                   // 杂料重
	NumGem      int                     `json:"num_gem"`                        // 主石数
	NumOther    int                     `json:"num_other"`                      // 杂料数
	ColorMetal  enums.ProductColor      `json:"color_metal"`                    // 金颜色
	ColorGem    enums.ProductColor      `json:"color_gem"`                      // 主石色
	Clarity     enums.ProductClarity    `json:"clarity"`                        // 净度
	RetailType  enums.ProductRetailType `json:"retail_type" binding:"required"` // 零售方式
	Class       enums.ProductClass      `json:"class" binding:"required"`       // 大类
	Supplier    enums.ProductSupplier   `json:"supplier" binding:"required"`    // 供应商
	Material    enums.ProductMaterial   `json:"material" binding:"required"`    // 材质
	Quality     enums.ProductQuality    `json:"quality" binding:"required"`     // 成色
	Gem         enums.ProductGem        `json:"gem" binding:"required"`         // 宝石
	Category    enums.ProductCategory   `json:"category" binding:"required"`    // 品类
	Brand       enums.ProductBrand      `json:"brand"`                          // 品牌
	Craft       enums.ProductCraft      `json:"craft"`                          // 工艺
	Style       string                  `json:"style"`                          // 款式
	Size        string                  `json:"size"`                           // 手寸
	Type        enums.ProductType       `json:"type"`                           // 类型

	Stock int64 `json:"stock" binding:"required"` // 库存

	IsSpecialOffer bool     `json:"is_special_offer"` // 是否特价
	Remark         string   `json:"remark"`           // 备注
	Certificate    []string `json:"certificate"`      // 证书
}

type ProductWhere struct {
	Code string `json:"code" label:"条码" find:"true" create:"true" update:"true" sort:"1" type:"string" input:"text" required:"true"` // 条码
	Name string `json:"name" label:"名称" find:"true" create:"true" update:"true" sort:"2" type:"string" input:"text" required:"true"` // 名称

	AccessFee decimal.Decimal `json:"access_fee" label:"入网费" find:"true" create:"true" update:"true" sort:"3" type:"float" input:"text" required:"true"` // 入网费
	Price     decimal.Decimal `json:"price" label:"价格" find:"true" create:"true" update:"true" sort:"4" type:"float" input:"text" required:"true"`       // 价格
	LaborFee  decimal.Decimal `json:"labor_fee" label:"工费" find:"true" create:"true" update:"true" sort:"5" type:"float" input:"text" required:"true"`   // 工费

	Weight      decimal.Decimal         `json:"weight" label:"总重量" find:"true" create:"true" update:"true" sort:"6" type:"float" input:"number"`                          // 总重量
	WeightMetal decimal.Decimal         `json:"weight_metal" label:"金重" find:"true" create:"true" update:"true" sort:"7" type:"float" input:"number"`                     // 金重
	WeightGem   decimal.Decimal         `json:"weight_gem" label:"主石重" find:"true" create:"true" update:"true" sort:"8" type:"float" input:"number"`                      // 主石重
	WeightOther decimal.Decimal         `json:"weight_other" label:"杂料重" find:"true" create:"true" update:"true" sort:"9" type:"float" input:"number"`                    // 杂料重
	NumGem      int                     `json:"num_gem" label:"主石数" find:"true" create:"true" update:"true" sort:"10" type:"number" input:"number"`                       // 主石数
	NumOther    int                     `json:"num_other" label:"杂料数" find:"true" create:"true" update:"true" sort:"11" type:"number" input:"number"`                     // 杂料数
	ColorMetal  enums.ProductColor      `json:"color_metal" label:"金颜色" find:"true" create:"true" update:"true" sort:"12" type:"number" input:"select" preset:"typeMap"`  // 金颜色
	ColorGem    enums.ProductColor      `json:"color_gem" label:"主石色" find:"true" create:"true" update:"true" sort:"13" type:"number" input:"select" preset:"typeMap"`    // 主石色
	Clarity     enums.ProductClarity    `json:"clarity" label:"净度" find:"true" create:"true" update:"true" sort:"14" type:"number" input:"select" preset:"typeMap"`       // 净度
	RetailType  enums.ProductRetailType `json:"retail_type" label:"零售方式" find:"true" create:"true" update:"true" sort:"15" type:"number" input:"select" preset:"typeMap"` // 零售方式
	Class       enums.ProductClass      `json:"class" label:"大类" find:"true" create:"true" update:"true" sort:"16" type:"number" input:"select" preset:"typeMap"`         // 大类
	Supplier    enums.ProductSupplier   `json:"supplier" label:"供应商" find:"true" create:"true" update:"true" sort:"17" type:"number" input:"select" preset:"typeMap"`     // 供应商
	Material    enums.ProductMaterial   `json:"material" label:"材质" find:"true" create:"true" update:"true" sort:"18" type:"number" input:"select" preset:"typeMap"`      // 材质
	Quality     enums.ProductQuality    `json:"quality" label:"成色" find:"true" create:"true" update:"true" sort:"19" type:"number" input:"select" preset:"typeMap"`       // 成色
	Gem         enums.ProductGem        `json:"gem" label:"宝石" find:"true" create:"true" update:"true" sort:"20" type:"number" input:"select" preset:"typeMap"`           // 宝石

	Category enums.ProductCategory `json:"category" label:"品类" find:"true" create:"true" update:"true" sort:"21" type:"number" input:"select" preset:"typeMap"` // 品类
	Brand    enums.ProductBrand    `json:"brand" label:"品牌" find:"true" create:"true" update:"true" sort:"22" type:"number" input:"select" preset:"typeMap"`    // 品牌
	Craft    enums.ProductCraft    `json:"craft" label:"工艺" find:"true" create:"true" update:"true" sort:"23" type:"number" input:"select" preset:"typeMap"`    // 工艺
	Style    string                `json:"style" label:"款式" find:"true" create:"true" update:"true" sort:"24" type:"string" input:"text"`                       // 款式
	Size     string                `json:"size" label:"手寸" find:"true" create:"true" update:"true" sort:"25" type:"string" input:"text"`                        // 手寸

	IsSpecialOffer bool                `json:"is_special_offer" label:"是否特价" find:"true" create:"true" update:"true" sort:"26" type:"bool" input:"switch"`          // 是否特价
	Remark         string              `json:"remark" label:"备注" find:"true" create:"true" update:"true" sort:"27" type:"string" input:"textarea"`                  // 备注
	Certificate    []string            `json:"certificate" label:"证书" find:"true" create:"true" update:"true" sort:"28" type:"string[]" input:"textarea"`           // 证书
	Status         enums.ProductStatus `json:"status" label:"状态" find:"true" create:"false" update:"false" sort:"29" type:"number" input:"select" preset:"typeMap"` // 状态
	Type           enums.ProductType   `json:"type" label:"类型" find:"true" create:"true" update:"false" sort:"30" type:"number" input:"select" preset:"typeMap"`    // 类型

	Stock int64 `json:"stock" label:"库存" find:"false" create:"true" update:"false" sort:"31" type:"number" input:"number"` // 库存

	StoreId string `json:"store_id" label:"门店" find:"true" create:"true" update:"false" sort:"32" type:"string" input:"text" binding:"required"` // 门店

	ProductEnterId string `json:"product_enter_id" label:"入库单" find:"true" sort:"2" type:"string" input:"text"` // 产品入库单ID
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

	AccessFee decimal.Decimal `json:"access_fee"` // 入网费
	Price     decimal.Decimal `json:"price"`      // 标签价
	LaborFee  decimal.Decimal `json:"labor_fee"`  // 工费

	Weight      decimal.Decimal         `json:"weight"`       // 总重量
	WeightMetal decimal.Decimal         `json:"weight_metal"` // 金重
	WeightGem   decimal.Decimal         `json:"weight_gem"`   // 主石重
	WeightOther decimal.Decimal         `json:"weight_other"` // 杂料重
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

type ProductConversionReq struct {
	Code string            `json:"code" binding:"required"` // 条码
	Type enums.ProductType `json:"type" binding:"required"` // 仓库类型
}

type ProductHistoryWhere struct {
	ProductId string              `json:"product_id" label:"产品" input:"text" type:"string" find:"true" sort:"1" required:"true"`                   // 产品
	StoreId   string              `json:"store_id" label:"门店" input:"text" type:"string" find:"false" sort:"2" required:"true" binding:"required"` // 门店
	Action    enums.ProductAction `json:"action" label:"操作" input:"select" type:"number" find:"true" sort:"3" required:"true" preset:"typeMap"`    // 操作
}

type ProductHistoryReq struct {
	PageReq
	Where ProductHistoryWhere `json:"where" binding:"required"`
}

type ProductAllocateCreateReq struct {
	Method      enums.ProductAllocateMethod `json:"method" binding:"required"` // 调拨方式
	Type        enums.ProductType           `json:"type" binding:"required"`   // 仓库类型
	Reason      enums.ProductAllocateReason `json:"reason" binding:"required"` // 调拨原因
	Remark      string                      `json:"remark"`                    // 备注
	FromStoreId string                      `json:"from_store_id"`             // 调出门店
	ToStoreId   string                      `json:"to_store_id"`               // 调入门店
}

func (req *ProductAllocateCreateReq) Validate() error {
	if req.Method == enums.ProductAllocateMethodStore && req.ToStoreId == "" {
		return errors.New("调拨门店不能为空")
	}

	return nil
}

type ProductAllocateWhere struct {
	Method      enums.ProductAllocateMethod `json:"method" label:"调拨方式" input:"select" type:"number" find:"true" sort:"1" required:"true" preset:"typeMap"` // 调拨方式
	Type        enums.ProductType           `json:"type" label:"仓库类型" input:"select" type:"number" find:"true" sort:"2" required:"true" preset:"typeMap"`   // 仓库类型
	Reason      enums.ProductAllocateReason `json:"reason" label:"调拨原因" input:"select" type:"number" find:"true" sort:"3" required:"true" preset:"typeMap"` // 调拨原因
	FromStoreId string                      `json:"from_store_id" label:"调出门店" input:"search" type:"string" find:"true" sort:"4" required:"false"`          // 调出门店
	ToStoreId   string                      `json:"to_store_id" label:"调入门店" input:"search" type:"string" find:"true" sort:"4" required:"false"`            // 调入门店
	Status      enums.ProductAllocateStatus `json:"status" label:"调拨状态" input:"select" type:"number" find:"true" sort:"5" required:"true" preset:"typeMap"` // 调拨状态

	StartTime *time.Time `json:"start_time" label:"开始时间" input:"date" type:"date" find:"true" sort:"6" required:"false"` // 开始时间
	EndTime   *time.Time `json:"end_time" label:"结束时间" input:"date" type:"date" find:"true" sort:"6" required:"false"`   // 结束时间
}

func (req *ProductAllocateWhere) Validate() error {
	if req.Method == enums.ProductAllocateMethodStore && req.ToStoreId == "" {
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

type ProductAllocateRemoveReq struct {
	Id   string `json:"id" binding:"required"`   // 调拨单ID
	Code string `json:"code" binding:"required"` // 产品条码
}

type ProductAllocateConfirmReq struct {
	Id string `json:"id" binding:"required"` // 调拨单ID
}

type ProductAllocateCancelReq struct {
	Id string `json:"id" binding:"required"` // 调拨单ID
}

type ProductAllocateCompleteReq struct {
	Id string `json:"id" binding:"required"` // 调拨单ID
}

type ProductEnterWhere struct {
	Id        string     `json:"id" label:"ID" input:"text" type:"string" find:"true" sort:"1" required:"false"`         // ID
	StartTime *time.Time `json:"start_time" label:"开始时间" input:"date" type:"date" find:"true" sort:"2" required:"false"` // 开始时间
	EndTime   *time.Time `json:"end_time" label:"结束时间" input:"date" type:"date" find:"true" sort:"3" required:"false"`   // 结束时间
}

type ProductEnterListReq struct {
	PageReq
	Where ProductEnterWhere `json:"where"`
}

type ProductEnterInfoReq struct {
	Id string `json:"id" binding:"required"`
}

type ProductInventoryWhere struct {
	Id      string `json:"id" label:"ID" input:"text" type:"string" find:"true" sort:"1" required:"false"`                                           // ID
	StoreId string `json:"store_id" label:"门店" input:"search" type:"string" find:"false" create:"false" sort:"2" required:"true" binding:"required"` // 门店

	InventoryPersonId string `json:"inventory_person_id" label:"盘点人" input:"search" type:"string" find:"true" create:"true" sort:"3" required:"true"` // 盘点人
	InspectorId       string `json:"inspector_id" label:"审核人" input:"search" type:"string" find:"true" create:"true" sort:"4" required:"true"`        // 监盘人

	Type   enums.ProductType            `json:"type" label:"仓库类型" input:"select" type:"number" find:"true" create:"true" sort:"5" required:"true" preset:"typeMap"` // 仓库类型
	Range  enums.ProductInventoryRange  `json:"range" label:"范围" input:"select" type:"number" find:"true" create:"true" sort:"6" required:"true" preset:"typeMap"`  // 范围
	Status enums.ProductInventoryStatus `json:"status" label:"状态" input:"select" type:"number" find:"false" sort:"7" required:"false" preset:"typeMap"`             // 状态

	Brand    enums.ProductBrand    `json:"brand" label:"品牌" input:"multiple" type:"number" find:"false" create:"true" sort:"8" required:"false" preset:"typeMap"`     // 品牌
	Class    enums.ProductClass    `json:"class" label:"系列" input:"multiple" type:"number" find:"false" create:"true" sort:"9" required:"false" preset:"typeMap"`     // 系列
	Category enums.ProductCategory `json:"category" label:"类别" input:"multiple" type:"number" find:"false" create:"true" sort:"10" required:"false" preset:"typeMap"` // 类别
	Craft    enums.ProductCraft    `json:"craft" label:"工艺" input:"multiple" type:"number" find:"false" create:"true" sort:"11" required:"false" preset:"typeMap"`    // 工艺
	Material enums.ProductMaterial `json:"material" label:"材质" input:"multiple" type:"number" find:"false" create:"true" sort:"12" required:"false" preset:"typeMap"` // 材质
	Quality  enums.ProductQuality  `json:"quality" label:"质地" input:"multiple" type:"number" find:"false" create:"true" sort:"13" required:"false" preset:"typeMap"`  // 质地
	Gem      enums.ProductGem      `json:"gem" label:"宝石" input:"multiple" type:"number" find:"false" create:"true" sort:"14" required:"false" preset:"typeMap"`      // 宝石

	Remark string `json:"remark" label:"备注" input:"textarea" type:"string" find:"false" create:"true" sort:"15" required:"false"` // 备注

	StartTime *time.Time `json:"start_time" label:"开始时间" input:"date" type:"date" find:"true" sort:"16" required:"false"` // 开始时间
	EndTime   *time.Time `json:"end_time" label:"结束时间" input:"date" type:"date" find:"true" sort:"17" required:"false"`   // 结束时间

	ProductStatus enums.ProductInventoryProductStatus `json:"product_status" label:"状态" input:"select" type:"number" find:"false" create:"false" sort:"18" required:"false" preset:"typeMap"` // 产品状态
}

type ProductInventoryCreateReq struct {
	StoreId string `json:"store_id" binding:"required"` // 门店ID

	InventoryPersonId string `json:"inventory_person_id" binding:"required"` // 盘点人
	InspectorId       string `json:"inspector_id" binding:"required"`        // 监盘人

	Type  enums.ProductType           `json:"type" binding:"required"`  // 仓库类型
	Range enums.ProductInventoryRange `json:"range" binding:"required"` // 盘点范围

	Brand    []enums.ProductBrand    `json:"brand"`    // 品牌
	Class    []enums.ProductClass    `json:"class"`    // 系列
	Category []enums.ProductCategory `json:"category"` // 类别
	Craft    []enums.ProductCraft    `json:"craft"`    // 工艺
	Material []enums.ProductMaterial `json:"material"` // 材质
	Quality  []enums.ProductQuality  `json:"quality"`  // 质地
	Gem      []enums.ProductGem      `json:"gem"`      // 宝石

	Remark string `json:"remark"`
}

func (req *ProductInventoryCreateReq) Validate() error {
	if err := req.Type.InMap(); err != nil {
		return err
	}
	if err := req.Range.InMap(); err != nil {
		return err
	}

	return nil
}

type ProductInventoryListReq struct {
	PageReq
	Where ProductInventoryWhere `json:"where"`
}

type ProductInventoryInfoReq struct {
	Id            string                              `json:"id" binding:"required"`
	ProductStatus enums.ProductInventoryProductStatus `json:"product_status"` // 产品状态
}

type ProductInventoryChangeReq struct {
	Id string `json:"id" binding:"required"`

	Status enums.ProductInventoryStatus `json:"status" binding:"required"`
}
