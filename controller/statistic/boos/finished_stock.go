package boos

import (
	"jdy/logic/statistic/boos/finished_stock"
	"jdy/utils"

	"github.com/gin-gonic/gin"
)

func (con BoosController) FinishedStockWhere(ctx *gin.Context) {
	where := utils.StructToWhere(finished_stock.Where{})

	con.Success(ctx, "ok", where)
}

// 成品库存统计
func (con BoosController) FinishedStockData(ctx *gin.Context) {
	var (
		req   finished_stock.DataReq
		logic = finished_stock.Logic{}
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
		logic.Ctx = ctx
	}

	res, err := logic.GetDatas(&req)
	if err != nil {
		con.Exception(ctx, "获取失败")
		return
	}

	con.Success(ctx, "ok", res)
}
