package region

import (
	"errors"
	"jdy/model"
	"jdy/types"

	"github.com/gin-gonic/gin"
)

type RegionStoreLogic struct {
	Ctx   *gin.Context
	Staff *model.Staff
}

// 区域门店列表
func (l *RegionStoreLogic) List(req *types.RegionStoreListReq) (*[]model.Store, error) {
	// 查询区域
	var (
		region model.Region
	)

	db := model.DB.Model(&model.Region{})
	db = region.Preloads(db)
	if err := db.First(&region, "id = ?", req.RegionId).Error; err != nil {
		return nil, errors.New("区域不存在")
	}

	return &region.Stores, nil
}

// 添加区域门店
func (l *RegionStoreLogic) Add(req *types.RegionStoreAddReq) error {
	// 查询区域
	var region model.Region
	if err := model.DB.First(&region, "id = ?", req.RegionId).Error; err != nil {
		return errors.New("区域不存在")
	}

	// 查询门店
	var store []model.Store
	if err := model.DB.Find(&store, "id = ?", req.StoreId).Error; err != nil {
		return errors.New("门店不存在")
	}

	// 添加区域门店
	if err := model.DB.Model(&region).Association("Stores").Append(&store); err != nil {
		return err
	}

	return nil
}

// 删除区域门店
func (l *RegionStoreLogic) Del(req *types.RegionStoreDelReq) error {
	// 查询区域
	var region model.Region
	if err := model.DB.First(&region, "id = ?", req.RegionId).Error; err != nil {
		return errors.New("区域不存在")
	}

	// 查询门店
	var store []model.Store
	if err := model.DB.Find(&store, "id = ?", req.StoreId).Error; err != nil {
		return errors.New("门店不存在")
	}

	// 删除区域门店
	if err := model.DB.Model(&region).Association("Stores").Delete(&store); err != nil {
		return err
	}
	return nil
}
