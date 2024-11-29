package store

import (
	"errors"
	"jdy/model"
	"jdy/types"

	"github.com/gin-gonic/gin"
)

type StoreLogic struct{}

// 门店列表
func (l *StoreLogic) List(ctx *gin.Context) ([]*model.Store, error) {
	var (
		store model.Store
	)

	list, err := store.GetTree(nil)
	if err != nil {
		return nil, errors.New("获取门店列表失败: " + err.Error())
	}

	return list, nil
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
