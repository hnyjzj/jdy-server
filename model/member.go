package model

import (
	"fmt"
	"jdy/enums"
	"jdy/types"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Member struct {
	SoftDelete

	Phone       *string      `json:"phone" gorm:"column:phone;uniqueIndex;size:255;comment:手机号;"` // 手机号
	Name        string       `json:"name" gorm:"column:name;size:255;not NULL;comment:姓名;"`       // 姓名
	Gender      enums.Gender `json:"gender" gorm:"column:gender;type:int(11);comment:性别;"`        // 性别
	Avatar      string       `json:"avatar" gorm:"column:avatar;type:text;comment:头像;"`           // 头像
	Birthday    string       `json:"birthday" gorm:"column:birthday;size:255;comment:生日;"`        // 生日
	Anniversary string       `json:"anniversary" gorm:"column:anniversary;size:255;comment:纪念日;"` // 纪念日
	Nickname    string       `json:"nickname" gorm:"column:nickname;size:255;comment:昵称;"`        // 昵称
	IDCard      string       `json:"id_card" gorm:"column:id_card;size:255;comment:身份证号;"`        // 身份证号

	Level      enums.MemberLevel `json:"level" gorm:"column:level;type:int(11);not NULL;default:0;comment:会员等级;"`             // 会员等级
	Integral   decimal.Decimal   `json:"integral" gorm:"column:integral;type:decimal(15,4);not NULL;default:0;comment:积分;"`   // 积分
	BuyCount   int               `json:"buy_count" gorm:"column:buy_count;type:int(11);not NULL;default:0;comment:购买次数;"`     // 购买次数
	EventCount int               `json:"event_count" gorm:"column:event_count;type:int(11);not NULL;default:0;comment:活动次数;"` // 活动次数

	Source   enums.MemberSource `json:"source" gorm:"column:source;type:int(11);not NULL;default:0;comment:来源;"` // 来源
	SourceId string             `json:"source_id" gorm:"column:source_id;size:255;not NULL;comment:来源id;"`       // 来源id

	ConsultantId string `json:"consultant_id" gorm:"index;column:consultant_id;size:255;not NULL;comment:顾问id;"` // 顾问id
	Consultant   Staff  `json:"consultant,omitempty" gorm:"foreignKey:ConsultantId;references:Id;"`              // 顾问

	StoreId string `json:"store_id" gorm:"index;column:store_id;size:255;not NULL;comment:入会门店id;"` // 入会门店id
	Store   Store  `json:"store,omitempty" gorm:"foreignKey:StoreId;references:Id;"`                // 门店

	Status enums.MemberStatus `json:"status" gorm:"column:status;type:int(11);not NULL;default:0;comment:状态;"` // 状态

	ExternalUserId string `json:"external_user_id" gorm:"index;column:external_user_id;size:255;not NULL;comment:外部用户id;"` // 外部用户id
}

func (Member) WhereCondition(db *gorm.DB, query *types.MemberWhere) *gorm.DB {
	if query.Phone != "" {
		db = db.Where("phone = ?", query.Phone)
	}
	if query.Name != "" {
		db = db.Where("name like ?", fmt.Sprintf("%%%s%%", query.Name))
	}
	if query.Gender != 0 {
		db = db.Where("gender = ?", query.Gender)
	}
	if query.Birthday != "" {
		db = db.Where("birthday = ?", query.Birthday)
	}
	if query.Anniversary != "" {
		db = db.Where("anniversary = ?", query.Anniversary)
	}
	if query.Nickname != "" {
		db = db.Where("nickname like ?", fmt.Sprintf("%%%s%%", query.Nickname))
	}
	if query.Level != 0 {
		db = db.Where("level = ?", query.Level)
	}
	if !query.Integral.IsZero() {
		db = db.Where("integral = ?", query.Integral)
	}
	if query.BuyCount != 0 {
		db = db.Where("buy_count = ?", query.BuyCount)
	}
	if query.EventCount != 0 {
		db = db.Where("event_count = ?", query.EventCount)
	}
	if query.Source != 0 {
		db = db.Where("source = ?", query.Source)
	}
	if query.ConsultantId != "" {
		db = db.Where("consultant_id = ?", query.ConsultantId)
	}
	if query.StoreId != "" {
		db = db.Where("store_id = ?", query.StoreId)
	}
	if query.Status != 0 {
		db = db.Where("status = ?", query.Status)
	}
	if query.ExternalUserId != "" {
		db = db.Where("external_user_id = ?", query.ExternalUserId)
	}

	return db
}

func (M *Member) IntegralChange(db *gorm.DB, change decimal.Decimal, types enums.MemberIntegralChangeType, more ...string) error {
	if change.IsZero() {
		return nil
	}

	// 新积分
	integral := M.Integral.Add(change)
	// 乐观锁更新
	if err := db.Model(&Member{}).Set("gorm:query_option", "FOR UPDATE").Where("id = ?", M.Id).Update("integral", integral).Error; err != nil {
		return err
	}

	log := &MemberIntegralLog{
		MemberId:   M.Id,
		Change:     change,
		ChangeType: types,
		Before:     M.Integral,
		After:      integral,
	}

	if len(more) > 0 && more[0] != "" {
		log.Remark = more[0]
	}
	if len(more) > 1 && more[1] != "" {
		log.OperatorId = more[1]
	}

	if err := db.Create(log).Error; err != nil {
		return err
	}

	return nil
}

func init() {
	// 注册模型
	RegisterModels(
		&Member{},
	)
	// 重置表
	RegisterRefreshModels(
	// &Member{},
	)
}
