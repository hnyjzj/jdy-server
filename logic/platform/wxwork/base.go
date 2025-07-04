package wxwork

import (
	"github.com/ArtisanCloud/PowerWeChat/v3/src/work"
	"github.com/gin-gonic/gin"
)

type WxWorkLogic struct {
	Ctx *gin.Context

	App      *work.Work
	Contacts *work.Work
}
