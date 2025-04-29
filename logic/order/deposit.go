package order

import (
	"jdy/enums"
	"jdy/errors"
	"jdy/model"
	"jdy/types"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderDepositLogic struct {
	Ctx   *gin.Context
	Staff *types.Staff
}

func (l *OrderDepositLogic) Create(req *types.OrderDepositCreateReq) (*model.OrderDeposit, error) {
	// 订单信息
	order := model.OrderDeposit{
		StoreId:    req.StoreId,
		Status:     enums.OrderDepositStatusWaitPay,
		MemberId:   req.MemberId,
		CashierId:  req.CashierId,
		ClerkId:    req.ClerkId,
		Remark:     req.Remark,
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

				if err := tx.Model(&product).Where("id = ?", product.Id).Updates(model.ProductFinished{
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
		}

		for _, p := range req.Payments {
			order.Price = order.Price.Add(p.Amount)
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
	db = db.Preload("Member")
	db = db.Preload("Store")
	db = db.Preload("Cashier")
	db = db.Preload("Clerk")
	db = db.Preload("Products", func(db *gorm.DB) *gorm.DB {
		return db.Preload("ProductFinished")
	})
	db = db.Preload("OrderSales")
	db = db.Preload("Payments")

	db = db.Order("created_at desc")
	db = model.PageCondition(db, req.Page, req.Limit)
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

	db = db.Preload("Member")
	db = db.Preload("Store")
	db = db.Preload("Cashier")
	db = db.Preload("Clerk")
	db = db.Preload("Products", func(db *gorm.DB) *gorm.DB {
		return db.Preload("ProductFinished")
	})
	db = db.Preload("OrderSales")
	db = db.Preload("Payments")

	db = db.Where("id = ?", req.Id)
	if err := db.First(&order).Error; err != nil {
		return nil, errors.New("获取订单详情失败")
	}

	return &order, nil
}
