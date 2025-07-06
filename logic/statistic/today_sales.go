package statistic

import (
	"database/sql"
	"errors"
	"jdy/enums"
	"jdy/model"
	"jdy/types"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type TodaySalesRes struct {
	GoldPrice      decimal.Decimal `json:"gold_price"`       // 金价
	SalesAmount    decimal.Decimal `json:"sales_amount"`     // 销售金额
	SalesCount     int64           `json:"sales_count"`      // 销售件数
	OldGoodsAmount decimal.Decimal `json:"old_goods_amount"` // 旧货抵值
	ReturnAmount   decimal.Decimal `json:"return_amount"`    // 退货金额
}

type TodaySalesLogic struct {
	*StatisticLogic

	Req *types.StatisticTodaySalesReq
	Res *TodaySalesRes
	Db  *gorm.DB

	clerk_query *gorm.DB
}

func (l *StatisticLogic) TodaySales(req *types.StatisticTodaySalesReq) (*TodaySalesRes, error) {
	logic := &TodaySalesLogic{
		StatisticLogic: l,
		Req:            req,
		Res:            &TodaySalesRes{},
		Db:             model.DB,
	}

	// 获取金价
	if err := logic.getGoldPrice(); err != nil {
		return nil, err
	}

	// 如果是店员，则仅查询该店员的订单
	if l.Staff.Identity == enums.IdentityClerk {
		logic.clerk_query = logic.Db.Model(&model.OrderSalesClerk{}).Where(&model.OrderSalesClerk{
			SalesmanId: l.Staff.Id,
		}).Scopes(model.DurationCondition(enums.DurationToday)).Select("order_id").Group("order_id")
	}

	// 获取今日销售数据
	if err := logic.getTodaySales(); err != nil {
		return nil, err
	}

	// 获取今日销售件数
	if err := logic.getTodaySalesCount(); err != nil {
		return nil, err
	}

	// 获取旧货抵值
	if err := logic.getOldGoodsAmount(); err != nil {
		return nil, err
	}

	// 获取退货金额
	if err := logic.getReturnAmount(); err != nil {
		return nil, err
	}

	return logic.Res, nil
}

// 获取金价
func (l *TodaySalesLogic) getGoldPrice() error {
	price, _ := model.GetGoldPrice(&types.GoldPriceOptions{
		StoreId: l.Req.StoreId,
	})
	l.Res.GoldPrice = price

	return nil
}

// 获取今日销售数据
func (l *TodaySalesLogic) getTodaySales() error {
	var (
		sales_amount sql.NullFloat64
		db           = l.Db.Model(&model.OrderSales{})
	)

	// 如果是店员，则仅查询该店员的订单
	if l.clerk_query != nil {
		db = db.Where("id IN (?)", l.clerk_query)
	}

	// 查询今日订单
	db = db.Where(&model.OrderSales{
		StoreId: l.Req.StoreId,
		Status:  enums.OrderSalesStatusComplete,
	}).Scopes(model.DurationCondition(enums.DurationToday))

	// 查询今日销售金额
	if err := db.Select("sum(price_pay) as sales_amount").Scan(&sales_amount).Error; err != nil {
		return errors.New("获取今日销售数据失败")
	}

	// 判断金额
	if sales_amount.Valid {
		l.Res.SalesAmount = decimal.NewFromFloat(sales_amount.Float64)
	} else {
		l.Res.SalesAmount = decimal.Zero
	}

	return nil
}

// 获取今日销售件数
func (l *TodaySalesLogic) getTodaySalesCount() error {
	var (
		db = l.Db.Model(&model.OrderSalesProductFinished{})
	)

	// 如果是店员，则仅查询该店员的订单
	if l.clerk_query != nil {
		db = db.Where("order_id IN (?)", l.clerk_query)
	}

	// 查询今日订单
	order_query := l.Db.Model(&model.OrderSales{}).Where(&model.OrderSales{
		StoreId: l.Req.StoreId,
		Status:  enums.OrderSalesStatusComplete,
	}).Scopes(model.DurationCondition(enums.DurationToday)).Select("id").Group("id")
	db = db.Where("order_id IN (?)", order_query)

	// 查询今日销售件数
	db = db.Where(&model.OrderSalesProductFinished{
		Status: enums.OrderSalesStatusComplete,
	}).Scopes(model.DurationCondition(enums.DurationToday))
	if err := db.Debug().Count(&l.Res.SalesCount).Error; err != nil {
		return errors.New("获取今日销售件数失败")
	}

	return nil
}

// 获取旧货抵值
func (l *TodaySalesLogic) getOldGoodsAmount() error {
	var (
		old_goods_amount sql.NullFloat64
		db               = l.Db.Model(&model.OrderSales{})
	)

	// 如果是店员，则仅查询该店员的订单
	if l.clerk_query != nil {
		db = db.Where("id IN (?)", l.clerk_query)
	}

	// 查询今日订单
	db = db.Scopes(model.DurationCondition(enums.DurationToday))

	// 查询本门店已完成的订单
	db = db.Where(&model.OrderSales{
		StoreId: l.Req.StoreId,
		Status:  enums.OrderSalesStatusComplete,
	})

	// 查询今日销售金额
	if err := db.Select("sum(product_old_price) as sales_amount").Scan(&old_goods_amount).Error; err != nil {
		return errors.New("获取今日销售数据失败")
	}

	// 判断金额
	if old_goods_amount.Valid {
		l.Res.OldGoodsAmount = decimal.NewFromFloat(old_goods_amount.Float64)
	} else {
		l.Res.OldGoodsAmount = decimal.Zero
	}

	return nil
}

// 获取退货金额
func (l *TodaySalesLogic) getReturnAmount() error {
	var (
		return_amount sql.NullFloat64
		db            = l.Db.Model(&model.OrderRefund{})
	)

	// 如果是店员，则仅查询该店员的订单
	if l.clerk_query != nil {
		db = db.Where("order_id IN (?)", l.clerk_query)
	}

	// 查询今日订单
	db = db.Scopes(model.DurationCondition(enums.DurationToday))
	// 查询本门店销售单的退货订单
	db = db.Where(&model.OrderRefund{
		StoreId:   l.Req.StoreId,
		OrderType: enums.OrderTypeSales,
	})

	if err := db.Select("sum(price) as return_amount").Scan(&return_amount).Error; err != nil {
		return errors.New("获取退货金额失败")
	}

	if return_amount.Valid {
		l.Res.ReturnAmount = decimal.NewFromFloat(return_amount.Float64)
	} else {
		l.Res.ReturnAmount = decimal.Zero
	}

	return nil
}
