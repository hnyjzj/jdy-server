package target

import (
	"errors"
	"jdy/enums"
	"jdy/model"
	"jdy/utils"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type achieveLogic struct {
	db *gorm.DB

	Target     *model.Target
	start_time string
	end_time   string

	Sales   []model.OrderSales
	Refunds []model.OrderRefund
}

// 获取目标完成情况
func (l *Logic) GetAchieve(target_id string) (*model.Target, error) {
	logic := &achieveLogic{
		db: model.DB,
	}

	// 获取目标
	if err := logic.getTarget(target_id); err != nil {
		return nil, err
	}

	// 获取销售单
	if err := logic.getSales(); err != nil {
		return nil, err
	}

	// 获取退款单
	if err := logic.getRefunds(); err != nil {
		return nil, err
	}

	// 计算业绩
	if err := logic.Calculate(); err != nil {
		return nil, err
	}

	// 退款单处理
	if err := logic.Refund(); err != nil {
		return nil, err
	}

	return logic.Target, nil
}

// 获取目标
func (l *achieveLogic) getTarget(target_id string) error {
	if l.Target != nil {
		return nil
	}
	var target model.Target
	db := l.db.Where("id = ?", target_id)
	db = target.Preloads(db)
	if err := db.First(&target).Error; err != nil {
		return errors.New("目标不存在")
	}

	l.Target = &target

	l.start_time = target.StartTime.Format(time.RFC3339)
	l.end_time = target.EndTime.Format(time.RFC3339)

	return nil
}

// 获取销售单
func (l *achieveLogic) getSales() error {
	var sales model.OrderSales
	db := l.db.Where(&model.OrderSales{
		StoreId: l.Target.StoreId,
	})
	db = db.Where("status in (?)", []enums.OrderSalesStatus{
		enums.OrderSalesStatusComplete,
		enums.OrderSalesStatusRefund,
	})
	db = db.Scopes(model.DurationCondition(enums.DurationCustom, "created_at", l.start_time, l.end_time))

	db = sales.Preloads(db)
	if err := db.Find(&l.Sales).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return errors.New("销售单获取失败")
		}
	}

	return nil
}

// 获取退款单
func (l *achieveLogic) getRefunds() error {
	var refunds model.OrderRefund
	db := l.db.Where(&model.OrderRefund{
		StoreId:   l.Target.StoreId,
		OrderType: enums.OrderTypeSales,
	})
	db = db.Scopes(model.DurationCondition(enums.DurationCustom, "created_at", l.start_time, l.end_time))
	db = refunds.Preloads(db)
	if err := db.Find(&l.Refunds).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return errors.New("退款单获取失败")
		}
	}

	return nil
}

// 计算销售单业绩
func (l *achieveLogic) Calculate() error {
	// 计算业绩
	for _, order := range l.Sales {
		for _, product := range order.Products {
			for _, clerk := range order.Clerks {
				switch product.Type {
				case enums.ProductTypeFinished:
					{
						amount, quantity := l.calculateFinished(&product.Finished.Product, &clerk, product.Finished.Price)
						l.addAchieve(clerk.SalesmanId, amount, quantity)
					}
				case enums.ProductTypeOld:
					{
						amount, quantity := l.calculateOld(&order, &product.Old.Product, &clerk, product.Old.RecyclePrice)
						l.addAchieve(clerk.SalesmanId, amount.Neg(), quantity.Neg())
					}
				case enums.ProductTypeAccessorie:
					{
						amount, quantity := l.calculateAccessorie(&clerk, product.Accessorie.Price, product.Accessorie.Quantity)
						l.addAchieve(clerk.SalesmanId, amount, quantity)
					}
				}
			}
		}
	}

	return nil
}

// 计算成品
func (l *achieveLogic) calculateFinished(finished *model.ProductFinished, clerk *model.OrderSalesClerk, price decimal.Decimal) (amount, quantity decimal.Decimal) {
	// 判断统计范围
	scopes := []enums.TargetScope{
		enums.TargetScopeClass,
		enums.TargetScopeCategory,
		enums.TargetScopeAll,
	}
	if in := utils.ArrayFindIn(scopes, l.Target.Scope); !in {
		return
	}
	switch l.Target.Scope {
	case enums.TargetScopeClass:
		{
			if in := utils.ArrayFindIn(l.Target.Class, finished.Class); !in {
				return
			}
		}
	case enums.TargetScopeCategory:
		{
			if len(l.Target.Material) > 0 {
				if in := utils.ArrayFindIn(l.Target.Material, finished.Material); !in {
					return
				}
			}
			if len(l.Target.Quality) > 0 {
				if in := utils.ArrayFindIn(l.Target.Quality, finished.Quality); !in {
					return
				}
			}
			if len(l.Target.Category) > 0 {
				if in := utils.ArrayFindIn(l.Target.Category, finished.Category); !in {
					return
				}
			}
			if len(l.Target.Gem) > 0 {
				if in := utils.ArrayFindIn(l.Target.Gem, finished.Gem); !in {
					return
				}
			}
			if len(l.Target.Craft) > 0 {
				if in := utils.ArrayFindIn(l.Target.Craft, finished.Craft); !in {
					return
				}
			}
		}
	}

	amount = price.Mul(clerk.PerformanceRate).Div(decimal.NewFromFloat(100))
	quantity = decimal.NewFromInt(1).Mul(clerk.PerformanceRate).Div(decimal.NewFromFloat(100))

	return
}

// 计算旧料（仅限当前订单范围）
func (l *achieveLogic) calculateOld(order *model.OrderSales, old *model.ProductOld, clerk *model.OrderSalesClerk, price decimal.Decimal) (amount, quantity decimal.Decimal) {
	// 判断统计范围
	scopes := []enums.TargetScope{
		enums.TargetScopeClass,
		enums.TargetScopeCategory,
		enums.TargetScopeAll,
	}
	if in := utils.ArrayFindIn(scopes, l.Target.Scope); !in {
		return
	}

	if old.RecycleType != enums.ProductRecycleTypeExchange {
		return
	}

	if len(old.ExchangeFinisheds) == 0 {
		return
	}

	for _, p := range order.Products {
		if p.Type != enums.ProductTypeFinished {
			continue
		}

		if in := utils.ArrayFindIn(old.ExchangeFinisheds, p.Code); !in {
			continue
		}

		a, q := l.calculateFinished(&p.Finished.Product, clerk, p.Finished.Price)
		if a.IsZero() && q.IsZero() {
			continue
		}

		amount = price.Div(decimal.NewFromInt(int64(len(old.ExchangeFinisheds)))).Mul(clerk.PerformanceRate).Div(decimal.NewFromInt(100))
		quantity = decimal.NewFromInt(1).Mul(clerk.PerformanceRate).Div(decimal.NewFromInt(100))

		break
	}

	return
}

// 计算配件
func (l *achieveLogic) calculateAccessorie(clerk *model.OrderSalesClerk, price decimal.Decimal, sum int64) (amount, quantity decimal.Decimal) {
	// 判断统计范围
	scopes := []enums.TargetScope{
		enums.TargetScopeAccessorie,
		enums.TargetScopeAll,
	}
	if in := utils.ArrayFindIn(scopes, l.Target.Scope); !in {
		return
	}

	amount = price.Mul(clerk.PerformanceRate).Div(decimal.NewFromFloat(100))
	quantity = decimal.NewFromInt(sum).Mul(clerk.PerformanceRate).Div(decimal.NewFromFloat(100))

	return
}

// 退款单处理
func (l *achieveLogic) Refund() error {
	for _, refund := range l.Refunds {
		var order model.OrderSales
		db := l.db.Model(&model.OrderSales{})
		db = order.Preloads(db)
		if err := db.First(&order, "id = ?", refund.OrderId).Error; err != nil {
			return errors.New("销售单获取失败")
		}

		switch refund.Type {
		case enums.ProductTypeFinished:
			{
				var product model.ProductFinished
				if err := l.db.Where(&model.ProductFinished{
					Code: refund.Code,
				}).First(&product).Error; err != nil {
					return errors.New("成品获取失败")
				}

				for _, clerk := range order.Clerks {
					amount, quantity := l.calculateFinished(&product, &clerk, refund.Price)
					l.addAchieve(clerk.SalesmanId, amount.Neg(), quantity.Neg())
				}
			}
		case enums.ProductTypeOld:
			{
				var product model.ProductOld
				if err := l.db.Where(&model.ProductOld{
					Code: refund.Code,
				}).First(&product).Error; err != nil {
					return errors.New("旧料获取失败")
				}

				for _, clerk := range order.Clerks {
					amount, quantity := l.calculateOld(&order, &product, &clerk, refund.Price)
					l.addAchieve(clerk.SalesmanId, amount, quantity)
				}
			}
		case enums.ProductTypeAccessorie:
			{
				for _, clerk := range order.Clerks {
					amount, quantity := l.calculateAccessorie(&clerk, refund.Price, refund.Quantity)
					l.addAchieve(clerk.SalesmanId, amount.Neg(), quantity.Neg())
				}
			}
		}
	}

	return nil
}

// 添加业绩
func (l *achieveLogic) addAchieve(staff_id string, amount, quantity decimal.Decimal) {
	if amount.IsZero() && quantity.IsZero() {
		return
	}
	for i, personal := range l.Target.Personals {
		if personal.StaffId == staff_id {
			switch l.Target.Method {
			case enums.TargetMethodAmount:
				{
					l.Target.Personals[i].Achieve = personal.Achieve.Add(amount.Round(2))
				}
			case enums.TargetMethodQuantity:
				{
					l.Target.Personals[i].Achieve = personal.Achieve.Add(quantity.Round(2))
				}
			}
		}
	}
}
