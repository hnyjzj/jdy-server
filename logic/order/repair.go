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

type OrderRepairLogic struct {
	Ctx   *gin.Context
	Staff *model.Staff
}

func (l *OrderRepairLogic) Create(req *types.OrderRepairCreateReq) (*model.OrderRepair, error) {
	// 订单信息
	order := model.OrderRepair{
		StoreId:        req.StoreId,
		Status:         enums.OrderRepairStatusWaitPay,
		ReceptionistId: req.ReceptionistId,
		CashierId:      req.CashierId,
		MemberId:       req.MemberId,
		Name:           req.Name,
		Desc:           req.Desc,
		DeliveryMethod: req.DeliveryMethod,
		Province:       req.Province,
		City:           req.City,
		Area:           req.Area,
		Address:        req.Address,
		Images:         req.Images,
		OperatorId:     l.Staff.Id,
		IP:             l.Ctx.ClientIP(),
		Expense:        req.Expense,
		Cost:           req.Cost,
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 商品
		for _, p := range req.Products {
			data := model.OrderRepairProduct{
				Status:      enums.OrderRepairStatusWaitPay,
				IsOur:       p.IsOur,
				Code:        p.Code,
				Name:        p.Name,
				LabelPrice:  p.LabelPrice,
				Brand:       p.Brand,
				Material:    p.Material,
				Quality:     p.Quality,
				Gem:         p.Gem,
				Category:    p.Category,
				Craft:       p.Craft,
				WeightMetal: p.WeightMetal,
				WeightTotal: p.WeightTotal,
				ColorGem:    p.ColorGem,
				WeightGem:   p.WeightGem,
				Clarity:     p.ClarityGem,
				Cut:         p.Cut,
				Remark:      p.Remark,
			}

			if p.IsOur {
				var product model.ProductFinished
				if err := tx.Model(&product).Where(&model.ProductFinished{
					StoreId: order.StoreId,
				}).First(&product, "id = ?", p.ProductId).Error; err != nil {
					return errors.New("获取商品信息失败")
				}
				if product.Status != enums.ProductStatusSold && product.Status != enums.ProductStatusNoStock {
					return errors.New("商品状态不正常")
				}

				data.ProductId = product.Id
			}

			order.Products = append(order.Products, data)
		}

		var expense decimal.Decimal
		for _, p := range req.Payments {
			expense = expense.Add(p.Amount)
			order.Payments = append(order.Payments, model.OrderPayment{
				StoreId:       order.StoreId,
				Type:          enums.FinanceTypeIncome,
				Source:        enums.FinanceSourceOtherReturn,
				PaymentMethod: p.PaymentMethod,
				Amount:        p.Amount,
				OrderType:     enums.OrderTypeRepair,
			})
		}
		if expense.Cmp(order.Expense) != 0 {
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

func (l *OrderRepairLogic) List(req *types.OrderRepairListReq) (*types.PageRes[model.OrderRepair], error) {
	var (
		order model.OrderRepair

		res types.PageRes[model.OrderRepair]
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
	db = model.PageCondition(db, req.Page, req.Limit)
	if err := db.Find(&res.List).Error; err != nil {
		return nil, errors.New("获取订单列表失败")
	}

	return &res, nil
}

func (l *OrderRepairLogic) Info(req *types.OrderRepairInfoReq) (*model.OrderRepair, error) {
	var (
		order model.OrderRepair
	)

	db := model.DB.Model(&order)
	db = order.Preloads(db)
	if err := db.First(&order, "id = ?", req.Id).Error; err != nil {
		return nil, errors.New("获取订单详情失败")
	}

	return &order, nil
}

func (l *OrderRepairLogic) Update(req *types.OrderRepairUpdateReq) error {
	var (
		order model.OrderRepair
	)

	db := model.DB.Model(&order)
	db = order.Preloads(db)
	if err := db.First(&order, "id = ?", req.Id).Error; err != nil {
		return errors.New("获取订单详情失败")
	}

	if order.Status == enums.OrderRepairStatusComplete || order.Status == enums.OrderRepairStatusCancel {
		return errors.New("订单状态不正确")
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 更新订单状态
		data, err := utils.StructToStruct[model.OrderRepair](req)
		if err != nil {
			return errors.New("验证信息失败")
		}

		if err := tx.Model(&model.OrderRepair{}).Where("id = ?", req.Id).Updates(data).Error; err != nil {
			return errors.New("更新失败")
		}

		return nil
	}); err != nil {
		return errors.New("更新订单状态失败")
	}

	return nil
}

// 操作
func (l *OrderRepairLogic) Operation(req *types.OrderRepairOperationReq) error {
	var (
		order model.OrderRepair
	)

	db := model.DB.Model(&order)
	db = order.Preloads(db)
	if err := db.First(&order, "id = ?", req.Id).Error; err != nil {
		return errors.New("获取订单详情失败")
	}

	// 判断状态
	if !order.Status.CanOperationTo(req.Operation) {
		return errors.New("订单状态不正确")
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 支付产品
		for _, product := range order.Products {
			// 更新订单状态
			if err := tx.Model(&product).Updates(&model.OrderRepairProduct{
				Status: req.Operation,
			}).Error; err != nil {
				return errors.New("更新订单产品状态失败")
			}
		}
		// 更新订单状态
		if err := tx.Model(&order).Updates(&model.OrderRepair{
			Status: req.Operation,
		}).Error; err != nil {
			return errors.New("更新订单状态失败")
		}

		return nil
	}); err != nil {
		return errors.New("更新订单状态失败")
	}

	return nil
}

// 撤销
func (l *OrderRepairLogic) Revoked(req *types.OrderRepairRevokedReq) error {
	var (
		order model.OrderRepair
	)

	db := model.DB.Model(&order)

	db = order.Preloads(db)
	if err := db.First(&order, "id = ?", req.Id).Error; err != nil {
		return errors.New("获取订单详情失败")
	}

	if order.Status != enums.OrderRepairStatusWaitPay {
		return errors.New("订单状态不正确")
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 撤销产品
		for _, product := range order.Products {
			// 更新订单状态
			if err := tx.Model(&product).Updates(&model.OrderRepairProduct{
				Status: enums.OrderRepairStatusCancel,
			}).Error; err != nil {
				return errors.New("更新订单产品状态失败")
			}
		}

		// 更新订单状态
		if err := tx.Model(&order).Updates(&model.OrderRepair{
			Status: enums.OrderRepairStatusCancel,
		}).Error; err != nil {
			return errors.New("撤销订单失败")
		}

		return nil
	}); err != nil {
		return errors.New("撤销订单失败")
	}

	return nil
}

func (l *OrderRepairLogic) Pay(req *types.OrderRepairPayReq) error {
	var (
		order model.OrderRepair
	)

	db := model.DB.Model(&order)
	db = order.Preloads(db)
	if err := db.First(&order, "id = ?", req.Id).Error; err != nil {
		return errors.New("获取订单详情失败")
	}

	inSuperiors := false
	for _, superior := range order.Store.Superiors {
		if superior.Id == l.Staff.Id {
			inSuperiors = true
			break
		}
	}
	if order.CashierId != l.Staff.Id && !inSuperiors {
		return errors.New("订单不是当前收银员操作")
	}
	if order.Status != enums.OrderRepairStatusWaitPay {
		return errors.New("订单状态不正确")
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 支付产品
		for _, product := range order.Products {
			// 更新订单状态
			if err := tx.Model(&product).Updates(&model.OrderRepairProduct{
				Status: enums.OrderRepairStatusStoreReceived,
			}).Error; err != nil {
				return errors.New("更新订单产品状态失败")
			}
		}

		// 更新订单状态
		if err := tx.Model(&order).Updates(&model.OrderRepair{
			Status: enums.OrderRepairStatusStoreReceived,
		}).Error; err != nil {
			return errors.New("支付订单失败")
		}

		return nil
	}); err != nil {
		return errors.New("支付订单失败")
	}

	return nil
}

// 退款
func (l *OrderRepairLogic) Refund(req *types.OrderRepairRefundReq) error {
	var (
		order model.OrderRepair
	)

	db := model.DB.Model(&order)

	db = order.Preloads(db)
	if err := db.First(&order, "id = ?", req.Id).Error; err != nil {
		return errors.New("获取订单详情失败")
	}

	if order.Status != enums.OrderRepairStatusStoreReceived {
		return errors.New("订单状态不正确")
	}

	data := model.OrderRefund{
		StoreId:    order.StoreId,
		OrderId:    order.Id,
		OrderType:  enums.OrderTypeRepair,
		MemberId:   order.MemberId,
		Remark:     req.Remark,
		OperatorId: l.Staff.Id,
		IP:         l.Ctx.ClientIP(),
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {

		for _, product := range order.Products {
			// 更新订单状态
			if err := tx.Model(&product).Updates(&model.OrderRepairProduct{
				Status: enums.OrderRepairStatusRefund,
			}).Error; err != nil {
				return errors.New("更新订单产品状态失败")
			}
		}

		if err := tx.Model(&order).Updates(&model.OrderRepair{
			Status: enums.OrderRepairStatusRefund,
		}).Error; err != nil {
			return errors.New("更新订单状态失败")
		}

		if err := tx.Create(&data).Error; err != nil {
			return errors.New("创建退款记录失败")
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
