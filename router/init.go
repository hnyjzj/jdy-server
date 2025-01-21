package router

import (
	"jdy/middlewares"

	"github.com/gin-gonic/gin"
)

func Init(r *gin.Engine) {
	// 正则参数验证
	r.Use(middlewares.CustomValidator())

	Base(r)
	Api(r)
	CallBack(r)
}
