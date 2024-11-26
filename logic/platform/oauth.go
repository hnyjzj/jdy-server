package platform

import (
	"errors"
	"fmt"
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

	platformType := fmt.Sprint(req.Platform)

	switch req.Platform {
	case types.PlatformTypeWxWork:
		wxwork := config.NewWechatService().JdyWork
		// 设置跳转地址
		wxwork.OAuth.Provider.WithRedirectURL(req.Uri)

		// 判断是否是微信浏览器
		if utils.IsWechat(req.Agent) {
			// 直接跳转授权页面
			wxwork.OAuth.Provider.WithState(platformType)
			res.RedirectURL, err = wxwork.OAuth.Provider.GetAuthURL()
		} else {
			// 跳转二维码页面
			wxwork.OAuth.Provider.WithState(platformType)
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
