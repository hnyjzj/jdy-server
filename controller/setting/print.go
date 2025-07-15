package setting

import (
	"jdy/controller"
	"jdy/errors"
	"jdy/logic/setting"
	"jdy/types"
	"log"

	"github.com/gin-gonic/gin"
)

type PrintController struct {
	controller.BaseController
}

func (con PrintController) Create(ctx *gin.Context) {
	var (
		req   types.PrintCreateReq
		logic = &setting.PrintLogic{}
	)

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		log.Printf("bind err: %v", err)
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}
	// 设置上下文
	logic.Ctx = ctx

	if _, err := con.GetStaff(ctx); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	if err := logic.Create(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}

func (con PrintController) List(ctx *gin.Context) {
	var (
		req   types.PrintListReq
		logic = &setting.PrintLogic{}
	)

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		log.Printf("bind err: %v", err)
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}
	// 设置上下文
	logic.Ctx = ctx

	if _, err := con.GetStaff(ctx); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	data, err := logic.List(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", data)
}

func (con PrintController) Info(ctx *gin.Context) {
	var (
		req   types.PrintInfoReq
		logic = &setting.PrintLogic{}
	)

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		log.Printf("bind err: %v", err)
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}
	// 设置上下文
	logic.Ctx = ctx

	if _, err := con.GetStaff(ctx); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	data, err := logic.Info(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", data)
}

func (con PrintController) Update(ctx *gin.Context) {
	var (
		req   types.PrintUpdateReq
		logic = &setting.PrintLogic{}
	)

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		log.Printf("bind err: %v", err)
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}
	// 设置上下文
	logic.Ctx = ctx

	if _, err := con.GetStaff(ctx); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	if err := logic.Update(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}

func (con PrintController) Delete(ctx *gin.Context) {
	var (
		req   types.PrintDeleteReq
		logic = &setting.PrintLogic{}
	)

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		log.Printf("bind err: %v", err)
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}
	// 设置上下文
	logic.Ctx = ctx

	if _, err := con.GetStaff(ctx); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	if err := logic.Delete(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}

func (con PrintController) Copy(ctx *gin.Context) {
	var (
		req   types.PrintCopyReq
		logic = &setting.PrintLogic{}
	)

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		log.Printf("bind err: %v", err)
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}
	// 设置上下文
	logic.Ctx = ctx

	if _, err := con.GetStaff(ctx); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	if err := logic.Copy(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}
