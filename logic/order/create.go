package order

import (
	"jdy/enums"
	"jdy/model"
	"jdy/types"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderCreateLogic struct {
	Ctx   *gin.Context
	Tx    *gorm.DB
	Staff *types.Staff

	Req *types.OrderCreateReq

	Order *model.Order
}

// 创建订单
func (c *OrderLogic) Create(req *types.OrderCreateReq) (*model.Order, error) {
	l := OrderCreateLogic{
		Ctx:   c.Ctx,
		Req:   req,
		Staff: c.Staff,
		Order: &model.Order{
			Type:      req.Type,
			Status:    enums.OrderStatusWaitPay,
			Source:    req.Source,
			Remark:    req.Remark,
			MemberId:  req.MemberId,
			StoreId:   req.StoreId,
			CashierId: req.CashierId,

			OperatorId: c.Staff.Id,
			IP:         c.Ctx.ClientIP(),
		},
	}

	// 开启事务
	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		l.Tx = tx

		// 创建订单
		if err := tx.Create(&l.Order).Error; err != nil {
			return err
		}

		// // 计算金额
		// if err := l.loopSales(); err != nil {
		// 	return err
		// }

		// // 添加优惠
		// if err := l.getDiscount(); err != nil {
		// 	return err
		// }

		// // 计算业绩
		// if err := l.getPerformance(); err != nil {
		// 	return err
		// }

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

// // 销售单金额
// func (l *OrderCreateLogic) loopSales() error {
// 	for _, p := range l.Req.Products {
// 		// 获取商品
// 		product, err := l.getProduct(p.ProductId)
// 		if err != nil {
// 			return err
// 		}

// 		old_product := *product
// 		log := model.ProductHistory{
// 			Action:     enums.ProductActionOrder,
// 			OldValue:   old_product,
// 			ProductId:  old_product.Id,
// 			StoreId:    old_product.StoreId,
// 			SourceId:   l.Order.Id,
// 			OperatorId: l.Staff.Id,
// 			IP:         l.Ctx.ClientIP(),
// 		}

// 		// 获取金价
// 		gold_price, err := model.GetGoldPrice(&types.GoldPriceOptions{
// 			StoreId:         l.Order.StoreId,
// 			ProductMaterial: product.Material,
// 			ProductType:     product.Type,
// 			ProductBrand:    []enums.ProductBrand{product.Brand},
// 			ProductQuality:  []enums.ProductQuality{product.Quality},
// 		})
// 		if err != nil {
// 			return err
// 		}

// 		var (
// 			price           decimal.Decimal                                                                         // 单价
// 			amount          decimal.Decimal                                                                         // 原价
// 			discount        decimal.Decimal = decimal.NewFromFloat(1).Sub(p.Discount.Div(decimal.NewFromFloat(10))) // 折扣
// 			amount_discount decimal.Decimal                                                                         // 折扣价
// 		)

// 		switch product.RetailType {
// 		case enums.ProductRetailTypePiece: // 一口价 = 标签价×数量
// 			{
// 				// 单价
// 				price = product.LabelPrice
// 				// 原价
// 				amount = price.Mul(decimal.NewFromInt(p.Quantity))
// 			}

// 		case enums.ProductRetailTypeGoldKg: // 计重论克 = (金价+工费)×克重×数量
// 			{
// 				// 单价 = 今日金价
// 				price = gold_price.Add(product.LaborFee)
// 				// 原价 = (金价+工费)×克重×数量
// 				amount = product.WeightMetal.Mul(price).Mul(decimal.NewFromInt(p.Quantity))
// 			}

// 		case enums.ProductRetailTypeGoldPiece: // 计重工费论件 = 金价×克重+工费
// 			{
// 				// 单价
// 				price = gold_price
// 				// 原价 = 金价×克重+工费×数量
// 				amount = gold_price.Mul(product.WeightMetal).Add(product.LaborFee).Mul(decimal.NewFromInt(p.Quantity))
// 			}
// 		default:
// 			{
// 				return errors.New("产品类型错误")
// 			}
// 		}

// 		// 折扣价
// 		amount_discount = amount.Mul(discount)

// 		// 添加订单商品
// 		order_product := model.OrderProduct{
// 			ProductId: product.Id,
// 			Status:    enums.OrderStatusWaitPay,

// 			Quantity:       p.Quantity,
// 			Price:          price,
// 			Amount:         amount.Sub(amount_discount),
// 			AmountOriginal: amount,

// 			Discount:       p.Discount,
// 			DiscountAmount: amount_discount,
// 		}
// 		l.Order.Products = append(l.Order.Products, order_product)

// 		// 更新商品状态
// 		product.Status = enums.ProductStatusSold
// 		if err := l.Tx.Save(&product).Error; err != nil {
// 			return err
// 		}

// 		// 添加记录
// 		log.NewValue = product
// 		if err := l.Tx.Create(&log).Error; err != nil {
// 			return err
// 		}

// 		// 计算总金额
// 		l.Order.Amount = l.Order.Amount.Add(order_product.Amount)
// 		l.Order.AmountOriginal = l.Order.AmountOriginal.Add(order_product.AmountOriginal)
// 	}

// 	return nil
// }

// // 获取商品
// func (l *OrderCreateLogic) getProduct(product_id string) (*model.Product, error) {
// 	// 获取商品信息
// 	var product model.Product
// 	db := l.Tx.Model(&model.Product{})
// 	db = db.Where("id = ?", product_id)
// 	db = db.Preload("Store")
// 	db = db.Preload("RecycleStore")

// 	if err := db.First(&product).Error; err != nil {
// 		return nil, errors.New("产品不存在")
// 	}

// 	// 判断商品状态
// 	if product.Status != enums.ProductStatusNormal {
// 		return nil, errors.New("产品当前不能销售")
// 	}

// 	return &product, nil
// }

// // 计算整单优惠
// func (l *OrderCreateLogic) getDiscount() error {
// 	// 判断整单折扣
// 	// 折扣
// 	l.Order.DiscountRate = decimal.NewFromFloat(1).Sub(l.Req.DiscountRate.Div(decimal.NewFromFloat(10)))
// 	// 折扣金额
// 	l.Order.DiscountAmount = l.Order.Amount.Mul(l.Order.DiscountRate)
// 	// 折扣后金额
// 	l.Order.Amount = l.Order.Amount.Sub(l.Order.DiscountAmount)

// 	// 抹零
// 	l.Order.AmountReduce = l.Req.AmountReduce
// 	l.Order.Amount = l.Order.Amount.Sub(l.Req.AmountReduce)

// 	return nil
// }

// // 计算业绩
// func (l *OrderCreateLogic) getPerformance() error {
// 	// 添加导购员业绩
// 	for _, s := range l.Req.Salesmans {
// 		var salesman model.Staff
// 		db := l.Tx.Model(&model.Staff{})
// 		db = db.Where("id = ?", s.SalesmanId)
// 		db = db.Where(&model.Staff{IsDisabled: false})
// 		if err := db.First(&salesman).Error; err != nil {
// 			return err
// 		}

// 		// 计算业绩 佣金 = 佣金率/100 * 订单金额
// 		performance := l.Order.Amount.Mul(s.PerformanceRate).Div(decimal.NewFromFloat(100))

// 		// 添加导购员业绩
// 		l.Order.Salesmans = append(l.Order.Salesmans, model.OrderSalesman{
// 			SalesmanId:        salesman.Id,
// 			PerformanceRate:   s.PerformanceRate,
// 			PerformanceAmount: performance,
// 			IsMain:            s.IsMain,
// 		})
// 	}

// 	return nil
// }
