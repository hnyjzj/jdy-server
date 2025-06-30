package setting

import (
	"errors"
	"jdy/logic"
	"jdy/model"
	"jdy/types"
	"log"

	"gorm.io/gorm"
)

type RoleLogic struct {
	logic.BaseLogic
	IP string
}

func (r *RoleLogic) Create(req *types.RoleCreateReq) (*model.Role, error) {
	role := model.Role{
		Name:       req.Name,
		Desc:       req.Desc,
		OperatorId: r.Staff.Id,
		IP:         r.IP,
	}
	if err := model.DB.Create(&role).Error; err != nil {
		return nil, err
	}

	return &role, nil
}

func (r *RoleLogic) List() ([]model.Role, error) {
	var (
		roles []model.Role
	)

	if err := model.DB.Find(&roles).Error; err != nil {
		return nil, err
	}

	return roles, nil
}

func (r *RoleLogic) Info(req *types.RoleInfoReq) (*model.Role, error) {
	var (
		role model.Role
	)
	db := model.DB.Model(&model.Role{})

	db = role.Preloads(db, nil)

	if err := db.First(&role, "id = ?", req.Id).Error; err != nil {
		log.Printf("查询角色失败: %v", err)
		return nil, errors.New("查询角色失败")
	}

	return &role, nil
}

func (r *RoleLogic) Update(req *types.RoleUpdateReq) error {
	var (
		role model.Role
	)
	if err := model.DB.First(&role, "id = ?", req.Id).Error; err != nil {
		log.Printf("查询角色失败: %v", err)
		return errors.New("查询角色失败")
	}

	// 更新角色信息
	data := model.Role{
		Name:       req.Name,
		Desc:       req.Desc,
		OperatorId: r.Staff.Id,
		IP:         r.IP,
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&role).Updates(data).Error; err != nil {
			return err
		}

		var routers []model.Router
		if err := tx.Model(&model.Router{}).Where("id in (?)", req.Routers).Find(&routers).Error; err != nil {
			return err
		}
		if err := tx.Model(&role).Association("Routers").Replace(routers); err != nil {
			return err
		}

		var Apis []model.Api
		if err := tx.Model(&model.Api{}).Where("id in (?)", req.Apis).Find(&Apis).Error; err != nil {
			return err
		}
		if err := tx.Model(&role).Association("Apis").Replace(Apis); err != nil {
			return err
		}

		var Stores []model.Store
		if err := tx.Model(&model.Store{}).Where("id in (?)", req.Stores).Find(&Stores).Error; err != nil {
			return err
		}
		if err := tx.Model(&role).Association("Stores").Replace(Stores); err != nil {
			return err
		}

		var Staffs []model.Staff
		if err := tx.Model(&model.Staff{}).Where("id in (?)", req.Staffs).Find(&Staffs).Error; err != nil {
			return err
		}
		if err := tx.Model(&role).Association("Staffs").Replace(Staffs); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return errors.New("更新角色失败")
	}
	return nil
}

func (r *RoleLogic) AddStaff(req *types.RoleAddStaffReq) error {
	var (
		role   model.Role
		staffs []model.Staff
	)

	if err := model.DB.Model(&model.Role{}).First(&role, "id = ?", req.Id).Error; err != nil {
		return errors.New("查询角色失败")
	}

	if err := model.DB.Model(&model.Staff{}).Where("username in (?)", req.Staffs).Find(&staffs).Error; err != nil {
		return errors.New("查询员工失败")
	}

	if err := model.DB.Model(&role).Association("Staffs").Append(staffs); err != nil {
		return errors.New("添加员工失败")
	}

	return nil
}

func (r *RoleLogic) Delete(req *types.RoleDeleteReq) error {
	var (
		role model.Role
	)
	if err := model.DB.First(&role, "id = ?", req.Id).Error; err != nil {
		return errors.New("查询角色失败")
	}

	if err := model.DB.Delete(&role).Error; err != nil {
		return errors.New("删除角色失败")
	}

	return nil
}

func (r *RoleLogic) Apis() (any, error) {
	path := "/api"
	list, err := model.Api{}.GetTree(&path, nil)
	if err != nil {
		return nil, errors.New("获取失败")
	}

	return list[0].Children, nil
}
