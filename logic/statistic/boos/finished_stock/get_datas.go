package finished_stock

import (
	"fmt"
	"jdy/enums"
	"jdy/logic/store"
	"jdy/model"
	"jdy/types"
	"time"

	"github.com/shopspring/decimal"
)

type dataLogic struct {
	*Logic

	endtime time.Time

	Stores []model.Store
}

func (l *Logic) GetDatas(req *DataReq) (any, error) {
	logic := dataLogic{
		Logic: l,
	}

	// 查询门店
	store_logic := store.StoreLogic{
		Staff: l.Staff,
		Ctx:   l.Ctx,
	}
	stores, err := store_logic.My(&types.StoreListMyReq{})
	if err != nil {
		return nil, err
	}
	if stores != nil {
		logic.Stores = *stores
	}

	_, endtime, err := req.Duration.GetTime(time.Now(), req.StartTime, req.EndTime)
	if err != nil {
		return nil, err
	}
	logic.endtime = endtime

	// 查询数据
	switch req.Type {
	case TypesCount:
		return logic.get_count()
	case TypesWeightMetal:
		return logic.get_weight_metal()
	case TypesLabelPrice:
		return logic.get_label_price()
	}

	return nil, nil
}

// 件数
func (r *dataLogic) get_count() (any, error) {
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
		db_total = db_total.Where("enter_time <= ?", r.endtime)

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
			db = db.Where("enter_time <= ?", r.endtime)

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
func (r *dataLogic) get_weight_metal() (any, error) {
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
		db_total = db_total.Where("enter_time <= ?", r.endtime)

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
			db = db.Where("enter_time <= ?", r.endtime)

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

func (r *dataLogic) get_label_price() (any, error) {
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
		db_total = db_total.Where("enter_time <= ?", r.endtime)

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
			db = db.Where("enter_time <= ?", r.endtime)

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
