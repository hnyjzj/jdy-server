package statistic

import (
	"jdy/logic/statistic"
	"jdy/utils"

	"github.com/gin-gonic/gin"
)

func (con StatisticController) ProductInventoryFinishedWhere(ctx *gin.Context) {
	where := utils.StructToWhere(statistic.ProductInventoryFinishedReq{})

	con.Success(ctx, "ok", where)
}

func (con StatisticController) ProductInventoryFinishedTitles(ctx *gin.Context) {

	var (
		logic = statistic.StatisticLogic{}
	)

	res := logic.ProductInventoryFinishedTitles()

	con.Success(ctx, "ok", res)
}

// 成品库存统计
func (con StatisticController) ProductInventoryFinishedData(ctx *gin.Context) {
	var (
		req   statistic.ProductInventoryFinishedReq
		logic = statistic.StatisticLogic{}
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

	res, err := logic.ProductInventoryFinishedData(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", res)
}
