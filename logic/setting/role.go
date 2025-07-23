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
	var (
		db = model.DB.Model(&model.Role{})
	)

	if req.IsDefault {
		var role model.Role
		if err := db.Where(&model.Role{
			Identity:  req.Identity,
			IsDefault: req.IsDefault,
		}).First(&role).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				return nil, errors.New("查询角色失败")
			}
		}
		if role.Id != "" {
			return nil, errors.New("每个身份只能有一个默认角色")
		}
	}

	role := model.Role{
		Name:      req.Name,
		Desc:      req.Desc,
		Identity:  req.Identity,
		IsDefault: req.IsDefault,

		OperatorId: r.Staff.Id,
		IP:         r.IP,
	}

	if err := db.Create(&role).Error; err != nil {
		return nil, errors.New("创建角色失败,名称可能有历史重复")
	}

	return &role, nil
}

func (r *RoleLogic) Copy(req *types.RoleCopyReq) (*model.Role, error) {
	var (
		role model.Role
		db   = model.DB
	)

	rdb := db.Model(&model.Role{})
	rdb = role.Preloads(rdb)
	if err := rdb.First(&role, "id = ?", req.Id).Error; err != nil {
		log.Printf("查询角色失败: %v", err)
		return nil, errors.New("查询角色失败")
	}

	data := model.Role{
		Name:       req.Name,
		Desc:       req.Desc,
		Identity:   role.Identity,
		IsDefault:  false,
		OperatorId: r.Staff.Id,
		IP:         r.IP,

		Apis:    role.Apis,
		Routers: role.Routers,
	}

	if err := db.Create(&data).Error; err != nil {
		return nil, errors.New("复制角色失败")
	}

	return &data, nil
}

func (r *RoleLogic) List(req *types.RoleListReq) ([]model.Role, error) {
	var (
		roles []model.Role
	)

	db := model.DB.Model(&model.Role{})
	db = model.Role{}.WhereCondition(db, &types.RoleWhere{Identity: req.Identity})

	db = db.Order("created_at desc").Order("is_default desc")
	if err := db.Find(&roles).Error; err != nil {
		return nil, err
	}

	return roles, nil
}

func (r *RoleLogic) Info(req *types.RoleInfoReq) (*model.Role, error) {
	var (
		role model.Role
	)
	db := model.DB.Model(&model.Role{})

	db = role.Preloads(db)

	if err := db.First(&role, "id = ?", req.Id).Error; err != nil {
		log.Printf("查询角色失败: %v", err)
		return nil, errors.New("查询角色失败")
	}

	return &role, nil
}

func (r *RoleLogic) Edit(req *types.RoleEditReq) error {
	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 查询角色信息
		var role model.Role
		if err := tx.First(&role, "id = ?", req.Id).Error; err != nil {
			log.Printf("查询角色失败: %v", err)
			return errors.New("查询角色失败")
		}
		// 如果是默认角色，则查询是否有其他默认角色，如果有则将其设置为非默认角色
		if *req.IsDefault {
			var def model.Role
			if err := tx.Where(&model.Role{
				Identity:  role.Identity,
				IsDefault: true,
			}).Where("id != ?", req.Id).First(&def).Error; err != nil {
				if err != gorm.ErrRecordNotFound {
					return errors.New("查询默认角色失败")
				}
			}
			if def.Id != "" {
				def.IsDefault = false
				if err := tx.Save(&def).Error; err != nil {
					return errors.New("更新默认角色失败")
				}
			}
		}

		// 更新角色信息
		data := model.Role{
			Name:      req.Name,
			Desc:      req.Desc,
			Identity:  req.Identity,
			IsDefault: *req.IsDefault,

			OperatorId: r.Staff.Id,
			IP:         r.IP,
		}

		if err := tx.Model(&model.Role{}).Where("id = ?", role.Id).Updates(data).Error; err != nil {
			return errors.New("更新角色失败")
		}

		if req.IsDefault != nil {
			role.IsDefault = *req.IsDefault
			if err := tx.Save(&role).Error; err != nil {
				return errors.New("更新角色默认失败")
			}
		}

		return nil
	}); err != nil {
		return errors.New(err.Error())
	}
	return nil
}

func (r *RoleLogic) Update(req *types.RoleUpdateReq) error {
	var (
		role model.Role
	)
	if err := model.DB.First(&role, "id = ?", req.Id).Error; err != nil {
		log.Printf("查询角色失败: %v", err)
		return errors.New("查询角色失败")
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 更新路由
		var routers []model.Router
		if err := tx.Model(&model.Router{}).Where("id in (?)", req.Routers).Find(&routers).Error; err != nil {
			return err
		}
		if err := tx.Model(&role).Association("Routers").Replace(routers); err != nil {
			return err
		}
		// 更新接口
		var Apis []model.Api
		if err := tx.Model(&model.Api{}).Where("id in (?)", req.Apis).Find(&Apis).Error; err != nil {
			return err
		}
		if err := tx.Model(&role).Association("Apis").Replace(Apis); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return errors.New("更新角色失败")
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

	var staff_count int64
	if err := model.DB.Model(&model.Staff{}).Where("role_id = ?", req.Id).Count(&staff_count).Error; err != nil {
		return errors.New("查询员工失败")
	}
	if staff_count > 0 {
		return errors.New("该角色下有员工，无法删除")
	}

	if role.IsDefault {
		return errors.New("默认角色无法删除")
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

	if len(list) == 0 {
		return []any{}, nil
	}

	return list[0].Children, nil
}
