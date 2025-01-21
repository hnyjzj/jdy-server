package order

import (
	"jdy/enums"
	"jdy/errors"
	"jdy/model"
	"jdy/types"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderCreateLogic struct {
	Ctx *gin.Context
	Tx  *gorm.DB

	Req *types.OrderCreateReq

	Order *model.Order

	GoldPrice float64
}

// 创建订单
func (c *OrderLogic) Create(req *types.OrderCreateReq) error {
	// 开启事务
	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		l := OrderCreateLogic{
			Ctx: c.Ctx,
			Tx:  tx,
			Req: req,
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

		// 获取今日金价
		if err := l.getGoldPrice(); err != nil {
			return err
		}

		// 计算金额
		if err := l.getAmount(); err != nil {
			return err
		}

		// 添加优惠
		if err := l.getDiscount(); err != nil {
			return err
		}

		// 创建订单
		if err := tx.Create(&l.Order).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}
	return nil
}

// 获取今日金价
func (l *OrderCreateLogic) getGoldPrice() error {
	gold_price, err := model.GetGoldPrice()
	if err != nil {
		return err
	}

	l.GoldPrice = gold_price

	return nil
}

// 计算金额
func (l *OrderCreateLogic) getAmount() error {
	switch l.Order.Type {
	case enums.OrderTypeSales:
		{
			if err := l.loopSales(); err != nil {
				return err
			}
		}
	default:
		{
			return errors.New("订单类型错误")
		}
	}
	return nil
}

// 销售单金额
func (l *OrderCreateLogic) loopSales() error {
	for _, p := range l.Req.Products {
		// 获取商品
		product, err := l.getProduct(p.ProductId)
		if err != nil {
			return err
		}

		// 判断单品折扣
		var discount_rate float64
		if p.Discount == nil {
			discount_rate = 10
		}

		var (
			price           float64                      // 单价
			amount          float64                      // 原价
			discount        float64 = discount_rate / 10 // 折扣
			amount_discount float64                      // 折扣价
		)

		switch product.RetailType {
		case enums.ProductRetailTypePiece: // 一口价 = 单价x数量
			{
				// 单价
				price = product.Price
				// 原价
				amount = product.Price * float64(p.Quantity)
				// 折扣价
				amount_discount = amount * discount
			}

		case enums.ProductRetailTypeGoldKg: // 计重论克 = (金价+工费)×克重×数量
			{
				// 单价
				price = l.GoldPrice
				// 原价
				amount = (l.GoldPrice + product.LaborFee) * product.WeightMetal * float64(p.Quantity)
				amount_discount = amount * discount
			}

		case enums.ProductRetailTypeGoldPiece: // 计重工费论件 = 金价×克重+工费
			{
				// 单价
				price = l.GoldPrice
				// 原价
				amount = (l.GoldPrice*product.WeightMetal + product.LaborFee) * float64(p.Quantity)
				amount_discount = amount * discount
			}
		default:
			{
				return errors.New("产品类型错误")
			}
		}

		// 添加记录
		l.Order.Products = append(l.Order.Products, model.OrderProduct{
			ProductId: product.Id,

			Quantity:       p.Quantity,
			Price:          price,
			Amount:         amount - amount_discount,
			AmountOriginal: amount,

			Discount:       discount,
			DiscountAmount: amount_discount,
		})

		// 计算总金额
		l.Order.Amount += amount - amount_discount
		l.Order.AmountOriginal += amount
	}

	return nil
}

// 获取商品
func (l *OrderCreateLogic) getProduct(product_id string) (*model.Product, error) {
	// 获取商品信息
	var product model.Product
	db := l.Tx.Model(&model.Product{})
	db = db.Where("id = ?", product_id)
	if err := db.First(&product).Error; err != nil {
		return nil, errors.New("产品不存在")
	}

	// 判断商品状态
	if product.Status != enums.ProductStatusNormal {
		return nil, errors.New("产品当前不能销售")
	}

	return &product, nil
}

// 计算整单优惠
func (l *OrderCreateLogic) getDiscount() error {
	// 判断整单折扣
	var discount_rate float64
	if l.Req.DiscountRate == nil {
		discount_rate = 10
	}
	// 折扣
	l.Order.DiscountRate = discount_rate / 10
	// 折扣金额
	l.Order.DiscountAmount = l.Order.Amount * l.Order.DiscountRate
	// 折扣后金额
	l.Order.Amount -= l.Order.DiscountAmount

	// 抹零
	l.Order.AmountReduce = l.Req.AmountReduce
	l.Order.Amount -= l.Order.AmountReduce

	return nil
}
