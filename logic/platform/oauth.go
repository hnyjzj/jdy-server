package platform

import (
	"errors"
	"jdy/config"
	"jdy/types"
	"jdy/utils"
)

// 获取授权链接
func (l *PlatformLogic) OauthUri(req *types.PlatformOAuthReq) (*types.PlatformOAuthRes, error) {
	var (
		res = types.PlatformOAuthRes{}
		err error
	)

	switch req.Platform {
	case "wxwork":
		wxwork := config.NewWechatService().JdyWork
		// 设置跳转地址
		wxwork.OAuth.Provider.WithRedirectURL(req.Uri)

		// 判断是否是微信浏览器
		if utils.IsWechat(req.Agent) {
			// 直接跳转授权页面
			wxwork.OAuth.Provider.WithState(req.Platform + "_auth")
			res.RedirectURL, err = wxwork.OAuth.Provider.GetAuthURL()
		} else {
			// 跳转二维码页面
			wxwork.OAuth.Provider.WithState(req.Platform + "_qrcode")
			res.RedirectURL, err = wxwork.OAuth.Provider.GetQrConnectURL()
		}

		if err != nil {
			return nil, err
		}

	default:
		return nil, errors.New("不支持的授权类型")
	}

	return &res, nil
}
