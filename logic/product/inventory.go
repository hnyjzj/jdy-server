package product

import (
	"errors"
	"jdy/enums"
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
	db = db.Preload("Store")
	db = db.Preload("InventoryPerson")
	db = db.Preload("Inspector")
	db = db.Preload("Products").Preload("Products.Product")

	db = db.Order("created_at desc")
	db = model.PageCondition(db, req.Page, req.Limit)
	if err := db.Find(&res.List).Error; err != nil {
		return nil, errors.New("获取列表失败")
	}

	return &res, nil
}
