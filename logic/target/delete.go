package target

import (
	"errors"
	"jdy/model"
	"jdy/types"
)

func (l *Logic) Delete(req *types.TargetDeleteReq) error {
	var target model.Target
	if err := model.DB.Where("id = ?", req.Id).First(&target).Error; err != nil {
		return errors.New("目标不存在")
	}

	// 删除目标
	if err := model.DB.Where("id = ?", req.Id).Delete(&target).Error; err != nil {
		return errors.New("删除目标失败")
	}

	return nil
}

func (l *Logic) DeleteGroup(req *types.TargetDeleteGroupReq) error {
	var (
		group model.TargetGroup
	)
	db := model.DB.Where("id = ?", req.Id)
	db = group.Preloads(db)
	if err := db.First(&group).Error; err != nil {
		return errors.New("分组不存在")
	}

	if len(group.Personals) > 0 {
		return errors.New("该分组下有个人目标，无法删除")
	}

	if err := model.DB.Where("id = ?", req.Id).Delete(&group).Error; err != nil {
		return errors.New("删除目标失败")
	}

	return nil
}

func (l *Logic) DeletePersonal(req *types.TargetDeletePersonalReq) error {
	var (
		personal model.TargetPersonal
	)

	db := model.DB.Where("id = ?", req.Id)
	if err := db.First(&personal).Error; err != nil {
		return errors.New("个人目标不存在")
	}

	if personal.Achieve.IsPositive() {
		return errors.New("有业绩，无法删除")
	}

	if err := model.DB.Where("id = ?", req.Id).Delete(&personal).Error; err != nil {
		return errors.New("删除个人目标失败")
	}

	return nil
}
