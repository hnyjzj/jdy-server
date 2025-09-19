package store

import (
	"errors"
	"jdy/model"
	"jdy/types"
	"jdy/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (l *StoreLogic) Create(ctx *gin.Context, req *types.StoreCreateReq) error {
	store, err := utils.StructToStruct[model.Store](req)
	if err != nil {
		return errors.New("验证信息失败")
	}
	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&store).Error; err != nil {
			return errors.New("创建失败")
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (l *StoreLogic) Update(ctx *gin.Context, req *types.StoreUpdateReq) error {

	// 查询门店信息
	var store model.Store
	if err := model.DB.First(&store, "id = ?", req.Id).Error; err != nil {
		return errors.New("门店不存在或已被删除")
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {

		data, err := utils.StructToStruct[model.Store](req)
		if err != nil {
			return errors.New("验证信息失败")
		}

		if err := tx.Model(&model.Store{}).Where("id = ?", store.Id).Updates(data).Error; err != nil {
			return errors.New("更新失败")
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

// 删除门店
func (l *StoreLogic) Delete(ctx *gin.Context, req *types.StoreDeleteReq) error {
	// 查询门店信息
	store := &model.Store{}
	if err := model.DB.First(store, "id = ?", req.Id).Error; err != nil {
		return errors.New("门店不存在或已被删除")
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(store).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
