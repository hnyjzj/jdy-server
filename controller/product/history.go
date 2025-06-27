package product

import (
	"jdy/controller"
	"jdy/errors"
	"jdy/logic/product"
	"jdy/types"
	"jdy/utils"

	"github.com/gin-gonic/gin"
)

type ProductHistoryController struct {
	controller.BaseController
}

// 产品操作记录条件
func (con ProductHistoryController) Where(ctx *gin.Context) {
	where := utils.StructToWhere(types.ProductHistoryWhere{})

	con.Success(ctx, "ok", where)
}

// 产品操作记录列表
func (con ProductHistoryController) List(ctx *gin.Context) {
	var (
		req types.ProductHistoryListReq

		logic = product.ProductHistoryLogic{
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

	res, err := logic.List(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", res)
}

// 产品操作记录详情
func (con ProductHistoryController) Info(ctx *gin.Context) {
	var (
		req types.ProductHistoryInfoReq

		logic = product.ProductHistoryLogic{
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
