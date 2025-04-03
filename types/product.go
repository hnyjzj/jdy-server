package types

import (
	"errors"
	"jdy/enums"
	"time"
)

type ProductDamageReq struct {
	Code   string `json:"code" binding:"required"`   // 条码
	Reason string `json:"reason" binding:"required"` // 损坏原因
}

type ProductConversionReq struct {
	Id     string            `json:"id" binding:"required"`   // 产品ID
	Type   enums.ProductType `json:"type" binding:"required"` // 仓库类型
	Remark string            `json:"remark"`                  // 备注
}

type ProductHistoryWhere struct {
	ProductId string              `json:"product_id" label:"产品" input:"text" type:"string" find:"true" sort:"1" required:"true"`                   // 产品
	StoreId   string              `json:"store_id" label:"门店" input:"text" type:"string" find:"false" sort:"2" required:"true" binding:"required"` // 门店
	Action    enums.ProductAction `json:"action" label:"操作" input:"select" type:"number" find:"true" sort:"3" required:"true" preset:"typeMap"`    // 操作
}

type ProductHistoryListReq struct {
	PageReq
	Where ProductHistoryWhere `json:"where" binding:"required"`
}

type ProductHistoryInfoReq struct {
	Id string `json:"id" binding:"required"`
}

type ProductInventoryWhere struct {
	Id      string `json:"id" label:"ID" input:"text" type:"string" find:"true" sort:"1" required:"false"`                                           // ID
	StoreId string `json:"store_id" label:"门店" input:"search" type:"string" find:"false" create:"false" sort:"2" required:"true" binding:"required"` // 门店

	InventoryPersonId string `json:"inventory_person_id" label:"盘点人" input:"search" type:"string" find:"true" create:"true" sort:"3" required:"true"` // 盘点人
	InspectorId       string `json:"inspector_id" label:"监盘人" input:"search" type:"string" find:"true" create:"true" sort:"4" required:"true"`        // 监盘人

	Type  enums.ProductType           `json:"type" label:"盘点仓库" input:"select" type:"number" find:"true" create:"true" sort:"5" required:"true" preset:"typeMap"`      // 盘点仓库
	Brand enums.ProductBrand          `json:"brand" label:"盘点品牌" input:"multiple" type:"number" find:"false" create:"true" sort:"6" required:"false" preset:"typeMap"` // 盘点品牌
	Range enums.ProductInventoryRange `json:"range" label:"盘点范围" input:"select" type:"number" find:"true" create:"true" sort:"7" required:"true" preset:"typeMap"`     // 盘点范围

	Class    enums.ProductClass    `json:"class" label:"大类" input:"multiple" type:"number" find:"false" create:"true" sort:"8" required:"false" preset:"typeMap" condition:"[{\"key\":\"range\",\"operator\":\"=\",\"value\":1}]"`     // 大类
	Category enums.ProductCategory `json:"category" label:"品类" input:"multiple" type:"number" find:"false" create:"true" sort:"9" required:"false" preset:"typeMap"`                                                                   // 品类
	Craft    enums.ProductCraft    `json:"craft" label:"工艺" input:"multiple" type:"number" find:"false" create:"true" sort:"10" required:"false" preset:"typeMap"`                                                                     // 工艺
	Material enums.ProductMaterial `json:"material" label:"材质" input:"multiple" type:"number" find:"false" create:"true" sort:"11" required:"false" preset:"typeMap" condition:"[{\"key\":\"range\",\"operator\":\"=\",\"value\":2}]"` // 材质
	Quality  enums.ProductQuality  `json:"quality" label:"成色" input:"multiple" type:"number" find:"false" create:"true" sort:"12" required:"false" preset:"typeMap" condition:"[{\"key\":\"range\",\"operator\":\"=\",\"value\":2}]"`  // 成色
	Gem      enums.ProductGem      `json:"gem" label:"主石" input:"multiple" type:"number" find:"false" create:"true" sort:"13" required:"false" preset:"typeMap" condition:"[{\"key\":\"range\",\"operator\":\"=\",\"value\":2}]"`      // 主石

	Remark string `json:"remark" label:"备注" input:"textarea" type:"string" find:"false" create:"true" sort:"14" required:"false"` // 备注

	Status enums.ProductInventoryStatus `json:"status" label:"状态" input:"select" type:"number" find:"false" sort:"15" required:"false" preset:"typeMap"` // 状态

	StartTime *time.Time `json:"start_time" label:"开始时间" input:"date" type:"date" find:"true" sort:"16" required:"false"` // 开始时间
	EndTime   *time.Time `json:"end_time" label:"结束时间" input:"date" type:"date" find:"true" sort:"17" required:"false"`   // 结束时间

	ProductStatus enums.ProductInventoryProductStatus `json:"product_status" label:"状态" input:"select" type:"number" find:"false" create:"false" sort:"18" required:"false" preset:"typeMap"` // 产品状态
}

type ProductInventoryCreateReq struct {
	StoreId string `json:"store_id" binding:"required"` // 门店ID

	InventoryPersonId string `json:"inventory_person_id" binding:"required"` // 盘点人
	InspectorId       string `json:"inspector_id" binding:"required"`        // 监盘人

	Type  enums.ProductType           `json:"type" binding:"required"`  // 盘点仓库
	Brand []enums.ProductBrand        `json:"brand"`                    // 盘点品牌
	Range enums.ProductInventoryRange `json:"range" binding:"required"` // 盘点范围

	Category []enums.ProductCategory `json:"category" binding:"required"` // 品类
	Craft    []enums.ProductCraft    `json:"craft" binding:"required"`    // 工艺
	Class    []enums.ProductClass    `json:"class"`                       // 大类
	Material []enums.ProductMaterial `json:"material"`                    // 材质
	Quality  []enums.ProductQuality  `json:"quality"`                     // 成色
	Gem      []enums.ProductGem      `json:"gem"`                         // 宝石

	Remark string `json:"remark"`
}

func (req *ProductInventoryCreateReq) Validate() error {
	if err := req.Type.InMap(); err != nil {
		return err
	}
	if err := req.Range.InMap(); err != nil {
		return err
	}

	switch req.Range {
	case enums.ProductInventoryRangeBigType:
		if len(req.Class) == 0 {
			return errors.New("大类是必填项")
		}
	case enums.ProductInventoryRangeMaterialType:
		if len(req.Material) == 0 {
			return errors.New("材质是必填项")
		}
		if len(req.Quality) == 0 {
			return errors.New("成色是必填项")
		}
		if len(req.Gem) == 0 {
			return errors.New("主石是必填项")
		}
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
