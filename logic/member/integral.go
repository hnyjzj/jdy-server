package member

import (
	"errors"
	"fmt"
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

func (l *MemberIntegralLogic) Change(req *types.MemberIntegralChangeReq) error {
	var (
		member model.Member
	)

	if err := model.DB.Model(&member).Where("id = ?", req.MemberId).First(&member).Error; err != nil {
		return errors.New("获取会员失败: " + err.Error())
	}

	reason := fmt.Sprintf("员工 %s 调整：%s", l.Staff.Phone, req.Reason)

	if err := member.IntegralChange(model.DB, req.Change, enums.MemberIntegralChangeTypeAdjust, reason); err != nil {
		return errors.New("积分变更失败: " + err.Error())
	}

	return nil
}
