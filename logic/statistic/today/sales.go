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

type SalesRes struct {
	GoldPrice      decimal.Decimal `json:"gold_price"`       // 金价
	SalesAmount    decimal.Decimal `json:"sales_amount"`     // 销售金额
	SalesCount     int64           `json:"sales_count"`      // 销售件数
	OldGoodsAmount decimal.Decimal `json:"old_goods_amount"` // 旧料抵值
	ReturnAmount   decimal.Decimal `json:"return_amount"`    // 退货金额
}

type SalesLogic struct {
	*ToDayLogic

	Req    *SalesReq
	Sales  []model.OrderSales
	Refund []model.OrderRefund

	Res *SalesRes
}

func (l *ToDayLogic) Sales(req *SalesReq, onlyself bool) (*SalesRes, error) {
	logic := &SalesLogic{
		ToDayLogic: l,
		Req:        req,
		Res:        &SalesRes{},
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
	l.Res.GoldPrice = price

	return nil
}

func (l *SalesLogic) get_sales_amount() error {
	for _, o := range l.Sales {
		l.Res.SalesAmount = l.Res.SalesAmount.Add(o.ProductFinishedPrice)
	}

	return nil
}

// 获取销售件数
func (l *SalesLogic) get_sales_count() error {
	for _, o := range l.Sales {
		for _, p := range o.Products {
			if p.Type == enums.ProductTypeFinished {
				l.Res.SalesCount++
			}
		}
	}

	return nil
}

// 获取旧料抵值
func (l *SalesLogic) get_old_goods_amount() error {
	for _, o := range l.Sales {
		l.Res.OldGoodsAmount = l.Res.OldGoodsAmount.Add(o.ProductOldPrice)
	}

	return nil
}

// 获取退货金额
func (l *SalesLogic) get_return_amount() error {
	for _, o := range l.Refund {
		l.Res.ReturnAmount = l.Res.ReturnAmount.Add(o.Price)
	}

	return nil
}
