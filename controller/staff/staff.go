package staff

import (
	"jdy/controller"
	"jdy/errors"
	"jdy/logic"
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

// 员工详情
func (con StaffController) Info(ctx *gin.Context) {
	var (
		req   types.StaffInfoReq
		logic = staff.StaffLogic{
			Base: logic.Base{
				Ctx: ctx,
			},
			Staff: con.GetStaff(ctx),
		}
	)

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	staffinfo, err := logic.Info(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", staffinfo)

}

// 获取我的员工信息
func (con StaffController) My(ctx *gin.Context) {
	var (
		logic = staff.StaffLogic{
			Base: logic.Base{
				Ctx: ctx,
			},
			Staff: con.GetStaff(ctx),
		}
	)

	staffinfo, err := logic.My()
	if err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	con.Success(ctx, "ok", staffinfo)
}
