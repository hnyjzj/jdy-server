package target

import (
	"jdy/model"

	"github.com/gin-gonic/gin"
)

type Logic struct {
	Ctx   *gin.Context
	Staff *model.Staff
}
