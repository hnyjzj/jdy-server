package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CallBack(g *gin.Engine) {
	c := g.Group("/callback")
	{
		ww := c.Group("workwechat")
		{
			ww.GET("/jdy", func(ctx *gin.Context) {
				ctx.JSON(http.StatusOK, gin.H{
					"code": 200,
					"msg":  "success",
				})
			})
		}
	}
}
