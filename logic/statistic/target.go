package statistic

import (
	"errors"
	"fmt"
	"jdy/enums"
	"jdy/model"
	"jdy/types"

	"github.com/shopspring/decimal"
)

func (l *StatisticLogic) Target(req *types.StatisticTargetReq) (*types.StatisticTargetResp, error) {
	var target model.Target
	db := model.DB.Model(&model.Target{})
	db = db.Where(&model.Target{
		StoreId:   req.StoreId,
		IsDefault: true,
	})

	tx := model.DB.Model(&model.TargetPersonal{})
	tx = tx.Where(&model.TargetPersonal{
		StaffId: l.Staff.Id,
	})
	tx = tx.Select("target_id")
	db = db.Where("id in (?)", tx)

	db = db.Preload("Personals")
	if err := db.First(&target).Error; err != nil {
		return nil, errors.New("目标获取失败")
	}

	resp := types.StatisticTargetResp{
		TargetId:  target.Id,
		StartTime: target.StartTime,
		EndTime:   target.EndTime,
	}

	var purpose, Achieve decimal.Decimal

	for _, personals := range target.Personals {
		purpose = purpose.Add(personals.Purpose)
		Achieve = Achieve.Add(personals.Achieve)
	}

	switch target.Method {
	case enums.TargetMethodAmount:
		{
			resp.Purpose = fmt.Sprintf("%s 元", purpose.StringFixed(2))
			resp.Remainder = fmt.Sprintf("%s 元", purpose.Sub(Achieve).StringFixed(2))
		}
	case enums.TargetMethodQuantity:
		{
			resp.Purpose = fmt.Sprintf("%s 件", purpose.StringFixed(2))
			resp.Remainder = fmt.Sprintf("%s 件", purpose.Sub(Achieve).StringFixed(2))
		}
	}

	resp.AchieveRate = fmt.Sprintf("%s %%", Achieve.Div(purpose).Mul(decimal.NewFromInt(100)).StringFixed(2))

	return &resp, nil
}
