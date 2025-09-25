package product

import (
	"jdy/controller"
	"jdy/errors"
	"jdy/logic/product"
	"jdy/types"
	"jdy/utils"

	"github.com/gin-gonic/gin"
)

type ProductAllocateController struct {
	controller.BaseController
}

// 创建产品调拨单
func (con ProductAllocateController) Create(ctx *gin.Context) {
	var (
		req types.ProductAllocateCreateReq

		logic = product.ProductAllocateLogic{
			Ctx: ctx,
		}
	)

	// 绑定请求参数
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
func (con ProductAllocateController) Where(ctx *gin.Context) {
	where := utils.StructToWhere(types.ProductAllocateWhere{})

	con.Success(ctx, "ok", where)
}

// 获取产品调拨单列表
func (con ProductAllocateController) List(ctx *gin.Context) {
	var (
		req types.ProductAllocateListReq

		logic = product.ProductAllocateLogic{
			Ctx: ctx,
		}
	)

	// 绑定请求参数
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

// 获取产品调拨单明细列表
func (con ProductAllocateController) Details(ctx *gin.Context) {
	var (
		req types.ProductAllocateDetailsReq

		logic = product.ProductAllocateLogic{
			Ctx: ctx,
		}
	)

	// 绑定请求参数
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

	// 校验参数
	if err := req.Where.Validate(); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	if req.Where.StartTime == nil || req.Where.EndTime == nil {
		con.Exception(ctx, "请选择时间范围")
		return
	}

	// 获取产品调拨单列表
	res, err := logic.Details(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", res)
}

// 获取产品调拨单详情
func (con ProductAllocateController) Info(ctx *gin.Context) {
	var (
		req types.ProductAllocateInfoReq

		logic = product.ProductAllocateLogic{
			Ctx: ctx,
		}
	)

	// 绑定请求参数
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

	// 获取产品调拨单详情
	res, err := logic.Info(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", res)
}

// 获取产品调拨单概览
func (con ProductAllocateController) InfoOverview(ctx *gin.Context) {
	var (
		req types.ProductAllocateInfoOverviewReq

		logic = product.ProductAllocateLogic{
			Ctx: ctx,
		}
	)

	// 绑定请求参数
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

	// 获取产品调拨单概览
	res, err := logic.InfoOverview(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", res)
}

// 添加调拨单产品
func (con ProductAllocateController) Add(ctx *gin.Context) {
	var (
		req types.ProductAllocateAddReq

		logic = product.ProductAllocateLogic{
			Ctx: ctx,
		}
	)

	// 绑定请求参数
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

	if err := logic.Add(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}

// 移除调拨单产品
func (con ProductAllocateController) Remove(ctx *gin.Context) {
	var (
		req types.ProductAllocateRemoveReq

		logic = product.ProductAllocateLogic{
			Ctx: ctx,
		}
	)

	// 绑定请求参数
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

	if err := logic.Remove(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}

// 清空调拨单产品
func (con ProductAllocateController) Clear(ctx *gin.Context) {
	var (
		req types.ProductAllocateClearReq

		logic = product.ProductAllocateLogic{
			Ctx: ctx,
		}
	)

	// 绑定请求参数
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

	if err := logic.Clear(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}

// 确认调拨
func (con ProductAllocateController) Confirm(ctx *gin.Context) {
	var (
		req types.ProductAllocateConfirmReq

		logic = product.ProductAllocateLogic{
			Ctx: ctx,
		}
	)

	// 绑定请求参数
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

	if err := logic.Confirm(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}

// 取消调拨
func (con ProductAllocateController) Cancel(ctx *gin.Context) {
	var (
		req types.ProductAllocateCancelReq

		logic = product.ProductAllocateLogic{
			Ctx: ctx,
		}
	)

	// 绑定请求参数
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

	if err := logic.Cancel(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}

// 完成调拨
func (con ProductAllocateController) Complete(ctx *gin.Context) {
	var (
		req types.ProductAllocateCompleteReq

		logic = product.ProductAllocateLogic{
			Ctx: ctx,
		}
	)

	// 绑定请求参数
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

	if err := logic.Complete(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}
