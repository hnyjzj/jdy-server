package region

import (
	"errors"
	"jdy/enums"
	"jdy/model"
	"jdy/types"

	"github.com/gin-gonic/gin"
)

type RegionStaffLogic struct {
	Ctx   *gin.Context
	Staff *model.Staff
}

// 区域员工列表
func (l *RegionStaffLogic) List(req *types.RegionStaffListReq) (*[]model.Staff, error) {
	// 查询区域
	var (
		region model.Region
	)

	db := model.DB.Model(&model.Region{})
	db = region.Preloads(db)
	if err := db.First(&region, "id = ?", req.RegionId).Error; err != nil {
		return nil, errors.New("区域不存在")
	}

	if l.Staff.Identity < enums.IdentityAdmin {
		if inRegion := region.InRegion(l.Staff.Id); !inRegion {
			return nil, errors.New("无权查看该门店员工列表")
		}
	}

	return &region.Staffs, nil
}
