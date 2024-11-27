package auth

import (
	"jdy/errors"
	"jdy/logic"
	"jdy/logic/common"
	"jdy/logic/platform/wxwork"
	"jdy/model"
	"jdy/types"

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
	var account model.Account
	if err := model.DB.
		Where(&model.Account{Phone: &req.Phone, Platform: types.PlatformTypeAccount}).
		Preload("Staff").
		First(&account).
		Error; err != nil {
		return nil, errors.ErrStaffNotFound
	}

	// 用户不存在
	if account.Staff == nil {
		return nil, errors.ErrStaffNotFound
	}

	// 密码错误
	if account.VerifyPassword(req.Password) != nil {
		return nil, errors.ErrPasswordIncorrect
	}

	// 生成token
	res, err := l.token.GenerateToken(ctx, account.Staff)
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
