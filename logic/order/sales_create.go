package order

import (
	"jdy/enums"
	"jdy/errors"
	"jdy/model"
	"jdy/types"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type OrderSalesCreateLogic struct {
	Ctx   *gin.Context
	Tx    *gorm.DB
	Staff *types.Staff

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
			Remark:     req.Remark,
			MemberId:   req.MemberId,
			StoreId:    req.StoreId,
			CashierId:  req.CashierId,
			OperatorId: c.Staff.Id,
			IP:         c.Ctx.ClientIP(),

			HasIntegral: req.HasIntegral,
		},
	}

	// 开启事务
	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		l.Tx = tx

		// 创建订单
		if err := tx.Create(&l.Order).Error; err != nil {
			return errors.New("创建订单失败")
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

		// 添加支付记录
		if err := l.setPayment(); err != nil {
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
			return errors.New("成品不存在")
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
			return errors.New("配件不存在")
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

		if err := l.Tx.Model(&order).Updates(model.OrderDeposit{
			Status: enums.OrderDepositStatusComplete,
		}).Error; err != nil {
			return errors.New("定金单更新失败")
		}
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
		OrderId: l.Order.Id,
		StoreId: finished.StoreId,
		Status:  enums.OrderSalesStatusWaitPay,
		Type:    enums.ProductTypeFinished,
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
	if err := l.Tx.Model(&finished).Updates(model.ProductFinished{
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
		OrderId: l.Order.Id,
		StoreId: old.StoreId,
		Status:  enums.OrderSalesStatusWaitPay,
		Type:    enums.ProductTypeOld,
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

	// 更新商品状态
	if err := l.Tx.Model(&old).Updates(model.ProductOld{
		Status: enums.ProductStatusNormal,
	}).Error; err != nil {
		return errors.New("旧料状态更新失败")
	}

	// 添加记录
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
		OrderId: l.Order.Id,
		StoreId: accessory.StoreId,
		Status:  enums.OrderSalesStatusWaitPay,
		Type:    enums.ProductTypeAccessorie,
		Accessorie: model.OrderSalesProductAccessorie{
			OrderId:   l.Order.Id,
			StoreId:   l.Order.StoreId,
			ProductId: old_product.Id,
			Quantity:  p.Quantity,
			Price:     p.Price,
		},
	}

	if l.Req.HasIntegral {
		order_product.Accessorie.Integral = p.Integral
	}

	l.Order.Products = append(l.Order.Products, order_product)
	l.Order.ProductAccessoriePrice = l.Order.ProductAccessoriePrice.Add(order_product.Accessorie.Price)

	status := accessory.Status
	// 判断库存
	if accessory.Stock-p.Quantity == 0 {
		status = enums.ProductStatusNoStock
	}

	// 更新商品状态
	if err := l.Tx.Model(&accessory).Updates(model.ProductAccessorie{
		Stock:  accessory.Stock - p.Quantity,
		Status: status,
	}).Error; err != nil {
		return errors.New("配件状态更新失败")
	}
	// 添加记录
	log.NewValue = accessory
	if err := l.Tx.Create(&log).Error; err != nil {
		return errors.New("配件记录添加失败")
	}
	// 计算总金额
	l.Order.Price = l.Order.Price.Add(order_product.Accessorie.Price)
	l.Order.PriceOriginal = l.Order.PriceOriginal.Add(order_product.Accessorie.Price)
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
			if err := l.Tx.Model(&old_product.ProductFinished).Updates(model.ProductFinished{
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
	db = db.Preload("Store")

	if err := db.First(&product).Error; err != nil {
		return nil, errors.New("产品不存在")
	}

	// 判断商品状态
	if product.Status != enums.ProductStatusNormal {
		return nil, errors.New("产品当前不能销售")
	}

	return &product, nil
}

func (l *OrderSalesCreateLogic) getProductOld(product_id string, p *types.OrderSalesCreateReqProductOld) (*model.ProductOld, error) {
	old := &model.ProductOld{
		Code:                    p.Code,
		Name:                    p.Name,
		Status:                  enums.ProductStatusNormal,
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
		StoreId:                 l.Order.StoreId,
		RecycleMethod:           p.RecycleMethod,
		RecycleType:             p.RecycleType,
		RecyclePriceGold:        p.RecyclePriceGold,
		RecyclePriceLabor:       p.RecyclePriceLabor,
		RecyclePriceLaborMethod: p.RecyclePriceLaborMethod,
		RecyclePrice:            p.RecyclePrice,
		QualityActual:           p.QualityActual,
		RecycleSource:           enums.ProductRecycleSourceHuiShou,
		RecycleSourceId:         l.Order.Id,
		RecycleStoreId:          l.Order.StoreId,
	}

	if !p.IsOur {
		old.IsOur = false
	} else {
		if p.ProductId == "" {
			return nil, errors.New("旧料ID不能为空")
		}
		// 获取商品信息
		var finished model.ProductFinished
		db := l.Tx.Model(&model.ProductFinished{})
		db = db.Where("id = ?", product_id)
		db = db.Where(&model.ProductFinished{
			Status: enums.ProductStatusSold,
		})
		db = db.Preload("Store")

		if err := db.First(&finished).Error; err != nil {
			return nil, errors.New("旧料不存在")
		}

		old.IsOur = true
	}

	// 添加商品
	old.Class = old.GetClass()
	if err := l.Tx.Create(old).Error; err != nil {
		return nil, errors.New("旧料添加失败")
	}

	return old, nil
}

// 获取配件
func (l *OrderSalesCreateLogic) getProductAccessory(product_id string, quantity int64) (*model.ProductAccessorie, error) {
	// 获取商品信息
	var product model.ProductAccessorie
	db := l.Tx.Model(&model.ProductAccessorie{})
	db = db.Where("id = ?", product_id)
	db = db.Preload("Store")
	db = db.Preload("Category")

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
	db = db.Preload("Products", func(db *gorm.DB) *gorm.DB {
		return db.Preload("ProductFinished")
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
			StoreId:       l.Order.StoreId,
			Type:          enums.FinanceTypeIncome,
			Source:        enums.FinanceSourceSaleReceive,
			OrderType:     enums.OrderTypeSales,
			PaymentMethod: p.PaymentMethod,
			Amount:        p.Amount,
		}

		l.Order.Payments = append(l.Order.Payments, payment)
	}

	return nil
}
