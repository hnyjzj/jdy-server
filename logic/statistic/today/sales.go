package today

import (
	"errors"
	"jdy/enums"
	"jdy/model"
	"jdy/types"

	"github.com/shopspring/decimal"
)

type SalesReq struct {
	DataReq
	StoreId string `json:"store_id"`
}

type SalesLogic struct {
	*ToDayLogic

	Req    *SalesReq
	Sales  []model.OrderSales
	Refund []model.OrderRefund

	Res map[string]any
}

func (l *ToDayLogic) Sales(req *SalesReq, onlyself bool) (map[string]any, error) {
	logic := &SalesLogic{
		ToDayLogic: l,
		Req:        req,
		Res:        make(map[string]any),
	}

	// 获取成品数据
	if err := logic.getSales(onlyself); err != nil {
		return nil, err
	}
	// 获取退款数据
	if err := logic.getRefund(onlyself); err != nil {
		return nil, err
	}

	// 获取金价
	if err := logic.get_gold_price(); err != nil {
		return nil, err
	}

	// 获取成品数据
	if err := logic.get_finisheds(); err != nil {
		return nil, err
	}

	// 获取旧料抵值
	if err := logic.get_olds(); err != nil {
		return nil, err
	}

	// 获取配件礼品
	if err := logic.get_accessories(); err != nil {
		return nil, err
	}

	// 获取退货金额
	if err := logic.get_return_amount(); err != nil {
		return nil, err
	}

	return logic.Res, nil
}

// 获取成品数据
func (l *SalesLogic) getSales(onlyself bool) error {
	db := model.DB.Model(&model.OrderSales{})
	db = db.Where(&model.OrderSales{
		StoreId: l.Req.StoreId,
	})
	db = db.Where("status in (?)", []enums.OrderSalesStatus{
		enums.OrderSalesStatusComplete,
		enums.OrderSalesStatusRefund,
	})
	db = db.Scopes(model.DurationCondition(l.Req.Duration, "created_at", l.Req.StartTime, l.Req.EndTime))

	if onlyself {
		self := model.DB.Model(&model.OrderSalesClerk{})
		self = self.Where(&model.OrderSalesClerk{
			SalesmanId: l.Staff.Id,
		})
		self = self.Select("order_id").Group("order_id")
		self = self.Scopes(model.DurationCondition(l.Req.Duration, "created_at", l.Req.StartTime, l.Req.EndTime))

		db = db.Where("id in (?)", self)
	}

	db = model.OrderSales{}.Preloads(db)
	if err := db.Find(&l.Sales).Error; err != nil {
		return errors.New("获取数据失败")
	}

	return nil
}

// 获取退款数据
func (l *SalesLogic) getRefund(onlyself bool) error {
	db := model.DB.Model(&model.OrderRefund{})
	db = db.Where(&model.OrderRefund{
		StoreId: l.Req.StoreId,
	})
	db = db.Scopes(model.DurationCondition(l.Req.Duration, "created_at", l.Req.StartTime, l.Req.EndTime))

	if onlyself {
		db = db.Where("operator_id = ?", l.Staff.Id)
	}

	db = model.OrderRefund{}.Preloads(db)
	if err := db.Find(&l.Refund).Error; err != nil {
		return errors.New("获取数据失败")
	}

	return nil
}

// 获取金价
func (l *SalesLogic) get_gold_price() error {
	price, _ := model.GetGoldPrice(&types.GoldPriceOptions{
		StoreId: l.Req.StoreId,
	})
	l.Res["金价"] = price

	return nil
}

func (l *SalesLogic) get_finisheds() error {
	price, ok := l.Res["成品金额"].(decimal.Decimal)
	if !ok {
		price = decimal.Zero
	}
	num, ok := l.Res["成品件数"].(int)
	if !ok {
		num = 0
	}

	for _, sale := range l.Sales {
		for _, product := range sale.Products {
			if product.Type == enums.ProductTypeFinished {
				price = price.Add(product.Finished.Price)
				num++
			}
			for _, refund := range l.Refund {
				if refund.Type != enums.ProductTypeFinished {
					continue
				}
				if refund.OrderId != sale.Id {
					continue
				}
				if refund.Code != product.Code {
					continue
				}
				price = price.Sub(refund.Price)
				num--
			}
		}
	}

	l.Res["成品金额"] = price
	l.Res["成品件数"] = num

	return nil
}

// 获取旧料抵值
func (l *SalesLogic) get_olds() error {
	price, ok := l.Res["旧料抵值"].(decimal.Decimal)
	if !ok {
		price = decimal.Zero
	}

	for _, sale := range l.Sales {
		for _, product := range sale.Products {
			if product.Type == enums.ProductTypeOld {
				price = price.Add(product.Old.RecyclePrice.Neg())
			}
			for _, refund := range l.Refund {
				if refund.Type != enums.ProductTypeOld {
					continue
				}
				if refund.OrderId != sale.Id {
					continue
				}
				if refund.Code != product.Code {
					continue
				}
				price = price.Sub(refund.Price.Neg())
			}
		}
	}

	l.Res["旧料抵值"] = price

	return nil
}

// 获取配件礼品
func (l *SalesLogic) get_accessories() error {
	price, ok := l.Res["配件礼品"].(decimal.Decimal)
	if !ok {
		price = decimal.Zero
	}

	for _, sale := range l.Sales {
		for _, product := range sale.Products {
			if product.Type == enums.ProductTypeAccessorie {
				price = price.Add(product.Accessorie.Price)
			}
			for _, refund := range l.Refund {
				if refund.Type != enums.ProductTypeAccessorie {
					continue
				}
				if refund.OrderId != sale.Id {
					continue
				}
				if refund.Name != product.Name {
					continue
				}
				price = price.Sub(refund.Price)
			}
		}
	}

	l.Res["配件礼品"] = price

	return nil
}

// 获取退货金额
func (l *SalesLogic) get_return_amount() error {
	price, ok := l.Res["退货金额"].(decimal.Decimal)
	if !ok {
		price = decimal.Zero
	}
	for _, refund := range l.Refund {
		switch refund.Type {
		case enums.ProductTypeFinished, enums.ProductTypeAccessorie:
			{
				price = price.Sub(refund.Price)
			}
		case enums.ProductTypeOld:
			{
				price = price.Add(refund.Price)
			}
		}
	}
	l.Res["退货金额"] = price

	return nil
}
