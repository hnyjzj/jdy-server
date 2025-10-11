package region

import (
	"errors"
	"jdy/enums"
	"jdy/model"
	"jdy/types"

	"github.com/gin-gonic/gin"
)

type RegionSuperiorLogic struct {
	Ctx   *gin.Context
	Staff *model.Staff
}

// 区域负责人列表
func (l *RegionSuperiorLogic) List(req *types.RegionSuperiorListReq) (*[]model.Staff, error) {
	// 查询区域
	var (
		region model.Region
	)

	db := model.DB.Model(&model.Region{})
	db = region.Preloads(db)

	if err := db.First(&region, "id = ?", req.RegionId).Error; err != nil {
		return nil, errors.New("区域不存在")
	}

	if l.Staff.Identity < enums.IdentityAdmin && !region.InRegion(l.Staff.Id) {
		return nil, errors.New("未入职该区域，无法查看负责人列表")
	}

	return &region.Superiors, nil
}
