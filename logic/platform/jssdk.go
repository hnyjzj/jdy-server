package platform

import (
	"jdy/logic/platform/wxwork"
	"jdy/types"
)

// 获取授权链接
func (l *PlatformLogic) GetJSSDK(req *types.PlatformJSSdkReq) (any, error) {
	var (
		wxwork wxwork.WxWorkLogic
	)

	res, err := wxwork.Jssdk(l.Ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
