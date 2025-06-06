package callback

import (
	"errors"
	"jdy/enums"
	"jdy/model"
	"log"
	"strings"

	"gorm.io/gorm"
)

type UserHandle struct {
	UserMessage
}

// 创建用户
func (l *EventChangeContactEvent) CreateUser() error {
	// 解析消息体
	var handler UserHandle
	if err := l.Handle.Event.ReadMessage(&handler.UserCreate); err != nil {
		return err
	}

	if handler.UserCreate.UserID == "" {
		return nil
	}

	var mobile *string
	if handler.UserCreate.Mobile == "" {
		mobile = nil
		log.Printf("%v,手机号为空", handler.UserCreate.UserID)
	} else {
		mobile = &handler.UserCreate.Mobile
	}

	var account model.Account
	if err := model.DB.Where(model.Account{
		Username: &handler.UserCreate.UserID,
		Platform: enums.PlatformTypeWxWork,
	}).Attrs(model.Account{
		Phone:    mobile,
		Nickname: &handler.UserCreate.Name,
		Avatar:   &handler.UserCreate.Avatar,
		Email:    &handler.UserCreate.Email,
		Gender:   enums.GenderUnknown.Convert(handler.UserCreate.Gender),
	}).FirstOrCreate(&account).Error; err != nil {
		return err
	}

	return nil
}

// 更新用户
func (l *EventChangeContactEvent) UpdateUser() error {
	// 解析消息体
	var handler UserHandle
	if err := l.Handle.Event.ReadMessage(&handler.UserUpdate); err != nil {
		return err
	}

	if handler.UserUpdate.UserID == "" {
		return nil
	}

	var account model.Account
	if err := model.DB.
		Where(model.Account{
			Username: &handler.UserUpdate.UserID,
			Platform: enums.PlatformTypeWxWork,
		}).
		Preload("Staff").
		First(&account).Error; err != nil {
		return err
	}

	uid := handler.UserUpdate.UserID
	if handler.UserUpdate.NewUserID != "" {
		uid = handler.UserUpdate.NewUserID
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 更新用户名
		if err := tx.Model(&account).Updates(model.Account{
			Username: &uid,
		}).Error; err != nil {
			return err
		}

		// 查询门店
		if account.Staff != nil {
			var stores []model.Store
			if err := tx.Where("id_wx in (?)", strings.Split(handler.UserUpdate.Department, ",")).Find(&stores).Error; err != nil {
				return err
			}

			if err := tx.Model(&account.Staff).Association("Stores").Replace(stores); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

// 删除用户
func (l *EventChangeContactEvent) DeleteUser() error {
	// 解析消息体
	var handler UserHandle
	if err := l.Handle.Event.ReadMessage(&handler.UserDelete); err != nil {
		return err
	}

	if handler.UserDelete.UserID == "" {
		return nil
	}

	var account model.Account
	if err := model.DB.Where(model.Account{
		Username: &handler.UserDelete.UserID,
		Platform: enums.PlatformTypeWxWork,
	}).Preload("Staff").First(&account).Error; err != nil {
		return errors.New(handler.UserDelete.UserID + "用户不存在")
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		if account.Staff != nil {
			if err := tx.Delete(&account.Staff).Error; err != nil {
				return errors.New(handler.UserDelete.UserID + "删除员工失败")
			}
			if err := tx.Where(model.Account{
				StaffId: account.StaffId,
			}).Delete(&model.Account{}).Error; err != nil {
				return errors.New(handler.UserDelete.UserID + "删除账号失败")
			}
		} else {
			if err := tx.Delete(&account).Error; err != nil {
				return errors.New(handler.UserDelete.UserID + "删除空账号失败")
			}
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
