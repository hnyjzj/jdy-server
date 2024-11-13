package auth_logic

import (
	"errors"
	"jdy/config"
	"jdy/logic"
	commonlogic "jdy/logic/common"
	usermodel "jdy/model/user"
	authtype "jdy/types/auth"

	"github.com/acmestack/gorm-plus/gplus"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LoginLogic struct {
	logic.Base

	captcha *commonlogic.CaptchaLogic
	token   TokenLogic
}

// Login 登录
func (l *LoginLogic) Login(ctx *gin.Context, req *authtype.LoginReq) (*authtype.TokenRes, error) {
	// 验证码校验
	if !l.captcha.VerifyCaptcha(req.CaptchaId, req.Captcha) {
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
	res, err := l.token.GenerateToken(ctx, user)
	if err != nil {
		return nil, err
	}

	return res, err
}

// 企业微信登录
func (l *LoginLogic) Oauth(ctx *gin.Context, req *authtype.LoginOAuthReq) (*authtype.TokenRes, error) {
	switch req.State {
	case "wxwork_auth":
		res, err := l.oauth_wxwork_auth(ctx, req.Code)
		if err != nil {
			return nil, err
		}
		return res, err
	case "wxwork_qrcode":
		res, err := l.oauth_wxwork_qrcode(ctx, req.Code)
		if err != nil {
			return nil, err
		}
		return res, err
	default:
		return nil, errors.New("错误的授权方式")
	}
}

// 企业微信授权登录
func (l *LoginLogic) oauth_wxwork_auth(ctx *gin.Context, code string) (*authtype.TokenRes, error) {
	var (
		wxwork = config.JdyAgent
	)

	// 获取用户信息
	userinfo, err := wxwork.OAuth.Provider.GetUserInfo(code)
	if err != nil {
		return nil, errors.New("获取企业微信用户信息失败")
	}
	// 获取用户详情
	detail, err := wxwork.OAuth.Provider.GetUserDetail(userinfo.UserTicket)
	if err != nil {
		return nil, errors.New("获取企业微信用户详情失败")
	}
	if detail.Mobile == "" {
		return nil, errors.New("非企业微信授权用户禁止登录")
	}

	// 查询用户
	query, u := gplus.NewQuery[usermodel.User]()
	query.Eq(&u.Phone, detail.Mobile)
	user, db := gplus.SelectOne(query)
	// 用户不存在
	if db.Error != nil || errors.Is(db.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("用户不存在")
	}

	// 生成token
	res, err := l.token.GenerateToken(ctx, user)
	if err != nil {
		return nil, err
	}

	return res, err
}

// 企业微信扫码登录
func (l *LoginLogic) oauth_wxwork_qrcode(ctx *gin.Context, code string) (*authtype.TokenRes, error) {
	var (
		wxwork = config.JdyAgent
	)

	// 获取用户信息
	userinfo, err := wxwork.OAuth.Provider.GetUserInfo(code)
	if err != nil {
		return nil, errors.New("获取企业微信用户信息失败")
	}
	if userinfo.UserID == "" {
		return nil, errors.New("非企业微信授权用户禁止登录")
	}

	// 查询用户
	query, u := gplus.NewQuery[usermodel.User]()
	query.Eq(&u.Username, userinfo.UserID)
	user, db := gplus.SelectOne(query)
	// 用户不存在
	if db.Error != nil || errors.Is(db.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("用户不存在")
	}

	// 生成token
	res, err := l.token.GenerateToken(ctx, user)
	if err != nil {
		return nil, err
	}

	return res, err
}
