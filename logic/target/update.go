package target

import (
	"errors"
	"jdy/enums"
	"jdy/model"
	"jdy/types"
	"jdy/utils"

	"gorm.io/gorm"
)

func (l *Logic) Update(req *types.TargetUpdateReq) error {
	var target model.Target
	if err := model.DB.Where("id = ?", req.Id).First(&target).Error; err != nil {
		return errors.New("目标不存在")
	}

	// 验证参数
	data, err := utils.StructToStruct[model.Target](req)
	if err != nil {
		return errors.New("验证参数失败")
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		if data.IsDefault {
			// 将所有目标设为非默认
			if err := tx.Model(&model.Target{}).Where(&model.Target{
				IsDefault: true,
				StoreId:   data.StoreId,
			}).Update("is_default", false).Error; err != nil {
				return errors.New("将所有目标设为非默认失败")
			}
		}

		// 更新目标
		if err := tx.Model(model.Target{}).Where("id = ?", data.Id).Updates(data).Error; err != nil {
			return errors.New("更新目标失败")
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (l *Logic) UpdateGroup(req *types.TargetUpdateGroupReq) error {
	var group model.TargetGroup
	if err := model.DB.Where("id = ?", req.Id).First(&group).Error; err != nil {
		return errors.New("目标不存在")
	}

	// 验证参数
	data, err := utils.StructToStruct[model.TargetGroup](req)
	if err != nil {
		return errors.New("验证参数失败")
	}

	if err := model.DB.Model(model.TargetGroup{}).Where("id = ?", data.Id).Updates(data).Error; err != nil {
		return errors.New("更新目标失败")
	}

	return nil
}

func (l *Logic) UpdatePersonal(req *types.TargetUpdatePersonalReq) error {
	var (
		personal model.TargetPersonal
	)

	db := model.DB.Model(model.TargetPersonal{})
	db = personal.Preloads(db)
	if err := db.First(&personal, "id = ?", req.Id).Error; err != nil {
		return errors.New("目标不存在")
	}

	// 验证参数
	data, err := utils.StructToStruct[model.TargetPersonal](req)
	if err != nil {
		return errors.New("验证参数失败")
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 如果目标类型为团队目标
		if personal.Target.Object == enums.TargetObjectGroup {
			// 将目标设为组长
			if data.IsLeader {
				// 将所有成员目标设为非组长
				if err := tx.Model(model.TargetPersonal{}).Where(&model.TargetPersonal{
					TargetId: personal.TargetId,
					GroupId:  personal.GroupId,
				}).Update("is_leader", false).Error; err != nil {
					return errors.New("将所有成员目标设为非组长失败")
				}
			}
		}

		if err := tx.Model(model.TargetPersonal{}).Where("id = ?", data.Id).Select([]string{
			"is_leader",
			"purpose",
		}).Updates(data).Error; err != nil {
			return errors.New("更新目标失败")
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
