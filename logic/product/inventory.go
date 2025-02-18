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
		for _, product := range products {
			data.CountShould++
			data.ContQuantity++
			data.CountPrice = data.CountPrice.Add(product.Price)
			data.CountWeightMetal = data.CountWeightMetal.Add(product.WeightMetal)

			data.Products = append(data.Products, model.ProductInventoryProduct{
				ProductId: product.Id,
				Product:   product,

				Status: enums.ProductInventoryProductStatusShould,
			})

			if err := product.Status.CanTransitionTo(enums.ProductStatusCheck); err != nil {
				return err
			}

			if err := tx.Model(&product).Updates(model.Product{Status: enums.ProductStatusCheck}).Error; err != nil {
				return err
			}
		}

		if err := tx.Create(&data).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}
	return &data, nil
}
