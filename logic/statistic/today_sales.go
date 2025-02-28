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
	Req *types.StatisticTodaySalesReq
	Res *TodaySalesRes
	Db  *gorm.DB
}

func (l *StatisticLogic) TodaySales(req *types.StatisticTodaySalesReq) (*TodaySalesRes, error) {
	logic := &TodaySalesLogic{
		Req: req,
		Res: &TodaySalesRes{},
		Db:  model.DB,
	}

	// 获取金价
	if err := logic.getGoldPrice(); err != nil {
		return nil, err
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
	price, err := model.GetGoldPrice()
	if err != nil {
		return errors.New("获取金价失败")
	}
	l.Res.GoldPrice = price

	return nil
}

// 获取今日销售数据
func (l *TodaySalesLogic) getTodaySales() error {
	var sales_amount sql.NullFloat64
	if err := l.Db.Model(&model.Order{}).
		Where(&model.Order{
			StoreId: l.Req.StoreId,
			Type:    enums.OrderTypeSales,
			Status:  enums.OrderStatusComplete,
		}).
		Select("sum(amount_pay) as sales_amount").
		Scan(&sales_amount).Error; err != nil {
		return errors.New("获取今日销售数据失败")
	}

	if sales_amount.Valid {
		l.Res.SalesAmount = decimal.NewFromFloat(sales_amount.Float64)
	} else {
		l.Res.SalesAmount = decimal.Zero
	}

	return nil
}

// 获取今日销售件数
func (l *TodaySalesLogic) getTodaySalesCount() error {
	var orders []model.Order
	if err := l.Db.Model(&model.Order{}).
		Where(&model.Order{
			StoreId: l.Req.StoreId,
			Type:    enums.OrderTypeSales,
			Status:  enums.OrderStatusComplete,
		}).
		Scopes(model.DurationCondition(enums.DurationToday)).
		Preload("Products").
		Find(&orders).Error; err != nil {
		return errors.New("获取今日销售件数失败")
	}

	sales_count := int64(0)
	for _, order := range orders {
		for _, product := range order.Products {
			sales_count += product.Quantity
		}
	}

	l.Res.SalesCount = sales_count

	return nil
}

// 获取旧货抵值
func (l *TodaySalesLogic) getOldGoodsAmount() error {
	var old_goods_amount sql.NullFloat64
	if err := l.Db.Model(&model.Order{}).
		Where(&model.Order{
			StoreId: l.Req.StoreId,
			Type:    enums.OrderTypeSales,
			Status:  enums.OrderStatusComplete,
		}).
		Scopes(model.DurationCondition(enums.DurationToday)).
		Select("sum(amount_old_material) as sales_amount").
		Scan(&old_goods_amount).Error; err != nil {
		return errors.New("获取今日销售数据失败")
	}

	if old_goods_amount.Valid {
		l.Res.OldGoodsAmount = decimal.NewFromFloat(old_goods_amount.Float64)
	} else {
		l.Res.OldGoodsAmount = decimal.Zero
	}

	return nil
}

// 获取退货金额
func (l *TodaySalesLogic) getReturnAmount() error {
	var return_amount sql.NullFloat64
	// if err := l.Db.Model(&model.OrderProduct{}).
	// 	Where(&model.OrderProduct{
	// 		Status: enums.OrderStatusRefund,
	// 	}).
	//     // 查询订单
	// 	Where("id IN (SELECT id FROM order_products WHERE order_id IN (SELECT id FROM orders WHERE store_id = ? ))", l.Req.StoreId).
	//  Scopes(model.DurationCondition(enums.DurationToday)).
	// 	Select("sum(amount) as return_amount").
	// 	Scan(&return_amount).Error; err != nil {
	// 	return errors.New("获取退货金额失败")
	// }

	if return_amount.Valid {
		l.Res.ReturnAmount = decimal.NewFromFloat(return_amount.Float64)
	} else {
		l.Res.ReturnAmount = decimal.Zero
	}

	return nil
}
