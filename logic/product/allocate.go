package product

import (
	"jdy/enums"
	"jdy/errors"
	"jdy/model"
	"jdy/types"

	"gorm.io/gorm"
)

type ProductAllocateLogic struct {
	ProductLogic
}

// 创建产品调拨单
func (l *ProductAllocateLogic) Create(req *types.ProductAllocateCreateReq) *errors.Errors {
	// 开启事务
	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		data := model.ProductAllocate{
			Method: req.Method,
			Type:   req.Type,
			Reason: req.Reason,
			Remark: req.Remark,

			OperatorId: l.Staff.Id,
			IP:         l.Ctx.ClientIP(),
		}
		if req.Method == enums.ProductAllocateMethodStore {
			data.StoreId = req.StoreId
		}
		// 添加报损记录
		if err := tx.Create(&data).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return errors.New("创建产品调拨单失败")
	}
	return nil
}

// 获取产品调拨单列表
func (p *ProductAllocateLogic) List(req *types.ProductAllocateListReq) (*types.PageRes[model.ProductAllocate], error) {
	var (
		allocate model.ProductAllocate

		res types.PageRes[model.ProductAllocate]
	)

	db := model.DB.Model(&allocate)
	db = allocate.WhereCondition(db, &req.Where)

	// 获取总数
	if err := db.Count(&res.Total).Error; err != nil {
		return nil, errors.New("获取产品调拨列表失败: " + err.Error())
	}

	// 获取列表
	db = db.Order("created_at desc")
	db = model.PageCondition(db, req.Page, req.Limit)
	if err := db.Find(&res.List).Error; err != nil {
		return nil, errors.New("获取产品调拨列表失败: " + err.Error())
	}

	return &res, nil
}
