package product

import (
	"jdy/controller"
	"jdy/errors"
	"jdy/logic/product"
	"jdy/types"
	"jdy/utils"

	"github.com/gin-gonic/gin"
)

type ProductOldAllocateController struct {
	controller.BaseController
}

// 创建产品调拨单
func (con ProductOldAllocateController) Create(ctx *gin.Context) {
	var (
		req types.ProductOldAllocateCreateReq

		logic = product.ProductOldAllocateLogic{
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
func (con ProductOldAllocateController) Where(ctx *gin.Context) {
	where := utils.StructToWhere(types.ProductOldAllocateWhere{})

	con.Success(ctx, "ok", where)
}

// 获取产品调拨单列表
func (con ProductOldAllocateController) List(ctx *gin.Context) {
	var (
		req types.ProductOldAllocateListReq

		logic = product.ProductOldAllocateLogic{
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
func (con ProductOldAllocateController) Info(ctx *gin.Context) {
	var (
		req types.ProductOldAllocateInfoReq

		logic = product.ProductOldAllocateLogic{
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
func (con ProductOldAllocateController) Add(ctx *gin.Context) {
	var (
		req types.ProductOldAllocateAddReq

		logic = product.ProductOldAllocateLogic{
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
func (con ProductOldAllocateController) Remove(ctx *gin.Context) {
	var (
		req types.ProductOldAllocateRemoveReq

		logic = product.ProductOldAllocateLogic{
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
func (con ProductOldAllocateController) Confirm(ctx *gin.Context) {
	var (
		req types.ProductOldAllocateConfirmReq

		logic = product.ProductOldAllocateLogic{
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
func (con ProductOldAllocateController) Cancel(ctx *gin.Context) {
	var (
		req types.ProductOldAllocateCancelReq

		logic = product.ProductOldAllocateLogic{
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
func (con ProductOldAllocateController) Complete(ctx *gin.Context) {
	var (
		req types.ProductOldAllocateCompleteReq

		logic = product.ProductOldAllocateLogic{
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
