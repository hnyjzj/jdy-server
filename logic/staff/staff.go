package staff

import (
	"jdy/errors"
	"jdy/logic"
	"jdy/model"
	"jdy/types"

	"github.com/gin-gonic/gin"
)

type StaffLogic struct {
	logic.Base
}

// 获取员工信息
func (l *StaffLogic) GetStaffInfo(ctx *gin.Context, user *string) (*types.StaffRes, error) {
	var saffRes types.StaffRes
	if err := model.DB.Model(&model.Staff{}).First(&saffRes, user).Error; err != nil {
		return nil, errors.ErrStaffNotFound
	}

	return &saffRes, nil
}
