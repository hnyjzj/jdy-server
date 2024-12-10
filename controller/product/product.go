package product

import (
	"jdy/controller"
	"jdy/model"
	"jdy/types"
	"jdy/utils"

	"github.com/gin-gonic/gin"
)

type ProductController struct {
	controller.BaseController
}

// 产品筛选条件
func (con ProductController) Where(ctx *gin.Context) {
	where := utils.ModelToWhere(model.Product{}, map[string]any{
		"ColorMetal":     types.ProductColorMap,          // 金颜色
		"ColorGem":       types.ProductColorMap,          // 主石色
		"Clarity":        types.ProductClarityMap,        // 净度
		"RetailType":     types.ProductRetailTypeMap,     // 零售方式
		"Class":          types.ProductClassMap,          // 大类
		"Supplier":       types.ProductSupplierMap,       // 供应商
		"Material":       types.ProductMaterialMap,       // 材质
		"Quality":        types.ProductQualityMap,        // 成色
		"Gem":            types.ProductGemMap,            // 宝石
		"Category":       types.ProductCategoryMap,       // 品类
		"Brand":          types.ProductBrandMap,          // 品牌
		"Craft":          types.ProductCraftMap,          // 工艺
		"IsSpecialOffer": map[int]string{1: "是", 0: "否"}, // 是否特价
		"Status":         types.ProductStatusMap,         // 状态
	})

	con.Success(ctx, "ok", where)
}
