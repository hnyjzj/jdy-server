package enums

import (
	"errors"
)

/* 时间区间 */
// 今天、昨天、本周、上周、本月、上月、本季度、上季度、今年、去年、自定义
type Duration int

const (
	DurationToday       Duration = iota + 1 // 今天
	DurationYesterday                       // 昨天
	DurationWeek                            // 本周
	DurationLastWeek                        // 上周
	DurationMonth                           // 本月
	DurationLastMonth                       // 上月
	DurationQuarter                         // 本季度
	DurationLastQuarter                     // 上季度
	DurationYear                            // 今年
	DurationLastYear                        // 去年
	DurationCustom                          // 自定义
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
	DurationCustom:      "自定义",
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
