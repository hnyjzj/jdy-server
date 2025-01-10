package staff

import (
	"jdy/controller"
	"jdy/errors"
	"jdy/logic/staff"
	"jdy/types"
	"jdy/utils"

	"github.com/gin-gonic/gin"
)

type StaffController struct {
	controller.BaseController
}

// 门店筛选条件
func (con StaffController) Where(ctx *gin.Context) {
	where := utils.StructToWhere(types.StaffWhere{})

	con.Success(ctx, "ok", where)
}

// 员工列表
func (con StaffController) List(ctx *gin.Context) {
	var (
		req   types.StaffListReq
		logic = &staff.StaffLogic{}
	)

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	logic.Ctx = ctx

	// 查询门店列表
	list, err := logic.List(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", list)
}

// 获取员工信息
func (con StaffController) Info(ctx *gin.Context) {
	var (
		logic = staff.StaffLogic{}
	)

	staff := con.GetStaff(ctx)

	staffinfo, err := logic.Info(ctx, &staff.Id)
	if err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	// 获取员工信息
	con.Success(ctx, "ok", staffinfo)
}
