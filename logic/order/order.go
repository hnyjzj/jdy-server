package order

import (
	"jdy/types"

	"github.com/gin-gonic/gin"
)

type OrderLogic struct {
	Ctx   *gin.Context
	Staff *types.Staff
}
