package target

import (
	"jdy/errors"
	"jdy/logic/target"
	"jdy/types"

	"github.com/gin-gonic/gin"
)

// 销售目标更新
func (con TargetController) Update(ctx *gin.Context) {
	var (
		req   types.TargetUpdateReq
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

	// 更新目标
	if err := logic.Update(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}

// 销售目标分组更新
func (con TargetController) UpdateGroup(ctx *gin.Context) {
	var (
		req   types.TargetUpdateGroupReq
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

	// 更新目标分组
	if err := logic.UpdateGroup(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}

// 销售目标个人更新
func (con TargetController) UpdatePersonal(ctx *gin.Context) {
	var (
		req   types.TargetUpdatePersonalReq
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

	// 验证参数
	if err := req.Validate(); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	// 更新个人
	if err := logic.UpdatePersonal(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}
