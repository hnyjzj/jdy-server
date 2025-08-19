package product

import (
	"jdy/errors"
	"jdy/logic/product"
	"jdy/types"

	"github.com/gin-gonic/gin"
)

type ProductFinishedBatchController struct {
	ProductFinishedController
}

// 批量更新
func (con ProductFinishedBatchController) Update(ctx *gin.Context) {
	var (
		req types.ProductFinishedUpdatesReq

		logic = product.ProductFinishedBatchLogic{
			ProductFinishedLogic: product.ProductFinishedLogic{
				Ctx: ctx,
			},
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

	// 校验参数
	if err := req.Validate(); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	// 调用逻辑层
	if err := logic.Update(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}

// 批量更新条码
func (con ProductFinishedBatchController) UpdateCode(ctx *gin.Context) {
	var (
		req types.ProductFinishedUpdateCodeReq

		logic = product.ProductFinishedBatchLogic{
			ProductFinishedLogic: product.ProductFinishedLogic{
				Ctx: ctx,
			},
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
	if err := logic.UpdateCode(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}

// 批量查找条码
func (con ProductFinishedBatchController) FindCode(ctx *gin.Context) {
	var (
		req types.ProductFinishedFindCodeReq

		logic = product.ProductFinishedBatchLogic{
			ProductFinishedLogic: product.ProductFinishedLogic{
				Ctx: ctx,
			},
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
	res, err := logic.FindCode(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", res)
}
