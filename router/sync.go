package router

import (
	"jdy/controller/sync"

	"github.com/gin-gonic/gin"
)

func Sync(g *gin.Engine) {
	c := g.Group("/sync")
	{
		api := c.Group("/api")
		{
			api.GET("/list", sync.ApiController{}.List)
		}
	}
}
