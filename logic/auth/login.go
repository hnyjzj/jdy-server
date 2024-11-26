package auth

import (
	"fmt"
	"jdy/config"
	"jdy/errors"
	"jdy/logic"
	"jdy/logic/common"
	"jdy/model"
	"jdy/types"
	"strconv"
	"time"

	"github.com/acmestack/gorm-plus/gplus"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LoginLogic struct {
	logic.Base

	captcha common.CaptchaLogic
	token   TokenLogic
}

// Login 登录
func (l *LoginLogic) Login(ctx *gin.Context, req *types.LoginReq) (*types.TokenRes, error) {
	// 验证码校验
	if !l.captcha.VerifyCaptcha(req.CaptchaId, req.Captcha) {
		return nil, errors.ErrInvalidCaptcha
	}

	// 查询用户
	query, u := gplus.NewQuery[model.Account]()
	query.Eq(&u.Phone, req.Phone).And().Eq(&u.Platform, types.PlatformTypeAccount)

	account, db := gplus.SelectOne(query)
	// 用户不存在
	if db.Error != nil {
		return nil, errors.ErrStaffNotFound
	}

	// 密码错误
	if account.VerifyPassword(req.Password) != nil {
		return nil, errors.ErrPasswordIncorrect
	}

	// 查询用户
	staff, db := gplus.SelectById[model.Staff](account.StaffId)
	if db.Error != nil {
		return nil, errors.ErrStaffNotFound
	}

	// 生成token
	res, err := l.token.GenerateToken(ctx, staff)
	if err != nil {
		return nil, err
	}

	return res, err
}

// 授权登录
func (l *LoginLogic) Oauth(ctx *gin.Context, req *types.LoginOAuthReq) (*types.TokenRes, error) {
	switch req.State {
	case string(types.PlatformTypeWxWork):
		res, err := l.oauth_wxwork_auth(ctx, req.Code)
		if err != nil {
			return nil, err
		}
		return res, err
	default:
		return nil, errors.New("错误的授权方式")
	}
}

// 企业微信授权登录
func (l *LoginLogic) oauth_wxwork_auth(ctx *gin.Context, code string) (*types.TokenRes, error) {
	var (
		wxwork = config.NewWechatService().JdyWork
	)

	// 获取用户信息
	user, err := wxwork.OAuth.Provider.GetUserInfo(code)
	if err != nil || user.UserID == "" {
		return nil, errors.New("获取企业微信用户信息失败")
	}

	// 查询账号
	aq, a := gplus.NewQuery[model.Account]()
	aq.Eq(&a.Username, user.UserID).And().Eq(&a.Platform, types.PlatformTypeWxWork)
	account, adb := gplus.SelectOne(aq)
	// 账号不存在
	if adb.Error != nil || errors.Is(adb.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("非授权用户禁止登录")
	}
	// 判断手机号
	if account.Phone == nil && user.UserTicket == "" {
		return nil, errors.New("首次登录需通过企业微信工作台打开并授权手机号")
	}

	// 获取用户详情
	detail, err := wxwork.OAuth.Provider.GetUserDetail(user.UserTicket)
	if err != nil {
		return nil, errors.New("获取企业微信用户详情失败")
	}
	// 获取不到手机号
	if detail.Mobile == "" {
		return nil, errors.New("非企业微信授权用户禁止登录")
	}

	// 账号没有手机号
	if account.Phone == nil {
		account.Phone = &detail.Mobile
	}
	// 手机号不一致
	if *account.Phone != detail.Mobile {
		return nil, errors.New("手机号不一致")
	}

	sq, s := gplus.NewQuery[model.Staff]()
	sq.Eq(&s.Phone, detail.Mobile)
	staff, sdb := gplus.SelectOne(sq)
	if sdb.Error != nil && !errors.Is(sdb.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("获取手机号授权失败，请重新授权登录")
	}
	// 用户不存在
	if sdb.Error != nil && errors.Is(sdb.Error, gorm.ErrRecordNotFound) {
		info, err := wxwork.User.Get(ctx, detail.UserID)
		if err != nil {
			return nil, errors.New("获取企业微信用户信息失败")
		}

		// 更新员工
		staff.Phone = &detail.Mobile
		staff.NickName = info.Name
		staff.Avatar = info.Avatar
		staff.Email = info.Email

		if err := gplus.Insert(staff); err.Error != nil {
			fmt.Printf("err.Error.Error(): %v\n", err.Error.Error())
			return nil, errors.New("更新信息失败")
		}

		// 更新账号
		account.StaffId = &staff.Id
	} else {
		if staff.NickName == "" {
			staff.NickName = *account.Nickname
		}
		if staff.Avatar == "" {
			staff.Avatar = *account.Avatar
		}
		if staff.Gender == 0 {
			// 0表示未定义，1表示男性，2表示女性。
			gender, err := strconv.Atoi(detail.Gender)
			if err != nil {
				gender = 0
			}
			staff.Gender = uint(gender)
		}
		if staff.Email == "" {
			staff.Email = detail.Email
		}
		if db := gplus.UpdateById(&staff); db.Error != nil {
			return nil, errors.New("更新账号信息失败")
		}
	}

	// 更新账号
	account.LastLoginIp = ctx.ClientIP()
	now := time.Now()
	account.LastLoginAt = &now
	if db := gplus.UpdateById(&account); db.Error != nil {
		return nil, errors.New("更新账号信息失败")
	}

	// 生成token
	res, err := l.token.GenerateToken(ctx, staff)
	if err != nil {
		return nil, err
	}

	return res, err
}
