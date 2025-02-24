package statistic

import (
	"jdy/controller"
	"jdy/logic/statistic"
	"jdy/types"

	"github.com/gin-gonic/gin"
)

type StatisticController struct {
	controller.BaseController
}

func (con StatisticController) Total(ctx *gin.Context) {
	var (
		logic = statistic.StatisticLogic{}
	)

	res, err := logic.Total()
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", res)
}

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

	res, err := logic.TodaySales(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", res)
}
