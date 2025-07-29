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
		// 撤销成品
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
					Status: enums.ProductStatusNormal,
				}).Error; err != nil {
					return errors.New("更新旧料状态失败")
				}
				// 添加旧料历史记录
				product.Old.Product.Status = enums.ProductStatusNormal
				log.NewValue = product.Old
				if err := tx.Create(&log).Error; err != nil {
					return errors.New("添加旧料历史记录失败")
				}
			case enums.ProductTypeAccessorie:
				log := model.ProductHistory{
					Action:     enums.ProductActionOrderCancel,
					Type:       enums.ProductTypeAccessorie,
					OldValue:   product.Accessorie.Product,
					ProductId:  product.Accessorie.Product.Id,
					StoreId:    product.Accessorie.Product.StoreId,
					SourceId:   product.Accessorie.Product.Id,
					OperatorId: l.Staff.Id,
					IP:         l.Ctx.ClientIP(),
				}
				// 更新配件状态
				stock := product.Accessorie.Product.Stock + product.Accessorie.Quantity
				if err := tx.Model(&model.ProductAccessorie{}).Where("id = ?", product.Accessorie.Product.Id).Updates(&model.ProductAccessorie{
					Stock: stock,
				}).Error; err != nil {
					return errors.New("更新配件状态失败")
				}
				// 添加配件历史记录
				log.NewValue = product.Accessorie.Product
				if err := tx.Create(&log).Error; err != nil {
					return errors.New("添加配件历史记录失败")
				}
			}
		}

		// 更新订单状态
		if err := tx.Model(&model.OrderSales{}).Where("id = ?", order.Id).Updates(&model.OrderSales{
			Status: enums.OrderSalesStatusCancel,
		}).Error; err != nil {
			return errors.New("撤销订单失败")
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
				// 更新成品状态
				if err := tx.Model(&model.ProductFinished{}).Where("id = ?", product.Finished.Product.Id).Updates(&model.ProductFinished{
					Status: enums.ProductStatusSold,
				}).Error; err != nil {
					return errors.New("更新成品状态失败")
				}
			case enums.ProductTypeOld:
				// 更新旧料状态
				if err := tx.Model(&model.ProductOld{}).Where("id = ?", product.Old.Product.Id).Updates(&model.ProductOld{
					Status: enums.ProductStatusNormal,
				}).Error; err != nil {
					return errors.New("更新旧料状态失败")
				}
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

func (l *OrderSalesLogic) Refund(req *types.OrderSalesRefundReq) error {
	var (
		order model.OrderSales
	)

	db := model.DB.Model(&order)

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
			data.Code = strings.ToUpper(p.Finished.Product.Code)
			data.Name = p.Finished.Product.Name
			data.Quantity = 1
			data.Price = req.Price
			data.PriceOriginal = p.Finished.Price

		case enums.ProductTypeOld:
			// 查询旧料
			var p model.OrderSalesProduct
			if err := tx.Preload("Old.Product").First(&p, "id = ?", req.ProductId).Error; err != nil {
				return errors.New("获取旧料订单详情失败")
			}
			if p.Status != enums.OrderSalesStatusComplete {
				return errors.New("旧料订单状态不正确")
			}

			// 更新订单旧料状态
			if err := tx.Model(&model.OrderSalesProduct{}).Where("id = ?", p.Id).Updates(&model.OrderSalesProduct{
				Status: enums.OrderSalesStatusReturn}).Error; err != nil {
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
				Status: enums.ProductStatusNormal,
			}).Error; err != nil {
				return errors.New("更新旧料状态失败")
			}

			log.NewValue = p.Old.Product
			if err := tx.Create(&log).Error; err != nil {
				return errors.New("创建旧料历史失败")
			}

			data.Type = enums.ProductTypeOld
			data.Code = strings.ToUpper(p.Old.Product.Code)
			data.Name = p.Old.Product.Name
			data.Quantity = 1
			data.Price = req.Price
			data.PriceOriginal = p.Old.RecyclePrice

		case enums.ProductTypeAccessorie:
			// 查询配件
			var p model.OrderSalesProduct
			if err := tx.Preload("Accessorie.Product.Category").First(&p, "id = ?", req.ProductId).Error; err != nil {
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
				ProductId:  p.Accessorie.Product.Id,
				StoreId:    p.Accessorie.Product.StoreId,
				SourceId:   p.Accessorie.Product.Id,
				OperatorId: l.Staff.Id,
				IP:         l.Ctx.ClientIP(),
			}

			// 更新配件状态
			if err := tx.Model(&model.ProductAccessorie{}).Where("id = ?", p.Accessorie.Product.Id).Updates(&model.ProductAccessorie{
				Status: enums.ProductStatusNormal,
				Stock:  p.Accessorie.Product.Stock + p.Accessorie.Quantity,
			}).Error; err != nil {
				return errors.New("更新配件状态失败")
			}
			log.NewValue = p.Accessorie.Product
			if err := tx.Create(&log).Error; err != nil {
				return errors.New("创建配件历史失败")
			}

			data.Type = enums.ProductTypeAccessorie
			data.Code = strings.ToUpper(p.Accessorie.Product.Code)
			data.Name = p.Accessorie.Product.Category.Name
			data.Quantity = p.Accessorie.Quantity
			data.Price = req.Price
			data.PriceOriginal = p.Accessorie.Price
		}

		if err := tx.Model(&model.OrderSales{}).Where("id = ?", order.Id).Updates(&model.OrderSales{
			Status: enums.OrderSalesStatusRefund,
		}).Error; err != nil {
			return errors.New("更新订单状态失败")
		}

		if err := tx.Create(&data).Error; err != nil {
			return errors.New("创建退货单失败")
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
