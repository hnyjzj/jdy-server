package callback

import (
	"fmt"
	"log"

	models1 "github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/models"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/work/server/handlers/models"
)

type UserMessage struct {
	UserCreate models.EventUserCreate // 新增成员
	UserUpdate models.EventUserUpdate // 更新成员
	UserDelete models.EventUserDelete // 删除成员
}

type PartyMessage struct {
	PartyCreate models.EventPartyCreate // 新增部门
	PartyUpdate models.EventPartyUpdate // 更新部门
	PartyDelete models.EventPartyDelete // 删除部门
}

// 分配
func (l *EventChangeContactEvent) Distribute() error {
	switch l.Message.ChangeType {
	case models.CALLBACK_EVENT_CHANGE_TYPE_CREATE_USER: // 新增成员
		return l.CreateUser()
	case models.CALLBACK_EVENT_CHANGE_TYPE_UPDATE_USER: // 更新成员
		return l.UpdateUser()
	case models.CALLBACK_EVENT_CHANGE_TYPE_DELETE_USER: // 删除成员
		return l.DeleteUser()
	case models.CALLBACK_EVENT_CHANGE_TYPE_CREATE_PARTY: // 新增部门
		return l.CreateParty()
	case models.CALLBACK_EVENT_CHANGE_TYPE_DELETE_PARTY: // 删除部门
		return l.DeleteParty()
	default:
		err := fmt.Errorf("不支持更改类型(%v)", l.Message.ChangeType)
		log.Printf(err.Error()+": %+v", l.Message)
		return err
	}
}

// 员工变更事件
type EventChangeContactEvent struct {
	Handle  *WxWork                       // 处理器
	Message models1.CallbackMessageHeader // 消息体
}

// 员工变更事件处理
func (Handle *WxWork) ChangeContactEvent() any {
	var (
		l = EventChangeContactEvent{
			Handle: Handle,
		}
	)
	// 解析消息体
	if err := l.Handle.Event.ReadMessage(&l.Message); err != nil {
		log.Printf("TemplateCardEvent.ReadMessage.Error(): %v\n", err.Error())
		return "error"
	}
	// 处理事件
	if err := l.Distribute(); err != nil {
		log.Printf("TemplateCardEvent.GetEventKey.Error(): %v\n", err.Error())
		return "error"
	}

	return nil
}
