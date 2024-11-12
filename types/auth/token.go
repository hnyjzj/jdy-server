package authtype

type TokenRes struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
}

func GetTokenName(phone string) string {
	return "token_" + phone
}
