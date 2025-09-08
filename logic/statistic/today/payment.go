package today

import (
	"errors"
	"jdy/enums"
	"jdy/model"

	"github.com/shopspring/decimal"
)

type PaymentReq struct {
	StoreId string `json:"store_id"` // 门店ID
}

type PaymentLogic struct {
	Req *PaymentReq
	Res map[string]any

	Payments []model.OrderPayment
}

func (l *ToDayLogic) Payment(req *PaymentReq) (map[string]any, error) {
	logic := &PaymentLogic{
		Req: req,
		Res: make(map[string]any),
	}

	if err := logic.get_payments(); err != nil {
		return nil, err
	}

	if err := logic.get_data(); err != nil {
		return nil, err
	}

	return logic.Res, nil
}

func (l *PaymentLogic) get_payments() error {
	db := model.DB.Model(&model.OrderPayment{})
	db = db.Where(&model.OrderPayment{
		StoreId: l.Req.StoreId,
		Status:  true,
	})
	db = db.Scopes(model.DurationCondition(enums.DurationYesterday))

	if err := db.Find(&l.Payments).Error; err != nil {
		return errors.New("获取数据失败")
	}

	return nil
}

func (l *PaymentLogic) get_data() error {
	if len(l.Payments) == 0 {
		l.Res["收入金额"] = decimal.Decimal{}
		l.Res["收入笔数"] = decimal.Decimal{}
		l.Res["收入订单数"] = decimal.Decimal{}
		l.Res["支出金额"] = decimal.Decimal{}
		l.Res["支出笔数"] = decimal.Decimal{}
		l.Res["支出订单数"] = decimal.Decimal{}

		return nil
	}

	// 收入订单号
	incomeNum := make(map[string]bool)
	// 支出订单号
	expenseNum := make(map[string]bool)

	for _, payment := range l.Payments {
		switch payment.Type {
		case enums.FinanceTypeIncome: // 收入
			{
				price, ok := l.Res["收入金额"].(decimal.Decimal)
				if !ok {
					price = decimal.Decimal{}
				}
				price = price.Add(payment.Amount)
				l.Res["收入金额"] = price

				num, ok := l.Res["收入笔数"].(decimal.Decimal)
				if !ok {
					num = decimal.Decimal{}
				}
				num = num.Add(decimal.NewFromInt(1))
				l.Res["收入笔数"] = num

				incomeNum[payment.OrderId] = true
			}
		case enums.FinanceTypeExpense: // 支出
			{
				price, ok := l.Res["支出金额"].(decimal.Decimal)
				if !ok {
					price = decimal.Decimal{}
				}
				price = price.Add(payment.Amount)
				l.Res["支出金额"] = price

				num, ok := l.Res["支出笔数"].(decimal.Decimal)
				if !ok {
					num = decimal.Decimal{}
				}
				num = num.Add(decimal.NewFromInt(1))
				l.Res["支出笔数"] = num

				expenseNum[payment.OrderId] = true
			}
		}
	}

	l.Res["收入订单数"] = decimal.NewFromInt(int64(len(incomeNum)))
	l.Res["支出订单数"] = decimal.NewFromInt(int64(len(expenseNum)))

	return nil
}
