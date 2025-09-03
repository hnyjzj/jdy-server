package statistic

import (
	"jdy/model"

	"github.com/gin-gonic/gin"
)

type StatisticLogic struct {
	Ctx   *gin.Context
	Staff *model.Staff
}
