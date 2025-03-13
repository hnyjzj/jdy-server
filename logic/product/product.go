package product

import (
	"jdy/enums"
	"jdy/errors"
	"jdy/model"
	"jdy/types"
	"jdy/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
		if err := tx.Model(&model.Product{}).Where("id = ?", req.Id).First(&product).Error; err != nil {
			return errors.New("获取产品信息失败")
		}

		old := product

		if err := tx.Model(&model.Product{}).Where("id = ?", req.Id).Updates(&data).Error; err != nil {
			return errors.New("更新产品信息失败")
		}

		// 添加记录
		if err := tx.Create(&model.ProductHistory{
			Action:     enums.ProductActionUpdate,
			OldValue:   old,
			NewValue:   data,
			ProductId:  data.Id,
			StoreId:    data.StoreId,
			SourceId:   data.Id,
			OperatorId: p.Staff.Id,
			IP:         p.Ctx.ClientIP(),
		}).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

// 产品操作记录列表
func (l *ProductLogic) History(req *types.ProductHistoryReq) (*types.PageRes[model.ProductHistory], error) {
	var (
		logs model.ProductHistory

		res types.PageRes[model.ProductHistory]
	)

	db := model.DB.Model(&logs)
	db = logs.WhereCondition(db, &req.Where)

	// 获取总数
	if err := db.Count(&res.Total).Error; err != nil {
		return nil, errors.New("获取总数失败")
	}

	// 获取列表
	db = db.Order("created_at desc")
	db = model.PageCondition(db, req.Page, req.Limit)
	db = db.Preload("Operator")
	db = db.Preload("Product")
	if err := db.Find(&res.List).Error; err != nil {
		return nil, errors.New("获取列表失败")
	}

	return &res, nil
}
