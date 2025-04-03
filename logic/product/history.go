package product

import (
	"jdy/errors"
	"jdy/model"
	"jdy/types"

	"github.com/gin-gonic/gin"
)

type ProductHistoryLogic struct {
	Ctx   *gin.Context
	Staff *types.Staff
}

// 产品操作记录列表
func (l *ProductHistoryLogic) List(req *types.ProductHistoryListReq) (*types.PageRes[model.ProductHistory], error) {
	var (
		logs model.ProductHistory

		res types.PageRes[model.ProductHistory]
	)

	db := model.DB.Model(&logs)
	db = logs.WhereCondition(db, &req.Where)

	// 获取总数
	if err := db.Count(&res.Total).Error; err != nil {
		return nil, errors.New("获取总数失败")
	}

	// 获取列表
	db = db.Order("created_at desc")
	db = model.PageCondition(db, req.Page, req.Limit)
	db = db.Preload("Operator")

	if err := db.Find(&res.List).Error; err != nil {
		return nil, errors.New("获取列表失败")
	}

	return &res, nil
}

// 产品操作记录详情
func (l *ProductHistoryLogic) Info(req *types.ProductHistoryInfoReq) (*model.ProductHistory, error) {
	var (
		logs model.ProductHistory
	)

	db := model.DB.Model(&logs)
	db = db.Where("id = ?", req.Id)
	db = db.Preload("Operator")

	if err := db.First(&logs).Error; err != nil {
		return nil, errors.New("获取详情失败")
	}

	return &logs, nil
}
