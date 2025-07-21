package statistic

import (
	"errors"
	"fmt"
	"jdy/enums"
	"jdy/logic/store"
	"jdy/model"
	"jdy/types"

	"github.com/shopspring/decimal"
)

type OrderPaymentTitle struct {
	Title     string `json:"title"`
	Key       string `json:"key"`
	Width     string `json:"width"`
	Fixed     string `json:"fixed"`
	ClassName string `json:"className"`
	Align     string `json:"align"`
}

func (l *StatisticLogic) OrderPaymentTitles() *[]OrderPaymentTitle {
	var titles []OrderPaymentTitle
	titles = append(titles, OrderPaymentTitle{
		Title:     "门店",
		Key:       "name",
		Width:     "80px",
		Fixed:     "left",
		ClassName: "age",
		Align:     "center",
	})
	titles = append(titles, OrderPaymentTitle{
		Title: "总",
		Key:   "total",
		Width: "80px",
		Fixed: "left",
		Align: "center",
	})

	for k, v := range enums.OrderPaymentMethodMap {
		titles = append(titles, OrderPaymentTitle{
			Title: v,
			Key:   fmt.Sprint(k),
			Width: "80px",
			Align: "center",
		})
	}

	return &titles
}

type OrderPaymentType int

const (
	OrderPaymentTypeIncome  OrderPaymentType = iota + 1 // 收入
	OrderPaymentTypeExpense                             // 支出
	OrderPaymentTypeSurplus                             // 结余
)

type OrderPaymentReq struct {
	Type      OrderPaymentType `json:"type" label:"类型" find:"true" required:"true" sort:"1" type:"number" input:"radio" preset:"typeMap"`           // 类型
	Duration  enums.Duration   `json:"duration" label:"时间范围" find:"true" required:"true" sort:"2" type:"number" input:"radio" preset:"durationMap"` // 时间范围
	StartTime string           `json:"startTime" label:"开始时间" find:"true" required:"false" sort:"3" type:"string" input:"date"`                     // 开始时间
	EndTime   string           `json:"endTime" label:"结束时间" find:"true" required:"false" sort:"4" type:"string" input:"date"`                       // 结束时间
}

type OrderPaymentLogic struct {
	*StatisticLogic

	Stores *[]model.Store
}

func (l *StatisticLogic) OrderPaymentData(req *OrderPaymentReq) (any, error) {
	logic := OrderPaymentLogic{
		StatisticLogic: l,
	}

	// 查询门店
	store_logic := store.StoreLogic{
		Staff: l.Staff,
	}
	stores, err := store_logic.My(&types.StoreListMyReq{})
	if err != nil {
		return nil, err
	}
	logic.Stores = stores

	// 查询数据
	switch req.Type {
	case OrderPaymentTypeIncome:
		return logic.OrderPaymentTypeIncomeData(req)
	case OrderPaymentTypeExpense:
		return logic.OrderPaymentTypeExpenseData(req)
	case OrderPaymentTypeSurplus:
		return logic.OrderPaymentTypeSurplusData(req)
	}

	return nil, nil
}

// 收入
func (r *OrderPaymentLogic) OrderPaymentTypeIncomeData(req *OrderPaymentReq) (any, error) {
	var data []map[string]any

	for _, store := range *r.Stores {
		item := map[string]any{
			"name": store.Name,
		}

		// 总计
		db_total := model.DB.Model(&model.OrderPayment{})
		db_total = db_total.Where(&model.OrderPayment{
			StoreId: store.Id,
			Type:    enums.FinanceTypeIncome,
		})
		db_total = db_total.Scopes(model.DurationCondition(req.Duration, req.StartTime, req.EndTime))
		var total decimal.Decimal
		if err := db_total.Select("SUM(amount) as total").Having("total <> 0").Scan(&total).Error; err != nil {
			return nil, err
		}
		item["total"] = total

		// 按支付方式
		for k := range enums.OrderPaymentMethodMap {
			db := model.DB.Model(&model.OrderPayment{})
			db = db.Where(&model.OrderPayment{
				StoreId:       store.Id,
				Type:          enums.FinanceTypeIncome,
				PaymentMethod: k,
			})
			db = db.Scopes(model.DurationCondition(req.Duration, req.StartTime, req.EndTime))
			var total decimal.Decimal
			if err := db.Select("SUM(amount) as total").Having("total <> 0").Scan(&total).Error; err != nil {
				return nil, err
			}
			item[fmt.Sprint(k)] = total
		}

		data = append(data, item)
	}

	return &data, nil
}

// 支出
func (r *OrderPaymentLogic) OrderPaymentTypeExpenseData(req *OrderPaymentReq) (any, error) {
	var data []map[string]any

	for _, store := range *r.Stores {
		item := map[string]any{
			"name": store.Name,
		}

		// 总计
		db_total := model.DB.Model(&model.OrderPayment{})
		db_total = db_total.Where(&model.OrderPayment{
			StoreId: store.Id,
			Type:    enums.FinanceTypeExpense,
		})
		db_total = db_total.Scopes(model.DurationCondition(req.Duration, req.StartTime, req.EndTime))
		var total decimal.Decimal
		if err := db_total.Select("SUM(amount) as total").Having("total <> 0").Scan(&total).Error; err != nil {
			return nil, err
		}
		item["total"] = total

		// 按支付方式
		for k := range enums.OrderPaymentMethodMap {
			db := model.DB.Model(&model.OrderPayment{})
			db = db.Where(&model.OrderPayment{
				StoreId:       store.Id,
				Type:          enums.FinanceTypeExpense,
				PaymentMethod: k,
			})
			db = db.Scopes(model.DurationCondition(req.Duration, req.StartTime, req.EndTime))
			var total decimal.Decimal
			if err := db.Select("SUM(amount) as total").Having("total <> 0").Scan(&total).Error; err != nil {
				return nil, err
			}
			item[fmt.Sprint(k)] = total
		}

		data = append(data, item)
	}

	return &data, nil
}

// 结余
func (r *OrderPaymentLogic) OrderPaymentTypeSurplusData(req *OrderPaymentReq) (any, error) {
	var data []map[string]any

	for _, store := range *r.Stores {
		item := map[string]any{
			"name": store.Name,
		}

		// 总计
		total_income_db := model.DB.Model(&model.OrderPayment{})
		total_income_db = total_income_db.Where(&model.OrderPayment{
			StoreId: store.Id,
			Type:    enums.FinanceTypeIncome,
		})
		total_income_db = total_income_db.Scopes(model.DurationCondition(req.Duration, req.StartTime, req.EndTime))
		var total_income decimal.Decimal
		if err := total_income_db.Select("SUM(amount) as total").Having("total <> 0").Scan(&total_income).Error; err != nil {
			return nil, err
		}

		total_expense_db := model.DB.Model(&model.OrderPayment{})
		total_expense_db = total_expense_db.Where(&model.OrderPayment{
			StoreId: store.Id,
			Type:    enums.FinanceTypeExpense,
		})
		total_expense_db = total_expense_db.Scopes(model.DurationCondition(req.Duration, req.StartTime, req.EndTime))
		var total_expense decimal.Decimal
		if err := total_expense_db.Select("SUM(amount) as total").Having("total <> 0").Scan(&total_expense).Error; err != nil {
			return nil, err
		}

		item["total"] = total_income.Sub(total_expense)

		for k := range enums.OrderPaymentMethodMap {
			income_db := model.DB.Model(&model.OrderPayment{})
			income_db = income_db.Where(&model.OrderPayment{
				StoreId:       store.Id,
				Type:          enums.FinanceTypeIncome,
				PaymentMethod: k,
			})
			income_db = income_db.Scopes(model.DurationCondition(req.Duration, req.StartTime, req.EndTime))
			var income decimal.Decimal
			if err := income_db.Select("SUM(amount) as total").Having("total <> 0").Scan(&income).Error; err != nil {
				return nil, err
			}

			expense_db := model.DB.Model(&model.OrderPayment{})
			expense_db = expense_db.Where(&model.OrderPayment{
				StoreId:       store.Id,
				Type:          enums.FinanceTypeExpense,
				PaymentMethod: k,
			})
			expense_db = expense_db.Scopes(model.DurationCondition(req.Duration, req.StartTime, req.EndTime))
			var expense decimal.Decimal
			if err := expense_db.Select("SUM(amount) as total").Having("total <> 0").Scan(&expense).Error; err != nil {
				return nil, err
			}

			item[fmt.Sprint(k)] = income.Sub(expense)
		}

		data = append(data, item)
	}

	return &data, nil
}

var OrderPaymentTypeMap = map[OrderPaymentType]string{
	OrderPaymentTypeIncome:  "收入",
	OrderPaymentTypeExpense: "支出",
	OrderPaymentTypeSurplus: "结余",
}

func (p OrderPaymentType) ToMap() any {
	return OrderPaymentTypeMap
}

func (p OrderPaymentType) InMap() error {
	if _, ok := OrderPaymentTypeMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}
