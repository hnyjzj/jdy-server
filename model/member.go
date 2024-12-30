package model

import (
	"jdy/enums"

	"gorm.io/gorm"
)

type Member struct {
	SoftDelete

	Phone       *string      `json:"phone" gorm:"column:phone;unique;size:255;not NULL;comment:手机号;"` // 手机号
	Name        string       `json:"name" gorm:"column:name;size:255;not NULL;comment:姓名;"`           // 姓名
	Gender      enums.Gender `json:"gender" gorm:"column:gender;type:tinyint(1);comment:性别;"`         // 性别
	Birthday    string       `json:"birthday" gorm:"column:birthday;size:255;comment:生日;"`            // 生日
	Anniversary string       `json:"anniversary" gorm:"column:anniversary;size:255;comment:纪念日;"`     // 纪念日
	Nickname    string       `json:"nickname" gorm:"column:nickname;size:255;comment:昵称;"`            // 昵称
	IDCard      string       `json:"id_card" gorm:"column:id_card;size:255;comment:身份证号;"`            // 身份证号

	Level      enums.MemberLevel `json:"level" gorm:"column:level;type:tinyint(1);not NULL;default:0;comment:会员等级;"`          // 会员等级
	Integral   float64           `json:"integral" gorm:"column:integral;type:decimal(10,2);not NULL;default:0;comment:积分;"`   // 积分
	BuyCount   int               `json:"buy_count" gorm:"column:buy_count;type:int(11);not NULL;default:0;comment:购买次数;"`     // 购买次数
	EventCount int               `json:"event_count" gorm:"column:event_count;type:int(11);not NULL;default:0;comment:活动次数;"` // 活动次数

	Source   enums.MemberSource `json:"source" gorm:"column:source;type:tinyint(1);not NULL;default:0;comment:来源;"` // 来源
	SourceId string             `json:"source_id" gorm:"column:source_id;size:255;not NULL;comment:来源id;"`          // 来源id

	ConsultantId string `json:"consultant_id" gorm:"column:consultant_id;size:255;not NULL;comment:顾问id;"` // 顾问id
	Consultant   Staff  `json:"consultant" gorm:"foreignKey:ConsultantId;references:Id;"`                  // 顾问

	StoreId string `json:"store_id" gorm:"column:store_id;size:255;not NULL;comment:入会门店id;"` // 入会门店id
	Store   Store  `json:"store" gorm:"foreignKey:StoreId;references:Id;"`                    // 门店

	Status enums.MemberStatus `json:"status" gorm:"column:status;type:tinyint(1);not NULL;default:0;comment:状态;"` // 状态
}

func (M *Member) IntegralChange(db *gorm.DB, change float64, types enums.MemberIntegralChangeType, remark ...string) error {
	if change == 0 {
		return nil
	}

	M.Integral += change
	if err := db.Save(M).Error; err != nil {
		return err
	}

	log := &MemberIntegralLog{
		MemberId:   M.Id,
		Before:     M.Integral - change,
		After:      M.Integral,
		ChangeType: types,
	}
	if len(remark) > 0 {
		log.Remark = remark[0]
	}
	if err := db.Create(log).Error; err != nil {
		return err
	}

	return nil
}

// 积分变动记录
type MemberIntegralLog struct {
	SoftDelete

	MemberId string `json:"memberId" gorm:"column:member_id;size:255;not NULL;comment:会员id;"` // 会员id
	Member   Member `json:"member" gorm:"foreignKey:MemberId;references:Id;"`                 // 会员

	Before     float64                        `json:"before" gorm:"column:before;type:decimal(10,2);not NULL;default:0;comment:变动前积分;"`       // 变动前积分
	After      float64                        `json:"after" gorm:"column:after;type:decimal(10,2);not NULL;default:0;comment:变动后积分;"`         // 变动后积分
	ChangeType enums.MemberIntegralChangeType `json:"change_type" gorm:"column:change_type;type:tinyint(1);not NULL;default:0;comment:变动类型;"` // 变动类型
	Remark     string                         `json:"remark" gorm:"column:remark;size:255;comment:备注;"`                                       // 备注
}

func init() {
	// 注册模型
	RegisterModels(
		&Member{},
		&MemberIntegralLog{},
	)
	// 重置表
	RegisterRefreshModels(
	// &Member{},
	// &MemberIntegralLog{},
	)
}