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

	Sales   map[string][]model.OrderSales
	Refunds map[string][]model.OrderRefund
	Stores  []model.Store
}

func (l *Logic) GetDatas(req *DataReq) (any, error) {
	logic := dataLogic{
		Logic: l,
	}

	if err := logic.get_stores(); err != nil {
		return nil, err
	}

	if err := logic.get_sales(req); err != nil {
		return nil, err
	}

	if err := logic.get_refunds(req); err != nil {
		return nil, err
	}

	return logic.get_data()
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
func (r *dataLogic) get_sales(req *DataReq) error {
	for _, store := range r.Stores {
		sales, ok := r.Sales[store.Id]
		if !ok {
			sales = []model.OrderSales{}
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

		if err := db.Find(&sales).Error; err != nil {
			return err
		}

		if r.Sales == nil {
			r.Sales = make(map[string][]model.OrderSales)
		}
		r.Sales[store.Id] = sales
	}

	return nil
}

// 获取退款列表
func (r *dataLogic) get_refunds(req *DataReq) error {
	for _, store := range r.Stores {
		refunds, ok := r.Refunds[store.Id]
		if !ok {
			refunds = []model.OrderRefund{}
		}

		// 查询退款
		db := model.DB.Model(&model.OrderRefund{})
		db = db.Where(&model.OrderRefund{
			StoreId: store.Id,
		})
		db = db.Where("order_type in (?)", []enums.OrderType{
			enums.OrderTypeSales,
		})
		db = db.Scopes(model.DurationCondition(req.Duration, "created_at", req.StartTime, req.EndTime))

		if err := db.Find(&refunds).Error; err != nil {
			return err
		}

		if r.Refunds == nil {
			r.Refunds = make(map[string][]model.OrderRefund)
		}
		r.Refunds[store.Id] = refunds
	}
	return nil
}

// 销售额
func (r *dataLogic) get_data() (any, error) {
	var data []map[string]any

	for _, store := range r.Stores {
		row := map[string]any{
			"name": store.Name,
		}

		total, ok := row["total"].(decimal.Decimal)
		if !ok {
			total = decimal.Decimal{}
		}

		for _, sale := range r.Sales[store.Id] {
			for _, product := range sale.Products {
				// 产品状态必须是已完成
				if product.Status != enums.OrderSalesStatusComplete {
					continue
				}
				// 判断产品类型
				switch product.Type {
				case enums.ProductTypeFinished:
					{
						k := product.Finished.Product.Class.String()
						if k == "" {
							k = "其他"
						}

						item, ok := row[k].(decimal.Decimal)
						if !ok {
							item = decimal.Decimal{}
						}

						item = item.Add(product.Finished.Price)
						total = total.Add(product.Finished.Price)

						for _, refund := range r.Refunds[store.Id] {
							if refund.Type != enums.ProductTypeFinished {
								continue
							}
							if refund.OrderId != sale.Id {
								continue
							}
							if refund.Code != product.Code {
								continue
							}
							item = item.Sub(refund.Price)
							total = total.Sub(refund.Price)
						}

						row[k] = item
					}
				case enums.ProductTypeOld:
					{
						finished := model.ProductFinished{
							Material: product.Old.Product.Material,
							Quality:  product.Old.Product.Quality,
							Gem:      product.Old.Product.Gem,
						}
						switch product.Old.Product.RecycleMethod {
						case enums.ProductRecycleMethod_KG:
							finished.RetailType = enums.ProductRetailTypeGoldKg
						case enums.ProductRecycleMethod_PIECE:
							finished.RetailType = enums.ProductRetailTypePiece
						}

						k := finished.GetClass().String()
						if k == "" {
							k = "其他"
						}

						switch product.Old.Product.RecycleType {
						case enums.ProductRecycleTypeRecycle:
							{
								k = k + "回收"
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

						item = item.Add(product.Old.RecyclePrice.Neg())

						for _, refund := range r.Refunds[store.Id] {
							if refund.Type != enums.ProductTypeOld {
								continue
							}
							if refund.OrderId != sale.Id {
								continue
							}
							if refund.Code != product.Code {
								continue
							}
							item = item.Sub(refund.Price.Neg())
						}

						switch product.Old.Product.RecycleType {
						case enums.ProductRecycleTypeExchange:
							{
								total = total.Add(product.Old.RecyclePrice.Neg())

								for _, refund := range r.Refunds[store.Id] {
									if refund.Type != enums.ProductTypeOld {
										continue
									}
									if refund.OrderId != sale.Id {
										continue
									}
									if refund.Code != product.Code {
										continue
									}
									total = total.Add(refund.Price.Neg())
								}
							}
						}

						row[k] = item
					}
				case enums.ProductTypeAccessorie:
					{
						k := "计件配件"
						item, ok := row[k].(decimal.Decimal)
						if !ok {
							item = decimal.Decimal{}
						}

						item = item.Add(product.Accessorie.Price)
						total = total.Add(product.Accessorie.Price)

						for _, refund := range r.Refunds[store.Id] {
							if refund.Type != enums.ProductTypeAccessorie {
								continue
							}
							if refund.OrderId != sale.Id {
								continue
							}
							if refund.Name != product.Name {
								continue
							}
							item = item.Sub(refund.Price)
							total = total.Sub(refund.Price)
						}

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
