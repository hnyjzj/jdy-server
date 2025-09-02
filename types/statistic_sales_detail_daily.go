package types

import (
	"errors"
	"time"

	"github.com/shopspring/decimal"
)

type StatisticSalesDetailDailyReq struct {
	StoreId    string `json:"store_id"`    // 店铺id
	SalesmanId string `json:"salesman_id"` // 业务员id

	StartTime *time.Time `json:"start_time" binding:"required"` // 开始时间
	EndTime   *time.Time `json:"end_time" binding:"required"`   // 结束时间
}

func (req *StatisticSalesDetailDailyReq) Validate() error {
	// 开始时间不能大于结束时间
	if req.StartTime.After(*req.EndTime) {
		return errors.New("开始时间不能大于结束时间")
	}

	return nil
}

type StatisticSalesDetailDailyResp struct {
	Summary      StatisticSalesDetailDailySummary   `json:"summary"`       // 汇总项
	Itemized     StatisticSalesDetailDailyItemized  `json:"itemized"`      // 分项
	Payment      []StatisticSalesDetailDailyPayment `json:"payment"`       // 支付项
	PaymentTotal StatisticSalesDetailDailyPayment   `json:"payment_total"` // 支付项汇总

	FinishedSales       map[string]map[string]StatisticSalesDetailDailyFinishedSales       `json:"finished_sales"`        // 成品销售
	OldSales            map[string][]StatisticSalesDetailDailyOldSales                     `json:"old_sales"`             // 旧料销售
	AccessorieSales     map[string]StatisticSalesDetailDailyAccessorieSales                `json:"accessorie_sales"`      // 配件销售
	FinishedSalesRefund map[string]map[string]StatisticSalesDetailDailyFinishedSalesRefund `json:"finished_sales_refund"` // 成品退货
}

type StatisticSalesDetailDailySummary struct {
	SalesReceivable decimal.Decimal `json:"sales_receivable"` // 销售应收
	SalesRefund     decimal.Decimal `json:"sales_refund"`     // 销售退款
	SalesReceived   decimal.Decimal `json:"sales_received"`   // 销售实收
	OtherIncome     decimal.Decimal `json:"other_income"`     // 其他收入
	OtherExpense    decimal.Decimal `json:"other_expense"`    // 其他支出
	DepositIncome   decimal.Decimal `json:"deposit_income"`   // 定金收入
	DepositRefund   decimal.Decimal `json:"deposit_refund"`   // 定金退款
}

type StatisticSalesDetailDailyItemized struct {
	FinishedReceivable       decimal.Decimal `json:"finished_receivable"`          // 成品应收
	FinishedQuantity         int64           `json:"finished_quantity"`            // 成品件数
	DepositDeduction         decimal.Decimal `json:"deposit_deduction"`            // 定金抵扣
	OldDeduction             decimal.Decimal `json:"old_deduction"`                // 旧料抵扣
	OldQuantity              int64           `json:"old_quantity"`                 // 旧料件数
	OldWeightMetal           decimal.Decimal `json:"old_weight_metal"`             // 旧料金重
	OldToFinishedQuantity    int64           `json:"old_to_finished_quantity"`     // 旧料转成品件数
	OldToFinishedDeduction   decimal.Decimal `json:"old_to_finished_deduction"`    // 旧料转成品抵扣
	OldToFinishedWeightMetal decimal.Decimal `json:"old_to_finished_weight_metal"` // 旧料转成品金重
	AccessoriePrice          decimal.Decimal `json:"accessorie_price"`             // 配件金额
	AccessorieQuantity       int64           `json:"accessorie_quantity"`          // 配件件数
}

type StatisticSalesDetailDailyPayment struct {
	Name     string          `json:"name"`     // 名称
	Income   decimal.Decimal `json:"income"`   // 收入
	Expense  decimal.Decimal `json:"expense"`  // 支出
	Received decimal.Decimal `json:"received"` // 实收
}

type StatisticSalesDetailDailyFinishedSales struct {
	Receivable  decimal.Decimal `json:"receivable"`   // 应收
	Price       decimal.Decimal `json:"price"`        // 标签价
	WeightMetal decimal.Decimal `json:"weight_metal"` // 金重
	LaborFee    decimal.Decimal `json:"labor_fee"`    // 工费
	Quantity    int64           `json:"quantity"`     // 件数
}

type StatisticSalesDetailDailyOldSales struct {
	Name                  string           `json:"name"`                               // 名称
	Deduction             decimal.Decimal  `json:"deduction"`                          // 抵值
	WeightMetal           decimal.Decimal  `json:"weight_metal"`                       // 金重
	WeightGem             decimal.Decimal  `json:"weight_gem"`                         // 主石重
	Quantity              int64            `json:"quantity"`                           // 件数
	LabelPrice            decimal.Decimal  `json:"label_price"`                        // 标签价
	LaborFee              decimal.Decimal  `json:"labor_fee"`                          // 工费
	ToFinishedDeduction   string           `json:"to_finished_deduction"`              // 转成品抵值
	ToFinishedWeightMetal *decimal.Decimal `json:"to_finished_weight_metal,omitempty"` // 转成品金重
	ToFinishedQuantity    *int64           `json:"to_finished_quantity,omitempty"`     // 转成品件数
	SurplusWeight         *decimal.Decimal `json:"surplus_weight,omitempty"`           // 剩余金重
}

type StatisticSalesDetailDailyAccessorieSales struct {
	Received   decimal.Decimal `json:"received"`   // 实收
	Receivable decimal.Decimal `json:"receivable"` // 应收
	Price      decimal.Decimal `json:"price"`      // 单价
	Quantity   int64           `json:"quantity"`   // 件数
}

type StatisticSalesDetailDailyFinishedSalesRefund struct {
	Refunded    decimal.Decimal `json:"refunded"`     // 退款
	Price       decimal.Decimal `json:"price"`        // 标签价
	WeightMetal decimal.Decimal `json:"weight_metal"` // 金重
	LaborFee    decimal.Decimal `json:"labor_fee"`    // 工费
	Quantity    int64           `json:"quantity"`     // 件数
}
