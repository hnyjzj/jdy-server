package crons

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

func init() {
	RegisterCrons(
		Crons{
			// // 每天晚上10 点半
			Spec: "0 30 22 * * *",
			Func: SendReportStatistic,
		},
	)
}

func SendReportStatistic() {
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
	sdb = sdb.Where("name <> ?", model.HeaderquartersName)
	sdb = sdb.Order("name desc")
	if err := sdb.Find(&allStore).Error; err != nil {
		log.Printf("SendReportStatistic error: %v", err.Error())
		return
	}

	allData := make(map[string]message.ReportStatisticMessage)
	for _, store := range allStore {
		// 今日订单
		orders_today, err := get_orders(store.Id, enums.DurationToday)
		if err != nil {
			log.Printf("SendReportStatistic error: %v", err.Error())
			continue
		}
		// 本月订单
		orders_month, err := get_orders(store.Id, enums.DurationMonth)
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

		for _, order := range orders_today {
			for _, product := range order.Products {
				switch product.Type {
				case enums.ProductTypeFinished:
					{
						req.TodayFinished = req.TodayFinished.Add(product.Finished.Price)
						class := product.Finished.Product.Class.String()
						req.TodayFinisheds[class] = req.TodayFinisheds[class].Add(product.Finished.Price)
					}
				case enums.ProductTypeOld:
					{
						req.TodayOld = req.TodayOld.Add(product.Old.RecyclePrice)
					}
				case enums.ProductTypeAccessorie:
					{
						req.TodayAcciessorie = req.TodayAcciessorie.Add(product.Accessorie.Price)
					}
				}
			}
		}

		for _, order := range orders_month {
			for _, product := range order.Products {
				switch product.Type {
				case enums.ProductTypeFinished:
					{
						req.MonthFinished = req.MonthFinished.Add(product.Finished.Price)
					}
				case enums.ProductTypeOld:
					{
						req.MonthOld = req.MonthOld.Add(product.Old.RecyclePrice)
					}
				case enums.ProductTypeAccessorie:
					{
						req.MonthAcciessorie = req.MonthAcciessorie.Add(product.Accessorie.Price)
					}
				}
			}
		}

		allData[store.Id] = req
	}

	for _, staff := range staffs {
		// 要查看的门店
		stores := make(map[string]model.Store, 0)
		if staff.Identity == enums.IdentitySuperAdmin {
			for _, store := range allStore {
				stores[store.Id] = store
			}
		} else {
			for _, store := range staff.Stores {
				stores[store.Id] = store
			}
			for _, store := range staff.StoreSuperiors {
				stores[store.Id] = store
			}
			for _, store := range staff.Regions {
				for _, s := range store.Stores {
					stores[s.Id] = s
				}
			}
			for _, store := range staff.RegionSuperiors {
				for _, s := range store.Stores {
					stores[s.Id] = s
				}
			}
		}

		// 将门店按名称排序
		keys := make([]string, 0, len(stores))
		for k := range stores {
			keys = append(keys, stores[k].Name)
		}
		sort.Strings(keys)

		// 发送消息
		for _, k := range keys {
			store := stores[k]
			msg := message.NewMessage(context.Background())
			req := allData[store.Id]
			req.ToUser = []string{staff.Username}
			msg.SendReportStatisticMessage(&req)
		}
	}
}

func get_orders(store_id string, duration enums.Duration) ([]model.OrderSales, error) {
	var orders []model.OrderSales
	db := model.DB.Model(&model.OrderSales{})
	db = db.Where(&model.OrderSales{StoreId: store_id})
	db = db.Where("status IN (?)", []enums.OrderSalesStatus{
		enums.OrderSalesStatusComplete,
		enums.OrderSalesStatusRefund,
	})
	db = model.OrderSales{}.Preloads(db)
	if err := db.Scopes(model.DurationCondition(duration)).Find(&orders).Error; err != nil {
		return nil, err
	}

	return orders, nil
}
