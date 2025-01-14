package product

import (
	"fmt"
	"jdy/enums"
	"jdy/errors"
	"jdy/model"
	"jdy/types"

	"gorm.io/gorm"
)

type ProductAllocateLogic struct {
	ProductLogic
}

// 创建产品调拨单
func (l *ProductAllocateLogic) Create(req *types.ProductAllocateCreateReq) *errors.Errors {
	// 开启事务
	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 创建调拨单
		data := model.ProductAllocate{
			Method: req.Method,
			Type:   req.Type,
			Reason: req.Reason,
			Remark: req.Remark,
			Status: enums.ProductAllocateStatusInventory,

			FromStoreId: req.FromStoreId,

			OperatorId: l.Staff.Id,
			IP:         l.Ctx.ClientIP(),
		}
		// 如果是调拨到门店，则添加门店ID
		if req.Method == enums.ProductAllocateMethodStore {
			data.ToStoreId = req.ToStoreId
		}
		// 创建调拨单
		if err := tx.Create(&data).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return errors.New("创建调拨单失败")
	}
	return nil
}

// 获取产品调拨单列表
func (p *ProductAllocateLogic) List(req *types.ProductAllocateListReq) (*types.PageRes[model.ProductAllocate], error) {
	var (
		allocate model.ProductAllocate

		res types.PageRes[model.ProductAllocate]
	)

	db := model.DB.Model(&allocate)
	db = allocate.WhereCondition(db, &req.Where)

	// 获取总数
	if err := db.Count(&res.Total).Error; err != nil {
		return nil, errors.New("获取调拨单数量失败")
	}

	// 获取列表
	db = db.Order("created_at desc")
	db = model.PageCondition(db, req.Page, req.Limit)
	if err := db.Find(&res.List).Error; err != nil {
		return nil, errors.New("获取调拨单列表失败")
	}

	return &res, nil
}

// 获取产品调拨单详情
func (p *ProductAllocateLogic) Info(req *types.ProductAllocateInfoReq) (*model.ProductAllocate, error) {
	var (
		allocate model.ProductAllocate
	)

	db := model.DB.Model(&allocate)

	db = db.Preload("Products")
	db = db.Preload("Operator")
	db = db.Preload("Store")

	if err := db.First(&allocate, req.Id).Error; err != nil {
		return nil, errors.New("获取调拨单详情失败")
	}

	return &allocate, nil
}

// 添加产品调拨单产品
func (p *ProductAllocateLogic) Add(req *types.ProductAllocateAddReq) *errors.Errors {
	var (
		allocate model.ProductAllocate
		product  model.Product
	)

	// 获取调拨单
	if err := model.DB.First(&allocate, req.Id).Error; err != nil {
		return errors.New("调拨单不存在")
	}

	// 获取产品
	if err := model.DB.Where(&model.Product{Code: req.Code}).First(&product).Error; err != nil {
		return errors.New("产品不存在")
	}

	if product.Status != enums.ProductStatusNormal {
		return errors.New("产品状态不正确")
	}

	// 添加产品
	if err := model.DB.Model(&allocate).Association("Products").Append(&product); err != nil {
		return errors.New("添加产品失败")
	}

	return nil
}

// 移除产品调拨单产品
func (p *ProductAllocateLogic) Remove(req *types.ProductAllocateRemoveReq) *errors.Errors {
	var (
		allocate model.ProductAllocate
		product  model.Product
	)

	// 获取调拨单
	if err := model.DB.First(&allocate, req.Id).Error; err != nil {
		return errors.New("调拨单不存在")
	}

	// 获取产品
	if err := model.DB.Where(&model.Product{Code: req.Code}).First(&product).Error; err != nil {
		return errors.New("产品不存在")
	}

	// 移除产品
	if err := model.DB.Model(&allocate).Association("Products").Delete(&product); err != nil {
		return errors.New("移除产品失败")
	}

	return nil
}

// 确认调拨
func (p *ProductAllocateLogic) Confirm(req *types.ProductAllocateConfirmReq) *errors.Errors {
	var (
		allocate model.ProductAllocate
	)

	// 获取调拨单
	if err := model.DB.Preload("Products").First(&allocate, req.Id).Error; err != nil {
		return errors.New("调拨单不存在")
	}

	if allocate.Status != enums.ProductAllocateStatusInventory {
		return errors.New("调拨单状态异常")
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 锁定产品
		for _, product := range allocate.Products {
			if product.Status != enums.ProductStatusNormal {
				return errors.New(fmt.Sprintf("【%s】%s 状态异常", product.Code, product.Name))
			}
			if err := tx.Model(&product).Update("status", enums.ProductStatusAllocate).Error; err != nil {
				return errors.New(fmt.Sprintf("【%s】%s 锁定失败", product.Code, product.Name))
			}
		}
		// 确认调拨
		allocate.Status = enums.ProductAllocateStatusAllocate
		if err := model.DB.Save(&allocate).Error; err != nil {
			return errors.New("更新调拨单失败")
		}

		return nil
	}); err != nil {
		return errors.New("调拨失败: " + err.Error())
	}

	return nil
}

// 取消调拨
func (p *ProductAllocateLogic) Cancel(req *types.ProductAllocateCancelReq) *errors.Errors {
	var (
		allocate model.ProductAllocate
	)

	// 获取调拨单
	if err := model.DB.Preload("Products").First(&allocate, req.Id).Error; err != nil {
		return errors.New("调拨单不存在")
	}

	if allocate.Status != enums.ProductAllocateStatusInventory && allocate.Status != enums.ProductAllocateStatusAllocate {
		return errors.New("调拨单状态异常")
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 取消调拨
		allocate.Status = enums.ProductAllocateStatusCanceled
		if err := model.DB.Save(&allocate).Error; err != nil {
			return errors.New("更新调拨单失败")
		}

		// 解锁产品
		for _, product := range allocate.Products {
			if product.Status != enums.ProductStatusAllocate {
				fmt.Printf("调拨单产品状态异常：【%s】%s\n", product.Code, product.Name)
				break
			}
			if err := model.DB.Model(&product).Update("status", enums.ProductStatusNormal).Error; err != nil {
				return errors.New(fmt.Sprintf("【%s】%s 解锁失败", product.Code, product.Name))
			}
		}
		return nil
	}); err != nil {
		return errors.New("调拨失败: " + err.Error())
	}
	return nil
}

// 完成调拨
func (p *ProductAllocateLogic) Complete(req *types.ProductAllocateCompleteReq) *errors.Errors {
	var (
		allocate model.ProductAllocate
	)

	// 获取调拨单
	if err := model.DB.Preload("Products").First(&allocate, req.Id).Error; err != nil {
		return errors.New("调拨单不存在")
	}

	if allocate.Status != enums.ProductAllocateStatusAllocate {
		return errors.New("调拨单状态异常")
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 解锁产品
		for _, product := range allocate.Products {
			// 判断产品状态是否为锁定状态
			if product.Status != enums.ProductStatusAllocate {
				return errors.New(fmt.Sprintf("【%s】%s 状态异常", product.Code, product.Name))
			}

			data := &model.Product{
				Status:  enums.ProductStatusNormal,
				StoreId: allocate.ToStoreId,
				Type:    allocate.Type,
			}

			// 解锁产品
			if err := model.DB.Model(&product).Updates(data).Error; err != nil {
				return errors.New(fmt.Sprintf("【%s】%s 解锁失败", product.Code, product.Name))
			}
		}

		// 确认调拨
		allocate.Status = enums.ProductAllocateStatusCompleted
		if err := model.DB.Save(&allocate).Error; err != nil {
			return errors.New("更新调拨单失败")
		}

		return nil
	}); err != nil {
		return errors.New("调拨失败: " + err.Error())
	}

	return nil
}
