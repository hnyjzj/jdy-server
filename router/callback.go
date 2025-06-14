package router

import (
	"jdy/controller/callback"

	"github.com/gin-gonic/gin"
)

func CallBack(g *gin.Engine) {
	c := g.Group("/callback")
	{
		apis := c.Group("/api")
		{
			apis.GET("/sync_api_list", callback.ApiController{}.SyncApiList)
		}

		ww := c.Group("/wxwork")
		{
			ww.GET("/jdy", callback.WxWorkCongtroller{}.JdyVerify)
			ww.POST("/jdy", callback.WxWorkCongtroller{}.JdyNotify)

			ww.GET("/contact", callback.WxWorkCongtroller{}.ContactsVerify)
			ww.POST("/contact", callback.WxWorkCongtroller{}.ContactsNotify)
		}
	}
}
