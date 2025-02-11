package router

import (
	"jdy/controller/auth"
	"jdy/controller/common"
	"jdy/controller/member"
	"jdy/controller/order"
	"jdy/controller/platform"
	"jdy/controller/product"
	"jdy/controller/setting"
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

			// 上传
			uploads := root.Group("/upload")
			{
				uploads.Use(middlewares.JWTMiddleware())
				{
					uploads.POST("/avatar", common.UploadController{}.Avatar)       // 上传头像
					uploads.POST("/workbench", common.UploadController{}.Workbench) // 上传工作台图片
					uploads.POST("/store", common.UploadController{}.Store)         // 上传门店图片
				}
			}
		}

		// 认证
		auths := r.Group("/auth")
		{
			auths.POST("/login", auth.LoginController{}.Login)                                // 登录
			auths.POST("/oauth", auth.LoginController{}.OAuth)                                // 授权登录
			auths.POST("/logout", middlewares.JWTMiddleware(), auth.LoginController{}.Logout) // 登出
		}

		// 员工
		staffs := r.Group("/staff")
		{
			staffs.GET("/where", staff.StaffController{}.Where) // 员工筛选
			staffs.Use(middlewares.JWTMiddleware())
			{
				staffs.POST("/list", staff.StaffController{}.List)     // 员工列表
				staffs.POST("/create", staff.StaffController{}.Create) // 创建账号
				staffs.POST("/info", staff.StaffController{}.Info)     // 员工详情
				staffs.GET("/my", staff.StaffController{}.My)          // 获取我的信息
				staffs.PUT("/update", staff.StaffController{}.Update)  // 更新员工信息
			}
		}

		// 工作台
		workbenchs := r.Group("/workbench")
		{
			workbenchs.GET("/list", workbench.WorkbenchController{}.List)      // 工作台列表
			workbenchs.POST("/search", workbench.WorkbenchController{}.Search) // 工作台搜索
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
			stores.GET("/where", store.StoreController{}.Where) // 门店筛选
			stores.Use(middlewares.JWTMiddleware())
			{
				stores.POST("/create", store.StoreController{}.Create)   // 创建门店
				stores.PUT("/update", store.StoreController{}.Update)    // 门店更新
				stores.DELETE("/delete", store.StoreController{}.Delete) // 门店删除
				stores.POST("/list", store.StoreController{}.List)       // 门店列表
				stores.POST("/my", store.StoreController{}.My)           // 我的门店
				stores.POST("/info", store.StoreController{}.Info)       // 门店详情

				staffs := stores.Group("/staff")
				{
					staffs.POST("/list", store.StoreStaffController{}.List) // 门店员工列表
					staffs.POST("/add", store.StoreStaffController{}.Add)   // 添加门店员工
				}
			}
		}

		// 产品
		products := r.Group("/product")
		{
			// 产品管理
			products = products.Group("/")
			{
				products.GET("/where", product.ProductController{}.Where) // 产品筛选
				products.Use(middlewares.JWTMiddleware())
				{
					products.POST("/list", product.ProductController{}.List)            // 产品列表
					products.POST("/info", product.ProductController{}.Info)            // 产品详情
					products.PUT("/update", product.ProductController{}.Update)         // 产品更新
					products.PUT("/damage", product.ProductController{}.Damage)         // 产品报损
					products.PUT("/conversion", product.ProductController{}.Conversion) // 产品转换
				}
			}

			// 产品入库
			enters := products.Group("/enter")
			{
				enters.GET("/where", product.ProductEnterController{}.Where) // 入库单筛选
				enters.Use(middlewares.JWTMiddleware())
				{
					enters.POST("/create", product.ProductEnterController{}.Create) // 创建入库单
					enters.POST("/list", product.ProductEnterController{}.List)     // 入库单列表
					enters.POST("/info", product.ProductEnterController{}.Info)     // 入库单详情
				}
			}

			// 产品调拨
			allocate := products.Group("/allocate")
			{
				allocate.GET("/where", product.ProductAllocateController{}.Where) // 调拨单筛选
				allocate.Use(middlewares.JWTMiddleware())
				{
					allocate.POST("/create", product.ProductAllocateController{}.Create)    // 创建调拨单
					allocate.POST("/list", product.ProductAllocateController{}.List)        // 调拨单列表
					allocate.POST("/info", product.ProductAllocateController{}.Info)        // 调拨单详情
					allocate.PUT("/add", product.ProductAllocateController{}.Add)           // 添加产品
					allocate.PUT("/remove", product.ProductAllocateController{}.Remove)     // 移除产品
					allocate.PUT("/confirm", product.ProductAllocateController{}.Confirm)   // 确认调拨
					allocate.PUT("/cancel", product.ProductAllocateController{}.Cancel)     // 取消调拨
					allocate.PUT("/complete", product.ProductAllocateController{}.Complete) // 完成调拨
				}
			}
		}

		// 会员
		members := r.Group("/member")
		{
			members.GET("/where", member.MemberController{}.Where) // 会员筛选
			members.Use(middlewares.JWTMiddleware())
			{
				members.POST("/create", member.MemberController{}.Create) // 创建会员
				members.POST("/list", member.MemberController{}.List)     // 会员列表
				members.POST("/info", member.MemberController{}.Info)     // 会员详情
				members.PUT("/update", member.MemberController{}.Update)  // 会员更新

			}

			integrals := members.Group("/integral")
			{
				integrals.GET("/where", member.MemberIntegralController{}.Where) // 积分变动筛选
				integrals.Use(middlewares.JWTMiddleware())
				{
					integrals.POST("/list", member.MemberIntegralController{}.List)     // 积分变动记录列表
					integrals.POST("/change", member.MemberIntegralController{}.Change) // 积分变动
				}
			}
		}

		// 订单
		orders := r.Group("/order")
		{
			orders.GET("/where", order.OrderController{}.Where) // 订单筛选
			orders.Use(middlewares.JWTMiddleware())
			{
				orders.POST("/create", order.OrderController{}.Create) // 创建订单
				orders.POST("/list", order.OrderController{}.List)     // 订单列表
				orders.POST("/info", order.OrderController{}.Info)     // 订单详情
			}
		}

		// 设置
		settings := r.Group("/setting")
		{
			gold_price := settings.Group("/gold_price")
			{
				gold_price.GET("/get", setting.GoldPriceController{}.Get) // 获取金价
				gold_price.Use(middlewares.JWTMiddleware())
				{
					gold_price.POST("/list", setting.GoldPriceController{}.List)     // 金价历史列表
					gold_price.POST("/create", setting.GoldPriceController{}.Create) // 创建金价
				}
			}
		}
	}
}
