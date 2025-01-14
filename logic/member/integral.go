package member

import (
	"errors"
	"jdy/model"
	"jdy/types"
)

func (l *MemberLogic) Integral(req *types.MemberIntegralListReq) (*types.PageRes[model.MemberIntegralLog], error) {
	var (
		integrals model.MemberIntegralLog

		res types.PageRes[model.MemberIntegralLog]
	)

	db := model.DB.Model(&integrals).Where(&model.MemberIntegralLog{MemberId: req.MemberId})

	// 获取总数
	if err := db.Count(&res.Total).Error; err != nil {
		return nil, errors.New("获取列表失败: " + err.Error())
	}

	// 获取列表
	db = db.Order("created_at desc")
	db = model.PageCondition(db, req.Page, req.Limit)
	if err := db.Find(&res.List).Error; err != nil {
		return nil, errors.New("获取列表失败: " + err.Error())
	}

	return &res, nil
}
