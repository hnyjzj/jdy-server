package product

import (
	"jdy/controller"
	"jdy/errors"
	"jdy/logic/product"
	"jdy/types"
	"jdy/utils"

	"github.com/gin-gonic/gin"
)

type ProductAccessorieEnterController struct {
	controller.BaseController
}

// 入库单筛选条件
func (con ProductAccessorieEnterController) Where(ctx *gin.Context) {
	where := utils.StructToWhere(types.ProductAccessorieEnterWhere{})

	con.Success(ctx, "ok", where)
}

// 入库单
func (con ProductAccessorieEnterController) Create(ctx *gin.Context) {
	var (
		req types.ProductAccessorieEnterCreateReq

		logic = product.ProductAccessorieEnterLogic{
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
	res, err := logic.Create(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", res)
}

// 入库单列表
func (con ProductAccessorieEnterController) List(ctx *gin.Context) {
	var (
		req types.ProductAccessorieEnterListReq

		logic = product.ProductAccessorieEnterLogic{
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
	res, err := logic.EnterList(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", res)
}

// 入库单详情
func (con ProductAccessorieEnterController) Info(ctx *gin.Context) {
	var (
		req types.ProductAccessorieEnterInfoReq

		logic = product.ProductAccessorieEnterLogic{
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
	res, err := logic.EnterInfo(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", res)
}

// 入库单筛选条件
func (con ProductAccessorieEnterController) WhereAddProduct(ctx *gin.Context) {
	where := utils.StructToWhere(types.ProductAccessorieEnterReqProduct{})

	con.Success(ctx, "ok", where)
}

// 入库单添加产品
func (con ProductAccessorieEnterController) AddProduct(ctx *gin.Context) {
	var (
		req types.ProductAccessorieEnterAddProductReq

		logic = product.ProductAccessorieEnterLogic{
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
	res, err := logic.AddProduct(&req)
	if err != nil {
		con.ExceptionWithResult(ctx, err.Error(), res)
		return
	}

	con.Success(ctx, "ok", res)
}

// 入库单编辑产品
func (con ProductAccessorieEnterController) EditProduct(ctx *gin.Context) {
	var (
		req types.ProductAccessorieEnterEditProductReq

		logic = product.ProductAccessorieEnterLogic{
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
	if err := logic.EditProduct(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}

// 入库单删除产品
func (con ProductAccessorieEnterController) DelProduct(ctx *gin.Context) {
	var (
		req types.ProductAccessorieEnterDelProductReq

		logic = product.ProductAccessorieEnterLogic{
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
	if err := logic.DelProduct(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}

// 入库单完成
func (con ProductAccessorieEnterController) Finish(ctx *gin.Context) {
	var (
		req types.ProductAccessorieEnterFinishReq

		logic = product.ProductAccessorieEnterLogic{
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
	if err := logic.Finish(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}

// 入库单取消
func (con ProductAccessorieEnterController) Cancel(ctx *gin.Context) {
	var (
		req types.ProductAccessorieEnterCancelReq

		logic = product.ProductAccessorieEnterLogic{
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
	if err := logic.Cancel(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}
