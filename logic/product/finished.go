package product

import (
	"jdy/enums"
	"jdy/errors"
	"jdy/model"
	"jdy/types"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProductFinishedLogic struct {
	Ctx   *gin.Context
	Staff *model.Staff
}

// 成品列表
func (p *ProductFinishedLogic) List(req *types.ProductFinishedListReq) (*types.ProductFinishedListRes[model.ProductFinished], error) {
	var (
		product model.ProductFinished

		res types.ProductFinishedListRes[model.ProductFinished]
	)

	// 获取总数
	db := model.DB.Model(&model.ProductFinished{})
	db = product.WhereCondition(db, &req.Where)
	if err := db.Count(&res.Total).Error; err != nil {
		return nil, errors.New("获取成品列表数量失败")
	}
	if res.Total == 0 {
		return &res, nil
	}

	// 获取列表
	db = db.Order("created_at desc")
	db = model.PageCondition(db, &req.PageReq)
	if err := db.Find(&res.List).Error; err != nil {
		return nil, errors.New("获取成品列表失败")
	}

	// 获取入网费
	adb := model.DB.Model(&model.ProductFinished{})
	adb = product.WhereCondition(adb, &req.Where)
	if err := adb.Select("SUM(access_fee) as access_fee").Scan(&res.AccessFee).Error; err != nil {
		return nil, errors.New("获取成品列表入网费失败")
	}

	// 获取标签价
	ldb := model.DB.Model(&model.ProductFinished{})
	ldb = product.WhereCondition(ldb, &req.Where)
	if err := ldb.Select("SUM(label_price) as label_price").Scan(&res.LabelPrice).Error; err != nil {
		return nil, errors.New("获取成品列表标签价失败")
	}

	// 获取金重
	wdb := model.DB.Model(&model.ProductFinished{})
	wdb = product.WhereCondition(wdb, &req.Where)
	if err := wdb.Select("SUM(weight_metal) as weight_metal").Scan(&res.WeightMetal).Error; err != nil {
		return nil, errors.New("获取成品列表金重失败")
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
			Code: strings.ToUpper(req.Code),
		}).
		Preload("Store").
		First(&product).Error; err != nil {
		return nil, errors.New("获取成品信息失败")
	}

	return &product, nil
}

// 成品检索
func (p *ProductFinishedLogic) Retrieval(req *types.ProductFinishedRetrievalReq) (*model.ProductFinished, error) {
	var (
		product model.ProductFinished
	)

	if err := model.DB.
		Where(model.ProductFinished{
			Code:    strings.ToUpper(req.Code),
			StoreId: req.StoreId,
		}).
		Preload("Store").
		First(&product).Error; err != nil {
		return nil, errors.New("获取成品信息失败")
	}

	return &product, nil
}

// 上传成品图片
func (p *ProductFinishedLogic) Upload(req *types.ProductFinishedUploadReq) error {
	if err := model.DB.Transaction(func(tx *gorm.DB) error {
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

		product.Images = req.Images
		if err := tx.Model(&model.ProductFinished{}).Clauses(clause.Returning{}).Where("id = ?", product.Id).Updates(&model.ProductFinished{
			Images: req.Images,
		}).Error; err != nil {
			return errors.New("上传成品图片失败")
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
