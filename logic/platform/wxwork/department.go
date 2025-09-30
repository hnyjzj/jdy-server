package wxwork

import (
	"errors"
	"jdy/config"
	"jdy/types"
	"log"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/work/department/request"
)

type DepartmentLogic struct {
	WxWorkLogic
}

func (d *DepartmentLogic) Create(req *types.PlatformDepartmentCreateReq) (int, error) {
	if d.Ctx == nil {
		return 0, errors.New("ctx is nil")
	}
	wxwork := config.NewWechatService().ContactsWork.Department
	params := &request.RequestDepartmentInsert{
		Name:     req.Name,
		NameEn:   req.NameEn,
		ParentID: req.ParentId,
		Order:    req.Order,
	}
	// 创建部门
	res, err := wxwork.Create(d.Ctx, params)
	if err != nil || (res != nil && res.ErrCode != 0) {
		log.Printf("创建部门失败: %v, %+v", err, res)
		return 0, errors.New("创建部门失败")
	}

	return res.ID, nil
}
