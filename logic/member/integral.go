package member

import (
	"errors"
	"jdy/enums"
	"jdy/model"
	"jdy/types"

	"github.com/gin-gonic/gin"
)

type MemberIntegralLogic struct {
	Ctx   *gin.Context
	Staff *types.Staff
}

func (l *MemberIntegralLogic) List(req *types.MemberIntegralListReq) (*types.PageRes[model.MemberIntegralLog], error) {
	var (
		integrals model.MemberIntegralLog

		res types.PageRes[model.MemberIntegralLog]
	)

	db := model.DB.Model(&integrals).Where(&model.MemberIntegralLog{MemberId: req.MemberId})
	db = integrals.WhereCondition(db, &req.Where)

	// 获取总数
	if err := db.Count(&res.Total).Error; err != nil {
		return nil, errors.New("获取总数失败")
	}

	// 获取列表
	db = db.Order("created_at desc")
	db = model.PageCondition(db, req.Page, req.Limit)
	if err := db.Find(&res.List).Error; err != nil {
		return nil, errors.New("获取列表失败")
	}

	return &res, nil
}

func (l *MemberIntegralLogic) Change(req *types.MemberIntegralChangeReq) error {
	var (
		member model.Member
	)

	if err := model.DB.Model(&member).Where("id = ?", req.MemberId).First(&member).Error; err != nil {
		return errors.New("获取会员失败")
	}

	if err := member.IntegralChange(model.DB, req.Change, enums.MemberIntegralChangeTypeAdjust, req.Remark, l.Staff.Id); err != nil {
		return errors.New("积分变更失败")
	}

	return nil
}
