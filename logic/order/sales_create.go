package order

import (
	"jdy/enums"
	"jdy/errors"
	"jdy/message"
	"jdy/model"
	"jdy/types"
	"jdy/utils"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type OrderSalesCreateLogic struct {
	Ctx   *gin.Context
	Tx    *gorm.DB
	Staff *model.Staff
	Store *model.Store

	Req *types.OrderSalesCreateReq

	Order *model.OrderSales
}

// 创建订单
func (c *OrderSalesLogic) Create(req *types.OrderSalesCreateReq) (*model.OrderSales, error) {
	l := OrderSalesCreateLogic{
		Ctx:   c.Ctx,
		Req:   req,
		Staff: c.Staff,
		Order: &model.OrderSales{
			Status:     enums.OrderSalesStatusWaitPay,
			Source:     req.Source,
			Remarks:    req.Remarks,
			MemberId:   req.MemberId,
			StoreId:    req.StoreId,
			CashierId:  req.CashierId,
			OperatorId: c.Staff.Id,
			Operator:   *c.Staff,
			IP:         c.Ctx.ClientIP(),

			HasIntegral: req.HasIntegral,
		},
	}

	// 开启事务
	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		l.Tx = tx

		// 获取门店
		if err := l.getStore(); err != nil {
			return err
		}

		// 创建订单
		if err := tx.Create(&l.Order).Error; err != nil {
			return errors.New("创建订单失败")
		}

		// 添加支付记录
		if err := l.setPayment(); err != nil {
			return err
		}

		// 计算金额
		if err := l.loopSales(); err != nil {
			return err
		}

		// 添加优惠
		if err := l.getDiscount(); err != nil {
			return err
		}

		// 计算业绩
		if err := l.getPerformance(); err != nil {
			return err
		}

		// 更新订单
		if err := tx.Save(&l.Order).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	// 发送通知
	go func() {
		msg := message.NewMessage(c.Ctx)
		if err := msg.SendOrderSalesCreateMessage(&message.OrderSalesMessage{
			OrderSales: l.Order,
		}); err != nil {
			log.Printf("发送订单创建通知失败: %v", err)
		}
	}()

	return l.Order, nil
}

// 销售单金额
func (l *OrderSalesCreateLogic) loopSales() error {
	for _, p := range l.Req.ProductFinisheds {
		if p.ProductId == "" {
			return errors.New("成品ID不能为空")
		}

		// 获取商品
		finished, err := l.getProductFinished(p.ProductId)
		if err != nil {
			return err
		}

		if err := l.loopFinished(&p, finished); err != nil {
			return err
		}
	}

	for _, p := range l.Req.ProductOlds {
		old, err := l.getProductOld(p.ProductId, &p)
		if err != nil {
			return err
		}

		if err := l.loopOld(&p, old); err != nil {
			return err
		}
	}

	for _, p := range l.Req.ProductAccessories {
		if p.ProductId == "" {
			return errors.New("配件ID不能为空")
		}

		// 获取配件
		accessory, err := l.getProductAccessory(p.ProductId, p.Quantity)
		if err != nil {
			return err
		}

		if err := l.loopAccessory(&p, accessory); err != nil {
			return err
		}
	}

	for _, p := range l.Req.OrderDepositIds {
		// 获取定金订单
		order, err := l.getOrderDeposit(p)
		if err != nil {
			return err
		}

		if err := l.loopOrderDeposit(order); err != nil {
			return err
		}

		if err := l.Tx.Model(&model.OrderDeposit{}).Where("id = ?", order.Id).Updates(model.OrderDeposit{
			Status: enums.OrderDepositStatusComplete,
		}).Error; err != nil {
			return errors.New("定金单更新失败")
		}
	}

	if l.Order.PricePay.Cmp(l.Order.Price) != 0 {
		return errors.New("支付方式与应付金额不一致")
	}

	// 计算优惠金额
	l.Order.PriceDiscount = l.Order.PriceOriginal.Sub(l.Order.Price)

	return nil
}

func (l *OrderSalesCreateLogic) loopFinished(p *types.OrderSalesCreateReqProductFinished, finished *model.ProductFinished) error {
	old_product := *finished
	log := model.ProductHistory{
		Action:     enums.ProductActionOrder,
		Type:       enums.ProductTypeFinished,
		OldValue:   old_product,
		ProductId:  old_product.Id,
		StoreId:    old_product.StoreId,
		SourceId:   l.Order.Id,
		OperatorId: l.Staff.Id,
		IP:         l.Ctx.ClientIP(),
	}

	// 添加订单商品
	order_product := model.OrderSalesProduct{
		OrderId:  l.Order.Id,
		StoreId:  finished.StoreId,
		Status:   enums.OrderSalesStatusWaitPay,
		Type:     enums.ProductTypeFinished,
		Code:     strings.ToUpper(finished.Code),
		MemberId: l.Order.MemberId,
		Finished: model.OrderSalesProductFinished{
			OrderId:           l.Order.Id,
			StoreId:           finished.StoreId,
			ProductId:         finished.Id,
			PriceGold:         p.PriceGold,
			LaborFee:          p.LaborFee,
			DiscountFixed:     p.DiscountFixed,
			IntegralDeduction: p.IntegralDeduction,
			DiscountMember:    p.DiscountMember,
			RoundOff:          p.RoundOff,
			PriceOriginal:     p.PriceOriginal,
			Price:             p.Price,
			DiscountFinal:     p.DiscountFinal,
		},
	}

	if l.Req.HasIntegral {
		order_product.Finished.Integral = p.Integral
	}

	l.Order.Products = append(l.Order.Products, order_product)
	l.Order.ProductFinishedPrice = l.Order.ProductFinishedPrice.Add(order_product.Finished.Price)

	// 更新商品状态
	if err := l.Tx.Model(&model.ProductFinished{}).Where("id = ?", finished.Id).Updates(model.ProductFinished{
		Status: enums.ProductStatusSold,
	}).Error; err != nil {
		return errors.New("成品状态更新失败")
	}

	// 添加记录
	log.NewValue = finished
	if err := l.Tx.Create(&log).Error; err != nil {
		return errors.New("成品记录添加失败")
	}

	// 计算总金额
	l.Order.Price = l.Order.Price.Add(order_product.Finished.Price)
	l.Order.PriceOriginal = l.Order.PriceOriginal.Add(order_product.Finished.PriceOriginal)
	l.Order.Integral = l.Order.Integral.Add(order_product.Finished.Integral)

	return nil
}

func (l *OrderSalesCreateLogic) loopOld(p *types.OrderSalesCreateReqProductOld, old *model.ProductOld) error {
	old_product := *old
	log := model.ProductHistory{
		Action:     enums.ProductActionOrder,
		Type:       enums.ProductTypeOld,
		OldValue:   old_product,
		ProductId:  old_product.Id,
		StoreId:    old_product.StoreId,
		SourceId:   l.Order.Id,
		OperatorId: l.Staff.Id,
		IP:         l.Ctx.ClientIP(),
	}

	// 添加订单商品
	order_product := model.OrderSalesProduct{
		OrderId:  l.Order.Id,
		StoreId:  old.StoreId,
		Status:   enums.OrderSalesStatusWaitPay,
		Type:     enums.ProductTypeOld,
		Code:     strings.ToUpper(old.Code),
		MemberId: l.Order.MemberId,
		Old: model.OrderSalesProductOld{
			OrderId:                 l.Order.Id,
			StoreId:                 old.StoreId,
			ProductId:               old.Id,
			WeightMetal:             p.WeightMetal,
			RecyclePriceGold:        p.RecyclePriceGold,
			RecyclePriceLabor:       p.RecyclePriceLabor,
			RecyclePriceLaborMethod: p.RecyclePriceLaborMethod,
			QualityActual:           p.QualityActual,
			RecyclePrice:            p.RecyclePrice,
			Integral:                p.Integral,
		},
	}

	if l.Req.HasIntegral {
		order_product.Old.Integral = p.Integral
	}

	l.Order.Products = append(l.Order.Products, order_product)
	l.Order.ProductOldPrice = l.Order.ProductOldPrice.Add(order_product.Old.RecyclePrice)

	// 添加记录
	log.NewValue = old
	if err := l.Tx.Create(&log).Error; err != nil {
		return errors.New("旧料记录添加失败")
	}
	// 计算总金额(减少)
	l.Order.Price = l.Order.Price.Sub(order_product.Old.RecyclePrice)
	l.Order.PriceOriginal = l.Order.PriceOriginal.Sub(order_product.Old.RecyclePrice)
	l.Order.Integral = l.Order.Integral.Sub(order_product.Old.Integral)

	return nil
}

func (l *OrderSalesCreateLogic) loopAccessory(p *types.OrderSalesCreateReqProductAccessorie, accessory *model.ProductAccessorie) error {
	old_product := *accessory
	log := model.ProductHistory{
		Action:     enums.ProductActionOrder,
		Type:       enums.ProductTypeAccessorie,
		OldValue:   old_product,
		ProductId:  old_product.Id,
		StoreId:    old_product.StoreId,
		SourceId:   l.Order.Id,
		OperatorId: l.Staff.Id,
		IP:         l.Ctx.ClientIP(),
	}

	// 添加订单商品
	order_product := model.OrderSalesProduct{
		OrderId:  l.Order.Id,
		StoreId:  l.Req.StoreId,
		Status:   enums.OrderSalesStatusWaitPay,
		Type:     enums.ProductTypeAccessorie,
		MemberId: l.Order.MemberId,
		Accessorie: model.OrderSalesProductAccessorie{
			OrderId:   l.Order.Id,
			StoreId:   l.Req.StoreId,
			ProductId: old_product.Id,
			Quantity:  p.Quantity,
			Price:     p.Price,
		},
	}

	if l.Req.HasIntegral {
		order_product.Accessorie.Integral = p.Integral
	}

	l.Order.Products = append(l.Order.Products, order_product)
	l.Order.ProductAccessoriePrice = l.Order.ProductAccessoriePrice.Add(order_product.Accessorie.Price.Mul(decimal.NewFromInt(order_product.Accessorie.Quantity)))

	// 判断库存
	stock := accessory.Stock - p.Quantity
	status := accessory.Status
	if stock == 0 {
		status = enums.ProductAccessorieStatusNoStock
	}

	// 更新商品状态
	if err := l.Tx.Model(&model.ProductAccessorie{}).Where("id = ?", accessory.Id).Updates(model.ProductAccessorie{
		Status: status,
	}).Update("stock", stock).Error; err != nil {
		return errors.New("配件状态更新失败")
	}
	// 添加记录
	log.NewValue = accessory
	if err := l.Tx.Create(&log).Error; err != nil {
		return errors.New("配件记录添加失败")
	}
	// 计算总金额
	l.Order.Price = l.Order.Price.Add(order_product.Accessorie.Price.Mul(decimal.NewFromInt(order_product.Accessorie.Quantity)))
	l.Order.PriceOriginal = l.Order.PriceOriginal.Add(order_product.Accessorie.Price.Mul(decimal.NewFromInt(order_product.Accessorie.Quantity)))
	l.Order.Integral = l.Order.Integral.Add(order_product.Accessorie.Integral)

	return nil
}

func (l *OrderSalesCreateLogic) loopOrderDeposit(order *model.OrderDeposit) error {
	// 循环产品
	for _, product := range order.Products {
		old_product := product
		if old_product.IsOur {
			log := model.ProductHistory{
				Action:     enums.ProductActionOrder,
				Type:       enums.ProductTypeFinished,
				OldValue:   old_product.ProductFinished,
				ProductId:  old_product.ProductFinished.Id,
				StoreId:    old_product.ProductFinished.StoreId,
				SourceId:   l.Order.Id,
				OperatorId: l.Staff.Id,
				IP:         l.Ctx.ClientIP(),
			}
			// 更新商品状态
			if err := l.Tx.Model(&model.ProductFinished{}).Where("id = ?", old_product.ProductFinished.Id).Updates(model.ProductFinished{
				Status: enums.ProductStatusSold,
			}).Error; err != nil {
				return errors.New("配件状态更新失败")
			}
			// 添加记录
			log.NewValue = old_product.ProductFinished
			if err := l.Tx.Create(&log).Error; err != nil {
				return errors.New("配件记录添加失败")
			}
		}

		// 添加关联
		l.Order.OrderDeposits = append(l.Order.OrderDeposits, *order)

		// 计算总金额
		l.Order.Price = l.Order.Price.Sub(product.Price)
		l.Order.PriceDeposit = l.Order.PriceDeposit.Add(product.Price)
	}

	return nil
}

// 获取商品
func (l *OrderSalesCreateLogic) getProductFinished(product_id string) (*model.ProductFinished, error) {
	// 获取商品信息
	var product model.ProductFinished
	db := l.Tx.Model(&model.ProductFinished{})
	db = db.Where("id = ?", product_id)
	db = db.Where(&model.ProductFinished{
		StoreId: l.Req.StoreId,
	})
	db = db.Preload("Store")

	if err := db.First(&product).Error; err != nil {
		return nil, errors.New("产品不存在")
	}

	// 判断商品状态
	if product.Status != enums.ProductStatusNormal && product.Status != enums.ProductStatusReturn {
		return nil, errors.New("产品当前不能销售")
	}

	return &product, nil
}

func (l *OrderSalesCreateLogic) getProductOld(product_id string, p *types.OrderSalesCreateReqProductOld) (*model.ProductOld, error) {
	old := model.ProductOld{
		Code:                    strings.ToUpper("JL" + utils.RandomCode(8)),
		CodeFinished:            strings.ToUpper(p.Code),
		Name:                    p.Name,
		LabelPrice:              p.LabelPrice,
		Brand:                   p.Brand,
		Material:                p.Material,
		Quality:                 p.Quality,
		Gem:                     p.Gem,
		Category:                p.Category,
		Craft:                   p.Craft,
		WeightMetal:             p.WeightMetal,
		WeightTotal:             p.WeightTotal,
		ColorGem:                p.ColorGem,
		WeightGem:               p.WeightGem,
		NumGem:                  p.NumGem,
		Clarity:                 p.ClarityGem,
		Cut:                     p.Cut,
		WeightOther:             p.WeightOther,
		NumOther:                p.NumOther,
		Remark:                  p.Remark,
		StoreId:                 l.Req.StoreId,
		Store:                   *l.Store,
		RecycleMethod:           p.RecycleMethod,
		RecycleType:             p.RecycleType,
		RecyclePriceGold:        p.RecyclePriceGold,
		RecyclePriceLabor:       p.RecyclePriceLabor,
		RecyclePriceLaborMethod: p.RecyclePriceLaborMethod,
		RecyclePrice:            p.RecyclePrice,
		QualityActual:           p.QualityActual,
		RecycleSource:           enums.ProductRecycleSourceHuiShou,
		RecycleSourceId:         l.Order.Id,
		RecycleStoreId:          l.Req.StoreId,
		RecycleStore:            *l.Store,
	}

	if p.IsOur {
		if product_id != "" {
			var finished model.ProductFinished
			// 获取商品信息
			db := l.Tx.Model(&model.ProductFinished{})
			db = db.Where("id = ?", product_id)
			db = db.Or("code = ?", strings.ToUpper(p.Code))
			db = db.Where(&model.ProductFinished{
				Status: enums.ProductStatusSold,
			})
			db = db.Preload("Store")

			if err := db.First(&finished).Error; err != nil {
				return nil, errors.New("旧料不存在")
			}
		}
		old.IsOur = true
		old.Status = enums.ProductStatusDraft
	} else {
		old.IsOur = false
		old.Status = enums.ProductStatusDraft
	}

	// 添加商品
	old.Class = old.GetClass()
	if err := l.Tx.Create(&old).Error; err != nil {
		return nil, errors.New("旧料添加失败")
	}

	return &old, nil
}

// 获取配件
func (l *OrderSalesCreateLogic) getProductAccessory(product_id string, quantity int64) (*model.ProductAccessorie, error) {
	// 获取商品信息
	var product model.ProductAccessorie
	db := l.Tx.Model(&model.ProductAccessorie{})
	db = db.Where("id = ?", product_id)
	db = db.Where(&model.ProductAccessorie{
		StoreId: l.Req.StoreId,
	})
	db = db.Preload("Store")

	if err := db.First(&product).Error; err != nil {
		return nil, errors.New("产品不存在")
	}

	// 判断商品状态
	if product.Stock < quantity {
		return nil, errors.New("配件库存不足")
	}

	return &product, nil
}

// 获取定金订单
func (l *OrderSalesCreateLogic) getOrderDeposit(order_id string) (*model.OrderDeposit, error) {
	// 获取商品信息
	var order model.OrderDeposit
	db := l.Tx.Model(&model.OrderDeposit{})
	db = db.Where("id = ?", order_id)
	db = db.Preload("Products", func(tx *gorm.DB) *gorm.DB {
		return tx.Preload("ProductFinished")
	})

	if err := db.First(&order).Error; err != nil {
		return nil, errors.New("定金订单不存在")
	}

	if order.Status != enums.OrderDepositStatusBooking {
		return nil, errors.New("定金单状态不正确")
	}

	return &order, nil
}

// 计算整单优惠
func (l *OrderSalesCreateLogic) getDiscount() error {
	// 整单折扣
	l.Order.DiscountRate = decimal.NewFromFloat(1).Sub(l.Req.DiscountRate.Div(decimal.NewFromFloat(10)))
	// 积分抵扣
	l.Order.IntegralDeduction = l.Req.IntegralDeduction
	// 抹零
	l.Order.RoundOff = l.Req.RoundOff
	// 优惠金额
	l.Order.PriceDiscount = l.Order.PriceOriginal.Sub(l.Order.Price)

	return nil
}

// 计算业绩
func (l *OrderSalesCreateLogic) getPerformance() error {
	// 添加导购员业绩
	for _, s := range l.Req.Clerks {
		var salesman model.Staff
		db := l.Tx.Model(&model.Staff{})
		db = db.Where("id = ?", s.SalesmanId)
		db = db.Where(&model.Staff{IsDisabled: false})
		if err := db.First(&salesman).Error; err != nil {
			return err
		}

		// 计算业绩 佣金 = 佣金率/100 * 订单金额
		performance := l.Order.Price.Mul(s.PerformanceRate).Div(decimal.NewFromFloat(100))

		// 添加导购员业绩
		l.Order.Clerks = append(l.Order.Clerks, model.OrderSalesClerk{
			SalesmanId:        salesman.Id,
			PerformanceRate:   s.PerformanceRate,
			PerformanceAmount: performance,
			IsMain:            s.IsMain,
		})
	}

	return nil
}

// 添加支付记录
func (l *OrderSalesCreateLogic) setPayment() error {
	// 添加支付记录
	for _, p := range l.Req.Payments {
		payment := model.OrderPayment{
			StoreId:       l.Req.StoreId,
			Type:          enums.FinanceTypeIncome,
			Source:        enums.FinanceSourceSaleReceive,
			OrderType:     enums.OrderTypeSales,
			PaymentMethod: p.PaymentMethod,
			Amount:        p.Amount,
		}

		l.Order.Payments = append(l.Order.Payments, payment)
		l.Order.PricePay = l.Order.PricePay.Add(p.Amount)
	}

	return nil
}

// 查询门店
func (l *OrderSalesCreateLogic) getStore() error {
	var store model.Store
	db := l.Tx.Model(&model.Store{})
	if err := db.First(&store, "id = ?", l.Req.StoreId).Error; err != nil {
		return errors.New("门店不存在")
	}

	l.Store = &store
	l.Order.Store = store

	return nil
}
