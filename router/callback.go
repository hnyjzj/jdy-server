package router

import (
	"jdy/controller/callback"

	"github.com/gin-gonic/gin"
)

func CallBack(g *gin.Engine) {
	c := g.Group("/callback")
	{
		ww := c.Group("/wxwork")
		{
			ww.GET("/jdy", callback.WxWorkCongtroller{}.JdyVerify)
			ww.POST("/jdy", callback.WxWorkCongtroller{}.JdyNotify)
		}
	}
}
