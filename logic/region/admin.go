package region

import (
	"errors"
	"jdy/model"
	"jdy/types"

	"github.com/gin-gonic/gin"
)

type RegionAdminLogic struct {
	Ctx   *gin.Context
	Staff *model.Staff
}

// 区域管理员列表
func (l *RegionAdminLogic) List(req *types.RegionAdminListReq) (*[]model.Staff, error) {
	// 查询区域
	var (
		region model.Region
	)

	db := model.DB.Model(&model.Region{})
	db = region.Preloads(db)

	if err := db.First(&region, "id = ?", req.RegionId).Error; err != nil {
		return nil, errors.New("区域不存在")
	}

	return &region.Admins, nil
}
