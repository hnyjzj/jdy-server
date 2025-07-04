package logic

import (
	"jdy/model"

	"github.com/gin-gonic/gin"
)

type BaseLogic struct {
	Ctx   *gin.Context
	Staff *model.Staff
}
