package product

import (
	"jdy/controller"
	"jdy/errors"
	"jdy/logic/product"
	"jdy/types"
	"jdy/utils"

	"github.com/gin-gonic/gin"
)

type ProductFinishedDamageController struct {
	controller.BaseController
}

// 产品筛选条件
func (con ProductFinishedDamageController) Where(ctx *gin.Context) {
	where := utils.StructToWhere(types.ProductFinishedWhere{})

	con.Success(ctx, "ok", where)
}

// 成品报损
func (con ProductFinishedDamageController) Damage(ctx *gin.Context) {
	var (
		req types.ProductDamageReq

		logic = product.ProductFinishedDamageLogic{
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

	if err := logic.Damage(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}

// 产品列表
func (con ProductFinishedDamageController) List(ctx *gin.Context) {
	var (
		req types.ProductFinishedListReq

		logic = product.ProductFinishedDamageLogic{
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

// 产品详情
func (con ProductFinishedDamageController) Info(ctx *gin.Context) {
	var (
		req types.ProductFinishedInfoReq

		logic = product.ProductFinishedDamageLogic{
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

// 产品转换
func (con ProductFinishedDamageController) Conversion(ctx *gin.Context) {
	var (
		req types.ProductConversionReq

		logic = product.ProductFinishedDamageLogic{
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

	if err := logic.Conversion(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}
