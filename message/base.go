package message

import (
	"context"
	"errors"
	"jdy/config"
	"log"

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

func (m *BaseMessage) Send(WXWork *work.Work, messages any) error {
	if res, err := WXWork.Message.Send(m.Ctx, messages); err != nil || res.ErrCode != 0 {
		log.Printf("res: %+v\n", res)
		log.Printf("err: %+v\n", err)
		log.Printf("messages: %+v\n", messages)

		return errors.New(res.ErrMsg)
	}

	return nil
}
