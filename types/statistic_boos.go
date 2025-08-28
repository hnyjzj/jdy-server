package types

import (
	"jdy/enums"
	"time"
)

type StatisticBoosReq struct {
	Duration  enums.Duration `json:"duration" binding:"required"` // 时间范围
	StartTime time.Time      `json:"startTime"`                   // 开始时间
	EndTime   time.Time      `json:"endTime"`                     // 结束时间
}

type BoosFinishedStockReq struct {
	Duration  enums.Duration `json:"duration" binding:"required"` // 时间范围
	StartTime time.Time      `json:"startTime"`                   // 开始时间
	EndTime   time.Time      `json:"endTime"`                     // 结束时间
}
