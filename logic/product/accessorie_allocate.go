package product

import (
	"fmt"
	"jdy/enums"
	"jdy/errors"
	"jdy/message"
	"jdy/model"
	"jdy/types"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductAccessorieAllocateLogic struct {
	Ctx   *gin.Context
	Staff *model.Staff
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

			InitiatorId: l.Staff.Id,
			InitiatorIP: l.Ctx.ClientIP(),
		}

		// 如果是调拨到门店，则添加门店ID
		if req.Method == enums.ProductAccessorieAllocateMethodStore {
			data.ToStoreId = req.ToStoreId
		}
		if req.Method == enums.ProductAccessorieAllocateMethodOut {
			store, err := model.Store{}.Headquarters()
			if err != nil {
				return err
			}
			data.ToStoreId = store.Id
		}
		if req.Method == enums.ProductAccessorieAllocateMethodRegion {
			data.ToRegionId = req.ToRegionId
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
func (l *ProductAccessorieAllocateLogic) List(req *types.ProductAccessorieAllocateListReq) (*types.PageRes[model.ProductAccessorieAllocate], error) {
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
	db = model.PageCondition(db, &req.PageReq)
	db = allocate.Preloads(db, nil)

	if err := db.Find(&res.List).Error; err != nil {
		return nil, errors.New("获取调拨单列表失败")
	}

	return &res, nil
}

// 获取配件调拨单明细
func (l *ProductAccessorieAllocateLogic) Details(req *types.ProductAccessorieAllocateDetailsReq) (*[]model.ProductAccessorieAllocate, error) {
	var (
		allocates []model.ProductAccessorieAllocate

		res []model.ProductAccessorieAllocate
	)

	db := model.DB.Model(&allocates)
	db = model.ProductAccessorieAllocate{}.WhereCondition(db, &req.Where)

	// 获取列表
	db = db.Order("created_at desc")
	db = model.ProductAccessorieAllocate{}.Preloads(db, nil)
	db = db.Preload("Products")
	if err := db.Find(&allocates).Error; err != nil {
		return nil, errors.New("获取调拨单列表失败")
	}

	for _, a := range allocates {
		for _, p := range a.Products {
			a.Product = p
			// 明细展开后不再需要整批产品，避免重复数据与响应体过大
			a.Products = nil
			res = append(res, a)
		}
	}

	return &res, nil
}

// 获取配件调拨单详情
func (l *ProductAccessorieAllocateLogic) Info(req *types.ProductAccessorieAllocateInfoReq) (*model.ProductAccessorieAllocate, error) {
	var (
		allocate model.ProductAccessorieAllocate
	)

	db := model.DB.Model(&allocate)

	db = allocate.Preloads(db, &req.PageReq)

	if err := db.First(&allocate, "id = ?", req.Id).Error; err != nil {
		return nil, errors.New("获取调拨单详情失败")
	}

	return &allocate, nil
}

// 添加配件调拨单产品
func (l *ProductAccessorieAllocateLogic) Add(req *types.ProductAccessorieAllocateAddReq) *errors.Errors {
	var (
		allocate model.ProductAccessorieAllocate
	)

	// 获取调拨单
	if err := model.DB.Preload("Products").First(&allocate, "id = ?", req.Id).Error; err != nil {
		return errors.New("调拨单不存在")
	}

	// 检查调拨单状态
	if allocate.Status != enums.ProductAllocateStatusDraft {
		return errors.New("调拨单状态异常")
	}

	// 所有配件ID
	names := make(map[string]model.ProductAccessorieAllocateProduct)
	for _, p := range allocate.Products {
		names[p.Name] = p
	}

	allocateData := model.ProductAccessorieAllocate{
		ProductCount: allocate.ProductCount,
		ProductTotal: allocate.ProductTotal,
	}

	// 添加配件
	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		for _, rp := range req.Products {
			// 检查数量
			if rp.Quantity <= 0 {
				return errors.New("调拨数量必须大于0")
			}

			// 获取要调拨的配件
			var accessorie model.ProductAccessorie
			if err := tx.Where(&model.ProductAccessorie{
				StoreId: allocate.FromStoreId,
				Name:    rp.Name,
				Status:  enums.ProductAccessorieStatusNormal,
			}).Preload("Store").First(&accessorie).Error; err != nil {
				return errors.New("配件不存在或状态异常")
			}

			// 检查调拨单是否已经存在该配件
			paap, ok := names[accessorie.Name]
			if ok {
				// 已存在，更新数量
				// 检查配件库存
				if accessorie.Stock < (rp.Quantity + paap.Stock) {
					return errors.New("配件库存不足1")
				}
				// 更新配件数量
				if err := tx.Model(&model.ProductAccessorieAllocateProduct{}).
					Where(&model.ProductAccessorieAllocateProduct{
						AllocateId: allocate.Id,
					}).Where("id = ?", paap.Id).Update("stock", gorm.Expr("stock + ?", rp.Quantity)).Error; err != nil {
					return errors.New("更新配件数量失败")
				}

				// 更新已存在的配件
				paap.Stock += rp.Quantity
				// 更新调拨单产品数量
				allocateData.ProductTotal += rp.Quantity

			} else {
				// 不存在，新增
				// 检查配件库存
				if accessorie.Stock < rp.Quantity {
					return errors.New("配件库存不足2")
				}
				data := model.ProductAccessorieAllocateProduct{
					AllocateId: allocate.Id,
					Name:       accessorie.Name,
					Type:       accessorie.Type,
					RetailType: accessorie.RetailType,
					Price:      accessorie.Price,
					Remark:     accessorie.Remark,
					Stock:      rp.Quantity,
					Status:     accessorie.Status,
				}
				// 添加配件
				if err := tx.Create(&data).Error; err != nil {
					return errors.New("添加配件失败")
				}

				// 更新已存在的配件
				names[accessorie.Name] = data
				// 更新调拨单产品数量
				allocateData.ProductCount++
				allocateData.ProductTotal += rp.Quantity
			}

			log := model.ProductHistory{
				Type:       enums.ProductTypeAccessorie,
				Action:     enums.ProductActionTransfer,
				OldValue:   accessorie,
				ProductId:  accessorie.Id,
				StoreId:    accessorie.StoreId,
				SourceId:   allocate.Id,
				Reason:     allocate.Remark,
				OperatorId: l.Staff.Id,
				IP:         l.Ctx.ClientIP(),
			}

			// 扣除配件库存
			accessorie.Stock -= rp.Quantity
			db := tx.Model(&model.ProductAccessorie{}).Where("id = ?", accessorie.Id).Update("stock", accessorie.Stock)
			if accessorie.Stock <= 0 {
				db = db.Update("status", enums.ProductAccessorieStatusNoStock)
			}
			if err := db.Error; err != nil {
				return errors.New("扣除配件库存失败")
			}

			// 添加历史记录
			log.NewValue = accessorie
			if err := tx.Create(&log).Error; err != nil {
				return errors.New("添加历史记录失败")
			}
		}

		// 更新调拨单
		if err := tx.Model(&model.ProductAccessorieAllocate{}).Where("id = ?", allocate.Id).Select([]string{
			"product_count",
			"product_total",
		}).Updates(allocateData).Error; err != nil {
			return errors.New("更新调拨单失败")
		}

		return nil
	}); err != nil {
		return errors.New("添加配件失败：" + err.Error())
	}

	return nil
}

// 移除配件调拨单配件
func (l *ProductAccessorieAllocateLogic) Remove(req *types.ProductAccessorieAllocateRemoveReq) *errors.Errors {
	var (
		allocate model.ProductAccessorieAllocate
	)

	// 获取调拨单
	if err := model.DB.First(&allocate, "id = ?", req.Id).Error; err != nil {
		return errors.New("调拨单不存在")
	}

	if allocate.Status != enums.ProductAllocateStatusDraft {
		return errors.New("调拨单状态异常")
	}

	allocateData := model.ProductAccessorieAllocate{
		ProductCount: allocate.ProductCount,
		ProductTotal: allocate.ProductTotal,
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 获取配件
		var product model.ProductAccessorieAllocateProduct
		where := &model.ProductAccessorieAllocateProduct{
			AllocateId: req.Id,
		}
		if err := tx.Where(where).Find(&product, "id = ?", req.ProductId).Error; err != nil {
			return errors.New("配件不存在")
		}

		// 更新调拨单
		allocateData.ProductCount--
		allocateData.ProductTotal -= product.Stock
		if err := tx.Model(&model.ProductAccessorieAllocate{}).Where("id = ?", allocate.Id).Select([]string{
			"product_count",
			"product_total",
		}).Updates(allocateData).Error; err != nil {
			return errors.New("更新调拨单失败")
		}

		// 移除配件
		if err := tx.Unscoped().Delete(&product).Error; err != nil {
			return errors.New("移除配件失败")
		}

		var accessorie model.ProductAccessorie
		if err := tx.Where(&model.ProductAccessorie{
			Name:    product.Name,
			StoreId: allocate.FromStoreId,
		}).Preload("Store").First(&accessorie).Error; err != nil {
			return errors.New("配件不存在或状态异常")
		}

		log := model.ProductHistory{
			Type:       enums.ProductTypeAccessorie,
			Action:     enums.ProductActionTransferCancel,
			OldValue:   accessorie,
			ProductId:  accessorie.Id,
			StoreId:    accessorie.StoreId,
			SourceId:   allocate.Id,
			Reason:     allocate.Remark,
			OperatorId: l.Staff.Id,
			IP:         l.Ctx.ClientIP(),
		}

		accessorie.Stock += product.Stock
		// 归还库存
		if err := tx.Model(&model.ProductAccessorie{}).Where("id = ?", accessorie.Id).Updates(&model.ProductAccessorie{
			Status: enums.ProductAccessorieStatusNormal,
		}).Update("stock", accessorie.Stock).Error; err != nil {
			return errors.New("归还库存失败")
		}

		// 添加历史记录
		log.NewValue = accessorie
		if err := tx.Create(&log).Error; err != nil {
			return errors.New("添加历史记录失败")
		}

		return nil
	}); err != nil {
		return errors.New(err.Error())
	}

	return nil
}

// 清空配件调拨单配件
func (l *ProductAccessorieAllocateLogic) Clear(req *types.ProductAccessorieAllocateClearReq) *errors.Errors {
	var (
		allocate model.ProductAccessorieAllocate
	)

	// 获取调拨单
	if err := model.DB.Preload("Products").First(&allocate, "id = ?", req.Id).Error; err != nil {
		return errors.New("调拨单不存在")
	}

	if allocate.Status != enums.ProductAllocateStatusDraft {
		return errors.New("调拨单状态异常")
	}

	allocateData := model.ProductAccessorieAllocate{
		ProductCount: allocate.ProductCount,
		ProductTotal: allocate.ProductTotal,
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 更新调拨单
		allocateData.ProductCount = 0
		allocateData.ProductTotal = 0
		if err := tx.Model(&model.ProductAccessorieAllocate{}).Where("id = ?", allocate.Id).Select([]string{
			"product_count",
			"product_total",
		}).Updates(allocateData).Error; err != nil {
			return errors.New("更新调拨单失败")
		}

		for _, product := range allocate.Products {
			// 移除配件
			if err := tx.Unscoped().Delete(&product).Error; err != nil {
				return errors.New("移除配件失败")
			}

			var accessorie model.ProductAccessorie
			if err := tx.Where(&model.ProductAccessorie{
				Name:    product.Name,
				StoreId: allocate.FromStoreId,
			}).Preload("Store").First(&accessorie).Error; err != nil {
				return errors.New("配件不存在或状态异常")
			}

			log := model.ProductHistory{
				Type:       enums.ProductTypeAccessorie,
				Action:     enums.ProductActionTransferCancel,
				OldValue:   accessorie,
				ProductId:  accessorie.Id,
				StoreId:    accessorie.StoreId,
				SourceId:   allocate.Id,
				Reason:     allocate.Remark,
				OperatorId: l.Staff.Id,
				IP:         l.Ctx.ClientIP(),
			}

			accessorie.Stock += product.Stock
			// 归还库存
			if err := tx.Model(&model.ProductAccessorie{}).Where("id = ?", accessorie.Id).Updates(&model.ProductAccessorie{
				Status: enums.ProductAccessorieStatusNormal,
			}).Update("stock", accessorie.Stock).Error; err != nil {
				return errors.New("归还库存失败")
			}

			// 添加历史记录
			log.NewValue = accessorie
			if err := tx.Create(&log).Error; err != nil {
				return errors.New("添加历史记录失败")
			}
		}

		return nil
	}); err != nil {
		return errors.New(err.Error())
	}

	return nil
}

// 确认调拨
func (l *ProductAccessorieAllocateLogic) Confirm(req *types.ProductAccessorieAllocateConfirmReq) *errors.Errors {
	var (
		allocate model.ProductAccessorieAllocate
	)

	// 获取调拨单
	if err := model.DB.
		Preload("Products").
		Preload("FromStore").Preload("ToStore").
		First(&allocate, "id = ?", req.Id).Error; err != nil {
		return errors.New("调拨单不存在")
	}

	if allocate.Status != enums.ProductAllocateStatusDraft {
		return errors.New("调拨单状态异常")
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 确认调拨
		allocate_status := enums.ProductAllocateStatusOnTheWay
		if allocate.Method == enums.ProductAccessorieAllocateMethodRegion {
			allocate_status = enums.ProductAllocateStatusCompleted
		}
		if err := tx.Model(&model.ProductAccessorieAllocate{}).Where("id = ?", allocate.Id).Updates(&model.ProductAccessorieAllocate{
			Status:      allocate_status,
			InitiatorId: l.Staff.Id,
			InitiatorIP: l.Ctx.ClientIP(),
		}).Error; err != nil {
			return errors.New("更新调拨单失败")
		}

		// 开始调拨
		if err := tx.Model(&model.ProductAccessorieAllocateProduct{}).Where(&model.ProductAccessorieAllocateProduct{
			AllocateId: allocate.Id,
		}).Updates(&model.ProductAccessorieAllocateProduct{
			Status: enums.ProductAccessorieStatusAllocate,
		}).Error; err != nil {
			return fmt.Errorf("更新调拨状态失败")
		}

		return nil
	}); err != nil {
		return errors.New("调拨失败: " + err.Error())
	}

	go func() {
		msg := message.NewMessage(l.Ctx)
		msg.SendProductAccessorieAllocateCreateMessage(&message.ProductAccessorieAllocateMessage{
			ProductAccessorieAllocate: &allocate,
		})
	}()

	return nil
}

// 取消调拨
func (l *ProductAccessorieAllocateLogic) Cancel(req *types.ProductAccessorieAllocateCancelReq) *errors.Errors {
	var (
		allocate model.ProductAccessorieAllocate
	)

	// 获取调拨单
	db := model.DB.Model(&allocate)
	db = db.Preload("Products")
	db = db.Preload("FromStore").Preload("ToStore")
	if err := db.First(&allocate, "id = ?", req.Id).Error; err != nil {
		return errors.New("调拨单不存在")
	}

	if allocate.Status != enums.ProductAllocateStatusDraft && allocate.Status != enums.ProductAllocateStatusOnTheWay {
		return errors.New("调拨单状态异常")
	}

	var sendmsg bool

	// 判断调拨单状态
	switch allocate.Status {
	case enums.ProductAllocateStatusDraft:
		sendmsg = false
	case enums.ProductAllocateStatusOnTheWay:
		sendmsg = true
	default:
		return errors.New("调拨单状态异常")
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 取消调拨
		if err := tx.Model(&model.ProductAccessorieAllocate{}).Where("id = ?", allocate.Id).Updates(&model.ProductAccessorieAllocate{
			Status:     enums.ProductAllocateStatusCancelled,
			ReceiverId: l.Staff.Id,
			ReceiverIP: l.Ctx.ClientIP(),
		}).Error; err != nil {
			return errors.New("更新调拨单失败")
		}

		// 更新调拨单配件状态
		if err := tx.Model(&model.ProductAccessorieAllocateProduct{}).Where(&model.ProductAccessorieAllocateProduct{
			AllocateId: allocate.Id,
		}).Updates(&model.ProductAccessorieAllocateProduct{
			Status: enums.ProductAccessorieStatusDraft,
		}).Error; err != nil {
			return fmt.Errorf("更新调拨状态失败")
		}

		// 归还库存
		for _, product := range allocate.Products {
			var accessorie model.ProductAccessorie
			if err := tx.Where(&model.ProductAccessorie{
				Name:    product.Name,
				StoreId: allocate.FromStoreId,
			}).Preload("Store").First(&accessorie).Error; err != nil {
				return errors.New("配件不存在或状态异常")
			}

			log := model.ProductHistory{
				Type:       enums.ProductTypeAccessorie,
				Action:     enums.ProductActionTransferCancel,
				OldValue:   accessorie,
				ProductId:  accessorie.Id,
				StoreId:    accessorie.StoreId,
				SourceId:   allocate.Id,
				Reason:     allocate.Remark,
				OperatorId: l.Staff.Id,
				IP:         l.Ctx.ClientIP(),
			}

			accessorie.Stock += product.Stock
			// 归还库存
			if err := tx.Model(&model.ProductAccessorie{}).Where("id = ?", accessorie.Id).Updates(&model.ProductAccessorie{
				Status: enums.ProductAccessorieStatusNormal,
			}).Update("stock", accessorie.Stock).Error; err != nil {
				return errors.New("归还库存失败")
			}

			// 添加历史记录
			log.NewValue = accessorie
			if err := tx.Create(&log).Error; err != nil {
				return errors.New("添加历史记录失败")
			}
		}

		return nil
	}); err != nil {
		return errors.New("调拨失败: " + err.Error())
	}

	go func() {
		if !sendmsg {
			return
		}
		msg := message.NewMessage(l.Ctx)
		msg.SendProductAccessorieAllocateCancelMessage(&message.ProductAccessorieAllocateMessage{
			ProductAccessorieAllocate: &allocate,
		})
	}()

	return nil
}

// 完成调拨
func (l *ProductAccessorieAllocateLogic) Complete(req *types.ProductAccessorieAllocateCompleteReq) *errors.Errors {
	var (
		allocate model.ProductAccessorieAllocate
	)

	// 获取调拨单
	db := model.DB.Model(&allocate)
	db = db.Preload("Products")
	db = db.Preload("FromStore").Preload("ToStore")
	if err := db.First(&allocate, "id = ?", req.Id).Error; err != nil {
		return errors.New("调拨单不存在")
	}

	if allocate.Status != enums.ProductAllocateStatusOnTheWay {
		return errors.New("调拨单状态异常")
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 确认调拨
		if err := tx.Model(&model.ProductAccessorieAllocate{}).Where("id = ?", allocate.Id).Updates(&model.ProductAccessorieAllocate{
			Status:     enums.ProductAllocateStatusCompleted,
			ReceiverId: l.Staff.Id,
			ReceiverIP: l.Ctx.ClientIP(),
		}).Error; err != nil {
			return errors.New("更新调拨单失败")
		}

		switch allocate.Method {
		case enums.ProductAccessorieAllocateMethodStore, enums.ProductAccessorieAllocateMethodOut:
			{
				// 接收配件
				for _, product := range allocate.Products {
					// 查询配件
					var accessorie model.ProductAccessorie
					if err := tx.Where(&model.ProductAccessorie{
						StoreId: allocate.ToStoreId,
						Name:    product.Name,
					}).Preload("Store").First(&accessorie).Error; err != nil {
						if err != gorm.ErrRecordNotFound {
							return errors.New("查询配件失败")
						}
					}

					log := model.ProductHistory{
						Type:       enums.ProductTypeAccessorie,
						Action:     enums.ProductActionDirectIn,
						OldValue:   nil,
						NewValue:   nil,
						ProductId:  "",
						StoreId:    allocate.ToStoreId,
						SourceId:   allocate.Id,
						Reason:     allocate.Remark,
						OperatorId: l.Staff.Id,
						IP:         l.Ctx.ClientIP(),
					}

					if accessorie.Id == "" {
						// 新增配件
						data := model.ProductAccessorie{
							StoreId:    allocate.ToStoreId,
							Name:       product.Name,
							Type:       product.Type,
							RetailType: product.RetailType,
							Price:      product.Price,
							Remark:     product.Remark,
							Stock:      product.Stock,
							Status:     enums.ProductAccessorieStatusNormal,
						}
						if err := tx.Create(&data).Error; err != nil {
							return fmt.Errorf("【%s】新增配件失败", product.Name)
						}

						// 更新记录
						data.Store = *allocate.ToStore
						log.ProductId = data.Id
						log.NewValue = data

					} else {
						// 更新记录
						log.ProductId = accessorie.Id
						log.OldValue = accessorie

						// 更新配件
						accessorie.Stock += product.Stock
						if err := tx.Model(&model.ProductAccessorie{}).Where("id = ?", accessorie.Id).Updates(&model.ProductAccessorie{
							Status: enums.ProductAccessorieStatusNormal,
						}).Update("stock", accessorie.Stock).Error; err != nil {
							return fmt.Errorf("【%s】更新配件失败", product.Name)
						}

						log.NewValue = accessorie
					}

					if err := tx.Create(&log).Error; err != nil {
						return errors.New("添加历史记录失败")
					}
				}
			}
		}

		return nil
	}); err != nil {
		return errors.New("调拨失败: " + err.Error())
	}

	go func() {
		msg := message.NewMessage(l.Ctx)
		msg.SendProductAccessorieAllocateCompleteMessage(&message.ProductAccessorieAllocateMessage{
			ProductAccessorieAllocate: &allocate,
		})
	}()

	return nil
}
