package boos

import (
	"jdy/logic/statistic/boos/finished_sales"
	"jdy/utils"

	"github.com/gin-gonic/gin"
)

func (con BoosController) FinishedSalesWhere(ctx *gin.Context) {
	where := utils.StructToWhere(finished_sales.Where{})

	con.Success(ctx, "ok", where)
}

func (con BoosController) FinishedSalesData(ctx *gin.Context) {
	var (
		req   finished_sales.DataReq
		logic = finished_sales.Logic{}
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
		con.Exception(ctx, "获取失败")
		return
	}

	con.Success(ctx, "ok", res)
}
