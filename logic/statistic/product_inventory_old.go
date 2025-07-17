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

type ProductInventoryOldTitle struct {
	Title     string `json:"title"`
	Key       string `json:"key"`
	Width     string `json:"width"`
	Fixed     string `json:"fixed"`
	ClassName string `json:"className"`
	Align     string `json:"align"`
}

func (l *StatisticLogic) ProductInventoryOldTitles() *[]ProductInventoryOldTitle {
	var titles []ProductInventoryOldTitle
	titles = append(titles, ProductInventoryOldTitle{
		Title:     "门店",
		Key:       "name",
		Width:     "100px",
		Fixed:     "left",
		ClassName: "age",
		Align:     "center",
	})
	titles = append(titles, ProductInventoryOldTitle{
		Title: "总",
		Key:   "total",
		Width: "100px",
		Fixed: "left",
		Align: "center",
	})

	for k, v := range enums.ProductClassOldMap {
		titles = append(titles, ProductInventoryOldTitle{
			Title: v,
			Key:   fmt.Sprint(k),
			Width: "100px",
			Align: "center",
		})
	}

	return &titles
}

type ProductInventoryOldType int

const (
	ProductInventoryOldTypeCount        ProductInventoryOldType = iota + 1 // 件数
	ProductInventoryOldTypeWeightMetal                                     // 金重
	ProductInventoryOldTypeRecyclePrice                                    // 抵值
)

type ProductInventoryOldReq struct {
	Type ProductInventoryOldType `json:"type" label:"类型" find:"true" required:"true" sort:"1" type:"number" input:"radio" preset:"typeMap"` // 类型
}

type ProductInventoryOldLogic struct {
	*StatisticLogic

	Stores *[]model.Store
}

func (l *StatisticLogic) ProductInventoryOldData(req *ProductInventoryOldReq) (any, error) {
	logic := ProductInventoryOldLogic{
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
	case ProductInventoryOldTypeCount:
		return logic.ProductInventoryOldCountData(req)
	case ProductInventoryOldTypeWeightMetal:
		return logic.ProductInventoryOldWeightMetalData(req)
	case ProductInventoryOldTypeRecyclePrice:
		return logic.ProductInventoryOldRecyclePriceData(req)
	}

	return nil, nil
}

// 件数
func (r *ProductInventoryOldLogic) ProductInventoryOldCountData(req *ProductInventoryOldReq) (any, error) {
	var data []map[string]any

	for _, store := range *r.Stores {
		item := map[string]any{
			"name": store.Name,
		}

		// 总
		db_total := model.DB.Model(&model.ProductOld{})
		db_total = db_total.Where(&model.ProductOld{
			StoreId: store.Id,
			Status:  enums.ProductStatusNormal,
		})
		var total int64
		if err := db_total.Count(&total).Error; err != nil {
			return nil, err
		}
		item["total"] = total

		for k := range enums.ProductClassOldMap {
			db := model.DB.Model(&model.ProductOld{})
			db = db.Where(&model.ProductOld{
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
func (r *ProductInventoryOldLogic) ProductInventoryOldWeightMetalData(req *ProductInventoryOldReq) (any, error) {
	var data []map[string]any

	for _, store := range *r.Stores {
		item := map[string]any{
			"name": store.Name,
		}

		// 总
		db_total := model.DB.Model(&model.ProductOld{})
		db_total = db_total.Where(&model.ProductOld{
			StoreId: store.Id,
			Status:  enums.ProductStatusNormal,
		})
		var total decimal.Decimal
		if err := db_total.Select("SUM(weight_metal) as total").Having("total > 0").Scan(&total).Error; err != nil {
			return nil, err
		}
		item["total"] = total

		for k := range enums.ProductClassOldMap {
			db := model.DB.Model(&model.ProductOld{})
			db = db.Where(&model.ProductOld{
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

func (r *ProductInventoryOldLogic) ProductInventoryOldRecyclePriceData(req *ProductInventoryOldReq) (any, error) {
	var data []map[string]any

	for _, store := range *r.Stores {
		item := map[string]any{
			"name": store.Name,
		}

		// 总
		db_total := model.DB.Model(&model.ProductOld{})
		db_total = db_total.Where(&model.ProductOld{
			StoreId: store.Id,
			Status:  enums.ProductStatusNormal,
		})
		var total decimal.Decimal
		if err := db_total.Select("SUM(recycle_price) as total").Having("total > 0").Scan(&total).Error; err != nil {
			return nil, err
		}
		item["total"] = total

		for k := range enums.ProductClassOldMap {
			db := model.DB.Model(&model.ProductOld{})
			db = db.Where(&model.ProductOld{
				StoreId: store.Id,
				Status:  enums.ProductStatusNormal,
				Class:   k,
			})
			var total decimal.Decimal
			if err := db.Select("SUM(recycle_price) as total").Having("total > 0").Scan(&total).Error; err != nil {
				return nil, err
			}
			item[fmt.Sprint(k)] = total
		}

		data = append(data, item)
	}

	return &data, nil
}

var ProductInventoryOldTypeMap = map[ProductInventoryOldType]string{
	ProductInventoryOldTypeCount:        "件数",
	ProductInventoryOldTypeWeightMetal:  "金重",
	ProductInventoryOldTypeRecyclePrice: "抵值",
}

func (p ProductInventoryOldType) ToMap() any {
	return ProductInventoryOldTypeMap
}

func (p ProductInventoryOldType) InMap() error {
	if _, ok := ProductInventoryOldTypeMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}
