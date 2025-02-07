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

func (con GoldPriceController) Get(ctx *gin.Context) {
	var (
		logic = &setting.GoldPriceLogic{}
	)

	data, err := logic.Get()
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", data)
}

func (con GoldPriceController) List(ctx *gin.Context) {
	var (
		req   types.GoldPriceListReq
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

	data, err := logic.List(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", data)
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
