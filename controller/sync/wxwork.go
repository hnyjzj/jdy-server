package sync

import (
	"jdy/logic/sync"
	"jdy/message"

	"github.com/gin-gonic/gin"
)

type WxworkController struct {
	SyncController
}

// 同步通讯录
func (con WxworkController) SyncContacts(ctx *gin.Context) {
	logic := sync.WxWorkLogic{}

	logic.Ctx = ctx

	go func(logic sync.WxWorkLogic) {
		data := &message.CaptureScreenMessage{
			Type:      "同步通讯录",
			Username:  "系统",
			Storename: "系统",
			Url:       "no_url",
		}

		if err := logic.Contacts(); err != nil {
			data.Title = "同步失败"
			data.Desc = err.Error()
		} else {
			data.Title = "同步成功"
		}

		_ = message.NewMessage(ctx).SendCaptureScreenMessage(data)

	}(logic)

	con.Success(ctx, "通讯录同步中", nil)
}
