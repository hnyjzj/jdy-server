package statistic

import (
	"jdy/logic/statistic"
	"jdy/types"

	"github.com/gin-gonic/gin"
)

// 销售明细日报
func (con StatisticController) SalesDetailDaily(ctx *gin.Context) {
	var (
		req   types.StatisticSalesDetailDailyReq
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

	res, err := logic.SalesDetailDaily(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", res)
}
