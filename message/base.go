package message

import (
	"jdy/config"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/work"
	"github.com/gin-gonic/gin"
)

type BaseMessage struct {
	Ctx *gin.Context

	WXWork *work.Work    `json:"wxwork"`
	App    *config.Agent `json:"app"`
}

func NewMessage(ctx *gin.Context) *BaseMessage {
	return &BaseMessage{
		Ctx:    ctx,
		WXWork: config.NewWechatService().JdyWork,
		App:    &config.Config.Wechat.Work.Jdy,
	}
}
