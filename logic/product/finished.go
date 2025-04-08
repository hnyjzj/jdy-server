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

type ProductFinishedLogic struct {
	Ctx   *gin.Context
	Staff *types.Staff
}

// 成品列表
func (p *ProductFinishedLogic) List(req *types.ProductFinishedListReq) (*types.PageRes[model.ProductFinished], error) {
	var (
		product model.ProductFinished

		res types.PageRes[model.ProductFinished]
	)

	db := model.DB.Model(&product)
	db = product.WhereCondition(db, &req.Where).Where(&model.ProductFinished{
		Status: enums.ProductStatusNormal,
	})

	// 获取总数
	if err := db.Count(&res.Total).Error; err != nil {
		return nil, errors.New("获取成品列表失败")
	}

	// 获取列表
	db = db.Order("created_at desc")
	db = model.PageCondition(db, req.Page, req.Limit)
	if err := db.Find(&res.List).Error; err != nil {
		return nil, errors.New("获取成品列表失败")
	}

	return &res, nil
}

// 成品详情
func (p *ProductFinishedLogic) Info(req *types.ProductFinishedInfoReq) (*model.ProductFinished, error) {
	var (
		product model.ProductFinished
	)

	if err := model.DB.
		Where(model.ProductFinished{
			Code: req.Code,
		}).
		Preload("Store").
		First(&product).Error; err != nil {
		return nil, errors.New("获取成品信息失败")
	}

	return &product, nil
}

// 更新成品信息
func (p *ProductFinishedLogic) Update(req *types.ProductFinishedUpdateReq) error {

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		data, err := utils.StructToStruct[model.ProductFinished](req)
		if err != nil {
			return errors.New("验证参数失败")
		}

		var product model.ProductFinished
		if err := tx.Model(&model.ProductFinished{}).
			Preload("Store").
			Where("id = ?", req.Id).First(&product).Error; err != nil {
			return errors.New("获取成品信息失败")
		}

		history := model.ProductHistory{
			Type:       enums.ProductTypeFinished,
			Action:     enums.ProductActionUpdate,
			OldValue:   product,
			ProductId:  product.Id,
			SourceId:   product.Id,
			StoreId:    product.StoreId,
			OperatorId: p.Staff.Id,
			IP:         p.Ctx.ClientIP(),
		}

		if err := tx.Model(&product).Clauses(clause.Returning{}).Where("id = ?", req.Id).Updates(&data).Error; err != nil {
			return errors.New("更新成品信息失败")
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
