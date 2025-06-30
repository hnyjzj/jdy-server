package platform

import (
	"jdy/logic/platform/wxwork"
	"jdy/types"
)

// 获取授权链接
func (l *PlatformLogic) OauthUri(req *types.PlatformOAuthReq) (*types.PlatformOAuthRes, error) {
	var (
		wxwork wxwork.WxWorkLogic
	)
	res, err := wxwork.OauthUri(req.Agent, req.Uri)
	if err != nil {
		return nil, err
	}
	result := &types.PlatformOAuthRes{
		RedirectURL: res,
	}

	return result, nil
}
