package region

import (
	"errors"
	"jdy/model"
	"jdy/types"

	"github.com/gin-gonic/gin"
)

type RegionAdminLogic struct {
	Ctx   *gin.Context
	Staff *model.Staff
}

// 区域管理员列表
func (l *RegionAdminLogic) List(req *types.RegionAdminListReq) (*[]model.Staff, error) {
	// 查询区域
	var (
		region model.Region
	)

	db := model.DB.Model(&model.Region{})
	db = region.Preloads(db)

	if err := db.First(&region, "id = ?", req.RegionId).Error; err != nil {
		return nil, errors.New("区域不存在")
	}

	return &region.Admins, nil
}

// 添加区域管理员
func (l *RegionAdminLogic) Add(req *types.RegionAdminAddReq) error {
	// 查询区域
	var region model.Region
	if err := model.DB.First(&region, "id = ?", req.RegionId).Error; err != nil {
		return errors.New("区域不存在")
	}

	// 查询管理员
	var staff []model.Staff
	if err := model.DB.Find(&staff, "id IN (?)", req.AdminId).Error; err != nil {
		return errors.New("管理员不存在")
	}

	// 添加区域管理员
	if err := model.DB.Model(&region).Association("Admins").Append(&staff); err != nil {
		return err
	}

	return nil
}

// 删除区域管理员
func (l *RegionAdminLogic) Del(req *types.RegionAdminDelReq) error {
	// 查询区域
	var region model.Region
	if err := model.DB.First(&region, "id = ?", req.RegionId).Error; err != nil {
		return errors.New("区域不存在")
	}

	// 查询管理员
	var staff []model.Staff
	if err := model.DB.Find(&staff, "id IN (?)", req.AdminId).Error; err != nil {
		return errors.New("管理员不存在")
	}

	// 删除区域管理员
	if err := model.DB.Model(&region).Association("Admins").Delete(&staff); err != nil {
		return err
	}
	return nil
}
