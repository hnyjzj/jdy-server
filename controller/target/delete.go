package target

import (
	"jdy/errors"
	"jdy/logic/target"
	"jdy/types"

	"github.com/gin-gonic/gin"
)

// 销售目标删除
func (con TargetController) Delete(ctx *gin.Context) {
	var (
		req   types.TargetDeleteReq
		logic target.Logic
	)

	if staff, err := con.GetStaff(ctx); err != nil {
		con.ExceptionWithAuth(ctx, err)
		return
	} else {
		logic.Staff = staff
		logic.Ctx = ctx
	}

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	// 删除目标
	if err := logic.Delete(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}

// 销售目标分组删除
func (con TargetController) DeleteGroup(ctx *gin.Context) {
	var (
		req   types.TargetDeleteGroupReq
		logic target.Logic
	)

	if staff, err := con.GetStaff(ctx); err != nil {
		con.ExceptionWithAuth(ctx, err)
		return
	} else {
		logic.Staff = staff
		logic.Ctx = ctx
	}

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	// 删除目标分组
	if err := logic.DeleteGroup(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}

// 销售目标个人删除
func (con TargetController) DeletePersonal(ctx *gin.Context) {
	var (
		req   types.TargetDeletePersonalReq
		logic target.Logic
	)

	if staff, err := con.GetStaff(ctx); err != nil {
		con.ExceptionWithAuth(ctx, err)
		return
	} else {
		logic.Staff = staff
		logic.Ctx = ctx
	}

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	// 删除个人
	if err := logic.DeletePersonal(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}
