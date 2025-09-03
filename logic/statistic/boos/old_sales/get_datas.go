package old_sales

import (
	"jdy/enums"
	"jdy/logic/store"
	"jdy/model"
	"jdy/types"

	"github.com/shopspring/decimal"
)

type dataLogic struct {
	*Logic

	Orders map[string][]model.OrderSales
	Stores []model.Store
}

func (l *Logic) GetDatas(req *DataReq) (any, error) {
	logic := dataLogic{
		Logic: l,
	}

	if err := logic.get_stores(); err != nil {
		return nil, err
	}

	if err := logic.get_orders(req); err != nil {
		return nil, err
	}

	// 查询数据
	switch req.Type {
	case TypesRecyclePrice:
		return logic.get_recycle_price()
	case TypesCount:
		return logic.get_count()
	case TypesWeightMetal:
		return logic.get_weight_metal()
	}

	return nil, nil
}

// 获取门店列表
func (r *dataLogic) get_stores() error {

	// 查询门店
	store_logic := store.StoreLogic{
		Staff: r.Staff,
		Ctx:   r.Ctx,
	}
	stores, err := store_logic.My(&types.StoreListMyReq{})
	if err != nil {
		return err
	}
	if stores != nil {
		r.Stores = *stores
	}

	return nil
}

// 获取订单列表
func (r *dataLogic) get_orders(req *DataReq) error {
	for _, store := range r.Stores {

		orders, ok := r.Orders[store.Id]
		if !ok {
			orders = []model.OrderSales{}
		}

		// 查询订单
		db := model.DB.Model(&model.OrderSales{})
		db = db.Where(&model.OrderSales{
			StoreId: store.Id,
		})
		db = db.Where("status in (?)", []enums.OrderSalesStatus{
			enums.OrderSalesStatusComplete,
			enums.OrderSalesStatusRefund,
		})
		db = db.Scopes(model.DurationCondition(req.Duration, "created_at", req.StartTime, req.EndTime))

		db = model.OrderSales{}.Preloads(db)

		if err := db.Find(&orders).Error; err != nil {
			return err
		}

		if r.Orders == nil {
			r.Orders = make(map[string][]model.OrderSales)
		}
		r.Orders[store.Id] = orders
	}

	return nil
}

// 销售额
func (r *dataLogic) get_recycle_price() (any, error) {
	var data []map[string]any

	for _, store := range r.Stores {
		item := map[string]any{
			"name": store.Name,
		}

		total, ok := item["total"].(decimal.Decimal)
		if !ok {
			total = decimal.Decimal{}
		}
		for _, order := range r.Orders[store.Id] {
			for _, p := range order.Products {
				if p.Type != enums.ProductTypeOld {
					continue
				}
				if p.Status != enums.OrderSalesStatusComplete {
					continue
				}

				total = total.Add(p.Old.RecyclePrice)
			}
		}
		item["total"] = total

		for k := range enums.ProductClassOldMap {
			total, ok := item[k.String()].(decimal.Decimal)
			if !ok {
				total = decimal.Decimal{}
			}
			for _, order := range r.Orders[store.Id] {
				for _, p := range order.Products {
					if p.Type != enums.ProductTypeOld {
						continue
					}
					if p.Status != enums.OrderSalesStatusComplete {
						continue
					}
					if p.Old.Product.Class != k {
						continue
					}

					total = total.Add(p.Old.RecyclePrice)
				}
			}
			item[k.String()] = total
		}

		data = append(data, item)
	}

	return &data, nil
}

// 件数
func (r *dataLogic) get_count() (any, error) {
	var data []map[string]any

	for _, store := range r.Stores {
		item := map[string]any{
			"name": store.Name,
		}

		total, ok := item["total"].(int64)
		if !ok {
			total = 0
		}
		for _, order := range r.Orders[store.Id] {
			for _, p := range order.Products {
				if p.Type != enums.ProductTypeOld {
					continue
				}
				if p.Status != enums.OrderSalesStatusComplete {
					continue
				}

				total = total + 1
			}
		}
		item["total"] = total

		for k := range enums.ProductClassOldMap {
			total, ok := item[k.String()].(int64)
			if !ok {
				total = 0
			}
			for _, order := range r.Orders[store.Id] {
				for _, p := range order.Products {
					if p.Type != enums.ProductTypeOld {
						continue
					}
					if p.Status != enums.OrderSalesStatusComplete {
						continue
					}
					if p.Old.Product.Class != k {
						continue
					}

					total = total + 1
				}
			}
			item[k.String()] = total
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

		total, ok := item["total"].(decimal.Decimal)
		if !ok {
			total = decimal.Decimal{}
		}
		for _, order := range r.Orders[store.Id] {
			for _, p := range order.Products {
				if p.Type != enums.ProductTypeOld {
					continue
				}
				if p.Status != enums.OrderSalesStatusComplete {
					continue
				}

				total = total.Add(p.Old.Product.WeightMetal)
			}
		}
		item["total"] = total

		for k := range enums.ProductClassOldMap {
			total, ok := item[k.String()].(decimal.Decimal)
			if !ok {
				total = decimal.Decimal{}
			}
			for _, order := range r.Orders[store.Id] {
				for _, p := range order.Products {
					if p.Type != enums.ProductTypeOld {
						continue
					}
					if p.Status != enums.OrderSalesStatusComplete {
						continue
					}
					if p.Old.Product.Class != k {
						continue
					}

					total = total.Add(p.Old.Product.WeightMetal)
				}
			}
			item[k.String()] = total
		}

		data = append(data, item)
	}

	return &data, nil
}
