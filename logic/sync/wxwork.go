package sync

import (
	"errors"
	"fmt"
	"jdy/config"
	"jdy/enums"
	"jdy/model"
	"log"
	"strings"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/work"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WxWorkLogic struct {
	Ctx *gin.Context

	App *work.Work
}

func (l *WxWorkLogic) Contacts() error {
	// 初始化应用实例
	l.App = config.NewWechatService().JdyWork
	// 获取通讯录
	id := 0
	list, err := l.App.Department.SimpleList(l.Ctx, id)
	if err != nil || list.ErrCode != 0 {
		log.Printf("获取通讯录失败: %+v, %+v", err, list)
		return errors.New("获取通讯录失败")
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 循环部门
		for _, dept := range list.DepartmentIDs {
			// 获取部门详情
			department, err := l.App.Department.Get(l.Ctx, dept.ID)
			if err != nil || department.ErrCode != 0 {
				log.Printf("获取部门失败: %+v, %+v", err, department)
				return errors.New("获取部门失败")
			}

			var (
				store  model.Store
				region model.Region
			)

			switch {
			case strings.Contains(department.Department.Name, "店"): // 如果是门店
				// 获取及创建门店
				if err := tx.Where("id_wx = ?", department.Department.ID).Attrs(model.Store{
					IdWx:  fmt.Sprintf("%d", department.Department.ID),
					Name:  department.Department.Name,
					Order: department.Department.Order,
				}).FirstOrCreate(&store).Error; err != nil {
					log.Printf("获取门店失败: %+v", err)
					return errors.New("获取门店失败")
				}
			case strings.Contains(department.Department.Name, "区域"): // 如果是区域
				// 获取及创建区域
				if err := tx.Where("id_wx = ?", department.Department.ID).Attrs(model.Region{
					IdWx:  fmt.Sprintf("%d", department.Department.ID),
					Name:  department.Department.Name,
					Order: department.Department.Order,
				}).FirstOrCreate(&region).Error; err != nil {
					log.Printf("获取区域失败: %+v", err)
					return errors.New("获取区域失败")
				}
			}

			// 获取部门成员
			users, err := l.App.User.GetDetailedDepartmentUsers(l.Ctx, department.Department.ID, 0)
			if err != nil || users.ErrCode != 0 {
				log.Printf("获取部门成员失败: %+v, %+v", err, users)
				return errors.New("获取部门成员失败")
			}

			// 循环成员
			for _, user := range users.UserList {
				// 获取及创建员工
				var staff model.Staff
				var gender enums.Gender
				data := model.Staff{
					Phone:      user.Mobile,
					Nickname:   user.Name,
					Avatar:     user.Avatar,
					Email:      user.Email,
					Gender:     gender.Convert(user.Gender),
					IsDisabled: true,
				}
				if err := tx.Where(&model.Staff{Username: user.UserID}).FirstOrCreate(&staff).Error; err != nil {
					log.Printf("获取员工失败: %+v", err)
					return errors.New("获取员工失败")
				}

				// 更新员工信息
				if err := tx.Model(&staff).Updates(data).Error; err != nil {
					log.Printf("更新员工失败: %+v", err)
					return errors.New("更新员工失败")
				}

				// 更新员工状态
				if user.Status == 1 {
					staff.IsDisabled = false
				}
				if err := tx.Save(&staff).Error; err != nil {
					log.Printf("更新员工状态失败: %+v", err)
					return errors.New("更新员工状态失败")
				}

				// 关联员工与门店/区域
				if region.IdWx != "" {
					if err := tx.Model(&staff).Association("Regions").Append(&region); err != nil {
						log.Printf("关联员工与区域失败: %+v", err)
						return errors.New("关联员工与区域失败")
					}
				}
				if store.IdWx != "" {
					if err := tx.Model(&staff).Association("Stores").Append(&store); err != nil {
						log.Printf("关联员工与门店失败: %+v", err)
						return errors.New("关联员工与门店失败")
					}
				}

				// 设置员工为门店/区域负责人
				if strings.Contains(strings.Join(department.Department.DepartmentLeaders, " "), user.UserID) {
					switch {
					case strings.Contains(department.Department.Name, "店"): // 如果是门店
						if err := tx.Model(&store).Association("Superiors").Append(&staff); err != nil {
							log.Printf("关联负责人与门店失败: %+v", err)
							return errors.New("关联负责人与门店失败")
						}
					case strings.Contains(department.Department.Name, "区域"): // 如果是区域
						if err := tx.Model(&region).Association("Superiors").Append(&staff); err != nil {
							log.Printf("关联负责人与区域失败: %+v", err)
							return errors.New("关联负责人与区域失败")
						}
					}
				}
			}
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
