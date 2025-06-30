package callback

import (
	"errors"
	"jdy/enums"
	"jdy/logic/platform/wxwork"
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
	wxlogic := &wxwork.WxWorkLogic{Ctx: l.Handle.Ctx}
	userinfo, err := wxlogic.GetUser(handler.UserCreate.UserID)
	if err != nil {
		return err
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		var staff model.Staff
		if err := tx.Where(model.Staff{
			Username: userinfo.Username,
		}).Attrs(model.Staff{
			Phone:    userinfo.Phone,
			Nickname: userinfo.Nickname,
			Avatar:   userinfo.Avatar,
			Email:    userinfo.Email,
			Gender:   enums.GenderUnknown.Convert(userinfo.Gender),
		}).FirstOrCreate(&staff).Error; err != nil {
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

	var staff model.Staff
	if err := model.DB.
		Where(model.Staff{
			Username: handler.UserUpdate.UserID,
		}).
		First(&staff).Error; err != nil {
		return errors.New("用户不存在：" + handler.UserUpdate.UserID)
	}

	uid := handler.UserUpdate.UserID
	if handler.UserUpdate.NewUserID != "" {
		uid = handler.UserUpdate.NewUserID
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 更新用户名
		if err := tx.Model(&staff).Updates(model.Staff{
			Username: uid,
		}).Error; err != nil {
			return err
		}

		// 查询门店
		var stores []model.Store
		if err := tx.Where("id_wx in (?)", strings.Split(handler.UserUpdate.Department, ",")).Find(&stores).Error; err != nil {
			return err
		}

		if err := tx.Model(&staff).Association("Stores").Replace(stores); err != nil {
			return err
		}

		var regions []model.Region
		if err := tx.Where("id_wx in (?)", strings.Split(handler.UserUpdate.Department, ",")).Find(&regions).Error; err != nil {
			return err
		}
		if err := tx.Model(&staff).Association("Regions").Replace(regions); err != nil {
			return err
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

	var staff model.Staff
	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 查询员工
		if err := tx.Where(model.Staff{
			Username: handler.UserDelete.UserID,
		}).First(&staff).Error; err != nil {
			return errors.New("用户不存在：" + handler.UserDelete.UserID)
		}
		// 删除员工
		if err := tx.Delete(&staff).Error; err != nil {
			return errors.New("删除员工失败：" + handler.UserDelete.UserID)
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
