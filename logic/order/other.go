package order

import (
	"jdy/enums"
	"jdy/errors"
	"jdy/model"
	"jdy/types"
	"jdy/utils"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type OrderOtherLogic struct {
	Ctx   *gin.Context
	Staff *model.Staff
}

func (l *OrderOtherLogic) Create(req *types.OrderOtherCreateReq) (*model.OrderOther, error) {
	// 订单信息
	order := model.OrderOther{
		StoreId:    req.StoreId,
		Type:       req.Type,
		Content:    req.Content,
		Source:     req.Source,
		ClerkId:    req.ClerkId,
		MemberId:   req.MemberId,
		Amount:     req.Amount,
		OrderId:    req.OrderId,
		OperatorId: l.Staff.Id,
		IP:         l.Ctx.ClientIP(),
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {

		var expense decimal.Decimal
		for _, p := range req.Payments {
			expense = expense.Add(p.Amount)
			order.Payments = append(order.Payments, model.OrderPayment{
				Status:        true,
				StoreId:       order.StoreId,
				Type:          req.Type,
				Source:        enums.FinanceSource(req.Source),
				OrderType:     enums.OrderTypeOthers,
				PaymentMethod: p.PaymentMethod,
				Amount:        p.Amount,
			})
		}
		if expense.Cmp(order.Amount) != 0 {
			return errors.New("支付金额不正确")
		}

		// 保存订单
		if err := tx.Create(&order).Error; err != nil {
			return errors.New("创建订单失败")
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &order, nil
}

func (l *OrderOtherLogic) List(req *types.OrderOtherListReq) (*types.PageRes[model.OrderOther], error) {
	var (
		order model.OrderOther

		res types.PageRes[model.OrderOther]
	)

	db := model.DB.Model(&order)
	db = order.WhereCondition(db, &req.Where)

	// 获取总数
	if err := db.Count(&res.Total).Error; err != nil {
		return nil, errors.New("获取订单总数失败")
	}

	// 获取列表
	db = order.Preloads(db)

	db = db.Order("created_at desc")
	db = model.PageCondition(db, &req.PageReq)
	if err := db.Find(&res.List).Error; err != nil {
		return nil, errors.New("获取订单列表失败")
	}

	return &res, nil
}

func (l *OrderOtherLogic) Info(req *types.OrderOtherInfoReq) (*model.OrderOther, error) {
	var (
		order model.OrderOther
	)

	db := model.DB.Model(&order)
	db = order.Preloads(db)
	if err := db.First(&order, "id = ?", req.Id).Error; err != nil {
		return nil, errors.New("获取订单详情失败")
	}

	return &order, nil
}

func (l *OrderOtherLogic) Update(req *types.OrderOtherUpdateReq) (*model.OrderOther, error) {
	var (
		order model.OrderOther
	)

	db := model.DB.Model(&order)
	if err := db.First(&order, "id = ?", req.Id).Error; err != nil {
		return nil, errors.New("获取订单详情失败")
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		data, err := utils.StructToStruct[model.OrderOther](req)
		if err != nil {
			return errors.New("转换数据失败")
		}

		if err := tx.Model(&model.OrderOther{}).Where("id = ?", order.Id).Updates(data).Error; err != nil {
			return errors.New("更新订单失败")
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &order, nil
}

func (l *OrderOtherLogic) Delete(req *types.OrderOtherDeleteReq) error {
	var (
		order model.OrderOther
	)

	db := model.DB.Model(&order)
	if err := db.First(&order, "id = ?", req.Id).Error; err != nil {
		return errors.New("获取订单详情失败")
	}

	if err := db.Delete(&order).Error; err != nil {
		return errors.New("删除订单失败")
	}

	return nil
}
