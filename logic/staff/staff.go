package staff

import (
	"encoding/json"
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
func (l *StaffLogic) GetStaffInfo(ctx *gin.Context, user *model.Staff) (*types.StaffRes, error) {
	userBytes, err := json.Marshal(user)
	if err != nil {
		return nil, errors.New("获取员工信息失败")
	}

	var saff types.StaffRes
	if err := json.Unmarshal(userBytes, &saff); err != nil {
		return nil, errors.New("获取员工信息失败")
	}

	return &saff, nil
}
