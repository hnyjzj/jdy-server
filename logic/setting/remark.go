package setting

import (
	"errors"
	"jdy/logic"
	"jdy/model"
	"jdy/types"
)

type RemarkLogic struct {
	logic.BaseLogic
}

func (l *RemarkLogic) Create(req *types.RemarkCreateReq) error {
	data := model.Remark{
		Content: req.Content,
		StoreId: req.StoreId,

		OperatorId: l.Staff.Id,
		IP:         l.Ctx.ClientIP(),
	}

	if err := model.DB.Create(&data).Error; err != nil {
		return errors.New("创建失败")
	}

	return nil
}

func (l *RemarkLogic) List(req *types.RemarkListReq) (*types.PageRes[model.Remark], error) {
	var (
		remark model.Remark

		res types.PageRes[model.Remark]
	)

	db := model.DB.Model(&remark)
	db = remark.WhereCondition(db, &req.Where)

	// 获取总数
	if err := db.Count(&res.Total).Error; err != nil {
		return nil, errors.New("获取总数失败")
	}

	// 获取列表
	db = db.Order("created_at desc")
	db = model.PageCondition(db, req.Page, req.Limit)
	db = remark.Preloads(db)
	if err := db.Find(&res.List).Error; err != nil {
		return nil, errors.New("获取列表失败")
	}

	return &res, nil
}

func (l *RemarkLogic) Update(req *types.RemarkUpdateReq) error {
	var (
		remark model.Remark
	)

	if err := model.DB.First(&remark, "id = ?", req.Id).Error; err != nil {
		return errors.New("获取备注失败")
	}

	if err := model.DB.Model(&model.Remark{}).Where("id = ?", remark.Id).Updates(model.Remark{
		Content:    req.Content,
		OperatorId: l.Staff.Id,
		IP:         l.Ctx.ClientIP(),
	}).Error; err != nil {
		return errors.New("更新失败")
	}

	return nil
}

func (l *RemarkLogic) Delete(req *types.RemarkDeleteReq) error {
	var (
		remark model.Remark
	)

	if err := model.DB.First(&remark, "id = ?", req.Id).Error; err != nil {
		return errors.New("获取备注失败")
	}

	if err := model.DB.Delete(&remark).Error; err != nil {
		return errors.New("删除失败")
	}

	return nil
}
