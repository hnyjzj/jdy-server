package router

import (
	"jdy/controller/auth"
	"jdy/controller/common"
	"jdy/controller/member"
	"jdy/controller/order"
	"jdy/controller/platform"
	"jdy/controller/product"
	"jdy/controller/region"
	"jdy/controller/setting"
	"jdy/controller/staff"
	"jdy/controller/statistic"
	"jdy/controller/statistic/boos"
	"jdy/controller/statistic/payment"
	"jdy/controller/statistic/sale"
	"jdy/controller/statistic/stock"
	"jdy/controller/statistic/today"
	"jdy/controller/store"
	"jdy/controller/target"
	"jdy/controller/workbench"
	"jdy/middlewares"

	"github.com/gin-gonic/gin"
)

func Api(g *gin.Engine) {
	// 跨域
	g.Use(middlewares.Cors())

	r := g.Group("/api")
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
				platforms.POST("/oauth", platform.PlatformController{}.OauthUri)   // 获取授权链接
				platforms.POST("/jssdk", platform.PlatformController{}.JSSDK)      // 获取JSSDK
				platforms.POST("/get_user", platform.PlatformController{}.GetUser) // 获取用户信息
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
					uploads.POST("/order", common.UploadController{}.Order)         // 上传订单图片
				}
			}

			// 记录
			logs := root.Group("/log")
			{
				logs.POST("/on_capture_screen", common.LogController{}.OnCaptureScreen) // 截屏
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
				staffs.POST("/list", staff.StaffController{}.List)       // 员工列表
				staffs.POST("/create", staff.StaffController{}.Create)   // 创建账号
				staffs.POST("/info", staff.StaffController{}.Info)       // 员工详情
				staffs.GET("/my", staff.StaffController{}.My)            // 获取我的信息
				staffs.PUT("/edit", staff.StaffController{}.Edit)        // 编辑员工信息
				staffs.PUT("/update", staff.StaffController{}.Update)    // 更新员工信息
				staffs.DELETE("/delete", staff.StaffController{}.Delete) // 删除员工
			}
		}

		// 统计
		statistics := r.Group("/statistic")
		{
			bosses := statistics.Group("/boos") // Boos看板
			{
				bosses.GET("/where", boos.BoosController{}.BoosWhere) // Boos看板筛选
				bosses.Use(middlewares.JWTMiddleware())
				{
					performance := bosses.Group("/performance") // 业绩统计
					{
						performance.GET("/where", boos.BoosController{}.PerformanceWhere) // 业绩统计筛选
						performance.POST("/data", boos.BoosController{}.PerformanceData)  // 业绩统计数据
					}

					finished_sales := bosses.Group("/finished_sales") // 成品销售
					{
						finished_sales.GET("/where", boos.BoosController{}.FinishedSalesWhere) // 成品销售筛选
						finished_sales.POST("/data", boos.BoosController{}.FinishedSalesData)  // 成品销售数据
					}

					finished_stock := bosses.Group("/finished_stock") // 成品库存
					{
						finished_stock.GET("/where", boos.BoosController{}.FinishedStockWhere) // 成品库存筛选
						finished_stock.POST("/data", boos.BoosController{}.FinishedStockData)  // 成品库存数据
					}

					old_stock := bosses.Group("/old_stock") // 旧料库存
					{
						old_stock.GET("/where", boos.BoosController{}.OldStockWhere) // 旧料库存筛选
						old_stock.POST("/data", boos.BoosController{}.OldStockData)  // 旧料库存数据
					}

					old_exchanges := bosses.Group("/old_exchange") // 旧料兑换
					{
						old_exchanges.GET("/where", boos.BoosController{}.OldExchangeWhere) // 旧料兑换筛选
						old_exchanges.POST("/data", boos.BoosController{}.OldExchangeData)  // 旧料兑换数据
					}

					old_recycles := bosses.Group("/old_recycle") // 旧料回收
					{
						old_recycles.GET("/where", boos.BoosController{}.OldRecycleWhere) // 旧料回收筛选
						old_recycles.POST("/data", boos.BoosController{}.OldRecycleData)  // 旧料回收数据
					}

					payments := bosses.Group("/payments") // 收支统计
					{
						payments.GET("/where", boos.BoosController{}.PaymentsWhere) // 收支统计筛选
						payments.POST("/data", boos.BoosController{}.PaymentsData)  // 收支统计数据
					}

				}
			}

			todays := statistics.Group("/today") // 今日统计
			{
				todays.Use(middlewares.JWTMiddleware())
				{
					todays.POST("/sales", today.ToDayController{}.Sales)     // 今日销售
					todays.POST("/product", today.ToDayController{}.Product) // 今日货品
					todays.POST("/payment", today.ToDayController{}.Payment) // 今日收支
				}
			}

			payments := statistics.Group("/payment") // 收支统计
			{
				payments.GET("/where", payment.PaymentController{}.Where) // 收支统计筛选
				payments.Use(middlewares.JWTMiddleware())
				{
					payments.POST("/data", payment.PaymentController{}.Data) // 收支统计数据
				}
			}

			stocks := statistics.Group("/stock") // 库存统计
			{
				stocks.GET("/where", stock.StockController{}.Where) // 库存统计筛选
				stocks.Use(middlewares.JWTMiddleware())
				{
					stocks.POST("/data", stock.StockController{}.Data) // 库存统计数据
				}
			}

			sales := statistics.Group("/sales") // 销售统计
			{
				sales.GET("/where", sale.SaleController{}.Where) // 销售统计筛选
				sales.Use(middlewares.JWTMiddleware())
				{
					sales.POST("/data", sale.SaleController{}.Data) // 销售统计数据
				}
			}

			statistics.Use(middlewares.JWTMiddleware())
			{
				statistics.POST("/sales_detail_daily", statistic.StatisticController{}.SalesDetailDaily) // 销售明细日报
			}
		}

		// 工作台
		workbenchs := r.Group("/workbench")
		{
			workbenchs.Use(middlewares.JWTMiddleware())
			{
				workbenchs.GET("/list", workbench.WorkbenchController{}.List)      // 工作台列表
				workbenchs.POST("/search", workbench.WorkbenchController{}.Search) // 工作台搜索
				workbenchs.POST("/add", workbench.WorkbenchController{}.Add)       // 工作台添加
				workbenchs.PUT("/update", workbench.WorkbenchController{}.Update)  // 工作台更新
				workbenchs.DELETE("/del", workbench.WorkbenchController{}.Del)     // 工作台删除
			}
		}

		// 门店
		stores := r.Group("/store")
		{
			root := stores.Group("/")
			{
				root.GET("/where", store.StoreController{}.Where) // 门店筛选
				root.Use(middlewares.JWTMiddleware())
				{
					root.POST("/create", store.StoreController{}.Create)   // 创建门店
					root.PUT("/update", store.StoreController{}.Update)    // 门店更新
					root.DELETE("/delete", store.StoreController{}.Delete) // 门店删除
					root.POST("/list", store.StoreController{}.List)       // 门店列表
					root.POST("/alias", store.StoreController{}.Alias)     // 门店别名
					root.POST("/my", store.StoreController{}.My)           // 我的门店
					root.POST("/info", store.StoreController{}.Info)       // 门店详情
				}

			}

			staffs := stores.Group("/staff")
			{
				staffs.Use(middlewares.JWTMiddleware())
				{
					staffs.POST("/list", store.StoreStaffController{}.List)  // 门店员工列表
					staffs.POST("/is_in", store.StoreStaffController{}.IsIn) // 是否在门店
				}
			}

			superiors := stores.Group("/superior")
			{
				superiors.Use(middlewares.JWTMiddleware())
				{
					superiors.POST("/list", store.StoreSuperiorController{}.List)  // 门店负责人列表
					superiors.POST("/is_in", store.StoreSuperiorController{}.IsIn) // 是否是负责人
				}
			}

			admins := stores.Group("/admin")
			{
				admins.Use(middlewares.JWTMiddleware())
				{
					admins.POST("/list", store.StoreAdminController{}.List)  // 门店管理员列表
					admins.POST("/is_in", store.StoreAdminController{}.IsIn) // 是否是管理员
				}
			}
		}

		// 区域
		regions := r.Group("/region")
		{
			root := regions.Group("/")
			{
				root.GET("/where", region.RegionController{}.Where) // 区域筛选
				root.Use(middlewares.JWTMiddleware())
				{
					root.POST("/create", region.RegionController{}.Create)   // 创建区域
					root.PUT("/update", region.RegionController{}.Update)    // 区域更新
					root.DELETE("/delete", region.RegionController{}.Delete) // 区域删除
					root.POST("/list", region.RegionController{}.List)       // 区域列表
					root.POST("/my", region.RegionController{}.My)           // 我的区域
					root.POST("/info", region.RegionController{}.Info)       // 区域详情
				}

			}

			stores := regions.Group("/store")
			{
				stores.Use(middlewares.JWTMiddleware())
				{
					stores.POST("/list", region.RegionStoreController{}.List) // 区域门店列表
				}
			}

			staffs := regions.Group("/staff")
			{
				staffs.Use(middlewares.JWTMiddleware())
				{
					staffs.POST("/list", region.RegionStaffController{}.List) // 区域员工列表
				}
			}

			superiors := regions.Group("/superior")
			{
				superiors.Use(middlewares.JWTMiddleware())
				{
					superiors.POST("/list", region.RegionSuperiorController{}.List) // 区域负责人列表
				}
			}

			admins := regions.Group("/admin")
			{
				admins.Use(middlewares.JWTMiddleware())
				{
					admins.POST("/list", region.RegionAdminController{}.List) // 区域管理员列表
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
						finished.POST("/list", product.ProductFinishedController{}.List)              // 成品列表
						finished.POST("/empty_image", product.ProductFinishedController{}.EmptyImage) // 空图片列表

						finished.POST("/info", product.ProductFinishedController{}.Info)           // 成品详情
						finished.POST("/retrieval", product.ProductFinishedController{}.Retrieval) // 成品检索
						finished.PUT("/update", product.ProductFinishedController{}.Update)        // 成品更新
						finished.PUT("/upload", product.ProductFinishedController{}.Upload)        // 成品图上传
						batch := finished.Group("/batch")                                          // 批量操作
						{
							batch.PUT("/update", product.ProductFinishedBatchController{}.Update)          // 批量更新
							batch.PUT("/update_code", product.ProductFinishedBatchController{}.UpdateCode) // 批量更新条码
							batch.PUT("/find_code", product.ProductFinishedBatchController{}.FindCode)     // 批量查找条码
						}
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

						enters.POST("/add_product", product.ProductFinishedEnterController{}.AddProduct)       // 添加产品
						enters.PUT("/edit_product", product.ProductFinishedEnterController{}.EditProduct)      // 编辑产品
						enters.DELETE("/del_product", product.ProductFinishedEnterController{}.DelProduct)     // 删除产品
						enters.DELETE("/clear_product", product.ProductFinishedEnterController{}.ClearProduct) // 清空产品

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
					old.GET("/where", product.ProductOldController{}.Where)              // 旧料筛选
					old.GET("/where_create", product.ProductOldController{}.WhereCreate) // 旧料筛选
					old.POST("/get_class", product.ProductOldController{}.GetClass)      // 获取旧料分类
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

				// 配件入库
				enters := accessories.Group("/enter")
				{
					enters.GET("/where", product.ProductAccessorieEnterController{}.Where)                       // 配件入库单筛选
					enters.GET("/where_add_product", product.ProductAccessorieEnterController{}.WhereAddProduct) // 配件入库单筛选
					enters.Use(middlewares.JWTMiddleware())
					{
						enters.POST("/create", product.ProductAccessorieEnterController{}.Create) // 创建配件入库单
						enters.POST("/list", product.ProductAccessorieEnterController{}.List)     // 配件入库单列表
						enters.POST("/info", product.ProductAccessorieEnterController{}.Info)     // 配件入库单详情

						enters.POST("/add_product", product.ProductAccessorieEnterController{}.AddProduct)       // 添加产品
						enters.DELETE("/del_product", product.ProductAccessorieEnterController{}.DelProduct)     // 删除产品
						enters.DELETE("/clear_product", product.ProductAccessorieEnterController{}.ClearProduct) // 清空产品
						enters.PUT("/edit_product", product.ProductAccessorieEnterController{}.EditProduct)      // 编辑产品

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
						allocate.POST("/details", product.ProductAccessorieAllocateController{}.Details)  // 调拨单明细
						allocate.POST("/info", product.ProductAccessorieAllocateController{}.Info)        // 调拨单详情
						allocate.PUT("/add", product.ProductAccessorieAllocateController{}.Add)           // 添加产品
						allocate.PUT("/remove", product.ProductAccessorieAllocateController{}.Remove)     // 移除产品
						allocate.PUT("/clear", product.ProductAccessorieAllocateController{}.Clear)       // 清空产品
						allocate.PUT("/confirm", product.ProductAccessorieAllocateController{}.Confirm)   // 确认调拨
						allocate.PUT("/cancel", product.ProductAccessorieAllocateController{}.Cancel)     // 取消调拨
						allocate.PUT("/complete", product.ProductAccessorieAllocateController{}.Complete) // 完成调拨
					}
				}
			}

			// 货品调拨
			allocate := products.Group("/allocate")
			{
				allocate.GET("/where", product.ProductAllocateController{}.Where) // 调拨单筛选
				allocate.Use(middlewares.JWTMiddleware())
				{
					allocate.POST("/create", product.ProductAllocateController{}.Create)              // 创建调拨单
					allocate.POST("/list", product.ProductAllocateController{}.List)                  // 调拨单列表
					allocate.POST("/details", product.ProductAllocateController{}.Details)            // 调拨单明细
					allocate.POST("/info", product.ProductAllocateController{}.Info)                  // 调拨单详情
					allocate.POST("/info_overview", product.ProductAllocateController{}.InfoOverview) // 调拨单概览
					allocate.PUT("/add", product.ProductAllocateController{}.Add)                     // 添加产品
					allocate.DELETE("/remove", product.ProductAllocateController{}.Remove)            // 移除产品
					allocate.DELETE("/clear", product.ProductAllocateController{}.Clear)              // 清空产品
					allocate.PUT("/confirm", product.ProductAllocateController{}.Confirm)             // 确认调拨
					allocate.PUT("/cancel", product.ProductAllocateController{}.Cancel)               // 取消调拨
					allocate.PUT("/complete", product.ProductAllocateController{}.Complete)           // 完成调拨
				}
			}

			// 货品盘点
			inventory := products.Group("/inventory")
			{
				inventory.GET("/where", product.ProductInventoryController{}.Where) // 盘点单筛选
				inventory.Use(middlewares.JWTMiddleware())
				{
					inventory.POST("/create", product.ProductInventoryController{}.Create)      // 创建盘点单
					inventory.POST("/list", product.ProductInventoryController{}.List)          // 盘点单列表
					inventory.POST("/info", product.ProductInventoryController{}.Info)          // 盘点单详情
					inventory.POST("/add", product.ProductInventoryController{}.Add)            // 添加产品
					inventory.POST("/add_batch", product.ProductInventoryController{}.AddBatch) // 批量添加产品
					inventory.PUT("/remove", product.ProductInventoryController{}.Remove)       // 移除产品
					inventory.PUT("/change", product.ProductInventoryController{}.Change)       // 盘点单变化
				}
			}

			// 货品记录
			historys := products.Group("/history")
			{
				historys.GET("/where", product.ProductHistoryController{}.Where)                      // 货品操作记录筛选
				historys.GET("/where_accessorie", product.ProductHistoryController{}.WhereAccessorie) // 配件操作记录筛选
				historys.Use(middlewares.JWTMiddleware())
				{
					historys.POST("/list", product.ProductHistoryController{}.List)                      // 货品操作记录列表
					historys.POST("/list_accessorie", product.ProductHistoryController{}.ListAccessorie) // 配件操作记录列表
					historys.POST("/info", product.ProductHistoryController{}.Info)                      // 货品操作记录详情
				}
			}
		}

		// 会员
		members := r.Group("/member")
		{
			root := members.Group("/")
			{
				root.GET("/where", member.MemberController{}.Where) // 会员筛选
				root.Use(middlewares.JWTMiddleware())
				{
					root.POST("/create", member.MemberController{}.Create)             // 创建会员
					root.POST("/list", member.MemberController{}.List)                 // 会员列表
					root.POST("/info", member.MemberController{}.Info)                 // 会员详情
					root.PUT("/update", member.MemberController{}.Update)              // 会员更新
					root.POST("/consumptions", member.MemberController{}.Consumptions) // 会员消费记录
				}
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
			// 销售单
			sales := orders.Group("/sales")
			{
				root := sales.Group("/")
				{
					root.GET("/where", order.OrderSalesController{}.Where) // 订单筛选
					root.Use(middlewares.JWTMiddleware())
					{
						root.POST("/create", order.OrderSalesController{}.Create)  // 创建订单
						root.POST("/list", order.OrderSalesController{}.List)      // 订单列表
						root.POST("/info", order.OrderSalesController{}.Info)      // 订单详情
						root.PUT("/revoked", order.OrderSalesController{}.Revoked) // 订单撤销
						root.PUT("/pay", order.OrderSalesController{}.Pay)         // 订单支付
						root.PUT("/refund", order.OrderSalesController{}.Refund)   // 退货
					}
				}

				details := sales.Group("/detail")
				{
					details.GET("/where", order.OrderSalesDetailController{}.Where) // 订单筛选
					details.Use(middlewares.JWTMiddleware())
					{
						details.POST("/list", order.OrderSalesDetailController{}.List) // 订单列表
						details.POST("/info", order.OrderSalesDetailController{}.Info) // 订单详情
					}
				}

				refunds := sales.Group("/refund")
				{
					refunds.GET("/where", order.OrderSalesRefundController{}.Where) // 订单筛选
					refunds.Use(middlewares.JWTMiddleware())
					{
						refunds.POST("/list", order.OrderSalesRefundController{}.List) // 订单列表
					}
				}
			}

			// 定金单
			deposits := orders.Group("/deposit")
			{
				deposits.GET("/where", order.OrderDepositController{}.Where) // 订单筛选
				deposits.Use(middlewares.JWTMiddleware())
				{
					deposits.POST("/create", order.OrderDepositController{}.Create)  // 创建订单
					deposits.POST("/list", order.OrderDepositController{}.List)      // 订单列表
					deposits.POST("/info", order.OrderDepositController{}.Info)      // 订单详情
					deposits.PUT("/revoked", order.OrderDepositController{}.Revoked) // 订单撤销
					deposits.PUT("/pay", order.OrderDepositController{}.Pay)         // 订单支付
					deposits.PUT("/refund", order.OrderDepositController{}.Refund)   // 退货
				}
			}

			// 维修单
			repairs := orders.Group("/repair")
			{
				repairs.GET("/where", order.OrderRepairController{}.Where)                // 订单筛选
				repairs.GET("/where_product", order.OrderRepairController{}.WhereProduct) // 订单筛选
				repairs.Use(middlewares.JWTMiddleware())
				{
					repairs.POST("/create", order.OrderRepairController{}.Create)      // 创建订单
					repairs.POST("/list", order.OrderRepairController{}.List)          // 订单列表
					repairs.POST("/info", order.OrderRepairController{}.Info)          // 订单详情
					repairs.PUT("/update", order.OrderRepairController{}.Update)       // 订单修改
					repairs.PUT("/operation", order.OrderRepairController{}.Operation) // 订单操作
					repairs.PUT("/revoked", order.OrderRepairController{}.Revoked)     // 订单撤销
					repairs.PUT("/pay", order.OrderRepairController{}.Pay)             // 订单支付
					repairs.PUT("/refund", order.OrderRepairController{}.Refund)       // 退款
				}
			}

			// 其他收支单
			other := orders.Group("/other")
			{
				other.GET("/where", order.OrderOtherController{}.Where) // 订单筛选
				other.Use(middlewares.JWTMiddleware())
				{
					other.POST("/create", order.OrderOtherController{}.Create)   // 创建订单
					other.POST("/list", order.OrderOtherController{}.List)       // 订单列表
					other.POST("/info", order.OrderOtherController{}.Info)       // 订单详情
					other.PUT("/update", order.OrderOtherController{}.Update)    // 订单修改
					other.DELETE("/delete", order.OrderOtherController{}.Delete) // 订单删除
				}
			}
		}

		// 销售目标
		targets := r.Group("/target")
		{
			targets.GET("/where", target.TargetController{}.Where)                  // 销售目标筛选
			targets.GET("/where_group", target.TargetController{}.WhereGroup)       // 销售目标分组筛选
			targets.GET("/where_personal", target.TargetController{}.WherePersonal) // 销售目标个人筛选
			targets.Use(middlewares.JWTMiddleware())
			{
				targets.POST("/create", target.TargetController{}.Create) // 创建销售目标
				targets.POST("/list", target.TargetController{}.List)     // 销售目标列表
			}
		}

		// 设置
		settings := r.Group("/setting")
		{
			// 金价设置
			gold_price := settings.Group("/gold_price")
			{
				gold_price.Use(middlewares.JWTMiddleware())
				{
					gold_price.POST("/list", setting.GoldPriceController{}.List)     // 金价历史列表
					gold_price.POST("/create", setting.GoldPriceController{}.Create) // 创建金价
				}
			}

			// 开单设置
			open_orders := settings.Group("/open_order")
			{
				open_orders.GET("/where", setting.OpenOrderController{}.Where) // 开单设置筛选
				open_orders.Use(middlewares.JWTMiddleware())
				{
					open_orders.POST("/info", setting.OpenOrderController{}.Info)    // 开单设置详情
					open_orders.PUT("/update", setting.OpenOrderController{}.Update) // 开单设置更新
				}
			}

			// 角色权限
			roles := settings.Group("/role")
			{
				roles.GET("/where", setting.RoleController{}.Where) // 角色权限筛选
				roles.Use(middlewares.JWTMiddleware())
				{
					roles.GET("/identity", setting.RoleController{}.GetIdentity) // 获取当前用户角色权限
					roles.POST("/create", setting.RoleController{}.Create)       // 创建角色
					roles.POST("/copy", setting.RoleController{}.Copy)           // 角色权限复制
					roles.POST("/list", setting.RoleController{}.List)           // 角色权限列表
					roles.POST("/info", setting.RoleController{}.Info)           // 角色权限详情
					roles.PUT("/edit", setting.RoleController{}.Edit)            // 角色权限编辑
					roles.PUT("/update", setting.RoleController{}.Update)        // 角色权限更新
					roles.DELETE("/delete", setting.RoleController{}.Delete)     // 角色权限删除

					roles.POST("/apis", setting.RoleController{}.Apis) // 角色权限API列表
				}
			}

			// 打印设置
			print_settings := settings.Group("/print")
			{
				print_settings.Use(middlewares.JWTMiddleware())
				{
					print_settings.POST("/create", setting.PrintController{}.Create)   // 创建打印设置
					print_settings.POST("/list", setting.PrintController{}.List)       // 打印设置列表
					print_settings.POST("/info", setting.PrintController{}.Info)       // 打印设置详情
					print_settings.PUT("/update", setting.PrintController{}.Update)    // 打印设置更新
					print_settings.DELETE("/delete", setting.PrintController{}.Delete) // 打印设置删除
					print_settings.PUT("/copy", setting.PrintController{}.Copy)        // 打印设置复制
				}
			}

			// 常用备注
			remarks := settings.Group("/remark")
			{
				remarks.GET("/where", setting.RemarkController{}.Where) // 常用备注筛选
				remarks.Use(middlewares.JWTMiddleware())
				{
					remarks.POST("/create", setting.RemarkController{}.Create)   // 创建备注
					remarks.POST("/list", setting.RemarkController{}.List)       // 备注列表
					remarks.PUT("/update", setting.RemarkController{}.Update)    // 备注更新
					remarks.DELETE("/delete", setting.RemarkController{}.Delete) // 备注删除
				}
			}
		}
	}
}
