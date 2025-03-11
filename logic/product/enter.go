package product

import (
	"jdy/enums"
	"jdy/errors"
	"jdy/model"
	"jdy/types"
	"jdy/utils"
	"time"

	"gorm.io/gorm"
)

type ProductEnterLogic struct {
	ProductLogic
}

// 产品入库
func (l *ProductEnterLogic) Create(req *types.ProductEnterCreateReq) (*model.ProductEnter, error) {
	enter := model.ProductEnter{
		StoreId:    req.StoreId,
		Remark:     req.Remark,
		Status:     enums.ProductEnterStatusDraft,
		OperatorId: l.Staff.Id,
		IP:         l.Ctx.ClientIP(),
	}

	if err := model.DB.Create(&enter).Error; err != nil {
		return nil, errors.New("创建入库单失败")
	}

	return &enter, nil
}

// 产品入库单列表
func (l *ProductEnterLogic) EnterList(req *types.ProductEnterListReq) (*types.PageRes[model.ProductEnter], error) {
	var (
		enter model.ProductEnter

		res types.PageRes[model.ProductEnter]
	)

	db := model.DB.Model(&enter)
	db = enter.WhereCondition(db, &req.Where)

	// 获取总数
	if err := db.Count(&res.Total).Error; err != nil {
		return nil, errors.New("获取产品列表失败: " + err.Error())
	}

	// 获取列表
	db = db.Preload("Products")
	db = db.Preload("Operator")
	db = db.Preload("Store")

	db = db.Order("created_at desc")
	db = model.PageCondition(db, req.Page, req.Limit)
	if err := db.Find(&res.List).Error; err != nil {
		return nil, errors.New("获取产品列表失败: " + err.Error())
	}

	return &res, nil
}

// 产品入库单详情
func (l *ProductEnterLogic) EnterInfo(req *types.ProductEnterInfoReq) (*model.ProductEnter, error) {
	var (
		enter model.ProductEnter
	)

	db := model.DB.Model(&enter)

	// 获取产品入库单详情
	db = db.Preload("Products")
	db = db.Preload("Operator")
	db = db.Preload("Store")

	if err := db.First(&enter, req.Id).Error; err != nil {
		return nil, errors.New("获取产品入库单详情失败")
	}

	return &enter, nil
}

// 产品入库单添加产品
func (l *ProductEnterLogic) AddProduct(req *types.ProductEnterAddProductReq) (*map[string]bool, error) {
	// 查询入库单
	var enter model.ProductEnter
	if err := model.DB.Where("id = ?", req.ProductEnterId).First(&enter).Error; err != nil {
		return nil, errors.New("入库单不存在")
	}

	if enter.Status != enums.ProductEnterStatusDraft {
		return nil, errors.New("入库单已结束")
	}

	if len(req.Products) == 0 {
		return nil, errors.New("请选择产品")
	}

	// 添加产品的结果
	products := map[string]bool{}
	success := 0
	// 添加产品入库
	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		for _, p := range req.Products {
			// 转换数据结构
			product, err := utils.StructToStruct[model.Product](p)
			if err != nil {
				return errors.New("产品录入失败: 参数错误")
			}

			products[product.Code] = false

			var p model.Product
			if err := tx.Where("code = ?", product.Code).First(&p).Error; err != nil {
				if err != gorm.ErrRecordNotFound {
					return err
				}
			}
			if p.Id != "" {
				continue
			}

			// 产品入库
			product.ProductEnterId = enter.Id
			product.StoreId = enter.StoreId
			product.EnterTime = time.Now()
			// 商品状态
			product.Status = enums.ProductStatusDraft

			if err := tx.Create(&product).Error; err != nil {
				continue
			}

			products[product.Code] = true
			success++
		}

		if success == 0 {
			return errors.New("产品录入失败：无产品录入成功")
		}

		if success != len(req.Products) {
			return errors.New("部分产品录入失败")
		}

		return nil
	}); err != nil {
		return nil, errors.New("产品录入失败")
	}

	return &products, nil
}

// 入库单编辑产品
func (l *ProductEnterLogic) EditProduct(req *types.ProductEnterEditProductReq) error {
	// 查询入库单
	var enter model.ProductEnter
	if err := model.DB.Where("id = ?", req.ProductEnterId).First(&enter).Error; err != nil {
		return errors.New("入库单不存在")
	}

	if enter.Status != enums.ProductEnterStatusDraft {
		return errors.New("入库单已结束")
	}

	var product model.Product
	if err := model.DB.Where("id = ?", req.ProductId).First(&product).Error; err != nil {
		return errors.New("产品不存在")
	}

	if product.ProductEnterId != enter.Id && product.StoreId != enter.StoreId && product.Status != enums.ProductStatusDraft {
		return errors.New("产品不属于该入库单")
	}

	// 转换数据结构
	product, err := utils.StructToStruct[model.Product](req.Product)
	if err != nil {
		return errors.New("产品录入失败: 参数错误")
	}

	// 更新产品
	if err := model.DB.Model(model.Product{}).Where("id = ?", req.ProductId).Updates(product).Error; err != nil {
		return errors.New("产品更新失败")
	}

	return nil
}

// 入库单删除产品
func (l *ProductEnterLogic) DelProduct(req *types.ProductEnterDelProductReq) error {
	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 查询入库单
		var enter model.ProductEnter
		if err := tx.Where("id = ?", req.ProductEnterId).First(&enter).Error; err != nil {
			return errors.New("入库单不存在")
		}

		if enter.Status != enums.ProductEnterStatusDraft {
			return errors.New("入库单已结束")
		}

		// 查询产品
		for _, id := range req.ProductIds {
			var product model.Product
			if err := tx.Where("id = ?", id).First(&product).Error; err != nil {
				return errors.New("产品不存在")
			}

			if product.ProductEnterId != enter.Id && product.StoreId != enter.StoreId && product.Status != enums.ProductStatusDraft {
				return errors.New("产品不属于该入库单")
			}

			// 删除产品(真实删除)
			if err := tx.Where("id = ?", product.Id).Unscoped().Delete(&model.Product{}).Error; err != nil {
				return errors.New("产品删除失败")
			}
		}

		return nil
	}); err != nil {
		return errors.New("产品删除失败：" + err.Error())
	}
	return nil
}

// 入库单完成
func (l *ProductEnterLogic) Finish(req *types.ProductEnterFinishReq) error {
	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 查询入库单
		var enter model.ProductEnter
		if err := tx.Where("id = ?", req.ProductEnterId).Preload("Products").First(&enter).Error; err != nil {
			return errors.New("入库单不存在")
		}

		if enter.Status != enums.ProductEnterStatusDraft {
			return errors.New("入库单已结束")
		}

		if len(enter.Products) == 0 {
			return errors.New("入库单没有产品")
		}

		// 更新产品状态
		for _, product := range enter.Products {
			if product.Status != enums.ProductStatusDraft {
				return errors.New("产品状态不正确")
			}
			product.Status = enums.ProductStatusNormal
			if err := tx.Save(&product).Error; err != nil {
				return errors.New("产品更新失败")
			}
		}

		// 更新入库单状态
		enter.Status = enums.ProductEnterStatusCompleted
		if err := tx.Save(&enter).Error; err != nil {
			return errors.New("入库单更新失败")
		}

		return nil
	}); err != nil {
		return errors.New("操作失败：" + err.Error())
	}

	return nil
}

// 入库单取消
func (l *ProductEnterLogic) Cancel(req *types.ProductEnterCancelReq) error {
	// 查询入库单
	var enter model.ProductEnter
	if err := model.DB.Where("id = ?", req.ProductEnterId).Preload("Products").First(&enter).Error; err != nil {
		return errors.New("入库单不存在")
	}

	if enter.Status != enums.ProductEnterStatusDraft {
		return errors.New("入库单已结束")
	}

	// 更新入库单状态
	enter.Status = enums.ProductEnterStatusCanceled
	if err := model.DB.Save(&enter).Error; err != nil {
		return errors.New("入库单更新失败")
	}

	return nil
}
