package setting

import (
	"jdy/controller"
	"jdy/errors"
	"jdy/logic/setting"
	"jdy/types"
	"log"

	"github.com/gin-gonic/gin"
)

type RoleController struct {
	controller.BaseController
}

func (con RoleController) List(ctx *gin.Context) {
	var (
		logic = &setting.RoleLogic{}
	)

	// 设置上下文
	logic.Ctx = ctx

	if staff, err := con.GetStaff(ctx); err != nil {
		con.ExceptionWithAuth(ctx, err)
		return
	} else {
		logic.Staff = staff
		logic.IP = ctx.ClientIP()
	}

	data, err := logic.List()
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", data)
}

func (con RoleController) Create(ctx *gin.Context) {
	var (
		req   types.RoleCreateReq
		logic = &setting.RoleLogic{}
	)

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	// 设置上下文
	if staff, err := con.GetStaff(ctx); err != nil {
		con.ExceptionWithAuth(ctx, err)
		return
	} else {
		logic.Staff = staff
		logic.Ctx = ctx
	}

	logic.IP = ctx.ClientIP()

	data, err := logic.Create(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", data)
}

func (con RoleController) Info(ctx *gin.Context) {
	var (
		req   types.RoleInfoReq
		logic = &setting.RoleLogic{}
	)

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	// 设置上下文
	if staff, err := con.GetStaff(ctx); err != nil {
		con.ExceptionWithAuth(ctx, err)
		return
	} else {
		logic.Staff = staff
		logic.Ctx = ctx
		logic.IP = ctx.ClientIP()
	}

	data, err := logic.Info(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", data)
}

func (con RoleController) Update(ctx *gin.Context) {
	var (
		req   types.RoleUpdateReq
		logic = &setting.RoleLogic{}
	)

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		log.Printf("err: %v", err.Error())
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	// 设置上下文
	if staff, err := con.GetStaff(ctx); err != nil {
		con.ExceptionWithAuth(ctx, err)
		return
	} else {
		logic.Staff = staff
		logic.Ctx = ctx
		logic.IP = ctx.ClientIP()
	}

	if err := logic.Update(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}

func (con RoleController) AddStaff(ctx *gin.Context) {
	var (
		req   types.RoleAddStaffReq
		logic = &setting.RoleLogic{}
	)

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		log.Printf("err: %v", err)
		return
	}

	// 设置上下文
	if staff, err := con.GetStaff(ctx); err != nil {
		con.ExceptionWithAuth(ctx, err)
		return
	} else {
		logic.Staff = staff
		logic.Ctx = ctx
		logic.IP = ctx.ClientIP()
	}

	if err := logic.AddStaff(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}

func (con RoleController) Delete(ctx *gin.Context) {
	var (
		req   types.RoleDeleteReq
		logic = &setting.RoleLogic{}
	)

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	// 设置上下文
	if staff, err := con.GetStaff(ctx); err != nil {
		con.ExceptionWithAuth(ctx, err)
		return
	} else {
		logic.Staff = staff
		logic.Ctx = ctx
		logic.IP = ctx.ClientIP()
	}

	if err := logic.Delete(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}

func (con RoleController) Apis(ctx *gin.Context) {
	var (
		logic = &setting.RoleLogic{}
	)

	// 设置上下文
	if staff, err := con.GetStaff(ctx); err != nil {
		con.ExceptionWithAuth(ctx, err)
		return
	} else {
		logic.Staff = staff
		logic.Ctx = ctx
		logic.IP = ctx.ClientIP()
	}

	data, err := logic.Apis()
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", data)
}
