package wxwork

import (
	"jdy/config"
	"jdy/enums"
	"jdy/errors"
	"jdy/model"
	"log"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/work"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type wxworkLoginLogic struct {
	Staff    *model.Staff
	UserInfo *model.Staff

	App *work.Work

	Ctx         *gin.Context
	Db          *gorm.DB
	Departments []int
}

// 扫码登录
func (w *WxWorkLogic) CodeLogin(code string) (*model.Staff, error) {
	l := &wxworkLoginLogic{
		Ctx: w.Ctx,
		App: config.NewWechatService().JdyWork,
	}
	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		l.Db = tx

		// 获取用户信息
		if err := l.getCodeUserInfo(code); err != nil {
			return err
		}

		// 获取员工信息
		if err := l.getStaff(); err != nil {
			return err
		}

		// 判断是否已注册
		if l.Staff.Phone == nil {
			return errors.New("首次登录需通过企业微信工作台打开并授权手机号")
		}

		// 更新登录信息
		if err := l.updateLogin(); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return l.Staff, nil
}

// 授权登录
func (w *WxWorkLogic) OauthLogin(code string, isRegister bool) (*model.Staff, error) {
	l := &wxworkLoginLogic{
		Ctx: w.Ctx,
		App: config.NewWechatService().JdyWork,
	}
	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		l.Db = tx

		// 获取用户信息
		if err := l.getOathUserInfo(code); err != nil {
			return err
		}

		// 获取员工信息
		if err := l.getStaff(); err != nil {
			return err
		}

		// 判断是否需要注册
		if err := l.register(); err != nil {
			return err
		}

		// 更新账号
		if err := l.updateAccount(); err != nil {
			return err
		}

		// 更新登录信息
		if err := l.updateLogin(); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return l.Staff, nil
}

func (l *wxworkLoginLogic) getCodeUserInfo(code string) error {
	// 获取用户信息
	user, err := l.App.OAuth.Provider.GetUserInfo(code)
	if err != nil || user.UserID == "" {
		return errors.New("获取企业微信用户信息失败")
	}

	l.UserInfo = &model.Staff{
		Username: &user.UserID,
	}

	return nil
}

func (l *wxworkLoginLogic) getOathUserInfo(code string) error {
	// 获取用户信息
	user, err := l.App.OAuth.Provider.GetUserInfo(code)
	if err != nil || user.UserID == "" || user.UserTicket == "" {
		log.Printf("获取企业微信用户信息失败: %+v, %+v", err, user)
		return errors.New("获取企业微信用户信息失败")
	}
	// 读取员工信息
	userinfo, err := l.App.User.Get(l.Ctx, user.UserID)
	if err != nil || userinfo.UserID == "" {
		log.Printf("读取员工信息失败: %+v, %+v", err, userinfo)
		return errors.New("读取员工信息失败")
	}
	// 获取用户详情
	detail, err := l.App.OAuth.Provider.GetUserDetail(user.UserTicket)
	if err != nil || detail.Mobile == "" {
		log.Printf("获取企业微信用户详情失败: %+v, %+v", err, detail)
		return errors.New("获取企业微信用户详情失败")
	}

	// 获取性别
	var gender enums.Gender

	l.UserInfo = &model.Staff{
		Phone:    &detail.Mobile,
		Username: &user.UserID,
		Nickname: userinfo.Name,
		Avatar:   detail.Avatar,
		Email:    detail.Email,
		Gender:   gender.Convert(detail.Gender),
	}
	l.Departments = userinfo.Department

	return nil
}

// 更新账号
func (l *wxworkLoginLogic) updateAccount() error {
	// 更新账号信息
	if err := l.Db.Model(&l.Staff).Updates(&model.Staff{
		Phone:    l.UserInfo.Phone,
		Nickname: l.UserInfo.Nickname,
		Avatar:   l.UserInfo.Avatar,
		Email:    l.UserInfo.Email,
		Gender:   l.UserInfo.Gender,
	}).Error; err != nil {
		return errors.New("更新账号失败")
	}

	return nil
}

func (l *wxworkLoginLogic) register() error {
	// 账号没有手机号
	if l.Staff.Phone == nil {
		l.Staff.Phone = l.UserInfo.Phone
	}

	// 手机号未授权
	if l.UserInfo.Phone == nil {
		return errors.New("手机号未授权")
	}

	// 手机号不一致
	if l.Staff.Phone != nil && *l.Staff.Phone != *l.UserInfo.Phone {
		return errors.New("手机号不一致")
	}

	return nil
}

// 获取员工
func (l *wxworkLoginLogic) getStaff() error {
	// 查询账号
	if err := l.Db.Where(&model.Staff{
		Username: l.UserInfo.Username,
	}).First(&l.Staff).Error; err != nil {
		return errors.ErrStaffNotFound
	}

	// 判断是否被禁用
	if l.Staff.IsDisabled {
		return errors.ErrStaffDisabled
	}

	return nil
}

// 更新登录信息
func (l *wxworkLoginLogic) updateLogin() error {
	// 更新账号
	l.Staff.UpdateLoginInfo(l.Ctx.ClientIP())

	// 更新员工
	if db := l.Db.Save(&l.Staff); db.Error != nil {
		return errors.New("登录失败")
	}

	return nil
}
