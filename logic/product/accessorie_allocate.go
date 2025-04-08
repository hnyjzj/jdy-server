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

type ProductAccessorieAllocateLogic struct {
	Ctx   *gin.Context
	Staff *types.Staff
}

// 创建配件调拨单
func (l *ProductAccessorieAllocateLogic) Create(req *types.ProductAccessorieAllocateCreateReq) *errors.Errors {
	// 开启事务
	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 创建调拨单
		data := model.ProductAccessorieAllocate{
			Method: req.Method,
			Remark: req.Remark,
			Status: enums.ProductAllocateStatusDraft,

			FromStoreId: req.FromStoreId,
			ToStoreId:   req.ToStoreId,

			OperatorId: l.Staff.Id,
			IP:         l.Ctx.ClientIP(),
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

// 获取配件调拨单列表
func (p *ProductAccessorieAllocateLogic) List(req *types.ProductAccessorieAllocateListReq) (*types.PageRes[model.ProductAccessorieAllocate], error) {
	var (
		allocate model.ProductAccessorieAllocate

		res types.PageRes[model.ProductAccessorieAllocate]
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
	db = db.Preload("Products.Product")

	if err := db.Find(&res.List).Error; err != nil {
		return nil, errors.New("获取调拨单列表失败")
	}

	return &res, nil
}

// 获取配件调拨单详情
func (p *ProductAccessorieAllocateLogic) Info(req *types.ProductAccessorieAllocateInfoReq) (*model.ProductAccessorieAllocate, error) {
	var (
		allocate model.ProductAccessorieAllocate
	)

	db := model.DB.Model(&allocate)

	db = db.Preload("Products.Product.Category")
	db = db.Preload("FromStore")
	db = db.Preload("ToStore")
	db = db.Preload("Operator")

	if err := db.First(&allocate, req.Id).Error; err != nil {
		return nil, errors.New("获取调拨单详情失败")
	}

	return &allocate, nil
}

// 添加配件调拨单产品
func (p *ProductAccessorieAllocateLogic) Add(req *types.ProductAccessorieAllocateAddReq) *errors.Errors {
	var (
		allocate model.ProductAccessorieAllocate
	)

	// 获取调拨单
	if err := model.DB.Preload("Products.Product.Category").First(&allocate, req.Id).Error; err != nil {
		return errors.New("调拨单不存在")
	}

	// 检查调拨单状态
	if allocate.Status != enums.ProductAllocateStatusDraft {
		return errors.New("调拨单状态异常")
	}

	// 所有配件ID
	products := make(map[string]model.ProductAccessorieAllocateProduct)
	for _, p := range allocate.Products {
		products[p.ProductId] = p
	}

	// 添加配件
	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		for _, rp := range req.Products {
			var product model.ProductAccessorie
			// 获取配件
			if err := tx.Where("id = ?", rp.ProductId).First(&product).Error; err != nil {
				return errors.New("配件不存在")
			}
			// 检查配件状态
			if product.Status != enums.ProductStatusNormal {
				return errors.New("配件状态异常")
			}

			// 检查配件是否已存在
			pap, ok := products[rp.ProductId]
			if ok { // 已存在，更新数量
				// 检查配件库存
				if product.Stock < (rp.Quantity + pap.Quantity) {
					return errors.New("配件库存不足")
				}
				// 更新配件数量
				if err := tx.Model(&model.ProductAccessorieAllocateProduct{}).
					Where(&model.ProductAccessorieAllocateProduct{
						AllocateId: req.Id,
						ProductId:  rp.ProductId,
					}).Update("quantity", gorm.Expr("quantity + ?", rp.Quantity)).Error; err != nil {
					return errors.New("更新配件数量失败")
				}
			} else { // 不存在，新增
				// 检查配件库存
				if product.Stock < rp.Quantity {
					return errors.New("配件库存不足")
				}
				data := model.ProductAccessorieAllocateProduct{
					ProductId: rp.ProductId,
					Quantity:  rp.Quantity,
				}
				// 添加配件
				if err := tx.Model(&allocate).Association("Products").Append(&data); err != nil {
					return errors.New("添加配件失败")
				}
			}
		}
		return nil
	}); err != nil {
		return errors.New("添加配件失败")
	}

	return nil
}

// 移除配件调拨单配件
func (p *ProductAccessorieAllocateLogic) Remove(req *types.ProductAccessorieAllocateRemoveReq) *errors.Errors {
	var (
		allocate model.ProductAccessorieAllocate
	)

	// 获取调拨单
	if err := model.DB.First(&allocate, req.Id).Error; err != nil {
		return errors.New("调拨单不存在")
	}

	var product model.ProductAccessorieAllocateProduct
	// 获取配件
	where := &model.ProductAccessorieAllocateProduct{
		ProductId:  req.ProductId,
		AllocateId: req.Id,
	}
	if err := model.DB.Where(where).First(&product).Error; err != nil {
		return errors.New("配件不存在")
	}

	// 移除配件
	if err := model.DB.Where(where).Delete(&product).Error; err != nil {
		return errors.New("移除配件失败")
	}

	return nil
}

// 确认调拨
func (p *ProductAccessorieAllocateLogic) Confirm(req *types.ProductAccessorieAllocateConfirmReq) *errors.Errors {
	var (
		allocate model.ProductAccessorieAllocate
	)

	// 获取调拨单
	if err := model.DB.
		Preload("Products.Product.Category").
		First(&allocate, "id = ?", req.Id).Error; err != nil {
		return errors.New("调拨单不存在")
	}

	if allocate.Status != enums.ProductAllocateStatusDraft {
		return errors.New("调拨单状态异常")
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 扣除配件库存
		for _, p := range allocate.Products {
			var product model.ProductAccessorieAllocateProduct
			if err := tx.Where("product_id = ?", p.ProductId).First(&product).Error; err != nil {
				return fmt.Errorf("【%s】%s 不存在", p.Product.Category.Code, p.Product.Category.Name)
			}
			if err := tx.Model(&product).Update("quantity", gorm.Expr("quantity - ?", p.Quantity)).Error; err != nil {
				return fmt.Errorf("【%s】%s 扣除库存失败", p.Product.Category.Code, p.Product.Category.Name)
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
func (p *ProductAccessorieAllocateLogic) Cancel(req *types.ProductAccessorieAllocateCancelReq) *errors.Errors {
	var (
		allocate model.ProductAccessorieAllocate
	)

	// 获取调拨单
	if err := model.DB.Preload("Products").First(&allocate, req.Id).Error; err != nil {
		return errors.New("调拨单不存在")
	}

	if allocate.Status != enums.ProductAllocateStatusDraft && allocate.Status != enums.ProductAllocateStatusOnTheWay {
		return errors.New("调拨单状态异常")
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 解锁配件
		if allocate.Status == enums.ProductAllocateStatusOnTheWay {
			for _, product := range allocate.Products {
				if err := tx.Model(&product).Update("quantity", gorm.Expr("quantity + ?", product.Quantity)).Error; err != nil {
					return fmt.Errorf("【%s】%s 扣除库存失败", product.Product.Category.Code, product.Product.Category.Name)
				}
			}
		}

		// 取消调拨
		allocate.Status = enums.ProductAllocateStatusCancelled
		if err := tx.Save(&allocate).Error; err != nil {
			return errors.New("更新调拨单失败")
		}

		return nil
	}); err != nil {
		return errors.New("调拨失败: " + err.Error())
	}
	return nil
}

// 完成调拨
func (p *ProductAccessorieAllocateLogic) Complete(req *types.ProductAccessorieAllocateCompleteReq) *errors.Errors {
	var (
		allocate model.ProductAccessorieAllocate
	)

	// 获取调拨单
	db := model.DB.Model(&allocate)
	db = db.Preload("Products.Product.Category", func(tx *gorm.DB) *gorm.DB {
		tx = tx.Preload("Product")
		return tx
	})
	if err := db.First(&allocate, req.Id).Error; err != nil {
		return errors.New("调拨单不存在")
	}

	if allocate.Status != enums.ProductAllocateStatusOnTheWay {
		return errors.New("调拨单状态异常")
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 接收配件
		for _, product := range allocate.Products {

			log := model.ProductHistory{
				OldValue:   product,
				ProductId:  product.ProductId,
				StoreId:    allocate.ToStoreId,
				SourceId:   allocate.Id,
				OperatorId: p.Staff.Id,
				IP:         p.Ctx.ClientIP(),
			}

			if allocate.Method == enums.ProductAllocateMethodStore {
				// 区分记录类型
				log.Action = enums.ProductActionTransfer
				// 查询配件
				var accessorie model.ProductAccessorie
				if err := tx.Where(&model.ProductAccessorie{
					StoreId: allocate.ToStoreId,
					Code:    product.Product.Code,
				}).First(&accessorie).Error; err != nil {
					if err != gorm.ErrRecordNotFound {
						return fmt.Errorf("【%s】%s 不存在", product.Product.Category.Code, product.Product.Category.Name)
					}
				}

				if accessorie.Id == "" {
					// 新增配件
					data := model.ProductAccessorie{
						StoreId:   allocate.ToStoreId,
						Code:      product.Product.Code,
						Stock:     product.Quantity,
						AccessFee: product.Product.AccessFee,
						Status:    enums.ProductStatusNormal,
					}
					if err := tx.Create(&data).Error; err != nil {
						return fmt.Errorf("【%s】%s 新增配件失败", product.Product.Category.Code, product.Product.Category.Name)
					}
				} else {
					// 更新配件
					if err := tx.Model(&accessorie).Where(&model.ProductAccessorie{
						StoreId: allocate.ToStoreId,
						Code:    product.Product.Code,
					}).Updates(&model.ProductAccessorie{
						Stock:     accessorie.Stock + product.Quantity,
						AccessFee: product.Product.AccessFee,
					}).Error; err != nil {
						return fmt.Errorf("【%s】%s 更新配件失败", product.Product.Category.Code, product.Product.Category.Name)
					}
				}
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
