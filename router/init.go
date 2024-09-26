package router

import "github.com/gin-gonic/gin"

func Init(r *gin.Engine) {
	Base(r)
	Api(r)
}
