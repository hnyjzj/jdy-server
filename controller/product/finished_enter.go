package product

import (
	"jdy/controller"
	"jdy/errors"
	"jdy/logic/product"
	"jdy/types"
	"jdy/utils"

	"github.com/gin-gonic/gin"
)

type ProductFinishedEnterController struct {
	controller.BaseController
}

// 入库单筛选条件
func (con ProductFinishedEnterController) Where(ctx *gin.Context) {
	where := utils.StructToWhere(types.ProductFinishedEnterWhere{})

	con.Success(ctx, "ok", where)
}

// 入库单
func (con ProductFinishedEnterController) Create(ctx *gin.Context) {
	var (
		req types.ProductFinishedEnterCreateReq

		logic = product.ProductFinishedEnterLogic{
			Ctx: ctx,
		}
	)

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	if staff, err := con.GetStaff(ctx); err != nil {
		con.ExceptionWithAuth(ctx, err)
		return
	} else {
		logic.Staff = staff
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
func (con ProductFinishedEnterController) List(ctx *gin.Context) {
	var (
		req types.ProductFinishedEnterListReq

		logic = product.ProductFinishedEnterLogic{
			Ctx: ctx,
		}
	)

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	if staff, err := con.GetStaff(ctx); err != nil {
		con.ExceptionWithAuth(ctx, err)
		return
	} else {
		logic.Staff = staff
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
func (con ProductFinishedEnterController) Info(ctx *gin.Context) {
	var (
		req types.ProductFinishedEnterInfoReq

		logic = product.ProductFinishedEnterLogic{
			Ctx: ctx,
		}
	)

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	if staff, err := con.GetStaff(ctx); err != nil {
		con.ExceptionWithAuth(ctx, err)
		return
	} else {
		logic.Staff = staff
	}

	// 调用逻辑层
	res, err := logic.EnterInfo(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", res)
}

// 入库单添加产品
func (con ProductFinishedEnterController) AddProduct(ctx *gin.Context) {
	var (
		req types.ProductFinishedEnterAddProductReq

		logic = product.ProductFinishedEnterLogic{
			Ctx: ctx,
		}
	)

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	if staff, err := con.GetStaff(ctx); err != nil {
		con.ExceptionWithAuth(ctx, err)
		return
	} else {
		logic.Staff = staff
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
func (con ProductFinishedEnterController) EditProduct(ctx *gin.Context) {
	var (
		req types.ProductFinishedEnterEditProductReq

		logic = product.ProductFinishedEnterLogic{
			Ctx: ctx,
		}
	)

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	if staff, err := con.GetStaff(ctx); err != nil {
		con.ExceptionWithAuth(ctx, err)
		return
	} else {
		logic.Staff = staff
	}

	// 调用逻辑层
	if err := logic.EditProduct(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}

// 入库单删除产品
func (con ProductFinishedEnterController) DelProduct(ctx *gin.Context) {
	var (
		req types.ProductFinishedEnterDelProductReq

		logic = product.ProductFinishedEnterLogic{
			Ctx: ctx,
		}
	)

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	if staff, err := con.GetStaff(ctx); err != nil {
		con.ExceptionWithAuth(ctx, err)
		return
	} else {
		logic.Staff = staff
	}

	// 调用逻辑层
	if err := logic.DelProduct(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}

// 入库单清空产品
func (con ProductFinishedEnterController) ClearProduct(ctx *gin.Context) {
	var (
		req types.ProductFinishedEnterClearProductReq

		logic = product.ProductFinishedEnterLogic{
			Ctx: ctx,
		}
	)

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	if staff, err := con.GetStaff(ctx); err != nil {
		con.ExceptionWithAuth(ctx, err)
		return
	} else {
		logic.Staff = staff
	}

	// 调用逻辑层
	if err := logic.ClearProduct(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}

// 入库单完成
func (con ProductFinishedEnterController) Finish(ctx *gin.Context) {
	var (
		req types.ProductFinishedEnterFinishReq

		logic = product.ProductFinishedEnterLogic{
			Ctx: ctx,
		}
	)

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	if staff, err := con.GetStaff(ctx); err != nil {
		con.ExceptionWithAuth(ctx, err)
		return
	} else {
		logic.Staff = staff
	}

	// 调用逻辑层
	if err := logic.Finish(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}

// 入库单取消
func (con ProductFinishedEnterController) Cancel(ctx *gin.Context) {
	var (
		req types.ProductFinishedEnterCancelReq

		logic = product.ProductFinishedEnterLogic{
			Ctx: ctx,
		}
	)

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	if staff, err := con.GetStaff(ctx); err != nil {
		con.ExceptionWithAuth(ctx, err)
		return
	} else {
		logic.Staff = staff
	}

	// 调用逻辑层
	if err := logic.Cancel(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}
