package statistic

import (
	"errors"
	"jdy/enums"
	"jdy/model"
	"jdy/types"
	"log"
	"strings"

	"gorm.io/gorm"
)

// 销售明细日报
func (l *StatisticLogic) SalesDetailDaily(req *types.StatisticSalesDetailDailyReq) (*types.StatisticSalesDetailDailyResp, error) {
	var (
		logic = &StatisticSalesDetailDailyLogic{
			p:    l,
			req:  req,
			resp: &types.StatisticSalesDetailDailyResp{},
		}
	)

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		logic.tx = tx

		if err := logic.getOrderSales(); err != nil {
			return err
		}

		if err := logic.getOrderDeposit(); err != nil {
			return err
		}

		if err := logic.getOrderOther(); err != nil {
			return err
		}

		// 开始统计
		logic.getSummary()
		logic.getItemized()
		logic.getPayment()
		logic.getFinishedSales()
		logic.getAccessorieSales()

		return nil
	}); err != nil {
		log.Printf("销售明细日报统计失败: %v", err)
		return nil, errors.New("统计失败")
	}

	return logic.resp, nil
}

type StatisticSalesDetailDailyLogic struct {
	p    *StatisticLogic
	tx   *gorm.DB
	req  *types.StatisticSalesDetailDailyReq
	resp *types.StatisticSalesDetailDailyResp

	OrderSales   []model.OrderSales
	OrderDeposit []model.OrderDeposit
	OrderOther   []model.OrderOther
}

// 获取销售单
func (l *StatisticSalesDetailDailyLogic) getOrderSales() error {

	db := l.tx.Model(&model.OrderSales{})

	if l.req.StoreId != "" {
		db = db.Where("store_id = ?", l.req.StoreId)
	}
	if l.req.SalesmanId != "" {
		db = db.Where("id in (select order_id from order_sales_clerks where salesman_id = ?)", l.req.SalesmanId)
	}
	if l.req.StartTime != nil {
		db = db.Where("created_at >= ?", l.req.StartTime)
	}
	if l.req.EndTime != nil {
		db = db.Where("created_at <= ?", l.req.EndTime)
	}

	db = model.OrderSales{}.Preloads(db)
	db = db.Where("status in (?)", []enums.OrderSalesStatus{
		enums.OrderSalesStatusComplete,
		enums.OrderSalesStatusRefund,
		enums.OrderSalesStatusReturn,
	})

	if err := db.Find(&l.OrderSales).Error; err != nil {
		return err
	}

	return nil
}

// 获取订金单
func (l *StatisticSalesDetailDailyLogic) getOrderDeposit() error {

	db := l.tx.Model(&model.OrderDeposit{})

	if l.req.StoreId != "" {
		db = db.Where("store_id = ?", l.req.StoreId)
	}
	if l.req.SalesmanId != "" {
		db = db.Where(&model.OrderDeposit{ClerkId: l.req.SalesmanId})
	}
	if l.req.StartTime != nil {
		db = db.Where("created_at >= ?", l.req.StartTime)
	}
	if l.req.EndTime != nil {
		db = db.Where("created_at <= ?", l.req.EndTime)
	}

	db = model.OrderDeposit{}.Preloads(db)
	db = db.Where("status in (?)", []enums.OrderDepositStatus{
		enums.OrderDepositStatusComplete,
		enums.OrderDepositStatusRefund,
		enums.OrderDepositStatusReturn,
	})

	if err := db.Find(&l.OrderDeposit).Error; err != nil {
		return err
	}

	return nil
}

// 获取其他收支单
func (l *StatisticSalesDetailDailyLogic) getOrderOther() error {

	db := l.tx.Model(&model.OrderOther{})

	if l.req.StoreId != "" {
		db = db.Where("store_id = ?", l.req.StoreId)
	}
	if l.req.SalesmanId != "" {
		db = db.Where(&model.OrderOther{ClerkId: l.req.SalesmanId})
	}
	if l.req.StartTime != nil {
		db = db.Where("created_at >= ?", l.req.StartTime)
	}
	if l.req.EndTime != nil {
		db = db.Where("created_at <= ?", l.req.EndTime)
	}

	db = model.OrderOther{}.Preloads(db)

	if err := db.Find(&l.OrderOther).Error; err != nil {
		return err
	}

	return nil
}

// 获取汇总项
func (l *StatisticSalesDetailDailyLogic) getSummary() {
	// 汇总
	res := types.StatisticSalesDetailDailySummary{}

	for _, order := range l.OrderSales {
		// 销售应收
		res.SalesReceivable = res.SalesReceivable.Add(order.Price)
		// 退货
		if order.Status == enums.OrderSalesStatusRefund {
			for _, refund := range order.Products {
				if refund.Status == enums.OrderSalesStatusRefund {
					switch refund.Type {
					case enums.ProductTypeFinished:
						{
							// 销售退款
							res.SalesRefund = res.SalesRefund.Add(refund.Finished.Price)
						}
					case enums.ProductTypeOld:
						{
							// 销售退款
							res.SalesRefund = res.SalesRefund.Add(refund.Old.RecyclePrice)
						}
					case enums.ProductTypeAccessorie:
						{
							// 销售退款
							res.SalesRefund = res.SalesRefund.Add(refund.Accessorie.Price)
						}
					}
				}
			}
		}
		// 销售实收
		res.SalesReceived = res.SalesReceived.Add(order.PricePay)
	}

	for _, order := range l.OrderDeposit {
		// 订金收入
		res.DepositIncome = res.DepositIncome.Add(order.PricePay)
		// 订金退款
		if order.Status == enums.OrderDepositStatusRefund {
			res.DepositRefund = res.DepositRefund.Add(order.PricePay)
		}
	}

	for _, order := range l.OrderOther {
		// 其他收入
		if order.Type == enums.FinanceTypeIncome {
			res.OtherIncome = res.OtherIncome.Add(order.Amount)
		}
		// 其他支出
		if order.Type == enums.FinanceTypeExpense {
			res.OtherExpense = res.OtherExpense.Add(order.Amount)
		}
	}

	l.resp.Summary = res
}

// 获取分项
func (l *StatisticSalesDetailDailyLogic) getItemized() {
	res := types.StatisticSalesDetailDailyItemized{}

	for _, order := range l.OrderSales {
		// 销售产品
		for _, product := range order.Products {
			switch product.Type {
			case enums.ProductTypeFinished:
				{
					// 成品应收
					res.FinishedReceivable = res.FinishedReceivable.Add(product.Finished.Price)
					// 成品件数
					res.FinishedQuantity++
				}
			case enums.ProductTypeOld:
				{
					// 旧料抵值
					res.OldDeduction = res.OldDeduction.Add(product.Old.RecyclePrice)
					// 旧料件数
					res.OldQuantity++
					// 旧料金重
					res.OldWeightMetal = res.OldWeightMetal.Add(product.Old.WeightMetal)

					// 旧料转成品
					if product.Old.Product.DeletedAt.Valid {
						// 旧料转成品件数
						res.OldToFinishedQuantity++
						// 旧料转成品金额
						res.OldToFinishedDeduction = res.OldToFinishedDeduction.Add(product.Old.RecyclePrice)
						// 旧料转成品重量
						res.OldToFinishedWeightMetal = res.OldToFinishedWeightMetal.Add(product.Old.WeightMetal)
					}
				}

			case enums.ProductTypeAccessorie:
				{
					// 配件金额
					res.AccessoriePrice = res.AccessoriePrice.Add(product.Accessorie.Price)
					// 配件件数
					res.AccessorieQuantity += product.Accessorie.Quantity
				}
			}
		}

		// 订金抵扣
		res.DepositDeduction = res.DepositDeduction.Add(order.PriceDeposit)
	}

	l.resp.Itemized = res
}

// 获取支付项
func (l *StatisticSalesDetailDailyLogic) getPayment() {
	res := map[string]types.StatisticSalesDetailDailyPayment{}

	for method, name := range enums.OrderPaymentMethodMap {
		for _, order := range l.OrderSales {
			for _, pay := range order.Payments {
				if method == pay.PaymentMethod {
					if _, ok := res[name]; !ok {
						res[name] = types.StatisticSalesDetailDailyPayment{}
					}
					payment := res[name]
					payment.Income = payment.Income.Add(pay.Amount)
					res[name] = payment
				}
			}
		}

		for _, order := range l.OrderDeposit {
			for _, pay := range order.Payments {
				if method == pay.PaymentMethod {
					if _, ok := res[name]; !ok {
						res[name] = types.StatisticSalesDetailDailyPayment{}
					}
					payment := res[name]
					payment.Income = payment.Income.Add(pay.Amount)
					res[name] = payment
				}
			}
		}

		for _, order := range l.OrderOther {
			switch order.Type {
			case enums.FinanceTypeIncome:
				{
					for _, pay := range order.Payments {
						if method == pay.PaymentMethod {
							if _, ok := res[name]; !ok {
								res[name] = types.StatisticSalesDetailDailyPayment{}
							}
							payment := res[name]
							payment.Income = payment.Income.Add(pay.Amount)
							res[name] = payment
						}
					}
				}
			case enums.FinanceTypeExpense:
				{
					for _, pay := range order.Payments {
						if method == pay.PaymentMethod {
							if _, ok := res[name]; !ok {
								res[name] = types.StatisticSalesDetailDailyPayment{}
							}
							payment := res[name]
							payment.Expense = payment.Expense.Add(pay.Amount)
							res[name] = payment
						}
					}
				}
			}
		}
	}

	total := types.StatisticSalesDetailDailyPayment{
		Name: "汇总",
	}
	for n, i := range res {
		i.Name = n
		i.Received = i.Income.Sub(i.Expense)
		l.resp.Payment = append(l.resp.Payment, i)

		total.Income = total.Income.Add(i.Income)
		total.Expense = total.Expense.Add(i.Expense)
		total.Received = total.Received.Add(i.Received)
	}

	l.resp.PaymentTotal = total
}

// 获取成品销售
func (l *StatisticSalesDetailDailyLogic) getFinishedSales() {
	res := map[string]map[string]types.StatisticSalesDetailDailyFinishedSalesCraftsCategory{}

	for _, order := range l.OrderSales {
		for _, product := range order.Products {
			if product.Type != enums.ProductTypeFinished {
				continue
			}
			// 大类
			blockName := product.Finished.Product.Class.String()
			if _, ok := res[blockName]; !ok {
				res[blockName] = map[string]types.StatisticSalesDetailDailyFinishedSalesCraftsCategory{}
			}
			// 块
			block := res[blockName]
			// 行名
			rowName := product.Finished.Product.Category.String() + product.Finished.Product.Craft.String()
			if rowName == "" {
				rowName = "其他"
			}
			rowTotalName := product.Finished.Product.Craft.String() + "合计"

			if _, ok := block[rowName]; !ok {
				block[rowName] = types.StatisticSalesDetailDailyFinishedSalesCraftsCategory{}
			}
			if _, ok := block[rowTotalName]; !ok {
				block[rowTotalName] = types.StatisticSalesDetailDailyFinishedSalesCraftsCategory{}
			}

			row := block[rowName]
			rowTotal := block[rowTotalName]

			// 应收
			row.Receivable = row.Receivable.Add(product.Finished.Price)
			rowTotal.Receivable = rowTotal.Receivable.Add(product.Finished.Price)
			// 标签价
			row.Price = row.Price.Add(product.Finished.Product.LabelPrice)
			rowTotal.Price = rowTotal.Price.Add(product.Finished.Product.LabelPrice)
			// 金重
			row.WeightMetal = row.WeightMetal.Add(product.Finished.Product.WeightMetal)
			rowTotal.WeightMetal = rowTotal.WeightMetal.Add(product.Finished.Product.WeightMetal)
			// 工费
			row.LaborFee = row.LaborFee.Add(product.Finished.Product.LaborFee)
			rowTotal.LaborFee = rowTotal.LaborFee.Add(product.Finished.Product.LaborFee)
			// 件数
			row.Quantity++
			rowTotal.Quantity++

			block[rowName] = row
			block[rowTotalName] = rowTotal

			res[blockName] = block
		}
	}

	blockName := "汇总统计"
	if _, ok := res[blockName]; !ok {
		res[blockName] = map[string]types.StatisticSalesDetailDailyFinishedSalesCraftsCategory{}
	}
	totalRow := res[blockName]
	rowName := "合计"
	if _, ok := totalRow[rowName]; !ok {
		totalRow[rowName] = types.StatisticSalesDetailDailyFinishedSalesCraftsCategory{}
	}
	totalRowTotal := totalRow[rowName]
	for ResblockName, block := range res {
		if ResblockName == blockName {
			continue
		}
		for name, row := range block {
			if strings.Contains(name, rowName) {
				totalRowTotal.Receivable = totalRowTotal.Receivable.Add(row.Receivable)
				totalRowTotal.Price = totalRowTotal.Price.Add(row.Price)
				totalRowTotal.WeightMetal = totalRowTotal.WeightMetal.Add(row.WeightMetal)
				totalRowTotal.LaborFee = totalRowTotal.LaborFee.Add(row.LaborFee)
				totalRowTotal.Quantity += row.Quantity
			}
		}
	}
	totalRow[rowName] = totalRowTotal
	res[blockName] = totalRow

	l.resp.FinishedSales = res
}

// 获取配件销售
func (l *StatisticSalesDetailDailyLogic) getAccessorieSales() {
	res := map[string]types.StatisticSalesDetailDailyAccessorieSales{}

	for _, order := range l.OrderSales {
		for _, product := range order.Products {
			if product.Type != enums.ProductTypeAccessorie {
				continue
			}

			name := product.Accessorie.Product.Name
			toolName := "合计"
			if name == "" {
				name = "其他"
			}
			if _, ok := res[name]; !ok {
				res[name] = types.StatisticSalesDetailDailyAccessorieSales{}
			}
			if _, ok := res[toolName]; !ok {
				res[toolName] = types.StatisticSalesDetailDailyAccessorieSales{}
			}

			accessorie := res[name]
			accessorie.Receivable = accessorie.Receivable.Add(product.Accessorie.Price)
			accessorie.Price = accessorie.Price.Add(product.Accessorie.Product.Price)
			accessorie.Quantity += product.Accessorie.Quantity

			accessorieTotal := res[toolName]
			accessorieTotal.Receivable = accessorieTotal.Receivable.Add(product.Accessorie.Price)
			accessorieTotal.Price = accessorieTotal.Price.Add(product.Accessorie.Product.Price)
			accessorieTotal.Quantity += product.Accessorie.Quantity
			res[toolName] = accessorieTotal

			res[name] = accessorie

		}
	}

	l.resp.AccessorieSales = res
}
