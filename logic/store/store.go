package store

import (
	"errors"
	"jdy/model"
	"jdy/types"

	"github.com/gin-gonic/gin"
)

type StoreLogic struct{}

type StoreListRes struct {
	Page types.PageRes `json:"page"`
	List []model.Store `json:"list"`
}

// 门店列表
func (l *StoreLogic) List(ctx *gin.Context, req *types.StoreListReq) (*StoreListRes, error) {
	var (
		store model.Store

		res StoreListRes
	)

	db := model.DB.Model(&store)
	db = store.WhereCondition(db, &req.Where)
	if err := db.Count(&res.Page.Total).Error; err != nil {
		return nil, errors.New("获取门店列表失败: " + err.Error())
	}

	db = db.Order("sort desc, created_at desc")
	db = model.PageCondition(db, req.Page, req.Limit)

	if err := db.Find(&res.List).Error; err != nil {
		return nil, errors.New("获取门店列表失败: " + err.Error())
	}

	return &res, nil
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
