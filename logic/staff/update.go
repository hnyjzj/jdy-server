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

// 更新信息
func (StaffLogic) StaffUpdate(ctx *gin.Context, uid string, req *types.StaffUpdateReq) error {
	// 更新逻辑
	l := &StaffUpdateLogic{
		ctx: ctx,
		uid: uid,
		req: req,
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		l.db = tx

		// 获取员工信息
		if err := l.getStaff(); err != nil {
			return err
		}

		// 更新员工信息
		if err := l.update(); err != nil {
			return err
		}

		// 退出登录
		if err := l.logout(); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

type StaffUpdateLogic struct {
	db  *gorm.DB
	ctx *gin.Context
	uid string
	req *types.StaffUpdateReq

	Staff *model.Staff
}

func (l *StaffUpdateLogic) getStaff() error {
	// 判断是否是微信登录
	if l.req.Code != "" {
		// 获取员工信息
		wxworklogic := wxwork.WxWorkLogic{
			Ctx: l.ctx,
		}
		staff, err := wxworklogic.OauthLogin(l.req.Code, true)
		if err != nil {
			return err
		}
		// 保存员工信息
		l.Staff = staff

		return nil
	} else {
		// 查询员工信息
		if err := l.db.Model(&model.Staff{}).Where("id = ?", l.uid).First(&l.Staff).Error; err != nil {
			return err
		}

		return nil
	}
}

func (l *StaffUpdateLogic) update() error {
	// 保存员工信息
	data := model.Staff{
		Nickname: l.req.Nickname,
		Avatar:   l.req.Avatar,
		Email:    l.req.Email,
		Gender:   l.req.Gender,
	}

	// 加密密码
	if l.req.Password != "" {
		password, err := model.Staff{}.HashPassword(&l.req.Password)
		if err != nil {
			return err
		}
		data.Password = password
	}

	if err := l.db.Model(&model.Staff{}).Where("id = ?", l.Staff.Id).Updates(&data).Error; err != nil { // 更新失败
		return err
	}

	return nil
}

func (l *StaffUpdateLogic) logout() error {
	// 退出登录
	if l.req.Password != "" && l.Staff.Phone != "" {
		// 退出登录
		auth := auth.LoginLogic{}
		if err := auth.Logout(l.ctx, l.Staff.Phone); err != nil {
			return errors.New("更新账号信息失败")
		}
	}

	return nil
}
