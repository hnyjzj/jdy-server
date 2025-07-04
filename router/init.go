package router

import (
	"jdy/middlewares"

	"github.com/gin-gonic/gin"
)

var Router *gin.Engine

func Init(r *gin.Engine) {
	Router = r
	// 正则参数验证
	Router.Use(middlewares.CustomValidator())

	Base(Router)
	Api(Router)
	CallBack(Router)
	Sync(Router)
}
