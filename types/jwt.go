package types

import (
	"github.com/golang-jwt/jwt/v5"
)

// 定义 token 中的数据结构
type Claims struct {
	jwt.RegisteredClaims

	Staff *Staff `json:"staff"`
}

type Staff struct {
	Id    string `json:"id"`
	Phone string `json:"phone"`

	IsDisabled bool `json:"is_disabled"`

	Platform PlatformType `json:"platform"`
}
