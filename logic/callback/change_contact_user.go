package callback

import (
	"errors"
	"jdy/enums"
	"jdy/model"
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

	// 获取用户信息
	userinfo, err := l.Handle.GetUser(handler.UserCreate.UserID)
	if err != nil {
		return err
	}

	var mobile *string
	if userinfo.Mobile == "" {
		mobile = nil
	} else {
		mobile = &userinfo.Mobile
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {

		var account model.Account
		if err := tx.Where(model.Account{
			Username: &userinfo.UserID,
			Platform: enums.PlatformTypeWxWork,
		}).Attrs(model.Account{
			Phone:    mobile,
			Nickname: &userinfo.Name,
			Avatar:   &userinfo.Avatar,
			Email:    &userinfo.Email,
			Gender:   enums.GenderUnknown.Convert(userinfo.Gender),
		}).FirstOrCreate(&account).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
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
		return errors.New("用户不存在：" + handler.UserDelete.UserID)
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		if account.Staff != nil {
			if err := tx.Delete(&account.Staff).Error; err != nil {
				return errors.New("删除员工失败：" + handler.UserDelete.UserID)
			}
			if err := tx.Where(model.Account{
				StaffId: account.StaffId,
			}).Delete(&model.Account{}).Error; err != nil {
				return errors.New("删除账号失败：" + handler.UserDelete.UserID)
			}
		} else {
			if err := tx.Delete(&account).Error; err != nil {
				return errors.New("删除空账号失败：" + handler.UserDelete.UserID)
			}
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
