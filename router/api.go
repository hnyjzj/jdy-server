package router

import (
	"jdy/controller/auth"
	"jdy/controller/common"

	"github.com/gin-gonic/gin"
)

func Api(g *gin.Engine) {
	r := g.Group("/")
	{
		root := r.Group("/")
		{
			root.GET("/get_captcha_image", common.CaptchaController{}.GetImage)
		}

		oauth := r.Group("/")
		{
			oauth.POST("/oauth", auth.OAuthController{}.GetUri)
		}

		login := r.Group("/login")
		{
			login.POST("/", auth.LoginController{}.Login)
			login.POST("/oauth", auth.LoginController{}.OAuth)
		}
	}
}
