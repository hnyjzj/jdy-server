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

type ProductInventoryFinishedTitle struct {
	Title     string `json:"title"`
	Key       string `json:"key"`
	Width     string `json:"width"`
	Fixed     string `json:"fixed"`
	ClassName string `json:"className"`
	Align     string `json:"align"`
}

func (l *StatisticLogic) ProductInventoryFinishedTitles() *[]ProductInventoryFinishedTitle {
	var titles []ProductInventoryFinishedTitle
	titles = append(titles, ProductInventoryFinishedTitle{
		Title:     "门店",
		Key:       "name",
		Width:     "100px",
		Fixed:     "left",
		ClassName: "age",
		Align:     "center",
	})
	titles = append(titles, ProductInventoryFinishedTitle{
		Title: "总",
		Key:   "total",
		Width: "100px",
		Fixed: "left",
		Align: "center",
	})

	for k, v := range enums.ProductClassFinishedMap {
		titles = append(titles, ProductInventoryFinishedTitle{
			Title: v,
			Key:   fmt.Sprint(k),
			Width: "100px",
			Align: "center",
		})
	}

	return &titles
}

type ProductInventoryFinishedType int

const (
	ProductInventoryFinishedTypeCount       ProductInventoryFinishedType = iota + 1 // 件数
	ProductInventoryFinishedTypeWeightMetal                                         // 金重
	ProductInventoryFinishedTypeLabelPrice                                          // 标价
	// ProductInventoryFinishedTypeCost                                                // 成本
)

type ProductInventoryFinishedReq struct {
	Type ProductInventoryFinishedType `json:"type" label:"类型" find:"true" required:"true" sort:"1" type:"number" input:"radio" preset:"typeMap"` // 类型
}

type ProductInventoryFinishedLogic struct {
	*StatisticLogic

	Stores *[]model.Store
}

func (l *StatisticLogic) ProductInventoryFinishedData(req *ProductInventoryFinishedReq) (any, error) {
	logic := ProductInventoryFinishedLogic{
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
	case ProductInventoryFinishedTypeCount:
		return logic.ProductInventoryFinishedCountData(req)
	case ProductInventoryFinishedTypeWeightMetal:
		return logic.ProductInventoryFinishedWeightMetalData(req)
	case ProductInventoryFinishedTypeLabelPrice:
		return logic.ProductInventoryFinishedLabelPriceData(req)
		// case ProductInventoryFinishedTypeCost:
		// 	return logic.ProductInventoryFinishedCostData(req)
	}

	return nil, nil
}

// 件数
func (r *ProductInventoryFinishedLogic) ProductInventoryFinishedCountData(req *ProductInventoryFinishedReq) (any, error) {
	var data []map[string]any

	for _, store := range *r.Stores {
		item := map[string]any{
			"name": store.Name,
		}

		// 总
		db_total := model.DB.Model(&model.ProductFinished{})
		db_total = db_total.Where(&model.ProductFinished{
			StoreId: store.Id,
			Status:  enums.ProductStatusNormal,
		})
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
func (r *ProductInventoryFinishedLogic) ProductInventoryFinishedWeightMetalData(req *ProductInventoryFinishedReq) (any, error) {
	var data []map[string]any

	for _, store := range *r.Stores {
		item := map[string]any{
			"name": store.Name,
		}

		// 总
		db_total := model.DB.Model(&model.ProductFinished{})
		db_total = db_total.Where(&model.ProductFinished{
			StoreId: store.Id,
			Status:  enums.ProductStatusNormal,
		})
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

func (r *ProductInventoryFinishedLogic) ProductInventoryFinishedLabelPriceData(req *ProductInventoryFinishedReq) (any, error) {
	var data []map[string]any

	for _, store := range *r.Stores {
		item := map[string]any{
			"name": store.Name,
		}

		// 总
		db_total := model.DB.Model(&model.ProductFinished{})
		db_total = db_total.Where(&model.ProductFinished{
			StoreId: store.Id,
			Status:  enums.ProductStatusNormal,
		})
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

func (r *ProductInventoryFinishedLogic) ProductInventoryFinishedCostData(req *ProductInventoryFinishedReq) (any, error) {
	var data []map[string]any

	for _, store := range *r.Stores {
		item := map[string]any{
			"name": store.Name,
		}

		// 总
		db_total := model.DB.Model(&model.ProductFinished{})
		db_total = db_total.Where(&model.ProductFinished{
			StoreId: store.Id,
			Status:  enums.ProductStatusNormal,
		})
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

var ProductInventoryFinishedTypeMap = map[ProductInventoryFinishedType]string{
	ProductInventoryFinishedTypeCount:       "件数",
	ProductInventoryFinishedTypeWeightMetal: "金重",
	ProductInventoryFinishedTypeLabelPrice:  "标价",
	// ProductInventoryFinishedTypeCost:        "成本",
}

func (p ProductInventoryFinishedType) ToMap() any {
	return ProductInventoryFinishedTypeMap
}

func (p ProductInventoryFinishedType) InMap() error {
	if _, ok := ProductInventoryFinishedTypeMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}
