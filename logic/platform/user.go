package platform

import (
	"jdy/logic/platform/wxwork"
	"jdy/model"
	"jdy/types"
)

// 获取授权链接
func (l *PlatformLogic) GetUser(req *types.PlatformGetUserReq) (*model.Staff, error) {
	var (
		wxwork = &wxwork.WxWorkLogic{
			Ctx: l.Ctx,
		}
	)

	res, err := wxwork.GetUser(req.UserId)
	if err != nil {
		return nil, err
	}

	return res, nil
}
