package auth

import (
	"jdy/config"
	"jdy/errors"
	"jdy/model"
	"jdy/service/redis"
	"jdy/types"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type TokenLogic struct{}

func (t *TokenLogic) GenerateToken(ctx *gin.Context, user *model.User) (*types.TokenRes, error) {
	var (
		conf = config.Config.JWT
	)

	if user.Phone == nil || *user.Phone == "" {
		return nil, errors.ErrUserNotFound
	}

	expires := time.Now().Add(time.Second * time.Duration(conf.Expire))

	claims := &types.Claims{
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
	if err := redis.Client.Set(ctx, types.GetTokenName(*user.Phone), token, time.Duration(conf.Expire)*time.Second).Err(); err != nil {
		return nil, err
	}

	res := types.TokenRes{
		Token:     token,
		ExpiresAt: expires.Unix(),
	}

	return &res, nil
}
