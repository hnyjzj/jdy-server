package product

import (
	"jdy/enums"
	"jdy/errors"
	"jdy/model"
	"jdy/types"
	"jdy/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type ProductFinishedDamageLogic struct {
	Ctx   *gin.Context
	Staff *model.Staff
}

// 成品报损
func (l *ProductFinishedDamageLogic) Damage(req *types.ProductDamageReq) *errors.Errors {
	// 查询商品信息
	var product model.ProductFinished
	if err := model.DB.Where(&model.ProductFinished{Code: strings.ToUpper(req.Code)}).
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
			Type:       enums.ProductTypeFinished,
			OldValue:   product,
			ProductId:  product.Id,
			StoreId:    product.StoreId,
			SourceId:   product.Id,
			OperatorId: l.Staff.Id,
			IP:         l.Ctx.ClientIP(),
		}

		// 更新商品状态
		product.Status = enums.ProductStatusDamage
		if err := tx.Model(&model.ProductFinished{}).Where("id = ?", product.Id).Updates(model.ProductFinished{
			Status: enums.ProductStatusDamage,
		}).Error; err != nil {
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
	req.Where.Status = enums.ProductStatusDamage
	db = product.WhereCondition(db, &req.Where)

	// 获取总数
	if err := db.Count(&res.Total).Error; err != nil {
		return nil, errors.New("获取产品列表失败")
	}

	// 获取列表
	db = db.Order("created_at desc")
	db = model.PageCondition(db, &req.PageReq)
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
			Code: strings.ToUpper(req.Code),
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
		if err := tx.Preload("Store").Where("id = ?", req.Id).First(&product).Error; err != nil {
			return errors.New("商品不存在")
		}

		// 判断产品状态
		if product.Status != enums.ProductStatusDamage && product.Status != enums.ProductStatusSold {
			return errors.New("产品不在库存中")
		}

		log := &model.ProductHistory{
			Type:       enums.ProductTypeFinished,
			OldValue:   product,
			ProductId:  product.Id,
			StoreId:    product.StoreId,
			SourceId:   product.Id,
			OperatorId: l.Staff.Id,
			IP:         l.Ctx.ClientIP(),
		}

		switch req.Type {
		case enums.ProductTypeUsedFinished:
			log.Action = enums.ProductActionDamageToNew
			product.Status = enums.ProductStatusNormal
			product.EnterTime = time.Now()
			if err := tx.Model(&model.ProductFinished{}).Where("id = ?", product.Id).Updates(model.ProductFinished{
				Status:    enums.ProductStatusNormal,
				EnterTime: time.Now(),
			}).Error; err != nil {
				return err
			}
			log.NewValue = product
		case enums.ProductTypeUsedOld:
			log.Action = enums.ProductActionDamageToOld
			if product.Status == enums.ProductStatusSold {
				log.Action = enums.ProductActionReturn
			}
			var old model.ProductOld
			if err := tx.Preload("Store").Where(&model.ProductOld{CodeFinished: strings.ToUpper(product.Code)}).First(&old).Error; err != nil {
				if err != gorm.ErrRecordNotFound {
					return errors.New("旧品不存在")
				}
			}

			data := model.ProductOld{
				Code:            strings.ToUpper("JL" + utils.RandomAlphanumericUpper(8)),
				CodeFinished:    strings.ToUpper(product.Code),
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

			data.Class = data.GetClass()

			if old.Id != "" {
				if old.DeletedAt.Valid {
					if err := tx.Model(&model.ProductOld{}).Unscoped().Where("id = ?", old.Id).Update("deleted_at", nil).Error; err != nil {
						return errors.New("恢复成品失败")
					}
				}
				if err := tx.Model(&model.ProductOld{}).Where("id = ?", old.Id).Updates(&data).Error; err != nil {
					return errors.New("转换失败")
				}
			} else {
				if err := tx.Create(&data).Error; err != nil {
					return errors.New("转换失败")
				}
			}
			data.Store = product.Store
			log.NewValue = data
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
