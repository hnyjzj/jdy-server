package callback

import (
	"log"

	models1 "github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/models"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/work/server/handlers/models"
)

// 分配
func (l *EventChangeContactEvent) Distribute() error {
	switch l.Message.ChangeType {
	case models.CALLBACK_EVENT_CHANGE_TYPE_CREATE_USER: // 新增成员
		return l.CreateUser(l.Message)
	case models.CALLBACK_EVENT_CHANGE_TYPE_UPDATE_USER: // 更新成员
		return l.UpdateUser(l.Message)
	case models.CALLBACK_EVENT_CHANGE_TYPE_DELETE_USER: // 删除成员
		return l.DeleteUser(l.Message)
	}
	return nil
}

// 员工变更事件
type EventChangeContactEvent struct {
	Handle  *WxWork                        // 处理器
	Message *models1.CallbackMessageHeader // 消息体

	UserCreate models.EventUserCreate // 新增成员
	UserUpdate models.EventUserUpdate // 更新成员
	UserDelete models.EventUserDelete // 删除成员
}

// 创建用户
func (l *EventChangeContactEvent) CreateUser(message *models1.CallbackMessageHeader) error {
	// 解析消息体
	if err := l.Handle.Event.ReadMessage(&l.UserCreate); err != nil {
		return err
	}

	return nil
}

// 更新用户
func (l *EventChangeContactEvent) UpdateUser(message *models1.CallbackMessageHeader) error {
	// 解析消息体
	if err := l.Handle.Event.ReadMessage(&l.UserUpdate); err != nil {
		return err
	}

	return nil
}

// 删除用户
func (l *EventChangeContactEvent) DeleteUser(message *models1.CallbackMessageHeader) error {
	// 解析消息体
	if err := l.Handle.Event.ReadMessage(&l.UserDelete); err != nil {
		return err
	}

	return nil
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
