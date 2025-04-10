package product

import (
	"fmt"
	"jdy/enums"
	"jdy/errors"
	"jdy/model"
	"jdy/types"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductOldAllocateLogic struct {
	Ctx   *gin.Context
	Staff *types.Staff
}

// 创建产品调拨单
func (l *ProductOldAllocateLogic) Create(req *types.ProductOldAllocateCreateReq) *errors.Errors {
	// 开启事务
	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 创建调拨单
		data := model.ProductOldAllocate{
			Method: req.Method,
			Type:   req.Type,
			Reason: req.Reason,
			Remark: req.Remark,
			Status: enums.ProductAllocateStatusDraft,

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
func (p *ProductOldAllocateLogic) List(req *types.ProductOldAllocateListReq) (*types.PageRes[model.ProductOldAllocate], error) {
	var (
		allocate model.ProductOldAllocate

		res types.PageRes[model.ProductOldAllocate]
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
func (p *ProductOldAllocateLogic) Info(req *types.ProductOldAllocateInfoReq) (*model.ProductOldAllocate, error) {
	var (
		allocate model.ProductOldAllocate
	)

	db := model.DB.Model(&allocate)

	db = db.Preload("Products")
	db = db.Preload("Operator")
	db = db.Preload("FromStore")
	db = db.Preload("ToStore")

	if err := db.First(&allocate, req.Id).Error; err != nil {
		return nil, errors.New("获取调拨单详情失败")
	}

	return &allocate, nil
}

// 添加产品调拨单产品
func (p *ProductOldAllocateLogic) Add(req *types.ProductOldAllocateAddReq) *errors.Errors {
	var (
		allocate model.ProductOldAllocate
	)

	// 获取调拨单
	if err := model.DB.First(&allocate, req.Id).Error; err != nil {
		return errors.New("调拨单不存在")
	}

	var product model.ProductOld
	// 获取产品
	if err := model.DB.Where(&model.ProductOld{Code: req.Code}).First(&product).Error; err != nil {
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
func (p *ProductOldAllocateLogic) Remove(req *types.ProductOldAllocateRemoveReq) *errors.Errors {
	var (
		allocate model.ProductOldAllocate
	)

	// 获取调拨单
	if err := model.DB.First(&allocate, req.Id).Error; err != nil {
		return errors.New("调拨单不存在")
	}

	var product model.ProductOld
	// 获取产品
	if err := model.DB.Where(&model.ProductOld{Code: req.Code}).First(&product).Error; err != nil {
		return errors.New("产品不存在")
	}

	// 移除产品
	if err := model.DB.Model(&allocate).Association("Products").Delete(&product); err != nil {
		return errors.New("移除产品失败")
	}

	return nil
}

// 确认调拨
func (p *ProductOldAllocateLogic) Confirm(req *types.ProductOldAllocateConfirmReq) *errors.Errors {
	var (
		allocate model.ProductOldAllocate
	)

	// 获取调拨单
	if err := model.DB.Preload("Products").First(&allocate, req.Id).Error; err != nil {
		return errors.New("调拨单不存在")
	}

	if allocate.Status != enums.ProductAllocateStatusDraft {
		return errors.New("调拨单状态异常")
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 锁定产品
		for _, product := range allocate.Products {
			if product.Status != enums.ProductStatusNormal {
				return errors.New(fmt.Sprintf("【%s】%s 状态异常", product.Code, product.Name))
			}
			product.Status = enums.ProductStatusAllocate
			if err := tx.Save(&product).Error; err != nil {
				return errors.New(fmt.Sprintf("【%s】%s 锁定失败", product.Code, product.Name))
			}
		}
		// 确认调拨
		allocate.Status = enums.ProductAllocateStatusOnTheWay
		if err := tx.Save(&allocate).Error; err != nil {
			return errors.New("更新调拨单失败")
		}

		return nil
	}); err != nil {
		return errors.New("调拨失败: " + err.Error())
	}

	return nil
}

// 取消调拨
func (p *ProductOldAllocateLogic) Cancel(req *types.ProductOldAllocateCancelReq) *errors.Errors {
	var (
		allocate model.ProductOldAllocate
	)

	// 获取调拨单
	if err := model.DB.Preload("Products").First(&allocate, req.Id).Error; err != nil {
		return errors.New("调拨单不存在")
	}

	if allocate.Status != enums.ProductAllocateStatusDraft && allocate.Status != enums.ProductAllocateStatusOnTheWay {
		return errors.New("调拨单状态异常")
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 取消调拨
		allocate.Status = enums.ProductAllocateStatusCancelled
		if err := tx.Save(&allocate).Error; err != nil {
			return errors.New("更新调拨单失败")
		}

		// 解锁产品
		for _, product := range allocate.Products {
			if product.Status != enums.ProductStatusAllocate {
				return fmt.Errorf("【%s】%s 状态异常", product.Code, product.Name)
			}
			if err := tx.Model(&product).Update("status", enums.ProductStatusNormal).Error; err != nil {
				return fmt.Errorf("【%s】%s 解锁失败", product.Code, product.Name)
			}
		}
		return nil
	}); err != nil {
		return errors.New("调拨失败: " + err.Error())
	}
	return nil
}

// 完成调拨
func (p *ProductOldAllocateLogic) Complete(req *types.ProductOldAllocateCompleteReq) *errors.Errors {
	var (
		allocate model.ProductOldAllocate
	)

	// 获取调拨单
	db := model.DB.Model(&allocate)
	db = db.Preload("Products", func(tx *gorm.DB) *gorm.DB {
		tx = tx.Preload("Store")
		return tx
	})
	if err := db.First(&allocate, req.Id).Error; err != nil {
		return errors.New("调拨单不存在")
	}

	if allocate.Status != enums.ProductAllocateStatusOnTheWay {
		return errors.New("调拨单状态异常")
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {

		// 解锁产品
		for _, product := range allocate.Products {
			log := model.ProductHistory{
				Type:       enums.ProductTypeOld,
				Action:     enums.ProductActionTransfer,
				OldValue:   product,
				ProductId:  product.Id,
				StoreId:    product.StoreId,
				SourceId:   allocate.Id,
				OperatorId: p.Staff.Id,
				IP:         p.Ctx.ClientIP(),
			}

			// 判断产品状态是否为锁定状态
			if product.Status != enums.ProductStatusAllocate {
				return errors.New(fmt.Sprintf("【%s】%s 状态异常", product.Code, product.Name))
			}

			data := &model.ProductOld{
				StoreId: allocate.ToStoreId,
			}

			// 解锁产品
			if err := tx.Model(&model.ProductOld{}).Where("id = ?", product.Id).Updates(data).Error; err != nil {
				return errors.New(fmt.Sprintf("【%s】%s 解锁失败", product.Code, product.Name))
			}

			// 更新商品状态
			product.Status = enums.ProductStatusNormal
			product.StoreId = allocate.ToStoreId
			if err := tx.Save(&product).Error; err != nil {
				return err
			}

			// 添加记录
			log.NewValue = product
			if err := tx.Create(&log).Error; err != nil {
				return err
			}
		}
		// 确认调拨
		allocate.Status = enums.ProductAllocateStatusCompleted
		if err := tx.Save(&allocate).Error; err != nil {
			return errors.New("更新调拨单失败")
		}

		return nil
	}); err != nil {
		return errors.New("调拨失败: " + err.Error())
	}

	return nil
}
