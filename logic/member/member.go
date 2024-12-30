package member

import (
	"jdy/types"

	"github.com/gin-gonic/gin"
)

type MemberLogic struct {
	Ctx   *gin.Context
	Staff *types.Staff
}
