package product

import (
	"jdy/controller"
	"jdy/errors"
	"jdy/logic/product"
	"jdy/types"
	"jdy/utils"

	"github.com/gin-gonic/gin"
)

type ProductFinishedAllocateController struct {
	controller.BaseController
}

// 创建产品调拨单
func (con ProductFinishedAllocateController) Create(ctx *gin.Context) {
	var (
		req types.ProductFinishedAllocateCreateReq

		logic = product.ProductFinishedAllocateLogic{
			Ctx: ctx,
		}
	)

	staff, err := con.GetStaff(ctx)
	if err != nil {
		con.ExceptionWithAuth(ctx, err.Error())
		return
	}
	logic.Staff = staff

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

	if err := logic.Create(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}

// 产品调拨单筛选条件
func (con ProductFinishedAllocateController) Where(ctx *gin.Context) {
	where := utils.StructToWhere(types.ProductFinishedAllocateWhere{})

	con.Success(ctx, "ok", where)
}

// 获取产品调拨单列表
func (con ProductFinishedAllocateController) List(ctx *gin.Context) {
	var (
		req types.ProductFinishedAllocateListReq

		logic = product.ProductFinishedAllocateLogic{
			Ctx: ctx,
		}
	)

	staff, err := con.GetStaff(ctx)
	if err != nil {
		con.ExceptionWithAuth(ctx, err.Error())
		return
	}
	logic.Staff = staff

	// 绑定请求参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	// 校验参数
	if err := req.Where.Validate(); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	// 获取产品调拨单列表
	res, err := logic.List(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", res)
}

// 获取产品调拨单详情
func (con ProductFinishedAllocateController) Info(ctx *gin.Context) {
	var (
		req types.ProductFinishedAllocateInfoReq

		logic = product.ProductFinishedAllocateLogic{
			Ctx: ctx,
		}
	)

	staff, err := con.GetStaff(ctx)
	if err != nil {
		con.ExceptionWithAuth(ctx, err.Error())
		return
	}
	logic.Staff = staff

	// 绑定请求参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	// 获取产品调拨单详情
	res, err := logic.Info(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", res)
}

// 添加调拨单产品
func (con ProductFinishedAllocateController) Add(ctx *gin.Context) {
	var (
		req types.ProductFinishedAllocateAddReq

		logic = product.ProductFinishedAllocateLogic{
			Ctx: ctx,
		}
	)

	staff, err := con.GetStaff(ctx)
	if err != nil {
		con.ExceptionWithAuth(ctx, err.Error())
		return
	}
	logic.Staff = staff

	// 绑定请求参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	if err := logic.Add(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}

// 移除调拨单产品
func (con ProductFinishedAllocateController) Remove(ctx *gin.Context) {
	var (
		req types.ProductFinishedAllocateRemoveReq

		logic = product.ProductFinishedAllocateLogic{
			Ctx: ctx,
		}
	)

	staff, err := con.GetStaff(ctx)
	if err != nil {
		con.ExceptionWithAuth(ctx, err.Error())
		return
	}
	logic.Staff = staff

	// 绑定请求参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	if err := logic.Remove(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}

// 确认调拨
func (con ProductFinishedAllocateController) Confirm(ctx *gin.Context) {
	var (
		req types.ProductFinishedAllocateConfirmReq

		logic = product.ProductFinishedAllocateLogic{
			Ctx: ctx,
		}
	)

	staff, err := con.GetStaff(ctx)
	if err != nil {
		con.ExceptionWithAuth(ctx, err.Error())
		return
	}
	logic.Staff = staff

	// 绑定请求参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	if err := logic.Confirm(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}

// 取消调拨
func (con ProductFinishedAllocateController) Cancel(ctx *gin.Context) {
	var (
		req types.ProductFinishedAllocateCancelReq

		logic = product.ProductFinishedAllocateLogic{
			Ctx: ctx,
		}
	)

	staff, err := con.GetStaff(ctx)
	if err != nil {
		con.ExceptionWithAuth(ctx, err.Error())
		return
	}
	logic.Staff = staff

	// 绑定请求参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	if err := logic.Cancel(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}

// 完成调拨
func (con ProductFinishedAllocateController) Complete(ctx *gin.Context) {
	var (
		req types.ProductFinishedAllocateCompleteReq

		logic = product.ProductFinishedAllocateLogic{
			Ctx: ctx,
		}
	)

	staff, err := con.GetStaff(ctx)
	if err != nil {
		con.ExceptionWithAuth(ctx, err.Error())
		return
	}
	logic.Staff = staff

	// 绑定请求参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	if err := logic.Complete(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}
