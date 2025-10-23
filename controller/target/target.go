package target

import (
	"jdy/controller"
	"jdy/errors"
	"jdy/logic/target"
	"jdy/types"
	"jdy/utils"

	"github.com/gin-gonic/gin"
)

type TargetController struct {
	controller.BaseController
}

// 筛选条件
func (con TargetController) Where(ctx *gin.Context) {
	where := utils.StructToWhere(types.TargetWhere{})

	con.Success(ctx, "ok", where)
}

// 销售目标分组筛选
func (con TargetController) WhereGroup(ctx *gin.Context) {
	where := utils.StructToWhere(types.TargetWhereGroup{})

	con.Success(ctx, "ok", where)
}

// 销售目标个人筛选
func (con TargetController) WherePersonal(ctx *gin.Context) {
	where := utils.StructToWhere(types.TargetWherePersonal{})

	con.Success(ctx, "ok", where)
}

// 列表
func (con TargetController) List(ctx *gin.Context) {
	var (
		req   types.TargetListReq
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

	// 获取列表
	list, err := logic.List(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", list)
}

// 详情
func (con TargetController) Info(ctx *gin.Context) {
	var (
		req   types.TargetInfoReq
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

	// 获取详情
	detail, err := logic.Info(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", detail)
}
