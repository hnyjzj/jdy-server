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

// 更新商品信息
func (con ProductFinishedBatchController) Code(ctx *gin.Context) {
	var (
		req types.ProductFinishedUpdateCodeReq

		logic = product.ProductFinishedLogic{
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
	if err := logic.Code(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}
