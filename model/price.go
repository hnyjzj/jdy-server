package model

import (
	"errors"
	"jdy/enums"
	"time"
)

type GoldPrice struct {
	SoftDelete

	Price float64 `json:"price" gorm:"type:decimal(10,2);comment:金价;"` // 金价

	InitiatorId string `json:"initiator_id" gorm:"type:varchar(255);not NULL;comment:发起人ID;"`      // 发起人ID
	Initiator   *Staff `json:"initiator" gorm:"foreignKey:InitiatorId;references:Id;comment:发起人;"` // 发起人
	IP          string `json:"ip" gorm:"type:varchar(255);not NULL;comment:IP地址;"`                 // IP地址

	Status enums.GoldPriceStatus `json:"status" gorm:"type:tinyint(1);not null;default:0;comment:状态;"` // 状态

	ApproverId string     `json:"approver_id" gorm:"type:varchar(255);not NULL;comment:审批人ID;"`     // 审批人ID
	Approver   *Staff     `json:"approver" gorm:"foreignKey:ApproverId;references:Id;comment:审批人;"` // 审批人
	ApprovedAt *time.Time `json:"approved_at" gorm:"type:datetime;default:NULL;comment:审批时间;"`      // 审批时间
}

func GetGoldPrice() (float64, error) {
	var price GoldPrice
	db := DB.Model(&GoldPrice{})
	// 排序最新一条
	db = db.Order("created_at desc")
	// 只查询状态为true的数据
	db = db.Where("status = ?", true)
	// 查询数据
	if err := db.First(&price).Error; err != nil {
		return 0, errors.New("获取今日金价失败")
	}

	return price.Price, nil
}

func init() {
	// 注册模型
	RegisterModels(
		&GoldPrice{},
	)
	// 重置表
	RegisterRefreshModels(
	// &GoldPrice{},
	)
}
