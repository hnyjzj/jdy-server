package router

import (
	"jdy/middlewares"

	"github.com/gin-gonic/gin"
)

func Init(r *gin.Engine) {
	// 正则参数验证
	r.Use(middlewares.RegexValidator())

	Base(r)
	Api(r)
}
