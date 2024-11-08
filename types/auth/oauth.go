package authtype

type OAuthWeChatWorkReq struct {
	Agent string `json:"agent"`

	Uri   string `json:"uri" binding:"required"`
	State string `json:"state" binding:"required"`
}

type OAuthWeChatWorkRes struct {
	RedirectURL string `json:"redirect_url"`
}
