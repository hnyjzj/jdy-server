package staff

import (
	"jdy/config"
	"jdy/enums"
	"jdy/errors"
	"jdy/model"
	"jdy/types"
	"log"

	"gorm.io/gorm"
)

// 删除员工
func (l *StaffLogic) StaffDelete(req *types.StaffDeleteReq) error {
	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 判断员工是否存在
		var staff model.Staff
		if err := tx.First(&staff, "id = ?", req.Id).Error; err != nil {
			return errors.New("员工不存在")
		}

		if err := tx.Model(&staff).Association("Stores").Clear(); err != nil {
			return errors.New("店铺离职失败")
		}
		if err := tx.Model(&staff).Association("StoreSuperiors").Clear(); err != nil {
			return errors.New("店铺负责离职失败")
		}
		if err := tx.Model(&staff).Association("StoreAdmins").Clear(); err != nil {
			return errors.New("店铺管理员离职失败")
		}
		if err := tx.Model(&staff).Association("Regions").Clear(); err != nil {
			return errors.New("区域离职失败")
		}
		if err := tx.Model(&staff).Association("RegionSuperiors").Clear(); err != nil {
			return errors.New("区域负责离职失败")
		}
		if err := tx.Model(&staff).Association("RegionAdmins").Clear(); err != nil {
			return errors.New("区域管理员离职失败")
		}

		// 删除员工
		if err := tx.Delete(&staff).Error; err != nil {
			return errors.New("删除员工失败")
		}

		// 删除企业微信用户
		wxwork := config.NewWechatService().ContactsWork
		res, err := wxwork.User.Delete(l.Ctx, staff.Username)
		if err != nil || (res != nil && res.ErrCode != 0) {
			log.Printf("删除企业微信用户失败, err: %+v, %+v", err, res)
		}

		// 添加记录
		log := model.StaffLog{
			Type:       enums.StaffLogTypeDelete,
			StaffId:    staff.Id,
			OldValue:   staff,
			Remark:     "手动删除",
			OperatorId: l.Staff.Id,
		}
		if err := tx.Create(&log).Error; err != nil {
			return errors.New("添加记录失败")
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
