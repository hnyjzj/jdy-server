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
	} else {
		return db
	}

	offset := (req.Page - 1) * req.Limit
	return db.Offset(offset).Limit(req.Limit)
}

func DurationCondition(duration enums.Duration, fields ...string) func(tx *gorm.DB) *gorm.DB {
	if len(fields) == 0 {
		fields = append(fields, "created_at")
	}

	var (
		now = time.Now()
	)

	return func(tx *gorm.DB) *gorm.DB {
		if err := duration.InMap(); err != nil {
			_ = tx.AddError(errors.New("duration not in enum"))
			return tx
		}

		switch duration {
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
				start, end := duration.GetTime(now)

				return tx.Where(fields[0]+" >= ? AND "+fields[0]+" < ?", start, end)
			}
		}
	}
}
