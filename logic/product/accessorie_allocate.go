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
	db = model.PageCondition(db, &req.PageReq)

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

	db = db.Preload("Products")
	db = db.Preload("FromStore").Preload("ToStore")
	db = db.Preload("Operator")

	if err := db.First(&allocate, "id = ?", req.Id).Error; err != nil {
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
			if err := tx.First(&accessorie, "id = ?", rp.ProductId).Error; err != nil {
				return errors.New("配件不存在")
			}
			// 检查配件状态
			if accessorie.Status != enums.ProductAccessorieStatusNormal {
				return errors.New("配件状态异常")
			}
			// 检查门店
			if accessorie.StoreId != allocate.FromStoreId {
				return errors.New("配件不在调拨单的门店中")
			}

			// 检查调拨单是否已经存在该配件
			paap, ok := names[accessorie.Name]
			if ok { // 已存在，更新数量
				// 检查配件库存
				if accessorie.Stock < (rp.Quantity + paap.Stock) {
					return errors.New("配件库存不足")
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

			} else { // 不存在，新增
				// 检查配件库存
				if accessorie.Stock < rp.Quantity {
					return errors.New("配件库存不足")
				}
				data := model.ProductAccessorieAllocateProduct{
					AllocateId: allocate.Id,
					ProductAccessorie: model.ProductAccessorie{
						StoreId:    accessorie.StoreId,
						Name:       accessorie.Name,
						Type:       accessorie.Type,
						RetailType: accessorie.RetailType,
						Remark:     accessorie.Remark,
						Stock:      rp.Quantity,
						Status:     accessorie.Status,
					},
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
func (p *ProductAccessorieAllocateLogic) Remove(req *types.ProductAccessorieAllocateRemoveReq) *errors.Errors {
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
		var product []model.ProductAccessorieAllocateProduct
		// 获取配件
		where := &model.ProductAccessorieAllocateProduct{
			AllocateId: req.Id,
		}

		if err := model.DB.Where(where).Find(&product, "id IN (?)", req.ProductIds).Error; err != nil {
			return errors.New("配件不存在")
		}

		for _, p := range product {
			// 更新调拨单产品数量
			allocateData.ProductCount--
			allocateData.ProductTotal -= p.Stock
		}

		// 更新调拨单
		if err := tx.Model(&model.ProductAccessorieAllocate{}).Where("id = ?", allocate.Id).Select([]string{
			"product_count",
			"product_total",
		}).Updates(allocateData).Error; err != nil {
			return errors.New("更新调拨单失败")
		}

		// 移除配件
		if err := model.DB.Where(where).Delete(&product).Error; err != nil {
			return errors.New("移除配件失败")
		}

		return nil
	}); err != nil {
		return errors.New(err.Error())
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
		Preload("Products").
		Preload("FromStore").Preload("ToStore").
		First(&allocate, "id = ?", req.Id).Error; err != nil {
		return errors.New("调拨单不存在")
	}

	if allocate.Status != enums.ProductAllocateStatusDraft {
		return errors.New("调拨单状态异常")
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 扣除配件库存
		for _, p := range allocate.Products {
			var product model.ProductAccessorie
			if err := tx.First(&product, "id = ?", p.Id).Error; err != nil {
				return fmt.Errorf("【%s】不存在", p.Name)
			}
			// 扣除库存
			product.Stock -= p.Stock
			if product.Stock < 0 {
				return fmt.Errorf("【%s】库存不足", p.Name)
			}
			if product.Stock == 0 {
				product.Status = enums.ProductAccessorieStatusNoStock
			}
			if err := tx.Model(&model.ProductAccessorie{}).Where("id = ?", product.Id).Select([]string{
				"stock",
				"status",
			}).Updates(product).Error; err != nil {
				return fmt.Errorf("【%s】扣除库存失败", p.Name)
			}

			if err := tx.Model(&model.ProductAccessorieAllocateProduct{}).Where("id = ?", p.Id).Updates(&model.ProductAccessorieAllocateProduct{
				ProductAccessorie: model.ProductAccessorie{
					Status: enums.ProductAccessorieStatusAllocate,
				},
			}).Error; err != nil {
				return fmt.Errorf("【%s】更新调拨单配件失败", p.Name)
			}
		}

		// 确认调拨
		if err := tx.Model(&model.ProductAccessorieAllocate{}).Where("id = ?", allocate.Id).Updates(&model.ProductAccessorieAllocate{
			Status: enums.ProductAllocateStatusOnTheWay,
		}).Error; err != nil {
			return errors.New("更新调拨单失败")
		}

		return nil
	}); err != nil {
		return errors.New("调拨失败: " + err.Error())
	}

	go func() {
		msg := message.NewMessage(p.Ctx)
		allocate.Operator = p.Staff
		msg.SendProductAccessorieAllocateCreateMessage(&message.ProductAccessorieAllocateMessage{
			ProductAccessorieAllocate: &allocate,
		})
	}()

	return nil
}

// 取消调拨
func (p *ProductAccessorieAllocateLogic) Cancel(req *types.ProductAccessorieAllocateCancelReq) *errors.Errors {
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

	sendmsg := false

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 判断调拨单状态
		if allocate.Status == enums.ProductAllocateStatusOnTheWay {
			sendmsg = true
			for _, product := range allocate.Products {
				// 获取调拨单配件
				var paap model.ProductAccessorieAllocateProduct
				if err := tx.First(&paap, "id = ?", product.Id).Error; err != nil {
					return fmt.Errorf("【%s】不存在", product.Name)
				}

				if paap.Status != enums.ProductAccessorieStatusAllocate {
					return fmt.Errorf("【%s】状态异常", product.Name)
				}

				// 获取配件
				var accessorie model.ProductAccessorie
				if err := tx.Where(&model.ProductAccessorie{
					StoreId: allocate.ToStoreId,
					Name:    product.Name,
				}).First(&accessorie).Error; err != nil {
					return fmt.Errorf("【%s】不存在", product.Name)
				}

				// 归还库存
				accessorie.Stock += product.Stock
				if err := tx.Model(&model.ProductAccessorie{}).Where("id = ?", accessorie.Id).Select([]string{
					"stock",
					"status",
				}).Updates(&model.ProductAccessorie{
					Stock:  accessorie.Stock,
					Status: enums.ProductAccessorieStatusNormal,
				}).Error; err != nil {
					return fmt.Errorf("【%s】归还库存失败", product.Name)
				}

				// 更新调拨单配件状态
				if err := tx.Model(&model.ProductAccessorieAllocateProduct{}).Where("id = ?", paap.Id).Updates(&model.ProductAccessorieAllocateProduct{
					ProductAccessorie: model.ProductAccessorie{
						Status: enums.ProductAccessorieStatusNormal,
					},
				}).Error; err != nil {
					return fmt.Errorf("【%s】更新调拨单配件失败", product.Name)
				}
			}
		}

		// 取消调拨
		if err := tx.Model(&model.ProductAccessorieAllocate{}).Where("id = ?", allocate.Id).Updates(&model.ProductAccessorieAllocate{
			Status: enums.ProductAllocateStatusCancelled,
		}).Error; err != nil {
			return errors.New("更新调拨单失败")
		}

		return nil
	}); err != nil {
		return errors.New("调拨失败: " + err.Error())
	}

	go func() {
		if !sendmsg {
			return
		}
		msg := message.NewMessage(p.Ctx)
		allocate.Operator = p.Staff
		msg.SendProductAccessorieAllocateCancelMessage(&message.ProductAccessorieAllocateMessage{
			ProductAccessorieAllocate: &allocate,
		})
	}()

	return nil
}

// 完成调拨
func (p *ProductAccessorieAllocateLogic) Complete(req *types.ProductAccessorieAllocateCompleteReq) *errors.Errors {
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
		// 接收配件
		for _, product := range allocate.Products {

			// 查询配件
			var accessorie model.ProductAccessorie
			if err := tx.Where(&model.ProductAccessorie{
				StoreId: allocate.ToStoreId,
				Name:    product.Name,
			}).First(&accessorie).Error; err != nil {
				if err != gorm.ErrRecordNotFound {
					return errors.New("查询配件失败")
				}
			}

			if accessorie.Id == "" {
				// 新增配件
				data := model.ProductAccessorie{
					StoreId:    allocate.ToStoreId,
					Name:       product.Name,
					Type:       product.Type,
					RetailType: product.RetailType,
					Remark:     product.Remark,
					Stock:      product.Stock,
					Status:     enums.ProductAccessorieStatusNormal,
				}
				if err := tx.Create(&data).Error; err != nil {
					return fmt.Errorf("【%s】新增配件失败", product.Name)
				}

			} else {
				// 更新配件
				if err := tx.Model(&model.ProductAccessorie{}).Where("id = ?", accessorie.Id).Updates(&model.ProductAccessorie{
					Stock:  accessorie.Stock + product.Stock,
					Status: enums.ProductAccessorieStatusNormal,
				}).Error; err != nil {
					return fmt.Errorf("【%s】更新配件失败", product.Name)
				}
			}
		}

		// 确认调拨
		if err := tx.Model(&model.ProductAccessorieAllocate{}).Where("id = ?", allocate.Id).Updates(&model.ProductAccessorieAllocate{
			Status: enums.ProductAllocateStatusCompleted,
		}).Error; err != nil {
			return errors.New("更新调拨单失败")
		}

		return nil
	}); err != nil {
		return errors.New("调拨失败: " + err.Error())
	}

	go func() {
		msg := message.NewMessage(p.Ctx)
		allocate.Operator = p.Staff
		msg.SendProductAccessorieAllocateCompleteMessage(&message.ProductAccessorieAllocateMessage{
			ProductAccessorieAllocate: &allocate,
		})
	}()

	return nil
}
