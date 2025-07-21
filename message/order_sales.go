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

	// 添加管理员
	to_user, err := req.getStoreSuperiors()
	if err != nil {
		log.Printf("获取门店用户失败: err=%v\n", err)
		return err
	}

	// 获取收银员
	if err := req.getOrderSalesCashier(); err != nil {
		log.Printf("获取收银员失败: err=%v\n", err)
		return err
	}

	// 添加收银员
	to_user = append(to_user, req.OrderSales.Cashier.Username)

	ToUser := strings.Join(to_user, "|")

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
					Value:   fmt.Sprintf("%s(元)", req.OrderSales.PriceOriginal.Round(2).String()),
				},
				{
					Type:    0,
					Keyname: "应付金额",
					Value:   fmt.Sprintf("%s(元)", req.OrderSales.Price.Round(2).String()),
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
				{
					Type:    3,
					Keyname: "开单人",
					Value:   req.OrderSales.Operator.Nickname,
					UserID:  req.OrderSales.Operator.Username,
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

// 发送销售单支付消息通知
func (M *BaseMessage) SendOrderSalesPayMessage(req *OrderSalesMessage) error {
	url := M.Url(OrderSalesInfoUrl, req.OrderSales.Id)

	// 接收消息的人
	var to_user []string
	for _, user := range req.OrderSales.Clerks {
		to_user = append(to_user, user.Salesman.Username)

	}
	// 添加管理员
	to_superiors, err := req.getStoreSuperiors()
	if err != nil {
		log.Printf("获取门店用户失败: err=%v\n", err)
		return err
	}
	to_user = append(to_user, to_superiors...)

	// 添加收银员
	to_user = append(to_user, req.OrderSales.Cashier.Username)

	ToUser := strings.Join(to_user, "|")

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
				Title: "销售单支付通知",
				Desc:  fmt.Sprintf("销售单【%s】，已支付，可前往打印出票", req.OrderSales.Id),
			},
			EmphasisContent: &request.TemplateCardEmphasisContent{
				Title: fmt.Sprintf("%d", len(req.OrderSales.Products)),
				Desc:  "总件数",
			},
			HorizontalContentList: []*request.TemplateCardHorizontalContentListItem{
				{
					Type:    0,
					Keyname: "原价",
					Value:   fmt.Sprintf("%s(元)", req.OrderSales.PriceOriginal.Round(2).String()),
				},
				{
					Type:    0,
					Keyname: "应付金额",
					Value:   fmt.Sprintf("%s(元)", req.OrderSales.Price.Round(2).String()),
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
				{
					Type:    3,
					Keyname: "操作人",
					Value:   req.OrderSales.Operator.Nickname,
					UserID:  req.OrderSales.Operator.Username,
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

// 发送销售单退款消息通知
func (M *BaseMessage) SendOrderSalesRefundMessage(req *OrderSalesMessage) error {
	url := M.Url(OrderSalesInfoUrl, req.OrderSales.Id)

	// 接收消息的人
	var to_user []string
	for _, user := range req.OrderSales.Clerks {
		to_user = append(to_user, user.Salesman.Username)

	}
	// 添加管理员
	to_superiors, err := req.getStoreSuperiors()
	if err != nil {
		log.Printf("获取门店用户失败: err=%v\n", err)
		return err
	}
	to_user = append(to_user, to_superiors...)

	ToUser := strings.Join(to_user, "|")

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
				Title: "销售单退款通知",
				Desc:  fmt.Sprintf("销售单【%s】，有退货，请注意查看", req.OrderSales.Id),
			},
			HorizontalContentList: []*request.TemplateCardHorizontalContentListItem{
				{
					Type:    3,
					Keyname: "操作人",
					Value:   req.OrderSales.Operator.Nickname,
					UserID:  req.OrderSales.Operator.Username,
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

// 发送销售单取消消息通知
func (M *BaseMessage) SendOrderSalesCancelMessage(req *OrderSalesMessage) error {
	url := M.Url(OrderSalesInfoUrl, req.OrderSales.Id)

	// 接收消息的人
	var to_user []string
	for _, user := range req.OrderSales.Clerks {
		to_user = append(to_user, user.Salesman.Username)

	}
	// 添加管理员
	to_superiors, err := req.getStoreSuperiors()
	if err != nil {
		log.Printf("获取门店用户失败: err=%v\n", err)
		return err
	}
	to_user = append(to_user, to_superiors...)

	// 添加收银员
	to_user = append(to_user, req.OrderSales.Cashier.Username)

	ToUser := strings.Join(to_user, "|")

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
				Title: "销售单取消通知",
				Desc:  fmt.Sprintf("销售单【%s】，已取消，请注意查看", req.OrderSales.Id),
			},
			HorizontalContentList: []*request.TemplateCardHorizontalContentListItem{
				{
					Type:    3,
					Keyname: "操作人",
					Value:   req.OrderSales.Operator.Nickname,
					UserID:  req.OrderSales.Operator.Username,
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

func (req *OrderSalesMessage) getStoreSuperiors() ([]string, error) {
	var (
		store model.Store
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

func (req *OrderSalesMessage) getOrderSalesCashier() error {
	if req.OrderSales.Cashier.Username != "" {
		return nil
	}

	var staff model.Staff
	if err := model.DB.First(&staff, "id = ?", req.OrderSales.CashierId).Error; err != nil {
		return err
	}

	req.OrderSales.Cashier = staff

	return nil
}
