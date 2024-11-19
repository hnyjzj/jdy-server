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
			// 验证码
			captchas := root.Group("/captcha")
			{
				captchas.GET("/image", common.CaptchaController{}.Image) // 获取验证码图片
			}

			// 平台
			platforms := r.Group("/")
			{
				platforms.POST("/oauth", auth.OAuthController{}.GetOauthUri) // 获取授权链接
			}
		}

		logins := r.Group("/login")
		{
			logins.POST("/", auth.LoginController{}.Login)      // 登录
			logins.POST("/oauth", auth.LoginController{}.OAuth) // 授权登录
		}

		users := r.Group("/user")
		{
			users.Use(middlewares.JWTMiddleware())
			{
				users.GET("/info", user.UserController{}.Info) // 获取用户信息
			}
		}
	}
}
