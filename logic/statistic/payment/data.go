package payment

import (
	"errors"
	"jdy/enums"
	"jdy/model"

	"github.com/shopspring/decimal"
)

type datawLogic struct {
	*StatisticPaymentLogic
	req      *DataReq
	Payments []model.OrderPayment
}

type DataRes struct {
	Overview    map[string]any   `json:"overview"`
	List        map[string]any   `json:"list"`
	Details     []map[string]any `json:"details"`
	DetailTitle []string         `json:"detail_title"`
}

func (l *StatisticPaymentLogic) Data(req *DataReq) (any, error) {
	logic := datawLogic{
		StatisticPaymentLogic: l,
		req:                   req,
	}

	if err := logic.get_payments(); err != nil {
		return nil, err
	}

	res := DataRes{
		Overview: logic.get_overview(),
		List:     logic.get_list(),
		DetailTitle: []string{
			"来源",
			"方式",
			"金额",
			"时间",
			"单号",
		},
		Details: logic.get_details(),
	}

	return res, nil
}

func (l *datawLogic) get_payments() error {
	db := model.DB.Model(&model.OrderPayment{})
	db = db.Where(&model.OrderPayment{
		StoreId: l.req.StoreId,
		Status:  true,
	})
	db = db.Scopes(model.DurationCondition(l.req.Duration, "created_at", l.req.StartTime, l.req.EndTime))

	db = db.Order("created_at desc")
	if err := db.Find(&l.Payments).Error; err != nil {
		return errors.New("获取数据失败")
	}

	return nil
}

func (l *datawLogic) get_overview() map[string]any {
	data := make(map[string]any)

	data["收入金额"] = decimal.Decimal{}
	data["收入笔数"] = decimal.Decimal{}
	data["支出金额"] = decimal.Decimal{}
	data["支出笔数"] = decimal.Decimal{}

	for _, payment := range l.Payments {
		switch payment.Type {
		case enums.FinanceTypeIncome: // 收入
			{
				price, ok := data["收入金额"].(decimal.Decimal)
				if !ok {
					price = decimal.Decimal{}
				}
				price = price.Add(payment.Amount)
				data["收入金额"] = price

				num, ok := data["收入笔数"].(decimal.Decimal)
				if !ok {
					num = decimal.Decimal{}
				}
				num = num.Add(decimal.NewFromInt(1))
				data["收入笔数"] = num
			}
		case enums.FinanceTypeExpense: // 支出
			{
				price, ok := data["支出金额"].(decimal.Decimal)
				if !ok {
					price = decimal.Decimal{}
				}
				price = price.Add(payment.Amount.Neg())
				data["支出金额"] = price

				num, ok := data["支出笔数"].(decimal.Decimal)
				if !ok {
					num = decimal.Decimal{}
				}
				num = num.Add(decimal.NewFromInt(1))
				data["支出笔数"] = num
			}
		}
	}

	return data
}

func (l *datawLogic) get_list() map[string]any {
	data := make(map[string]any)

	totalIncome := decimal.Zero
	totalOutcome := decimal.Zero

	if len(l.Payments) == 0 {
		data["合计"] = map[string]any{
			"收入": decimal.Zero,
			"支出": decimal.Zero,
			"结余": decimal.Zero,
		}
		return data
	}

	for _, payment := range l.Payments {
		k := payment.PaymentMethod.String()

		row, _ := data[k].(map[string]any)
		if row == nil {
			row = map[string]any{
				"收入": decimal.Zero,
				"支出": decimal.Zero,
				"结余": decimal.Zero,
			}
		}

		switch payment.Type {
		case enums.FinanceTypeIncome:
			row["收入"] = row["收入"].(decimal.Decimal).Add(payment.Amount)
			totalIncome = totalIncome.Add(payment.Amount)
		case enums.FinanceTypeExpense:
			row["支出"] = row["支出"].(decimal.Decimal).Add(payment.Amount.Neg())
			totalOutcome = totalOutcome.Add(payment.Amount.Neg())
		}

		// Update this row’s balance
		row["结余"] = row["收入"].(decimal.Decimal).Add(row["支出"].(decimal.Decimal))
		data[k] = row
	}

	data["合计"] = map[string]any{
		"收入": totalIncome,
		"支出": totalOutcome,
		"结余": totalIncome.Add(totalOutcome),
	}

	return data
}

func (l *datawLogic) get_details() []map[string]any {
	res := []map[string]any{}

	for _, payment := range l.Payments {
		data := map[string]any{
			"来源": payment.Source.String(),
			"方式": payment.PaymentMethod.String(),
			"时间": payment.CreatedAt.Format("2006-01-02 15:04:05"),
			"单号": payment.OrderId,
		}
		switch payment.Type {
		case enums.FinanceTypeIncome:
			data["金额"] = payment.Amount
		case enums.FinanceTypeExpense:
			data["金额"] = payment.Amount.Neg()
		}

		res = append(res, data)
	}

	return res
}
