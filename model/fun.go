package model

import (
	"gorm.io/gorm"
)

// PageCondition applies pagination conditions to a Gorm database query. It normalizes the page and limit parameters
// to ensure reasonable defaults and calculates the appropriate offset for database pagination.
// If page is zero, it defaults to the first page. If limit exceeds 100, it is capped at 100,
// and if limit is zero or negative, it defaults to 10 records per page.
// Returns the modified database query with offset and limit applied.
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
