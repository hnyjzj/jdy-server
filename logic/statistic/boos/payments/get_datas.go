package payments

import (
	"fmt"
	"jdy/enums"
	"jdy/logic/store"
	"jdy/model"
	"jdy/types"

	"github.com/shopspring/decimal"
)

type dataLogic struct {
	*Logic

	Stores *[]model.Store
}

func (l *Logic) GetDatas(req *DataReq) (any, error) {
	logic := dataLogic{
		Logic: l,
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
	case TypesIncome:
		return logic.get_income(req)
	case TypesExpense:
		return logic.get_expense(req)
	case TypesSurplus:
		return logic.get_surplus(req)
	}

	return nil, nil
}

// 收入
func (r *dataLogic) get_income(req *DataReq) (any, error) {
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
		db_total = db_total.Scopes(model.DurationCondition(req.Duration, "created_at", req.StartTime, req.EndTime))
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
			db = db.Scopes(model.DurationCondition(req.Duration, "created_at", req.StartTime, req.EndTime))
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
func (r *dataLogic) get_expense(req *DataReq) (any, error) {
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
		db_total = db_total.Scopes(model.DurationCondition(req.Duration, "created_at", req.StartTime, req.EndTime))
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
			db = db.Scopes(model.DurationCondition(req.Duration, "created_at", req.StartTime, req.EndTime))
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
func (r *dataLogic) get_surplus(req *DataReq) (any, error) {
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
		total_income_db = total_income_db.Scopes(model.DurationCondition(req.Duration, "created_at", req.StartTime, req.EndTime))
		var total_income decimal.Decimal
		if err := total_income_db.Select("SUM(amount) as total").Having("total <> 0").Scan(&total_income).Error; err != nil {
			return nil, err
		}

		total_expense_db := model.DB.Model(&model.OrderPayment{})
		total_expense_db = total_expense_db.Where(&model.OrderPayment{
			StoreId: store.Id,
			Type:    enums.FinanceTypeExpense,
		})
		total_expense_db = total_expense_db.Scopes(model.DurationCondition(req.Duration, "created_at", req.StartTime, req.EndTime))
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
			income_db = income_db.Scopes(model.DurationCondition(req.Duration, "created_at", req.StartTime, req.EndTime))
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
			expense_db = expense_db.Scopes(model.DurationCondition(req.Duration, "created_at", req.StartTime, req.EndTime))
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
