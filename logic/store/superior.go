package store

import (
	"errors"
	"jdy/enums"
	"jdy/model"
	"jdy/types"

	"github.com/gin-gonic/gin"
)

type StoreSuperiorLogic struct {
	Ctx   *gin.Context
	Staff *model.Staff
}

// 门店员工列表
func (l *StoreSuperiorLogic) List(req *types.StoreSuperiorListReq) (*[]model.Staff, error) {
	// 查询门店
	var (
		store model.Store
	)

	db := model.DB.Model(&model.Store{})
	db = store.Preloads(db)

	if err := db.First(&store, "id = ?", req.StoreId).Error; err != nil {
		return nil, errors.New("门店不存在")
	}

	if l.Staff.Identity < enums.IdentityAdmin && !store.InStore(l.Staff.Id) {
		return nil, errors.New("未入职该门店，无法查看员工列表")
	}

	return &store.Superiors, nil
}

// 是否在门店
func (l *StoreSuperiorLogic) IsIn(req *types.StoreSuperiorIsInReq) (bool, error) {
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
	for _, staff := range store.Superiors {
		if staff.Id == staff_id {
			res = true
			break
		}
	}

	return res, nil
}
