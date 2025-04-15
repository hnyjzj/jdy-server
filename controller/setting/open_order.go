package setting

import (
	"jdy/controller"
	"jdy/errors"
	"jdy/logic/setting"
	"jdy/types"
	"jdy/utils"

	"github.com/gin-gonic/gin"
)

type OpenOrderController struct {
	controller.BaseController
}

func (con OpenOrderController) Where(ctx *gin.Context) {
	where := utils.StructToWhere(types.OpenOrderWhere{})

	con.Success(ctx, "ok", where)
}

func (con OpenOrderController) Info(ctx *gin.Context) {
	var (
		req   types.OpenOrderInfoReq
		logic = &setting.OpenOrderLogic{}
	)

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	// 设置上下文
	logic.Ctx = ctx
	staff, err := con.GetStaff(ctx)
	if err != nil {
		con.ExceptionWithAuth(ctx, err.Error())
		return
	}
	logic.Staff = staff

	res, err := logic.Info(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", res)
}

func (con OpenOrderController) Update(ctx *gin.Context) {
	var (
		req   types.OpenOrderUpdateReq
		logic = &setting.OpenOrderLogic{}
	)

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	// 设置上下文
	logic.Ctx = ctx
	staff, err := con.GetStaff(ctx)
	if err != nil {
		con.ExceptionWithAuth(ctx, err.Error())
		return
	}
	logic.Staff = staff

	if err = logic.Update(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}
