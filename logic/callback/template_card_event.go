package callback

import (
	"errors"
	"log"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/work/server/handlers/models"
)

type TemplateCardEvent struct {
	Handle  *WxWork                        // 处理器
	Message *models.EventTemplateCardEvent // 消息体
}

func (l *TemplateCardEvent) GetEventKey() error {
	switch l.Handle.Event.GetEventKey() {
	case "":

	default:
		return errors.New("TemplateCardEvent.GetEventKey(): event key not found")
	}

	return nil
}

func (Handle *WxWork) TemplateCardEvent() any {
	var (
		l = TemplateCardEvent{
			Handle: Handle,
		}
	)

	// 获取员工信息
	if err := Handle.GetStaff(); err != nil {
		log.Printf("TemplateCardEvent.GetStaff.Error(): %v\n", err.Error())
		return "error"
	}

	// 解析消息体
	if err := l.Handle.Event.ReadMessage(&l.Message); err != nil {
		log.Printf("TemplateCardEvent.ReadMessage.Error(): %v\n", err.Error())
		return "error"
	}

	return nil
}
