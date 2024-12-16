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
		product.Status = enums.ProductStatusDamage
		if err := tx.Save(&product).Error; err != nil {
			return err
		}

		// 添加报损记录
		if err := tx.Create(&model.ProductDamage{
			ProductId:  product.Id,
			OperatorId: l.Staff.Id,
			Reason:     req.Reason,
			IP:         l.Ctx.ClientIP(),
		}).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return errors.New("报损失败")
	}

	return nil
}

// 产品调拨
func (l *ProductLogic) Allocate(req *types.ProductAllocateReq) *errors.Errors {
	// 开启事务
	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		data := model.ProductAllocate{
			Method: req.Method,
			Type:   req.Type,
			Reason: req.Reason,
			Remark: req.Remark,

			OperatorId: l.Staff.Id,
			IP:         l.Ctx.ClientIP(),
		}
		if req.Method == enums.ProductAllocateMethodStore {
			data.StoreId = req.StoreId
		}
		// 添加报损记录
		if err := tx.Create(&data).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return errors.New("报损失败")
	}
	return nil
}
