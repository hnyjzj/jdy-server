package types

import (
	"errors"
	"jdy/enums"

	"github.com/shopspring/decimal"
)

type ProductOldWhere struct {
	Code         string                `json:"code" label:"旧料条码" find:"true" create:"false" update:"false" sort:"4" type:"string" input:"text" required:"false"`                       // 条码
	CodeFinished string                `json:"code_finished" label:"成品条码" find:"true" create:"true" update:"false" sort:"4" type:"string" input:"text" required:"false"`               // 成品条码
	Name         string                `json:"name" label:"名称" find:"true" create:"true" update:"false" sort:"5" type:"string" input:"text" required:"false"`                          // 名称
	Class        enums.ProductClassOld `json:"class" label:"大类" find:"true" create:"false" update:"false" sort:"3" type:"number" input:"select" required:"false" preset:"typeMap"`     // 大类
	Status       enums.ProductStatus   `json:"status" label:"状态" find:"true" create:"false" update:"false" sort:"13" type:"number" input:"select" required:"false" preset:"typeMap"`   // 状态
	LabelPrice   *decimal.Decimal      `json:"label_price" label:"标签价" find:"true" create:"true" update:"false" sort:"12" type:"number" input:"text" required:"false"`                 // 标签价
	Brand        enums.ProductBrand    `json:"brand" label:"品牌" find:"true" create:"true" update:"false" sort:"16" type:"number" input:"select" required:"false" preset:"typeMap"`     // 品牌
	Material     enums.ProductMaterial `json:"material" label:"材质" find:"true" create:"true" update:"false" sort:"6" type:"number" input:"select" required:"true" preset:"typeMap"`    // 材质
	Quality      enums.ProductQuality  `json:"quality" label:"成色" find:"true" create:"true" update:"false" sort:"7" type:"number" input:"select" required:"true" preset:"typeMap"`     // 成色
	Gem          enums.ProductGem      `json:"gem" label:"主石" find:"true" create:"true" update:"false" sort:"8" type:"number" input:"select" required:"true" preset:"typeMap"`         // 主石
	Category     enums.ProductCategory `json:"category" label:"品类" find:"true" create:"true" update:"false" sort:"9" type:"number" input:"select" required:"false" preset:"typeMap"`   // 品类
	Craft        enums.ProductCraft    `json:"craft" label:"工艺" find:"true" create:"true" update:"false" sort:"10" type:"number" input:"select" required:"false" preset:"typeMap"`     // 工艺
	WeightMetal  *decimal.Decimal      `json:"weight_metal" label:"金重" find:"true" create:"true" update:"false" sort:"11" type:"number" input:"text" required:"true"`                  // 金重
	WeightTotal  *decimal.Decimal      `json:"weight_total" label:"总重" find:"true" create:"true" update:"false" sort:"24" type:"number" input:"text" required:"false"`                 // 总重
	ColorGem     enums.ProductColor    `json:"color_gem" label:"颜色" find:"true" create:"true" update:"false" sort:"18" type:"number" input:"select" required:"false" preset:"typeMap"` // 颜色
	WeightGem    *decimal.Decimal      `json:"weight_gem" label:"主石重" find:"true" create:"true" update:"false" sort:"17" type:"number" input:"text" required:"false"`                  // 主石重
	NumGem       int                   `json:"num_gem" label:"主石数" find:"true" create:"true" update:"false" sort:"21" type:"number" input:"number" required:"false"`                   // 主石数
	Clarity      enums.ProductClarity  `json:"clarity" label:"净度" find:"true" create:"true" update:"false" sort:"19" type:"number" input:"select" required:"false" preset:"typeMap"`   // 净度
	Cut          enums.ProductCut      `json:"cut" label:"切工" find:"true" create:"true" update:"false" sort:"20" type:"number" input:"select" required:"false" preset:"typeMap"`       // 切工
	WeightOther  *decimal.Decimal      `json:"weight_other" label:"副石重" find:"true" create:"true" update:"false" sort:"22" type:"number" input:"text" required:"false"`                // 副石重
	NumOther     int                   `json:"num_other" label:"副石数" find:"true" create:"true" update:"false" sort:"23" type:"number" input:"number" required:"false"`                 // 副石数
	Remark       string                `json:"remark" label:"备注" find:"true" create:"true" update:"false" sort:"26" type:"string" input:"textarea" required:"false"`                   // 备注

	StoreId string `json:"store_id" label:"所属店铺" find:"false" create:"true" update:"false" sort:"28" type:"string" input:"text" required:"false"`

	IsOur                   *bool                      `json:"is_our" label:"是否自有" find:"true" create:"true" update:"false" sort:"1" type:"boolean" input:"switch" required:"true"`                                         // 是否自有
	RecycleMethod           enums.ProductRecycleMethod `json:"recycle_method" label:"回收方式" find:"true" create:"true" update:"false" sort:"2" type:"number" input:"select" required:"true" preset:"typeMap"`                 // 回收方式
	RecycleType             enums.ProductRecycleType   `json:"recycle_type" label:"回收类型" find:"true" create:"true" update:"false" sort:"3" type:"number" input:"select" required:"true" preset:"typeMap"`                   // 回收类型
	RecyclePrice            *decimal.Decimal           `json:"recycle_price" label:"回收金额" find:"true" create:"true" update:"false" sort:"27" type:"number" input:"text" required:"false"`                                   // 回收金额
	RecyclePriceGold        *decimal.Decimal           `json:"recycle_price_gold" label:"回收金价" find:"true" create:"true" update:"false" sort:"13" type:"number" input:"text" required:"false"`                              // 回收金价
	RecyclePriceLabor       *decimal.Decimal           `json:"recycle_price_labor" label:"回收工费" find:"true" create:"true" update:"false" sort:"15" type:"number" input:"text" required:"false"`                             // 回收工费
	RecyclePriceLaborMethod enums.ProductRecycleMethod `json:"recycle_price_labor_method" label:"回收工费方式" find:"true" create:"true" update:"false" sort:"14" type:"number" input:"select" required:"false" preset:"typeMap"` // 回收工费方式
	QualityActual           *decimal.Decimal           `json:"quality_actual" label:"实际成色" find:"true" create:"true" update:"false" sort:"25" type:"number" input:"text" required:"true"`                                   // 实际成色
	RecycleSource           enums.ProductRecycleSource `json:"recycle_source" label:"回收来源" find:"true" create:"true" update:"false" sort:"4" type:"number" input:"select" required:"true" preset:"typeMap"`                 // 回收来源
	RecycleSourceId         string                     `json:"recycle_source_id" label:"回收来源" find:"true" create:"false" update:"false" sort:"5" type:"string" input:"text" required:"false"`                               // 回收来源
	RecycleStoreId          string                     `json:"recycle_store_id" label:"回收店铺" find:"true" create:"false" update:"false" sort:"6" type:"string" input:"text" required:"false"`                                // 回收店铺
}

type ProductOldCreateWhere struct {
	IsOur                   *bool                      `json:"is_our" label:"是否为本公司货品" info:"true" find:"true" create:"true" update:"false" sort:"1" type:"boolean" input:"switch" required:"true"`                         // 是否为本公司货品
	Code                    string                     `json:"code" label:"旧料条码" find:"true" create:"false" update:"false" sort:"2" type:"string" input:"text" required:"false"`                                            // 条码
	CodeFinished            string                     `json:"code_finished" label:"成品条码" find:"true" create:"true" update:"false" sort:"3" type:"string" input:"text" required:"false"`                                    // 成品条码
	RecycleMethod           enums.ProductRecycleMethod `json:"recycle_method" label:"回收方式" find:"true" create:"true" update:"false" sort:"3" type:"number" input:"select" required:"true" preset:"typeMap"`                 // 回收方式
	RecycleType             enums.ProductRecycleType   `json:"recycle_type" label:"回收类型" find:"true" create:"true" update:"false" sort:"4" type:"number" input:"select" required:"true" preset:"typeMap"`                   // 回收类型
	Material                enums.ProductMaterial      `json:"material" label:"材质" find:"true" create:"true" update:"false" sort:"5" type:"number" input:"select" required:"true" preset:"typeMap"`                         // 材质
	Gem                     enums.ProductGem           `json:"gem" label:"主石" find:"true" create:"true" update:"false" sort:"6" type:"number" input:"select" required:"true" preset:"typeMap"`                              // 主石
	Quality                 enums.ProductQuality       `json:"quality" label:"成色" find:"true" create:"true" update:"false" sort:"7" type:"number" input:"select" required:"true" preset:"typeMap"`                          // 成色
	QualityActual           *decimal.Decimal           `json:"quality_actual" label:"实际成色" find:"true" create:"true" update:"false" sort:"8" type:"number" input:"text" required:"true"`                                    // 实际成色
	WeightMetal             *decimal.Decimal           `json:"weight_metal" label:"金重" find:"true" create:"true" update:"false" sort:"9" type:"number" input:"text" required:"true"`                                        // 金重
	RecyclePriceGold        *decimal.Decimal           `json:"recycle_price_gold" label:"回收金价" find:"true" create:"true" update:"false" sort:"10" type:"number" input:"text" required:"false"`                              // 回收金价
	RecyclePriceLabor       *decimal.Decimal           `json:"recycle_price_labor" label:"回收工费" find:"true" create:"true" update:"false" sort:"11" type:"number" input:"text" required:"false"`                             // 回收工费
	RecyclePriceLaborMethod enums.ProductRecycleMethod `json:"recycle_price_labor_method" label:"回收工费方式" find:"true" create:"true" update:"false" sort:"12" type:"number" input:"select" required:"false" preset:"typeMap"` // 回收工费方式
	LabelPrice              *decimal.Decimal           `json:"label_price" label:"标签价" find:"true" create:"true" update:"false" sort:"13" type:"number" input:"text" required:"false"`                                      // 标签价
	Category                enums.ProductCategory      `json:"category" label:"品类" find:"true" create:"true" update:"false" sort:"14" type:"number" input:"select" required:"false" preset:"typeMap"`                       // 品类
	Craft                   enums.ProductCraft         `json:"craft" label:"工艺" find:"true" create:"true" update:"false" sort:"15" type:"number" input:"select" required:"false" preset:"typeMap"`                          // 工艺
	Brand                   enums.ProductBrand         `json:"brand" label:"品牌" find:"true" create:"true" update:"false" sort:"16" type:"number" input:"select" required:"false" preset:"typeMap"`                          // 品牌
	WeightGem               *decimal.Decimal           `json:"weight_gem" label:"主石重" find:"true" create:"true" update:"false" sort:"17" type:"number" input:"text" required:"false"`                                       // 主石重
	ColorGem                enums.ProductColor         `json:"color_gem" label:"主石颜色" find:"true" create:"true" update:"false" sort:"18" type:"number" input:"select" required:"false" preset:"typeMap"`                    // 主石颜色
	Clarity                 enums.ProductClarity       `json:"clarity" label:"主石净度" find:"true" create:"true" update:"false" sort:"19" type:"number" input:"select" required:"false" preset:"typeMap"`                      // 主石净度
	Cut                     enums.ProductCut           `json:"cut" label:"主石切工" find:"true" create:"true" update:"false" sort:"20" type:"number" input:"select" required:"false" preset:"typeMap"`                          // 主石切工
	NumGem                  int                        `json:"num_gem" label:"主石数量" find:"true" create:"true" update:"false" sort:"21" type:"number" input:"number" required:"false"`                                       // 主石数量
	WeightOther             *decimal.Decimal           `json:"weight_other" label:"副石重" find:"true" create:"true" update:"false" sort:"22" type:"number" input:"text" required:"false"`                                     // 副石重
	NumOther                int                        `json:"num_other" label:"副石数" find:"true" create:"true" update:"false" sort:"23" type:"number" input:"number" required:"false"`                                      // 副石数
	WeightTotal             *decimal.Decimal           `json:"weight_total" label:"总重" find:"true" create:"true" update:"false" sort:"24" type:"number" input:"text" required:"false"`                                      // 总重
	Remark                  string                     `json:"remark" label:"备注" find:"true" create:"true" update:"false" sort:"25" type:"string" input:"textarea" required:"false"`                                        // 备注
	Name                    string                     `json:"name" label:"货品名称" find:"true" create:"true" update:"false" sort:"26" type:"string" input:"text" required:"false"`                                            // 货品名称
	RecyclePrice            *decimal.Decimal           `json:"recycle_price" label:"回收金额" find:"true" create:"true" update:"false" sort:"27" type:"number" input:"text" required:"false"`                                   // 回收金额
}

type ProductOldListReq struct {
	PageReq
	Where ProductOldWhere `json:"where" binding:"required"`
}

type ProductOldInfoReq struct {
	Id   string `json:"id"`   // ID
	Code string `json:"code"` // 条码
}

func (r *ProductOldInfoReq) Validate() error {
	if r.Id == "" && r.Code == "" {
		return errors.New("条件不能为空")
	}

	return nil
}

type ProductOldUpdateReq struct {
	Id string `json:"id" binding:"required"` // ID
	ProductOldWhere
}

type ProductOldGetClassReq struct {
	Material enums.ProductMaterial `json:"material" binding:"required"` // 材质
	Quality  enums.ProductQuality  `json:"quality" binding:"required"`  // 成色
	Gem      enums.ProductGem      `json:"gem" binding:"required"`      // 主石
}

type ProductOldGetClassRes struct {
	Value enums.ProductClassOld `json:"value"` // 大类
	Label string                `json:"label"` // 大类名称
}
