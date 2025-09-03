package performance

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

	return logic.get_sales()
}

// 获取门店列表
func (r *dataLogic) get_stores() error {
	// 查询门店
	store_logic := store.StoreLogic{
		Staff: r.Staff,
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
func (r *dataLogic) get_sales() (any, error) {
	var data []map[string]any

	for _, store := range r.Stores {
		row := map[string]any{
			"name": store.Name,
		}

		total, ok := row["total"].(decimal.Decimal)
		if !ok {
			total = decimal.Decimal{}
		}

		for _, order := range r.Orders[store.Id] {
			for _, p := range order.Products {
				// 产品状态必须是已完成
				if p.Status != enums.OrderSalesStatusComplete {
					continue
				}
				// 判断产品类型
				switch p.Type {
				case enums.ProductTypeFinished:
					{
						k := p.Finished.Product.Class.String()
						if k == "" {
							k = "其他"
						}

						item, ok := row[k].(decimal.Decimal)
						if !ok {
							item = decimal.Decimal{}
						}

						item = item.Add(p.Finished.Price)

						// 合计
						total = total.Add(p.Finished.Price)

						row[k] = item
					}
				case enums.ProductTypeOld:
					{
						finished := model.ProductFinished{
							Material: p.Old.Product.Material,
							Quality:  p.Old.Product.Quality,
							Gem:      p.Old.Product.Gem,
						}
						switch p.Old.Product.RecycleMethod {
						case enums.ProductRecycleMethod_KG:
							finished.RetailType = enums.ProductRetailTypeGoldKg
						case enums.ProductRecycleMethod_PIECE:
							finished.RetailType = enums.ProductRetailTypePiece
						}

						k := finished.GetClass().String()
						if k == "" {
							k = "其他"
						}

						switch p.Old.Product.RecycleType {
						case enums.ProductRecycleTypeRecycle:
							{
								k = k + "回收旧料抵扣"
							}
						case enums.ProductRecycleTypeExchange:
							{
								k = k + "兑换旧料抵扣"
							}
						}
						item, ok := row[k].(decimal.Decimal)
						if !ok {
							item = decimal.Decimal{}
						}

						item = item.Add(p.Old.RecyclePrice)

						// 合计
						total = total.Add(p.Old.RecyclePrice)

						row[k] = item
					}
				case enums.ProductTypeAccessorie:
					{
						k := "计件配件"
						item, ok := row[k].(decimal.Decimal)
						if !ok {
							item = decimal.Decimal{}
						}

						item = item.Add(p.Accessorie.Price)

						// 合计
						total = total.Add(p.Accessorie.Price)

						row[k] = item
					}
				}
			}
		}

		row["total"] = total

		data = append(data, row)
	}

	return &data, nil
}
