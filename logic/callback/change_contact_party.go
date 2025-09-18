package callback

import (
	"errors"
	"fmt"
	"jdy/model"
	"log"
	"strconv"
	"strings"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/models"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/work/department/response"
	"gorm.io/gorm"
)

type PartyHandle struct {
	PartyMessage
}

type PartyCreateHandle struct {
	PartyHandle

	Party *models.Department
}

// 创建部门
func (l *EventChangeContactEvent) CreateParty() error {
	// 解析消息体
	var handler PartyHandle
	if err := l.Handle.Event.ReadMessage(&handler.PartyCreate); err != nil {
		return err
	}

	if handler.PartyCreate.ID == "" {
		return nil
	}

	party, err := handler.getInfo(l, handler.PartyCreate.ID)
	if err != nil {
		return err
	}

	// 名称不含有"店"，则返回
	create_handle := PartyCreateHandle{
		PartyHandle: handler,
		Party:       party.Department,
	}

	switch {
	case strings.HasSuffix(party.Department.Name, model.StorePrefix):
		return create_handle.isStore(l)
	case strings.HasSuffix(party.Department.Name, model.RegionPrefix):
		return create_handle.isRegion(l)
	case strings.HasSuffix(party.Department.Name, model.HeaderquartersPrefix):
		return create_handle.isStore(l)
	default:
		return nil
	}
}

// 删除部门
func (l *EventChangeContactEvent) DeleteParty() error {
	// 解析消息体
	var msg PartyHandle
	if err := l.Handle.Event.ReadMessage(&msg.PartyDelete); err != nil {
		return err
	}

	if msg.PartyDelete.ID == "" {
		return nil
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where(&model.Store{
			IdWx: msg.PartyDelete.ID,
		}).Delete(&model.Store{}).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				return err
			}
		}
		if err := tx.Where(&model.Region{
			IdWx: msg.PartyDelete.ID,
		}).Delete(&model.Region{}).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				return err
			}
		}

		return nil
	}); err != nil {
		fmt.Printf("删除门店失败：%s\n", err.Error())
		return err
	}
	return nil
}

// 门店
func (h *PartyCreateHandle) isStore(l *EventChangeContactEvent) error {
	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 查询是否存在
		var have model.Store
		if err := tx.Where(&model.Store{
			IdWx: h.PartyCreate.ID,
		}).First(&have).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				return err
			}
		}
		if have.Id != "" {
			return errors.New(h.PartyCreate.ID + "部门已存在")
		}

		// 创建部门
		store := model.Store{
			IdWx:  fmt.Sprint(h.Party.ID),
			Alias: h.Party.NameEN,
			Name:  h.Party.Name,
			Order: h.Party.Order,
		}

		if len(h.Party.DepartmentLeaders) > 0 {
			var superiors []model.Staff
			if err := tx.
				Where("username IN (?)", h.Party.DepartmentLeaders).
				Find(&superiors).Error; err != nil {
				return err
			}
			if len(superiors) > 0 {
				store.Superiors = append(store.Superiors, superiors...)
				store.Admins = append(store.Admins, superiors...)
			}
		}

		if err := tx.Create(&store).Error; err != nil {
			return err
		}

		// 获取父级部门
		if h.Party.ParentID != 0 {
			// 获取父级部门信息
			parent, err := h.getInfo(l, h.Party.ParentID)
			if err != nil {
				return err
			}
			// 判断父级部门是否为区域
			if strings.HasSuffix(parent.Department.Name, model.RegionPrefix) {
				// 获取父级部门ID
				var region model.Region
				if err := tx.Where(&model.Region{
					IdWx: fmt.Sprint(parent.Department.ID),
				}).First(&region).Error; err != nil {
					if err != gorm.ErrRecordNotFound {
						return err
					}
				}
				// 添加门店到区域
				if region.Id != "" {
					region.Stores = append(region.Stores, store)
					if err := tx.Save(&region).Error; err != nil {
						return err
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

// 区域
func (h *PartyCreateHandle) isRegion(l *EventChangeContactEvent) error {
	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 查询是否存在
		var have model.Region
		if err := tx.Where(&model.Region{
			IdWx: h.PartyCreate.ID,
		}).First(&have).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				return err
			}
		}
		if have.Id != "" {
			return errors.New(h.PartyCreate.ID + "区域已存在")
		}

		simpleList, err := h.getSimpleList(l, h.PartyCreate.ID)
		if err != nil {
			return err
		}

		// 创建区域
		region := model.Region{
			IdWx:  fmt.Sprint(h.Party.ID),
			Name:  h.Party.Name,
			Alias: h.Party.NameEN,
			Order: h.Party.Order,
		}

		// 获取子部门
		var simpleListDepartmentIDs []string
		for _, department := range simpleList.DepartmentIDs {
			simpleListDepartmentIDs = append(simpleListDepartmentIDs, fmt.Sprint(department.ID))
		}
		var stores []model.Store
		if err := tx.Where("id_wx IN (?)", simpleListDepartmentIDs).Find(&stores).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				return err
			}
		}
		region.Stores = append(region.Stores, stores...)

		if err := tx.Create(&region).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

// 获取部门信息
func (h *PartyHandle) getInfo(l *EventChangeContactEvent, Id any) (*response.ResponseDepartmentGet, error) {
	// 转换部门ID
	id, err := strconv.Atoi(fmt.Sprint(Id))
	if err != nil {
		log.Printf("部门ID转换失败: %v\n", err)
		return nil, err
	}

	// 获取部门信息
	res, err := l.Handle.Wechat.JdyWork.Department.Get(l.Handle.Ctx, id)
	if err != nil || (res != nil && res.ErrCode != 0) {
		log.Printf("获取部门失败: %+v\n", res)
		if err == nil {
			err = fmt.Errorf("wechat api error: %d %s", res.ErrCode, res.ErrMsg)
		}
		return nil, err
	}

	return res, nil
}

// 获取子部门列表
func (h *PartyHandle) getSimpleList(l *EventChangeContactEvent, Id string) (*response.ResponseDepartmentIDList, error) {
	// 转换部门ID
	id, err := strconv.Atoi(fmt.Sprint(Id))
	if err != nil {
		log.Printf("部门ID转换失败: %v\n", err)
		return nil, err
	}

	// 获取部门信息
	res, err := l.Handle.Wechat.JdyWork.Department.SimpleList(l.Handle.Ctx, id)
	if err != nil || (res != nil && res.ErrCode != 0) {
		log.Printf("获取部门失败: %+v\n", res)
		if err == nil {
			err = fmt.Errorf("wechat api error: %d %s", res.ErrCode, res.ErrMsg)
		}
		return nil, err
	}

	return res, nil
}
