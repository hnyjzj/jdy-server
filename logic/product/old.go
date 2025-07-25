package product

import (
	"jdy/enums"
	"jdy/errors"
	"jdy/model"
	"jdy/types"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductOldLogic struct {
	Ctx   *gin.Context
	Staff *model.Staff
}

// 旧料列表
func (p *ProductOldLogic) List(req *types.ProductOldListReq) (*types.PageRes[model.ProductOld], error) {
	var (
		product model.ProductOld

		res types.PageRes[model.ProductOld]
	)

	db := model.DB.Model(&product)
	db = product.WhereCondition(db, &req.Where)

	// 获取总数
	if err := db.Count(&res.Total).Error; err != nil {
		return nil, errors.New("获取旧料列表失败")
	}

	// 获取列表
	db = db.Order("created_at desc")
	db = model.PageCondition(db, &req.PageReq)
	if err := db.Find(&res.List).Error; err != nil {
		return nil, errors.New("获取旧料列表失败")
	}

	return &res, nil
}

// 旧料详情
func (p *ProductOldLogic) Info(req *types.ProductOldInfoReq) (*model.ProductOld, error) {
	var (
		product model.ProductOld
	)

	db := model.DB.Model(&model.ProductOld{})
	db = db.Preload("Store")
	if req.Id != "" {
		db = db.Where("id = ?", req.Id)
	}
	if req.Code != "" {
		db = db.Where(&model.ProductOld{Code: strings.ToUpper(req.Code)})
	}

	if err := db.First(&product).Error; err != nil {
		return nil, errors.New("获取旧料信息失败")
	}

	return &product, nil
}

// 旧料转换
func (l *ProductOldLogic) Conversion(req *types.ProductConversionReq) *errors.Errors {
	// 开启事务
	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 查询商品信息
		var old_product model.ProductOld
		if err := tx.Unscoped().Preload("Store").Where("id = ?", req.Id).First(&old_product).Error; err != nil {
			return errors.New("商品不存在")
		}

		// 判断旧料状态
		if !old_product.IsOur {
			return errors.New("非自有旧料无法转换")
		}

		log := &model.ProductHistory{
			Action:     enums.ProductActionOldToNew,
			Type:       enums.ProductTypeFinished,
			OldValue:   old_product,
			ProductId:  old_product.Id,
			StoreId:    old_product.StoreId,
			SourceId:   old_product.Id,
			OperatorId: l.Staff.Id,
			IP:         l.Ctx.ClientIP(),
		}

		// 转换
		var finished_product model.ProductFinished
		if err := tx.Unscoped().Preload("Store").Where("code = ?", strings.ToUpper(old_product.Code)).First(&finished_product).Error; err != nil {
			return errors.New("成品不在库中")
		}
		if finished_product.Status != enums.ProductStatusDamage && finished_product.Status != enums.ProductStatusSold { // 判断成品状态
			return errors.New("成品不在库中")
		}
		// 更新成品状态,如果被删除了，则恢复
		if finished_product.DeletedAt.Valid {
			if err := tx.Model(&model.ProductFinished{}).Unscoped().Where("id = ?", finished_product.Id).Update("deleted_at", nil).Error; err != nil {
				return errors.New("恢复成品失败")
			}
		}

		// 更新旧料状态
		if err := tx.Model(&model.ProductFinished{}).Where("id = ?", finished_product.Id).Update("status", enums.ProductStatusNormal).Error; err != nil {
			return errors.New("更新成品状态失败")
		}

		// 添加日志
		log.NewValue = finished_product
		if err := tx.Create(log).Error; err != nil {
			return errors.New("添加日志失败")
		}

		// 删除旧料
		if err := tx.Where("id = ?", old_product.Id).Delete(&old_product).Error; err != nil {
			return errors.New("删除旧料失败")
		}

		return nil
	}); err != nil {
		return errors.New("转换失败：" + err.Error())
	}

	return nil
}

// 获取大类
func (p *ProductOldLogic) GetClass(req *types.ProductOldGetClassReq) types.ProductOldGetClassRes {
	old := model.ProductOld{
		Material: req.Material,
		Quality:  req.Quality,
		Gem:      req.Gem,
	}
	class := old.GetClass()

	return types.ProductOldGetClassRes{
		Value: class,
		Label: enums.ProductClassOldMap[class],
	}
}
