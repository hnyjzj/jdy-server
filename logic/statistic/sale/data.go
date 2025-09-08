package sale

import (
	"errors"
	"jdy/enums"
	"jdy/model"
	"time"

	"github.com/shopspring/decimal"
)

type datawLogic struct {
	*StatisticSaleLogic
	req *DataReq

	Sales  []model.OrderSales
	Refund []model.OrderRefund
}

type DataRes struct {
	Overview map[string]any            `json:"overview"` // 总览
	Trend    map[string]map[string]any `json:"trend"`    // 趋势
}

func (l *StatisticSaleLogic) Data(req *DataReq, onlyself bool) (any, error) {
	logic := datawLogic{
		StatisticSaleLogic: l,
		req:                req,
	}

	if err := logic.get_sales(onlyself); err != nil {
		return nil, err
	}
	if err := logic.get_refund(onlyself); err != nil {
		return nil, err
	}

	res := DataRes{
		Overview: logic.get_overview(),
		Trend:    logic.get_trend(),
	}

	return res, nil
}

func (l *datawLogic) get_sales(onlyself bool) error {
	db := model.DB.Model(&model.OrderSales{})
	db = db.Where(&model.OrderSales{
		StoreId: l.req.StoreId,
	})
	db = db.Where("status in (?)", []enums.OrderSalesStatus{
		enums.OrderSalesStatusComplete,
		enums.OrderSalesStatusRefund,
	})
	db = db.Scopes(model.DurationCondition(l.req.Duration, "created_at", l.req.StartTime, l.req.EndTime))

	if onlyself {
		self := model.DB.Model(&model.OrderSalesClerk{})
		self = self.Where(&model.OrderSalesClerk{
			SalesmanId: l.Staff.Id,
		})
		self = self.Select("order_id").Group("order_id")
		self = self.Scopes(model.DurationCondition(l.req.Duration, "created_at", l.req.StartTime, l.req.EndTime))

		db = db.Where("id in (?)", self)
	}

	db = model.OrderSales{}.Preloads(db)
	if err := db.Find(&l.Sales).Error; err != nil {
		return errors.New("获取数据失败")
	}

	return nil
}
func (l *datawLogic) get_refund(onlyself bool) error {
	db := model.DB.Model(&model.OrderRefund{})
	db = db.Where(&model.OrderRefund{
		StoreId: l.req.StoreId,
	})
	db = db.Scopes(model.DurationCondition(l.req.Duration, "created_at", l.req.StartTime, l.req.EndTime))

	if onlyself {
		db = db.Where("operator_id = ?", l.Staff.Id)
	}

	db = model.OrderRefund{}.Preloads(db)
	if err := db.Find(&l.Refund).Error; err != nil {
		return errors.New("获取数据失败")
	}

	return nil
}

func (l *datawLogic) get_overview() map[string]any {
	data := make(map[string]any)

	if len(l.Sales) == 0 {
		data["销售金额"] = decimal.Zero
		data["销售件数"] = decimal.Zero
		data["旧料抵值"] = decimal.Zero
		data["配件礼品"] = decimal.Zero
	}
	if len(l.Refund) == 0 {
		data["退款金额"] = decimal.Zero
		data["退款件数"] = decimal.Zero
	}

	for _, s := range l.Sales {
		price, ok := data["销售金额"].(decimal.Decimal)
		if !ok {
			price = decimal.Zero
		}
		price = price.Add(s.ProductFinishedPrice)
		data["销售金额"] = price

		count, ok := data["销售件数"].(int64)
		if !ok {
			count = 0
		}
		for _, p := range s.Products {
			if p.Type == enums.ProductTypeFinished {
				count = count + 1
			}
		}
		data["销售件数"] = count

		old, ok := data["旧料抵值"].(decimal.Decimal)
		if !ok {
			old = decimal.Zero
		}
		old = old.Add(s.ProductOldPrice)
		data["旧料抵值"] = old

		accessorie, ok := data["配件礼品"].(decimal.Decimal)
		if !ok {
			accessorie = decimal.Zero
		}
		accessorie = accessorie.Add(s.ProductAccessoriePrice)
		data["配件礼品"] = accessorie
	}

	for _, r := range l.Refund {
		price, ok := data["退款金额"].(decimal.Decimal)
		if !ok {
			price = decimal.Zero
		}
		price = price.Add(r.Price)
		data["退款金额"] = price

		count, ok := data["退款件数"].(int64)
		if !ok {
			count = 0
		}
		count = count + 1
		data["退款件数"] = count
	}

	return data
}

func (l *datawLogic) get_trend() map[string]map[string]any {
	data := make(map[string]map[string]any)

	now := time.Now()
	start, end, err := l.req.Duration.GetTime(now, l.req.StartTime, l.req.EndTime)
	if err != nil {
		return data
	}

	_, list := get_date_format(start, end, start)
	for _, v := range list {
		if _, ok := data[v]; !ok {
			data[v] = map[string]any{
				"销售额": decimal.Zero,
				"件数":  0,
			}
		}
	}

	for _, order := range l.Sales {
		k, _ := get_date_format(start, end, *order.CreatedAt)
		if _, ok := data[k]; !ok {
			data[k] = make(map[string]any)
		}

		price, ok := data[k]["销售额"].(decimal.Decimal)
		if !ok {
			price = decimal.Zero
		}
		price = price.Add(order.ProductFinishedPrice)
		data[k]["销售额"] = price
		num, ok := data[k]["件数"].(int)
		if !ok {
			num = 0
		}
		for _, product := range order.Products {
			if product.Type == enums.ProductTypeFinished {
				num = num + 1
			}
		}
		data[k]["件数"] = num
	}

	return data
}

func get_date_format(start, end, field time.Time) (string, []string) {
	var (
		k    string
		list []string
	)

	now := time.Now()
	span := end.Sub(start)

	switch {
	case span.Hours() <= 24: // 每小时
		{
			k = field.Format("15:00:00")
			for i := start; i.Before(end) && i.Before(now); i = i.Add(time.Hour) {
				list = append(list, i.Format("15:00:00"))
			}
		}
	case span.Hours() > 24 && span.Hours() <= 24*180: // 每天
		{
			k = field.Format(time.DateOnly)
			for i := start; i.Before(end) && i.Before(now); i = i.Add(time.Hour * 24) {
				list = append(list, i.Format(time.DateOnly))
			}
		}
	default: // 每月
		{
			k = field.Format("2006-01")
			current := start
			for current.Before(end) && current.Before(now) {
				nextYear := current.Year()
				nextMonth := current.Month() + 1
				if nextMonth > 12 {
					nextYear++
					nextMonth = 1
				}
				next := time.Date(nextYear, nextMonth, 1, 0, 0, 0, 0, current.Location())
				list = append(list, current.Format("2006-01"))
				current = next
			}
		}
	}

	return k, list
}
