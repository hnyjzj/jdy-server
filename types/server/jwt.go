package servertype

import (
	usermodel "jdy/model/user"

	"github.com/golang-jwt/jwt/v5"
)

// 定义 token 中的数据结构
type Claims struct {
	jwt.RegisteredClaims
	User usermodel.User `json:"user"`
}
