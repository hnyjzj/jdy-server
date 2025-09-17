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
	db = store.Preloads(db)
	if err := db.Find(&res.List).Error; err != nil {
		return nil, errors.New("获取门店列表失败")
	}

	return &res, nil
}

// 门店别名列表
func (l *StoreLogic) Alias(ctx *gin.Context, req *types.StoreAliasReq) ([]model.Store, error) {
	var (
		store model.Store

		res []model.Store
	)

	db := model.DB.Model(&store)

	if req.IsHeadquarters {
		db = db.Where("name LIKE ?", "%"+model.HeaderquartersPrefix+"%")
	} else {
		db = db.Where("name NOT LIKE ?", "%"+model.HeaderquartersPrefix+"%")
		db = db.Where("`alias` <> '' OR `alias` IS NOT NULL")
		db = db.Omit("name")
	}

	if err := db.Find(&res).Error; err != nil {
		return nil, errors.New("获取门店列表失败")
	}

	return res, nil
}

// 门店列表
func (l *StoreLogic) My(req *types.StoreListMyReq) (*[]model.Store, error) {

	var (
		stores []model.Store
		db     = model.DB.Model(&model.Store{})
	)
	if l.Staff.Identity < enums.IdentityAdmin {
		db = db.Where("id in (?)", l.Staff.StoreIds)
	}

	db = db.Order("name desc")
	if err := db.Find(&stores).Error; err != nil {
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
