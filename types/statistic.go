package types

import "jdy/enums"

type StatisticStoreSalesTotalReq struct {
	Duration enums.Duration `json:"duration" required:"true"`
}

type StatisticTodaySalesReq struct {
	StoreId string `json:"store_id" required:"true"`
}
