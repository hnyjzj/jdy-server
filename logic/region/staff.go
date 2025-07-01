package region

import (
	"errors"
	"jdy/model"
	"jdy/types"

	"github.com/gin-gonic/gin"
)

type RegionStaffLogic struct {
	Ctx   *gin.Context
	Staff *model.Staff
}

// 区域员工列表
func (l *RegionStaffLogic) List(req *types.RegionStaffListReq) (*[]model.Staff, error) {
	// 查询区域
	var (
		region   model.Region
		inRegion = false
	)

	db := model.DB.Model(&model.Region{})
	db = region.Preloads(db)
	if err := db.First(&region, "id = ?", req.RegionId).Error; err != nil {
		return nil, errors.New("区域不存在")
	}

	for _, staff := range region.Staffs {
		if staff.Id == l.Staff.Id {
			inRegion = true
			break
		}
	}

	if !inRegion {
		return nil, errors.New("无权限访问")
	}

	return &region.Staffs, nil
}

// 添加区域员工
func (l *RegionStaffLogic) Add(req *types.RegionStaffAddReq) error {
	// 查询区域
	var region model.Region
	if err := model.DB.First(&region, "id = ?", req.RegionId).Error; err != nil {
		return errors.New("区域不存在")
	}

	// 查询员工
	var staff []model.Staff
	if err := model.DB.Find(&staff, "id IN (?)", req.StaffId).Error; err != nil {
		return errors.New("员工不存在")
	}

	// 添加区域员工
	if err := model.DB.Model(&region).Association("Staffs").Append(&staff); err != nil {
		return err
	}

	return nil
}

// 删除区域员工
func (l *RegionStaffLogic) Del(req *types.RegionStaffDelReq) error {
	// 查询区域
	var region model.Region
	if err := model.DB.First(&region, "id = ?", req.RegionId).Error; err != nil {
		return errors.New("区域不存在")
	}

	// 查询员工
	var staff []model.Staff
	if err := model.DB.Find(&staff, "id IN (?)", req.StaffId).Error; err != nil {
		return errors.New("员工不存在")
	}

	// 删除区域员工
	if err := model.DB.Model(&region).Association("Staffs").Delete(&staff); err != nil {
		return err
	}
	return nil
}
