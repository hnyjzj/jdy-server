package product

import (
	"jdy/controller"
	"jdy/errors"
	"jdy/logic/product"
	"jdy/types"
	"jdy/utils"

	"github.com/gin-gonic/gin"
)

type ProductAccessorieCategoryController struct {
	controller.BaseController
}

// 配件条目筛选
func (con ProductAccessorieCategoryController) Where(ctx *gin.Context) {
	where := utils.StructToWhere(types.ProductAccessorieCategoryWhere{})

	con.Success(ctx, "ok", where)
}

// 配件条目列表
func (con ProductAccessorieCategoryController) List(ctx *gin.Context) {
	var (
		req types.ProductAccessorieCategoryListReq

		logic = product.ProductAccessorieCategoryLogic{
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

// 配件条目详情
func (con ProductAccessorieCategoryController) Info(ctx *gin.Context) {
	var (
		req types.ProductAccessorieCategoryInfoReq

		logic = product.ProductAccessorieCategoryLogic{
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

// 新增配件条目
func (con ProductAccessorieCategoryController) Create(ctx *gin.Context) {
	var (
		req types.ProductAccessorieCategoryCreateReq

		logic = product.ProductAccessorieCategoryLogic{
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
	if err := logic.Create(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}

// 更新配件条目
func (con ProductAccessorieCategoryController) Update(ctx *gin.Context) {
	var (
		req types.ProductAccessorieCategoryUpdateReq

		logic = product.ProductAccessorieCategoryLogic{
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
	if err := logic.Update(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}

// 删除配件条目
func (con ProductAccessorieCategoryController) Delete(ctx *gin.Context) {
	var (
		req types.ProductAccessorieCategoryDeleteReq

		logic = product.ProductAccessorieCategoryLogic{
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
	if err := logic.Delete(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}
