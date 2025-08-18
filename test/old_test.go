package test

import (
	"jdy/config"
	"jdy/model"
	"jdy/utils"
	"strings"
	"testing"

	"gorm.io/gorm"
)

func init() {
	// 初始化配置
	config.Init()
	// 初始化数据库
	model.Init()
}

func TestOldCreateCode(t *testing.T) {
	var olds []model.ProductOld

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Unscoped().Where("code = '' or code is null or code = ' '").Find(&olds).Error; err != nil {
			t.Error(err)
		}

		for _, old := range olds {
			t.Log(old.Id)

			if err := tx.Unscoped().Model(&model.ProductOld{}).Where("id = ?", old.Id).Update("code", strings.ToUpper("JL"+utils.RandomCode(8))).Error; err != nil {
				t.Error(err)
			}
		}

		return nil
	}); err != nil {
		t.Error(err)
	}
}
