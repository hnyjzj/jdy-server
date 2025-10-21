package target

import (
	"errors"
	"jdy/enums"
	"jdy/model"
	"jdy/types"
	"jdy/utils"

	"gorm.io/gorm"
)

func (l *Logic) Create(req *types.TargetCreateReq) error {

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

		// 创建目标
		if err := tx.Create(&data).Error; err != nil {
			return errors.New("创建目标失败")
		}

		// 根据目标对象类型创建群组或个人
		switch data.Object {
		case enums.TargetObjectGroup:
			{
				// 创建群组
				var groups []model.TargetGroup
				for _, group := range req.Groups {
					groups = append(groups, model.TargetGroup{
						TargetId: data.Id,
						Name:     group.Name,
					})
				}
				if err := tx.Create(&groups).Error; err != nil {
					return errors.New("创建群组失败")
				}
				// 创建员工
				for _, group := range groups {
					for _, rg := range req.Groups {
						for i, personal := range req.Personals {
							if personal.GroupId == rg.Id {
								req.Personals[i].GroupId = group.Id
							}
						}
					}
				}
				// 创建个人
				var personal []model.TargetPersonal
				for _, p := range req.Personals {
					personal = append(personal, model.TargetPersonal{
						TargetId: data.Id,
						StaffId:  p.StaffId,
						GroupId:  p.GroupId,
						IsLeader: p.IsLeader,
						Purpose:  *p.Purpose,
					})
				}
				if err := tx.Create(&personal).Error; err != nil {
					return errors.New("创建群组员工失败")
				}
			}
		case enums.TargetObjectPersonal:
			{
				// 创建个人
				var personal []model.TargetPersonal
				for _, p := range req.Personals {
					personal = append(personal, model.TargetPersonal{
						TargetId: data.Id,
						StaffId:  p.StaffId,
						Purpose:  *p.Purpose,
					})
				}
				if err := tx.Create(&personal).Error; err != nil {
					return errors.New("创建个人失败")
				}
			}
		default:
			{
				return errors.New("目标对象类型错误")
			}
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
