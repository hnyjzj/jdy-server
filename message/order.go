package message

import (
	"fmt"
	"jdy/model"
	"log"
	"strings"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/work/message/request"
)

type OrderSalesMessage struct {
	OrderSales *model.OrderSales
}

// 发送销售单创建消息
func (M *BaseMessage) SendOrderSalesCreateMessage(req *OrderSalesMessage) error {
	url := M.Url(OrderSalesInfoUrl, req.OrderSales.Id)

	// 接收消息的人
	to_user, err := req.getStoreUser()
	if err != nil {
		log.Printf("获取门店用户失败: err=%v\n", err)
		return err
	}
	ToUser := strings.Join(append(to_user, req.OrderSales.Cashier.Username), "|")

	// 消息内容
	messages := &request.RequestMessageSendTemplateCard{
		RequestMessageSend: request.RequestMessageSend{
			ToUser:  ToUser,
			MsgType: "template_card",
			AgentID: M.Config.Jdy.Id,
		},
		TemplateCard: &request.RequestTemplateCard{
			CardType: "text_notice",
			MainTitle: &request.TemplateCardMainTitle{
				Title: "销售单创建通知",
				Desc:  fmt.Sprintf("新销售单【%s】，请及时处理", req.OrderSales.Id),
			},
			EmphasisContent: &request.TemplateCardEmphasisContent{
				Title: fmt.Sprintf("%d", len(req.OrderSales.Products)),
				Desc:  "总件数",
			},
			HorizontalContentList: []*request.TemplateCardHorizontalContentListItem{
				{
					Type:    0,
					Keyname: "原价",
					Value:   fmt.Sprintf("%s(g)", req.OrderSales.PriceOriginal.Round(2).String()),
				},
				{
					Type:    0,
					Keyname: "应付金额",
					Value:   fmt.Sprintf("%s(g)", req.OrderSales.Price.Round(2).String()),
				},
				{
					Type:    0,
					Keyname: "实付金额",
					Value:   fmt.Sprintf("%s(元)", req.OrderSales.PricePay.Round(2).String()),
				},
				{
					Type:    0,
					Keyname: "优惠金额",
					Value:   fmt.Sprintf("%s(元)", req.OrderSales.PriceDiscount.Round(2).String()),
				},
			},
			CardAction: &request.TemplateCardAction{
				Type: 1,
				Url:  url,
			},
			JumpList: []*request.TemplateCardJumpListItem{
				{
					Type:  1,
					Title: "查看详情",
					Url:   url,
				},
			},
		},
	}

	if res, err := M.WXWork.Message.SendTemplateCard(M.Ctx, messages); err != nil || (res != nil && res.ErrCode != 0) {
		log.Printf("发送消息失败: err=%v, res=%+v\n", err, res)
	}
	return nil
}

func (req *OrderSalesMessage) getStoreUser() ([]string, error) {
	var (
		store *model.Store
		db    = model.DB.Model(&model.Store{})
	)
	db = store.Preloads(db)
	if err := db.First(&store, "id = ?", req.OrderSales.StoreId).Error; err != nil {
		return nil, err
	}

	var user []string
	for _, u := range store.Superiors {
		user = append(user, u.Username)
	}

	return user, nil
}
