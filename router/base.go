package router

import (
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
}
