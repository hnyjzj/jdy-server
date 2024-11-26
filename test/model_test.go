package test

import (
	"jdy/config"
	"jdy/model"
	"testing"
)

func init() {
	// 初始化配置
	config.Init()
	// 初始化数据库
	model.Init()
}

func TestStaff(t *testing.T) {
	var (
	// user usermodel.Staff
	)
	// phone := "18503857576"
	// user.Phone = &phone
	// tx := user.Models().Begin()
	// if err := tx.Create(&user); err != nil {
	// 	t.Error(err)
	// 	tx.Rollback()
	// }
	// tx.Commit()
	// t.Log("\n", "user:", user, "\n")

	// // if err := tx.Delete(&user); err != nil {
	// // 	t.Error(err)
	// // 	tx.DB.Rollback()
	// // }

	// res, err := user.Models().GetList()
	// if err != nil {
	// 	t.Error(err)
	// }
	t.Log("\n", "res:", "\n")
}
