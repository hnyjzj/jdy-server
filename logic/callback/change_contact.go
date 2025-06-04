package callback

import (
	"errors"
	"fmt"
	"jdy/enums"
	"jdy/model"
	"log"

	models1 "github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/models"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/work/server/handlers/models"
	"gorm.io/gorm"
)

// 分配
func (l *EventChangeContactEvent) Distribute() error {
	switch l.Message.ChangeType {
	case models.CALLBACK_EVENT_CHANGE_TYPE_CREATE_USER: // 新增成员
		return l.CreateUser(&l.Message)
	case models.CALLBACK_EVENT_CHANGE_TYPE_UPDATE_USER: // 更新成员
		return l.UpdateUser(&l.Message)
	case models.CALLBACK_EVENT_CHANGE_TYPE_DELETE_USER: // 删除成员
		return l.DeleteUser(&l.Message)
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

	if l.UserCreate.UserID == "" {
		return nil
	}

	var mobile *string
	if l.UserCreate.Mobile == "" {
		mobile = nil
		log.Printf("%v,手机号为空", l.UserCreate.UserID)
	} else {
		mobile = &l.UserCreate.Mobile
	}

	var account model.Account
	if err := model.DB.Where(model.Account{
		Username: &l.UserCreate.UserID,
		Platform: enums.PlatformTypeWxWork,
	}).Attrs(model.Account{
		Phone:    mobile,
		Nickname: &l.UserCreate.Name,
		Avatar:   &l.UserCreate.Avatar,
		Email:    &l.UserCreate.Email,
		Gender:   enums.GenderUnknown.Convert(l.UserCreate.Gender),
		Info:     &l.UserCreate,
	}).FirstOrCreate(&account).Error; err != nil {
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

	if l.UserUpdate.UserID == "" {
		return nil
	}

	if l.UserUpdate.Mobile == "" {
		log.Printf("%v,手机号为空", l.UserUpdate.UserID)
	}

	var account model.Account
	if err := model.DB.Where(model.Account{
		Username: &l.UserUpdate.UserID,
		Platform: enums.PlatformTypeWxWork,
	}).First(&account).Error; err != nil {
		return err
	}

	uid := l.UserUpdate.UserID
	if l.UserUpdate.NewUserID != "" {
		uid = l.UserUpdate.NewUserID
	}

	if err := model.DB.Model(&account).Updates(model.Account{
		Username: &uid,
		Info:     &l.UserUpdate,
	}).Error; err != nil {
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

	if l.UserDelete.UserID == "" {
		return nil
	}

	var account model.Account
	if err := model.DB.Where(model.Account{
		Username: &l.UserDelete.UserID,
		Platform: enums.PlatformTypeWxWork,
	}).Preload("Staff").First(&account).Error; err != nil {
		return errors.New(l.UserDelete.UserID + "用户不存在")
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		if account.Staff != nil {
			if err := tx.Delete(&account.Staff).Error; err != nil {
				return errors.New(l.UserDelete.UserID + "删除员工失败")
			}
			if err := tx.Where(model.Account{
				StaffId: account.StaffId,
			}).Delete(&model.Account{}).Error; err != nil {
				return errors.New(l.UserDelete.UserID + "删除账号失败")
			}
		} else {
			if err := tx.Delete(&account).Error; err != nil {
				return errors.New(l.UserDelete.UserID + "删除空账号失败")
			}
		}

		return nil
	}); err != nil {
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
