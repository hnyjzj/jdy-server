package product

import (
	"jdy/enums"
	"jdy/errors"
	"jdy/model"
	"jdy/types"
	"jdy/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProductAccessorieEnterLogic struct {
	Ctx   *gin.Context
	Staff *model.Staff
}

// 配件入库
func (l *ProductAccessorieEnterLogic) Create(req *types.ProductAccessorieEnterCreateReq) (*model.ProductAccessorieEnter, error) {
	enter := model.ProductAccessorieEnter{
		StoreId:    req.StoreId,
		Remark:     req.Remark,
		Status:     enums.ProductEnterStatusDraft,
		OperatorId: l.Staff.Id,
		IP:         l.Ctx.ClientIP(),
	}

	if err := model.DB.Create(&enter).Error; err != nil {
		return nil, errors.New("创建入库单失败")
	}

	return &enter, nil
}

// 配件入库单列表
func (l *ProductAccessorieEnterLogic) EnterList(req *types.ProductAccessorieEnterListReq) (*types.PageRes[model.ProductAccessorieEnter], error) {
	var (
		enter model.ProductAccessorieEnter

		res types.PageRes[model.ProductAccessorieEnter]
	)

	db := model.DB.Model(&enter)
	db = enter.WhereCondition(db, &req.Where)

	// 获取总数
	if err := db.Count(&res.Total).Error; err != nil {
		return nil, errors.New("获取配件总数失败")
	}

	// 获取列表
	db = enter.Preloads(db, nil)
	db = db.Order("created_at desc")
	db = model.PageCondition(db, &req.PageReq)

	if err := db.Find(&res.List).Error; err != nil {
		return nil, errors.New("获取入库单失败")
	}

	return &res, nil
}

// 配件入库单详情
func (l *ProductAccessorieEnterLogic) EnterInfo(req *types.ProductAccessorieEnterInfoReq) (*model.ProductAccessorieEnter, error) {
	var (
		enter model.ProductAccessorieEnter
	)

	db := model.DB.Model(&enter)

	// 获取配件入库单详情
	db = enter.Preloads(db, &req.PageReq)

	if err := db.First(&enter, "id = ?", req.Id).Error; err != nil {
		return nil, errors.New("获取配件入库单详情失败")
	}

	return &enter, nil
}

// 配件入库单添加配件
func (l *ProductAccessorieEnterLogic) AddProduct(req *types.ProductAccessorieEnterAddProductReq) error {
	// 查询入库单
	var enter model.ProductAccessorieEnter
	if err := model.DB.First(&enter, "id = ?", req.EnterId).Error; err != nil {
		return errors.New("入库单不存在")
	}

	if enter.Status != enums.ProductEnterStatusDraft {
		return errors.New("入库单已结束")
	}

	if len(req.Products) == 0 {
		return errors.New("请添加配件")
	}

	// 所有配件ID
	names := make(map[string]model.ProductAccessorieEnterProduct)
	for _, p := range enter.Products {
		names[strings.TrimSpace(strings.ToUpper(p.Name))] = p
	}
	// 添加配件的结果
	enterData := model.ProductAccessorieEnter{
		ProductCount: enter.ProductCount,
		ProductTotal: enter.ProductTotal,
	}
	// 添加配件入库
	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		for _, p := range req.Products {
			// 检查调拨单是否已经存在该配件
			paap, ok := names[strings.TrimSpace(strings.ToUpper(p.Name))]
			if ok {
				// 配件已存在，更新配件
				if err := tx.Model(&model.ProductAccessorieEnterProduct{}).Where("id = ?", paap.Id).Update("stock", gorm.Expr("stock + ?", p.Stock)).Error; err != nil {
					return errors.New("[" + p.Name + "]入库失败")
				}
				// 更新入库单配件数量
				enterData.ProductTotal += p.Stock
			} else {
				data := model.ProductAccessorieEnterProduct{
					EnterId:    enter.Id,
					Name:       strings.TrimSpace(strings.ToUpper(p.Name)),
					Type:       p.Type,
					RetailType: p.RetailType,
					Price:      p.Price,
					Remark:     p.Remark,
					Stock:      p.Stock,
					Status:     enums.ProductAccessorieStatusDraft,
				}

				if err := tx.Create(&data).Error; err != nil {
					return errors.New("[" + p.Name + "]入库失败")
				}

				// 更新入库单配件数量
				enterData.ProductCount++
				enterData.ProductTotal += p.Stock
			}

		}

		// 更新入库单配件数量
		if err := tx.Model(&model.ProductAccessorieEnter{}).Where("id = ?", enter.Id).Select([]string{
			"product_count",
			"product_total",
		}).Updates(enterData).Error; err != nil {
			return errors.New("入库单更新失败")
		}

		return nil
	}); err != nil {
		return errors.New("配件录入失败: " + err.Error())
	}

	return nil
}

// 入库单编辑配件
func (l *ProductAccessorieEnterLogic) EditProduct(req *types.ProductAccessorieEnterEditProductReq) error {
	// 编辑配件逻辑存在潜在并发风险，使用事务包裹并加锁
	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 查询入库单
		var enter model.ProductAccessorieEnter
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&enter, "id = ?", req.EnterId).Error; err != nil {
			return errors.New("入库单不存在")
		}

		if enter.Status != enums.ProductEnterStatusDraft {
			return errors.New("入库单已结束")
		}

		var product model.ProductAccessorieEnterProduct
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where(&model.ProductAccessorieEnterProduct{
			EnterId: enter.Id,
		}).First(&product, "id = ?", req.ProductId).Error; err != nil {
			return errors.New("配件不存在")
		}

		if product.Status != enums.ProductAccessorieStatusDraft {
			return errors.New("配件状态异常")
		}

		// 转换数据结构
		data, err := utils.StructToStruct[model.ProductAccessorieEnterProduct](req.Product)
		if err != nil {
			return errors.New("配件录入失败: 参数错误")
		}

		// 更新入库单
		enter.ProductTotal += product.Stock - data.Stock
		if enter.ProductTotal < 0 {
			return errors.New("配件数量不能小于0")
		}
		if err := tx.Model(&model.ProductAccessorieEnter{}).Where("id = ?", enter.Id).Select([]string{
			"product_total",
		}).Updates(enter).Error; err != nil {
			return errors.New("入库单更新失败")
		}

		// 更新配件
		if err := tx.Model(&model.ProductAccessorieEnterProduct{}).Where("id = ?", product.Id).Updates(data).Error; err != nil {
			return errors.New("配件更新失败")
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

// 入库单删除配件
func (l *ProductAccessorieEnterLogic) DelProduct(req *types.ProductAccessorieEnterDelProductReq) error {
	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 查询入库单
		var enter model.ProductAccessorieEnter
		if err := tx.First(&enter, "id = ?", req.EnterId).Error; err != nil {
			return errors.New("入库单不存在")
		}

		if enter.Status != enums.ProductEnterStatusDraft {
			return errors.New("入库单已结束")
		}

		// 查询配件
		var product model.ProductAccessorieEnterProduct
		if err := tx.Where(&model.ProductAccessorieEnterProduct{
			EnterId: enter.Id,
		}).First(&product, "id = ?", req.ProductId).Error; err != nil {
			return errors.New("配件不存在")
		}

		if product.Status != enums.ProductAccessorieStatusDraft {
			return errors.New("配件状态异常")
		}

		// 更新入库单
		enter.ProductCount--
		enter.ProductTotal -= product.Stock
		if err := tx.Model(&model.ProductAccessorieEnter{}).Where("id = ?", enter.Id).Select([]string{
			"product_count",
			"product_total",
		}).Updates(enter).Error; err != nil {
			return errors.New("入库单更新失败")
		}

		// 删除配件(真实删除)
		if err := tx.Unscoped().Where("id = ?", product.Id).Delete(&model.ProductAccessorieEnterProduct{}).Error; err != nil {
			return errors.New("配件删除失败")
		}

		return nil
	}); err != nil {
		return errors.New("配件删除失败：" + err.Error())
	}
	return nil
}

// 入库单清空配件
func (l *ProductAccessorieEnterLogic) ClearProduct(req *types.ProductAccessorieEnterClearProductReq) error {
	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 查询入库单
		var enter model.ProductAccessorieEnter
		if err := tx.First(&enter, "id = ?", req.EnterId).Error; err != nil {
			return errors.New("入库单不存在")
		}

		if enter.Status != enums.ProductEnterStatusDraft {
			return errors.New("入库单已结束")
		}

		// 更新入库单
		enter.ProductCount = 0
		enter.ProductTotal = 0
		if err := tx.Model(&model.ProductAccessorieEnter{}).Where("id = ?", enter.Id).Select([]string{
			"product_count",
			"product_total",
		}).Updates(enter).Error; err != nil {
			return errors.New("入库单更新失败")
		}

		// 删除配件(真实删除)
		if err := tx.Unscoped().Where(&model.ProductAccessorieEnterProduct{
			EnterId: enter.Id,
		}).Delete(&model.ProductAccessorieEnterProduct{}).Error; err != nil {
			return errors.New("配件删除失败")
		}

		return nil
	}); err != nil {
		return errors.New("配件删除失败：" + err.Error())
	}
	return nil
}

// 入库单完成
func (l *ProductAccessorieEnterLogic) Finish(req *types.ProductAccessorieEnterFinishReq) error {
	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 查询入库单
		var enter model.ProductAccessorieEnter
		if err := tx.Preload("Products").Preload("Store").First(&enter, "id = ?", req.EnterId).Error; err != nil {
			return errors.New("入库单不存在")
		}

		if enter.Status != enums.ProductEnterStatusDraft {
			return errors.New("入库单已结束")
		}

		if len(enter.Products) == 0 {
			return errors.New("入库单没有配件")
		}

		// 更新配件状态
		for _, product := range enter.Products {
			// 查询配件
			if err := tx.Where("id = ?", product.Id).
				Where(&model.ProductAccessorieEnterProduct{
					Status: enums.ProductAccessorieStatusDraft,
				}).
				First(&product).Error; err != nil {
				return errors.New("配件状态不正确或配件不存在")
			}

			var accessorie model.ProductAccessorie
			if err := tx.Model(&model.ProductAccessorie{}).
				Where(&model.ProductAccessorie{
					StoreId: enter.StoreId,
					Name:    product.Name,
				}).First(&accessorie).Error; err != nil {
				if err != gorm.ErrRecordNotFound {
					return errors.New("配件查询失败")
				}
			}

			// 添加记录
			history := model.ProductHistory{
				Type:       enums.ProductTypeAccessorie,
				Action:     enums.ProductActionEntry,
				OldValue:   nil,
				NewValue:   nil,
				ProductId:  "",
				StoreId:    enter.StoreId,
				SourceId:   enter.Id,
				OperatorId: l.Staff.Id,
				IP:         l.Ctx.ClientIP(),
			}

			if accessorie.Id == "" {
				// 添加配件
				accessorie = model.ProductAccessorie{
					StoreId:    enter.StoreId,
					Name:       product.Name,
					Type:       product.Type,
					RetailType: product.RetailType,
					Price:      product.Price,
					Remark:     product.Remark,
					Stock:      product.Stock,
					Status:     enums.ProductAccessorieStatusNormal,
				}
				if err := tx.Create(&accessorie).Error; err != nil {
					return errors.New("配件添加失败")
				}

				// 更新记录
				accessorie.Store = *enter.Store
				history.NewValue = accessorie
				history.ProductId = accessorie.Id

			} else {
				// 更新记录
				accessorie.Store = *enter.Store
				history.OldValue = accessorie
				history.ProductId = accessorie.Id

				// 更新库存
				accessorie.Stock += product.Stock
				if err := tx.Model(&model.ProductAccessorie{}).Where("id = ?", accessorie.Id).Updates(&model.ProductAccessorie{
					Status: enums.ProductAccessorieStatusNormal,
				}).Update("stock", accessorie.Stock).Error; err != nil {
					return errors.New("归还库存失败")
				}

				// 更新记录
				history.NewValue = accessorie
			}

			if err := tx.Create(&history).Error; err != nil {
				return errors.New("配件记录添加失败")
			}

			// 更新配件状态
			if err := tx.Model(&model.ProductAccessorieEnterProduct{}).Where("id = ?", product.Id).Updates(model.ProductAccessorieEnterProduct{
				Status: enums.ProductAccessorieStatusNormal,
			}).Error; err != nil {
				return errors.New("配件状态更新失败")
			}
		}

		// 更新入库单状态
		enter.Status = enums.ProductEnterStatusCompleted
		if err := tx.Model(&model.ProductAccessorieEnter{}).Where("id = ?", enter.Id).Updates(model.ProductAccessorieEnter{
			Status: enter.Status,
		}).Error; err != nil {
			return errors.New("入库单更新失败")
		}

		return nil
	}); err != nil {
		return errors.New("操作失败：" + err.Error())
	}

	return nil
}

// 入库单取消
func (l *ProductAccessorieEnterLogic) Cancel(req *types.ProductAccessorieEnterCancelReq) error {
	// 查询入库单
	var enter model.ProductAccessorieEnter
	if err := model.DB.Preload("Products").First(&enter, "id = ?", req.EnterId).Error; err != nil {
		return errors.New("入库单不存在")
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		switch enter.Status {
		case enums.ProductEnterStatusDraft:
			// 草稿直接取消
			if err := tx.Model(&model.ProductAccessorieEnter{}).Where("id = ?", enter.Id).Updates(model.ProductAccessorieEnter{
				Status: enums.ProductEnterStatusCanceled,
			}).Error; err != nil {
				return errors.New("入库单取消失败")
			}

		case enums.ProductEnterStatusCompleted:
			// 已完成的入库单，需要将配件状态还原
			for _, product := range enter.Products {
				// 判断产品状态
				if product.Status != enums.ProductAccessorieStatusNormal {
					return errors.New("配件状态异常")
				}

				// 查询配件
				var accessorie model.ProductAccessorie
				if err := tx.Model(&model.ProductAccessorie{}).Preload("Store").Where(&model.ProductAccessorie{
					StoreId: enter.StoreId,
					Name:    product.Name,
					Status:  enums.ProductAccessorieStatusNormal,
				}).First(&accessorie).Error; err != nil {
					return errors.New("配件查询失败")
				}
				if accessorie.Stock < product.Stock || accessorie.Stock-product.Stock < 0 {
					return errors.New("配件库存不足")
				}

				accessorie.Stock -= product.Stock
				status := enums.ProductAccessorieStatusNormal
				if accessorie.Stock == 0 {
					status = enums.ProductAccessorieStatusNoStock
				}
				if err := tx.Model(&model.ProductAccessorie{}).Where("id = ?", accessorie.Id).
					Updates(model.ProductAccessorie{
						Status: status,
					}).Update("stock", accessorie.Stock).Error; err != nil {
					return errors.New("配件库存更新失败")
				}

				// 添加记录
				history := model.ProductHistory{
					Type:       enums.ProductTypeAccessorie,
					OldValue:   product,
					NewValue:   accessorie,
					Action:     enums.ProductActionEntryCancel,
					ProductId:  accessorie.Id,
					StoreId:    enter.StoreId,
					SourceId:   enter.Id,
					OperatorId: l.Staff.Id,
					IP:         l.Ctx.ClientIP(),
				}
				if err := tx.Create(&history).Error; err != nil {
					return errors.New("配件记录添加失败")
				}
			}

			// 更新入库单状态
			if err := tx.Model(&model.ProductAccessorieEnter{}).Where("id = ?", enter.Id).Updates(model.ProductAccessorieEnter{
				Status: enums.ProductEnterStatusCanceled,
			}).Error; err != nil {
				return errors.New("入库单取消失败")
			}

		default:
			return errors.New("入库单状态不支持取消")
		}

		return nil

	}); err != nil {
		return errors.New(err.Error())
	}

	return nil
}
