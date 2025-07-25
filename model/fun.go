package model

import (
	"errors"
	"jdy/enums"
	"jdy/types"
	"time"

	"gorm.io/gorm"
)

// 分页条件
func PageCondition(db *gorm.DB, req *types.PageReq) *gorm.DB {
	if req == nil {
		_ = db.AddError(errors.New("分页参数不能为空"))
		return db
	}
	if !req.All {
		if req.Page == 0 {
			req.Page = 1
		}

		switch {
		case req.Limit <= 0:
			req.Limit = 10
		}
	}

	offset := (req.Page - 1) * req.Limit
	return db.Offset(offset).Limit(req.Limit)
}

func DurationCondition(duration enums.Duration, fields ...string) func(tx *gorm.DB) *gorm.DB {
	if len(fields) == 0 {
		fields = append(fields, "created_at")
	}
	return func(tx *gorm.DB) *gorm.DB {
		var (
			now = time.Now()
		)

		switch duration {
		case enums.DurationToday: // 今日
			{
				start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
				end := start.Add(24*time.Hour - time.Nanosecond)

				return tx.Where(fields[0]+" >= ? AND "+fields[0]+" < ?", start, end)
			}
		case enums.DurationYesterday: // 昨日
			{
				start := time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, now.Location())
				end := start.Add(24*time.Hour - time.Nanosecond)

				return tx.Where(fields[0]+" >= ? AND "+fields[0]+" < ?", start, end)
			}
		case enums.DurationWeek: // 本周（周一为周起点）
			{
				start := time.Date(now.Year(), now.Month(), now.Day()-int(now.Weekday())+1, 0, 0, 0, 0, now.Location())
				end := start.Add(7*24*time.Hour - time.Nanosecond)

				return tx.Where(fields[0]+" >= ? AND "+fields[0]+" < ?", start, end)
			}
		case enums.DurationLastWeek: // 上周
			{
				start := time.Date(now.Year(), now.Month(), now.Day()-int(now.Weekday())+1-7, 0, 0, 0, 0, now.Location())
				end := start.Add(7*24*time.Hour - time.Nanosecond)

				return tx.Where(fields[0]+" >= ? AND "+fields[0]+" < ?", start, end)
			}
		case enums.DurationMonth: // 本月
			{
				start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
				end := start.Add(time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location()).Sub(start)).Add(-time.Nanosecond)

				return tx.Where(fields[0]+" >= ? AND "+fields[0]+" < ?", start, end)
			}
		case enums.DurationLastMonth: // 上月
			{
				start := time.Date(now.Year(), now.Month()-1, 1, 0, 0, 0, 0, now.Location())
				end := start.Add(time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).Sub(start)).Add(-time.Nanosecond)

				return tx.Where(fields[0]+" >= ? AND "+fields[0]+" < ?", start, end)
			}
		case enums.DurationQuarter: // 本季度
			{
				start := time.Date(now.Year(), time.Month(int(now.Month())/3*3+1), 1, 0, 0, 0, 0, now.Location())
				end := start.Add(time.Date(now.Year(), time.Month(int(now.Month())/3*3+4), 1, 0, 0, 0, 0, now.Location()).Sub(start)).Add(-time.Nanosecond)

				return tx.Where(fields[0]+" >= ? AND "+fields[0]+" < ?", start, end)
			}
		case enums.DurationLastQuarter: // 上季度
			{
				start := time.Date(now.Year(), time.Month(int(now.Month())/3*3+1)-3, 1, 0, 0, 0, 0, now.Location())
				end := start.Add(time.Date(now.Year(), time.Month(int(now.Month())/3*3+4)-3, 1, 0, 0, 0, 0, now.Location()).Sub(start)).Add(-time.Nanosecond)

				return tx.Where(fields[0]+" >= ? AND "+fields[0]+" < ?", start, end)
			}
		case enums.DurationYear: // 今年
			{
				start := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
				end := start.Add(time.Date(now.Year()+1, 1, 1, 0, 0, 0, 0, now.Location()).Sub(start)).Add(-time.Nanosecond)

				return tx.Where(fields[0]+" >= ? AND "+fields[0]+" < ?", start, end)
			}
		case enums.DurationLastYear: // 去年
			{
				start := time.Date(now.Year()-1, 1, 1, 0, 0, 0, 0, now.Location())
				end := start.Add(time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location()).Sub(start)).Add(-time.Nanosecond)

				return tx.Where(fields[0]+" >= ? AND "+fields[0]+" < ?", start, end)
			}
		case enums.DurationCustom: // 自定义
			{
				if len(fields) < 3 || fields[1] == "" || fields[2] == "" {
					_ = tx.AddError(errors.New("start or end time is empty"))
					return tx
				}

				start, err := time.ParseInLocation(time.RFC3339, fields[1], now.Location())
				if err != nil {
					_ = tx.AddError(errors.New("start time format error"))
					return tx
				}

				end, err := time.ParseInLocation(time.RFC3339, fields[2], now.Location())
				if err != nil {
					_ = tx.AddError(errors.New("end time format error"))
					return tx
				}

				return tx.Where(fields[0]+" >= ? AND "+fields[0]+" < ?", start, end)
			}
		default:
			{
				_ = tx.AddError(errors.New("duration not in enum"))
				return tx
			}
		}
	}
}
