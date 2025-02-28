package middlewares

import (
	"jdy/config"
	"jdy/errors"
	"jdy/service/redis"
	"jdy/types"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	redisv9 "github.com/redis/go-redis/v9"
)

// 验证 token
func verifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &types.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		return []byte(config.Config.JWT.Secret), nil
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
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

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
		claims, ok := token.Claims.(*types.Claims)

		// 如果 token 验证通过，则将用户信息注入到 context 中，并继续处理请求
		if ok && token.Valid {
			// 获取缓存中的 token
			res, err := redis.Client.Get(c, types.GetTokenName(*claims.Staff.Phone)).Result()
			// 如果缓存中没有 token 或者 token 不一致，则返回 401
			if err == redisv9.Nil || res != tokenString {
				c.JSON(http.StatusUnauthorized, gin.H{
					"code":    http.StatusUnauthorized,
					"message": "令牌不符",
				})

				c.Abort()
				return
			} else if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    http.StatusInternalServerError,
					"message": "服务器错误",
				})
				c.Abort()
				return
			}

			// 将用户信息注入到 context 中
			c.Set("staff", claims.Staff)
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
