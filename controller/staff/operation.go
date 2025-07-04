package staff

import (
	"jdy/errors"
	"jdy/logic"
	"jdy/logic/staff"
	"jdy/types"

	"github.com/gin-gonic/gin"
)

// 创建员工
func (con StaffController) Create(ctx *gin.Context) {
	var (
		req types.StaffReq

		logic = staff.StaffLogic{}
	)

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	// 获取当前用户
	if staff, err := con.GetStaff(ctx); err != nil {
		con.ExceptionWithAuth(ctx, err)
		return
	} else {
		logic.Staff = staff
	}

	// 创建员工
	err := logic.StaffCreate(ctx, &req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	// 返回结果
	con.Success(ctx, "ok", nil)
}

func (con StaffController) Edit(ctx *gin.Context) {
	var (
		req types.StaffEditReq

		logic = staff.StaffLogic{
			BaseLogic: logic.BaseLogic{
				Ctx: ctx,
			},
		}
	)

	// 解析参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	// 获取当前用户
	if staff, err := con.GetStaff(ctx); err != nil {
		con.ExceptionWithAuth(ctx, err)
		return
	} else {
		logic.Staff = staff
	}

	if err := logic.StaffEdit(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	// 返回结果
	con.Success(ctx, "ok", logic.Staff)
}

func (con StaffController) Update(ctx *gin.Context) {
	var (
		req types.StaffUpdateReq

		logic = staff.StaffLogic{}
	)

	// 解析参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	// 获取当前用户
	if staff, err := con.GetStaff(ctx); err != nil {
		con.ExceptionWithAuth(ctx, err)
		return
	} else {
		logic.Staff = staff
	}

	if err := logic.StaffUpdate(ctx, logic.Staff.Id, &req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	// 返回结果
	con.Success(ctx, "ok", nil)
}
