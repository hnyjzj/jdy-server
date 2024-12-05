package router

import (
	"jdy/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Base(g *gin.Engine) {
	r := g.Group("/")
	{
		r.GET("/", func(c *gin.Context) {
			c.String(http.StatusOK, "Hello World")
		})
	}

	// 访问上传文件
	dir := config.Config.Storage.Local.Root
	if dir != "" {
		r.Static("/uploads", dir)
	}
}
