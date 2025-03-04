package types

import "jdy/enums"

type PlatformReq struct {
	Platform enums.PlatformType `json:"platform" binding:"required"` // 平台类型，可选值：wxwork
}

type PlatformOAuthReq struct {
	Agent string `json:"agent"` // 用户浏览器环境
	PlatformReq

	Uri string `json:"uri" binding:"required"` // 授权后重定向的回调链接地址，请使用urlencode对链接进行处理
}

type PlatformOAuthRes struct {
	RedirectURL string `json:"redirect_url"` // 重定向链接
}

type PlatformJSSdkReq struct {
	Agent string `json:"agent"` // 用户浏览器环境
	PlatformReq

	Type string `json:"type" binding:"required"` // jsapi_ticket类型，可选值：jsapi, agent
}

type PlatformJSSdkRes struct {
	Ticket string `json:"ticket"` // jsapi_ticket
}
