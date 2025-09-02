package boos

import (
	"jdy/logic/statistic/boos/performance"
	"jdy/utils"

	"github.com/gin-gonic/gin"
)

func (con BoosController) PerformanceWhere(ctx *gin.Context) {
	where := utils.StructToWhere(performance.Where{})

	con.Success(ctx, "ok", where)
}

func (con BoosController) PerformanceData(ctx *gin.Context) {
	var (
		req   performance.DataReq
		logic = performance.Logic{}
	)

	// 获取请求参数
	if err := ctx.ShouldBindJSON(&req); err != nil {
		con.Exception(ctx, "参数错误")
		return
	}

	// 获取当前登录用户
	if staff, err := con.GetStaff(ctx); err != nil {
		con.Exception(ctx, "无法获取")
		return
	} else {
		logic.Staff = staff
	}

	res, err := logic.GetDatas(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", res)
}
