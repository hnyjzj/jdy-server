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

type ProductAccessorieEnterLogic struct {
	Ctx   *gin.Context
	Staff *types.Staff
}

// 配件入库
func (l *ProductAccessorieEnterLogic) Create(req *types.ProductAccessorieEnterCreateReq) (*model.ProductAccessorieEnter, error) {
	enter := model.ProductAccessorieEnter{
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

// 配件入库单列表
func (l *ProductAccessorieEnterLogic) EnterList(req *types.ProductAccessorieEnterListReq) (*types.PageRes[model.ProductAccessorieEnter], error) {
	var (
		enter model.ProductAccessorieEnter

		res types.PageRes[model.ProductAccessorieEnter]
	)

	db := model.DB.Model(&enter)
	db = enter.WhereCondition(db, &req.Where)

	// 获取总数
	if err := db.Count(&res.Total).Error; err != nil {
		return nil, errors.New("获取配件列表失败: " + err.Error())
	}

	// 获取列表
	db = db.Preload("Products.Category")
	db = db.Preload("Operator")
	db = db.Preload("Store")

	db = db.Order("created_at desc")
	db = model.PageCondition(db, req.Page, req.Limit)
	if err := db.Find(&res.List).Error; err != nil {
		return nil, errors.New("获取配件列表失败: " + err.Error())
	}

	return &res, nil
}

// 配件入库单详情
func (l *ProductAccessorieEnterLogic) EnterInfo(req *types.ProductAccessorieEnterInfoReq) (*model.ProductAccessorieEnter, error) {
	var (
		enter model.ProductAccessorieEnter
	)

	db := model.DB.Model(&enter)

	// 获取配件入库单详情
	db = db.Preload("Products.Category")
	db = db.Preload("Operator")
	db = db.Preload("Store")

	if err := db.First(&enter, "id = ?", req.Id).Error; err != nil {
		return nil, errors.New("获取配件入库单详情失败")
	}

	return &enter, nil
}

// 配件入库单添加配件
func (l *ProductAccessorieEnterLogic) AddProduct(req *types.ProductAccessorieEnterAddProductReq) (*map[string]string, error) {
	// 查询入库单
	var enter model.ProductAccessorieEnter
	if err := model.DB.Where("id = ?", req.EnterId).First(&enter).Error; err != nil {
		return nil, errors.New("入库单不存在")
	}

	if enter.Status != enums.ProductEnterStatusDraft {
		return nil, errors.New("入库单已结束")
	}

	if len(req.Products) == 0 {
		return nil, errors.New("请选择配件")
	}

	// 添加配件的结果
	products := map[string]string{}
	success := 0
	// 添加配件入库
	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		for _, p := range req.Products {
			var category model.ProductAccessorieCategory
			if err := tx.Where("id = ?", p.Code).First(&category).Error; err != nil {
				products[p.Code] = "配件条目不存在"
				continue
			}

			var product model.ProductAccessorie
			if err := tx.Where(&model.ProductAccessorie{
				EnterId: enter.Id,
				Code:    p.Code,
			}).First(&product).Error; err != nil {
				if err != gorm.ErrRecordNotFound {
					return errors.New("配件条目错误")
				}
			}

			if product.Id != "" {
				// 配件已存在，更新配件
				if err := tx.Model(model.ProductAccessorie{}).Where("id = ?", product.Id).Updates(model.ProductAccessorie{
					Stock:     product.Stock + p.Stock,
					AccessFee: *p.AccessFee,
				}).Error; err != nil {
					products[p.Code] = "配件更新失败"
					continue
				}
			} else {
				data := model.ProductAccessorie{
					StoreId:   enter.StoreId,
					Code:      p.Code,
					Stock:     p.Stock,
					AccessFee: *p.AccessFee,
					Status:    enums.ProductStatusDraft,
					EnterId:   enter.Id,
				}

				if err := tx.Create(&data).Error; err != nil {
					products[p.Code] = "入库失败"
					continue
				}
			}

			success++
		}

		if success == 0 {
			return errors.New("无配件录入成功")
		}

		if success != len(req.Products) {
			return errors.New("部分配件录入失败")
		}

		return nil
	}); err != nil {
		return &products, errors.New("配件录入失败: " + err.Error())
	}

	return &products, nil
}

// 入库单编辑配件
func (l *ProductAccessorieEnterLogic) EditProduct(req *types.ProductAccessorieEnterEditProductReq) error {

	// 编辑配件逻辑存在潜在并发风险，使用事务包裹并加锁
	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 查询入库单
		var enter model.ProductAccessorieEnter
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", req.EnterId).First(&enter).Error; err != nil {
			return errors.New("入库单不存在")
		}

		if enter.Status != enums.ProductEnterStatusDraft {
			return errors.New("入库单已结束")
		}

		var product model.ProductAccessorie
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", req.ProductId).First(&product).Error; err != nil {
			return errors.New("配件不存在")
		}

		if product.EnterId != enter.Id || product.StoreId != enter.StoreId || product.Status != enums.ProductStatusDraft {
			return errors.New("配件不属于该入库单")
		}

		// 转换数据结构
		data, err := utils.StructToStruct[model.ProductAccessorie](req.Product)
		if err != nil {
			return errors.New("配件录入失败: 参数错误")
		}

		// 更新配件
		if err := tx.Model(model.ProductAccessorie{}).Where("id = ?", req.ProductId).Updates(data).Error; err != nil {
			return errors.New("配件更新失败")
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

// 入库单删除配件
func (l *ProductAccessorieEnterLogic) DelProduct(req *types.ProductAccessorieEnterDelProductReq) error {
	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 查询入库单
		var enter model.ProductAccessorieEnter
		if err := tx.Where("id = ?", req.EnterId).First(&enter).Error; err != nil {
			return errors.New("入库单不存在")
		}

		if enter.Status != enums.ProductEnterStatusDraft {
			return errors.New("入库单已结束")
		}

		// 查询配件
		for _, id := range req.ProductIds {
			var product model.ProductAccessorie
			if err := tx.Where("id = ?", id).First(&product).Error; err != nil {
				return errors.New("配件不存在")
			}

			if product.EnterId != enter.Id || product.StoreId != enter.StoreId || product.Status != enums.ProductStatusDraft {
				return errors.New("配件不属于该入库单")
			}

			// 删除配件(真实删除)
			if err := tx.Unscoped().Where("id = ?", product.Id).Delete(&model.ProductAccessorie{}).Error; err != nil {
				return errors.New("配件删除失败")
			}
		}

		return nil
	}); err != nil {
		return errors.New("配件删除失败：" + err.Error())
	}
	return nil
}

// 入库单完成
func (l *ProductAccessorieEnterLogic) Finish(req *types.ProductAccessorieEnterFinishReq) error {
	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 查询入库单
		var enter model.ProductAccessorieEnter
		if err := tx.Where("id = ?", req.EnterId).Preload("Products", func(tx *gorm.DB) *gorm.DB {
			tx = tx.Preload("Category")
			tx = tx.Preload("Store")
			return tx
		}).First(&enter).Error; err != nil {
			return errors.New("入库单不存在")
		}

		if enter.Status != enums.ProductEnterStatusDraft {
			return errors.New("入库单已结束")
		}

		if len(enter.Products) == 0 {
			return errors.New("入库单没有配件")
		}

		// 更新配件状态
		for _, product := range enter.Products {
			// 加锁查询配件
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				Where("id = ?", product.Id).
				Where(&model.ProductAccessorie{
					Status: enums.ProductStatusDraft,
				}).
				First(&product).Error; err != nil {
				return errors.New("配件状态不正确或配件不存在")
			}
			product.Status = enums.ProductStatusNormal
			if err := tx.Save(&product).Error; err != nil {
				return errors.New("配件更新失败")
			}

			// 添加记录
			history := model.ProductHistory{
				Type:       enums.ProductTypeAccessorie,
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
				return errors.New("配件记录添加失败")
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
func (l *ProductAccessorieEnterLogic) Cancel(req *types.ProductAccessorieEnterCancelReq) error {
	// 查询入库单
	var enter model.ProductAccessorieEnter
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
