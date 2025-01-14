package staff

import (
	"jdy/errors"
	"jdy/logic"
	"jdy/model"
	"jdy/types"
)

type StaffLogic struct {
	logic.Base
	Staff *types.Staff
}

// 员工列表
func (l *StaffLogic) List(req *types.StaffListReq) (*types.PageRes[model.Staff], error) {
	var (
		staff model.Staff

		res types.PageRes[model.Staff]
	)

	db := model.DB.Model(&staff)
	db = staff.WhereCondition(db, &req.Where)
	if err := db.Count(&res.Total).Error; err != nil {
		return nil, errors.New("获取员工列表数量失败")
	}

	db = db.Preload("Stores")
	db = db.Order("created_at desc")
	db = model.PageCondition(db, req.Page, req.Limit)

	if err := db.Find(&res.List).Error; err != nil {
		return nil, errors.New("获取员工列表失败")
	}

	return &res, nil
}

// 员工详情
func (l *StaffLogic) Info(req *types.StaffInfoReq) (*model.Staff, error) {
	var staff model.Staff
	if err := model.DB.Preload("Stores").First(&staff, req.Id).Error; err != nil {
		return nil, errors.ErrStaffNotFound
	}

	return &staff, nil
}

// 获取员工信息
func (l *StaffLogic) My() (*types.StaffRes, error) {
	var staffRes types.StaffRes
	if err := model.DB.Model(&model.Staff{}).First(&staffRes, l.Staff.Id).Error; err != nil {
		return nil, errors.ErrStaffNotFound
	}

	return &staffRes, nil
}
