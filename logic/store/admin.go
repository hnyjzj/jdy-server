package store

import (
	"errors"
	"jdy/model"
	"jdy/types"

	"github.com/gin-gonic/gin"
)

type StoreAdminLogic struct {
	Ctx   *gin.Context
	Staff *model.Staff
}

// 门店管理员列表
func (l *StoreAdminLogic) List(req *types.StoreAdminListReq) (*[]model.Staff, error) {
	// 查询门店
	var (
		store model.Store
	)

	db := model.DB.Model(&model.Store{})
	db = store.Preloads(db)

	if err := db.First(&store, "id = ?", req.StoreId).Error; err != nil {
		return nil, errors.New("门店不存在")
	}

	if in := store.InStore(l.Staff.Id); !in {
		return nil, errors.New("未入职该门店，无法查看管理员列表")
	}

	return &store.Admins, nil
}

// 是否在门店
func (l *StoreAdminLogic) IsIn(req *types.StoreAdminIsInReq) (bool, error) {
	var (
		store model.Store
		db    = model.DB

		staff_id string
		res      = false
	)

	db = store.Preloads(db)
	if err := db.First(&store, "id = ?", req.StoreId).Error; err != nil {
		return false, errors.New("门店不存在")
	}

	if req.StaffId != "" {
		staff_id = req.StaffId
	} else {
		staff_id = l.Staff.Id
	}
	for _, staff := range store.Admins {
		if staff.Id == staff_id {
			res = true
			break
		}
	}

	return res, nil
}
