package product

import (
	"jdy/controller"
	"jdy/errors"
	"jdy/logic/product"
	"jdy/types"
	"jdy/utils"

	"github.com/gin-gonic/gin"
)

type ProductAccessorieController struct {
	controller.BaseController
}

// 配件筛选条件
func (con ProductAccessorieController) Where(ctx *gin.Context) {
	where := utils.StructToWhere(types.ProductAccessorieWhere{})

	con.Success(ctx, "ok", where)
}

// 配件列表
func (con ProductAccessorieController) List(ctx *gin.Context) {
	var (
		req types.ProductAccessorieListReq

		logic = product.ProductAccessorieLogic{
			Ctx: ctx,
		}
	)

	if staff, err := con.GetStaff(ctx); err != nil {
		con.ExceptionWithAuth(ctx, err)
		return
	} else {
		logic.Staff = staff
	}

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	res, err := logic.List(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", res)
}

// 配件详情
func (con ProductAccessorieController) Info(ctx *gin.Context) {
	var (
		req types.ProductAccessorieInfoReq

		logic = product.ProductAccessorieLogic{
			Ctx: ctx,
		}
	)

	if staff, err := con.GetStaff(ctx); err != nil {
		con.ExceptionWithAuth(ctx, err)
		return
	} else {
		logic.Staff = staff
	}

	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	res, err := logic.Info(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", res)
}

// 更新商品信息
func (con ProductAccessorieController) Update(ctx *gin.Context) {
	var (
		req types.ProductAccessorieUpdateReq

		logic = product.ProductAccessorieLogic{
			Ctx: ctx,
		}
	)

	if staff, err := con.GetStaff(ctx); err != nil {
		con.ExceptionWithAuth(ctx, err)
		return
	} else {
		logic.Staff = staff
	}

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	// 调用逻辑层
	if err := logic.Update(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}
