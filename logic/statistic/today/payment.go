package today

import (
	"errors"
	"jdy/enums"
	"jdy/model"

	"github.com/shopspring/decimal"
)

type PaymentReq struct {
	DataReq
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
	db = db.Scopes(model.DurationCondition(l.Req.Duration, "created_at", l.Req.StartTime, l.Req.EndTime))

	if err := db.Find(&l.Payments).Error; err != nil {
		return errors.New("获取数据失败")
	}

	return nil
}

func (l *PaymentLogic) get_data() error {
	if len(l.Payments) == 0 {
		l.Res["收入金额"] = decimal.Decimal{}
		l.Res["收入笔数"] = decimal.Decimal{}
		l.Res["支出金额"] = decimal.Decimal{}
		l.Res["支出笔数"] = decimal.Decimal{}

		return nil
	}

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
			}
		case enums.FinanceTypeExpense: // 支出
			{
				price, ok := l.Res["支出金额"].(decimal.Decimal)
				if !ok {
					price = decimal.Decimal{}
				}
				price = price.Add(payment.Amount.Neg())
				l.Res["支出金额"] = price

				num, ok := l.Res["支出笔数"].(decimal.Decimal)
				if !ok {
					num = decimal.Decimal{}
				}
				num = num.Add(decimal.NewFromInt(1))
				l.Res["支出笔数"] = num
			}
		}
	}

	return nil
}
