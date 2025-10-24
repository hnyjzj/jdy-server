package target

import (
	"errors"
	"jdy/enums"
	"jdy/model"
	"jdy/types"
	"jdy/utils"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (l *Logic) Create(req *types.TargetCreateReq) (*model.Target, error) {

	// 验证参数
	data, err := utils.StructToStruct[model.Target](req)
	if err != nil {
		return nil, errors.New("验证参数失败")
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

		// 创建目标(跳过关联)
		if err := tx.Omit(clause.Associations).Create(&data).Error; err != nil {
			return errors.New("创建目标失败")
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &data, nil
}

func (l *Logic) CreateGroup(req *types.TargetGroupCreateReq) (*model.TargetGroup, error) {
	// 验证参数
	data, err := utils.StructToStruct[model.TargetGroup](req)
	if err != nil {
		return nil, errors.New("验证参数失败")
	}
	var (
		target model.Target
	)
	if err := model.DB.Preload("Groups").First(&target, "id = ?", req.TargetId).Error; err != nil {
		return nil, errors.New("目标不存在")
	}
	if target.Object != enums.TargetObjectGroup {
		return nil, errors.New("非群组目标不能创建群组")
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		for _, group := range target.Groups {
			if group.Name == data.Name {
				return errors.New("目标组名称已存在")
			}
		}

		// 创建群组(跳过关联)
		if err := tx.Omit(clause.Associations).Create(&data).Error; err != nil {
			return errors.New("创建群组失败")
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &data, nil
}

func (l *Logic) CreatePersonal(req *types.TargetPersonalCreateReq) (*model.TargetPersonal, error) {
	data, err := utils.StructToStruct[model.TargetPersonal](req)
	if err != nil {
		return nil, errors.New("验证参数失败")
	}
	var (
		target model.Target
	)
	if err := model.DB.Preload("Personals").First(&target, "id = ?", req.TargetId).Error; err != nil {
		return nil, errors.New("目标不存在")
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 如果目标类型为团队目标
		if target.Object == enums.TargetObjectGroup {
			// 获取目标组
			var group model.TargetGroup
			if err := tx.First(&group, "id = ?", req.GroupId).Error; err != nil {
				return errors.New("群组不存在")
			}
			// 将目标设为组长时
			if data.IsLeader {
				// 将所有成员目标设为非组长
				if err := tx.Model(model.TargetPersonal{}).Where(&model.TargetPersonal{
					TargetId: target.Id,
					GroupId:  group.Id,
				}).Update("is_leader", false).Error; err != nil {
					return errors.New("将所有成员设为非组长失败")
				}
			}
		}

		// 个人不能重复
		for _, personal := range target.Personals {
			if personal.StaffId == data.StaffId {
				return errors.New("个人已存在")
			}
		}

		// 创建个人(跳过关联)
		if err := tx.Omit(clause.Associations).Create(&data).Error; err != nil {
			return errors.New("创建个人失败")
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &data, nil
}
