package order

import (
	"jdy/errors"
	"jdy/model"
	"jdy/types"

	"github.com/gin-gonic/gin"
)

type OrderSalesDetailLogic struct {
	Ctx   *gin.Context
	Staff *types.Staff
}

func (l *OrderSalesDetailLogic) List(req *types.OrderSalesDetailListReq) (*types.PageRes[[]model.OrderSalesProduct], error) {
	var (
		product model.OrderSalesProduct

		res types.PageRes[[]model.OrderSalesProduct]
	)

	db := model.DB.Model(&product)
	db = product.WhereCondition(db, &req.Where)

	// 获取总数
	if err := db.Count(&res.Total).Error; err != nil {
		return nil, errors.New("获取订单总数失败")
	}

	// 获取列表
	db = product.Preloads(db)

	db = db.Order("created_at desc")
	db = model.PageCondition(db, req.Page, req.Limit)
	if err := db.Find(&res.List).Error; err != nil {
		return nil, errors.New("获取订单列表失败")
	}

	return &res, nil
}

func (l *OrderSalesDetailLogic) Info(req *types.OrderSalesDetailInfoReq) (*model.OrderSalesProduct, error) {
	var (
		product model.OrderSalesProduct
	)

	db := model.DB.Model(&model.OrderSalesProduct{})

	db = product.Preloads(db)

	db = db.Where("id = ?", req.Id)
	if err := db.First(&product).Error; err != nil {
		return nil, errors.New("获取订单详情失败")
	}

	return &product, nil
}
