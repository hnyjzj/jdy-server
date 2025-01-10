package staff

import (
	"jdy/errors"
	"jdy/logic"
	"jdy/model"
	"jdy/types"

	"github.com/gin-gonic/gin"
)

type StaffLogic struct {
	logic.Base
}

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

// 获取员工信息
func (l *StaffLogic) Info(ctx *gin.Context, user *string) (*types.StaffRes, error) {
	var saffRes types.StaffRes
	if err := model.DB.Model(&model.Staff{}).First(&saffRes, user).Error; err != nil {
		return nil, errors.ErrStaffNotFound
	}

	return &saffRes, nil
}
