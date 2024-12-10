package product

import (
	"jdy/controller"
	"jdy/logic/product"
	"jdy/types"
	"jdy/utils"

	"github.com/gin-gonic/gin"
)

type ProductController struct {
	controller.BaseController
}

// 产品筛选条件
func (con ProductController) Where(ctx *gin.Context) {
	where := utils.ModelToWhere(types.ProductWhere{}, map[string]any{
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

// 产品列表
func (con ProductController) List(ctx *gin.Context) {
	var (
		req types.ProductListReq

		logic = product.ProductLogic{
			Ctx:   ctx,
			Staff: con.GetStaff(ctx),
		}
	)

	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, "参数错误")
		return
	}

	res, err := logic.List(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", res)
}
