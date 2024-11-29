package store

import (
	"errors"
	"jdy/logic/platform/wxwork"
	"jdy/model"
	"jdy/types"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (l *StoreLogic) Create(ctx *gin.Context, req *types.StoreCreateReq) error {
	var (
		wxwork = wxwork.WxWorkLogic{}
	)

	store := &model.Store{
		ParentId: req.ParentId,

		Name:     req.Name,
		Address:  req.Address,
		Contact:  req.Contact,
		Logo:     req.Logo,
		Sort:     req.Sort,
		Province: req.Province,
		City:     req.City,
		District: req.District,
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 创建企业微信部门
		if req.SyncWxwork {
			id, err := wxwork.StoreCreate(ctx, req)
			if err != nil {
				return errors.New("同步企业微信失败: " + err.Error())
			}
			store.WxworkId = id
		}

		if err := tx.Create(store).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		if req.SyncWxwork && store.WxworkId != 0 {
			if err := wxwork.StoreDelete(ctx, store.WxworkId); err != nil {
				return errors.New("同步企业微信失败: " + err.Error())
			}
		}
		return err
	}

	return nil
}
