package callback

import (
	"jdy/enums"
	"jdy/logic"
	"jdy/logic/setting"
	"jdy/message"
	"jdy/types"
	"log"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/work/server/handlers/models"
)

type TemplateCardEvent struct {
	Handle  *WxWork                        // 处理器
	Message *models.EventTemplateCardEvent // 消息体
}

func (Handle *WxWork) TemplateCardEvent() any {
	var (
		l = TemplateCardEvent{
			Handle: Handle,
		}
	)

	// 获取员工信息
	if err := Handle.GetStaff(); err != nil {
		log.Printf("TemplateCardEvent.GetStaff.Error(): %v\n", err.Error())
		return "error"
	}

	// 解析消息体
	if err := l.Handle.Event.ReadMessage(&l.Message); err != nil {
		log.Printf("TemplateCardEvent.ReadMessage.Error(): %v\n", err.Error())
		return "error"
	}

	switch Handle.Event.GetEventKey() {
	case string(enums.GoldPriceReviewApproved):
		{
			if err := l.gold_price_review_approved(); err != nil {
				return "error"
			}
		}
	case string(enums.GoldPriceReviewRejected):
		{
			if err := l.gold_price_review_rejected(); err != nil {
				return "error"
			}
		}
	}

	return nil
}

// 审批通过
func (l *TemplateCardEvent) gold_price_review_approved() error {
	var (
		logic = setting.GoldPriceLogic{
			BaseLogic: logic.BaseLogic{
				Ctx:   l.Handle.Ctx,
				Staff: l.Handle.Staff,
			},
		}
	)

	if err := logic.Update(&types.GoldPriceUpdateReq{
		Id:     l.Message.TaskID,
		Status: enums.GoldPriceStatusApproved,
	}); err != nil {
		log.Printf("gold_price_review_approved.Error(): %v\n", err.Error())
		return err
	}

	go func() {
		m := message.NewMessage(logic.Ctx)
		m.UpdateGoldPriceMessage(&message.UpdateGoldPriceMessage{
			Code:    l.Message.ResponseCode,
			Message: "已通过",
		})
	}()

	return nil
}

// 审批拒绝
func (l *TemplateCardEvent) gold_price_review_rejected() error {
	var (
		logic = setting.GoldPriceLogic{
			BaseLogic: logic.BaseLogic{
				Ctx:   l.Handle.Ctx,
				Staff: l.Handle.Staff,
			},
		}
	)

	if err := logic.Update(&types.GoldPriceUpdateReq{
		Id:     l.Message.TaskID,
		Status: enums.GoldPriceStatusRejected,
	}); err != nil {
		log.Printf("gold_price_review_rejected.Error(): %v\n", err.Error())
		return err
	}

	go func() {
		m := message.NewMessage(logic.Ctx)
		m.UpdateGoldPriceMessage(&message.UpdateGoldPriceMessage{
			Code:    l.Message.ResponseCode,
			Message: "已拒绝",
		})
	}()

	return nil
}
