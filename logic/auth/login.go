package auth

import (
	"jdy/errors"
	"jdy/logic"
	"jdy/logic/common"
	"jdy/logic/platform/wxwork"
	"jdy/model"
	"jdy/types"

	"github.com/acmestack/gorm-plus/gplus"
	"github.com/gin-gonic/gin"
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
	case types.PlatformTypeWxWork:
		var (
			wxwork wxwork.WxWorkLogic
		)

		staff, err := wxwork.Login(ctx, req.Code)
		if err != nil {
			return nil, err
		}

		return l.token.GenerateToken(ctx, staff)
	default:
		return nil, errors.New("错误的授权方式")
	}
}
