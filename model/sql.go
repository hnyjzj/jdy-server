package model

import (
	"gorm.io/gorm"
)

type Sql[T any] struct {
	db *gorm.DB
}

// 初始化模型
func (Sql[T]) Models() Sql[T] {
	// 初始化结构
	var sql Sql[T]
	// 初始化模型
	sql.db = DB.Model(new(T))

	return sql
}

// 获取单条数据
func (sql Sql[T]) Get() (*T, error) {
	var res T
	err := sql.db.First(&res).Error
	return &res, err
}

// 获取多条数据
func (sql Sql[T]) GetList() ([]*T, error) {
	var list []*T
	err := sql.db.Find(&list).Error
	return list, err
}

// 创建数据
func (sql Sql[T]) Create(v *T) error {
	return sql.db.Create(&v).Error
}

// 更新数据
func (sql Sql[T]) Update(v *T) error {
	return sql.db.Save(&v).Error
}

// 删除数据
func (sql Sql[T]) Delete(v *T, real ...bool) error {
	if len(real) > 0 && real[0] {
		return sql.db.Unscoped().Delete(&v).Error
	} else {
		return sql.db.Delete(&v).Error
	}
}
