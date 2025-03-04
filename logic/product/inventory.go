package product

import (
	"errors"
	"jdy/enums"
	"jdy/message"
	"jdy/model"
	"jdy/types"
	"jdy/utils"

	"gorm.io/gorm"
)

type ProductInventoryLogic struct {
	ProductLogic
}

func (l *ProductInventoryLogic) Create(req *types.ProductInventoryCreateReq) (*model.ProductInventory, error) {
	// 查询应盘列表
	var (
		products []model.Product
		db       = model.DB
	)
	pdb := db.Model(&model.Product{}).Where(&model.Product{Status: enums.ProductStatusNormal, Type: req.Type})
	if len(req.Brand) > 0 {
		pdb = pdb.Where("brand in (?)", req.Brand)
	}
	if len(req.Class) > 0 {
		pdb = pdb.Where("class in (?)", req.Class)
	}
	if len(req.Category) > 0 {
		pdb = pdb.Where("category in (?)", req.Category)
	}
	if len(req.Craft) > 0 {
		pdb = pdb.Where("craft in (?)", req.Craft)
	}
	if len(req.Material) > 0 {
		pdb = pdb.Where("material in (?)", req.Material)
	}
	if len(req.Quality) > 0 {
		pdb = pdb.Where("quality in (?)", req.Quality)
	}
	if len(req.Gem) > 0 {
		pdb = pdb.Where("gem in (?)", req.Gem)
	}

	if err := pdb.Find(&products).Error; err != nil {
		return nil, err
	}

	if len(products) == 0 {
		return nil, errors.New("没有符合条件的产品")
	}

	// 转换参数
	data, err := utils.StructToStruct[model.ProductInventory](req)
	if err != nil {
		return nil, err
	}

	// 创建应盘记录
	if err := db.Transaction(func(tx *gorm.DB) error {
		// 计算并添加产品
		for _, product := range products {
			// 添加产品
			data.Products = append(data.Products, model.ProductInventoryProduct{
				ProductId: product.Id,
				Product:   product,

				Status: enums.ProductInventoryProductStatusShould,
			})
			// 产品总数
			data.CountShould++
			// 总件数
			data.ContQuantity++
			// 总价格
			data.CountPrice = data.CountPrice.Add(product.Price)
			// 总重量
			data.CountWeightMetal = data.CountWeightMetal.Add(product.WeightMetal)
			// 判断是否可以转态
			if err := product.Status.CanTransitionTo(enums.ProductStatusCheck); err != nil {
				return errors.New("产品状态不正确")
			}
			// 更新产品状态
			if err := tx.Model(&product).Updates(model.Product{Status: enums.ProductStatusCheck}).Error; err != nil {
				return errors.New("更新产品状态失败")
			}
		}

		// 设置状态
		data.Status = enums.ProductInventoryStatusDraft

		// 创建记录
		if err := tx.Create(&data).Error; err != nil {
			return errors.New("创建失败")
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
		pdb = product_inventory.Preloads(pdb, nil)
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
	db = inventory.Preloads(db, &req.Where)
	db = db.Order("created_at desc")
	db = model.PageCondition(db, req.Page, req.Limit)
	if err := db.Find(&res.List).Error; err != nil {
		return nil, errors.New("获取列表失败")
	}

	return &res, nil
}

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

	db = inventory.Preloads(db, &where)

	if err := db.First(&res).Error; err != nil {
		return nil, errors.New("获取失败")
	}

	return &res, nil
}

func (l *ProductInventoryLogic) Change(req *types.ProductInventoryChangeReq) error {
	var (
		inventory model.ProductInventory
	)
	db := model.DB.Where("id = ?", req.Id)
	db = inventory.Preloads(db, nil)

	if err := db.First(&inventory).Error; err != nil {
		return errors.New("获取失败")
	}

	if err := inventory.Status.CanTransitionTo(req.Status); err != nil {
		return errors.New("当前状态不允许这样操作")
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
