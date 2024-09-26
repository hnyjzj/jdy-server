package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Api(g *gin.Engine) {
	a := g.Group("/api")
	{
		a.GET("/", func(c *gin.Context) {
			c.String(http.StatusOK, "Hello World")
		})
	}
}
