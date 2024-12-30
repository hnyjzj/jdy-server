package member

import (
	"errors"
	"jdy/enums"
	"jdy/model"
	"jdy/types"
	"jdy/utils"

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

func (l *MemberLogic) Update(req *types.MemberUpdateReq) error {

	var member model.Member
	if err := model.DB.First(&member, req.Id).Error; err != nil || err == gorm.ErrRecordNotFound {
		return errors.New("更新会员失败")
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {

		data, err := utils.StructToStruct[model.Member](req)
		if err != nil {
			return errors.New("验证信息失败")
		}

		if err := tx.Model(&model.Member{}).Where("id = ?", req.Id).Updates(data).Error; err != nil {
			return errors.New("更新失败")
		}

		return nil
	}); err != nil {
		return errors.New("更新失败")
	}
	return nil
}
