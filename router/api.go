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
	"jdy/controller/statistic"
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
					uploads.POST("/product", common.UploadController{}.Product)     // 上传商品图片
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

		// 统计
		statistics := r.Group("/statistic")
		{
			statistics.Use(middlewares.JWTMiddleware())
			{
				statistics.POST("/store_sales_total", statistic.StatisticController{}.StoreSalesTotal) // 门店销售总览

				statistics.POST("/today_sales", statistic.StatisticController{}.TodaySales)     // 今日销售
				statistics.POST("/today_product", statistic.StatisticController{}.TodayProduct) // 今日货品
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
			storesm := stores.Group("/")
			{
				storesm.GET("/where", store.StoreController{}.Where) // 门店筛选
				storesm.Use(middlewares.JWTMiddleware())
				{
					storesm.POST("/create", store.StoreController{}.Create)   // 创建门店
					storesm.PUT("/update", store.StoreController{}.Update)    // 门店更新
					storesm.DELETE("/delete", store.StoreController{}.Delete) // 门店删除
					storesm.POST("/list", store.StoreController{}.List)       // 门店列表
					storesm.POST("/my", store.StoreController{}.My)           // 我的门店
					storesm.POST("/info", store.StoreController{}.Info)       // 门店详情
				}

			}

			staffs := stores.Group("/staff")
			{
				staffs.Use(middlewares.JWTMiddleware())
				{
					staffs.POST("/list", store.StoreStaffController{}.List) // 门店员工列表
					staffs.POST("/add", store.StoreStaffController{}.Add)   // 添加门店员工
					staffs.DELETE("/del", store.StoreStaffController{}.Del) // 删除门店员工
				}
			}
		}

		// 产品
		products := r.Group("/product")
		{
			// 成品
			finisheds := products.Group("/finished")
			{
				// 成品管理
				finished := finisheds.Group("/")
				{
					finished.GET("/where", product.ProductFinishedController{}.Where) // 成品筛选
					finished.Use(middlewares.JWTMiddleware())
					{
						finished.POST("/list", product.ProductFinishedController{}.List)    // 成品列表
						finished.POST("/info", product.ProductFinishedController{}.Info)    // 成品详情
						finished.PUT("/update", product.ProductFinishedController{}.Update) // 成品更新
					}
				}

				// 成品入库
				enters := finisheds.Group("/enter")
				{
					enters.GET("/where", product.ProductFinishedEnterController{}.Where) // 成品入库单筛选
					enters.Use(middlewares.JWTMiddleware())
					{
						enters.POST("/create", product.ProductFinishedEnterController{}.Create) // 创建成品入库单
						enters.POST("/list", product.ProductFinishedEnterController{}.List)     // 成品入库单列表
						enters.POST("/info", product.ProductFinishedEnterController{}.Info)     // 成品入库单详情

						enters.POST("/add_product", product.ProductFinishedEnterController{}.AddProduct)   // 添加产品
						enters.DELETE("/del_product", product.ProductFinishedEnterController{}.DelProduct) // 删除产品
						enters.PUT("/edit_product", product.ProductFinishedEnterController{}.EditProduct)  // 编辑产品

						enters.PUT("/finish", product.ProductFinishedEnterController{}.Finish) // 完成入库
						enters.PUT("/cancel", product.ProductFinishedEnterController{}.Cancel) // 取消入库
					}
				}

				// 报损单管理
				damages := finisheds.Group("/damage")
				{
					damages.GET("/where", product.ProductFinishedDamageController{}.Where) // 报损单筛选
					damages.Use(middlewares.JWTMiddleware())
					{
						damages.PUT("/create", product.ProductFinishedDamageController{}.Damage)         // 成品报损
						damages.POST("/list", product.ProductFinishedDamageController{}.List)            // 报损单列表
						damages.POST("/info", product.ProductFinishedDamageController{}.Info)            // 报损单详情
						damages.PUT("/conversion", product.ProductFinishedDamageController{}.Conversion) // 成品转换
					}
				}
			}

			// 旧料
			olds := products.Group("/old")
			{
				// 旧料管理
				old := olds.Group("/")
				{
					old.GET("/where", product.ProductOldController{}.Where)         // 旧料筛选
					old.POST("/get_class", product.ProductOldController{}.GetClass) // 获取旧料分类
					old.Use(middlewares.JWTMiddleware())
					{
						old.POST("/list", product.ProductOldController{}.List)            // 旧料列表
						old.POST("/info", product.ProductOldController{}.Info)            // 旧料详情
						old.PUT("/conversion", product.ProductOldController{}.Conversion) // 旧料转换

					}
				}
			}

			// 配件
			accessories := products.Group("/accessorie")
			{
				// 配件管理
				accessorie := accessories.Group("/")
				{
					accessorie.GET("/where", product.ProductAccessorieController{}.Where) // 配件筛选
					accessorie.Use(middlewares.JWTMiddleware())
					{
						accessorie.POST("/list", product.ProductAccessorieController{}.List) // 配件列表
						accessorie.POST("/info", product.ProductAccessorieController{}.Info) // 配件详情
					}
				}

				// 配件条目管理
				categorys := accessories.Group("/category")
				{
					categorys.GET("/where", product.ProductAccessorieCategoryController{}.Where) // 配件条目筛选
					categorys.Use(middlewares.JWTMiddleware())
					{
						categorys.POST("/create", product.ProductAccessorieCategoryController{}.Create) // 创建配件条目
						categorys.PUT("/update", product.ProductAccessorieCategoryController{}.Update)  // 更新配件条目
						categorys.DELETE("/del", product.ProductAccessorieCategoryController{}.Delete)  // 删除配件条目
						categorys.POST("/list", product.ProductAccessorieCategoryController{}.List)     // 配件条目列表
						categorys.POST("/info", product.ProductAccessorieCategoryController{}.Info)     // 配件条目详情
					}
				}

				// 配件入库
				enters := accessories.Group("/enter")
				{
					enters.GET("/where", product.ProductAccessorieEnterController{}.Where) // 配件入库单筛选
					enters.Use(middlewares.JWTMiddleware())
					{
						enters.POST("/create", product.ProductAccessorieEnterController{}.Create) // 创建配件入库单
						enters.POST("/list", product.ProductAccessorieEnterController{}.List)     // 配件入库单列表
						enters.POST("/info", product.ProductAccessorieEnterController{}.Info)     // 配件入库单详情

						enters.POST("/add_product", product.ProductAccessorieEnterController{}.AddProduct)   // 添加产品
						enters.DELETE("/del_product", product.ProductAccessorieEnterController{}.DelProduct) // 删除产品
						enters.PUT("/edit_product", product.ProductAccessorieEnterController{}.EditProduct)  // 编辑产品

						enters.PUT("/finish", product.ProductAccessorieEnterController{}.Finish) // 完成入库
						enters.PUT("/cancel", product.ProductAccessorieEnterController{}.Cancel) // 取消入库
					}
				}

				// 配件调拨
				allocate := accessories.Group("/allocate")
				{
					allocate.GET("/where", product.ProductAccessorieAllocateController{}.Where) // 调拨单筛选
					allocate.Use(middlewares.JWTMiddleware())
					{
						allocate.POST("/create", product.ProductAccessorieAllocateController{}.Create)    // 创建调拨单
						allocate.POST("/list", product.ProductAccessorieAllocateController{}.List)        // 调拨单列表
						allocate.POST("/info", product.ProductAccessorieAllocateController{}.Info)        // 调拨单详情
						allocate.PUT("/add", product.ProductAccessorieAllocateController{}.Add)           // 添加产品
						allocate.PUT("/remove", product.ProductAccessorieAllocateController{}.Remove)     // 移除产品
						allocate.PUT("/confirm", product.ProductAccessorieAllocateController{}.Confirm)   // 确认调拨
						allocate.PUT("/cancel", product.ProductAccessorieAllocateController{}.Cancel)     // 取消调拨
						allocate.PUT("/complete", product.ProductAccessorieAllocateController{}.Complete) // 完成调拨
					}
				}
			}

			// 成品调拨
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

			// 货品盘点
			inventory := products.Group("/inventory")
			{
				inventory.GET("/where", product.ProductInventoryController{}.Where) // 盘点单筛选
				inventory.Use(middlewares.JWTMiddleware())
				{
					inventory.POST("/create", product.ProductInventoryController{}.Create) // 创建盘点单
					inventory.POST("/list", product.ProductInventoryController{}.List)     // 盘点单列表
					inventory.POST("/info", product.ProductInventoryController{}.Info)     // 盘点单详情
					inventory.PUT("/change", product.ProductInventoryController{}.Change)  // 盘点单变化
				}
			}

			// 货品操作记录
			historys := products.Group("/history")
			{
				historys.GET("/where", product.ProductHistoryController{}.Where) // 货品操作记录筛选
				historys.Use(middlewares.JWTMiddleware())
				{
					historys.POST("/list", product.ProductHistoryController{}.List) // 货品操作记录列表
					historys.POST("/info", product.ProductHistoryController{}.Info) // 货品操作记录详情
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
				integral := integrals.Group("/")
				{
					integral.GET("/where", member.MemberIntegralController{}.Where) // 积分变动筛选
					integral.Use(middlewares.JWTMiddleware())
					{
						integral.POST("/list", member.MemberIntegralController{}.List)     // 积分变动记录列表
						integral.POST("/change", member.MemberIntegralController{}.Change) // 积分变动
					}
				}

				rule := integrals.Group("/rule")
				{
					rule.Use(middlewares.JWTMiddleware())
					{
						rule.POST("/finished", member.MemberIntegralRuleController{}.Finished)     // 成品积分规则
						rule.POST("/old", member.MemberIntegralRuleController{}.Old)               // 旧料积分规则
						rule.POST("/accessorie", member.MemberIntegralRuleController{}.Accessorie) // 配件积分规则
					}
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
				gold_price.Use(middlewares.JWTMiddleware())
				{
					gold_price.POST("/list", setting.GoldPriceController{}.List)     // 金价历史列表
					gold_price.POST("/create", setting.GoldPriceController{}.Create) // 创建金价
				}
			}
		}
	}
}
