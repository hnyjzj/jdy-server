package store

import (
	"errors"
	"jdy/model"

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
