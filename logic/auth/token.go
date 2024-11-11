package auth_logic

import (
	"jdy/config"
	usermodel "jdy/model/user"
	"jdy/service/redis"
	authtype "jdy/types/auth"
	servertype "jdy/types/server"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type TokenLogic struct{}

func (t *TokenLogic) GenerateToken(ctx *gin.Context, user *usermodel.User) (*authtype.TokenRes, error) {
	var (
		conf = config.Config.JWT
	)

	expires := time.Now().Add(time.Second * time.Duration(conf.Expire))

	claims := &servertype.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			NotBefore: jwt.NewNumericDate(time.Now().Add(time.Second * -10)),
			ExpiresAt: jwt.NewNumericDate(expires),
			Issuer:    "jdy",
		},
		User: *user,
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(conf.Secret))
	if err != nil {
		return nil, err
	}

	// 存入redis
	if err := redis.Client.Set(ctx, "token:"+*user.Phone, token, time.Duration(conf.Expire)*time.Second).Err(); err != nil {
		return nil, err
	}

	res := authtype.TokenRes{
		Token:     token,
		ExpiresAt: expires.Unix(),
	}

	return &res, nil
}
