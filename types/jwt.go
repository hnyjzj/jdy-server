package types

import (
	"jdy/model"

	"github.com/golang-jwt/jwt/v5"
)

// 定义 token 中的数据结构
type Claims struct {
	jwt.RegisteredClaims
	User model.User `json:"user"`
}
