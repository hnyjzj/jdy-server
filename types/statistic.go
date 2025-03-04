package types

import "jdy/enums"

type StatisticStoreSalesTotalReq struct {
	Duration enums.Duration `json:"duration" binding:"required"`
}

type StatisticTodaySalesReq struct {
	StoreId string `json:"store_id" binding:"required"`
}

type StatisticTodayProductReq struct {
	StoreId string `json:"store_id" binding:"required"`
}
