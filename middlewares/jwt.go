package middlewares

import (
	"errors"
	"jdy/config"
	"jdy/service/redis"
	servertype "jdy/types/server"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// 验证 token
func verifyToken(tokenString string) (*jwt.Token, error) {
	tokenString = strings.ReplaceAll(tokenString, "Bearer ", "")

	token, err := jwt.ParseWithClaims(tokenString, &servertype.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		return config.Config.JWT.Secret, nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

// 拦截器
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求中获取 token
		authorizationHeader := c.GetHeader("Authorization")
		tokenQuery := c.Query("token")

		// 判断token
		var tokenString string
		if authorizationHeader != "" {
			tokenString = authorizationHeader
		} else if tokenQuery != "" {
			tokenString = "Bearer " + tokenQuery
		}

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "令牌不存在",
			})
			c.Abort()
			return
		}

		// 解析并验证 token
		token, err := verifyToken(tokenString)
		if err != nil {
			var e error
			switch err {
			case jwt.ErrSignatureInvalid:
				e = errors.New("令牌签名无效")
			default:
				e = err
			}
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": e.Error(),
			})
			c.Abort()
			return
		}

		// 将 token 中的数据解析出来
		claims, ok := token.Claims.(*servertype.Claims)

		// 获取缓存中的 token
		if _, err := redis.Client.Get(c, *claims.User.Phone).Result(); err != nil {
			ok = false
		}

		// 如果 token 验证通过，则将用户信息注入到 context 中，并继续处理请求
		if ok && token.Valid {
			c.Set("user", claims.User)
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "令牌无效",
			})
			c.Abort()
			return
		}
	}
}
