package target

import (
	"errors"
	"jdy/enums"
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

	switch target.Object {
	case enums.TargetObjectPersonal:
		{
			// 按照达成量排序
			sort.Slice(target.Personals, func(i, j int) bool {
				return target.Personals[i].Achieve.GreaterThan(target.Personals[j].Achieve)
			})
		}
	case enums.TargetObjectGroup:
		{
			for _, p := range target.Personals {
				for gi := range target.Groups {
					for pi := range target.Groups[gi].Personals {
						if target.Groups[gi].Personals[pi].Id == p.Id {
							target.Groups[gi].Personals[pi].Achieve = p.Achieve
							break
						}
					}
				}
			}

			for gi := range target.Groups {
				// 按照达成量排序
				sort.Slice(target.Groups[gi].Personals, func(i, j int) bool {
					return target.Groups[gi].Personals[i].Achieve.GreaterThan(target.Groups[gi].Personals[j].Achieve)
				})
			}
		}
	}

	return &target, nil
}
