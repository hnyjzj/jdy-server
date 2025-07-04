package order

import (
	"jdy/errors"
	"jdy/model"
	"jdy/types"

	"github.com/gin-gonic/gin"
)

type OrderSalesRefundLogic struct {
	Ctx   *gin.Context
	Staff *model.Staff
}

func (l *OrderSalesRefundLogic) List(req *types.OrderSalesRefundListReq) (*types.PageRes[model.OrderRefund], error) {
	var (
		product model.OrderRefund

		res types.PageRes[model.OrderRefund]
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
