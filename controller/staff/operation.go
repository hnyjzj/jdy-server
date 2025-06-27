package staff

import (
	"jdy/errors"
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

	// 使用自定义验证器进行验证
	if err := req.Validate(); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
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

func (con StaffController) Update(ctx *gin.Context) {
	var (
		req types.StaffUpdateReq

		logic = staff.StaffLogic{}
	)

	// 解析参数
	if staff, err := con.GetStaff(ctx); err != nil {
		con.ExceptionWithAuth(ctx, err)
		return
	} else {
		logic.Staff = staff
	}

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	// 使用自定义验证器进行验证
	if err := req.Validate(); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	if err := logic.StaffUpdate(ctx, logic.Staff.Id, &req); err != nil {
		con.ErrorLogic(ctx, err)
		return
	}

	// 返回结果
	con.Success(ctx, "ok", nil)
}
