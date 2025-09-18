package store

import (
	"errors"
	"jdy/model"
	"jdy/types"

	"github.com/gin-gonic/gin"
)

type StoreAdminLogic struct {
	Ctx   *gin.Context
	Staff *model.Staff
}

// 门店管理员列表
func (l *StoreAdminLogic) List(req *types.StoreAdminListReq) (*[]model.Staff, error) {
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
		return nil, errors.New("未入职该门店，无法查看管理员列表")
	}

	return &store.Admins, nil
}

// 添加门店管理员
func (l *StoreAdminLogic) Add(req *types.StoreAdminAddReq) error {
	// 查询门店
	var store model.Store
	if err := model.DB.First(&store, "id = ?", req.StoreId).Error; err != nil {
		return errors.New("门店不存在")
	}

	// 查询管理员
	var Admins []model.Staff
	if err := model.DB.Find(&Admins, "id IN (?)", req.AdminId).Error; err != nil {
		return errors.New("管理员不存在")
	}

	// 添加门店管理员
	if err := model.DB.Model(&store).Association("Admins").Append(&Admins); err != nil {
		return err
	}

	return nil
}

// 删除门店管理员
func (l *StoreAdminLogic) Del(req *types.StoreAdminDelReq) error {
	// 查询门店
	var store model.Store
	if err := model.DB.First(&store, "id = ?", req.StoreId).Error; err != nil {
		return errors.New("门店不存在")
	}

	// 查询管理员
	var Admins []model.Staff
	if err := model.DB.Find(&Admins, "id IN (?)", req.AdminId).Error; err != nil {
		return errors.New("管理员不存在")
	}

	// 删除门店管理员
	if err := model.DB.Model(&store).Association("Admins").Delete(&Admins); err != nil {
		return err
	}
	return nil
}

// 是否在门店
func (l *StoreAdminLogic) IsIn(req *types.StoreAdminIsInReq) (bool, error) {
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
	for _, staff := range store.Admins {
		if staff.Id == staff_id {
			res = true
			break
		}
	}

	return res, nil
}
