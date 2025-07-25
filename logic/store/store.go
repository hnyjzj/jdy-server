package store

import (
	"errors"
	"jdy/enums"
	"jdy/model"
	"jdy/types"

	"github.com/gin-gonic/gin"
)

type StoreLogic struct {
	Ctx   *gin.Context
	Staff *model.Staff
}

// 门店列表
func (l *StoreLogic) List(ctx *gin.Context, req *types.StoreListReq) (*types.PageRes[model.Store], error) {
	var (
		store model.Store

		res types.PageRes[model.Store]
	)

	db := model.DB.Model(&store)
	db = store.WhereCondition(db, &req.Where)
	if err := db.Count(&res.Total).Error; err != nil {
		return nil, errors.New("获取门店列表数量失败")
	}

	db = model.PageCondition(db, &req.PageReq)

	if err := db.Find(&res.List).Error; err != nil {
		return nil, errors.New("获取门店列表失败")
	}

	return &res, nil
}

// 门店列表
func (l *StoreLogic) My(req *types.StoreListMyReq) (*[]model.Store, error) {

	var (
		staff model.Staff
	)

	db := model.DB.Model(&staff)
	db = db.Where("id = ?", l.Staff.Id)
	db = staff.Preloads(db)

	if err := db.First(&staff).Error; err != nil {
		return nil, errors.New("获取门店列表失败")
	}

	var store_ids []string
	for _, v := range staff.Stores {
		store_ids = append(store_ids, v.Id)
	}
	for _, v := range staff.StoreSuperiors {
		store_ids = append(store_ids, v.Id)
	}
	for _, v := range staff.Regions {
		for _, store := range v.Stores {
			store_ids = append(store_ids, store.Id)
		}
	}
	for _, v := range staff.RegionSuperiors {
		for _, store := range v.Stores {
			store_ids = append(store_ids, store.Id)
		}
	}

	var (
		stores []model.Store
		sdb    = model.DB.Model(&model.Store{})
	)
	if l.Staff.Identity < enums.IdentityAdmin {
		sdb = sdb.Where("id in (?)", store_ids)
	}

	sdb = sdb.Order("`order` asc")
	if err := sdb.Find(&stores).Error; err != nil {
		return nil, errors.New("获取门店列表失败")
	}

	if len(stores) >= 2 {
		def := model.Store{}.Default(l.Staff.Identity)
		if def != nil {
			stores = append([]model.Store{*def}, stores...)
		}
	}

	return &stores, nil
}

// 门店详情
func (l *StoreLogic) Info(ctx *gin.Context, req *types.StoreInfoReq) (*model.Store, error) {
	var (
		store model.Store
	)

	if err := model.DB.
		Preload("Staffs").
		Preload("Superiors").
		First(&store, "id = ?", req.Id).Error; err != nil {
		return nil, errors.New("获取门店详情失败")
	}

	return &store, nil
}
