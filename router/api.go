package router

import (
	"jdy/config"
	"jdy/controller"
	"jdy/controller/auth"
	"jdy/controller/common"
	"jdy/controller/user"
	"jdy/controller/workbench"
	"jdy/middlewares"
	"net/http"

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
			platforms := root.Group("/")
			{
				platforms.POST("/oauth", auth.OAuthController{}.GetOauthUri) // 获取授权链接
			}

			root.GET("/jssdk/wxwork", func(c *gin.Context) {
				var (
					wxwork = config.NewWechatService().JdyWork
					client = wxwork.JSSDK.Client
				)
				wxwork.GetAccessToken()
				client.TicketEndpoint = "/cgi-bin/get_jsapi_ticket"
				jsapi, err := wxwork.JSSDK.Client.GetTicket(c, false, "jsapi")
				if err != nil {
					c.JSON(http.StatusInternalServerError, err)
				}
				client.TicketEndpoint = "/cgi-bin/ticket/get"
				agent, err := wxwork.JSSDK.Client.GetTicket(c, false, "agent_config")
				if err != nil {
					c.JSON(http.StatusInternalServerError, err)
				}

				controller.BaseController{}.Success(c, "ok", gin.H{
					"jsapi": jsapi.Get("ticket"),
					"agent": agent.Get("ticket"),
				})
			})
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
				users.POST("/create", user.UserController{}.Create) // 创建用户
				users.GET("/info", user.UserController{}.Info)      // 获取用户信息
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
