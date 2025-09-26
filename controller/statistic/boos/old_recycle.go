package boos

import (
	"jdy/logic/statistic/boos/old_recycle"
	"jdy/utils"

	"github.com/gin-gonic/gin"
)

func (con BoosController) OldRecycleWhere(ctx *gin.Context) {
	where := utils.StructToWhere(old_recycle.Where{})

	con.Success(ctx, "ok", where)
}

func (con BoosController) OldRecycleData(ctx *gin.Context) {
	var (
		req   old_recycle.DataReq
		logic = old_recycle.Logic{}
	)

	// 获取请求参数
	if err := ctx.ShouldBindJSON(&req); err != nil {
		con.Exception(ctx, "参数错误")
		return
	}

	// 获取当前登录用户
	if staff, err := con.GetStaff(ctx); err != nil {
		con.ExceptionWithAuth(ctx, err)
		return
	} else {
		logic.Staff = staff
		logic.Ctx = ctx
	}

	res, err := logic.GetDatas(&req)
	if err != nil {
		con.Exception(ctx, "获取失败")
		return
	}

	con.Success(ctx, "ok", res)
}
