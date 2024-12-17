package staff

import (
	"jdy/errors"
	"jdy/logic/auth"
	"jdy/logic/platform/wxwork"
	"jdy/model"
	"jdy/types"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type StaffUpdateLogic struct {
	ctx *gin.Context
	uid string
	req *types.StaffUpdateReq
}

// 更新信息
func (StaffLogic) StaffUpdate(ctx *gin.Context, uid string, req *types.StaffUpdateReq) *errors.Errors {
	// 更新逻辑
	l := &StaffUpdateLogic{
		ctx: ctx,
		uid: uid,
		req: req,
	}
	// 判断平台
	switch req.Platform {
	case types.PlatformTypeAccount:
		if err := l.account(); err != nil {
			return errors.New(err.Error())
		}
	case types.PlatformTypeWxWork:

		if err := l.wxwork(); err != nil {
			return errors.New(err.Error())
		}
	default:
		return errors.New("平台类型错误")
	}

	return nil
}

// 更新账号信息
func (l *StaffUpdateLogic) account() error {
	var (
		req = l.req.Account
	)

	var staff model.Staff
	if err := model.DB.
		Preload("Account", func(db *gorm.DB) *gorm.DB {
			return db.Where(&model.Account{Platform: types.PlatformTypeAccount})
		}).
		Preload("Accounts").
		First(&staff, l.uid).Error; err != nil {
		return err
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 更新员工信息
		if err := tx.Model(&staff).Updates(model.Staff{
			Nickname: req.Nickname,
			Avatar:   req.Avatar,
			Email:    req.Email,
			Gender:   req.Gender,
		}).Error; err != nil {
			return err
		}

		// 更新账号信息
		account := model.Account{
			Platform: types.PlatformTypeAccount,
			Phone:    staff.Phone,

			Nickname: &req.Nickname,
			Avatar:   &req.Avatar,
			Email:    &req.Email,
			Gender:   req.Gender,
		}

		// 加密密码
		if req.Password != "" {
			password, err := model.Account{}.HashPassword(&req.Password)
			if err != nil {
				return err
			}
			account.Password = &password
		}

		if staff.Account == nil {
			staff.Account = &account
			// 创建账号
			if err := tx.Save(&staff).Error; err != nil {
				return err
			}
		} else {
			// 更新账号信息
			if err := tx.Model(&staff.Account).Updates(account).Error; err != nil {
				return err
			}
		}

		if req.Password != "" && staff.Account.Phone != nil {
			// 退出登录
			auth := auth.LoginLogic{}
			if err := auth.Logout(l.ctx, *staff.Account.Phone); err != nil {
				return errors.New("更新账号信息失败")
			}
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

// 更新企业微信信息
func (l *StaffUpdateLogic) wxwork() error {
	var (
		wxwork = wxwork.WxWorkLogic{}
	)

	staff, err := wxwork.OauthLogin(l.ctx, l.req.WxWork.Code, true)
	if err != nil {
		return err
	}

	if err := model.DB.Model(&model.Staff{}).
		Preload("Account", func(db *gorm.DB) *gorm.DB {
			return db.Where(&model.Account{Platform: types.PlatformTypeWxWork})
		}).
		First(&staff, l.uid).Error; err != nil {
		return err
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&staff).Updates(model.Staff{
			Phone:    staff.Account.Phone,
			Nickname: *staff.Account.Nickname,
			Avatar:   *staff.Account.Avatar,
			Email:    *staff.Account.Email,
			Gender:   staff.Account.Gender,
		}).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
