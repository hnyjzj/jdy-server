package product

import (
	"jdy/controller"
	"jdy/errors"
	"jdy/logic/product"
	"jdy/types"
	"jdy/utils"

	"github.com/gin-gonic/gin"
)

type ProductAllocateController struct {
	controller.BaseController
}

// 创建产品调拨单
func (con ProductAllocateController) Create(ctx *gin.Context) {
	var (
		req types.ProductAllocateCreateReq

		logic = product.ProductAllocateLogic{
			ProductLogic: product.ProductLogic{
				Ctx:   ctx,
				Staff: con.GetStaff(ctx),
			},
		}
	)

	// 绑定请求参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	// 校验参数
	if err := req.Validate(); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	if err := logic.Create(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}

// 产品调拨单筛选条件
func (con ProductAllocateController) Where(ctx *gin.Context) {
	where := utils.StructToWhere(types.ProductAllocateWhere{})

	con.Success(ctx, "ok", where)
}

// 获取产品调拨单列表
func (con ProductAllocateController) List(ctx *gin.Context) {
	var (
		req types.ProductAllocateListReq

		logic = product.ProductAllocateLogic{
			ProductLogic: product.ProductLogic{
				Ctx:   ctx,
				Staff: con.GetStaff(ctx),
			},
		}
	)

	// 绑定请求参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	// 校验参数
	if err := req.Where.Validate(); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	// 获取产品调拨单列表
	res, err := logic.List(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", res)

}
