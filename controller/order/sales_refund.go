package order

import (
	"jdy/controller"
	"jdy/errors"
	"jdy/logic/order"
	"jdy/types"
	"jdy/utils"

	"github.com/gin-gonic/gin"
)

type OrderSalesRefundController struct {
	controller.BaseController
}

// 订单筛选条件
func (con OrderSalesRefundController) Where(ctx *gin.Context) {
	where := utils.StructToWhere(types.OrderSalesRefundWhere{})

	con.Success(ctx, "ok", where)
}

// 订单列表
func (con OrderSalesRefundController) List(ctx *gin.Context) {
	var (
		req types.OrderSalesRefundListReq

		logic = order.OrderSalesRefundLogic{
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
