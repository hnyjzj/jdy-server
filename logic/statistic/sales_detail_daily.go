package statistic

import (
	"errors"
	"jdy/enums"
	"jdy/model"
	"jdy/types"
	"jdy/utils"
	"log"
	"strings"

	"github.com/shopspring/decimal"
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

		if err := logic.getOrderRefund(); err != nil {
			return err
		}

		if err := logic.getOrderPayment(); err != nil {
			return err
		}

		// 开始统计
		logic.getSummary()
		logic.getItemized()
		logic.getPayment()
		logic.getFinishedSales()
		logic.getOldSales()
		logic.getAccessorieSales()
		logic.getFinishedSalesRefund()
		logic.getOldSalesRefund()
		logic.getAccessorieSalesRefund()

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
	OrderRefund  []model.OrderRefund
	OrderPayment []model.OrderPayment
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
		enums.OrderDepositStatusBooking,
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

// 获取退货单
func (l *StatisticSalesDetailDailyLogic) getOrderRefund() error {
	db := l.tx.Model(&model.OrderRefund{})

	if l.req.StoreId != "" {
		db = db.Where("store_id = ?", l.req.StoreId)
	}
	if l.req.StartTime != nil {
		db = db.Where("created_at >= ?", l.req.StartTime)
	}
	if l.req.EndTime != nil {
		db = db.Where("created_at <= ?", l.req.EndTime)
	}

	if err := db.Find(&l.OrderRefund).Error; err != nil {
		return err
	}

	for i, refund := range l.OrderRefund {
		var order model.OrderSalesProduct
		pdb := l.tx.Model(&model.OrderSalesProduct{})
		switch refund.Type {
		case enums.ProductTypeFinished:
			{
				pdb = pdb.Where(&model.OrderSalesProduct{
					OrderId: refund.OrderId,
					Type:    enums.ProductTypeFinished,
					Code:    refund.Code,
				})
			}
		case enums.ProductTypeOld:
			{
				pdb = pdb.Where(&model.OrderSalesProduct{
					OrderId: refund.OrderId,
					Type:    enums.ProductTypeOld,
					Code:    refund.Code,
				})
			}
		case enums.ProductTypeAccessorie:
			{
				pdb = pdb.Where(&model.OrderSalesProduct{
					OrderId: refund.OrderId,
					Type:    enums.ProductTypeAccessorie,
					Name:    refund.Name,
				})
			}
		}

		pdb = order.Preloads(pdb)
		if err := pdb.First(&order).Error; err != nil {
			return err
		}

		switch refund.Type {
		case enums.ProductTypeFinished:
			{
				l.OrderRefund[i].Product = order.Finished.Product
			}
		case enums.ProductTypeOld:
			{
				l.OrderRefund[i].Product = order.Old.Product
			}
		case enums.ProductTypeAccessorie:
			{
				l.OrderRefund[i].Product = order.Accessorie.Product
			}
		}
	}

	return nil
}

// 获取支付数据
func (l *StatisticSalesDetailDailyLogic) getOrderPayment() error {
	db := l.tx.Model(&model.OrderPayment{})

	if l.req.StoreId != "" {
		db = db.Where("store_id = ?", l.req.StoreId)
	}
	if l.req.StartTime != nil {
		db = db.Where("created_at >= ?", l.req.StartTime)
	}
	if l.req.EndTime != nil {
		db = db.Where("created_at <= ?", l.req.EndTime)
	}

	db = db.Where(&model.OrderPayment{
		Status: true,
	})

	if err := db.Find(&l.OrderPayment).Error; err != nil {
		return err
	}

	return nil
}

// 获取汇总项
func (l *StatisticSalesDetailDailyLogic) getSummary() {
	// 汇总
	res := types.StatisticSalesDetailDailySummary{}

	for _, sales := range l.OrderSales {
		// 销售应收
		res.SalesReceivable = res.SalesReceivable.Add(sales.Price)
		// 销售实收
		res.SalesReceived = res.SalesReceived.Add(sales.PricePay)
		// 排除退款
		for _, refund := range l.OrderRefund {
			if refund.OrderId != sales.Id {
				continue
			}
			// 销售退款
			switch refund.Type {
			case enums.ProductTypeOld:
				{
					res.SalesReceived = res.SalesReceived.Add(refund.Price)
				}
			default:
				{
					res.SalesReceived = res.SalesReceived.Sub(refund.Price)
				}
			}
		}
	}

	for _, refund := range l.OrderRefund {
		// 销售退款
		switch refund.Type {
		case enums.ProductTypeOld:
			{
				res.SalesRefund = res.SalesRefund.Add(refund.Price)
			}
		default:
			{
				res.SalesRefund = res.SalesRefund.Sub(refund.Price)
			}
		}
	}

	for _, deposit := range l.OrderDeposit {
		// 订金收入
		res.DepositIncome = res.DepositIncome.Add(deposit.PricePay)
		// 订金退款
		if deposit.Status == enums.OrderDepositStatusRefund {
			for _, refund := range deposit.Products {
				if refund.Status != enums.OrderDepositStatusReturn && refund.Status != enums.OrderDepositStatusRefund {
					continue
				}
				res.DepositRefund = res.DepositRefund.Add(refund.Price)
			}
		}
	}

	for _, other := range l.OrderOther {
		// 其他收入
		if other.Type == enums.FinanceTypeIncome {
			res.OtherIncome = res.OtherIncome.Add(other.Amount)
		}
		// 其他支出
		if other.Type == enums.FinanceTypeExpense {
			res.OtherExpense = res.OtherExpense.Add(other.Amount)
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

					// 排除退款
					for _, refund := range l.OrderRefund {
						if refund.Type != enums.ProductTypeFinished {
							continue
						}
						if refund.OrderId != order.Id {
							continue
						}
						if refund.Code != product.Code {
							continue
						}

						// 成品应收
						res.FinishedReceivable = res.FinishedReceivable.Sub(product.Finished.Price)
						// 成品件数
						res.FinishedQuantity--
					}
				}
			case enums.ProductTypeOld:
				{
					// 旧料抵值
					res.OldDeduction = res.OldDeduction.Sub(product.Old.RecyclePrice)
					// 旧料件数
					res.OldQuantity++
					// 旧料金重
					res.OldWeightMetal = res.OldWeightMetal.Add(product.Old.WeightMetal)

					// 排除退款
					for _, refund := range l.OrderRefund {
						if refund.Type != enums.ProductTypeOld {
							continue
						}
						if refund.OrderId != order.Id {
							continue
						}
						if refund.Code != product.Code {
							continue
						}

						// 旧料抵值
						res.OldDeduction = res.OldDeduction.Add(product.Old.RecyclePrice)
						// 旧料件数
						res.OldQuantity--
						// 旧料重量
						res.OldWeightMetal = res.OldWeightMetal.Sub(product.Old.WeightMetal)
					}

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

					// 排除退款
					for _, refund := range l.OrderRefund {
						if refund.Type != enums.ProductTypeAccessorie {
							continue
						}
						if refund.OrderId != order.Id {
							continue
						}
						if refund.Name != product.Name {
							continue
						}

						// 配件金额
						res.AccessoriePrice = res.AccessoriePrice.Sub(refund.Price)
						// 配件件数
						res.AccessorieQuantity -= refund.Quantity
					}
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
		for _, order := range l.OrderPayment {
			if method == order.PaymentMethod {
				if _, ok := res[name]; !ok {
					res[name] = types.StatisticSalesDetailDailyPayment{}
				}
				payment := res[name]
				switch order.Type {
				case enums.FinanceTypeIncome:
					{
						payment.Income = payment.Income.Add(order.Amount)
					}
				case enums.FinanceTypeExpense:
					{
						payment.Expense = payment.Expense.Sub(order.Amount)
					}
				}
				res[name] = payment
			}
		}
	}

	total := types.StatisticSalesDetailDailyPayment{
		Name: "汇总",
	}
	for n, i := range res {
		i.Name = n
		i.Received = i.Income.Add(i.Expense)
		l.resp.Payment = append(l.resp.Payment, i)

		total.Income = total.Income.Add(i.Income)
		total.Expense = total.Expense.Add(i.Expense)
		total.Received = total.Received.Add(i.Received)
	}

	l.resp.PaymentTotal = total
}

// 获取成品销售
func (l *StatisticSalesDetailDailyLogic) getFinishedSales() {
	res := map[string]map[string]types.StatisticSalesDetailDailyFinishedSales{}

	for _, order := range l.OrderSales {
		for _, product := range order.Products {
			if product.Type != enums.ProductTypeFinished {
				continue
			}
			// 大类
			blockName := product.Finished.Product.Class.String()
			if _, ok := res[blockName]; !ok {
				res[blockName] = map[string]types.StatisticSalesDetailDailyFinishedSales{}
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
				block[rowName] = types.StatisticSalesDetailDailyFinishedSales{}
			}
			if _, ok := block[rowTotalName]; !ok {
				block[rowTotalName] = types.StatisticSalesDetailDailyFinishedSales{}
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
		res[blockName] = map[string]types.StatisticSalesDetailDailyFinishedSales{}
	}
	totalRow := res[blockName]
	rowName := "合计"
	if _, ok := totalRow[rowName]; !ok {
		totalRow[rowName] = types.StatisticSalesDetailDailyFinishedSales{}
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

// 获取旧料回收
func (l *StatisticSalesDetailDailyLogic) getOldSales() {
	res := map[string][]types.StatisticSalesDetailDailyOldSales{}

	for _, order := range l.OrderSales {
		for _, product := range order.Products {
			if product.Type != enums.ProductTypeOld {
				continue
			}
			// 大类
			blockName := product.Old.Product.RecycleMethod.String()
			if _, ok := res[blockName]; !ok {
				res[blockName] = []types.StatisticSalesDetailDailyOldSales{}
			}
			// 块
			block := res[blockName]
			if block == nil {
				block = make([]types.StatisticSalesDetailDailyOldSales, 0)
			}
			// 行名
			name := []string{
				product.Old.Product.Material.String(),
				product.Old.Product.Quality.String(),
				product.Old.Product.Gem.String(),
			}
			// 行名
			rowName := strings.Join(name, "")
			if rowName == "" {
				rowName = "其他"
			}
			// 统计行名
			totalName := strings.Join(append([]string{
				product.Old.Product.RecycleType.String(),
			}, name...), "-")

			// 行
			row := types.StatisticSalesDetailDailyOldSales{
				Name: rowName,
			}
			// 查找统计行
			total, index, err := utils.ArrayFind(block, func(item types.StatisticSalesDetailDailyOldSales) bool {
				return item.Name == totalName
			})
			// 如果没有找到统计行
			if err != nil {
				total = types.StatisticSalesDetailDailyOldSales{
					Name: totalName,
				}
			} else {
				// 如果找到了统计行,则删除统计行
				block = utils.ArrayDeleteOfIndex(block, index)
			}

			// 抵值
			row.Deduction = product.Old.RecyclePrice
			total.Deduction = total.Deduction.Add(product.Old.RecyclePrice)
			// 金重
			row.WeightMetal = product.Old.Product.WeightMetal
			total.WeightMetal = total.WeightMetal.Add(product.Old.Product.WeightMetal)
			// 宝石重
			row.WeightGem = product.Old.Product.WeightGem
			total.WeightGem = total.WeightGem.Add(product.Old.Product.WeightGem)
			// 件数
			row.Quantity++
			total.Quantity++
			// 标签价
			row.LabelPrice = product.Old.Product.LabelPrice
			total.LabelPrice = total.LabelPrice.Add(product.Old.Product.LabelPrice)
			// 工费
			row.LaborFee = product.Old.Product.RecyclePriceLabor
			total.LaborFee = total.LaborFee.Add(product.Old.Product.RecyclePriceLabor)
			// 转成品抵值
			row.ToFinishedDeduction = product.Old.Product.Code
			total.ToFinishedDeduction = "0"
			// 转成品金重
			toFinishedWeightMetal := decimal.NewFromInt(0)
			if total.ToFinishedWeightMetal == nil {
				total.ToFinishedWeightMetal = &toFinishedWeightMetal
			}
			if product.Old.Product.DeletedAt.Valid {
				toFinishedWeightMetal = total.ToFinishedWeightMetal.Add(product.Old.Product.WeightMetal)
			}
			total.ToFinishedWeightMetal = &toFinishedWeightMetal
			// 转成品件数
			var toFinishedQuantity int64 = 0
			if total.ToFinishedQuantity == nil {
				total.ToFinishedQuantity = &toFinishedQuantity
			}
			if product.Old.Product.DeletedAt.Valid {
				toFinishedQuantity++
			}
			total.ToFinishedQuantity = &toFinishedQuantity
			// 剩余金重
			surplusWeight := decimal.NewFromInt(0)
			if total.SurplusWeight == nil {
				total.SurplusWeight = &surplusWeight
			}
			surplusWeight = total.SurplusWeight.Add(row.WeightMetal)
			total.SurplusWeight = &surplusWeight

			if index == -1 {
				block = append(block, total)
			} else {
				// 将 block 放回原处
				block = append(block[:index], append([]types.StatisticSalesDetailDailyOldSales{total}, block[index:]...)...)
			}
			block = append(block, row)

			res[blockName] = block
		}
	}

	blockName := "合计"
	if _, ok := res[blockName]; !ok {
		res[blockName] = []types.StatisticSalesDetailDailyOldSales{}
	}
	totalRow := res[blockName]
	total := types.StatisticSalesDetailDailyOldSales{}
	for ResblockName, block := range res {
		if ResblockName == blockName {
			continue
		}
		for _, row := range block {
			if strings.Contains(row.Name, "-") {
				total.Deduction = total.Deduction.Add(row.Deduction)
				total.WeightMetal = total.WeightMetal.Add(row.WeightMetal)
				total.WeightGem = total.WeightGem.Add(row.WeightGem)
				total.Quantity += row.Quantity
				total.LabelPrice = total.LabelPrice.Add(row.LabelPrice)
				total.LaborFee = total.LaborFee.Add(row.LaborFee)

				ToFinishedDeduction, err := decimal.NewFromString(total.ToFinishedDeduction)
				if err != nil {
					ToFinishedDeduction = decimal.NewFromInt(0)
				}
				toFinishedDeduction, err := decimal.NewFromString(row.ToFinishedDeduction)
				if err != nil {
					toFinishedDeduction = decimal.NewFromInt(0)
				}
				total.ToFinishedDeduction = ToFinishedDeduction.Add(toFinishedDeduction).String()

				toFinishedWeightMetal := decimal.NewFromInt(0)
				if total.ToFinishedWeightMetal == nil {
					total.ToFinishedWeightMetal = &toFinishedWeightMetal
				}
				if row.ToFinishedWeightMetal != nil {
					toFinishedWeightMetal = total.ToFinishedWeightMetal.Add(*row.ToFinishedWeightMetal)
					total.ToFinishedWeightMetal = &toFinishedWeightMetal
				}

				var toFinishedQuantity int64 = 0
				if total.ToFinishedQuantity == nil {
					total.ToFinishedQuantity = &toFinishedQuantity
				}
				if row.ToFinishedQuantity != nil {
					toFinishedQuantity = *total.ToFinishedQuantity + *row.ToFinishedQuantity
					total.ToFinishedQuantity = &toFinishedQuantity
				}

				surplusWeight := decimal.NewFromInt(0)
				if total.SurplusWeight == nil {
					total.SurplusWeight = &surplusWeight
				}
				if row.SurplusWeight != nil {
					surplusWeight = total.SurplusWeight.Add(*row.SurplusWeight)
					total.SurplusWeight = &surplusWeight
				}
			}
		}
	}
	totalRow = append(totalRow, total)
	res[blockName] = totalRow

	l.resp.OldSales = res
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
			if name == "" {
				name = "其他"
			}
			if _, ok := res[name]; !ok {
				res[name] = types.StatisticSalesDetailDailyAccessorieSales{}
			}

			accessorie := res[name]
			accessorie.Received = accessorie.Received.Add(product.Accessorie.Price)
			accessorie.Receivable = accessorie.Receivable.Add(product.Accessorie.PriceOriginal)
			accessorie.Price = accessorie.Price.Add(product.Accessorie.Product.Price)
			accessorie.Quantity += product.Accessorie.Quantity

			res[name] = accessorie

		}
	}

	toolName := "合计"
	if _, ok := res[toolName]; !ok {
		res[toolName] = types.StatisticSalesDetailDailyAccessorieSales{}
	}
	total := res[toolName]
	for name, accessorie := range res {
		if name == toolName {
			continue
		}
		total.Received = total.Received.Add(accessorie.Received)
		total.Receivable = total.Receivable.Add(accessorie.Receivable)
		total.Price = total.Price.Add(accessorie.Price)
		total.Quantity += accessorie.Quantity
	}
	res[toolName] = total

	l.resp.AccessorieSales = res
}

// 获取成品退货
func (l *StatisticSalesDetailDailyLogic) getFinishedSalesRefund() {
	res := map[string]map[string]types.StatisticSalesDetailDailyFinishedSalesRefund{}

	for _, order := range l.OrderRefund {
		if order.Type != enums.ProductTypeFinished {
			continue
		}

		product, ok := order.Product.(model.ProductFinished)
		if !ok {
			log.Printf("order: %+v \n product: %+v", order, product)
			continue
		}

		// 大类
		blockName := product.Class.String()
		if _, ok := res[blockName]; !ok {
			res[blockName] = map[string]types.StatisticSalesDetailDailyFinishedSalesRefund{}
		}
		// 块
		block := res[blockName]
		// 行名
		rowName := product.Category.String() + product.Craft.String()
		if rowName == "" {
			rowName = "其他"
		}
		rowTotalName := product.Craft.String() + "合计"

		if _, ok := block[rowName]; !ok {
			block[rowName] = types.StatisticSalesDetailDailyFinishedSalesRefund{}
		}
		if _, ok := block[rowTotalName]; !ok {
			block[rowTotalName] = types.StatisticSalesDetailDailyFinishedSalesRefund{}
		}

		row := block[rowName]
		rowTotal := block[rowTotalName]

		// 退款
		row.Refunded = row.Refunded.Add(order.Price.Neg())
		rowTotal.Refunded = rowTotal.Refunded.Add(order.Price.Neg())
		// 标签价
		row.Price = row.Price.Add(product.LabelPrice)
		rowTotal.Price = rowTotal.Price.Add(product.LabelPrice)
		// 金重
		row.WeightMetal = row.WeightMetal.Add(product.WeightMetal)
		rowTotal.WeightMetal = rowTotal.WeightMetal.Add(product.WeightMetal)
		// 工费
		row.LaborFee = row.LaborFee.Add(product.LaborFee)
		rowTotal.LaborFee = rowTotal.LaborFee.Add(product.LaborFee)
		// 件数
		row.Quantity++
		rowTotal.Quantity++

		block[rowName] = row
		block[rowTotalName] = rowTotal

		res[blockName] = block

	}
	blockName := "汇总统计"
	if _, ok := res[blockName]; !ok {
		res[blockName] = map[string]types.StatisticSalesDetailDailyFinishedSalesRefund{}
	}
	totalRow := res[blockName]
	rowName := "合计"
	if _, ok := totalRow[rowName]; !ok {
		totalRow[rowName] = types.StatisticSalesDetailDailyFinishedSalesRefund{}
	}
	totalRowTotal := totalRow[rowName]
	for ResblockName, block := range res {
		if ResblockName == blockName {
			continue
		}
		for name, row := range block {
			if strings.Contains(name, rowName) {
				totalRowTotal.Refunded = totalRowTotal.Refunded.Add(row.Refunded)
				totalRowTotal.Price = totalRowTotal.Price.Add(row.Price)
				totalRowTotal.WeightMetal = totalRowTotal.WeightMetal.Add(row.WeightMetal)
				totalRowTotal.LaborFee = totalRowTotal.LaborFee.Add(row.LaborFee)
				totalRowTotal.Quantity += row.Quantity
			}
		}
	}
	totalRow[rowName] = totalRowTotal
	res[blockName] = totalRow

	l.resp.FinishedSalesRefund = res
}

// 获取旧料退货
func (l *StatisticSalesDetailDailyLogic) getOldSalesRefund() {
	res := map[string][]types.StatisticSalesDetailDailyOldSalesRefund{}

	for _, refund := range l.OrderRefund {
		if refund.Type != enums.ProductTypeOld {
			continue
		}
		product, ok := refund.Product.(model.ProductOld)
		if !ok {
			log.Printf("order: %+v \n product: %+v", refund, product)
			continue
		}
		// 大类
		blockName := product.RecycleMethod.String()
		if _, ok := res[blockName]; !ok {
			res[blockName] = []types.StatisticSalesDetailDailyOldSalesRefund{}
		}
		// 块
		block := res[blockName]
		if block == nil {
			block = make([]types.StatisticSalesDetailDailyOldSalesRefund, 0)
		}
		// 行名
		name := []string{
			product.Material.String(),
			product.Quality.String(),
			product.Gem.String(),
		}
		// 行名
		rowName := strings.Join(name, "")
		if rowName == "" {
			rowName = "其他"
		}
		// 统计行名
		totalName := strings.Join(append([]string{
			product.RecycleType.String(),
		}, name...), "-")

		// 行
		row := types.StatisticSalesDetailDailyOldSalesRefund{
			Name: rowName,
		}
		// 查找统计行
		total, index, err := utils.ArrayFind(block, func(item types.StatisticSalesDetailDailyOldSalesRefund) bool {
			return item.Name == totalName
		})
		// 如果没有找到统计行
		if err != nil {
			total = types.StatisticSalesDetailDailyOldSalesRefund{
				Name: totalName,
			}
		} else {
			// 如果找到了统计行,则删除统计行
			block = utils.ArrayDeleteOfIndex(block, index)
		}

		// 退款
		row.Refunded = refund.Price.Neg()
		total.Refunded = total.Refunded.Add(refund.Price.Neg())
		// 金重
		row.WeightMetal = product.WeightMetal
		total.WeightMetal = total.WeightMetal.Add(product.WeightMetal)
		// 宝石重
		row.WeightGem = product.WeightGem
		total.WeightGem = total.WeightGem.Add(product.WeightGem)
		// 件数
		row.Quantity++
		total.Quantity++
		// 标签价
		row.LabelPrice = product.LabelPrice
		total.LabelPrice = total.LabelPrice.Add(product.LabelPrice)
		// 工费
		row.LaborFee = product.RecyclePriceLabor
		total.LaborFee = total.LaborFee.Add(product.RecyclePriceLabor)
		// 条码
		row.Code = product.Code

		if index == -1 {
			block = append(block, total)
		} else {
			// 将 block 放回原处
			block = append(block[:index], append([]types.StatisticSalesDetailDailyOldSalesRefund{total}, block[index:]...)...)
		}
		block = append(block, row)

		res[blockName] = block

	}

	blockName := "合计"
	if _, ok := res[blockName]; !ok {
		res[blockName] = []types.StatisticSalesDetailDailyOldSalesRefund{}
	}
	totalRow := res[blockName]
	total := types.StatisticSalesDetailDailyOldSalesRefund{}
	for ResblockName, block := range res {
		if ResblockName == blockName {
			continue
		}
		for _, row := range block {
			if strings.Contains(row.Name, "-") {
				total.Refunded = total.Refunded.Add(row.Refunded)
				total.WeightMetal = total.WeightMetal.Add(row.WeightMetal)
				total.WeightGem = total.WeightGem.Add(row.WeightGem)
				total.Quantity += row.Quantity
				total.LabelPrice = total.LabelPrice.Add(row.LabelPrice)
				total.LaborFee = total.LaborFee.Add(row.LaborFee)
			}
		}
	}
	totalRow = append(totalRow, total)
	res[blockName] = totalRow

	l.resp.OldSalesRefund = res
}

// 获取配件退货
func (l *StatisticSalesDetailDailyLogic) getAccessorieSalesRefund() {

	res := map[string]types.StatisticSalesDetailDailyAccessorieSalesRefund{}

	for _, refund := range l.OrderRefund {
		if refund.Type != enums.ProductTypeAccessorie {
			continue
		}

		name := refund.Name
		if name == "" {
			name = "其他"
		}
		if _, ok := res[name]; !ok {
			res[name] = types.StatisticSalesDetailDailyAccessorieSalesRefund{}
		}
		accessorie := res[name]

		accessorie.Refunded = accessorie.Refunded.Add(refund.Price.Neg())
		accessorie.Quantity += refund.Quantity

		res[name] = accessorie

	}

	toolName := "合计"
	if _, ok := res[toolName]; !ok {
		res[toolName] = types.StatisticSalesDetailDailyAccessorieSalesRefund{}
	}
	total := res[toolName]
	for name, accessorie := range res {
		if name == toolName {
			continue
		}
		total.Refunded = total.Refunded.Add(accessorie.Refunded)
		total.Quantity += accessorie.Quantity
	}
	res[toolName] = total

	l.resp.AccessorieSalesRefund = res
}
