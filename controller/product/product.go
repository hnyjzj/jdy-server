package product

import (
	"jdy/controller"
	"jdy/errors"
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
	where := utils.StructToWhere(types.ProductWhere{})

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
func (con ProductController) Info(ctx *gin.Context) {
	var (
		req types.ProductInfoReq

		logic = product.ProductLogic{
			Ctx:   ctx,
			Staff: con.GetStaff(ctx),
		}
	)

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
func (con ProductController) Update(ctx *gin.Context) {
	var (
		req types.ProductUpdateReq

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
	if err := logic.Update(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}

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

// 产品转换
func (con ProductController) Conversion(ctx *gin.Context) {
	var (
		req types.ProductConversionReq

		logic = product.ProductLogic{
			Ctx:   ctx,
			Staff: con.GetStaff(ctx),
		}
	)

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
