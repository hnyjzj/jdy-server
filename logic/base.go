package logic

import (
	"jdy/types"

	"github.com/gin-gonic/gin"
)

type BaseLogic struct {
	Ctx   *gin.Context
	Staff *types.Staff
}
