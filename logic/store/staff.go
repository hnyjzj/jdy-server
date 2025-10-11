package store

import (
	"errors"
	"jdy/enums"
	"jdy/model"
	"jdy/types"

	"github.com/gin-gonic/gin"
)

type StoreStaffLogic struct {
	Ctx   *gin.Context
	Staff *model.Staff
}

// 门店员工列表
func (l *StoreStaffLogic) List(req *types.StoreStaffListReq) (*[]model.Staff, error) {
	// 查询门店
	var (
		store model.Store
	)
	db := model.DB.Model(&store)
	db = store.Preloads(db)

	if err := db.First(&store, "id = ?", req.StoreId).Error; err != nil {
		return nil, errors.New("门店不存在")
	}

	if l.Staff.Identity < enums.IdentityAdmin {
		if inStore := store.InStore(l.Staff.Id); !inStore {
			return nil, errors.New("无权查看该门店员工列表")
		}
	}

	return &store.Staffs, nil
}

// 是否在门店
func (l *StoreStaffLogic) IsIn(req *types.StoreStaffIsInReq) (bool, error) {
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

	for _, staff := range store.Staffs {
		if staff.Id == staff_id {
			res = true
			break
		}
	}

	return res, nil
}
