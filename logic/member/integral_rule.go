package member

import (
	"errors"
	"jdy/enums"
	"jdy/model"
	"jdy/types"

	"github.com/gin-gonic/gin"
)

type MemberIntegraRulelLogic struct {
	Ctx *gin.Context
}

func (l *MemberIntegraRulelLogic) Finished(req *types.MemberIntegralRuleReq) (*model.MemberIntegralRule, error) {
	var (
		rule model.MemberIntegralRule
	)

	if enums.ProductClass.InMap(enums.ProductClass(req.Class)) != nil {
		return nil, errors.New("类型错误")
	}

	if err := model.DB.Where(&model.MemberIntegralRule{
		Type:  enums.MemberIntegralRuleTypeFinished,
		Class: req.Class,
	}).First(&rule).Error; err != nil {
		return nil, errors.New("没有找到积分规则")
	}

	return &rule, nil
}
func (l *MemberIntegraRulelLogic) Old(req *types.MemberIntegralRuleReq) (*model.MemberIntegralRule, error) {
	var (
		rule model.MemberIntegralRule
	)

	if enums.ProductOldClass.InMap(enums.ProductOldClass(req.Class)) != nil {
		return nil, errors.New("类型错误")
	}

	if err := model.DB.Where(&model.MemberIntegralRule{
		Type:  enums.MemberIntegralRuleTypeOld,
		Class: req.Class,
	}).First(&rule).Error; err != nil {
		return nil, errors.New("没有找到积分规则")
	}

	return &rule, nil
}
