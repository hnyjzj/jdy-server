package product

import (
	"jdy/errors"
	"jdy/logic/product"
	"jdy/types"

	"github.com/gin-gonic/gin"
)

// 入库单
func (con ProductController) Enter(ctx *gin.Context) {
	var (
		req types.ProductEnterReq

		logic = product.ProductLogic{
			Ctx:   ctx,
			Staff: con.GetStaff(ctx),
		}
	)

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	// 调用逻辑层
	res, err := logic.Enter(&req)
	if err != nil {
		con.ErrorLogic(ctx, err)
		return
	}

	con.Success(ctx, "ok", res)
}
