package router

import (
	"jdy/controller/auth"
	"jdy/controller/common"
	"jdy/controller/user"
	"jdy/controller/workbench"
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

		// 认证
		auths := r.Group("/auth")
		{
			auths.POST("/login", auth.LoginController{}.Login) // 登录
			auths.POST("/oauth", auth.LoginController{}.OAuth) // 授权登录
		}

		users := r.Group("/user")
		{
			users.Use(middlewares.JWTMiddleware())
			{
				users.GET("/info", user.UserController{}.Info) // 获取用户信息
			}
		}

		// 工作台
		workbenchs := r.Group("/workbench")
		{
			workbenchs.Use(middlewares.JWTMiddleware())
			{
				workbenchs.POST("/add", workbench.WorkbenchController{}.Add) // 工作台添加
			}
			workbenchs.GET("/list", workbench.WorkbenchController{}.List) // 工作台列表

		}
	}
}
