package today

import (
	"errors"
	"jdy/enums"
	"jdy/logic/statistic"
	"time"
)

type ToDayLogic struct {
	statistic.StatisticLogic
}

type DataReq struct {
	Duration  enums.Duration `json:"duration"`
	StartTime string         `json:"startTime"`
	EndTime   string         `json:"endTime"`
}

func (req *DataReq) Validate() error {
	if err := req.Duration.InMap(); err != nil {
		req.Duration = enums.DurationToday
	}
	if req.Duration == enums.DurationCustom {
		if req.StartTime == "" || req.EndTime == "" {
			return errors.New("时间格式错误")
		}
	}

	start, end, err := req.Duration.GetTime(time.Now())
	if err != nil {
		return err
	}
	req.StartTime = start.Format(time.RFC3339)
	req.EndTime = end.Format(time.RFC3339)

	return nil
}
