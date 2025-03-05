package scripts

import (
	"context"
	"fmt"
	"jdy/enums"
	"jdy/message"
	"jdy/model"

	"gorm.io/gorm"
)

// 发送金价设置提醒
func SendGoldPriceSetMessage() {
	// 查询所有门店
	var stores []model.Store
	if err := model.DB.Preload("Staffs", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Account", func(db *gorm.DB) *gorm.DB {
			return db.Where(&model.Account{Platform: enums.PlatformTypeWxWork})
		})
	}).Find(&stores).Error; err != nil {
		fmt.Printf("err.Error(): %v\n", err.Error())
		return
	}

	for _, v := range stores {
		if v.Staffs != nil {
			var receiver []string
			for _, staff := range v.Staffs {
				if staff.Account != nil && staff.Account.Username != nil {
					receiver = append(receiver, *staff.Account.Username)
				}
			}
			m := message.NewMessage(context.Background())
			m.SendGoldPriceSetMessage(&message.GoldPriceMessage{
				ToUser:    receiver,
				StoreName: v.Name,
			})
		}
	}
}
