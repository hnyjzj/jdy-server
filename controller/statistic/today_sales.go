package statistic

import (
	"jdy/logic/statistic"
	"jdy/types"

	"github.com/gin-gonic/gin"
)

// 今日销售统计
func (con StatisticController) TodaySales(ctx *gin.Context) {
	var (
		req   types.StatisticTodaySalesReq
		logic = statistic.StatisticLogic{}
	)

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
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

	res, err := logic.TodaySales(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", res)
}
