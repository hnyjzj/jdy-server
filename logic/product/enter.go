package product

import (
	"jdy/enums"
	"jdy/errors"
	"jdy/model"
	"jdy/types"
	"jdy/utils"

	"gorm.io/gorm"
)

// 产品入库
func (l *ProductLogic) Enter(req *types.ProductEnterReq) (*map[string]bool, *errors.Errors) {
	// 添加产品的结果
	products := map[string]bool{}
	success := 0
	// 添加产品入库
	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 转换数据结构
		data, err := utils.StructToStruct[[]model.Product](req.Products)
		if err != nil {
			return nil
		}

		if len(data) == 0 {
			return errors.New("产品录入失败")
		}

		enter := model.ProductEnter{
			OperatorId: l.Staff.Id,
			IP:         l.Ctx.ClientIP(),
		}
		if err := tx.Create(&enter).Error; err != nil {
			return err
		}

		for _, v := range data {
			products[v.Code] = false

			var p model.Product
			if err := tx.Where("code = ?", v.Code).First(&p).Error; err != nil {
				if err != gorm.ErrRecordNotFound {
					return err
				}
			}
			if p.Id != "" {
				continue
			}

			// 产品入库
			v.Status = enums.ProductStatusNormal
			v.ProductEnterId = enter.Id
			if v.Stock == 0 {
				v.Stock = 1
			}

			if err := tx.Create(&v).Error; err != nil {
				continue
			}

			products[v.Code] = true
			success++
		}

		if success == 0 {
			return errors.New("产品录入失败")
		}

		if success != len(data) {
			return errors.New("部分产品录入失败")
		}

		return nil
	}); err != nil {
		return nil, errors.New("产品录入失败")
	}

	return &products, nil
}

// 产品入库单列表
func (l *ProductLogic) EnterList(req *types.ProductEnterListReq) (*types.PageRes[model.ProductEnter], error) {
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
	db = db.Order("created_at desc")
	db = model.PageCondition(db, req.Page, req.Limit)
	if err := db.Find(&res.List).Error; err != nil {
		return nil, errors.New("获取产品列表失败: " + err.Error())
	}

	return &res, nil
}

// 产品入库单详情
func (l *ProductLogic) EnterInfo(req *types.ProductEnterInfoReq) (*model.ProductEnter, error) {
	var (
		enter model.ProductEnter
	)

	db := model.DB.Model(&enter)

	// 获取产品入库单详情
	db = db.Preload("Products")
	db = db.Preload("Operator")

	if err := db.First(&enter, req.Id).Error; err != nil {
		return nil, errors.New("获取产品入库单详情失败")
	}

	return &enter, nil
}
