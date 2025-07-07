package store

import (
	"errors"
	"jdy/model"
	"jdy/types"

	"github.com/gin-gonic/gin"
)

type StoreStaffLogic struct {
	Ctx   *gin.Context
	Staff *model.Staff
}

// 门店员工列表
func (l *StoreStaffLogic) List(req *types.StoreStaffListReq) (*[]model.Staff, error) {
	// 查询门店
	var (
		store   model.Store
		inStore = false
	)
	if err := model.DB.Preload("Staffs").First(&store, "id = ?", req.StoreId).Error; err != nil {
		return nil, errors.New("门店不存在")
	}
	for _, staff := range store.Staffs {
		if staff.Id == l.Staff.Id {
			inStore = true
			break
		}
	}
	if !inStore {
		return nil, errors.New("未入职该门店，无法查看员工列表")
	}

	return &store.Staffs, nil
}

// 添加门店员工
func (l *StoreStaffLogic) Add(req *types.StoreStaffAddReq) error {
	// 查询门店
	var store model.Store
	if err := model.DB.First(&store, "id = ?", req.StoreId).Error; err != nil {
		return errors.New("门店不存在")
	}

	// 查询员工
	var staff []model.Staff
	if err := model.DB.Find(&staff, "id IN (?)", req.StaffId).Error; err != nil {
		return errors.New("员工不存在")
	}

	// 添加门店员工
	if err := model.DB.Model(&store).Association("Staffs").Append(&staff); err != nil {
		return err
	}

	return nil
}

// 删除门店员工
func (l *StoreStaffLogic) Del(req *types.StoreStaffDelReq) error {
	// 查询门店
	var store model.Store
	if err := model.DB.First(&store, "id = ?", req.StoreId).Error; err != nil {
		return errors.New("门店不存在")
	}

	// 查询员工
	var staff []model.Staff
	if err := model.DB.Find(&staff, "id IN (?)", req.StaffId).Error; err != nil {
		return errors.New("员工不存在")
	}

	// 删除门店员工
	if err := model.DB.Model(&store).Association("Staffs").Delete(&staff); err != nil {
		return err
	}
	return nil
}
