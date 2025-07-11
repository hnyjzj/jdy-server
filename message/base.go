package message

import (
	"context"
	"errors"
	"fmt"
	"jdy/config"
	"log"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/work"
)

type MessageUrl string

const (
	ProductInventoryInfoUrl          MessageUrl = "/product/check/info?id=%s"               // 货品盘点详情
	ProductAllocateInfoUrl           MessageUrl = "/product/allocate/info?id=%s"            // 货品调拨详情
	ProductAccessorieAllocateInfoUrl MessageUrl = "/product/accessorie/allocate/info?id=%s" // 配件调拨详情
	OrderSalesInfoUrl                MessageUrl = "/sale/sales/order?id=%s"                 // 销售订单详情
)

type BaseMessage struct {
	Ctx context.Context

	WXWork *work.Work         `json:"wxwork"`
	Config *config.WechatWork `json:"config"`
}

func NewMessage(ctx context.Context) *BaseMessage {
	return &BaseMessage{
		Ctx:    ctx,
		WXWork: config.NewWechatService().JdyWork,
		Config: &config.Config.Wechat.Work,
	}
}

func (m *BaseMessage) Send(WXWork *work.Work, messages any) error {
	if res, err := WXWork.Message.Send(m.Ctx, messages); err != nil || (res != nil && res.ErrCode != 0) {
		log.Printf("res: %+v\n", res)
		log.Printf("err: %+v\n", err)
		log.Printf("messages: %+v\n", messages)

		return errors.New(res.ErrMsg)
	}

	return nil
}

func (m *BaseMessage) Url(url MessageUrl, params ...any) string {
	path := m.Config.Jdy.Home + fmt.Sprintf(string(url), params...)

	return path
}
