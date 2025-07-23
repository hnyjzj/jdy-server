package setting

import (
	"errors"
	"jdy/logic"
	"jdy/model"
	"jdy/types"
	"jdy/utils"

	"gorm.io/gorm"
)

type OpenOrderLogic struct {
	logic.BaseLogic
}

func (OpenOrderLogic) Info(req *types.OpenOrderInfoReq) (*model.OpenOrder, error) {
	var (
		open_order model.OpenOrder
	)

	if err := model.DB.Where(model.OpenOrder{
		StoreId: req.StoreId,
	}).First(&open_order).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, errors.New("获取失败")
		}
	}
	if open_order.Id == "" {
		return open_order.Default(), nil
	}

	return &open_order, nil
}

func (OpenOrderLogic) Update(req *types.OpenOrderUpdateReq) error {

	var open_order model.OpenOrder
	if err := model.DB.Model(&model.OpenOrder{}).Where(model.OpenOrder{
		StoreId: req.StoreId,
	}).First(&open_order).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return errors.New("查询失败")
		}
	}

	if open_order.Id == "" {
		// 数据转换
		data, err := utils.StructToStruct[model.OpenOrder](req)
		if err != nil {
			return errors.New("数据转换失败")
		}
		if err := model.DB.Create(&data).Error; err != nil {
			return errors.New("创建失败")
		}
	} else {
		if err := model.DB.Model(&model.OpenOrder{}).Where("id = ?", open_order.Id).Select("*").Updates(req).Error; err != nil {
			return errors.New("更新失败")
		}
	}

	return nil
}
