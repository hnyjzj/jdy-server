package types

import (
	"time"
)

type StatisticTargetReq struct {
	StoreId string `json:"store_id" binding:"required"`
}

type StatisticTargetResp struct {
	TargetId    string    `json:"target_id"`    // 目标ID
	StartTime   time.Time `json:"start_time"`   // 开始时间
	EndTime     time.Time `json:"end_time"`     // 结束时间
	Purpose     string    `json:"purpose"`      // 目标值
	Achieve     string    `json:"achieve"`      // 实际值
	Remainder   string    `json:"remainder"`    // 剩余值
	AchieveRate string    `json:"achieve_rate"` // 达成率
}
