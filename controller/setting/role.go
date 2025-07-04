package setting

import (
	"jdy/controller"
	"jdy/errors"
	"jdy/logic/setting"
	"jdy/types"
	"jdy/utils"
	"log"

	"github.com/gin-gonic/gin"
)

type RoleController struct {
	controller.BaseController
}

func (con RoleController) Where(ctx *gin.Context) {
	where := utils.StructToWhere(types.RoleWhere{})

	con.Success(ctx, "ok", where)
}

func (con RoleController) GetIdentity(ctx *gin.Context) {
	// 获取当前用户
	staff, err := con.GetStaff(ctx)
	if err != nil {
		con.ExceptionWithAuth(ctx, err)
		return
	}

	data := staff.Identity.GetMinMap()

	con.Success(ctx, "ok", data)
}

func (con RoleController) List(ctx *gin.Context) {
	var (
		req types.RoleListReq

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
		logic.Ctx = ctx
		logic.Staff = staff
		logic.IP = ctx.ClientIP()
	}

	data, err := logic.List(&req)
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
		logic.IP = ctx.ClientIP()
	}

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

func (con RoleController) Edit(ctx *gin.Context) {
	var (
		req   types.RoleEditReq
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

	if err := logic.Edit(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
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
