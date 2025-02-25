package enums

import (
	"errors"
)

/* 时间区间 */
// 今天、昨天、本周、上周、本月、上月、本季度、上季度、今年、去年
type Duration string

const (
	DurationToday       Duration = "today"        // 今天
	DurationYesterday   Duration = "yesterday"    // 昨天
	DurationWeek        Duration = "week"         // 本周
	DurationLastWeek    Duration = "last_week"    // 上周
	DurationMonth       Duration = "month"        // 本月
	DurationLastMonth   Duration = "last_month"   // 上月
	DurationQuarter     Duration = "quarter"      // 本季度
	DurationLastQuarter Duration = "last_quarter" // 上季度
	DurationYear        Duration = "year"         // 今年
	DurationLastYear    Duration = "last_year"    // 去年
)

var DurationMap = map[Duration]string{
	DurationToday:       "今天",
	DurationYesterday:   "昨天",
	DurationWeek:        "本周",
	DurationLastWeek:    "上周",
	DurationMonth:       "本月",
	DurationLastMonth:   "上月",
	DurationQuarter:     "本季度",
	DurationLastQuarter: "上季度",
	DurationYear:        "今年",
	DurationLastYear:    "去年",
}

func (p Duration) ToMap() any {
	return DurationMap
}

func (p Duration) InMap() error {
	if _, ok := DurationMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}
