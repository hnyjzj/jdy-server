package product

import (
	"jdy/types"

	"github.com/gin-gonic/gin"
)

type ProductLogic struct {
	Ctx   *gin.Context
	Staff *types.Staff
}
