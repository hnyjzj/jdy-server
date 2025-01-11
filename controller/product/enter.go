package product

import (
	"jdy/controller"
	"jdy/errors"
	"jdy/logic/product"
	"jdy/types"
	"jdy/utils"

	"github.com/gin-gonic/gin"
)

type ProductEnterController struct {
	controller.BaseController
}

// 入库单筛选条件
func (con ProductEnterController) Where(ctx *gin.Context) {
	where := utils.StructToWhere(types.ProductEnterWhere{})

	con.Success(ctx, "ok", where)
}

// 入库单
func (con ProductEnterController) Create(ctx *gin.Context) {
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

// 入库单列表
func (con ProductEnterController) List(ctx *gin.Context) {
	var (
		req types.ProductEnterListReq

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
	res, err := logic.EnterList(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", res)
}
