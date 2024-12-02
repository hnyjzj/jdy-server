package router

import (
	"jdy/controller/auth"
	"jdy/controller/common"
	"jdy/controller/platform"
	"jdy/controller/staff"
	"jdy/controller/store"
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
			platforms := root.Group("/platform")
			{
				platforms.POST("/oauth", platform.PlatformController{}.OauthUri) // 获取授权链接
				platforms.POST("/jssdk", platform.PlatformController{}.JSSDK)    // 获取JSSDK
			}
		}

		// 认证
		auths := r.Group("/auth")
		{
			auths.POST("/login", auth.LoginController{}.Login) // 登录
			auths.POST("/oauth", auth.LoginController{}.OAuth) // 授权登录
		}

		// 员工
		staffs := r.Group("/staff")
		{
			staffs.Use(middlewares.JWTMiddleware())
			{
				staffs.POST("/create", staff.StaffController{}.Create) // 创建账号
				staffs.GET("/info", staff.StaffController{}.Info)      // 获取员工信息
			}
		}

		// 工作台
		workbenchs := r.Group("/workbench")
		{
			workbenchs.GET("/list", workbench.WorkbenchController{}.List) // 工作台列表
			workbenchs.Use(middlewares.JWTMiddleware())
			{
				workbenchs.POST("/add", workbench.WorkbenchController{}.Add)      // 工作台添加
				workbenchs.PUT("/update", workbench.WorkbenchController{}.Update) // 工作台更新
				workbenchs.DELETE("/del", workbench.WorkbenchController{}.Del)    // 工作台删除
			}
		}

		// 门店
		stores := r.Group("/store")
		{
			stores.Use(middlewares.JWTMiddleware())
			{
				stores.POST("/create", store.StoreController{}.Create)   // 创建门店
				stores.PUT("/update", store.StoreController{}.Update)    // 门店更新
				stores.DELETE("/delete", store.StoreController{}.Delete) // 门店删除
				stores.POST("/list", store.StoreController{}.List)       // 门店列表
				stores.GET("/info", store.StoreController{}.Info)        // 门店详情
			}
		}
	}
}
