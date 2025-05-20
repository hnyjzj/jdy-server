package product

import (
	"jdy/enums"
	"jdy/errors"
	"jdy/model"
	"jdy/types"
	"jdy/utils"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProductFinishedEnterLogic struct {
	Ctx   *gin.Context
	Staff *types.Staff
}

// 产品入库
func (l *ProductFinishedEnterLogic) Create(req *types.ProductFinishedEnterCreateReq) (*model.ProductFinishedEnter, error) {
	enter := model.ProductFinishedEnter{
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
func (l *ProductFinishedEnterLogic) EnterList(req *types.ProductFinishedEnterListReq) (*types.PageRes[model.ProductFinishedEnter], error) {
	var (
		enter model.ProductFinishedEnter

		res types.PageRes[model.ProductFinishedEnter]
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
func (l *ProductFinishedEnterLogic) EnterInfo(req *types.ProductFinishedEnterInfoReq) (*model.ProductFinishedEnter, error) {
	var (
		enter model.ProductFinishedEnter
	)

	db := model.DB.Model(&enter)

	// 获取产品入库单详情
	db = db.Preload("Products", func(tx *gorm.DB) *gorm.DB {
		return tx.Unscoped()
	})
	db = db.Preload("Operator")
	db = db.Preload("Store")

	if err := db.First(&enter, "id = ?", req.Id).Error; err != nil {
		return nil, errors.New("获取产品入库单详情失败")
	}

	return &enter, nil
}

// 产品入库单添加产品
func (l *ProductFinishedEnterLogic) AddProduct(req *types.ProductFinishedEnterAddProductReq) (*map[string]string, error) {
	// 查询入库单
	var enter model.ProductFinishedEnter
	if err := model.DB.Where("id = ?", req.EnterId).First(&enter).Error; err != nil {
		return nil, errors.New("入库单不存在")
	}

	if enter.Status != enums.ProductEnterStatusDraft {
		return nil, errors.New("入库单已结束")
	}

	if len(req.Products) == 0 {
		return nil, errors.New("请选择产品")
	}

	// 添加产品的结果
	products := map[string]string{}
	success := 0
	// 添加产品入库
	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		for _, p := range req.Products {
			// 转换数据结构
			product, err := utils.StructToStruct[model.ProductFinished](p)
			if err != nil {
				return errors.New("参数错误")
			}

			var p model.ProductFinished
			if err := tx.Where("code = ?", product.Code).First(&p).Error; err != nil {
				if err != gorm.ErrRecordNotFound {
					return errors.New("产品不存在")
				}
			}
			if p.Id != "" {
				products[product.Code] = "条码已存在"
				continue
			}

			// 产品信息
			product.EnterId = enter.Id
			product.StoreId = enter.StoreId
			product.Class = product.GetClass()
			product.Status = enums.ProductStatusDraft
			product.EnterTime = time.Now()

			if err := tx.Create(&product).Error; err != nil {
				products[product.Code] = "入库失败"
				continue
			}

			success++
		}

		if success == 0 {
			return errors.New("无产品录入成功")
		}

		if success != len(req.Products) {
			return errors.New("部分产品录入失败")
		}

		return nil
	}); err != nil {
		return &products, errors.New("产品录入失败: " + err.Error())
	}

	return &products, nil
}

// 入库单编辑产品
func (l *ProductFinishedEnterLogic) EditProduct(req *types.ProductFinishedEnterEditProductReq) error {

	// 编辑产品逻辑存在潜在并发风险，使用事务包裹并加锁
	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 查询入库单
		var enter model.ProductFinishedEnter
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", req.EnterId).First(&enter).Error; err != nil {
			return errors.New("入库单不存在")
		}

		if enter.Status != enums.ProductEnterStatusDraft {
			return errors.New("入库单已结束")
		}

		var product model.ProductFinished
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", req.ProductId).First(&product).Error; err != nil {
			return errors.New("产品不存在")
		}

		if product.EnterId != enter.Id || product.StoreId != enter.StoreId || product.Status != enums.ProductStatusDraft {
			return errors.New("产品不属于该入库单")
		}

		// 转换数据结构
		product, err := utils.StructToStruct[model.ProductFinished](req.Product)
		if err != nil {
			return errors.New("产品录入失败: 参数错误")
		}

		// 更新产品
		if err := tx.Model(model.ProductFinished{}).Where("id = ?", req.ProductId).Updates(product).Error; err != nil {
			return errors.New("产品更新失败")
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

// 入库单删除产品
func (l *ProductFinishedEnterLogic) DelProduct(req *types.ProductFinishedEnterDelProductReq) error {
	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 查询入库单
		var enter model.ProductFinishedEnter
		if err := tx.Where("id = ?", req.EnterId).First(&enter).Error; err != nil {
			return errors.New("入库单不存在")
		}

		if enter.Status != enums.ProductEnterStatusDraft {
			return errors.New("入库单已结束")
		}

		// 查询产品
		for _, id := range req.ProductIds {
			var product model.ProductFinished
			if err := tx.Where("id = ?", id).First(&product).Error; err != nil {
				return errors.New("产品不存在")
			}

			if product.EnterId != enter.Id || product.StoreId != enter.StoreId || product.Status != enums.ProductStatusDraft {
				return errors.New("产品不属于该入库单")
			}

			// 删除产品(真实删除)
			if err := tx.Where("id = ?", product.Id).Unscoped().Delete(&model.ProductFinished{}).Error; err != nil {
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
func (l *ProductFinishedEnterLogic) Finish(req *types.ProductFinishedEnterFinishReq) error {
	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 查询入库单
		var enter model.ProductFinishedEnter
		if err := tx.Where("id = ?", req.EnterId).Preload("Products", func(tx *gorm.DB) *gorm.DB {
			return tx.Preload("Store")
		}).First(&enter).Error; err != nil {
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
			// 加锁
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ? AND status = ?", product.Id, enums.ProductStatusDraft).First(&product).Error; err != nil {
				return errors.New("产品状态不正确或产品不存在")
			}
			product.Status = enums.ProductStatusNormal
			if err := tx.Save(&product).Error; err != nil {
				return errors.New("产品更新失败")
			}

			// 添加记录
			history := model.ProductHistory{
				Type:       enums.ProductTypeFinished,
				OldValue:   nil,
				NewValue:   product,
				Action:     enums.ProductActionEntry,
				ProductId:  product.Id,
				StoreId:    enter.StoreId,
				SourceId:   enter.Id,
				OperatorId: l.Staff.Id,
				IP:         l.Ctx.ClientIP(),
			}
			if err := tx.Save(&history).Error; err != nil {
				return errors.New("产品记录添加失败")
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
func (l *ProductFinishedEnterLogic) Cancel(req *types.ProductFinishedEnterCancelReq) error {
	// 查询入库单
	var enter model.ProductFinishedEnter
	if err := model.DB.Where("id = ?", req.EnterId).Preload("Products").First(&enter).Error; err != nil {
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
