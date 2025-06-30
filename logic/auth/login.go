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
	logic.BaseLogic

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
	var staff model.Staff
	if err := model.DB.
		Where(&model.Staff{Phone: req.Phone}).
		First(&staff).
		Error; err != nil {
		return nil, errors.ErrStaffUnauthorized
	}

	// 密码错误
	if staff.VerifyPassword(req.Password) != nil {
		return nil, errors.ErrPasswordIncorrect
	}

	// 生成token
	res, err := l.token.GenerateToken(ctx, &types.Staff{
		Id:         staff.Id,
		Phone:      staff.Phone,
		Nickname:   staff.Nickname,
		IsDisabled: staff.IsDisabled,
	})
	if err != nil {
		return nil, err
	}

	// 更新登录时间
	staff.UpdateLoginInfo(ctx.ClientIP())

	// 更新账号
	if err := model.DB.Save(&staff).Error; err != nil {
		return nil, errors.New("更新账号信息失败")
	}

	return res, err
}

// 授权登录
func (l *LoginLogic) Oauth(ctx *gin.Context, req *types.LoginOAuthReq) (*types.TokenRes, error) {
	var (
		staff        *model.Staff
		err          error
		wxwork_logic = &wxwork.WxWorkLogic{
			Ctx: ctx,
		}
	)
	switch req.State {
	case wxwork.WxWorkOauth:
		staff, err = wxwork_logic.OauthLogin(req.Code, false)
		if err != nil {
			return nil, err
		}

	case wxwork.WxWorkCode:
		staff, err = wxwork_logic.CodeLogin(req.Code)
		if err != nil {
			return nil, err
		}

	default:
		return nil, errors.New("错误的授权方式")
	}

	return l.token.GenerateToken(ctx, &types.Staff{
		Id:         staff.Id,
		Phone:      staff.Phone,
		Nickname:   staff.Nickname,
		IsDisabled: staff.IsDisabled,
	})
}

// 登出
func (l *LoginLogic) Logout(ctx *gin.Context, phone string) error {
	return l.token.RevokeToken(ctx, phone)
}
