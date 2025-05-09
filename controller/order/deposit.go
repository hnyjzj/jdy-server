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

type OrderDepositController struct {
	controller.BaseController
}

// 订单筛选条件
func (con OrderDepositController) Where(ctx *gin.Context) {
	where := utils.StructToWhere(types.OrderDepositWhere{})

	con.Success(ctx, "ok", where)
}

// 创建订单
func (con OrderDepositController) Create(ctx *gin.Context) {
	var (
		req types.OrderDepositCreateReq

		logic = order.OrderDepositLogic{
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
func (con OrderDepositController) List(ctx *gin.Context) {
	var (
		req types.OrderDepositListReq

		logic = order.OrderDepositLogic{
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
	data, err := logic.List(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", data)
}

// 订单详情
func (con OrderDepositController) Info(ctx *gin.Context) {
	var (
		req types.OrderDepositInfoReq

		logic = order.OrderDepositLogic{
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

// 订单撤销
func (con OrderDepositController) Revoked(ctx *gin.Context) {
	var (
		req types.OrderDepositRevokedReq

		logic = order.OrderDepositLogic{
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
	err = logic.Revoked(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}

// 订单支付
func (con OrderDepositController) Pay(ctx *gin.Context) {
	var (
		req types.OrderDepositPayReq

		logic = order.OrderDepositLogic{
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
	err = logic.Pay(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}

func (con OrderDepositController) Refund(ctx *gin.Context) {
	var (
		req types.OrderDepositRefundReq

		logic = order.OrderDepositLogic{
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
	err = logic.Refund(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}
