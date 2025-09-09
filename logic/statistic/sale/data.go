package sale

import (
	"errors"
	"jdy/enums"
	"jdy/model"
	"time"

	"github.com/shopspring/decimal"
)

type dataLogic struct {
	*StatisticSaleLogic
	req *DataReq

	Sales  []model.OrderSales
	Refund []model.OrderRefund
}

type DataRes struct {
	Overview         map[string]any            `json:"overview"`          // 总览
	Trend            map[string]map[string]any `json:"trend"`             // 趋势
	FinishedClass    map[string]any            `json:"finished_class"`    // 成品大类
	FinishedCategory map[string]map[string]any `json:"finished_category"` // 成品品类
	OldClass         map[string]any            `json:"old_class"`         // 旧料大类
	Accessorie       map[string]any            `json:"accessorie"`        // 配件
	List             map[string]any            `json:"list"`              // 列表
}

func (l *StatisticSaleLogic) Data(req *DataReq, onlyself bool) (any, error) {
	logic := dataLogic{
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
		Overview:         logic.get_overview(),
		Trend:            logic.get_trend(),
		FinishedClass:    logic.get_finished_class(),
		FinishedCategory: logic.get_finished_category(),
		OldClass:         logic.get_old_class(),
		Accessorie:       logic.get_accessorie(),
		List:             logic.get_list(),
	}

	return res, nil
}

func (l *dataLogic) get_sales(onlyself bool) error {
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
func (l *dataLogic) get_refund(onlyself bool) error {
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

func (l *dataLogic) get_overview() map[string]any {
	data := make(map[string]any)

	if len(l.Sales) == 0 {
		data["销售金额"] = decimal.Zero
		data["销售件数"] = 0
		data["旧料抵值"] = decimal.Zero
		data["配件礼品"] = decimal.Zero
	}
	if len(l.Refund) == 0 {
		data["退款金额"] = decimal.Zero
		data["退款件数"] = 0
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

func (l *dataLogic) get_trend() map[string]map[string]any {
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

// 获取成品大类数据
func (l *dataLogic) get_finished_class() map[string]any {
	data := make(map[string]any)

	for _, order := range l.Sales {
		for _, product := range order.Products {
			if product.Type != enums.ProductTypeFinished {
				continue
			}
			k := product.Finished.Product.Class.String()
			if k == "" {
				k = "其他"
			}

			num_row, ok := data["件数"].(map[string]any)
			if !ok {
				num_row = make(map[string]any, 0)
			}
			num, ok := num_row[k].(int64)
			if !ok {
				num = 0
			}
			num = num + 1
			num_row[k] = num
			data["件数"] = num_row

			price_row, ok := data["销售额"].(map[string]any)
			if !ok {
				price_row = make(map[string]any, 0)
			}
			price, ok := price_row[k].(decimal.Decimal)
			if !ok {
				price = decimal.NewFromInt(0)
			}
			price = price.Add(product.Finished.Price)
			price_row[k] = price
			data["销售额"] = price_row
		}
	}

	return data
}

// 获取成品品类数据
func (l *dataLogic) get_finished_category() map[string]map[string]any {

	data := make(map[string]map[string]any)

	for _, class := range enums.ProductClassFinishedMap {
		data[class] = map[string]any{
			"件数":  map[string]any{},
			"销售额": map[string]any{},
			"金重":  map[string]any{},
		}
	}

	for _, order := range l.Sales {
		for _, product := range order.Products {
			if product.Type != enums.ProductTypeFinished {
				continue
			}

			c := product.Finished.Product.Class.String()
			if c == "" {
				c = "其他"
			}
			k := product.Finished.Product.Category.String()
			if k == "" {
				k = "其他"
			}

			num_item, ok := data[c]["件数"].(map[string]any)
			if !ok {
				num_item = make(map[string]any, 0)
			}
			num, ok := num_item[k].(decimal.Decimal)
			if !ok {
				num = decimal.Zero
			}
			num = num.Add(decimal.NewFromInt(1))
			num_item[k] = num
			data[c]["件数"] = num_item

			price_item, ok := data[c]["销售额"].(map[string]any)
			if !ok {
				price_item = make(map[string]any, 0)
			}
			price, ok := price_item[k].(decimal.Decimal)
			if !ok {
				price = decimal.Zero
			}
			price = price.Add(product.Finished.Price)
			price_item[k] = price
			data[c]["销售额"] = price_item

			weight_item, ok := data[c]["金重"].(map[string]any)
			if !ok {
				weight_item = make(map[string]any, 0)
			}
			weight, ok := weight_item[k].(decimal.Decimal)
			if !ok {
				weight = decimal.Zero
			}
			weight = weight.Add(product.Finished.Product.WeightMetal)
			weight_item[k] = weight
			data[c]["金重"] = weight_item
		}
	}

	return data
}

// 获取旧料大类数据
func (l *dataLogic) get_old_class() map[string]any {
	data := make(map[string]any)

	for _, order := range l.Sales {
		for _, product := range order.Products {
			if product.Type != enums.ProductTypeOld {
				continue
			}
			k := product.Old.Product.Class.String()
			if k == "" {
				k = "其他"
			}

			price_row, ok := data["抵值"].(map[string]any)
			if !ok {
				price_row = make(map[string]any, 0)
			}
			price, ok := price_row[k].(decimal.Decimal)
			if !ok {
				price = decimal.NewFromInt(0)
			}
			price = price.Add(product.Old.RecyclePrice)
			price_row[k] = price
			data["抵值"] = price_row

			weight_item, ok := data["金重"].(map[string]any)
			if !ok {
				weight_item = make(map[string]any, 0)
			}
			weight, ok := weight_item[k].(decimal.Decimal)
			if !ok {
				weight = decimal.Zero
			}
			weight = weight.Add(product.Old.WeightMetal)
			weight_item[k] = weight
			data["金重"] = weight_item

			num_row, ok := data["件数"].(map[string]any)
			if !ok {
				num_row = make(map[string]any, 0)
			}
			num, ok := num_row[k].(int64)
			if !ok {
				num = 0
			}
			num = num + 1
			num_row[k] = num
			data["件数"] = num_row

		}
	}

	return data
}

// 获取配件数据
func (l *dataLogic) get_accessorie() map[string]any {
	data := make(map[string]any)

	for _, order := range l.Sales {
		for _, product := range order.Products {
			if product.Type != enums.ProductTypeAccessorie {
				continue
			}
			k := product.Accessorie.Product.Name + "(" + product.Accessorie.Product.Type.String() + ")"

			price_row, ok := data["销售额"].(map[string]any)
			if !ok {
				price_row = make(map[string]any, 0)
			}
			price, ok := price_row[k].(decimal.Decimal)
			if !ok {
				price = decimal.NewFromInt(0)
			}
			price = price.Add(product.Accessorie.Price)
			price_row[k] = price
			data["销售额"] = price_row

			num_row, ok := data["件数"].(map[string]any)
			if !ok {
				num_row = make(map[string]any, 0)
			}
			num, ok := num_row[k].(int64)
			if !ok {
				num = 0
			}
			num = num + product.Accessorie.Quantity
			num_row[k] = num
			data["件数"] = num_row

		}
	}

	return data
}

// 获取列表
func (l *dataLogic) get_list() map[string]any {
	data := make(map[string]any)

	for _, order := range l.Sales {
		for _, clerk := range order.Clerks {
			k := clerk.Salesman.Nickname

			row, ok := data[k].(map[string]any)
			if !ok {
				row = make(map[string]any, 0)
			}
			finished_price, ok := row["成品销售额"].(decimal.Decimal)
			if !ok {
				finished_price = decimal.NewFromInt(0)
			}

			accessorie_price, ok := row["配件销售额"].(decimal.Decimal)
			if !ok {
				accessorie_price = decimal.NewFromInt(0)
			}

			finished_num, ok := row["成品件数"].(int64)
			if !ok {
				finished_num = 0
			}

			accessorie_num, ok := row["配件件数"].(int64)
			if !ok {
				accessorie_num = 0
			}

			for _, product := range order.Products {
				switch product.Type {
				case enums.ProductTypeFinished:
					finished_price = finished_price.Add(product.Finished.Price.Mul(clerk.PerformanceRate).Div(decimal.NewFromFloat(100)))
					finished_num = finished_num + 1
				case enums.ProductTypeAccessorie:
					accessorie_price = accessorie_price.Add(product.Accessorie.Price.Mul(clerk.PerformanceRate).Div(decimal.NewFromFloat(100)))
					accessorie_num = accessorie_num + product.Accessorie.Quantity
				}
			}
			row["成品销售额"] = finished_price
			row["配件销售额"] = accessorie_price
			row["成品件数"] = finished_num
			row["配件件数"] = accessorie_num
			data[k] = row
		}
	}

	return data
}
