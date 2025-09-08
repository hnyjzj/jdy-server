package payment

import (
	"jdy/controller/statistic"
	"jdy/logic/statistic/payment"
	"jdy/utils"

	"github.com/gin-gonic/gin"
)

type PaymentController struct {
	statistic.StatisticController
}

func (con PaymentController) Where(ctx *gin.Context) {
	where := utils.StructToWhere(payment.Where{})

	con.Success(ctx, "ok", where)
}
