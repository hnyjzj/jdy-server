package message

import (
	"fmt"
	"jdy/model"
	"log"
	"strings"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/work/message/request"
)

// 盘点单通知
type ProductInventoryMessage struct {
	ProductInventory *model.ProductInventory
}

// 发送盘点单创建通知
func (M *BaseMessage) SendProductInventoryCreateMessage(req *ProductInventoryMessage) {
	url := M.Url(ProductInventoryInfoUrl, req.ProductInventory.Id)
	var InventoryPersonIds []string
	for _, v := range req.ProductInventory.InventoryPersons {
		InventoryPersonIds = append(InventoryPersonIds, v.Username)
	}
	ToUser := strings.Join(append(InventoryPersonIds, req.ProductInventory.Inspector.Username), "|")
	messages := &request.RequestMessageSendTemplateCard{
		RequestMessageSend: request.RequestMessageSend{
			ToUser:  ToUser,
			MsgType: "template_card",
			AgentID: M.Config.Jdy.Id,
		},
		TemplateCard: &request.RequestTemplateCard{
			CardType: "text_notice",
			MainTitle: &request.TemplateCardMainTitle{
				Title: "盘点单创建通知",
				Desc:  fmt.Sprintf("新盘点单【%s】，请及时处理", req.ProductInventory.Id),
			},
			EmphasisContent: &request.TemplateCardEmphasisContent{
				Title: fmt.Sprintf("%d", req.ProductInventory.ShouldCount),
				Desc:  "应盘数量",
			},
			HorizontalContentList: []*request.TemplateCardHorizontalContentListItem{
				{
					Type:    0,
					Keyname: "总重量",
					Value:   fmt.Sprintf("%s(g)", req.ProductInventory.CountWeightMetal.Round(2).String()),
				},
				{
					Type:    0,
					Keyname: "总价值",
					Value:   fmt.Sprintf("%s(元)", req.ProductInventory.CountPrice.Round(2).String()),
				},
				{
					Type:    0,
					Keyname: "总件数",
					Value:   fmt.Sprintf("%d", req.ProductInventory.CountQuantity),
				},
				{
					Type:    3,
					Keyname: "监盘人",
					Value:   req.ProductInventory.Inspector.Nickname,
					UserID:  req.ProductInventory.Inspector.Username,
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

// 发送盘点单更新通知
func (M *BaseMessage) SendProductInventoryUpdateMessage(req *ProductInventoryMessage) {
	url := M.Url(ProductInventoryInfoUrl, req.ProductInventory.Id)
	var InventoryPersonIds []string
	for _, v := range req.ProductInventory.InventoryPersons {
		InventoryPersonIds = append(InventoryPersonIds, v.Username)
	}
	ToUser := strings.Join(append(InventoryPersonIds, req.ProductInventory.Inspector.Username), "|")
	messages := &request.RequestMessageSendTemplateCard{
		RequestMessageSend: request.RequestMessageSend{
			ToUser:  ToUser,
			MsgType: "template_card",
			AgentID: M.Config.Jdy.Id,
		},
		TemplateCard: &request.RequestTemplateCard{
			CardType: "text_notice",
			MainTitle: &request.TemplateCardMainTitle{
				Title: "盘点单更新通知",
				Desc:  fmt.Sprintf("盘点单【%s】状态更新，请及时处理", req.ProductInventory.Id),
			},
			EmphasisContent: &request.TemplateCardEmphasisContent{
				Title: req.ProductInventory.Status.String(),
				Desc:  "当前状态",
			},
			HorizontalContentList: []*request.TemplateCardHorizontalContentListItem{
				{
					Type:    0,
					Keyname: "应盘数量",
					Value:   fmt.Sprintf("%d", req.ProductInventory.ShouldCount),
				},
				{
					Type:    0,
					Keyname: "实盘数量",
					Value:   fmt.Sprintf("%d", req.ProductInventory.ActualCount),
				},
				{
					Type:    0,
					Keyname: "盘盈数量",
					Value:   fmt.Sprintf("%d", req.ProductInventory.ExtraCount),
				},
				{
					Type:    0,
					Keyname: "盘亏数量",
					Value:   fmt.Sprintf("%d", req.ProductInventory.LossCount),
				},
				{
					Type:    3,
					Keyname: "监盘人",
					Value:   req.ProductInventory.Inspector.Nickname,
					UserID:  req.ProductInventory.Inspector.Username,
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
