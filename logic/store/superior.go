package store

import (
	"errors"
	"jdy/model"
	"jdy/types"

	"github.com/gin-gonic/gin"
)

type StoreSuperiorLogic struct {
	Ctx   *gin.Context
	Staff *model.Staff
}

// 门店员工列表
func (l *StoreSuperiorLogic) List(req *types.StoreSuperiorListReq) (*[]model.Staff, error) {
	// 查询门店
	var (
		store model.Store
	)

	db := model.DB.Model(&model.Store{})
	db = store.Preloads(db)

	if err := db.First(&store, "id = ?", req.StoreId).Error; err != nil {
		return nil, errors.New("门店不存在")
	}

	if in := store.InStore(l.Staff.Id); !in {
		return nil, errors.New("未入职该门店，无法查看员工列表")
	}

	return &store.Superiors, nil
}

// 添加门店员工
func (l *StoreSuperiorLogic) Add(req *types.StoreSuperiorAddReq) error {
	// 查询门店
	var store model.Store
	if err := model.DB.First(&store, "id = ?", req.StoreId).Error; err != nil {
		return errors.New("门店不存在")
	}

	// 查询员工
	var superiors []model.Staff
	if err := model.DB.Find(&superiors, "id IN (?)", req.SuperiorId).Error; err != nil {
		return errors.New("员工不存在")
	}

	// 添加门店员工
	if err := model.DB.Model(&store).Association("Superiors").Append(&superiors); err != nil {
		return err
	}

	return nil
}

// 删除门店员工
func (l *StoreSuperiorLogic) Del(req *types.StoreSuperiorDelReq) error {
	// 查询门店
	var store model.Store
	if err := model.DB.First(&store, "id = ?", req.StoreId).Error; err != nil {
		return errors.New("门店不存在")
	}

	// 查询员工
	var superiors []model.Staff
	if err := model.DB.Find(&superiors, "id IN (?)", req.SuperiorId).Error; err != nil {
		return errors.New("员工不存在")
	}

	// 删除门店员工
	if err := model.DB.Model(&store).Association("Superiors").Delete(&superiors); err != nil {
		return err
	}
	return nil
}

// 是否在门店
func (l *StoreSuperiorLogic) IsIn(req *types.StoreSuperiorIsInReq) (bool, error) {
	var (
		store model.Store
		db    = model.DB

		staff_id string
		res      = false
	)

	db = store.Preloads(db)
	if err := db.First(&store, "id = ?", req.StoreId).Error; err != nil {
		return false, errors.New("门店不存在")
	}

	if req.StaffId != "" {
		staff_id = req.StaffId
	} else {
		staff_id = l.Staff.Id
	}
	for _, staff := range store.Superiors {
		if staff.Id == staff_id {
			res = true
			break
		}
	}

	return res, nil
}
