package region

import (
	"errors"
	"jdy/model"
	"jdy/types"

	"github.com/gin-gonic/gin"
)

type RegionStoreLogic struct {
	Ctx   *gin.Context
	Staff *model.Staff
}

// 区域门店列表
func (l *RegionStoreLogic) List(req *types.RegionStoreListReq) (*[]model.Store, error) {
	// 查询区域
	var (
		region model.Region
	)

	db := model.DB.Model(&model.Region{})
	db = region.Preloads(db)
	if err := db.First(&region, "id = ?", req.RegionId).Error; err != nil {
		return nil, errors.New("区域不存在")
	}

	return &region.Stores, nil
}
