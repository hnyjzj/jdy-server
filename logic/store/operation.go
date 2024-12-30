package store

import (
	"errors"
	"jdy/model"
	"jdy/types"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (l *StoreLogic) Create(ctx *gin.Context, req *types.StoreCreateReq) error {
	store := &model.Store{
		ParentId: req.ParentId,

		Name:     req.Name,
		Address:  req.Address,
		Contact:  req.Contact,
		Logo:     req.Logo,
		Sort:     req.Sort,
		Province: req.Province,
		City:     req.City,
		District: req.District,
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(store).Error; err != nil {
			return err
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
	if err := model.DB.First(&store, req.Id).Error; err != nil {
		return errors.New("门店不存在或已被删除")
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {

		if err := tx.Model(&store).Updates(model.Store{
			ParentId: req.ParentId,

			Name:     req.Name,
			Address:  req.Address,
			Contact:  req.Contact,
			Logo:     req.Logo,
			Sort:     req.Sort,
			Province: req.Province,
			City:     req.City,
			District: req.District,
		}).Error; err != nil {
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
	if err := model.DB.First(store, req.Id).Error; err != nil {
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
