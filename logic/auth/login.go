package auth_logic

import (
	"errors"
	"jdy/config"
	commonlogic "jdy/logic/common"
	usermodel "jdy/model/user"
	authtype "jdy/types/auth"

	"github.com/acmestack/gorm-plus/gplus"
	"github.com/gin-gonic/gin"
)

type LoginLogic struct{}

// Login 登录
func (l *LoginLogic) Login(ctx *gin.Context, req *authtype.LoginReq) (*authtype.TokenRes, error) {
	var (
		res        = &authtype.TokenRes{}
		err        error
		captcha    = &commonlogic.CaptchaLogic{}
		tokenlogic = &TokenLogic{}
	)

	// 验证码校验
	if !captcha.VerifyCaptcha(req.CaptchaId, req.Captcha) {
		return nil, errors.New("验证码错误")
	}

	// 查询用户
	query, u := gplus.NewQuery[usermodel.User]()
	query.Eq(&u.Phone, req.Phone)
	user, db := gplus.SelectOne(query)
	// 用户不存在
	if db.Error != nil {
		return nil, errors.New("用户不存在")
	}

	// 密码错误
	if user.VerifyPassword(req.Password) != nil {
		return nil, errors.New("密码不正确")
	}

	// 生成token
	res, err = tokenlogic.GenerateToken(ctx, user)
	if err != nil {
		return nil, err
	}

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
