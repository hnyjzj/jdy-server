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

	// 获取销售数据
	if err := logic.getSales(onlyself); err != nil {
		return nil, err
	}
	if err := logic.getRefund(onlyself); err != nil {
		return nil, err
	}

	// 获取金价
	if err := logic.get_gold_price(); err != nil {
		return nil, err
	}

	// 获取销售数据
	if err := logic.get_sales_amount(); err != nil {
		return nil, err
	}

	// 获取销售件数
	if err := logic.get_sales_count(); err != nil {
		return nil, err
	}

	// 获取旧料抵值
	if err := logic.get_old_goods_amount(); err != nil {
		return nil, err
	}

	// 获取退货金额
	if err := logic.get_return_amount(); err != nil {
		return nil, err
	}

	return logic.Res, nil
}

// 获取销售数据
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

func (l *SalesLogic) get_sales_amount() error {
	price, ok := l.Res["销售金额"].(decimal.Decimal)
	if !ok {
		price = decimal.Zero
	}
	for _, o := range l.Sales {
		price = price.Add(o.ProductFinishedPrice)
	}
	l.Res["销售金额"] = price

	return nil
}

// 获取销售件数
func (l *SalesLogic) get_sales_count() error {
	count, ok := l.Res["销售件数"].(int)
	if !ok {
		count = 0
	}
	for _, o := range l.Sales {
		for _, p := range o.Products {
			if p.Type == enums.ProductTypeFinished {
				count++
			}
		}
	}
	l.Res["销售件数"] = count

	return nil
}

// 获取旧料抵值
func (l *SalesLogic) get_old_goods_amount() error {
	price, ok := l.Res["旧料抵值"].(decimal.Decimal)
	if !ok {
		price = decimal.Zero
	}
	for _, o := range l.Sales {
		price = price.Add(o.ProductOldPrice)
	}
	l.Res["旧料抵值"] = price

	return nil
}

// 获取退货金额
func (l *SalesLogic) get_return_amount() error {
	price, ok := l.Res["退货金额"].(decimal.Decimal)
	if !ok {
		price = decimal.Zero
	}
	for _, o := range l.Refund {
		price = price.Add(o.Price)
	}
	l.Res["退货金额"] = price

	return nil
}
