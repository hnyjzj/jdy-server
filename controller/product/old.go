package product

import (
	"jdy/controller"
	"jdy/errors"
	"jdy/logic/product"
	"jdy/types"
	"jdy/utils"

	"github.com/gin-gonic/gin"
)

type ProductOldController struct {
	controller.BaseController
}

// 产品筛选条件
func (con ProductOldController) Where(ctx *gin.Context) {
	where := utils.StructToWhere(types.ProductOldWhere{})

	con.Success(ctx, "ok", where)
}

// 产品列表
func (con ProductOldController) List(ctx *gin.Context) {
	var (
		req types.ProductOldListReq

		logic = product.ProductOldLogic{
			Ctx: ctx,
		}
	)

	staff, err := con.GetStaff(ctx)
	if err != nil {
		con.ExceptionWithAuth(ctx, err.Error())
		return
	}
	logic.Staff = staff

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

// 产品详情
func (con ProductOldController) Info(ctx *gin.Context) {
	var (
		req types.ProductOldInfoReq

		logic = product.ProductOldLogic{
			Ctx: ctx,
		}
	)

	staff, err := con.GetStaff(ctx)
	if err != nil {
		con.ExceptionWithAuth(ctx, err.Error())
		return
	}
	logic.Staff = staff

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

// 产品转换
func (con ProductOldController) Conversion(ctx *gin.Context) {
	var (
		req types.ProductConversionReq

		logic = product.ProductOldLogic{
			Ctx: ctx,
		}
	)

	staff, err := con.GetStaff(ctx)
	if err != nil {
		con.ExceptionWithAuth(ctx, err.Error())
		return
	}
	logic.Staff = staff

	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	if err := logic.Conversion(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}

// 获取大类
func (con ProductOldController) GetClass(ctx *gin.Context) {
	var (
		req types.ProductOldGetClassReq

		logic = product.ProductOldLogic{
			Ctx: ctx,
		}
	)

	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	class := logic.GetClass(&req)

	con.Success(ctx, "ok", class)
}
