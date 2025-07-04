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

type OrderSalesDetailController struct {
	controller.BaseController
}

// 订单筛选条件
func (con OrderSalesDetailController) Where(ctx *gin.Context) {
	where := utils.StructToWhere(types.OrderSalesDetailWhere{})

	con.Success(ctx, "ok", where)
}

// 订单列表
func (con OrderSalesDetailController) List(ctx *gin.Context) {
	var (
		req types.OrderSalesDetailListReq

		logic = order.OrderSalesDetailLogic{
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
		log.Printf("参数绑定失败: %v", err.Error())
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
func (con OrderSalesDetailController) Info(ctx *gin.Context) {
	var (
		req types.OrderSalesDetailInfoReq

		logic = order.OrderSalesDetailLogic{
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
