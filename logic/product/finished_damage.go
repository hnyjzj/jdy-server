package product

import (
	"jdy/enums"
	"jdy/errors"
	"jdy/model"
	"jdy/types"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type ProductFinishedDamageLogic struct {
	Ctx   *gin.Context
	Staff *types.Staff
}

// 成品报损
func (l *ProductFinishedDamageLogic) Damage(req *types.ProductDamageReq) *errors.Errors {
	// 查询商品信息
	var product model.ProductFinished
	if err := model.DB.Where(&model.ProductFinished{Code: req.Code}).
		Preload("Store").
		First(&product).Error; err != nil {
		return errors.New("商品不存在")
	}

	// 判断成品状态
	if product.Status != enums.ProductStatusNormal {
		return errors.New("成品不在库存中")
	}
	// 判断是否可以报损
	if err := product.Status.CanTransitionTo(enums.ProductStatusDamage); err != nil {
		return errors.New("成品状态不允许报损")
	}

	// 开启事务
	if err := model.DB.Transaction(func(tx *gorm.DB) error {

		history := model.ProductHistory{
			Action:     enums.ProductActionDamage,
			OldValue:   product,
			ProductId:  product.Id,
			StoreId:    product.StoreId,
			SourceId:   product.Id,
			OperatorId: l.Staff.Id,
			IP:         l.Ctx.ClientIP(),
		}

		// 更新商品状态
		product.Status = enums.ProductStatusDamage
		if err := tx.Save(&product).Error; err != nil {
			return err
		}

		// 添加记录
		history.NewValue = product
		if err := tx.Create(&history).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return errors.New("报损失败")
	}

	return nil
}

// 产品列表
func (p *ProductFinishedDamageLogic) List(req *types.ProductFinishedListReq) (*types.PageRes[model.ProductFinished], error) {
	var (
		product model.ProductFinished

		res types.PageRes[model.ProductFinished]
	)

	db := model.DB.Model(&product)
	db = product.WhereCondition(db, &req.Where).Where(&model.ProductFinished{
		Status: enums.ProductStatusDamage,
	})

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
func (p *ProductFinishedDamageLogic) Info(req *types.ProductFinishedInfoReq) (*model.ProductFinished, error) {
	var (
		product model.ProductFinished
	)

	if err := model.DB.
		Where(model.ProductFinished{
			Code: req.Code,
		}).
		Preload("Store").
		First(&product).Error; err != nil {
		return nil, errors.New("获取产品信息失败")
	}

	return &product, nil
}

// 产品转换
func (l *ProductFinishedDamageLogic) Conversion(req *types.ProductConversionReq) *errors.Errors {
	// 开启事务
	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 查询商品信息
		var product model.ProductFinished
		if err := tx.Where("id = ?", req.Id).First(&product).Error; err != nil {
			return errors.New("商品不存在")
		}

		// 判断产品状态
		if product.Status != enums.ProductStatusDamage {
			return errors.New("产品不在库存中")
		}

		log := &model.ProductHistory{
			OldValue:   product,
			ProductId:  product.Id,
			StoreId:    product.StoreId,
			SourceId:   product.Id,
			OperatorId: l.Staff.Id,
			IP:         l.Ctx.ClientIP(),
		}

		switch req.Type {
		case enums.ProductTypeFinished:
			log.Action = enums.ProductActionDamageToNew
			product.Status = enums.ProductStatusNormal
			if err := tx.Save(&product).Error; err != nil {
				return errors.New("转换失败")
			}
			log.NewValue = product
		case enums.ProductTypeOld:
			log.Action = enums.ProductActionDamageToOld
			if err := tx.Delete(&model.ProductFinished{}, "id = ?", req.Id).Error; err != nil {
				return errors.New("删除失败")
			}
			var old model.ProductOld
			if err := tx.Unscoped().Where(&model.ProductOld{Code: product.Code}).First(&old).Error; err != nil {
				if err != gorm.ErrRecordNotFound {
					return errors.New("旧品不存在")
				}
			}

			if old.Id != "" {
				return errors.New("旧料已存在")
			}

			old = model.ProductOld{
				Code:            product.Code,
				Name:            product.Name,
				Status:          enums.ProductStatusNormal,
				LabelPrice:      product.LabelPrice,
				Brand:           product.Brand,
				Material:        product.Material,
				Quality:         product.Quality,
				Gem:             product.Gem,
				Category:        product.Category,
				Craft:           product.Craft,
				WeightMetal:     product.WeightMetal,
				WeightTotal:     product.WeightTotal,
				ColorGem:        product.ColorGem,
				WeightGem:       product.WeightGem,
				NumGem:          product.NumGem,
				Clarity:         product.Clarity,
				WeightOther:     product.WeightOther,
				NumOther:        product.NumOther,
				Remark:          req.Remark,
				StoreId:         product.StoreId,
				IsOur:           true,
				QualityActual:   decimal.NewFromInt(1),
				RecycleSource:   enums.ProductRecycleSourceBaoSongZhenGold,
				RecycleSourceId: product.Id,
			}
			if err := tx.Create(&old).Error; err != nil {
				return errors.New("转换失败")
			}
			log.NewValue = old

			if err := tx.Delete(&model.ProductFinished{}, "id = ?", req.Id).Error; err != nil {
				return errors.New("删除失败")
			}

		default:
			return errors.New("转换类型错误")
		}

		if err := tx.Create(&log).Error; err != nil {
			return errors.New("添加记录失败")
		}

		return nil
	}); err != nil {
		return errors.New("转换失败：" + err.Error())
	}

	return nil
}
