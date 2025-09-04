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
	var (
		now       = time.Now()
		def_field = "created_at"
	)

	return func(tx *gorm.DB) *gorm.DB {
		if err := duration.InMap(); err != nil {
			_ = tx.AddError(errors.New("时间范围不合法"))
			return tx
		}

		switch duration {
		case enums.DurationCustom: // 自定义
			{
				if len(fields) < 3 || fields[1] == "" || fields[2] == "" {
					_ = tx.AddError(errors.New("自定义时间范围格式不正确"))
					return tx
				}

				field, stime, etime := fields[0], fields[1], fields[2]
				if field == "" {
					field = def_field
				}

				start, end, err := duration.GetTime(now, stime, etime)
				if err != nil {
					_ = tx.AddError(err)
					return tx
				}

				return tx.Where(field+" >= ? AND "+field+" < ?", start, end)
			}
		default:
			{
				if len(fields) == 0 {
					fields = append(fields, def_field)
				}

				start, end, err := duration.GetTime(now)
				if err != nil {
					_ = tx.AddError(err)
					return tx
				}

				return tx.Where(fields[0]+" >= ? AND "+fields[0]+" < ?", start, end)
			}
		}
	}
}
