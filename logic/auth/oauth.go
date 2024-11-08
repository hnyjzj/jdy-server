package auth_logic

import (
	"jdy/config"
	authtype "jdy/types/auth"
	"jdy/utils"
)

type OAuthLogic struct{}

// oauth 登录
func (l *OAuthLogic) GetUri(req *authtype.OAuthWeChatWorkReq) (*authtype.OAuthWeChatWorkRes, error) {

	var (
		wxwork = config.JdyAgent
		res    = authtype.OAuthWeChatWorkRes{}
		err    error
	)

	switch req.State {
	case "wxwork":
		// 设置跳转地址
		wxwork.OAuth.Provider.WithRedirectURL(req.Uri)

		// 判断是否是微信浏览器
		if utils.IsWechat(req.Agent) {
			// 直接跳转授权页面
			wxwork.OAuth.Provider.WithState(req.State + "_auth")
			res.RedirectURL, err = wxwork.OAuth.Provider.GetAuthURL()
		} else {
			// 跳转二维码页面
			wxwork.OAuth.Provider.WithState(req.State + "_qrcode")
			res.RedirectURL, err = wxwork.OAuth.Provider.GetQrConnectURL()
		}

		if err != nil {
			return nil, err
		}

	default:
		return nil, err
	}

	return &res, err
}
