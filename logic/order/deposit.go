package order

import (
	"jdy/enums"
	"jdy/errors"
	"jdy/model"
	"jdy/types"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type OrderDepositLogic struct {
	Ctx   *gin.Context
	Staff *model.Staff
}

func (l *OrderDepositLogic) Create(req *types.OrderDepositCreateReq) (*model.OrderDeposit, error) {
	// 订单信息
	order := model.OrderDeposit{
		StoreId:    req.StoreId,
		Status:     enums.OrderDepositStatusWaitPay,
		MemberId:   req.MemberId,
		CashierId:  req.CashierId,
		ClerkId:    req.ClerkId,
		Remarks:    req.Remarks,
		OperatorId: l.Staff.Id,
		IP:         l.Ctx.ClientIP(),
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 商品
		for _, p := range req.Products {
			data := model.OrderDepositProduct{
				Status:    enums.OrderDepositStatusWaitPay,
				ProductId: p.ProductId,
				PriceGold: p.PriceGold,
				Price:     p.Price,
				IsOur:     p.IsOur,
			}

			if p.IsOur {
				var product model.ProductFinished
				if err := tx.Model(&product).Where(&model.ProductFinished{
					StoreId: order.StoreId,
				}).First(&product, "id = ?", p.ProductId).Error; err != nil {
					return errors.New("获取商品信息失败")
				}
				if product.Status != enums.ProductStatusNormal {
					return errors.New("商品状态不正常")
				}

				data.ProductId = product.Id

				if err := tx.Model(&model.ProductFinished{}).Where("id = ?", product.Id).Updates(model.ProductFinished{
					Status: enums.ProductStatusReturn,
				}).Error; err != nil {
					return errors.New("更新商品状态失败")
				}
			} else {
				data.ProductDemand = model.ProductFinished{
					Name:        p.Name,
					LabelPrice:  p.LabelPrice,
					LaborFee:    p.LaborFee,
					WeightMetal: p.WeightMetal,
					RetailType:  p.RetailType,
					ColorGem:    p.ColorGem,
					Clarity:     p.ClarityGem,
				}
			}

			order.Products = append(order.Products, data)
			order.Price = order.Price.Add(data.Price)
		}

		for _, p := range req.Payments {
			order.PricePay = order.PricePay.Add(p.Amount)
			order.Payments = append(order.Payments, model.OrderPayment{
				StoreId:       order.StoreId,
				Type:          enums.FinanceTypeIncome,
				Source:        enums.FinanceSourceDepositReceive,
				PaymentMethod: p.PaymentMethod,
				Amount:        p.Amount,
				OrderType:     enums.OrderTypeDeposit,
			})
		}

		// 保存订单
		if err := tx.Create(&order).Error; err != nil {
			return errors.New("创建订单失败")
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &order, nil
}

func (l *OrderDepositLogic) List(req *types.OrderDepositListReq) (*types.PageRes[model.OrderDeposit], error) {
	var (
		order model.OrderDeposit

		res types.PageRes[model.OrderDeposit]
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

func (l *OrderDepositLogic) Info(req *types.OrderDepositInfoReq) (*model.OrderDeposit, error) {
	var (
		order model.OrderDeposit
	)

	db := model.DB.Model(&order)

	db = order.Preloads(db)

	db = db.Where("id = ?", req.Id)
	if err := db.First(&order).Error; err != nil {
		return nil, errors.New("获取订单详情失败")
	}

	return &order, nil
}

// 撤销
func (l *OrderDepositLogic) Revoked(req *types.OrderDepositRevokedReq) error {
	var (
		order model.OrderDeposit
	)

	db := model.DB.Model(&order)

	db = db.Where("id = ?", req.Id)
	db = order.Preloads(db)
	if err := db.First(&order).Error; err != nil {
		return errors.New("获取订单详情失败")
	}

	if order.Status != enums.OrderDepositStatusWaitPay {
		return errors.New("订单状态不正确")
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 撤销成品
		for _, product := range order.Products {
			// 更新订单状态
			if err := tx.Model(&model.OrderDepositProduct{}).Where("id = ?", product.Id).Updates(&model.OrderDepositProduct{
				Status: enums.OrderDepositStatusCancel,
			}).Error; err != nil {
				return errors.New("更新订单成品状态失败")
			}

			if product.IsOur && product.ProductFinished.Id != "" {
				// 更新成品状态
				if err := tx.Model(&model.ProductFinished{}).Where("id = ?", product.ProductFinished.Id).Updates(&model.ProductFinished{
					Status: enums.ProductStatusNormal,
				}).Error; err != nil {
					return errors.New("更新成品状态失败")
				}
			}
		}

		// 更新订单状态
		if err := tx.Model(&model.OrderDeposit{}).Where("id = ?", order.Id).Updates(&model.OrderDeposit{
			Status: enums.OrderDepositStatusCancel,
		}).Error; err != nil {
			return errors.New("撤销订单失败")
		}

		return nil
	}); err != nil {
		return errors.New("撤销订单失败")
	}

	return nil
}

func (l *OrderDepositLogic) Pay(req *types.OrderDepositPayReq) error {
	var (
		order model.OrderDeposit
	)

	db := model.DB.Model(&order)

	db = db.Where("id = ?", req.Id)
	db = order.Preloads(db)
	if err := db.First(&order).Error; err != nil {
		return errors.New("获取订单详情失败")
	}

	if order.Status != enums.OrderDepositStatusWaitPay {
		return errors.New("订单状态不正确")
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 支付成品
		for _, product := range order.Products {
			// 更新订单状态
			if err := tx.Model(&model.OrderDepositProduct{}).Where("id = ?", product.Id).Updates(&model.OrderDepositProduct{
				Status: enums.OrderDepositStatusBooking,
			}).Error; err != nil {
				return errors.New("更新订单成品状态失败")
			}
		}

		// 更新订单状态
		if err := tx.Model(&model.OrderDeposit{}).Where("id = ?", order.Id).Updates(&model.OrderDeposit{
			Status: enums.OrderDepositStatusBooking,
		}).Error; err != nil {
			return errors.New("支付订单失败")
		}

		return nil
	}); err != nil {
		return errors.New("支付订单失败")
	}

	return nil
}

func (l *OrderDepositLogic) Refund(req *types.OrderDepositRefundReq) error {
	var (
		order model.OrderDeposit
	)

	db := model.DB.Model(&order)

	db = db.Where("id = ?", req.Id)
	if err := db.First(&order).Error; err != nil {
		return errors.New("获取订单详情失败")
	}

	if order.Status != enums.OrderDepositStatusBooking && order.Status != enums.OrderDepositStatusRefund {
		return errors.New("订单状态不正确")
	}

	data := model.OrderRefund{
		StoreId:    order.StoreId,
		OrderId:    order.Id,
		OrderType:  enums.OrderTypeDeposit,
		MemberId:   order.MemberId,
		Remark:     req.Remark,
		OperatorId: l.Staff.Id,
		IP:         l.Ctx.ClientIP(),
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 查询成品
		var p model.OrderDepositProduct
		if err := tx.Preload("ProductFinished").First(&p, "id = ?", req.ProductId).Error; err != nil {
			return errors.New("获取成品订单详情失败")
		}

		if p.Status != enums.OrderDepositStatusBooking {
			return errors.New("成品订单状态不正确")
		}

		// 更新订单成品状态
		if err := tx.Model(&model.OrderDepositProduct{}).Where("id = ?", p.Id).Updates(&model.OrderDepositProduct{
			Status: enums.OrderDepositStatusReturn,
		}).Error; err != nil {
			return errors.New("更新订单成品状态失败")
		}

		var refund_price decimal.Decimal
		for _, payment := range req.Payments {
			refund_price = refund_price.Add(payment.Amount)

			data := model.OrderPayment{
				StoreId:       order.StoreId,
				OrderId:       order.Id,
				Type:          enums.FinanceTypeExpense,
				Source:        enums.FinanceSourceDepositRefund,
				OrderType:     enums.OrderTypeDeposit,
				PaymentMethod: payment.PaymentMethod,
				Amount:        payment.Amount,
			}
			if err := tx.Create(&data).Error; err != nil {
				return errors.New("创建退款记录失败")
			}
		}
		if refund_price.Cmp(p.Price) != 0 {
			return errors.New("退款金额不正确")
		}

		if !p.IsOur {
			data.Name = p.ProductDemand.Name
			data.Code = strings.ToUpper(p.ProductDemand.Code)
		} else {
			data.Name = p.ProductFinished.Name
			data.Code = strings.ToUpper(p.ProductFinished.Code)
			// 添加历史
			log := model.ProductHistory{
				Action:     enums.ProductActionReturn,
				Type:       enums.ProductTypeFinished,
				OldValue:   p.ProductFinished,
				ProductId:  p.ProductFinished.Id,
				StoreId:    p.ProductFinished.StoreId,
				SourceId:   p.ProductFinished.Id,
				OperatorId: l.Staff.Id,
				IP:         l.Ctx.ClientIP(),
			}
			if err := tx.Model(&model.ProductFinished{}).Where("id = ?", p.ProductFinished.Id).Updates(&model.ProductFinished{
				Status: enums.ProductStatusNormal,
			}).Error; err != nil {
				return errors.New("更新成品状态失败")
			}

			log.NewValue = p.ProductFinished
			if err := tx.Create(&log).Error; err != nil {
				return errors.New("创建成品历史失败")
			}
		}

		data.Type = enums.ProductTypeOld
		data.Quantity = 1
		data.Price = p.Price
		data.PriceOriginal = p.Price

		if err := tx.Model(&model.OrderDeposit{}).Where("id = ?", order.Id).Updates(&model.OrderDeposit{
			Status: enums.OrderDepositStatusRefund,
		}).Error; err != nil {
			return errors.New("更新订单状态失败")
		}

		if err := tx.Create(&data).Error; err != nil {
			return errors.New("创建退款记录失败")
		}

		return nil
	}); err != nil {
		return errors.New(err.Error())
	}

	return nil
}
