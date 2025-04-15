package platform

import (
	"errors"
	"jdy/enums"
	"jdy/logic/platform/wxwork"
	"jdy/types"
)

// 获取授权链接
func (l *PlatformLogic) OauthUri(req *types.PlatformOAuthReq) (*types.PlatformOAuthRes, error) {
	switch req.Platform {
	case enums.PlatformTypeWxWork:

		var (
			wxwork wxwork.WxWorkLogic
		)
		res, err := wxwork.OauthUri(req.Agent, req.Platform.String(), req.Uri)
		if err != nil {
			return nil, err
		}
		result := &types.PlatformOAuthRes{
			RedirectURL: res,
		}

		return result, nil

	default:
		return nil, errors.New("不支持的授权类型")
	}
}
