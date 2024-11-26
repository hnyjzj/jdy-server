package wxwork

import (
	"jdy/config"
	"jdy/utils"
)

// 获取授权链接
func (w *WxWorkLogic) OauthUri(agent string, state string, uri string) (string, error) {
	wxwork := config.NewWechatService().JdyWork
	// 设置跳转地址
	wxwork.OAuth.Provider.WithRedirectURL(uri)

	// 判断是否是微信浏览器
	if utils.IsWechat(agent) {
		// 直接跳转授权页面
		wxwork.OAuth.Provider.WithState(state)
		redirect_url, err := wxwork.OAuth.Provider.GetAuthURL()
		if err != nil {
			return "", err
		}
		return redirect_url, nil

	} else {
		// 跳转二维码页面
		wxwork.OAuth.Provider.WithState(state)
		redirect_url, err := wxwork.OAuth.Provider.GetQrConnectURL()
		if err != nil {
			return "", err
		}
		return redirect_url, nil
	}
}
