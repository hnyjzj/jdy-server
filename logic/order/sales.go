package order

import (
	"jdy/errors"
	"jdy/model"
	"jdy/types"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderSalesLogic struct {
	Ctx   *gin.Context
	Staff *types.Staff
}

func (l *OrderSalesLogic) List(req *types.OrderSalesListReq) (*types.PageRes[model.OrderSales], error) {
	var (
		order model.OrderSales

		res types.PageRes[model.OrderSales]
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
	db = db.Preload("Clerks", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Salesman")
	})
	db = db.Preload("ProductFinisheds", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Product")
	})
	db = db.Preload("ProductOlds", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Product")
	})
	db = db.Preload("ProductAccessories", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Product", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Category")
		})
	})
	db = db.Preload("Payments")

	db = db.Order("created_at desc")
	db = model.PageCondition(db, req.Page, req.Limit)
	if err := db.Find(&res.List).Error; err != nil {
		return nil, errors.New("获取订单列表失败")
	}

	return &res, nil
}

func (l *OrderSalesLogic) Info(req *types.OrderSalesInfoReq) (*model.OrderSales, error) {
	var (
		order model.OrderSales
	)

	db := model.DB.Model(&order)

	db = db.Preload("Member")
	db = db.Preload("Store")
	db = db.Preload("Cashier")
	db = db.Preload("Clerks", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Salesman")
	})
	db = db.Preload("ProductFinisheds", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Product")
	})
	db = db.Preload("ProductOlds", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Product")
	})
	db = db.Preload("ProductAccessories", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Product", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Category")
		})
	})
	db = db.Preload("OrderDeposits", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Products", func(db *gorm.DB) *gorm.DB {
			return db.Preload("ProductFinished")
		})
	})
	db = db.Preload("Payments")

	db = db.Where("id = ?", req.Id)
	if err := db.First(&order).Error; err != nil {
		return nil, errors.New("获取订单详情失败")
	}

	return &order, nil
}
