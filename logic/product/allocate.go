package product

import (
	"fmt"
	"jdy/enums"
	"jdy/errors"
	"jdy/message"
	"jdy/model"
	"jdy/types"
	"jdy/utils"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductAllocateLogic struct {
	Ctx   *gin.Context
	Staff *model.Staff
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
			Status: enums.ProductAllocateStatusDraft,

			FromStoreId: req.FromStoreId,

			OperatorId: l.Staff.Id,
			IP:         l.Ctx.ClientIP(),
		}
		// 如果是调拨到门店，则添加门店ID
		if req.Method == enums.ProductAllocateMethodStore {
			data.ToStoreId = req.ToStoreId
		}
		if req.Method == enums.ProductAllocateMethodOut {
			var store model.Store
			if err := tx.Where(&model.Store{
				Name: "总部",
			}).First(&store).Error; err != nil {
				return err
			}
			data.ToStoreId = store.Id
		}

		// 创建调拨单
		if err := tx.Create(&data).Error; err != nil {
			return err
		}

		// 判断是不是成品整单调拨
		if req.EnterId != "" && req.Type == enums.ProductTypeFinished {
			// 获取产品
			var enter model.ProductFinishedEnter
			if err := tx.Preload("Products", func(tx *gorm.DB) *gorm.DB {
				return tx.Where(&model.ProductFinished{
					StoreId: req.FromStoreId,
					Status:  enums.ProductStatusNormal,
				})
			}).Where("id = ?", req.EnterId).First(&enter).Error; err != nil {
				return errors.New("获取入库单失败")
			}
			var allocate model.ProductAllocate
			for _, p := range enter.Products {
				if p.Status != enums.ProductStatusNormal {
					return errors.New("入库单产品状态异常")
				}
				allocate.ProductCount++
				allocate.ProductTotalWeightMetal = allocate.ProductTotalWeightMetal.Add(p.WeightMetal)
				allocate.ProductTotalLabelPrice = allocate.ProductTotalLabelPrice.Add(p.LabelPrice)
				allocate.ProductTotalAccessFee = allocate.ProductTotalAccessFee.Add(p.AccessFee)
			}

			// 更新调拨单
			if err := tx.Model(&model.ProductAllocate{}).Where("id = ?", data.Id).Updates(&allocate).Error; err != nil {
				return err
			}

			// 添加产品
			for _, p := range enter.Products {
				if err := tx.Table("product_allocate_finished_products").Create(map[string]any{
					"product_allocate_id": data.Id,
					"product_finished_id": p.Id,
				}).Error; err != nil {
					return err
				}
			}
		}

		return nil
	}); err != nil {
		log.Printf("创建调拨单失败: %v", err)
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
	db = allocate.Preloads(db)
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

	db = allocate.Preloads(db)

	db = db.Preload("ProductFinisheds", func(tx *gorm.DB) *gorm.DB {
		tx = model.PageCondition(tx, req.Page, req.Limit)
		return tx
	})
	db = db.Preload("ProductOlds", func(tx *gorm.DB) *gorm.DB {
		tx = model.PageCondition(tx, req.Page, req.Limit)
		return tx
	})

	if err := db.First(&allocate, "id = ?", req.Id).Error; err != nil {
		return nil, errors.New("获取调拨单详情失败")
	}

	return &allocate, nil
}

// 添加产品调拨单产品
func (p *ProductAllocateLogic) Add(req *types.ProductAllocateAddReq) *errors.Errors {
	var (
		allocate model.ProductAllocate
	)

	// 获取调拨单
	if err := model.DB.First(&allocate, "id = ?", req.Id).Error; err != nil {
		return errors.New("调拨单不存在")
	}

	if allocate.Status != enums.ProductAllocateStatusDraft {
		return errors.New("调拨单状态异常")
	}

	data := model.ProductAllocate{
		ProductCount:            allocate.ProductCount,
		ProductTotalWeightMetal: allocate.ProductTotalWeightMetal,
		ProductTotalLabelPrice:  allocate.ProductTotalLabelPrice,
		ProductTotalAccessFee:   allocate.ProductTotalAccessFee,
	}
	// 添加产品
	if err := model.DB.Transaction(func(tx *gorm.DB) error {

		switch allocate.Type {
		case enums.ProductTypeFinished:
			var product []model.ProductFinished
			// 获取产品
			if err := tx.Where("id in (?)", req.ProductId).Find(&product).Error; err != nil {
				return errors.New("产品不存在")
			}

			for _, p := range product {
				if p.Status != enums.ProductStatusNormal {
					return errors.New("产品状态不正确")
				}

				data.ProductCount++
				data.ProductTotalWeightMetal = data.ProductTotalWeightMetal.Add(p.WeightMetal)
				data.ProductTotalLabelPrice = data.ProductTotalLabelPrice.Add(p.LabelPrice)
				data.ProductTotalAccessFee = data.ProductTotalAccessFee.Add(p.AccessFee)
			}
			// 添加产品
			for _, p := range product {
				if err := tx.Table("product_allocate_finished_products").Create(map[string]any{
					"product_allocate_id": data.Id,
					"product_finished_id": p.Id,
				}).Error; err != nil {
					return err
				}
			}
		case enums.ProductTypeOld:
			var product []model.ProductOld
			// 获取产品
			if err := tx.Where("id in (?)", req.ProductId).Find(&product).Error; err != nil {
				return errors.New("产品不存在")
			}

			for _, p := range product {
				if p.Status != enums.ProductStatusNormal {
					return errors.New("产品状态不正确")
				}

				data.ProductCount++
				data.ProductTotalWeightMetal = data.ProductTotalWeightMetal.Add(p.WeightMetal)
				data.ProductTotalLabelPrice = data.ProductTotalLabelPrice.Add(p.LabelPrice)
			}

			// 添加产品
			for _, p := range product {
				if err := tx.Table("product_allocate_old_products").Create(map[string]any{
					"product_allocate_id": data.Id,
					"product_old_id":      p.Id,
				}).Error; err != nil {
					return err
				}
			}
		}

		// 更新调拨单
		if err := tx.Model(&model.ProductAllocate{}).Where("id = ?", allocate.Id).Updates(&data).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return errors.New("创建调拨单失败")
	}

	return nil
}

// 移除产品调拨单产品
func (p *ProductAllocateLogic) Remove(req *types.ProductAllocateRemoveReq) *errors.Errors {
	var (
		allocate model.ProductAllocate
	)

	// 获取调拨单
	if err := model.DB.First(&allocate, "id = ?", req.Id).Error; err != nil {
		return errors.New("调拨单不存在")
	}

	if allocate.Status != enums.ProductAllocateStatusDraft {
		return errors.New("调拨单状态异常")
	}

	data := model.ProductAllocate{
		ProductCount:            allocate.ProductCount,
		ProductTotalWeightMetal: allocate.ProductTotalWeightMetal,
		ProductTotalLabelPrice:  allocate.ProductTotalLabelPrice,
		ProductTotalAccessFee:   allocate.ProductTotalAccessFee,
	}
	// 添加产品
	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		switch allocate.Type {
		case enums.ProductTypeFinished:
			var product model.ProductFinished
			// 获取产品
			if err := tx.First(&product, "id = ?", req.ProductId).Error; err != nil {
				return errors.New("产品不存在")
			}

			data.ProductCount--
			data.ProductTotalWeightMetal = data.ProductTotalWeightMetal.Sub(product.WeightMetal)
			data.ProductTotalLabelPrice = data.ProductTotalLabelPrice.Sub(product.LabelPrice)
			data.ProductTotalAccessFee = data.ProductTotalAccessFee.Sub(product.AccessFee)

			// 移除产品
			if err := tx.Model(&allocate).Association("ProductFinisheds").Delete(&product); err != nil {
				return errors.New("移除产品失败")
			}
		case enums.ProductTypeOld:
			var product model.ProductOld
			// 获取产品
			if err := tx.First(&product, "id = ?", req.ProductId).Error; err != nil {
				return errors.New("产品不存在")
			}

			data.ProductCount--
			data.ProductTotalWeightMetal = data.ProductTotalWeightMetal.Sub(product.WeightMetal)
			data.ProductTotalLabelPrice = data.ProductTotalLabelPrice.Sub(product.LabelPrice)

			// 移除产品
			if err := tx.Model(&allocate).Association("ProductOlds").Delete(&product); err != nil {
				return errors.New("移除产品失败")
			}
		}

		// 更新调拨单
		if err := tx.Model(&model.ProductAllocate{}).Where("id = ?", allocate.Id).Updates(&data).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return errors.New("创建调拨单失败")
	}

	return nil
}

// 确认调拨
func (p *ProductAllocateLogic) Confirm(req *types.ProductAllocateConfirmReq) *errors.Errors {
	var (
		allocate model.ProductAllocate
	)

	// 获取调拨单
	db := model.DB
	db = allocate.Preloads(db)
	db = db.Preload("ProductFinisheds")
	db = db.Preload("ProductOlds")
	if err := db.First(&allocate, "id = ?", req.Id).Error; err != nil {
		return errors.New("调拨单不存在")
	}

	if allocate.Status != enums.ProductAllocateStatusDraft {
		return errors.New("调拨单状态异常")
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 锁定产品
		for _, product := range allocate.ProductFinisheds {
			if product.Status != enums.ProductStatusNormal {
				return errors.New(fmt.Sprintf("【%s】%s 状态异常", product.Code, product.Name))
			}
			product.Status = enums.ProductStatusAllocate
			product.EnterTime = time.Now()
			if err := tx.Save(&product).Error; err != nil {
				return errors.New(fmt.Sprintf("【%s】%s 锁定失败", product.Code, product.Name))
			}
		}
		for _, product := range allocate.ProductOlds {
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

	go func() {
		msg := message.NewMessage(p.Ctx)
		allocate.Operator = p.Staff
		msg.SendProductAllocateCreateMessage(&message.ProductAllocateMessage{
			ProductAllocate: &allocate,
		})
	}()

	return nil
}

// 取消调拨
func (p *ProductAllocateLogic) Cancel(req *types.ProductAllocateCancelReq) *errors.Errors {
	var (
		allocate model.ProductAllocate
	)

	// 获取调拨单
	db := model.DB
	db = allocate.Preloads(db)
	db = db.Preload("ProductFinisheds")
	db = db.Preload("ProductOlds")
	if err := db.First(&allocate, "id = ?", req.Id).Error; err != nil {
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
		for _, product := range allocate.ProductFinisheds {
			if product.Status == enums.ProductStatusNormal {
				continue
			}
			if product.Status != enums.ProductStatusAllocate {
				return fmt.Errorf("【%s】%s 状态异常", product.Code, product.Name)
			}
			if err := tx.Model(&product).Where("id = ?", product.Id).Update("status", enums.ProductStatusNormal).Error; err != nil {
				return fmt.Errorf("【%s】%s 解锁失败", product.Code, product.Name)
			}
		}
		for _, product := range allocate.ProductOlds {
			if product.Status == enums.ProductStatusNormal {
				continue
			}
			if product.Status != enums.ProductStatusAllocate {
				return fmt.Errorf("【%s】%s 状态异常", product.Code, product.Name)
			}
			if err := tx.Model(&product).Where("id = ?", product.Id).Update("status", enums.ProductStatusNormal).Error; err != nil {
				return fmt.Errorf("【%s】%s 解锁失败", product.Code, product.Name)
			}
		}
		return nil
	}); err != nil {
		return errors.New("调拨失败: " + err.Error())
	}

	go func() {
		msg := message.NewMessage(p.Ctx)
		allocate.Operator = p.Staff
		msg.SendProductAllocateCancelMessage(&message.ProductAllocateMessage{
			ProductAllocate: &allocate,
		})
	}()

	return nil
}

// 完成调拨
func (p *ProductAllocateLogic) Complete(req *types.ProductAllocateCompleteReq) *errors.Errors {
	var (
		allocate model.ProductAllocate
	)

	// 获取调拨单
	db := model.DB
	db = allocate.Preloads(db)
	db = db.Preload("ProductFinisheds")
	db = db.Preload("ProductOlds")
	if err := db.First(&allocate, "id = ?", req.Id).Error; err != nil {
		return errors.New("调拨单不存在")
	}

	if allocate.Status != enums.ProductAllocateStatusOnTheWay {
		return errors.New("调拨单状态异常")
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 解锁产品
		for _, product := range allocate.ProductFinisheds {
			log := model.ProductHistory{
				Action:     enums.ProductActionTransfer,
				Type:       enums.ProductTypeFinished,
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

			data := model.ProductFinished{
				Status:  enums.ProductStatusNormal,
				StoreId: allocate.ToStoreId,
			}

			// 解锁产品
			if err := tx.Model(&model.ProductFinished{}).Where("id = ?", product.Id).Updates(&data).Error; err != nil {
				return errors.New(fmt.Sprintf("【%s】%s 解锁失败", product.Code, product.Name))
			}

			// 添加记录
			if err := utils.StructMerge(&product, data); err != nil {
				return errors.New(fmt.Sprintf("【%s】%s 更新失败", product.Code, product.Name))
			}
			log.NewValue = product
			if err := tx.Create(&log).Error; err != nil {
				return err
			}
		}

		for _, product := range allocate.ProductOlds {
			log := model.ProductHistory{
				Action:     enums.ProductActionTransfer,
				Type:       enums.ProductTypeOld,
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

			data := model.ProductOld{
				Status:  enums.ProductStatusNormal,
				StoreId: allocate.ToStoreId,
			}

			// 解锁产品
			if err := tx.Model(&model.ProductOld{}).Where("id = ?", product.Id).Updates(&data).Error; err != nil {
				return errors.New(fmt.Sprintf("【%s】%s 解锁失败", product.Code, product.Name))
			}

			// 添加记录
			if err := utils.StructMerge(&product, data); err != nil {
				return errors.New(fmt.Sprintf("【%s】%s 更新失败", product.Code, product.Name))
			}
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

	go func() {
		msg := message.NewMessage(p.Ctx)
		allocate.Operator = p.Staff
		msg.SendProductAllocateCompleteMessage(&message.ProductAllocateMessage{
			ProductAllocate: &allocate,
		})
	}()

	return nil
}
