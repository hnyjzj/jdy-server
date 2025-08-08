package product

import (
	"jdy/enums"
	"jdy/errors"
	"jdy/model"
	"jdy/types"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// 更新成品信息
func (p *ProductFinishedLogic) Code(req *types.ProductFinishedUpdateCodeReq) error {
	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 验证成品信息
		for _, r := range req.Data {
			var product model.ProductFinished
			if err := tx.Model(&model.ProductFinished{}).Clauses(clause.Locking{Strength: "UPDATE"}).Where(&model.ProductFinished{
				Code: r.Code,
			}).First(&product).Error; err != nil {
				return errors.New("查找[" + r.Code + "]信息失败")
			}

			history := model.ProductHistory{
				Type:       enums.ProductTypeFinished,
				Action:     enums.ProductActionUpdate,
				OldValue:   product,
				ProductId:  product.Id,
				SourceId:   product.Id,
				StoreId:    product.StoreId,
				OperatorId: p.Staff.Id,
				IP:         p.Ctx.ClientIP(),
			}

			product.Code = r.NewCode
			if err := tx.Model(&model.ProductFinished{}).Where("id = ?", product.Id).Update("code", product.Code).Error; err != nil {
				return errors.New("更新[" + r.Code + "]为[" + r.NewCode + "]失败")
			}

			// 添加记录
			history.NewValue = product
			if err := tx.Create(&history).Error; err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
