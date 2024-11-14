package router

import (
	"jdy/controller/auth"
	"jdy/controller/common"
	"jdy/controller/user"
	"jdy/middlewares"

	"github.com/gin-gonic/gin"
)

func Api(g *gin.Engine) {
	// 跨域
	g.Use(middlewares.Cors())

	r := g.Group("/")
	{
		root := r.Group("/")
		{
			root.GET("/get_captcha_image", common.CaptchaController{}.GetImage) // 获取验证码图片
		}

		oauth := r.Group("/")
		{
			oauth.POST("/oauth", auth.OAuthController{}.GetUri) // 获取授权链接
		}

		login := r.Group("/login")
		{
			login.POST("/", auth.LoginController{}.Login)      // 登录
			login.POST("/oauth", auth.LoginController{}.OAuth) // 授权登录
		}

		r.Use(middlewares.JWTMiddleware())
		{
			users := r.Group("/user")
			{
				users.GET("/info", user.UserController{}.Info) // 获取用户信息
			}
		}
	}
}
