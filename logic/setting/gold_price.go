package setting

import (
	"errors"
	"jdy/enums"
	"jdy/logic"
	"jdy/message"
	"jdy/model"
	"jdy/types"
	"time"
)

type GoldPriceLogic struct {
	logic.BaseLogic

	IP string
}

func (l *GoldPriceLogic) Get() (*types.GoldPriceGetRes, error) {
	var res types.GoldPriceGetRes

	price, err := model.GetGoldPrice()
	if err != nil {
		return &res, err
	}

	res.Price = price

	return &res, nil
}

func (l *GoldPriceLogic) List(req *types.GoldPriceListReq) (*types.PageRes[model.GoldPrice], error) {
	var (
		data []model.GoldPrice

		res types.PageRes[model.GoldPrice]
	)

	db := model.DB.Model(&data)
	db = db.Where(&model.GoldPrice{Status: enums.GoldPriceStatusApproved})

	// 获取总数
	if err := db.Count(&res.Total).Error; err != nil {
		return nil, errors.New("获取金价历史总数失败")
	}

	// 获取列表
	db = db.Order("created_at desc")
	db = db.Preload("Initiator")
	db = db.Preload("Approver")
	db = model.PageCondition(db, req.Page, req.Limit)
	if err := db.Find(&res.List).Error; err != nil {
		return nil, errors.New("获取金价历史列表失败")
	}

	return &res, nil
}

// 创建金价审批
func (l *GoldPriceLogic) Create(req *types.GoldPriceCreateReq) error {
	data := &model.GoldPrice{
		Price: req.Price,

		InitiatorId: l.Staff.Id,
		IP:          l.IP,

		Status: enums.GoldPriceStatusPending,
	}

	if err := model.DB.Create(data).Error; err != nil {
		return err
	}

	// 发送审批消息
	go func() {
		var initiator model.Staff
		if err := model.DB.Where("id = ?", data.InitiatorId).First(&initiator).Error; err != nil {
			return
		}
		m := message.NewMessage(l.Ctx)
		m.SendGoldPriceApprovalMessage(&message.GoldPriceApprovalMessage{
			Id:        data.Id,
			Price:     req.Price,
			Initiator: initiator.Nickname,
		})
	}()

	return nil
}

// 更新金价
func (l *GoldPriceLogic) Update(req *types.GoldPriceUpdateReq) error {
	// 查询记录
	var (
		price model.GoldPrice
		db    = model.DB.Where("id = ?", req.Id)
	)

	// 查询操作人
	db = db.Preload("Initiator") // 发起人
	db = db.Preload("Approver")  // 审批人

	if err := db.First(&price).Error; err != nil {
		return err
	}

	// 判断审批状态
	if price.Status != enums.GoldPriceStatusPending {
		return errors.New("请勿重复审批")
	}

	// 获取当前时间
	now := time.Now()
	// 更新审批
	if err := model.DB.Model(&price).Updates(model.GoldPrice{
		ApproverId: l.Staff.Id,
		ApprovedAt: &now,
		Status:     enums.GoldPriceStatusApproved,
	}).Error; err != nil {
		return err
	}

	// 发送更新消息
	go func() {
		m := message.NewMessage(l.Ctx)
		m.SendGoldPriceMessage(&message.GoldPriceMessage{
			Price:     price.Price,
			Initiator: price.Initiator.Nickname,
			Approver:  price.Approver.Nickname,
		})
	}()

	return nil
}
