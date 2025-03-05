package message

import (
	"context"
	"jdy/config"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/work"
)

type BaseMessage struct {
	Ctx context.Context

	WXWork *work.Work    `json:"wxwork"`
	App    *config.Agent `json:"app"`
}

func NewMessage(ctx context.Context) *BaseMessage {
	return &BaseMessage{
		Ctx:    ctx,
		WXWork: config.NewWechatService().JdyWork,
		App:    &config.Config.Wechat.Work.Jdy,
	}
}
