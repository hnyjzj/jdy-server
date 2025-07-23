package callback

import (
	"errors"
	"fmt"
	"jdy/enums"
	"jdy/logic/platform/wxwork"
	"jdy/model"
	"log"

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
			Identity: enums.IdentityClerk,
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

	// 获取用户信息
	user, err := l.Handle.Wechat.JdyWork.User.Get(l.Handle.Ctx, handler.UserUpdate.UserID)
	if err != nil || user.UserID == "" {
		log.Printf("读取员工信息失败: %+v, %+v", err, user)
		return errors.New("读取员工信息失败")
	}

	// 查询员工
	staff, err := model.Staff{}.Get(nil, &user.UserID)
	if err != nil {
		return errors.New("用户不存在：" + user.UserID)
	}

	uid := handler.UserUpdate.UserID
	if handler.UserUpdate.NewUserID != "" {
		uid = handler.UserUpdate.NewUserID
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 更新用户名
		if err := tx.Model(&model.Staff{}).Where("id = ?", staff.Id).Updates(model.Staff{
			Username: uid,
			Identity: enums.IdentityClerk,
		}).Error; err != nil {
			return err
		}

		// 关联门店
		var stores []model.Store
		if err := tx.Where("id_wx in (?)", user.Department).Find(&stores).Error; err != nil {
			return err
		}
		if err := tx.Model(&staff).Association("Stores").Replace(stores); err != nil {
			return err
		}

		// 关联区域
		var regions []model.Region
		if err := tx.Where("id_wx in (?)", user.Department).Find(&regions).Error; err != nil {
			return err
		}
		if err := tx.Model(&staff).Association("Regions").Replace(regions); err != nil {
			return err
		}

		// 关联门店负责人
		var StoreSuperiorsIds []string
		for i := range user.Department {
			if user.IsLeaderInDept[i] == 1 {
				StoreSuperiorsIds = append(StoreSuperiorsIds, fmt.Sprint(user.Department[i]))
			}
		}
		var StoreSuperiors []model.Store
		if err := tx.Where("id_wx in (?)", StoreSuperiorsIds).Find(&StoreSuperiors).Error; err != nil {
			return err
		}
		if err := tx.Model(&staff).Association("StoreSuperiors").Replace(StoreSuperiors); err != nil {
			return err
		}
		// 关联区域负责人
		var RegionSuperiors []model.Region
		if err := tx.Where("id_wx in (?)", StoreSuperiorsIds).Find(&RegionSuperiors).Error; err != nil {
			return err
		}
		if err := tx.Model(&staff).Association("RegionSuperiors").Replace(RegionSuperiors); err != nil {
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
