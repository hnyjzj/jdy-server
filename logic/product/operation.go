package product

import (
	"jdy/enums"
	"jdy/errors"
	"jdy/model"
	"jdy/types"

	"gorm.io/gorm"
)

// 产品报损
func (l *ProductLogic) Damage(req *types.ProductDamageReq) *errors.Errors {
	// 查询商品信息
	var product model.Product
	if err := model.DB.Where(&model.Product{Code: req.Code}).
		Preload("Store").
		Preload("RecycleStore").
		First(&product).Error; err != nil {
		return errors.New("商品不存在")
	}

	// 判断产品状态
	if product.Status != enums.ProductStatusNormal {
		return errors.New("产品不在库存中")
	}
	// 判断是否可以报损
	if err := product.Status.CanTransitionTo(enums.ProductStatusDamage); err != nil {
		return errors.New("产品状态不允许报损")
	}

	// 开启事务
	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		log := &model.ProductDamage{
			ProductId:  product.Id,
			OperatorId: l.Staff.Id,
			Reason:     req.Reason,
			IP:         l.Ctx.ClientIP(),
		}

		history := model.ProductHistory{
			Action:     enums.ProductActionDamage,
			OldValue:   product,
			ProductId:  product.Id,
			StoreId:    product.StoreId,
			SourceId:   log.Id,
			OperatorId: l.Staff.Id,
			IP:         l.Ctx.ClientIP(),
		}

		// 添加报损记录
		if err := tx.Create(&log).Error; err != nil {
			return err
		}

		// 更新商品状态
		product.Status = enums.ProductStatusDamage
		if err := tx.Save(&product).Error; err != nil {
			return err
		}

		// 添加记录
		history.NewValue = product
		if err := tx.Create(&history).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return errors.New("报损失败")
	}

	return nil
}

// 产品转换
func (l *ProductLogic) Conversion(req *types.ProductConversionReq) *errors.Errors {
	// 查询商品信息
	var product model.Product
	if err := model.DB.Where(&model.Product{Code: req.Code}).First(&product).Error; err != nil {
		return errors.New("商品不存在")
	}

	// 判断产品状态
	if product.Status != enums.ProductStatusNormal {
		return errors.New("产品不在库存中")
	}

	// 开启事务
	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 更新产品状态
		product.Type = req.Type
		if err := tx.Model(&product).
			Where(&model.Product{Code: req.Code}).
			Updates(model.Product{Type: req.Type}).
			Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return errors.New("转换失败")
	}

	return nil
}
