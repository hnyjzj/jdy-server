package callback

import (
	"errors"
	"jdy/config"
	"jdy/enums"
	"jdy/message"
	"jdy/model"
	"log"

	models1 "github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/models"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/work/server/handlers/models"
	"gorm.io/gorm"
)

// 分配事件
func (l *EventChangeExternalContact) Distribute() error {
	switch l.Message.ChangeType {
	case models.CALLBACK_EVENT_CHANGE_TYPE_ADD_EXTERNAL_CONTACT:
		if err := l.GetExternalContact(); err != nil {
			log.Printf("TemplateCardEvent.GetExternalContact.Error(): %v\n", err.Error())
			return err
		}

	default:
		return errors.New("ChangeExternalContactEvent.ChangeType not supported")
	}

	return nil
}

// 外部联系人事件
type EventChangeExternalContact struct {
	Handle  *WxWork                       // 处理器
	Message models1.CallbackMessageHeader // 消息体

	ExternalUserAdd models.EventExternalUserAdd // 外部联系人添加
}

// 获取外部联系人
func (l *EventChangeExternalContact) GetExternalContact() error {
	// 解析消息体
	if err := l.Handle.Event.ReadMessage(&l.ExternalUserAdd); err != nil {
		return err
	}

	// 获取外部联系人信息
	app := config.NewWechatService().JdyWork
	user, err := app.ExternalContact.Get(l.Handle.Ctx, l.ExternalUserAdd.ExternalUserID, "")
	if err != nil || user == nil || user.ExternalContact == nil {
		return errors.New("获取外部联系人信息失败: " + err.Error())
	}

	// 查找员工
	var account model.Account
	if err := model.DB.Where(model.Account{Username: &l.ExternalUserAdd.UserID, Platform: enums.PlatformTypeWxWork}).Preload("Staff").First(&account).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return errors.New("查询会员失败: " + err.Error())
		}
	}
	if account.Staff == nil {
		account.Staff = &model.Staff{}
		account.Username = &l.ExternalUserAdd.UserID
	}

	// 查找会员
	var member model.Member
	var gender enums.Gender
	if err := model.DB.Where(model.Member{ExternalUserId: l.ExternalUserAdd.ExternalUserID}).Attrs(model.Member{
		Name:           user.ExternalContact.Name,
		Gender:         gender.Convert(user.ExternalContact.Gender),
		Nickname:       user.ExternalContact.Name,
		Level:          enums.MemberLevelNone,
		SourceId:       account.Staff.Id,
		ConsultantId:   account.Staff.Id,
		Status:         enums.MemberStatusPending,
		ExternalUserId: l.ExternalUserAdd.ExternalUserID,
	}).FirstOrCreate(&member).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return errors.New("查询会员失败: " + err.Error())
		}
	}

	// 发送消息
	m := message.NewMessage(l.Handle.Ctx)
	m.SendMemberCreateMessage(&message.MemberCreateMessage{
		ToUser:         *account.Username,
		ExternalUserID: l.ExternalUserAdd.ExternalUserID,
		Name:           user.ExternalContact.Name,
	})

	return nil
}

// 外部联系人事件处理
func (Handle *WxWork) ChangeExternalContactEvent() any {
	var (
		l = EventChangeExternalContact{
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
