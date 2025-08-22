package types

type StatisticTodaySalesReq struct {
	StoreId string `json:"store_id"` // 店铺id
}

type StatisticTodayProductReq struct {
	StoreId string `json:"store_id"` // 店铺id
}
