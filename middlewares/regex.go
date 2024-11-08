package middlewares

import (
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// 正则表达式验证器
func RegexValidator() gin.HandlerFunc {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("regex", func(fl validator.FieldLevel) bool {
			// 从传递给验证器的 tag 获取正则表达式。
			regex := fl.Param()

			// 将字段值转换为字符串并与正则表达式进行匹配。
			field := fl.Field().String()
			matched, _ := regexp.MatchString(regex, field)

			return matched
		})
	}
	// 返回 Gin 中间件。
	return func(c *gin.Context) {
		c.Next()
	}
}
