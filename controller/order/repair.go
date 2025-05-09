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

type OrderRepairController struct {
	controller.BaseController
}

// 订单筛选条件
func (con OrderRepairController) Where(ctx *gin.Context) {
	where := utils.StructToWhere(types.OrderRepairWhere{})

	con.Success(ctx, "ok", where)
}

func (con OrderRepairController) WhereProduct(ctx *gin.Context) {
	where := utils.StructToWhere(types.OrderRepairWhereProduct{})

	con.Success(ctx, "ok", where)
}

// 创建订单
func (con OrderRepairController) Create(ctx *gin.Context) {
	var (
		req types.OrderRepairCreateReq

		logic = order.OrderRepairLogic{
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
func (con OrderRepairController) List(ctx *gin.Context) {
	var (
		req types.OrderRepairListReq

		logic = order.OrderRepairLogic{
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

	// 调用逻辑层
	data, err := logic.List(&req)
	if err != nil {
		log.Printf("err.Error(): %v\n", err.Error())
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", data)
}

// 订单详情
func (con OrderRepairController) Info(ctx *gin.Context) {
	var (
		req types.OrderRepairInfoReq

		logic = order.OrderRepairLogic{
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

// 订单修改
func (con OrderRepairController) Update(ctx *gin.Context) {
	var (
		req types.OrderRepairUpdateReq

		logic = order.OrderRepairLogic{
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
	err = logic.Update(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}

// 订单操作
func (con OrderRepairController) Operation(ctx *gin.Context) {
	var (
		req types.OrderRepairOperationReq

		logic = order.OrderRepairLogic{
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
	err = logic.Operation(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}

// 订单撤销
func (con OrderRepairController) Revoked(ctx *gin.Context) {
	var (
		req types.OrderRepairRevokedReq

		logic = order.OrderRepairLogic{
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
func (con OrderRepairController) Pay(ctx *gin.Context) {
	var (
		req types.OrderRepairPayReq

		logic = order.OrderRepairLogic{
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

// 退款
func (con OrderRepairController) Refund(ctx *gin.Context) {
	var (
		req types.OrderRepairRefundReq

		logic = order.OrderRepairLogic{
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
