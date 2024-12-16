package product

import (
	"jdy/errors"
	"jdy/logic/product"
	"jdy/types"

	"github.com/gin-gonic/gin"
)

// 产品报损
func (con ProductController) Damage(ctx *gin.Context) {
	var (
		req types.ProductDamageReq

		logic = product.ProductLogic{
			Ctx:   ctx,
			Staff: con.GetStaff(ctx),
		}
	)

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

// 产品调拨
func (con ProductController) Allocate(ctx *gin.Context) {
	var (
		req types.ProductAllocateReq

		logic = product.ProductLogic{
			Ctx:   ctx,
			Staff: con.GetStaff(ctx),
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

	if err := logic.Allocate(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}
