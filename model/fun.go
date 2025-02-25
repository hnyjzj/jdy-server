package model

import (
	"jdy/enums"
	"time"

	"gorm.io/gorm"
)

func PageCondition(db *gorm.DB, page, limit int) *gorm.DB {
	if page == 0 {
		page = 1
	}

	switch {
	case limit > 100:
		limit = 100
	case limit <= 0:
		limit = 10
	}

	offset := (page - 1) * limit
	return db.Offset(offset).Limit(limit)
}

func DurationCondition(duration enums.Duration, fields ...string) func(db *gorm.DB) *gorm.DB {
	if len(fields) == 0 {
		fields = append(fields, "created_at")
	}
	return func(db *gorm.DB) *gorm.DB {
		var (
			now   = time.Now()
			start time.Time
			end   time.Time
		)
		switch duration {
		case enums.DurationToday: // 今日
			{
				start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
				end = start.Add(24 * time.Hour)
			}
		case enums.DurationYesterday: // 昨日
			{
				start = time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, now.Location())
				end = start.Add(24 * time.Hour)
			}
		case enums.DurationWeek: // 本周（周一为周起点）
			{
				start = time.Date(now.Year(), now.Month(), now.Day()-int(now.Weekday())+1, 0, 0, 0, 0, now.Location())
				end = start.Add(7 * 24 * time.Hour)
			}
		case enums.DurationLastWeek: // 上周
			{
				start = time.Date(now.Year(), now.Month(), now.Day()-int(now.Weekday())+1-7, 0, 0, 0, 0, now.Location())
				end = start.Add(7 * 24 * time.Hour)
			}
		case enums.DurationMonth: // 本月
			{
				start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
				end = start.Add(time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location()).Sub(start))
			}
		case enums.DurationLastMonth: // 上月
			{
				start = time.Date(now.Year(), now.Month()-1, 1, 0, 0, 0, 0, now.Location())
				end = start.Add(time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).Sub(start))
			}
		case enums.DurationQuarter: // 本季度
			{
				start = time.Date(now.Year(), time.Month(int(now.Month())/3*3+1), 1, 0, 0, 0, 0, now.Location())
				end = start.Add(time.Date(now.Year(), time.Month(int(now.Month())/3*3+4), 1, 0, 0, 0, 0, now.Location()).Sub(start))
			}
		case enums.DurationYear: // 本年
			{
				start = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
				end = start.Add(time.Date(now.Year()+1, 1, 1, 0, 0, 0, 0, now.Location()).Sub(start))
			}
		case enums.DurationLastYear: // 去年
			{
				start = time.Date(now.Year()-1, 1, 1, 0, 0, 0, 0, now.Location())
				end = start.Add(time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location()).Sub(start))
			}
		default:
			{
				return db
			}
		}

		return db.Where(fields[0]+" >= ? AND "+fields[0]+" < ?", start, end)
	}
}
