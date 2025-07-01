package region

import (
	"errors"
	"jdy/model"
	"jdy/types"

	"github.com/gin-gonic/gin"
)

type RegionSuperiorLogic struct {
	Ctx   *gin.Context
	Staff *model.Staff
}

// 区域负责人列表
func (l *RegionSuperiorLogic) List(req *types.RegionSuperiorListReq) (*[]model.Staff, error) {
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

	for _, staff := range region.Superiors {
		if staff.Id == l.Staff.Id {
			inRegion = true
			break
		}
	}

	if !inRegion {
		return nil, errors.New("无权限访问")
	}

	return &region.Superiors, nil
}

// 添加区域负责人
func (l *RegionSuperiorLogic) Add(req *types.RegionSuperiorAddReq) error {
	// 查询区域
	var region model.Region
	if err := model.DB.First(&region, "id = ?", req.RegionId).Error; err != nil {
		return errors.New("区域不存在")
	}

	// 查询负责人
	var staff []model.Staff
	if err := model.DB.Find(&staff, "id IN (?)", req.SuperiorId).Error; err != nil {
		return errors.New("负责人不存在")
	}

	// 添加区域负责人
	if err := model.DB.Model(&region).Association("Superiors").Append(&staff); err != nil {
		return err
	}

	return nil
}

// 删除区域负责人
func (l *RegionSuperiorLogic) Del(req *types.RegionSuperiorDelReq) error {
	// 查询区域
	var region model.Region
	if err := model.DB.First(&region, "id = ?", req.RegionId).Error; err != nil {
		return errors.New("区域不存在")
	}

	// 查询负责人
	var staff []model.Staff
	if err := model.DB.Find(&staff, "id IN (?)", req.SuperiorId).Error; err != nil {
		return errors.New("负责人不存在")
	}

	// 删除区域负责人
	if err := model.DB.Model(&region).Association("Superiors").Delete(&staff); err != nil {
		return err
	}
	return nil
}
