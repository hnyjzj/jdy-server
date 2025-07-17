package member

import (
	"errors"
	"jdy/model"
	"jdy/types"

	"github.com/gin-gonic/gin"
)

type MemberLogic struct {
	Ctx   *gin.Context
	Staff *model.Staff
}

func (l *MemberLogic) List(req *types.MemberListReq) (*types.PageRes[model.Member], error) {
	var (
		member model.Member

		res types.PageRes[model.Member]
	)

	db := model.DB.Model(&member)
	db = member.WhereCondition(db, &req.Where)

	// 获取总数
	if err := db.Count(&res.Total).Error; err != nil {
		return nil, errors.New("获取列表失败: " + err.Error())
	}

	// 获取列表
	db = db.Preload("Store").Preload("Consultant")
	db = db.Order("created_at desc")
	db = model.PageCondition(db, req.Page, req.Limit)
	if err := db.Find(&res.List).Error; err != nil {
		return nil, errors.New("获取列表失败: " + err.Error())
	}

	return &res, nil
}

func (l *MemberLogic) Info(req *types.MemberInfoReq) (*model.Member, error) {
	var (
		member model.Member
	)

	db := model.DB.Model(&member)

	if req.Id != "" {
		db = db.Where("id = ?", req.Id)
	}
	if req.ExternalUserId != "" {
		db = db.Or(&model.Member{ExternalUserId: req.ExternalUserId})
	}

	db = db.Preload("Store")
	db = db.Preload("Consultant")

	if err := db.First(&member).Error; err != nil {
		return nil, errors.New("获取信息失败")
	}

	return &member, nil
}

func (l *MemberLogic) Consumptions(req *types.MemberConsumptionsReq) (*types.PageRes[model.Member], error) {

	return nil, nil
}
