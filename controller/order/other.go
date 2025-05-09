package order

import (
	"jdy/controller"
	"jdy/errors"
	"jdy/logic/order"
	"jdy/types"
	"jdy/utils"
	"log"

	"github.com/gin-gonic/gin"
)

type OrderOtherController struct {
	controller.BaseController
}

// 订单筛选条件
func (con OrderOtherController) Where(ctx *gin.Context) {
	where := utils.StructToWhere(types.OrderOtherWhere{})

	con.Success(ctx, "ok", where)
}

// 创建订单
func (con OrderOtherController) Create(ctx *gin.Context) {
	var (
		req types.OrderOtherCreateReq

		logic = order.OrderOtherLogic{
			Ctx: ctx,
		}
	)

	staff, err := con.GetStaff(ctx)
	if err != nil {
		con.ExceptionWithAuth(ctx, err.Error())
		return
	}
	logic.Staff = staff

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		log.Printf("err.Error(): %v\n", err.Error())
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	// 校验参数
	if err := req.Validate(); err != nil {
		log.Printf("err.Error(): %v\n", err.Error())
		con.Exception(ctx, err.Error())
		return
	}

	// 调用逻辑层
	order, err := logic.Create(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", order)
}

// 订单列表
func (con OrderOtherController) List(ctx *gin.Context) {
	var (
		req types.OrderOtherListReq

		logic = order.OrderOtherLogic{
			Ctx: ctx,
		}
	)

	staff, err := con.GetStaff(ctx)
	if err != nil {
		con.ExceptionWithAuth(ctx, err.Error())
		return
	}
	logic.Staff = staff

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		log.Printf("参数错误: %v\n", err.Error())
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	// 调用逻辑层
	data, err := logic.List(&req)
	if err != nil {
		log.Printf("参数错误: %v\n", err.Error())
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", data)
}

// 订单详情
func (con OrderOtherController) Info(ctx *gin.Context) {
	var (
		req types.OrderOtherInfoReq

		logic = order.OrderOtherLogic{
			Ctx: ctx,
		}
	)

	staff, err := con.GetStaff(ctx)
	if err != nil {
		con.ExceptionWithAuth(ctx, err.Error())
		return
	}
	logic.Staff = staff

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	// 调用逻辑层
	data, err := logic.Info(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", data)
}

// 更新订单
func (con OrderOtherController) Update(ctx *gin.Context) {
	var (
		req types.OrderOtherUpdateReq

		logic = order.OrderOtherLogic{
			Ctx: ctx,
		}
	)

	staff, err := con.GetStaff(ctx)
	if err != nil {
		con.ExceptionWithAuth(ctx, err.Error())
		return
	}
	logic.Staff = staff

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	// 调用逻辑层
	order, err := logic.Update(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", order)
}

// 删除订单
func (con OrderOtherController) Delete(ctx *gin.Context) {
	var (
		req types.OrderOtherDeleteReq

		logic = order.OrderOtherLogic{
			Ctx: ctx,
		}
	)

	staff, err := con.GetStaff(ctx)
	if err != nil {
		con.ExceptionWithAuth(ctx, err.Error())
		return
	}
	logic.Staff = staff

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	// 调用逻辑层
	err = logic.Delete(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}
