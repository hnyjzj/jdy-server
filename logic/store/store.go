package store

import (
	"errors"
	"jdy/model"
	"jdy/types"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type StoreLogic struct {
	Ctx   *gin.Context
	Staff *types.Staff
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

	db = db.Order("sort desc, created_at desc")
	db = model.PageCondition(db, req.Page, req.Limit)

	if err := db.Find(&res.List).Error; err != nil {
		return nil, errors.New("获取门店列表失败")
	}

	return &res, nil
}

// 门店列表
func (l *StoreLogic) My(req *types.StoreListMyReq) (*[]model.Store, error) {
	var (
		store model.Store
	)

	var staff model.Staff
	if err := model.DB.Preload("Stores", func(tx *gorm.DB) *gorm.DB {
		db := store.WhereCondition(tx, &req.Where)

		return db
	}).First(&staff, l.Staff.Id).Error; err != nil {
		return nil, errors.New("获取门店列表失败")
	}

	// 如果是管理员
	// if staff.IsAdmin {
	// admin := model.Store{
	// 	SoftDelete: model.SoftDelete{
	// 		Model: model.Model{
	// 			BaseModel: model.BaseModel{
	// 				Id: "",
	// 			},
	// 		},
	// 	},
	// 	Name: "总部",
	// }
	// staff.Stores = append([]model.Store{admin}, staff.Stores...)
	// }

	return &staff.Stores, nil
}

// 门店详情
func (l *StoreLogic) Info(ctx *gin.Context, req *types.StoreInfoReq) (*model.Store, error) {
	var (
		store model.Store
	)

	if err := model.DB.
		Preload("Staffs").
		First(&store, req.Id).Error; err != nil {
		return nil, errors.New("获取门店详情失败")
	}

	return &store, nil
}
