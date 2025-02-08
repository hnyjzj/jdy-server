package order

import (
	"jdy/controller"
	"jdy/errors"
	"jdy/logic/order"
	"jdy/types"
	"jdy/utils"

	"github.com/gin-gonic/gin"
)

type OrderController struct {
	controller.BaseController
}

// 订单筛选条件
func (con OrderController) Where(ctx *gin.Context) {
	where := utils.StructToWhere(types.OrderWhere{})

	con.Success(ctx, "ok", where)
}

// 创建订单
func (con OrderController) Create(ctx *gin.Context) {
	var (
		req types.OrderCreateReq

		logic = order.OrderLogic{
			Ctx:   ctx,
			Staff: con.GetStaff(ctx),
		}
	)

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	// 校验参数
	if err := req.Validate(); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	// 调用逻辑层
	if err := logic.Create(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}

// 订单列表
func (con OrderController) List(ctx *gin.Context) {
	var (
		req types.OrderListReq

		logic = order.OrderLogic{
			Ctx:   ctx,
			Staff: con.GetStaff(ctx),
		}
	)

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
