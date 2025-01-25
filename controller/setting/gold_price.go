package setting

import (
	"jdy/controller"
	"jdy/errors"
	"jdy/logic/setting"
	"jdy/types"

	"github.com/gin-gonic/gin"
)

type GoldPriceController struct {
	controller.BaseController
}

func (con GoldPriceController) Create(ctx *gin.Context) {
	var (
		req   types.GoldPriceCreateReq
		logic = &setting.GoldPriceLogic{}
	)

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	// 设置上下文
	logic.Ctx = ctx
	logic.Staff = con.GetStaff(ctx)
	logic.IP = ctx.ClientIP()

	if err := logic.Create(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}

func (con GoldPriceController) Update(ctx *gin.Context) {
	var (
		req types.GoldPriceUpdateReq

		logic = &setting.GoldPriceLogic{}
	)

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	// 设置上下文
	logic.Ctx = ctx
	logic.Staff = con.GetStaff(ctx)
	logic.IP = ctx.ClientIP()

	if err := logic.Update(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}
