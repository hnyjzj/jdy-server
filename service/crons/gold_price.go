package crons

import (
	"context"
	"jdy/message"
	"jdy/model"
	"log"
)

// 发送金价设置提醒
func SendGoldPriceSetMessage() {
	// 查询所有门店
	var stores []model.Store
	if err := model.DB.Preload("Staffs").Find(&stores).Error; err != nil {
		log.Printf("SendGoldPriceSetMessage: %v\n", err.Error())
		return
	}

	for _, v := range stores {
		if v.Staffs != nil {
			var receiver []string
			for _, staff := range v.Staffs {
				receiver = append(receiver, *staff.Username)
			}
			if len(receiver) == 0 {
				log.Printf("门店 %s 没有有效的接收者，跳过消息发送", v.Name)
				continue
			}
			m := message.NewMessage(context.Background())
			err := m.SendGoldPriceSetMessage(&message.GoldPriceMessage{
				ToUser:    receiver,
				StoreId:   v.Id,
				StoreName: v.Name,
			})
			if err != nil {
				log.Printf("SendGoldPriceSetMessage: %v\n", err.Error())
			} else {
				log.Printf("成功向门店 %s 的 %d 位员工发送金价设置提醒", v.Name, len(receiver))
			}
		}
	}
}
