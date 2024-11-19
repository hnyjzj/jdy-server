package types

type TokenRes struct {
	Token     string `json:"token"`      // token
	ExpiresAt int64  `json:"expires_at"` // 过期时间
}

// 获取 token 名字
func GetTokenName(phone string) string {
	return "token_" + phone
}
