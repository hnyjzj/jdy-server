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

type OrderSalesController struct {
	controller.BaseController
}

// 订单筛选条件
func (con OrderSalesController) Where(ctx *gin.Context) {
	where := utils.StructToWhere(types.OrderSalesWhere{})

	con.Success(ctx, "ok", where)
}

// 创建订单
func (con OrderSalesController) Create(ctx *gin.Context) {
	var (
		req types.OrderSalesCreateReq

		logic = order.OrderSalesLogic{
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
func (con OrderSalesController) List(ctx *gin.Context) {
	var (
		req types.OrderSalesListReq

		logic = order.OrderSalesLogic{
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
	data, err := logic.List(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", data)
}

// 订单详情
func (con OrderSalesController) Info(ctx *gin.Context) {
	var (
		req types.OrderSalesInfoReq

		logic = order.OrderSalesLogic{
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
	data, err := logic.Info(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", data)
}

// 订单撤销
func (con OrderSalesController) Revoked(ctx *gin.Context) {
	var (
		req types.OrderSalesRevokedReq

		logic = order.OrderSalesLogic{
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
	err := logic.Revoked(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}

// 订单支付
func (con OrderSalesController) Pay(ctx *gin.Context) {
	var (
		req types.OrderSalesPayReq

		logic = order.OrderSalesLogic{
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
	err := logic.Pay(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}

// 订单退货
func (con OrderSalesController) Refund(ctx *gin.Context) {
	var (
		req types.OrderSalesRefundReq

		logic = order.OrderSalesLogic{
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

	// 判断参数
	if err := req.Validate(); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	// 调用逻辑层
	err := logic.Refund(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}
