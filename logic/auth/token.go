package auth

import (
	"jdy/config"
	"jdy/errors"
	"jdy/service/redis"
	"jdy/types"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type TokenLogic struct{}

func (t *TokenLogic) GenerateToken(ctx *gin.Context, staff *types.Staff) (*types.TokenRes, error) {
	var (
		conf = config.Config.JWT
	)

	if staff.Phone == nil || *staff.Phone == "" {
		return nil, errors.ErrStaffNotFound
	}

	// 保存 ip
	staff.IP = ctx.ClientIP()

	expires := time.Now().Add(time.Second * time.Duration(conf.Expire))
	countdown_timer := time.Until(expires)

	claims := &types.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			NotBefore: jwt.NewNumericDate(time.Now().Add(time.Second * -10)),
			ExpiresAt: jwt.NewNumericDate(expires),
			Issuer:    "jdy",
		},
		Staff: staff,
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(conf.Secret))
	if err != nil {
		return nil, err
	}

	// 存入redis
	if err := redis.Client.Set(ctx, types.GetTokenName(*staff.Phone), token, countdown_timer).Err(); err != nil {
		return nil, err
	}

	res := types.TokenRes{
		Token:     token,
		ExpiresAt: expires.Unix(),
	}

	return &res, nil
}

func (t *TokenLogic) RevokeToken(ctx *gin.Context, phone string) error {
	return redis.Client.Del(ctx, types.GetTokenName(phone)).Err()
}
