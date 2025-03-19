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

	staff, err := con.GetStaff(ctx)
	if err != nil {
		con.ExceptionWithAuth(ctx, err.Error())
		return
	}
	logic.Staff = staff

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
	// 验证参数
	if err := req.Validate(); err != nil {
		con.Exception(ctx, err.Error())
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

	logic.IP = ctx.ClientIP()

	if err := logic.Create(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}
