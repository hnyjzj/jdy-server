package boos

import (
	"jdy/controller/statistic"
	"jdy/logic/statistic/boos"
	"jdy/utils"

	"github.com/gin-gonic/gin"
)

type BoosController struct {
	statistic.StatisticController
}

func (con BoosController) BoosWhere(ctx *gin.Context) {
	where := utils.StructToWhere(boos.Where{})

	con.Success(ctx, "ok", where)
}
