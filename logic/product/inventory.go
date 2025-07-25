package product

import (
	"errors"
	"jdy/enums"
	"jdy/message"
	"jdy/model"
	"jdy/types"
	"jdy/utils"
	"slices"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductInventoryLogic struct {
	Ctx   *gin.Context
	Staff *model.Staff
}

// 创建产品盘点单
func (l *ProductInventoryLogic) Create(req *types.ProductInventoryCreateReq) (*model.ProductInventory, error) {
	var (
		db = model.DB
	)

	// 转换参数
	data, err := utils.StructToStruct[model.ProductInventory](req)
	if err != nil {
		return nil, err
	}

	// 设置创建人
	data.CreatorId = l.Staff.Id

	// 创建应盘记录
	if err := db.Transaction(func(tx *gorm.DB) error {
		switch req.Type {
		case enums.ProductTypeUsedFinished:
			var products []model.ProductFinished
			pdb := tx.Model(&model.ProductFinished{})
			pdb = model.CreateProductInventoryCondition(pdb, req)

			if err := pdb.Where(&model.ProductFinished{
				Status:  enums.ProductStatusNormal,
				StoreId: req.StoreId,
			}).Find(&products).Error; err != nil {
				return err
			}
			if len(products) == 0 {
				return errors.New("没有符合条件的产品")
			}

			// 计算并添加产品
			for _, product := range products {
				// 添加产品
				data.ShouldProducts = append(data.ShouldProducts, model.ProductInventoryProduct{
					ProductType: enums.ProductTypeUsedFinished,
					ProductCode: strings.ToUpper(product.Code),
					Status:      enums.ProductInventoryProductStatusShould,
				})
				// 产品总数
				data.ShouldCount++
				// 总件数
				data.CountQuantity++
				// 总标价
				data.CountPrice = data.CountPrice.Add(product.LabelPrice)
				// 总重量
				data.CountWeightMetal = data.CountWeightMetal.Add(product.WeightMetal)
				// 判断是否可以转态
				if err := product.Status.CanTransitionTo(enums.ProductStatusCheck); err != nil {
					return errors.New("产品状态不正确")
				}
			}

			// 设置状态
			data.Status = enums.ProductInventoryStatusDraft

			// 创建记录
			if err := tx.Create(&data).Error; err != nil {
				return errors.New("创建失败")
			}

		case enums.ProductTypeUsedOld:
			var products []model.ProductOld
			pdb := tx.Model(&model.ProductOld{})
			pdb = model.CreateProductInventoryCondition(pdb, req)

			if err := pdb.Where(&model.ProductOld{
				IsOur:   true,
				Status:  enums.ProductStatusNormal,
				StoreId: req.StoreId,
			}).Find(&products).Error; err != nil {
				return err
			}
			if len(products) == 0 {
				return errors.New("没有符合条件的产品")
			}
			// 计算并添加产品
			for _, product := range products {
				// 添加产品
				data.ShouldProducts = append(data.ShouldProducts, model.ProductInventoryProduct{
					ProductType: enums.ProductTypeUsedOld,
					ProductCode: strings.ToUpper(product.Code),
					Status:      enums.ProductInventoryProductStatusShould,
				})
				// 产品总数
				data.ShouldCount++
				// 总件数
				data.CountQuantity++
				// 总标价
				data.CountPrice = data.CountPrice.Add(product.LabelPrice)
				// 总重量
				data.CountWeightMetal = data.CountWeightMetal.Add(product.WeightMetal)
				// 判断是否可以转态
				if err := product.Status.CanTransitionTo(enums.ProductStatusCheck); err != nil {
					return errors.New("产品状态不正确")
				}
			}

			// 设置状态
			data.Status = enums.ProductInventoryStatusDraft

			// 创建记录
			if err := tx.Create(&data).Error; err != nil {
				return errors.New("创建失败")
			}

		default:
			return errors.New("类型不正确")
		}

		var InventoryPersons []model.Staff
		if err := tx.Where("id in (?)", req.InventoryPersonIds).Find(&InventoryPersons).Error; err != nil {
			return err
		}
		if err := tx.Model(&data).Association("InventoryPersons").Append(InventoryPersons); err != nil {
			return errors.New("添加盘点人员失败")
		}

		return nil
	}); err != nil {
		return nil, err
	}

	go func() {
		var (
			product_inventory model.ProductInventory
			pdb               = model.DB.Model(&product_inventory)
		)
		pdb = product_inventory.Preloads(pdb, nil, false)
		// 查询记录
		if err := pdb.First(&product_inventory, "id = ?", data.Id).Error; err != nil {
			return
		}
		msg := message.NewMessage(l.Ctx)
		msg.SendProductInventoryCreateMessage(&message.ProductInventoryMessage{
			ProductInventory: &product_inventory,
		})
	}()

	return &data, nil
}

// 搜索产品盘点单
func (l *ProductInventoryLogic) List(req *types.ProductInventoryListReq) (*types.PageRes[model.ProductInventory], error) {

	var (
		inventory model.ProductInventory

		res types.PageRes[model.ProductInventory]
	)

	db := model.DB.Model(&inventory)
	db = inventory.WhereCondition(db, &req.Where)

	// 获取总数
	if err := db.Count(&res.Total).Error; err != nil {
		return nil, errors.New("获取总数失败")
	}

	// 获取列表
	db = inventory.Preloads(db, &req.Where, false)
	db = db.Order("created_at desc")
	db = model.PageCondition(db, &req.PageReq)
	if err := db.Find(&res.List).Error; err != nil {
		return nil, errors.New("获取列表失败")
	}

	return &res, nil
}

// 获取产品盘点单详情
func (l *ProductInventoryLogic) Info(req *types.ProductInventoryInfoReq) (*model.ProductInventory, error) {
	var (
		inventory model.ProductInventory

		res model.ProductInventory
	)
	db := model.DB.Model(&inventory)

	db = db.Where("id = ?", req.Id)

	where := types.ProductInventoryWhere{
		PageReqNon: types.PageReqNon{
			Page:  req.Page,
			Limit: req.Limit,
		},
		ProductStatus: req.ProductStatus,
	}

	db = inventory.Preloads(db, &where, false)

	if err := db.First(&res).Error; err != nil {
		return nil, errors.New("获取失败")
	}

	if res.Status.IsOver() {
		db = inventory.Preloads(db, &where, true)
		if err := db.First(&res).Error; err != nil {
			return nil, errors.New("获取失败")
		}
	}

	return &res, nil
}

// 添加产品盘点单产品
func (l *ProductInventoryLogic) Add(req *types.ProductInventoryAddReq) error {
	var (
		inventory model.ProductInventory
	)
	db := model.DB.Model(&model.ProductInventory{})
	db = db.Preload("ActualProducts")
	db = inventory.Preloads(db, nil, false)
	if err := db.First(&inventory, "id = ?", req.Id).Error; err != nil {
		return errors.New("获取失败")
	}

	var InventoryPersonIds []string
	for _, staff := range inventory.InventoryPersons {
		InventoryPersonIds = append(InventoryPersonIds, staff.Id)
	}

	if can := inventory.Status.CanEdit(enums.ProductInventoryStatusInventorying, l.Staff.Id, InventoryPersonIds, inventory.InspectorId); !can {
		return errors.New("当前状态不允许这样操作")
	}

	now := time.Now()

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 计算并添加产品
		for _, code := range req.Codes {
			for _, product := range inventory.ActualProducts {
				if product.ProductCode == strings.ToUpper(code) {
					return errors.New(code + "产品已存在")
				}
			}

			switch inventory.Type {
			case enums.ProductTypeUsedFinished:
				var finished model.ProductFinished
				if err := tx.Unscoped().Where(&model.ProductFinished{
					Code:    code,
					StoreId: inventory.StoreId,
				}).First(&finished).Error; err != nil {
					if err == gorm.ErrRecordNotFound {
						return errors.New("[" + code + "] 不存在")
					}

					return errors.New("[" + code + "] 查询失败")
				}

			case enums.ProductTypeUsedOld:
				var old model.ProductOld
				if err := tx.Unscoped().Where(&model.ProductOld{
					Code:    code,
					StoreId: inventory.StoreId,
				}).First(&old).Error; err != nil {
					if err == gorm.ErrRecordNotFound {
						return errors.New("[" + code + "] 不存在")
					}

					return errors.New("[" + code + "] 查询失败")
				}

			default:
				return errors.New("[" + code + "]不存在")
			}

			// 添加产品
			if err := tx.Create(&model.ProductInventoryProduct{
				ProductInventoryId: req.Id,
				ProductType:        inventory.Type,
				ProductCode:        strings.ToUpper(code),
				Status:             enums.ProductInventoryProductStatusActual,
				InventoryTime:      &now,
			}).Error; err != nil {
				return errors.New("[" + code + "]添加失败")
			}
			// 产品总数
			if err := tx.Model(&model.ProductInventory{}).
				Where("id = ?", req.Id).
				Update("actual_count", gorm.Expr("actual_count + ?", 1)).Error; err != nil {
				return errors.New("[" + code + "]添加失败")
			}
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (l *ProductInventoryLogic) Remove(req *types.ProductInventoryRemoveReq) error {
	var (
		inventory model.ProductInventory
	)
	db := model.DB.Model(&model.ProductInventory{})
	db = inventory.Preloads(db, nil, false)
	if err := db.First(&inventory, "id = ?", req.Id).Error; err != nil {
		return errors.New("获取失败")
	}

	var InventoryPersonIds []string
	for _, staff := range inventory.InventoryPersons {
		InventoryPersonIds = append(InventoryPersonIds, staff.Id)
	}
	if can := inventory.Status.CanEdit(enums.ProductInventoryStatusInventorying, l.Staff.Id, InventoryPersonIds, inventory.InspectorId); !can {
		return errors.New("当前状态不允许这样操作")
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		var product model.ProductInventoryProduct
		if err := tx.Where("product_inventory_id = ? and id = ?", req.Id, req.ProductId).First(&product).Error; err != nil {
			return errors.New("获取失败")
		}

		del := tx.Delete(&product)
		if err := del.Error; err != nil {
			return errors.New("删除失败")
		}
		if del.RowsAffected == 0 {
			return errors.New("删除失败")
		}

		// 产品总数
		if err := tx.Model(&model.ProductInventory{}).
			Where("id = ?", req.Id).
			Update("actual_count", gorm.Expr("actual_count - ?", 1)).Error; err != nil {
			return errors.New("删除失败")
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

// 切换盘点单状态
func (l *ProductInventoryLogic) Change(req *types.ProductInventoryChangeReq) error {
	var (
		inventory model.ProductInventory
	)
	db := model.DB.Where("id = ?", req.Id)
	db = inventory.Preloads(db, nil, false)

	if err := db.First(&inventory).Error; err != nil {
		return errors.New("获取失败")
	}

	if err := inventory.Status.CanTransitionTo(req.Status); err != nil {
		return errors.New("当前状态不允许这样操作")
	}

	var InventoryPersonIds []string
	for _, staff := range inventory.InventoryPersons {
		InventoryPersonIds = append(InventoryPersonIds, staff.Id)
	}
	if can := inventory.Status.CanEdit(req.Status, l.Staff.Id, InventoryPersonIds, inventory.InspectorId); !can {
		return errors.New("处理人不一致")
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 设置状态
		res := model.ProductInventory{
			Status: req.Status,
		}
		inventory.Status = req.Status
		// 如果是完成状态
		if req.Status == enums.ProductInventoryStatusToBeVerified || req.Status == enums.ProductInventoryStatusCancelled {
			// 计算盘盈、盘亏
			var data, products []model.ProductInventoryProduct
			if err := tx.Where("product_inventory_id = ?", req.Id).Find(&products).Error; err != nil {
				return errors.New("获取失败")
			}

			var ShouldCodes []string                           // 应盘
			var ShouldProducts []model.ProductInventoryProduct // 应盘
			var ActualCodes []string                           // 实盘
			var ActualProducts []model.ProductInventoryProduct // 实盘
			for _, product := range products {
				switch product.Status {
				case enums.ProductInventoryProductStatusShould:
					ShouldCodes = append(ShouldCodes, strings.ToUpper(product.ProductCode))
					ShouldProducts = append(ShouldProducts, product)
				case enums.ProductInventoryProductStatusActual:
					ActualCodes = append(ActualCodes, strings.ToUpper(product.ProductCode))
					ActualProducts = append(ActualProducts, product)
				}
			}

			for _, actual := range ShouldProducts {
				if !slices.Contains(ActualCodes, actual.ProductCode) {
					loss := model.ProductInventoryProduct{
						ProductInventoryId: req.Id,
						ProductType:        inventory.Type,
						ProductCode:        strings.ToUpper(actual.ProductCode),
						Status:             enums.ProductInventoryProductStatusLoss,
					}
					data = append(data, loss)
					res.LossCount++
				}
			}
			for _, actual := range ActualProducts {
				if !slices.Contains(ShouldCodes, actual.ProductCode) {
					extra := model.ProductInventoryProduct{
						ProductInventoryId: req.Id,
						ProductType:        inventory.Type,
						ProductCode:        strings.ToUpper(actual.ProductCode),
						Status:             enums.ProductInventoryProductStatusExtra,
					}
					data = append(data, extra)
					res.ExtraCount++
				}
			}

			if len(data) > 0 {
				if err := tx.Create(&data).Error; err != nil {
					return errors.New("添加失败")
				}
			}
		}

		if err := tx.Model(&model.ProductInventory{}).Where("id = ?", inventory.Id).Updates(&res).Error; err != nil {
			return errors.New("更新失败")
		}

		return nil
	}); err != nil {
		return err
	}

	go func() {
		msg := message.NewMessage(l.Ctx)
		msg.SendProductInventoryUpdateMessage(&message.ProductInventoryMessage{
			ProductInventory: &inventory,
		})
	}()

	return nil
}
