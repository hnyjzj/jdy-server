package platform

import (
	"jdy/logic/platform/wxwork"
	"jdy/types"
)

type DepartmentLogic struct {
	PlatformLogic
}

// 获取授权链接
func (l *DepartmentLogic) Create(req *types.PlatformDepartmentCreateReq) (int, error) {
	var (
		wxwork wxwork.DepartmentLogic
	)
	wxwork.Ctx = l.Ctx

	return wxwork.Create(req)
}
