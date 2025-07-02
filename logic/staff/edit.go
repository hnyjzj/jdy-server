package staff

import (
	"jdy/model"
	"jdy/types"

	"gorm.io/gorm"
)

// 修改员工
func (l *StaffLogic) StaffEdit(req *types.StaffEditReq) error {
	var (
		staff model.Staff
	)
	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 判断员工是否存在
		if err := tx.First(&staff, "id = ?", req.Id).Error; err != nil {
			return err
		}

		data := model.Staff{
			Phone:      req.Phone,
			Username:   req.Username,
			Nickname:   req.Nickname,
			Avatar:     req.Avatar,
			Email:      req.Email,
			Gender:     req.Gender,
			IsDisabled: req.IsDisabled,
			Identity:   req.Identity,
			RoleId:     req.RoleId,
		}

		if req.Password != "" {
			password, err := staff.HashPassword(&req.Password)
			if err != nil {
				return err
			}
			data.Password = password
		}

		// 修改员工信息
		if err := tx.Model(&staff).Updates(data).Error; err != nil {
			return err
		}

		// 关联门店
		if len(req.StoreIds) > 0 {
			var stores []model.Store
			if err := tx.Where("id in (?)", req.StoreIds).Find(&stores).Error; err != nil {
				return err
			}
			if err := tx.Model(&staff).Association("Stores").Replace(stores); err != nil {
				return err
			}
		}

		// 关联负责门店
		if len(req.StoreSuperiorIds) > 0 {
			var stores []model.Store
			if err := tx.Where("id in (?)", req.StoreSuperiorIds).Find(&stores).Error; err != nil {
				return err
			}
			if err := tx.Model(&staff).Association("StoreSuperiors").Replace(stores); err != nil {
				return err
			}
		}

		// 关联区域
		if len(req.RegionIds) > 0 {
			var regions []model.Region
			if err := tx.Where("id in (?)", req.RegionIds).Find(&regions).Error; err != nil {
				return err
			}
			if err := tx.Model(&staff).Association("Regions").Replace(regions); err != nil {
				return err
			}
		}

		// 关联负责区域
		if len(req.RegionSuperiorIds) > 0 {
			var regions []model.Region
			if err := tx.Where("id in (?)", req.RegionSuperiorIds).Find(&regions).Error; err != nil {
				return err
			}
			if err := tx.Model(&staff).Association("RegionSuperiors").Replace(regions); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
