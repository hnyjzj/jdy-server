package auth_logic

import (
	"errors"
	"jdy/config"
	commonlogic "jdy/logic/common"
	authtype "jdy/types/auth"
)

type LoginLogic struct{}

// Login 登录
func (l *LoginLogic) Login(req *authtype.LoginReq) (*authtype.LoginRes, error) {
	var (
		res     = &authtype.LoginRes{}
		err     error
		captcha = &commonlogic.CaptchaLogic{}
	)

	if captcha.VerifyCaptcha(req.CaptchaId, req.Captcha) {
		return nil, errors.New("验证码错误")
	}

	res.Token = "1234567890"

	return res, err
}

// 企业微信登录
func (l *LoginLogic) Oauth(req *authtype.LoginOAuthReq) (*authtype.LoginOAuthRes, error) {
	var (
		wxwork = config.JdyAgent
		res    = &authtype.LoginOAuthRes{}
		err    error
	)

	switch req.State {
	case "wxwork_auth":

		user, err := wxwork.OAuth.Provider.GetUserInfo(req.Code)
		if err != nil {
			return nil, err
		}
		userDetail, err := wxwork.OAuth.Provider.GetUserDetail(user.UserTicket)
		if err != nil {
			return nil, err
		}

		res.Res = userDetail
	case "wxwork_qrcode":

		user, err := wxwork.OAuth.Provider.ContactFromCode(req.Code)
		if err != nil {
			return nil, err
		}
		res.Res = user
	default:
		return nil, errors.New("state错误")
	}

	res.Token = "1234567890"

	return res, err
}
