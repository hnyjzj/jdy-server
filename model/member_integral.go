package model

import (
	"jdy/enums"
	"jdy/types"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// 积分变动记录
type MemberIntegralLog struct {
	SoftDelete

	MemberId string `json:"member_id" gorm:"column:member_id;size:255;not NULL;comment:会员id;"` // 会员id
	Member   Member `json:"member,omitempty" gorm:"foreignKey:MemberId;references:Id;"`        // 会员

	Change     decimal.Decimal                `json:"change" gorm:"column:change;type:decimal(10,2);not NULL;default:0;comment:变动积分;"`     // 变动积分
	ChangeType enums.MemberIntegralChangeType `json:"change_type" gorm:"column:change_type;type:int(11);not NULL;default:0;comment:变动类型;"` // 变动类型
	Before     decimal.Decimal                `json:"before" gorm:"column:before;type:decimal(10,2);not NULL;default:0;comment:变动前积分;"`    // 变动前积分
	After      decimal.Decimal                `json:"after" gorm:"column:after;type:decimal(10,2);not NULL;default:0;comment:变动后积分;"`      // 变动后积分
	Remark     string                         `json:"remark" gorm:"column:remark;size:255;comment:备注;"`                                    // 备注

	OperatorId string `json:"operator_id" gorm:"type:varchar(255);not NULL;comment:操作员ID;"`     // 操作员ID
	Operator   Staff  `json:"operator" gorm:"foreignKey:OperatorId;references:Id;comment:操作员;"` // 操作员
}

func (MemberIntegralLog) WhereCondition(db *gorm.DB, query *types.MemberIntegralWhere) *gorm.DB {
	if query.MemberId != "" {
		db = db.Where("member_id = ?", query.MemberId)
	}
	if query.ChangeType != 0 {
		db = db.Where("change_type = ?", query.ChangeType)
	}

	return db
}

type MemberIntegralRule struct {
	SoftDelete

	Type enums.MemberIntegralRuleType `json:"type" gorm:"column:type;type:int(11);not NULL;default:0;comment:类型;"` // 类型

	Class int             `json:"class" gorm:"column:class;type:int(11);NULL;comment:大类;"`               // 大类
	Rate  decimal.Decimal `json:"rate" gorm:"column:rate;type:decimal(10,2);NULL;default:0;comment:比例;"` // 比例

	OperatorId string `json:"operator_id" gorm:"type:varchar(255);not NULL;comment:操作员ID;"`     // 操作员ID
	Operator   Staff  `json:"operator" gorm:"foreignKey:OperatorId;references:Id;comment:操作员;"` // 操作员
}

func init() {
	// 注册模型
	RegisterModels(
		&MemberIntegralLog{},
		&MemberIntegralRule{},
	)
	// 重置表
	RegisterRefreshModels(
	// &MemberIntegralLog{},
	// &MemberIntegralRule{},
	)
}
