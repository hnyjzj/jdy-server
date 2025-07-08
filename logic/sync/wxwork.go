package sync

import (
	"errors"
	"fmt"
	"jdy/config"
	"jdy/enums"
	"jdy/model"
	"log"
	"strconv"
	"strings"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/work"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WxWorkLogic struct {
	Ctx *gin.Context

	App *work.Work
}

type SyncWxWorkContacts struct {
	*WxWorkLogic

	db *gorm.DB

	staff  model.Staff
	store  model.Store
	region model.Region
}

func (l *WxWorkLogic) Contacts() error {
	// 初始化应用实例
	l.App = config.NewWechatService().JdyWork
	// 获取通讯录
	id := 0
	list, err := l.App.Department.SimpleList(l.Ctx, id)
	if err != nil || (list != nil && list.ErrCode != 0) {
		log.Printf("获取通讯录失败: %+v, %+v", err, list)
		return errors.New("获取通讯录失败")
	}

	// 开始事务
	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 循环部门，同步门店和区域、同步员工
		for _, dept := range list.DepartmentIDs {
			// 初始化逻辑
			logic := &SyncWxWorkContacts{
				WxWorkLogic: l,
				db:          tx,
			}
			// 获取部门详情
			department, err := l.App.Department.Get(l.Ctx, dept.ID)
			if err != nil || (department != nil && department.ErrCode != 0) {
				log.Printf("获取部门失败: %+v, %+v", err, department)
				return errors.New("获取部门失败")
			}

			switch {
			case strings.Contains(department.Department.Name, "店"): // 如果是门店
				store := model.Store{
					IdWx:  fmt.Sprintf("%d", department.Department.ID),
					Name:  department.Department.Name,
					Order: department.Department.Order,
				}
				if err := logic.getStore(store); err != nil {
					return err
				}
			case department.Department.Name == "总部": // 如果是总部
				store := model.Store{
					IdWx:  fmt.Sprintf("%d", department.Department.ID),
					Name:  department.Department.Name,
					Order: department.Department.Order,
				}
				if err := logic.getStore(store); err != nil {
					return err
				}
			case strings.Contains(department.Department.Name, "区域"): // 如果是区域
				region := model.Region{
					IdWx:  fmt.Sprintf("%d", department.Department.ID),
					Name:  department.Department.Name,
					Order: department.Department.Order,
				}
				if err := logic.getRegion(region); err != nil {
					return err
				}
			}

			// 获取部门成员
			users, err := l.App.User.GetDetailedDepartmentUsers(l.Ctx, department.Department.ID, 0)
			if err != nil || (users != nil && users.ErrCode != 0) {
				log.Printf("获取部门成员失败: %+v, %+v", err, users)
				return errors.New("获取部门成员失败")
			}

			// 循环成员
			for _, user := range users.UserList {
				// 获取员工
				var gender enums.Gender
				if err := logic.getStaff(model.Staff{
					Username:   user.UserID,
					Phone:      user.Mobile,
					Nickname:   user.Name,
					Avatar:     user.Avatar,
					Email:      user.Email,
					Gender:     gender.Convert(user.Gender),
					IsDisabled: true,
					Identity:   enums.IdentityClerk,
				}); err != nil {
					return err
				}

				// 关联区域
				if err := logic.appendRegion(); err != nil {
					return err
				}
				// 关联门店
				if err := logic.appendStore(); err != nil {
					return err
				}
				// 设置负责人
				if err := logic.appendSuperior(department.Department.Name, department.Department.DepartmentLeaders); err != nil {
					return err
				}

				// 判断是否禁用
				if user.Status == 1 {
					logic.staff.IsDisabled = false
				}

				// 保存员工
				if err := logic.saveStaff(); err != nil {
					return err
				}
			}
		}

		return nil
	}); err != nil {
		return err
	}

	// 同步区域下的门店
	if err := l.syncRegionStore(); err != nil {
		return err
	}

	return nil
}

// 同步区域下的门店
func (l *WxWorkLogic) syncRegionStore() error {
	var regions []model.Region
	if err := model.DB.Find(&regions).Error; err != nil {
		log.Printf("获取区域失败: %+v", err)
		return errors.New("获取区域失败")
	}

	for _, region := range regions {
		if region.IdWx == "" {
			continue
		}

		id, err := strconv.Atoi(region.IdWx)
		if err != nil {
			log.Printf("获取区域失败: %+v", err)
			return errors.New("获取区域失败")
		}

		list, err := l.App.Department.List(l.Ctx, id)
		if err != nil || (list != nil && list.ErrCode != 0) {
			log.Printf("获取通讯录失败: %+v, %+v", err, list)
			return errors.New("获取通讯录失败")
		}

		var ids []string
		for _, dept := range list.Departments {
			if strings.Contains(dept.Name, "店") {
				ids = append(ids, fmt.Sprintf("%d", dept.ID))
			}
		}

		var stores []model.Store
		if err := model.DB.Where("id_wx in (?)", ids).Find(&stores).Error; err != nil {
			log.Printf("获取门店失败: %+v", err)
			return errors.New("获取门店失败")
		}

		if err := model.DB.Model(&region).Association("Stores").Replace(&stores); err != nil {
			log.Printf("关联员工与门店失败: %+v", err)
			return errors.New("关联员工与门店失败")
		}
	}

	return nil
}

// 门店（获取/创建）
func (s *SyncWxWorkContacts) getStore(where model.Store) error {
	var store model.Store
	// 获取及创建门店
	if err := s.db.Where("id_wx = ?", where.IdWx).Attrs(model.Store{
		IdWx:  where.IdWx,
		Name:  where.Name,
		Order: where.Order,
	}).FirstOrCreate(&store).Error; err != nil {
		log.Printf("获取门店失败: %+v", err)
		return errors.New("获取门店失败")
	}

	s.store = store

	return nil
}

// 区域（获取/创建）
func (s *SyncWxWorkContacts) getRegion(where model.Region) error {
	var region model.Region
	// 获取及创建区域
	if err := s.db.Where("id_wx = ?", where.IdWx).Attrs(model.Region{
		IdWx:  where.IdWx,
		Name:  where.Name,
		Order: where.Order,
	}).FirstOrCreate(&region).Error; err != nil {
		log.Printf("获取区域失败: %+v", err)
		return errors.New("获取区域失败")
	}

	s.region = region

	return nil
}

// 员工（获取/创建/更新）
func (s *SyncWxWorkContacts) getStaff(where model.Staff) error {
	var staff model.Staff
	// 获取及创建员工
	if err := s.db.Where(model.Staff{
		Username: where.Username,
	}).Attrs(where).FirstOrCreate(&staff).Error; err != nil {
		log.Printf("获取员工失败: %+v", err)
		return errors.New("获取员工失败")
	}

	// 更新员工信息
	if err := s.db.Model(&staff).Updates(where).Error; err != nil {
		log.Printf("更新员工信息失败: %+v", err)
		return errors.New("更新员工信息失败")
	}

	s.staff = staff

	return nil
}

// 关联门店
func (s *SyncWxWorkContacts) appendStore() error {
	// 关联员工与门店/区域
	if s.store.IdWx != "" {
		if err := s.db.Model(&s.staff).Association("Stores").Append(&s.store); err != nil {
			log.Printf("关联员工与门店失败: %+v", err)
			return errors.New("关联员工与门店失败")
		}
		s.staff.Identity = enums.IdentityClerk
	}

	return nil
}

// 关联区域
func (s *SyncWxWorkContacts) appendRegion() error {
	// 关联员工与门店/区域
	if s.region.IdWx != "" {
		if err := s.db.Model(&s.staff).Association("Regions").Append(&s.region); err != nil {
			log.Printf("关联员工与区域失败: %+v", err)
			return errors.New("关联员工与区域失败")
		}
		s.staff.Identity = enums.IdentityAreaManager
	}

	return nil
}

// 设置负责人
func (s *SyncWxWorkContacts) appendSuperior(name string, ids []string) error {
	// 设置员工为门店/区域负责人
	if strings.Contains(strings.Join(ids, " "), s.staff.Username) {
		switch {
		case strings.Contains(name, "店"): // 如果是门店
			if err := s.db.Model(&s.store).Association("Superiors").Append(&s.staff); err != nil {
				log.Printf("关联负责人与门店失败: %+v", err)
				return errors.New("关联负责人与门店失败")
			}
			s.staff.Identity = enums.IdentityShopkeeper
		case strings.Contains(name, "区域"): // 如果是区域
			if err := s.db.Model(&s.region).Association("Superiors").Append(&s.staff); err != nil {
				log.Printf("关联负责人与区域失败: %+v", err)
				return errors.New("关联负责人与区域失败")
			}
			s.staff.Identity = enums.IdentityAreaManager
		}
	}

	return nil
}

// 保存员工
func (s *SyncWxWorkContacts) saveStaff() error {
	if err := s.db.Save(&s.staff).Error; err != nil {
		log.Printf("更新员工状态失败: %+v", err)
		return errors.New("更新员工状态失败")
	}

	return nil
}
