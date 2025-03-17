package callback

import (
	"errors"
	"jdy/config"
	"jdy/enums"
	"jdy/message"
	"jdy/model"
	"log"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/work/server/handlers/models"
	"gorm.io/gorm"
)

type EventChangeExternalContact struct {
	Handle  *WxWork                      // 处理器
	Message *models.EventExternalUserAdd // 消息体
}

func (Handle *WxWork) ChangeExternalContactEvent() any {
	var (
		l = EventChangeExternalContact{
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

	// 处理事件
	if err := l.Distribute(); err != nil {
		log.Printf("TemplateCardEvent.GetEventKey.Error(): %v\n", err.Error())
		return "error"
	}

	return nil
}

func (l *EventChangeExternalContact) Distribute() error {
	switch l.Handle.Event.GetChangeType() {
	case models.CALLBACK_EVENT_CHANGE_TYPE_ADD_EXTERNAL_CONTACT:
		if err := l.GetExternalContact(); err != nil {
			log.Printf("TemplateCardEvent.GetExternalContact.Error(): %v\n", err.Error())
			return err
		}

	default:
		return errors.New("TemplateCardEvent.GetEventKey(): event key not found")
	}

	return nil
}

func (l *EventChangeExternalContact) GetExternalContact() error {
	// 获取外部联系人信息
	var msg models.EventExternalUserAdd
	if err := l.Handle.Event.ReadMessage(&msg); err != nil {
		return errors.New("解析消息体失败: " + err.Error())
	}

	// 查找员工
	var account model.Account
	if err := model.DB.Where(model.Account{Username: &msg.UserID, Platform: enums.PlatformTypeWxWork}).Preload("Staff").First(&account).Error; err != nil {
		return errors.New("查询员工失败: " + err.Error())
	}
	if account.Id == "" || account.Staff == nil {
		return errors.New("员工不存在")
	}

	// 获取外部联系人信息
	app := config.NewWechatService().JdyWork
	user, err := app.ExternalContact.Get(l.Handle.Ctx, msg.ExternalUserID, "")
	if err != nil {
		return errors.New("获取外部联系人信息失败: " + err.Error())
	}

	// 查找会员
	var member model.Member
	if err := model.DB.Where(model.Member{ExternalUserID: msg.ExternalUserID}).Attrs(model.Member{
		Name:           user.ExternalContact.Name,
		Gender:         enums.GenderUnknown.Convert(user.ExternalContact.Gender),
		Nickname:       user.ExternalContact.Name,
		Level:          enums.MemberLevelNone,
		SourceId:       account.Staff.Id,
		ConsultantId:   account.Staff.Id,
		Status:         enums.MemberStatusPending,
		ExternalUserID: msg.ExternalUserID,
	}).FirstOrCreate(&member).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return errors.New("查询会员失败: " + err.Error())
		}
	}

	// 发送消息
	m := message.NewMessage(l.Handle.Ctx)
	m.SendMemberCreateMessage(&message.MemberCreateMessage{
		ToUser:         *account.Username,
		ExternalUserID: msg.ExternalUserID,
		Name:           user.ExternalContact.Name,
	})

	return nil
}
