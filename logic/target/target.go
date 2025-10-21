package target

import (
	"errors"
	"jdy/model"
	"jdy/types"

	"github.com/gin-gonic/gin"
)

type Logic struct {
	Ctx   *gin.Context
	Staff *model.Staff
}

func (l *Logic) List(req *types.TargetListReq) (*types.PageRes[model.Target], error) {
	var (
		target model.Target

		res types.PageRes[model.Target]
	)

	db := model.DB.Model(&target)
	db = target.WhereCondition(db, &req.Where)
	if err := db.Count(&res.Total).Error; err != nil {
		return nil, errors.New("获取门店列表数量失败")
	}

	db = model.PageCondition(db, &req.PageReq)
	db = target.Preloads(db)
	if err := db.Find(&res.List).Error; err != nil {
		return nil, errors.New("获取门店列表失败")
	}

	return &res, nil
}
