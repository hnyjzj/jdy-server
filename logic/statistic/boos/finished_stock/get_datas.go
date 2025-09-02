package finished_stock

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

	Stores []model.Store
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
	logic.Stores = *stores

	// 查询数据
	switch req.Type {
	case TypesCount:
		return logic.get_count(req)
	case TypesWeightMetal:
		return logic.get_weight_metal(req)
	case TypesLabelPrice:
		return logic.get_label_price(req)
	}

	return nil, nil
}

// 件数
func (r *dataLogic) get_count(req *DataReq) (any, error) {
	var data []map[string]any

	for _, store := range r.Stores {
		item := map[string]any{
			"name": store.Name,
		}

		// 合计
		db_total := model.DB.Model(&model.ProductFinished{})
		db_total = db_total.Where(&model.ProductFinished{
			StoreId: store.Id,
			Status:  enums.ProductStatusNormal,
		})
		db_total = db_total.Scopes(model.DurationCondition(req.Duration, "created_at", req.StartTime, req.EndTime))

		var total int64
		if err := db_total.Count(&total).Error; err != nil {
			return nil, err
		}
		item["total"] = total

		for k := range enums.ProductClassFinishedMap {
			db := model.DB.Model(&model.ProductFinished{})
			db = db.Where(&model.ProductFinished{
				StoreId: store.Id,
				Status:  enums.ProductStatusNormal,
				Class:   k,
			})
			db = db.Scopes(model.DurationCondition(req.Duration, "created_at", req.StartTime, req.EndTime))

			var count int64
			if err := db.Count(&count).Error; err != nil {
				return nil, err
			}
			item[fmt.Sprint(k)] = count
		}

		data = append(data, item)
	}

	return &data, nil
}

// 金重
func (r *dataLogic) get_weight_metal(req *DataReq) (any, error) {
	var data []map[string]any

	for _, store := range r.Stores {
		item := map[string]any{
			"name": store.Name,
		}

		// 合计
		db_total := model.DB.Model(&model.ProductFinished{})
		db_total = db_total.Where(&model.ProductFinished{
			StoreId: store.Id,
			Status:  enums.ProductStatusNormal,
		})
		db_total = db_total.Scopes(model.DurationCondition(req.Duration, "created_at", req.StartTime, req.EndTime))

		var total decimal.Decimal
		if err := db_total.Select("SUM(weight_metal) as total").Having("total > 0").Scan(&total).Error; err != nil {
			return nil, err
		}
		item["total"] = total

		for k := range enums.ProductClassFinishedMap {
			db := model.DB.Model(&model.ProductFinished{})
			db = db.Where(&model.ProductFinished{
				StoreId: store.Id,
				Status:  enums.ProductStatusNormal,
				Class:   k,
			})
			db = db.Scopes(model.DurationCondition(req.Duration, "created_at", req.StartTime, req.EndTime))

			var total decimal.Decimal
			if err := db.Select("SUM(weight_metal) as total").Having("total > 0").Scan(&total).Error; err != nil {
				return nil, err
			}
			item[fmt.Sprint(k)] = total
		}

		data = append(data, item)
	}

	return &data, nil
}

func (r *dataLogic) get_label_price(req *DataReq) (any, error) {
	var data []map[string]any

	for _, store := range r.Stores {
		item := map[string]any{
			"name": store.Name,
		}

		// 合计
		db_total := model.DB.Model(&model.ProductFinished{})
		db_total = db_total.Where(&model.ProductFinished{
			StoreId: store.Id,
			Status:  enums.ProductStatusNormal,
		})
		db_total = db_total.Scopes(model.DurationCondition(req.Duration, "created_at", req.StartTime, req.EndTime))

		var total decimal.Decimal
		if err := db_total.Select("SUM(label_price) as total").Having("total > 0").Scan(&total).Error; err != nil {
			return nil, err
		}
		item["total"] = total

		for k := range enums.ProductClassFinishedMap {
			db := model.DB.Model(&model.ProductFinished{})
			db = db.Where(&model.ProductFinished{
				StoreId: store.Id,
				Status:  enums.ProductStatusNormal,
				Class:   k,
			})
			db = db.Scopes(model.DurationCondition(req.Duration, "created_at", req.StartTime, req.EndTime))

			var total decimal.Decimal
			if err := db.Select("SUM(label_price) as total").Having("total > 0").Scan(&total).Error; err != nil {
				return nil, err
			}
			item[fmt.Sprint(k)] = total
		}

		data = append(data, item)
	}

	return &data, nil
}
