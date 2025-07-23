package region

import (
	"errors"
	"jdy/model"
	"jdy/types"

	"github.com/gin-gonic/gin"
)

type RegionLogic struct {
	Ctx   *gin.Context
	Staff *model.Staff
}

// 区域列表
func (l *RegionLogic) List(ctx *gin.Context, req *types.RegionListReq) (*types.PageRes[model.Region], error) {
	var (
		region model.Region

		res types.PageRes[model.Region]
	)

	db := model.DB.Model(&region)
	db = region.WhereCondition(db, &req.Where)
	if err := db.Count(&res.Total).Error; err != nil {
		return nil, errors.New("获取区域列表数量失败")
	}

	db = db.Order("created_at desc")
	db = region.Preloads(db)
	db = model.PageCondition(db, req.Page, req.Limit)

	if err := db.Find(&res.List).Error; err != nil {
		return nil, errors.New("获取区域列表失败")
	}

	return &res, nil
}

// 区域列表
func (l *RegionLogic) My(req *types.RegionListMyReq) (*[]model.Region, error) {

	var (
		staff model.Staff
	)

	db := model.DB.Model(&staff)
	db = db.Where("id = ?", l.Staff.Id)
	db = db.Preload("Regions")
	db = db.Preload("RegionSuperiors")

	if err := db.First(&staff).Error; err != nil {
		return nil, errors.New("获取区域列表失败")
	}

	var region_ids []string
	for _, v := range staff.Regions {
		region_ids = append(region_ids, v.Id)
	}
	for _, v := range staff.RegionSuperiors {
		region_ids = append(region_ids, v.Id)
	}

	var regions []model.Region
	if err := model.DB.Where("id in (?)", region_ids).Find(&regions).Error; err != nil {
		return nil, errors.New("获取区域列表失败")
	}

	def := model.Region{}.Default(l.Staff.Identity)
	if def != nil {
		regions = append([]model.Region{*def}, regions...)
	}

	return &regions, nil
}

// 区域详情
func (l *RegionLogic) Info(ctx *gin.Context, req *types.RegionInfoReq) (*model.Region, error) {
	var (
		region model.Region
	)

	db := model.DB.Model(&model.Region{})

	db = region.Preloads(db)

	if err := db.First(&region, "id = ?", req.Id).Error; err != nil {
		return nil, errors.New("获取区域详情失败")
	}

	return &region, nil
}
