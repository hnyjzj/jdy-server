package product

import (
	"jdy/errors"
	"jdy/model"
	"jdy/types"
	"jdy/utils"

	"gorm.io/gorm"
)

// 产品入库
func (l *ProductLogic) Enter(req *types.ProductEnterReq) (*map[string]bool, *errors.Errors) {
	// 添加产品的结果
	products := map[string]bool{}
	// 添加产品入库
	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 转换数据结构
		data, err := utils.StructToStruct[[]model.Product](req.Products)
		if err != nil {
			return nil
		}
		if len(data) == 0 {
			return errors.New("产品录入失败")
		}

		enter := model.ProductEnter{
			OperatorId: l.Staff.Id,
		}
		if err := tx.Create(&enter).Error; err != nil {
			return err
		}

		for i, v := range data {
			data[i].ProductEnterId = enter.Id
			products[v.Code] = false

			var p model.Product
			if err := tx.Where("code = ?", v.Code).First(&p).Error; err != nil {
				if err != gorm.ErrRecordNotFound {
					return err
				}
			}
			if p.Id != "" {
				continue
			}
			if err := tx.Create(&v).Error; err != nil {
				continue
			}

			products[v.Code] = true
		}
		return nil
	}); err != nil {
		return nil, errors.New("产品录入失败")
	}

	return &products, nil
}
