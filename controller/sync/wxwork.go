package sync

import (
	"jdy/logic/sync"

	"github.com/gin-gonic/gin"
)

type WxworkController struct {
	SyncController
}

// 同步通讯录
func (con WxworkController) SyncContacts(ctx *gin.Context) {
	logic := sync.WxWorkLogic{}

	logic.Ctx = ctx

	if err := logic.Contacts(); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}
