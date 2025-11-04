package target

import (
	"errors"
	"jdy/model"
	"jdy/types"
	"sort"

	"github.com/gin-gonic/gin"
)

type Logic struct {
	Ctx   *gin.Context
	Staff *model.Staff

	Target *model.Target
}

func (l *Logic) List(req *types.TargetListReq) (*types.PageRes[model.Target], error) {
	var (
		target model.Target

		res types.PageRes[model.Target]
	)

	db := model.DB.Model(&target)
	db = target.WhereCondition(db, &req.Where)
	if err := db.Count(&res.Total).Error; err != nil {
		return nil, errors.New("获取目标列表数量失败")
	}

	db = model.PageCondition(db, &req.PageReq)
	db = target.Preloads(db)

	db = db.Order("created_at desc")
	db = db.Order("start_time desc")
	db = db.Order("end_time desc")

	if err := db.Find(&res.List).Error; err != nil {
		return nil, errors.New("获取目标列表失败")
	}

	return &res, nil
}

func (l *Logic) Info(req *types.TargetInfoReq) (*model.Target, error) {
	var (
		target model.Target
	)

	db := model.DB.Model(&target)
	db = target.Preloads(db)
	if err := db.First(&target, "id = ?", req.Id).Error; err != nil {
		return nil, errors.New("获取目标详情失败")
	}

	l.Target = &target
	tg, err := l.GetAchieve(target.Id)
	if err != nil {
		return nil, errors.New("获取目标详情失败")
	}
	target = *tg

	// 按照达成量、目标量排序
	sort.Slice(target.Personals, func(i, j int) bool {
		// 先按照目标量排序，再按照达成量排序
		if target.Personals[i].Purpose.Equal(target.Personals[j].Purpose) {
			return target.Personals[i].Achieve.GreaterThan(target.Personals[j].Achieve)
		}
		return target.Personals[i].Purpose.GreaterThan(target.Personals[j].Purpose)
	})

	return &target, nil
}
