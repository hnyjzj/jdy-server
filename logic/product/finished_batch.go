package product

import (
	"jdy/enums"
	"jdy/errors"
	"jdy/model"
	"jdy/types"
	"jdy/utils"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProductFinishedBatchLogic struct {
	ProductFinishedLogic
}

// 批量更新成品信息
func (p *ProductFinishedBatchLogic) Code(req *types.ProductFinishedUpdateCodeReq) error {
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
				return errors.New("添加记录失败")
			}
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

// 批量更新成品信息
func (p *ProductFinishedBatchLogic) Update(req *types.ProductFinishedUpdatesReq) error {
	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		for _, r := range req.Data {
			// 验证参数
			data, err := utils.StructToStruct[model.ProductFinished](r)
			if err != nil {
				return errors.New("验证参数失败")
			}

			var product model.ProductFinished
			if err := tx.Model(&model.ProductFinished{}).
				Clauses(clause.Locking{Strength: "UPDATE"}).
				Where(&model.ProductFinished{
					Code: data.Code,
				}).First(&product).Error; err != nil {
				return errors.New("获取【" + data.Code + "】条码信息失败，可能不存在")
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

			data.Class = data.GetClass()

			if err := tx.Model(&model.ProductFinished{}).Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", product.Id).Omit(
				"id", "created_at", "updated_at", "deleted_at",
				"code", "status", "images", "store_id", "enter_id",
			).Updates(&data).Error; err != nil {
				return errors.New("更新【" + data.Code + "】信息失败")
			}

			// 添加记录
			history.NewValue = data
			if err := tx.Create(&history).Error; err != nil {
				return errors.New("添加记录失败")
			}

		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
