package platform

import (
	"errors"
	"jdy/enums"
	"jdy/logic/platform/wxwork"
	"jdy/types"

	"github.com/gin-gonic/gin"
)

// 获取授权链接
func (l *PlatformLogic) GetJSSDK(ctx *gin.Context, req *types.PlatformJSSdkReq) (res any, err error) {
	switch req.Platform {
	case enums.PlatformTypeWxWork:
		var (
			wxwork wxwork.WxWorkLogic
		)
		res, err = wxwork.Jssdk(ctx, req)
		if err != nil {
			return nil, err
		}

	default:
		return nil, errors.New("state error")
	}

	return res, nil
}
