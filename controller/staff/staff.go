package staff

import (
	"jdy/controller"
	"jdy/logic/staff"
	"jdy/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

type StaffController struct {
	controller.BaseController
}

// 创建员工
func (con StaffController) Create(ctx *gin.Context) {
	var (
		req types.StaffReq

		logic = staff.StaffLogic{}
	)

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	// 使用自定义验证器进行验证
	if err := req.ValidateStaffReq(&req); err != nil {
		con.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// 创建员工
	err := logic.CreateAccount(ctx, &req)
	if err != nil {
		con.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// 返回结果
	con.Success(ctx, "ok", nil)
}

// 获取员工信息
func (con StaffController) Info(ctx *gin.Context) {
	var (
		logic = staff.StaffLogic{}
	)

	staff := con.GetStaff(ctx)

	staffinfo, err := logic.GetStaffInfo(ctx, staff)
	if err != nil {
		con.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// 获取员工信息
	con.Success(ctx, "ok", staffinfo)
}
