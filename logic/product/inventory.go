package product

import (
	"errors"
	"jdy/enums"
	"jdy/message"
	"jdy/model"
	"jdy/types"
	"jdy/utils"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductInventoryLogic struct {
	Ctx   *gin.Context
	Staff *types.Staff
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
		case enums.ProductTypeFinished:
			var products []model.ProductFinished
			pdb := tx.Model(&model.ProductFinished{})
			pdb = model.CreateProductInventoryCondition(pdb, req)

			if err := pdb.Where(&model.ProductFinished{
				Status: enums.ProductStatusNormal,
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
					ProductType: enums.ProductTypeFinished,
					ProductCode: product.Code,
					Status:      enums.ProductInventoryProductStatusShould,
				})
				// 产品总数
				data.ShouldCount++
				// 总件数
				data.ContQuantity++
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

		case enums.ProductTypeOld:
			var products []model.ProductOld
			pdb := tx.Model(&model.ProductOld{})
			pdb = model.CreateProductInventoryCondition(pdb, req)

			if err := pdb.Where(&model.ProductOld{
				IsOur:  true,
				Status: enums.ProductStatusNormal,
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
					ProductType: enums.ProductTypeOld,
					ProductCode: product.Code,
					Status:      enums.ProductInventoryProductStatusShould,
				})
				// 产品总数
				data.ShouldCount++
				// 总件数
				data.ContQuantity++
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
		msg.SendProductInventoryCreateMessage(&message.ProductInventoryCreate{
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
	db = model.PageCondition(db, req.Page, req.Limit)
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
		ProductStatus: req.ProductStatus,
	}

	db = inventory.Preloads(db, &where, false)

	if err := db.First(&res).Error; err != nil {
		return nil, errors.New("获取失败")
	}

	if inventory.Status.IsOver() {
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

	if err := model.DB.First(&inventory, "id = ?", req.Id).Error; err != nil {
		return errors.New("获取失败")
	}

	if err := inventory.Status.CanEdit(enums.ProductInventoryStatusInventorying, l.Staff.Id, inventory.InventoryPersonId, inventory.InspectorId); !err {
		return errors.New("当前状态不允许这样操作")
	}

	now := time.Now()

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		switch inventory.Type {
		case enums.ProductTypeFinished:
			var products []model.ProductFinished
			pdb := tx.Model(&model.ProductFinished{})
			pdb = pdb.Where("code in (?)", req.Codes)
			pdb = pdb.Where(&model.ProductFinished{
				Status: enums.ProductStatusNormal,
			})
			if err := pdb.Find(&products).Error; err != nil {
				return err
			}
			if len(products) == 0 || len(products) != len(req.Codes) {
				return errors.New("没有符合条件的产品")
			}
			// 计算并添加产品
			for _, product := range products {
				// 添加产品
				inventory.ActualProducts = append(inventory.ActualProducts, model.ProductInventoryProduct{
					ProductType:   enums.ProductTypeFinished,
					ProductCode:   product.Code,
					Status:        enums.ProductInventoryProductStatusActual,
					InventoryTime: &now,
				})
				// 产品总数
				inventory.ActualCount++
			}
		case enums.ProductTypeOld:
			var products []model.ProductOld
			pdb := tx.Model(&model.ProductOld{})
			pdb = pdb.Where("code in (?)", req.Codes)
			pdb = pdb.Where(&model.ProductOld{
				IsOur:  true,
				Status: enums.ProductStatusNormal,
			})
			if err := pdb.Find(&products).Error; err != nil {
				return err
			}
			if len(products) == 0 || len(products) != len(req.Codes) {
				return errors.New("没有符合条件的产品")
			}
			// 计算并添加产品
			for _, product := range products {
				// 添加产品
				inventory.ActualProducts = append(inventory.ActualProducts, model.ProductInventoryProduct{
					ProductType:   enums.ProductTypeOld,
					ProductCode:   product.Code,
					Status:        enums.ProductInventoryProductStatusActual,
					InventoryTime: &now,
				})
				// 产品总数
				inventory.ActualCount++
			}
		}

		if err := tx.Save(&inventory).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return errors.New("添加失败")
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

	if can := inventory.Status.CanEdit(req.Status, l.Staff.Id, inventory.InventoryPersonId, inventory.InspectorId); !can {
		return errors.New("处理人不一致")
	}

	if err := db.Model(&inventory).Updates(model.ProductInventory{Status: req.Status}).Error; err != nil {
		return errors.New("更新失败")
	}

	go func() {
		msg := message.NewMessage(l.Ctx)
		msg.SendProductInventoryUpdateMessage(&message.ProductInventoryUpdate{
			ProductInventory: &inventory,
		})
	}()

	return nil
}
