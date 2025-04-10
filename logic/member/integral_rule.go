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

	if enums.ProductClassFinished.InMap(enums.ProductClassFinished(req.Class)) != nil {
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

	if enums.ProductClassOld.InMap(enums.ProductClassOld(req.Class)) != nil {
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
func (l *MemberIntegraRulelLogic) Accessorie(req *types.MemberIntegralRuleReq) ([]model.MemberIntegralRule, error) {
	var (
		rule []model.MemberIntegralRule
	)

	for _, v := range req.Classes {
		if enums.ProductTypePart.InMap(enums.ProductTypePart(v)) != nil {
			return nil, errors.New("类型错误")
		}
	}

	if err := model.DB.
		Where(&model.MemberIntegralRule{
			Type: enums.MemberIntegralRuleTypeAccessorie,
		}).
		Where("class in (?)", req.Classes).
		Find(&rule).Error; err != nil {
		return nil, errors.New("没有找到积分规则")
	}

	return rule, nil
}
