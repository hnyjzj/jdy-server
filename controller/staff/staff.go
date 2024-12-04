package staff

import (
	"jdy/controller"
	"jdy/logic/staff"
	"net/http"

	"github.com/gin-gonic/gin"
)

type StaffController struct {
	controller.BaseController
}

// 获取员工信息
func (con StaffController) Info(ctx *gin.Context) {
	var (
		logic = staff.StaffLogic{}
	)

	staff := con.GetStaff(ctx)

	staffinfo, err := logic.GetStaffInfo(ctx, &staff.Id)
	if err != nil {
		con.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// 获取员工信息
	con.Success(ctx, "ok", staffinfo)
}
