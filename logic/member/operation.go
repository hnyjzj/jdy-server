package member

import (
	"errors"
	"jdy/enums"
	"jdy/model"
	"jdy/types"

	"gorm.io/gorm"
)

// 创建会员
func (l *MemberLogic) Create(req *types.MemberCreateReq) error {

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		m := model.Member{
			Phone:       req.Phone,
			Name:        req.Name,
			Gender:      req.Gender,
			Birthday:    req.Birthday,
			Anniversary: req.Anniversary,
			Nickname:    req.Nickname,
			IDCard:      req.IDCard,

			Level:      enums.MemberLevelNone,
			Integral:   0,
			BuyCount:   0,
			EventCount: 0,

			Source:       enums.MemberSourceStaff,
			SourceId:     l.Staff.Id,
			ConsultantId: req.ConsultantId,
			StoreId:      req.StoreId,

			Status: enums.MemberStatusNormal,
		}
		if err := tx.Create(&m).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return errors.New("创建会员失败")
	}

	return nil
}
