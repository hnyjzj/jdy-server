package order

import (
	"jdy/errors"
	"jdy/model"
	"jdy/types"

	"github.com/gin-gonic/gin"
)

type OrderLogic struct {
	Ctx   *gin.Context
	Staff *types.Staff
}

func (l *OrderLogic) List(req *types.OrderListReq) (*types.PageRes[model.Order], error) {
	var (
		order model.Order

		res types.PageRes[model.Order]
	)

	db := model.DB.Model(&order)
	db = order.WhereCondition(db, &req.Where)

	// 获取总数
	if err := db.Count(&res.Total).Error; err != nil {
		return nil, errors.New("获取订单总数失败")
	}

	// 获取列表
	db = db.Preload("Member")
	db = db.Preload("Store")
	db = db.Preload("Cashier")
	db = db.Preload("Salesmens").Preload("Salesmens.Salesman")
	db = db.Preload("Products").Preload("Products.Product")

	db = db.Order("created_at desc")
	db = model.PageCondition(db, req.Page, req.Limit)
	if err := db.Find(&res.List).Error; err != nil {
		return nil, errors.New("获取订单列表失败")
	}

	return &res, nil
}

func (l *OrderLogic) Info(req *types.OrderInfoReq) (*model.Order, error) {
	var (
		order model.Order
	)

	db := model.DB.Model(&order)
	db = db.Preload("Member")
	db = db.Preload("Store")
	db = db.Preload("Cashier")
	db = db.Preload("Salesmens").Preload("Salesmens.Salesman")
	db = db.Preload("Products").Preload("Products.Product")

	db = db.Where("id = ?", req.Id)
	if err := db.First(&order).Error; err != nil {
		return nil, errors.New("获取订单详情失败")
	}

	return &order, nil
}
