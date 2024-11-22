package types

type PlatformOAuthReq struct {
	Agent string `json:"agent"` // 用户浏览器环境

	Uri      string `json:"uri" binding:"required"`      // 授权后重定向的回调链接地址，请使用urlencode对链接进行处理
	Platform string `json:"platform" binding:"required"` // 平台类型，可选值：wxwork
}

type PlatformOAuthRes struct {
	RedirectURL string `json:"redirect_url"` // 重定向链接
}

type PlatformJSSdkReq struct {
	Agent    string `json:"agent"`                       // 用户浏览器环境
	Platform string `json:"platform" binding:"required"` // 平台类型，可选值：wxwork
}

type PlatformJSSdkRes struct {
	Ticket string `json:"ticket"` // jsapi_ticket
}
