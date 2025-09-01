package product

import (
	"jdy/enums"
	"jdy/errors"
	"jdy/model"
	"jdy/types"
	"jdy/utils"
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

	db = db.Where("status IN (?)", []enums.ProductStatus{
		enums.ProductStatusNormal,
		enums.ProductStatusAllocate,
	})

	if err := db.Count(&res.Total).Error; err != nil {
		return nil, errors.New("获取成品列表数量失败")
	}
	if res.Total == 0 {
		return &res, nil
	}

	// 获取列表
	db = db.Order("created_at desc")
	db = model.PageCondition(db, &req.PageReq)
	db = product.Preloads(db)
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

// 空图片列表
func (ProductFinishedLogic) EmptyImage(req *types.ProductFinishedEmptyImageReq) ([]model.ProductFinished, error) {
	var (
		product model.ProductFinished

		res []model.ProductFinished
	)

	// 获取总数
	db := model.DB.Model(&model.ProductFinished{})
	db = product.WhereCondition(db, &types.ProductFinishedWhere{
		StoreId: req.StoreId,
	})
	db = db.Where("status IN (?)", []enums.ProductStatus{
		enums.ProductStatusNormal,
		enums.ProductStatusAllocate,
	})

	db = db.Where("images is null or images = ''")
	db = product.Preloads(db)

	if err := db.Find(&res).Error; err != nil {
		return nil, errors.New("获取成品列表失败")
	}

	return res, nil
}

// 成品详情
func (p *ProductFinishedLogic) Info(req *types.ProductFinishedInfoReq) (*model.ProductFinished, error) {
	var (
		product model.ProductFinished
	)

	db := model.DB.Model(&model.ProductFinished{})

	db = db.Where(model.ProductFinished{
		Code: strings.ToUpper(req.Code),
	})
	db = product.Preloads(db)

	if err := db.First(&product).Error; err != nil {
		return nil, errors.New("获取成品信息失败")
	}

	return &product, nil
}

// 成品检索
func (p *ProductFinishedLogic) Retrieval(req *types.ProductFinishedRetrievalReq) (*model.ProductFinished, error) {
	var (
		product model.ProductFinished
	)

	db := model.DB.Model(&model.ProductFinished{})

	db = db.Where(model.ProductFinished{
		Code: strings.ToUpper(req.Code),
	})
	if req.StoreId != "" {
		db = db.Where(model.ProductFinished{
			StoreId: req.StoreId,
		})
	}
	db = product.Preloads(db)

	if err := db.First(&product).Error; err != nil {
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
			Clauses(clause.Locking{Strength: "UPDATE"}).
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

		data.Class = data.GetClass()
		// 统一 code 大小写
		if data.Code != "" {
			data.Code = strings.ToUpper(data.Code)
		}

		// 若 code 发生变更，则同步更新旧料的 code_finished，保持关联一致
		if data.Code != "" && data.Code != product.Code {
			if err := tx.Model(&model.ProductOld{}).Where(&model.ProductOld{
				CodeFinished: product.Code,
			}).Update("code_finished", data.Code).Error; err != nil {
				return errors.New("同步旧料成品条码失败")
			}
		}

		if err := tx.Model(&model.ProductFinished{}).Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", product.Id).Omit(
			"id", "created_at", "deleted_at",
			"status", "store_id", "enter_id",
		).Updates(&data).Error; err != nil {
			return errors.New("更新成品信息失败")
		}

		// 添加记录
		history.NewValue = data
		if err := tx.Create(&history).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
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
		if err := tx.Model(&model.ProductFinished{}).Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", product.Id).Updates(&model.ProductFinished{
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
