package sync

import (
	"errors"
	"fmt"
	"jdy/config"
	"jdy/enums"
	"jdy/model"
	"jdy/utils"
	"log"
	"strconv"
	"strings"

	"github.com/ArtisanCloud/PowerSocialite/v3/src/models"
	kmodels "github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/models"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/work"
	"gorm.io/gorm"
)

type WxWorkLogic struct {
	SyncLogic

	App *work.Work
}

type SyncWxWorkContacts struct {
	*WxWorkLogic

	db *gorm.DB
}

func (l *WxWorkLogic) Contacts() error {
	// 初始化应用实例
	l.App = config.NewWechatService().JdyWork
	// 获取通讯录
	list, err := l.App.Department.SimpleList(l.Ctx, 0)
	if err != nil || (list != nil && list.ErrCode != 0) {
		log.Printf("获取通讯录失败: %+v, %+v", err, list)
		return errors.New("获取通讯录失败")
	}

	// 开始事务
	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 初始化逻辑
		logic := &SyncWxWorkContacts{
			WxWorkLogic: l,
			db:          tx,
		}

		// 获取标签
		tags, err := logic.getTags()
		if err != nil {
			return err
		}

		// 循环部门，同步门店和区域、同步员工
		for _, dept := range list.DepartmentIDs {

			// 获取部门详情
			department, err := l.App.Department.Get(l.Ctx, dept.ID)
			if err != nil || (department != nil && department.ErrCode != 0) {
				log.Printf("获取部门失败: %+v, %+v", err, department)
				return errors.New("获取部门失败")
			}
			res, err := logic.getDepartment(department.Department)
			if err != nil {
				return err
			}
			if res == nil {
				continue
			}
			res.Tags = tags

			// 获取部门成员
			users, err := l.App.User.GetDetailedDepartmentUsers(l.Ctx, department.Department.ID, 0)
			if err != nil || (users != nil && users.ErrCode != 0) {
				log.Printf("获取部门成员失败: %+v, %+v", err, users)
				return errors.New("获取部门成员失败")
			}
			if err := logic.getUsers(department.Department, users.UserList, res); err != nil {
				return err
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

func (l *SyncWxWorkContacts) getTags() (map[enums.Identity]model.StaffTag, error) {
	tags := make(map[enums.Identity]model.StaffTag)

	res, err := l.App.UserTag.List(l.Ctx)
	if err != nil || (res != nil && res.ErrCode != 0) {
		log.Printf("获取标签失败: %+v, %+v", err, res)
		return nil, errors.New("获取标签失败")
	}
	for identity, name := range enums.IdentityMap {
		for _, tag := range res.TagList {
			if tag.TagName == name {
				tags[identity] = model.StaffTag{
					TagId: int64(tag.TagID),
					Name:  name,
				}
			}
		}
		if _, ok := tags[identity]; !ok {
			// 如果tags中不存在identity，则创建一个新标签
			cres, err := l.App.UserTag.Create(l.Ctx, name, 0)
			if err != nil || (cres != nil && cres.ErrCode != 0) {
				return nil, err
			}
			tags[identity] = model.StaffTag{
				TagId: cres.TagID,
				Name:  name,
			}
		}
	}

	for _, tag := range tags {
		var data model.StaffTag
		if err := l.db.Where(model.StaffTag{TagId: tag.TagId}).Assign(tag).FirstOrCreate(&data).Error; err != nil {
			return nil, errors.New("创建或更新标签失败")
		}
	}

	return tags, nil
}

func (l *SyncWxWorkContacts) updateTags(username string, tagid int64) error {
	res, err := l.App.UserTag.TagUsers(l.Ctx, tagid, []string{username})
	if err != nil || (res != nil && res.ErrCode != 0) {
		log.Printf("更新标签失败: %+v, %+v", err, res)
		return errors.New("更新标签失败")
	}

	return nil
}

type getDepartmentRes struct {
	Tags map[enums.Identity]model.StaffTag

	Department enums.Department
	Store      *model.Store
	Region     *model.Region
}

// 分析部门
func (l *SyncWxWorkContacts) getDepartment(dept *kmodels.Department) (*getDepartmentRes, error) {
	switch {
	case strings.HasSuffix(dept.Name, enums.DepartmentStore.String()): // 如果是门店
		{
			data := model.Store{
				IdWx:  fmt.Sprintf("%d", dept.ID),
				Name:  dept.Name,
				Alias: dept.NameEN,
				Order: dept.Order,
			}

			store, err := l.getStore(data)
			if err != nil {
				return nil, err
			}

			res := &getDepartmentRes{
				Department: enums.DepartmentStore,
				Store:      store,
				Region:     nil,
			}

			return res, nil
		}
	case strings.HasSuffix(dept.Name, enums.DepartmentRegion.String()): // 如果是区域
		{
			data := model.Region{
				IdWx:  fmt.Sprintf("%d", dept.ID),
				Name:  dept.Name,
				Alias: dept.NameEN,
				Order: dept.Order,
			}

			region, err := l.getRegion(data)
			if err != nil {
				return nil, err
			}

			res := &getDepartmentRes{
				Department: enums.DepartmentRegion,
				Store:      nil,
				Region:     region,
			}

			return res, nil
		}
	case strings.HasSuffix(dept.Name, enums.DepartmentHeaderquarters.String()): // 如果是总部
		{
			data := model.Store{
				IdWx:  fmt.Sprintf("%d", dept.ID),
				Name:  dept.Name,
				Alias: dept.NameEN,
				Order: dept.Order,
			}
			store, err := l.getStore(data)
			if err != nil {
				return nil, err
			}
			res := &getDepartmentRes{
				Department: enums.DepartmentHeaderquarters,
				Store:      store,
				Region:     nil,
			}
			return res, nil
		}
	}

	return nil, nil
}

// 分析员工
func (l *SyncWxWorkContacts) getUsers(dept *kmodels.Department, users []*models.Employee, res *getDepartmentRes) error {
	// 循环成员
	for _, user := range users {
		data := model.Staff{
			Username:   user.UserID,
			Nickname:   user.Name,
			IsDisabled: true,
			Identity:   enums.IdentityClerk,
		}

		// 判断是否禁用
		if user.Status == 1 {
			data.IsDisabled = false
		}

		if len(user.DirectLeader) > 0 {
			data.LeaderName = user.DirectLeader[0]
		}

		// 获取员工
		staff, logs, err := l.getStaff(data)
		if err != nil {
			return err
		}

		// 判断是否是门店/区域负责人
		_, index, _ := utils.ArrayFind(dept.DepartmentLeaders, func(item string) bool {
			return item == user.UserID
		})

		switch res.Department {
		case enums.DepartmentStore:
			{
				staff.Identity = enums.IdentityClerk
				if err := l.db.Model(&res.Store).Association("Staffs").Append(&staff); err != nil {
					log.Printf("关联员工与门店失败: %+v", err)
					return errors.New("关联员工与门店失败")
				}
				if index != -1 {
					staff.Identity = enums.IdentityShopkeeper
					if err := l.db.Model(&res.Store).Association("Superiors").Append(&staff); err != nil {
						log.Printf("关联负责人与门店失败: %+v", err)
						return errors.New("关联负责人与门店失败")
					}
				}
			}
		case enums.DepartmentRegion:
			{
				staff.Identity = enums.IdentityAreaManager
				if err := l.db.Model(&res.Region).Association("Staffs").Append(&staff); err != nil {
					log.Printf("关联员工与区域失败: %+v", err)
					return errors.New("关联员工与区域失败")
				}
				if index != -1 {
					if err := l.db.Model(&res.Region).Association("Superiors").Append(&staff); err != nil {
						log.Printf("关联负责人与区域失败: %+v", err)
						return errors.New("关联负责人与区域失败")
					}
				}
			}
		case enums.DepartmentHeaderquarters:
			{
				staff.Identity = enums.IdentityHeadquarters
				if err := l.db.Model(&res.Store).Association("Staffs").Append(&staff); err != nil {
					log.Printf("关联员工与总部失败: %+v", err)
					return errors.New("关联员工与总部失败")
				}
				if index != -1 {
					if err := l.db.Model(&res.Store).Association("Superiors").Append(&staff); err != nil {
						log.Printf("关联负责人与门店失败: %+v", err)
						return errors.New("关联负责人与门店失败")
					}
				}
			}
		}

		if err := l.db.Model(&staff).Update("identity", staff.Identity).Error; err != nil {
			log.Printf("更新员工身份失败: %+v", err)
			return errors.New("更新员工身份失败")
		}
		if err := l.updateTags(staff.Username, res.Tags[staff.Identity].TagId); err != nil {
			return err
		}

		logs.NewValue = *staff
		logs.StaffId = staff.Id
		if err := l.db.Create(&logs).Error; err != nil {
			log.Printf("创建员工日志失败: %+v", err)
			return errors.New("创建员工日志失败")
		}
	}

	return nil
}

// 门店
func (l *SyncWxWorkContacts) getStore(where model.Store) (*model.Store, error) {
	// 获取及创建门店
	var store model.Store
	if err := l.db.Where("id_wx = ?", where.IdWx).Assign(where).FirstOrCreate(&store).Error; err != nil {
		return nil, errors.New("获取门店失败")
	}

	return &store, nil
}

// 区域
func (l *SyncWxWorkContacts) getRegion(where model.Region) (*model.Region, error) {
	// 获取及创建区域
	var region model.Region
	if err := l.db.Where("id_wx = ?", where.IdWx).Assign(where).FirstOrCreate(&region).Error; err != nil {
		return nil, errors.New("获取区域失败")
	}

	return &region, nil
}

// 员工
func (l *SyncWxWorkContacts) getStaff(where model.Staff) (*model.Staff, *model.StaffLog, error) {
	var staff model.Staff
	// 获取及创建员工
	if err := l.db.Where(model.Staff{Username: where.Username}).First(&staff).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, nil, errors.New("获取员工失败")
		}
	}

	var logs model.StaffLog
	if staff.Id == "" {
		logs = model.StaffLog{
			Type: enums.StaffLogTypeCreate,
		}
	} else {
		logs = model.StaffLog{
			Type:     enums.StaffLogTypeUpdate,
			OldValue: staff,
		}
	}
	if where.IsDisabled {
		logs.Type = enums.StaffLogTypeDisable
	}

	if err := l.db.Where("username = ?", where.Username).Assign(where).FirstOrCreate(&staff).Error; err != nil {
		return nil, nil, errors.New("获取员工失败")
	}

	logs.NewValue = staff

	return &staff, &logs, nil
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
			if strings.HasSuffix(dept.Name, enums.DepartmentStore.String()) {
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
