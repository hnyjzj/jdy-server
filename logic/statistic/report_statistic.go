package statistic

import (
	"context"
	"jdy/enums"
	"jdy/message"
	"jdy/model"
	"log"
	"sort"
	"time"

	"github.com/shopspring/decimal"
)

func (con StatisticLogic) SendReportStatistic() {
	// 获取发送人
	var staffs []model.Staff
	db := model.DB.Model(&model.Staff{})
	db = db.Where("identity IN (?)", []enums.Identity{
		enums.IdentitySuperAdmin,
		enums.IdentityAreaManager,
	})
	db = model.Staff{}.Preloads(db)
	if err := db.Find(&staffs).Error; err != nil {
		log.Println(err)
		return
	}

	// 获取所有门店
	var allStore []model.Store
	sdb := model.DB.Model(&model.Store{})
	sdb = sdb.Where("name NOT LIKE (?)", "%"+model.HeaderquartersPrefix+"%")
	sdb = sdb.Order("name desc")
	if err := sdb.Find(&allStore).Error; err != nil {
		log.Printf("SendReportStatistic error: %v", err.Error())
		return
	}

	allData := make(map[string]message.ReportStatisticMessage)
	for _, store := range allStore {
		// 今日销售
		today_sales, err := get_sales(store.Id, enums.DurationToday)
		if err != nil {
			log.Printf("SendReportStatistic error: %v", err.Error())
			continue
		}
		// 今日退货
		today_refunds, err := get_refunds(store.Id, enums.DurationToday)
		if err != nil {
			log.Printf("SendReportStatistic error: %v", err.Error())
			continue
		}
		// 本月销售
		month_sales, err := get_sales(store.Id, enums.DurationMonth)
		if err != nil {
			log.Printf("SendReportStatistic error: %v", err.Error())
			continue
		}
		// 本月退货
		month_refunds, err := get_refunds(store.Id, enums.DurationMonth)
		if err != nil {
			log.Printf("SendReportStatistic error: %v", err.Error())
			continue
		}

		// 发送
		req := message.ReportStatisticMessage{
			StoreName:       store.Name,
			StatisticalTime: time.Now(),
			TodayFinisheds:  make(map[string]decimal.Decimal),
		}

		for _, sale := range today_sales {
			for _, product := range sale.Products {
				switch product.Type {
				case enums.ProductTypeFinished:
					{
						price := product.Finished.Price

						for _, refund := range today_refunds {
							if refund.Type != enums.ProductTypeFinished {
								continue
							}
							if refund.OrderId != sale.Id {
								continue
							}
							if refund.Code != product.Code {
								continue
							}

							price = price.Sub(refund.Price)
						}

						if price.Equal(decimal.Zero) {
							continue
						}

						req.TodayFinished = req.TodayFinished.Add(price)
						class := product.Finished.Product.Class.String()
						req.TodayFinisheds[class] = req.TodayFinisheds[class].Add(price)
					}
				case enums.ProductTypeOld:
					{
						var price decimal.Decimal
						switch product.Old.Product.RecycleType {
						case enums.ProductRecycleTypeExchange:
							{
								price = product.Old.RecyclePrice

								for _, refund := range today_refunds {
									if refund.Type != enums.ProductTypeOld {
										continue
									}
									if refund.OrderId != sale.Id {
										continue
									}
									if refund.Code != product.Code {
										continue
									}
									price = price.Sub(refund.Price)
								}
							}
						}
						req.TodayOld = req.TodayOld.Add(price)
					}
				case enums.ProductTypeAccessorie:
					{
						price := product.Accessorie.Price

						for _, refund := range today_refunds {
							if refund.Type != enums.ProductTypeAccessorie {
								continue
							}
							if refund.OrderId != sale.Id {
								continue
							}
							if refund.Name != product.Name {
								continue
							}

							price = price.Sub(refund.Price)
						}

						req.TodayAcciessorie = req.TodayAcciessorie.Add(price)
					}
				}
			}
		}

		for _, sale := range month_sales {
			for _, product := range sale.Products {
				switch product.Type {
				case enums.ProductTypeFinished:
					{
						price := product.Finished.Price

						for _, refund := range month_refunds {
							if refund.Type != enums.ProductTypeFinished {
								continue
							}
							if refund.OrderId != sale.Id {
								continue
							}
							if refund.Code != product.Code {
								continue
							}

							price = price.Sub(refund.Price)
						}

						req.MonthFinished = req.MonthFinished.Add(price)
					}
				case enums.ProductTypeOld:
					{
						var price decimal.Decimal
						switch product.Old.Product.RecycleType {
						case enums.ProductRecycleTypeExchange:
							{
								price = product.Old.RecyclePrice

								for _, refund := range month_refunds {
									if refund.Type != enums.ProductTypeOld {
										continue
									}
									if refund.OrderId != sale.Id {
										continue
									}
									if refund.Code != product.Code {
										continue
									}
									price = price.Sub(refund.Price)
								}
							}
						}
						req.MonthOld = req.MonthOld.Add(price)
					}
				case enums.ProductTypeAccessorie:
					{
						price := product.Accessorie.Price

						for _, refund := range month_refunds {
							if refund.Type != enums.ProductTypeAccessorie {
								continue
							}
							if refund.OrderId != sale.Id {
								continue
							}
							if refund.Name != product.Name {
								continue
							}

							price = price.Sub(refund.Price)
						}

						req.MonthAcciessorie = req.MonthAcciessorie.Add(price)
					}
				}
			}
		}

		allData[store.Name] = req
	}

	for _, staff := range staffs {
		// 要查看的门店
		stores := make(map[string]model.Store, 0)
		if staff.Identity == enums.IdentitySuperAdmin {
			for _, store := range allStore {
				stores[store.Name] = store
			}
		} else {
			for _, store := range staff.Stores {
				stores[store.Name] = store
			}
			for _, store := range staff.StoreSuperiors {
				stores[store.Name] = store
			}
			for _, region := range staff.Regions {
				for _, store := range region.Stores {
					stores[store.Name] = store
				}
			}
			for _, region := range staff.RegionSuperiors {
				for _, store := range region.Stores {
					stores[store.Name] = store
				}
			}
		}

		// 将门店按名称排序
		names := make([]string, 0, len(stores))
		for k := range stores {
			names = append(names, stores[k].Name)
		}
		sort.Strings(names)

		// 发送消息
		for _, name := range names {
			msg := message.NewMessage(context.Background())
			req, ok := allData[name]
			if !ok {
				continue
			}
			req.ToUser = []string{staff.Username}
			msg.SendReportStatisticMessage(&req)
		}
	}
}

func get_sales(store_id string, duration enums.Duration) ([]model.OrderSales, error) {
	var sales []model.OrderSales
	db := model.DB.Model(&model.OrderSales{})
	db = db.Where(&model.OrderSales{StoreId: store_id})
	db = db.Where("status IN (?)", []enums.OrderSalesStatus{
		enums.OrderSalesStatusComplete,
		enums.OrderSalesStatusRefund,
	})
	db = model.OrderSales{}.Preloads(db)
	if err := db.Scopes(model.DurationCondition(duration)).Find(&sales).Error; err != nil {
		return nil, err
	}

	return sales, nil
}

func get_refunds(store_id string, duration enums.Duration) ([]model.OrderRefund, error) {
	var refunds []model.OrderRefund
	db := model.DB.Model(&model.OrderRefund{})
	db = db.Where(&model.OrderRefund{StoreId: store_id})
	db = model.OrderRefund{}.Preloads(db)
	if err := db.Scopes(model.DurationCondition(duration)).Find(&refunds).Error; err != nil {
		return nil, err
	}

	return refunds, nil
}
