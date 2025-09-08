package sale

import (
	"jdy/controller/statistic"
	"jdy/logic/statistic/sale"
	"jdy/utils"

	"github.com/gin-gonic/gin"
)

type SaleController struct {
	statistic.StatisticController
}

func (con SaleController) Where(ctx *gin.Context) {
	where := utils.StructToWhere(sale.Where{})

	con.Success(ctx, "ok", where)
}
