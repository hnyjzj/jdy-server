package enums

import (
	"errors"
	"time"
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

func (p Duration) GetTime(now time.Time) (start, end time.Time) {
	if now.IsZero() {
		now = time.Now()
	}

	switch p {
	case DurationToday: // 今天
		start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		end = start.Add(24*time.Hour - time.Nanosecond)

	case DurationYesterday: // 昨天
		start = time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, now.Location())
		end = start.Add(24*time.Hour - time.Nanosecond)

	case DurationWeek: // 本周
		start = time.Date(now.Year(), now.Month(), now.Day()-int(now.Weekday())+1, 0, 0, 0, 0, now.Location())
		end = start.Add(7*24*time.Hour - time.Nanosecond)

	case DurationLastWeek: // 上周
		start = time.Date(now.Year(), now.Month(), now.Day()-int(now.Weekday())+1-7, 0, 0, 0, 0, now.Location())
		end = start.Add(7*24*time.Hour - time.Nanosecond)

	case DurationMonth: // 本月
		start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		end = start.Add(time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location()).Sub(start)).Add(-time.Nanosecond)

	case DurationLastMonth: // 上月
		start = time.Date(now.Year(), now.Month()-1, 1, 0, 0, 0, 0, now.Location())
		end = start.Add(time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).Sub(start)).Add(-time.Nanosecond)

	case DurationQuarter: // 本季度
		start = time.Date(now.Year(), time.Month(int(now.Month()-1)/3*3+1), 1, 0, 0, 0, 0, now.Location())
		end = start.Add(time.Date(now.Year(), time.Month(int(now.Month()-1)/3*3+4), 1, 0, 0, 0, 0, now.Location()).Sub(start)).Add(-time.Nanosecond)

	case DurationLastQuarter: // 上季度
		start = time.Date(now.Year(), time.Month(int(now.Month()-1)/3*3+1)-3, 1, 0, 0, 0, 0, now.Location())
		end = start.Add(time.Date(now.Year(), time.Month(int(now.Month()-1)/3*3+4)-3, 1, 0, 0, 0, 0, now.Location()).Sub(start)).Add(-time.Nanosecond)

	case DurationYear: // 今年
		start = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
		end = start.Add(time.Date(now.Year()+1, 1, 1, 0, 0, 0, 0, now.Location()).Sub(start)).Add(-time.Nanosecond)

	case DurationLastYear: // 去年
		start = time.Date(now.Year()-1, 1, 1, 0, 0, 0, 0, now.Location())
		end = start.Add(time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location()).Sub(start)).Add(-time.Nanosecond)
	}

	return start, end
}
