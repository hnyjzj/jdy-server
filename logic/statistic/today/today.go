package today

import (
	"jdy/enums"
	"jdy/logic/statistic"
)

type ToDayLogic struct {
	statistic.StatisticLogic
}

type DataReq struct {
	Duration  enums.Duration `json:"duration"`
	StartTime string         `json:"startTime"`
	EndTime   string         `json:"endTime"`
}
