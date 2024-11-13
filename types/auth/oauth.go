package authtype

type OAuthWeChatWorkReq struct {
	Agent string `json:"agent"` // 用户浏览器环境

	Uri   string `json:"uri" binding:"required"`   // 授权后重定向的回调链接地址，请使用urlencode对链接进行处理
	State string `json:"state" binding:"required"` // 重定向后会带上state参数，企业可以填写a-zA-Z0-9的参数值，长度不可超过128个字节
}

type OAuthWeChatWorkRes struct {
	RedirectURL string `json:"redirect_url"` // 重定向链接
}
