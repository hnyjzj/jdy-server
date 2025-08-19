package message

import (
	"fmt"
	"jdy/model"
	"log"
	"strings"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/work/message/request"
)

// 调拨单通知
type ProductAllocateMessage struct {
	ProductAllocate *model.ProductAllocate
}

// 发送调拨单创建通知
func (M *BaseMessage) SendProductAllocateCreateMessage(req *ProductAllocateMessage) {
	// 跳转地址
	url := M.Url(ProductAllocateInfoUrl, req.ProductAllocate.Id)
	// 接收消息的人
	to_user, err := req.getToStoreUser()
	if err != nil {
		log.Printf("获取门店用户失败: err=%v\n", err)
		return
	}
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
				Title: "调拨单创建通知",
				Desc:  fmt.Sprintf("新调拨单【%s】，请及时处理", req.ProductAllocate.Id),
			},
			EmphasisContent: &request.TemplateCardEmphasisContent{
				Title: fmt.Sprintf("%d", req.ProductAllocate.ProductCount),
				Desc:  "总件数",
			},
			HorizontalContentList: []*request.TemplateCardHorizontalContentListItem{
				{
					Type:    0,
					Keyname: "总重量",
					Value:   fmt.Sprintf("%s(g)", req.ProductAllocate.ProductTotalWeightMetal.Round(2).String()),
				},
				{
					Type:    0,
					Keyname: "总标签价",
					Value:   fmt.Sprintf("%s(元)", req.ProductAllocate.ProductTotalLabelPrice.Round(2).String()),
				},
				{
					Type:    0,
					Keyname: "来源门店",
					Value:   req.ProductAllocate.FromStore.Name,
				},
				{
					Type:    0,
					Keyname: "目标门店",
					Value:   req.ProductAllocate.ToStore.Name,
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
}

// 发送调拨单取消通知
func (M *BaseMessage) SendProductAllocateCancelMessage(req *ProductAllocateMessage) {
	// 跳转地址
	url := M.Url(ProductAllocateInfoUrl, req.ProductAllocate.Id)
	// 接收消息的人
	to_user, err := req.getToStoreUser()
	if err != nil {
		log.Printf("获取门店用户失败: err=%v\n", err)
		return
	}
	from_user, err := req.getFromStoreUser()
	if err != nil {
		log.Printf("获取门店用户失败: err=%v\n", err)
		return
	}
	ToUser := strings.Join(append(to_user, from_user...), "|")
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
				Title: "调拨单取消通知",
				Desc:  fmt.Sprintf("调拨单【%s】，被取消，请及时处理", req.ProductAllocate.Id),
			},
			EmphasisContent: &request.TemplateCardEmphasisContent{
				Title: fmt.Sprintf("%d", req.ProductAllocate.ProductCount),
				Desc:  "总件数",
			},
			HorizontalContentList: []*request.TemplateCardHorizontalContentListItem{
				{
					Type:    0,
					Keyname: "总重量",
					Value:   fmt.Sprintf("%s(g)", req.ProductAllocate.ProductTotalWeightMetal.Round(2).String()),
				},
				{
					Type:    0,
					Keyname: "总标签价",
					Value:   fmt.Sprintf("%s(元)", req.ProductAllocate.ProductTotalLabelPrice.Round(2).String()),
				},
				{
					Type:    0,
					Keyname: "来源门店",
					Value:   req.ProductAllocate.FromStore.Name,
				},
				{
					Type:    0,
					Keyname: "目标门店",
					Value:   req.ProductAllocate.ToStore.Name,
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
}

// 发送调拨单完成通知
func (M *BaseMessage) SendProductAllocateCompleteMessage(req *ProductAllocateMessage) {
	// 跳转地址
	url := M.Url(ProductAllocateInfoUrl, req.ProductAllocate.Id)
	// 接收消息的人
	to_user, err := req.getToStoreUser()
	if err != nil {
		log.Printf("获取门店用户失败: err=%v\n", err)
		return
	}
	from_user, err := req.getFromStoreUser()
	if err != nil {
		log.Printf("获取门店用户失败: err=%v\n", err)
		return
	}
	ToUser := strings.Join(append(to_user, from_user...), "|")
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
				Title: "调拨单完成通知",
				Desc:  fmt.Sprintf("调拨单【%s】，已完成接收", req.ProductAllocate.Id),
			},
			EmphasisContent: &request.TemplateCardEmphasisContent{
				Title: fmt.Sprintf("%d", req.ProductAllocate.ProductCount),
				Desc:  "总件数",
			},
			HorizontalContentList: []*request.TemplateCardHorizontalContentListItem{
				{
					Type:    0,
					Keyname: "总重量",
					Value:   fmt.Sprintf("%s(g)", req.ProductAllocate.ProductTotalWeightMetal.Round(2).String()),
				},
				{
					Type:    0,
					Keyname: "总标签价",
					Value:   fmt.Sprintf("%s(元)", req.ProductAllocate.ProductTotalLabelPrice.Round(2).String()),
				},
				{
					Type:    0,
					Keyname: "来源门店",
					Value:   req.ProductAllocate.FromStore.Name,
				},
				{
					Type:    0,
					Keyname: "目标门店",
					Value:   req.ProductAllocate.ToStore.Name,
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
}

func (P *ProductAllocateMessage) getToStoreUser() ([]string, error) {
	var (
		to_store model.Store
		db       = model.DB.Model(&model.Store{})
	)
	db = to_store.Preloads(db)
	if err := db.First(&to_store, "id = ?", P.ProductAllocate.ToStoreId).Error; err != nil {
		return nil, err
	}

	var to_user []string
	for _, user := range to_store.Superiors {
		to_user = append(to_user, user.Username)
	}

	return to_user, nil
}

func (P *ProductAllocateMessage) getFromStoreUser() ([]string, error) {
	var (
		from_store model.Store
		db         = model.DB.Model(&model.Store{})
	)
	db = from_store.Preloads(db)
	if err := db.First(&from_store, "id = ?", P.ProductAllocate.FromStoreId).Error; err != nil {
		return nil, err
	}

	var from_user []string
	for _, user := range from_store.Superiors {
		from_user = append(from_user, user.Username)
	}

	return from_user, nil
}
