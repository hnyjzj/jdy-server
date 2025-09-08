package stock

import (
	"jdy/enums"
	"jdy/model"
	"time"

	"github.com/shopspring/decimal"
)

type dataLogic struct {
	*StatisticStockLogic

	storeId string
	endtime time.Time

	Finisheds []model.ProductFinished
	Olds      []model.ProductOld
}

type DataRes struct {
	Overview         map[string]any            `json:"overview"`          // 概览
	FinishedClass    map[string]any            `json:"finished_class"`    // 成品分类
	FinishedCategory map[string]map[string]any `json:"finished_category"` // 成品品类
	FinishedAge      map[string]map[string]any `json:"finished_age"`      // 成品件数
	OldClass         map[string]any            `json:"old_class"`         // 旧料分类
}

func (l *StatisticStockLogic) Data(req *DataReq) (any, error) {
	logic := dataLogic{
		StatisticStockLogic: l,
	}

	day, err := time.ParseInLocation(time.DateOnly, req.Day, time.Local)
	if err != nil {
		return nil, err
	}

	logic.storeId = req.StoreId
	logic.endtime = day.AddDate(0, 0, 1)

	if err := logic.get_finisheds(); err != nil {
		return nil, err
	}
	if err := logic.get_olds(); err != nil {
		return nil, err
	}

	res := DataRes{
		Overview:         logic.get_overview(),
		FinishedClass:    logic.get_finished_class(),
		FinishedCategory: logic.get_finished_category(),
		FinishedAge:      logic.get_finished_age(),
		OldClass:         logic.get_old_class(),
	}

	return res, nil
}

// 获取成品库存
func (l *dataLogic) get_finisheds() error {
	db := model.DB.Model(&model.ProductFinished{})
	db = db.Where(&model.ProductFinished{
		StoreId: l.storeId,
	})
	db = db.Where("enter_time <= ?", l.endtime)
	db = db.Where("status IN (?)", []enums.ProductStatus{
		enums.ProductStatusNormal,
		enums.ProductStatusAllocate,
	})

	if err := db.Find(&l.Finisheds).Error; err != nil {
		return err
	}

	return nil
}

// 获取旧料库存
func (l *dataLogic) get_olds() error {
	db := model.DB.Model(&model.ProductOld{})
	db = db.Where(&model.ProductOld{
		StoreId: l.storeId,
	})
	db = db.Where("created_at <= ?", l.endtime)
	db = db.Where("status IN (?)", []enums.ProductStatus{
		enums.ProductStatusNormal,
		enums.ProductStatusAllocate,
	})

	if err := db.Find(&l.Olds).Error; err != nil {
		return err
	}

	return nil
}

// 获取成品库存数据
func (l *dataLogic) get_overview() map[string]any {
	data := make(map[string]any)

	if len(l.Finisheds) == 0 {
		data["成品总件数"] = decimal.Zero
		data["成品总金重"] = decimal.Zero
		data["成品总标签价"] = decimal.Zero
		data["成品滞销件数"] = decimal.Zero
	}
	for _, finished := range l.Finisheds {
		number, ok := data["成品总件数"].(decimal.Decimal)
		if !ok {
			number = decimal.Zero
		}
		number = number.Add(decimal.NewFromInt(1))
		data["成品总件数"] = number

		weight, ok := data["成品总金重"].(decimal.Decimal)
		if !ok {
			weight = decimal.Zero
		}
		weight = weight.Add(finished.WeightMetal)
		data["成品总金重"] = weight

		price, ok := data["成品总标签价"].(decimal.Decimal)
		if !ok {
			price = decimal.Zero
		}
		price = price.Add(finished.LabelPrice)
		data["成品总标签价"] = price

		unsalable, ok := data["成品滞销件数"].(decimal.Decimal)
		if !ok {
			unsalable = decimal.Zero
		}
		if finished.IsUnsalable(l.endtime) {
			unsalable = unsalable.Add(decimal.NewFromInt(1))
		}
		data["成品滞销件数"] = unsalable
	}

	if len(l.Olds) == 0 {
		data["旧料总件数"] = decimal.Zero
		data["旧料总金重"] = decimal.Zero
		data["旧料总抵值"] = decimal.Zero
	}

	for _, old := range l.Olds {
		number, ok := data["旧料总件数"].(decimal.Decimal)
		if !ok {
			number = decimal.Zero
		}
		number = number.Add(decimal.NewFromInt(1))
		data["旧料总件数"] = number

		weight, ok := data["旧料总金重"].(decimal.Decimal)
		if !ok {
			weight = decimal.Zero
		}
		weight = weight.Add(old.WeightMetal)
		data["旧料总金重"] = weight

		price, ok := data["旧料总抵值"].(decimal.Decimal)
		if !ok {
			price = decimal.Zero
		}
		price = price.Add(old.RecyclePrice)
		data["旧料总抵值"] = price
	}

	return data
}

// 获取成品大类数据
func (l *dataLogic) get_finished_class() map[string]any {
	data := make(map[string]any)

	for _, finished := range l.Finisheds {
		k := finished.Class.String()

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

		weight_row, ok := data["金重"].(map[string]any)
		if !ok {
			weight_row = make(map[string]any, 0)
		}
		weight, ok := weight_row[k].(decimal.Decimal)
		if !ok {
			weight = decimal.NewFromInt(0)
		}
		weight = weight.Add(finished.WeightMetal)
		weight_row[k] = weight
		data["金重"] = weight_row

		price_row, ok := data["标价"].(map[string]any)
		if !ok {
			price_row = make(map[string]any, 0)
		}
		price, ok := price_row[k].(decimal.Decimal)
		if !ok {
			price = decimal.NewFromInt(0)
		}
		price = price.Add(finished.LabelPrice)
		price_row[k] = price
		data["标价"] = price_row
	}

	return data
}

// 获取成品品类数据
func (l *dataLogic) get_finished_category() map[string]map[string]any {

	data := make(map[string]map[string]any)

	for _, class := range enums.ProductClassFinishedMap {
		data[class] = map[string]any{
			"件数": map[string]any{},
			"金重": map[string]any{},
			"标价": map[string]any{},
		}
	}

	for _, finished := range l.Finisheds {
		c := finished.Class.String()
		if c == "" {
			c = "其他"
		}
		k := finished.Category.String()
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

		weight_item, ok := data[c]["金重"].(map[string]any)
		if !ok {
			weight_item = make(map[string]any, 0)
		}
		weight, ok := weight_item[k].(decimal.Decimal)
		if !ok {
			weight = decimal.Zero
		}
		weight = weight.Add(finished.WeightMetal)
		weight_item[k] = weight
		data[c]["金重"] = weight_item

		price_item, ok := data[c]["标价"].(map[string]any)
		if !ok {
			price_item = make(map[string]any, 0)
		}
		price, ok := price_item[k].(decimal.Decimal)
		if !ok {
			price = decimal.Zero
		}
		price = price.Add(finished.LabelPrice)
		price_item[k] = price
		data[c]["标价"] = price_item
	}

	return data
}

// 获取成品库龄
func (l *dataLogic) get_finished_age() map[string]map[string]any {
	data := make(map[string]map[string]any)

	for _, class := range enums.ProductClassFinishedMap {
		data[class] = map[string]any{
			"件数": map[string]any{},
			"金重": map[string]any{},
			"标价": map[string]any{},
		}
	}

	for _, finished := range l.Finisheds {
		var k string
		// 库龄(天)，根据入库时间计算，按整天计算
		age := l.endtime.Sub(finished.EnterTime).Hours() / 24
		switch {
		case age <= 30:
			k = "30天内"
		case age > 30 && age <= 90:
			k = "31-90天"
		case age > 90 && age <= 180:
			k = "91-180天"
		case age > 180 && age <= 360:
			k = "181-360天"
		case age > 360 && age <= 720:
			k = "361-720天"
		default:
			k = "720天以上"
		}

		c := finished.Class.String()

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

		weight_item, ok := data[c]["金重"].(map[string]any)
		if !ok {
			weight_item = make(map[string]any, 0)
		}
		weight, ok := weight_item[k].(decimal.Decimal)
		if !ok {
			weight = decimal.Zero
		}
		weight = weight.Add(finished.WeightMetal)
		weight_item[k] = weight
		data[c]["金重"] = weight_item

		price_item, ok := data[c]["标价"].(map[string]any)
		if !ok {
			price_item = make(map[string]any, 0)
		}
		price, ok := price_item[k].(decimal.Decimal)
		if !ok {
			price = decimal.Zero
		}
		price = price.Add(finished.LabelPrice)
		price_item[k] = price
		data[c]["标价"] = price_item
	}

	return data
}

// 获取旧料数据
func (l *dataLogic) get_old_class() map[string]any {
	data := make(map[string]any)

	if len(l.Olds) == 0 {
		data["件数"] = map[string]any{}
		data["金重"] = map[string]any{}
		data["标价"] = map[string]any{}
	}
	for _, old := range l.Olds {
		k := old.Class.String()

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

		weight_row, ok := data["金重"].(map[string]any)
		if !ok {
			weight_row = make(map[string]any, 0)
		}
		weight, ok := weight_row[k].(decimal.Decimal)
		if !ok {
			weight = decimal.NewFromInt(0)
		}
		weight = weight.Add(old.WeightMetal)
		weight_row[k] = weight
		data["金重"] = weight_row

		price_row, ok := data["标价"].(map[string]any)
		if !ok {
			price_row = make(map[string]any, 0)
		}
		price, ok := price_row[k].(decimal.Decimal)
		if !ok {
			price = decimal.NewFromInt(0)
		}
		price = price.Add(old.LabelPrice)
		price_row[k] = price
		data["标价"] = price_row
	}

	return data
}
