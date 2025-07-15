package types

import "jdy/enums"

type StatisticStoreSalesTotalReq struct {
	Duration enums.Duration `json:"duration" binding:"required"`
}

type StatisticTodaySalesReq struct {
	StoreId string `json:"store_id"`
}

type StatisticTodayProductReq struct {
	StoreId string `json:"store_id"`
}
