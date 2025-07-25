package staff

import (
	"jdy/errors"
	"jdy/logic"
	"jdy/model"
	"jdy/types"
)

type StaffLogic struct {
	logic.BaseLogic
	Staff *model.Staff
}

// 员工列表
func (l *StaffLogic) List(req *types.StaffListReq) (*types.PageRes[model.Staff], error) {
	var (
		staff model.Staff

		res types.PageRes[model.Staff]
	)

	db := model.DB.Model(&staff)
	db = staff.WhereCondition(db, &req.Where)

	if req.Where.StoreId != "" {
		var (
			store model.Store
			sdb   = model.DB.Model(&store)
		)

		sdb = store.Preloads(sdb)
		if err := sdb.First(&store, "id = ?", req.Where.StoreId).Error; err != nil {
			return nil, errors.New("获取门店信息失败")
		}

		var staffIds []string
		for _, s := range store.Staffs {
			staffIds = append(staffIds, s.Id)
		}
		for _, s := range store.Superiors {
			staffIds = append(staffIds, s.Id)
		}

		db = db.Where("id in (?)", staffIds)
	}

	if err := db.Count(&res.Total).Error; err != nil {
		return nil, errors.New("获取员工列表数量失败")
	}

	db = staff.Preloads(db)
	db = db.Order("created_at desc")
	db = model.PageCondition(db, &req.PageReq)

	if err := db.Find(&res.List).Error; err != nil {
		return nil, errors.New("获取员工列表失败")
	}

	return &res, nil
}

// 员工详情
func (l *StaffLogic) Info(req *types.StaffInfoReq) (*model.Staff, error) {
	var staff model.Staff
	db := model.DB.Model(&staff)
	db = staff.Preloads(db)
	if err := db.First(&staff, "id = ?", req.Id).Error; err != nil {
		return nil, errors.ErrStaffNotFound
	}

	return &staff, nil
}

// 获取员工信息
func (l *StaffLogic) My() (*types.StaffRes, error) {
	var staffRes types.StaffRes
	db := model.DB.Model(&model.Staff{})
	db = db.Where("id = ?", l.Staff.Id)

	if err := db.First(&staffRes).Error; err != nil {
		return nil, errors.ErrStaffNotFound
	}

	return &staffRes, nil
}
