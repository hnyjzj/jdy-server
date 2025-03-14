package product

import (
	"jdy/enums"
	"jdy/errors"
	"jdy/model"
	"jdy/types"
	"jdy/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProductLogic struct {
	Ctx   *gin.Context
	Staff *types.Staff
}

// 产品列表
func (p *ProductLogic) List(req *types.ProductListReq) (*types.PageRes[model.Product], error) {
	var (
		product model.Product

		res types.PageRes[model.Product]
	)

	db := model.DB.Model(&product)
	db = product.WhereCondition(db, &req.Where)

	// 获取总数
	if err := db.Count(&res.Total).Error; err != nil {
		return nil, errors.New("获取产品列表失败")
	}

	// 获取列表
	db = db.Order("created_at desc")
	db = model.PageCondition(db, req.Page, req.Limit)
	if err := db.Find(&res.List).Error; err != nil {
		return nil, errors.New("获取产品列表失败")
	}

	return &res, nil
}

// 产品详情
func (p *ProductLogic) Info(req *types.ProductInfoReq) (*model.Product, error) {
	var (
		product model.Product
	)

	if err := model.DB.
		Where(model.Product{
			Code: req.Code,
		}).
		Preload("Store").
		First(&product).Error; err != nil {
		return nil, errors.New("获取产品信息失败")
	}

	return &product, nil
}

// 更新产品信息
func (p *ProductLogic) Update(req *types.ProductUpdateReq) error {

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		data, err := utils.StructToStruct[model.Product](req)
		if err != nil {
			return errors.New("验证参数失败")
		}

		var product model.Product
		if err := tx.Model(&model.Product{}).
			Preload("Store").
			Preload("RecycleStore").
			Where("id = ?", req.Id).First(&product).Error; err != nil {
			return errors.New("获取产品信息失败")
		}

		history := model.ProductHistory{
			Action:     enums.ProductActionUpdate,
			OldValue:   product,
			ProductId:  product.Id,
			SourceId:   product.Id,
			StoreId:    product.StoreId,
			OperatorId: p.Staff.Id,
			IP:         p.Ctx.ClientIP(),
		}

		if err := tx.Model(&product).Clauses(clause.Returning{}).Where("id = ?", req.Id).Updates(&data).Error; err != nil {
			return errors.New("更新产品信息失败")
		}

		// 添加记录
		history.NewValue = product
		if err := tx.Create(&history).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
