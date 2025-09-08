package stock

import (
	"jdy/controller/statistic"
	"jdy/logic/statistic/stock"
	"jdy/utils"

	"github.com/gin-gonic/gin"
)

type StockController struct {
	statistic.StatisticController
}

func (con StockController) Where(ctx *gin.Context) {
	where := utils.StructToWhere(stock.Where{})

	con.Success(ctx, "ok", where)
}
