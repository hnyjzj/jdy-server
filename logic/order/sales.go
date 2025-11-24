package order

import (
	"fmt"
	"jdy/enums"
	"jdy/errors"
	"jdy/logic/product"
	"jdy/message"
	"jdy/model"
	"jdy/types"
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderSalesLogic struct {
	Ctx   *gin.Context
	Staff *model.Staff
}

func (l *OrderSalesLogic) List(req *types.OrderSalesListReq) (*types.PageRes[model.OrderSales], error) {
	var (
		order model.OrderSales

		res types.PageRes[model.OrderSales]
	)

	db := model.DB.Model(&order)
	db = order.WhereCondition(db, &req.Where)

	// 获取总数
	if err := db.Count(&res.Total).Error; err != nil {
		return nil, errors.New("获取订单总数失败")
	}

	// 获取列表
	db = order.Preloads(db)

	db = db.Order("created_at desc")
	db = model.PageCondition(db, &req.PageReq)
	if err := db.Find(&res.List).Error; err != nil {
		return nil, errors.New("获取订单列表失败")
	}

	return &res, nil
}

func (l *OrderSalesLogic) Info(req *types.OrderSalesInfoReq) (*model.OrderSales, error) {
	var (
		order model.OrderSales
	)

	db := model.DB.Model(&order)

	db = order.Preloads(db)
	db = db.Preload("OrderRefunds")

	db = db.Where("id = ?", req.Id)
	if err := db.First(&order).Error; err != nil {
		return nil, errors.New("获取订单详情失败")
	}

	return &order, nil
}

// 撤销订单
func (l *OrderSalesLogic) Revoked(req *types.OrderSalesRevokedReq) error {
	var (
		order model.OrderSales
	)

	db := model.DB.Model(&order)

	db = db.Where("id = ?", req.Id)
	db = order.Preloads(db)
	if err := db.First(&order).Error; err != nil {
		return errors.New("获取订单详情失败")
	}

	if order.Status != enums.OrderSalesStatusWaitPay {
		return errors.New("订单状态不正确")
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 更新订单状态
		if err := tx.Model(&model.OrderSales{}).Where("id = ?", order.Id).Updates(&model.OrderSales{
			Status: enums.OrderSalesStatusCancel,
		}).Error; err != nil {
			return errors.New("撤销订单失败")
		}

		// 处理订单产品
		for _, product := range order.Products {
			// 更新订单状态
			if err := tx.Model(&model.OrderSalesProduct{}).Where("id = ?", product.Id).Updates(&model.OrderSalesProduct{
				Status: enums.OrderSalesStatusCancel,
			}).Error; err != nil {
				return errors.New("更新订单成品状态失败")
			}
			switch product.Type {
			case enums.ProductTypeFinished:
				log := model.ProductHistory{
					Action:     enums.ProductActionOrderCancel,
					Type:       enums.ProductTypeFinished,
					OldValue:   product.Finished.Product,
					ProductId:  product.Finished.Product.Id,
					StoreId:    product.Finished.Product.StoreId,
					SourceId:   product.Finished.Product.Id,
					OperatorId: l.Staff.Id,
					IP:         l.Ctx.ClientIP(),
				}
				// 更新成品状态
				if err := tx.Model(&model.ProductFinished{}).Where("id = ?", product.Finished.Product.Id).Updates(&model.ProductFinished{
					Status: enums.ProductStatusNormal,
				}).Error; err != nil {
					return errors.New("更新成品状态失败")
				}
				// 添加成品历史记录
				product.Finished.Product.Status = enums.ProductStatusNormal
				log.NewValue = product.Finished
				if err := tx.Create(&log).Error; err != nil {
					return errors.New("添加成品历史记录失败")
				}
			case enums.ProductTypeOld:
				log := model.ProductHistory{
					Action:     enums.ProductActionOrderCancel,
					Type:       enums.ProductTypeOld,
					OldValue:   product.Old.Product,
					ProductId:  product.Old.Product.Id,
					StoreId:    product.Old.Product.StoreId,
					SourceId:   product.Old.Product.Id,
					OperatorId: l.Staff.Id,
					IP:         l.Ctx.ClientIP(),
				}
				// 更新旧料状态
				if err := tx.Model(&model.ProductOld{}).Where("id = ?", product.Old.Product.Id).Updates(&model.ProductOld{
					Status: enums.ProductStatusDraft,
				}).Error; err != nil {
					return errors.New("更新旧料状态失败")
				}
				// 添加旧料历史记录
				product.Old.Product.Status = enums.ProductStatusDraft
				log.NewValue = product.Old
				if err := tx.Create(&log).Error; err != nil {
					return errors.New("添加旧料历史记录失败")
				}
			case enums.ProductTypeAccessorie:
				log := model.ProductHistory{
					Action:     enums.ProductActionOrderCancel,
					Type:       enums.ProductTypeAccessorie,
					OldValue:   product.Accessorie.Product,
					ProductId:  product.Accessorie.ProductId,
					StoreId:    product.Accessorie.Product.StoreId,
					SourceId:   product.Accessorie.ProductId,
					OperatorId: l.Staff.Id,
					IP:         l.Ctx.ClientIP(),
				}
				// 更新配件状态
				if err := tx.Model(&model.ProductAccessorie{}).Where("id = ?", product.Accessorie.ProductId).Updates(&model.ProductAccessorie{
					Status: enums.ProductAccessorieStatusNormal,
				}).Update("stock", gorm.Expr("stock + ?", product.Accessorie.Quantity)).Error; err != nil {
					return errors.New("更新配件状态失败")
				}
				// 添加配件历史记录
				log.NewValue = product.Accessorie.Product
				if err := tx.Create(&log).Error; err != nil {
					return errors.New("添加配件历史记录失败")
				}
			}
		}

		// 处理定金单
		for _, deposit := range order.OrderDeposits {
			// 循环产品
			for _, product := range deposit.Products {
				old_product := product
				if old_product.IsOur {
					log := model.ProductHistory{
						Action:     enums.ProductActionOrderCancel,
						Type:       enums.ProductTypeFinished,
						OldValue:   old_product.ProductFinished,
						ProductId:  old_product.ProductFinished.Id,
						StoreId:    old_product.ProductFinished.StoreId,
						SourceId:   order.Id,
						OperatorId: l.Staff.Id,
						IP:         l.Ctx.ClientIP(),
					}
					// 更新商品状态
					if err := tx.Model(&model.ProductFinished{}).Where("id = ?", old_product.ProductFinished.Id).Updates(model.ProductFinished{
						Status: enums.ProductStatusReturn,
					}).Error; err != nil {
						return errors.New("配件状态更新失败")
					}
					// 添加记录
					old_product.ProductFinished.Status = enums.ProductStatusReturn
					log.NewValue = old_product.ProductFinished
					if err := tx.Create(&log).Error; err != nil {
						return errors.New("配件记录添加失败")
					}
				}
			}

			// 更新定金单状态
			if err := tx.Model(&model.OrderDeposit{}).Where("id = ?", deposit.Id).Updates(model.OrderDeposit{
				Status: enums.OrderDepositStatusBooking,
			}).Error; err != nil {
				return errors.New("更新定金单状态失败")
			}

			// 清空与销售单的关联
			if err := tx.Model(&deposit).Association("OrderSales").Clear(); err != nil {
				return errors.New("清空定金单关联失败")
			}
		}

		return nil
	}); err != nil {
		return errors.New("撤销订单失败")
	}

	// 发送通知
	go func() {
		msg := message.NewMessage(l.Ctx)
		order.Operator = *l.Staff
		if err := msg.SendOrderSalesCancelMessage(&message.OrderSalesMessage{
			OrderSales: &order,
		}); err != nil {
			log.Printf("发送订单撤销通知失败: %v", err)
		}
	}()

	return nil
}

// 支付
func (l *OrderSalesLogic) Pay(req *types.OrderSalesPayReq) error {
	var (
		order model.OrderSales
	)

	db := model.DB.Model(&order)

	db = db.Where("id = ?", req.Id)
	db = order.Preloads(db)
	if err := db.First(&order).Error; err != nil {
		return errors.New("获取订单详情失败")
	}

	if order.Status != enums.OrderSalesStatusWaitPay {
		return errors.New("订单状态不正确")
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 支付成品
		for _, product := range order.Products {
			// 更新订单状态
			if err := tx.Model(&model.OrderSalesProduct{}).Where("id = ?", product.Id).Updates(&model.OrderSalesProduct{
				Status: enums.OrderSalesStatusComplete,
			}).Error; err != nil {
				return errors.New("更新订单成品状态失败")
			}
			switch product.Type {
			case enums.ProductTypeFinished:
				{
					// 更新成品状态
					if err := tx.Model(&model.ProductFinished{}).Where("id = ?", product.Finished.Product.Id).Updates(&model.ProductFinished{
						Status: enums.ProductStatusSold,
					}).Error; err != nil {
						return errors.New("更新成品状态失败")
					}
				}
			case enums.ProductTypeOld:
				{
					// 更新旧料状态
					if err := tx.Model(&model.ProductOld{}).Where("id = ?", product.Old.Product.Id).Updates(&model.ProductOld{
						Status: enums.ProductStatusNormal,
					}).Error; err != nil {
						return errors.New("更新旧料状态失败")
					}
				}
			}
		}

		for _, payment := range order.Payments {
			// 更新支付状态
			if err := tx.Model(&model.OrderPayment{}).Where("id = ?", payment.Id).Updates(&model.OrderPayment{
				Status: true,
			}).Error; err != nil {
				return errors.New("更新支付状态失败")
			}
		}

		// 更新订单状态
		if err := tx.Model(&model.OrderSales{}).Where("id = ?", order.Id).Updates(&model.OrderSales{
			Status: enums.OrderSalesStatusComplete,
		}).Error; err != nil {
			return errors.New("支付订单失败")
		}

		return nil
	}); err != nil {
		return errors.New("支付订单失败")
	}

	// 发送通知
	go func() {
		msg := message.NewMessage(l.Ctx)
		order.Operator = *l.Staff
		if err := msg.SendOrderSalesPayMessage(&message.OrderSalesMessage{
			OrderSales: &order,
		}); err != nil {
			log.Printf("发送订单支付通知失败: %v", err)
		}
	}()

	return nil
}

// 退货
func (l *OrderSalesLogic) Refund(req *types.OrderSalesRefundReq) error {
	var (
		order model.OrderSales
	)

	db := model.DB.Model(&order)
	db = order.Preloads(db)
	db = db.Where("id = ?", req.Id)
	if err := db.First(&order).Error; err != nil {
		return errors.New("获取订单详情失败")
	}

	if order.Status != enums.OrderSalesStatusComplete && order.Status != enums.OrderSalesStatusRefund {
		return errors.New("订单状态不正确")
	}

	data := model.OrderRefund{
		StoreId:    order.StoreId,
		OrderId:    order.Id,
		OrderType:  enums.OrderTypeSales,
		MemberId:   order.MemberId,
		Remark:     req.Remark,
		OperatorId: l.Staff.Id,
		IP:         l.Ctx.ClientIP(),
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 查询产品
		switch req.ProductType {
		case enums.ProductTypeFinished:
			{
				// 查询成品
				var p model.OrderSalesProduct
				if err := tx.Preload("Finished.Product").First(&p, "id = ?", req.ProductId).Error; err != nil {
					return errors.New("获取成品订单详情失败")
				}

				if p.Status != enums.OrderSalesStatusComplete {
					return errors.New("成品订单状态不正确")
				}

				// 更新订单成品状态
				if err := tx.Model(&model.OrderSalesProduct{}).Where("id = ?", p.Id).Updates(&model.OrderSalesProduct{
					Status: enums.OrderSalesStatusReturn,
				}).Error; err != nil {
					return errors.New("更新订单成品状态失败")
				}

				// 根据入库类型更新成品
				switch req.Method {
				case enums.ProductTypeUsedFinished:
					// 添加历史
					log := model.ProductHistory{
						Action:     enums.ProductActionReturn,
						Type:       enums.ProductTypeFinished,
						OldValue:   p.Finished.Product,
						ProductId:  p.Finished.Product.Id,
						StoreId:    p.Finished.Product.StoreId,
						SourceId:   p.Finished.Product.Id,
						OperatorId: l.Staff.Id,
						IP:         l.Ctx.ClientIP(),
					}
					if err := tx.Model(&model.ProductFinished{}).Where("id = ?", p.Finished.Product.Id).Updates(&model.ProductFinished{
						Status: enums.ProductStatusNormal,
					}).Error; err != nil {
						return errors.New("更新成品状态失败")
					}

					p.Finished.Product.Status = enums.ProductStatusNormal
					log.NewValue = p.Finished.Product
					if err := tx.Create(&log).Error; err != nil {
						return errors.New("创建成品历史失败")
					}

				case enums.ProductTypeUsedOld:
					if err := tx.Model(&model.ProductFinished{}).Where("id = ?", p.Finished.Product.Id).Updates(&model.ProductFinished{
						Status: enums.ProductStatusDamage,
					}).Error; err != nil {
						return errors.New("更新成品状态失败")
					}

					// 成品转旧料
					damage := product.ProductFinishedDamageLogic{
						Ctx:   l.Ctx,
						Staff: l.Staff,
					}
					if err := damage.Conversion(&types.ProductConversionReq{
						Id:     p.Finished.Product.Id,
						Type:   enums.ProductTypeUsedOld,
						Remark: fmt.Sprintf("销售单退货(%s): %s", order.Id, req.Remark),
					}); err != nil {
						return err
					}
				}

				data.Type = enums.ProductTypeFinished
				data.Code = strings.TrimSpace(strings.ToUpper(p.Finished.Product.Code))
				data.Name = p.Finished.Product.Name
				data.Quantity = 1
				data.Price = req.Price
				data.PriceOriginal = p.Finished.Price
			}

		case enums.ProductTypeOld:
			{
				// 查询旧料
				var p model.OrderSalesProduct
				if err := tx.Preload("Old.Product").First(&p, "id = ?", req.ProductId).Error; err != nil {
					return errors.New("获取旧料订单详情失败")
				}
				if p.Status != enums.OrderSalesStatusComplete {
					return errors.New("旧料订单状态不正确")
				}
				if p.Old.Product.Status != enums.ProductStatusNormal {
					return errors.New("旧料已不在库，已无法退货")
				}

				// 更新订单旧料状态
				if err := tx.Model(&model.OrderSalesProduct{}).Where("id = ?", p.Id).Updates(&model.OrderSalesProduct{
					Status: enums.OrderSalesStatusReturn,
				}).Error; err != nil {
					return errors.New("更新订单旧料状态失败")
				}

				// 添加历史
				log := model.ProductHistory{
					Action:     enums.ProductActionReturn,
					Type:       enums.ProductTypeOld,
					OldValue:   p.Old.Product,
					ProductId:  p.Old.Product.Id,
					StoreId:    p.Old.Product.StoreId,
					SourceId:   p.Old.Product.Id,
					OperatorId: l.Staff.Id,
					IP:         l.Ctx.ClientIP(),
				}
				// 更新旧料状态
				if err := tx.Model(&model.ProductOld{}).Where("id = ?", p.Old.Product.Id).Updates(&model.ProductOld{
					Status: enums.ProductStatusNoStock,
				}).Error; err != nil {
					return errors.New("更新旧料状态失败")
				}

				p.Old.Product.Status = enums.ProductStatusNoStock
				log.NewValue = p.Old.Product
				if err := tx.Create(&log).Error; err != nil {
					return errors.New("创建旧料历史失败")
				}

				data.Type = enums.ProductTypeOld
				data.Code = strings.TrimSpace(strings.ToUpper(p.Old.Product.Code))
				data.Name = p.Old.Product.Name
				data.Quantity = 1
				data.Price = req.Price
				data.PriceOriginal = p.Old.RecyclePrice
			}
		case enums.ProductTypeAccessorie:
			{
				// 查询配件
				var p model.OrderSalesProduct
				if err := tx.Preload("Accessorie.Product").First(&p, "id = ?", req.ProductId).Error; err != nil {
					return errors.New("获取配件订单详情失败")
				}
				if p.Status != enums.OrderSalesStatusComplete {
					return errors.New("配件订单状态不正确")
				}

				// 更新订单配件状态
				if err := tx.Model(&model.OrderSalesProduct{}).Where("id = ?", p.Id).Updates(&model.OrderSalesProduct{
					Status: enums.OrderSalesStatusReturn,
				}).Error; err != nil {
					return errors.New("更新订单配件状态失败")
				}
				// 添加历史
				log := model.ProductHistory{
					Action:     enums.ProductActionReturn,
					Type:       enums.ProductTypeAccessorie,
					OldValue:   p.Accessorie.Product,
					ProductId:  p.Accessorie.ProductId,
					StoreId:    p.Accessorie.Product.StoreId,
					SourceId:   p.Accessorie.ProductId,
					OperatorId: l.Staff.Id,
					IP:         l.Ctx.ClientIP(),
				}

				// 更新配件状态
				if err := tx.Model(&model.ProductAccessorie{}).Where("id = ?", p.Accessorie.Product.Id).Updates(&model.ProductAccessorie{
					Status: enums.ProductAccessorieStatusNormal,
				}).Update("stock", gorm.Expr("stock + ?", p.Accessorie.Quantity)).Error; err != nil {
					return errors.New("更新配件状态失败")
				}

				p.Accessorie.Product.Status = enums.ProductAccessorieStatusNormal
				p.Accessorie.Product.Stock += p.Accessorie.Quantity

				log.NewValue = p.Accessorie.Product
				if err := tx.Create(&log).Error; err != nil {
					return errors.New("创建配件历史失败")
				}

				data.Type = enums.ProductTypeAccessorie
				data.Name = p.Accessorie.Product.Name
				data.Quantity = p.Accessorie.Quantity
				data.Price = req.Price
				data.PriceOriginal = p.Accessorie.Price
			}
		}

		if err := tx.Model(&model.OrderSales{}).Where("id = ?", order.Id).Updates(&model.OrderSales{
			Status: enums.OrderSalesStatusRefund,
		}).Error; err != nil {
			return errors.New("更新订单状态失败")
		}

		if err := tx.Create(&data).Error; err != nil {
			return errors.New("创建退货单失败")
		}

		switch req.ProductType {
		case enums.ProductTypeOld:
			{
				// 添加退款方式
				for _, payment := range req.Payments {
					payment := model.OrderPayment{
						Status:        true,
						StoreId:       order.StoreId,
						OrderId:       data.Id,
						Type:          enums.FinanceTypeIncome,
						Source:        enums.FinanceSourceSaleRefund,
						OrderType:     enums.OrderTypeReturn,
						PaymentMethod: payment.PaymentMethod,
						Amount:        payment.Amount,
					}
					if err := tx.Create(&payment).Error; err != nil {
						return errors.New("添加退款方式失败")
					}
				}
			}
		default:
			{
				// 添加退款方式
				for _, payment := range req.Payments {
					payment := model.OrderPayment{
						Status:        true,
						StoreId:       order.StoreId,
						OrderId:       data.Id,
						Type:          enums.FinanceTypeExpense,
						Source:        enums.FinanceSourceSaleRefund,
						OrderType:     enums.OrderTypeReturn,
						PaymentMethod: payment.PaymentMethod,
						Amount:        payment.Amount,
					}
					if err := tx.Create(&payment).Error; err != nil {
						return errors.New("添加退款方式失败")
					}
				}
			}
		}

		return nil
	}); err != nil {
		return errors.New(err.Error())
	}

	// 发送通知
	go func() {
		msg := message.NewMessage(l.Ctx)
		order.Operator = *l.Staff
		if err := msg.SendOrderSalesRefundMessage(&message.OrderSalesMessage{
			OrderSales: &order,
		}); err != nil {
			log.Printf("发送通知失败: %s", err.Error())
		}
	}()

	return nil
}

func (l *OrderSalesLogic) Retreat(req *types.OrderSalesRetreatReq) error {
	var (
		order model.OrderSales
	)

	db := model.DB.Model(&order)
	db = order.Preloads(db)
	db = db.Where("id = ?", req.Id)
	if err := db.First(&order).Error; err != nil {
		return errors.New("获取订单详情失败")
	}

	if order.Status != enums.OrderSalesStatusComplete {
		return errors.New("订单状态不正确")
	}

	if order.CreatedAt.Format(time.DateOnly) != time.Now().Format(time.DateOnly) {
		return errors.New("只能退当天订单")
	}

	if order.OperatorId != l.Staff.Id && order.CashierId != l.Staff.Id {
		return errors.New("订单只能由开单人或收银员退回")
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.OrderSales{}).Where("id = ?", order.Id).Updates(&model.OrderSales{
			Status: enums.OrderSalesStatusRevoke,
		}).Error; err != nil {
			return errors.New("更新订单状态失败")
		}

		// 处理订单产品
		for _, p := range order.Products {
			if p.Status != enums.OrderSalesStatusComplete {
				return errors.New("产品订单状态不正确")
			}
			// 更新订单成品状态
			if err := tx.Model(&model.OrderSalesProduct{}).Where("id = ?", p.Id).Updates(&model.OrderSalesProduct{
				Status: enums.OrderSalesStatusRevoke,
			}).Error; err != nil {
				return errors.New("更新订单成品状态失败")
			}

			switch p.Type {
			case enums.ProductTypeFinished:
				{
					// 添加历史
					log := model.ProductHistory{
						Action:     enums.ProductActionOrderCancel,
						Type:       enums.ProductTypeFinished,
						OldValue:   p.Finished.Product,
						ProductId:  p.Finished.Product.Id,
						StoreId:    p.Finished.Product.StoreId,
						SourceId:   p.Finished.Product.Id,
						OperatorId: l.Staff.Id,
						IP:         l.Ctx.ClientIP(),
					}
					if err := tx.Model(&model.ProductFinished{}).Where("id = ?", p.Finished.Product.Id).Updates(&model.ProductFinished{
						Status: enums.ProductStatusNormal,
					}).Error; err != nil {
						return errors.New("更新成品状态失败")
					}

					p.Finished.Product.Status = enums.ProductStatusNormal
					log.NewValue = p.Finished.Product
					if err := tx.Create(&log).Error; err != nil {
						return errors.New("创建成品历史失败")
					}
				}
			case enums.ProductTypeOld:
				{
					// 添加历史
					log := model.ProductHistory{
						Action:     enums.ProductActionOrderCancel,
						Type:       enums.ProductTypeOld,
						OldValue:   p.Old.Product,
						ProductId:  p.Old.Product.Id,
						StoreId:    p.Old.Product.StoreId,
						SourceId:   p.Old.Product.Id,
						OperatorId: l.Staff.Id,
						IP:         l.Ctx.ClientIP(),
					}

					if err := tx.Model(&model.ProductOld{}).Where("id = ?", p.Old.Product.Id).Updates(&model.ProductOld{
						Status: enums.ProductStatusNoStock,
					}).Error; err != nil {
						return errors.New("更新旧料状态失败")
					}

					p.Old.Product.Status = enums.ProductStatusNoStock
					log.NewValue = p.Old.Product
					if err := tx.Create(&log).Error; err != nil {
						return errors.New("创建旧料历史失败")
					}
				}
			case enums.ProductTypeAccessorie:
				{
					// 添加历史
					log := model.ProductHistory{
						Action:     enums.ProductActionOrderCancel,
						Type:       enums.ProductTypeAccessorie,
						OldValue:   p.Accessorie.Product,
						ProductId:  p.Accessorie.Product.Id,
						StoreId:    p.Accessorie.Product.StoreId,
						SourceId:   p.Accessorie.Product.Id,
						OperatorId: l.Staff.Id,
						IP:         l.Ctx.ClientIP(),
					}

					// 更新配件状态
					if err := tx.Model(&model.ProductAccessorie{}).Where("id = ?", p.Accessorie.Product.Id).Updates(&model.ProductAccessorie{
						Status: enums.ProductAccessorieStatusNormal,
					}).Update("stock", gorm.Expr("stock + ?", p.Accessorie.Quantity)).Error; err != nil {
						return errors.New("更新配件状态失败")
					}

					p.Accessorie.Product.Status = enums.ProductAccessorieStatusNormal
					p.Accessorie.Product.Stock += p.Accessorie.Quantity
					log.NewValue = p.Accessorie.Product
					if err := tx.Create(&log).Error; err != nil {
						return errors.New("创建配件历史失败")
					}
				}
			}
		}

		// 处理定金单
		for _, deposit := range order.OrderDeposits {
			// 循环产品
			for _, product := range deposit.Products {
				old_product := product
				if old_product.IsOur {
					log := model.ProductHistory{
						Action:     enums.ProductActionOrderCancel,
						Type:       enums.ProductTypeFinished,
						OldValue:   old_product.ProductFinished,
						ProductId:  old_product.ProductFinished.Id,
						StoreId:    old_product.ProductFinished.StoreId,
						SourceId:   order.Id,
						OperatorId: l.Staff.Id,
						IP:         l.Ctx.ClientIP(),
					}
					// 更新商品状态
					if err := tx.Model(&model.ProductFinished{}).Where("id = ?", old_product.ProductFinished.Id).Updates(model.ProductFinished{
						Status: enums.ProductStatusReturn,
					}).Error; err != nil {
						return errors.New("配件状态更新失败")
					}
					// 添加记录
					old_product.ProductFinished.Status = enums.ProductStatusReturn
					log.NewValue = old_product.ProductFinished
					if err := tx.Create(&log).Error; err != nil {
						return errors.New("配件记录添加失败")
					}
				}
			}

			// 更新定金单状态
			if err := tx.Model(&model.OrderDeposit{}).Where("id = ?", deposit.Id).Updates(model.OrderDeposit{
				Status: enums.OrderDepositStatusBooking,
			}).Error; err != nil {
				return errors.New("更新定金单状态失败")
			}

			// 清空与销售单的关联
			if err := tx.Model(&deposit).Association("OrderSales").Clear(); err != nil {
				return errors.New("清空定金单关联失败")
			}
		}

		// 处理支付记录
		if err := tx.Model(&model.OrderPayment{}).
			Where(&model.OrderPayment{OrderId: order.Id}).
			Select("status").
			Updates(model.OrderPayment{
				Status: false,
			}).Error; err != nil {
			return errors.New("更新支付记录失败")
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
