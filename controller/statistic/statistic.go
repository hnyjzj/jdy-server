package statistic

import (
	"jdy/controller"
	"jdy/logic/statistic"

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
