package product

import (
	"jdy/errors"
	"jdy/model"
	"jdy/types"
	"jdy/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProductAccessorieCategoryLogic struct {
	Ctx   *gin.Context
	Staff *model.Staff
}

// 配件条目列表
func (p *ProductAccessorieCategoryLogic) List(req *types.ProductAccessorieCategoryListReq) (*types.PageRes[model.ProductAccessorieCategory], error) {
	var (
		product model.ProductAccessorieCategory

		res types.PageRes[model.ProductAccessorieCategory]
	)

	db := model.DB.Model(&product)
	db = product.WhereCondition(db, &req.Where)

	// 获取总数
	if err := db.Count(&res.Total).Error; err != nil {
		return nil, errors.New("获取配件条目总数失败")
	}

	// 获取列表
	db = db.Order("created_at desc")
	db = db.Preload("Products")
	db = model.PageCondition(db, req.Page, req.Limit)
	if err := db.Find(&res.List).Error; err != nil {
		return nil, errors.New("获取配件条目列表失败")
	}

	return &res, nil
}

// 配件条目详情
func (p *ProductAccessorieCategoryLogic) Info(req *types.ProductAccessorieCategoryInfoReq) (*model.ProductAccessorieCategory, error) {
	var (
		product model.ProductAccessorieCategory
	)

	if err := model.DB.
		Where("id = ?", req.Id).
		First(&product).Error; err != nil {
		return nil, errors.New("获取配件条目信息失败")
	}

	return &product, nil
}

// 新增配件条目
func (p *ProductAccessorieCategoryLogic) Create(req *types.ProductAccessorieCategoryCreateReq) error {

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		data, err := utils.StructToStruct[[]model.ProductAccessorieCategory](req.List)
		if err != nil {
			return errors.New("验证参数失败")
		}

		if err := tx.Create(&data).Error; err != nil {
			return errors.New("新增配件条目失败")
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

// 更新配件条目信息
func (p *ProductAccessorieCategoryLogic) Update(req *types.ProductAccessorieCategoryUpdateReq) error {

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		data, err := utils.StructToStruct[model.ProductAccessorieCategory](req)
		if err != nil {
			return errors.New("验证参数失败")
		}

		var product model.ProductAccessorieCategory
		if err := tx.Model(&model.ProductAccessorieCategory{}).
			Preload("Store").
			Where("id = ?", req.Id).First(&product).Error; err != nil {
			return errors.New("获取配件条目信息失败")
		}

		if err := tx.Model(&product).Clauses(clause.Returning{}).Where("id = ?", req.Id).Updates(&data).Error; err != nil {
			return errors.New("更新配件条目信息失败")
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

// 删除配件条目
func (p *ProductAccessorieCategoryLogic) Delete(req *types.ProductAccessorieCategoryDeleteReq) error {
	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		var product model.ProductAccessorieCategory
		if err := tx.Model(&model.ProductAccessorieCategory{}).
			Preload("Store").Where("id = ?", req.Id).First(&product).Error; err != nil {
			return errors.New("获取配件条目信息失败")
		}

		if err := tx.Where("id = ?", req.Id).Delete(&product).Error; err != nil {
			return errors.New("删除配件条目失败")
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
