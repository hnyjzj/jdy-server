package stock

import (
	"jdy/logic/statistic/stock"

	"github.com/gin-gonic/gin"
)

// 成品库存统计
func (con StockController) Data(ctx *gin.Context) {
	var (
		req   stock.DataReq
		logic = stock.StatisticStockLogic{}
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

	res, err := logic.Data(&req)
	if err != nil {
		con.Exception(ctx, "获取失败")
		return
	}

	con.Success(ctx, "ok", res)
}
