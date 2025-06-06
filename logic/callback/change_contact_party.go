package callback

import (
	"errors"
	"fmt"
	"jdy/enums"
	"jdy/model"
	"log"
	"strconv"
	"strings"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/models"
	"gorm.io/gorm"
)

type PartyHandle struct {
	PartyMessage
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

	var stroe model.Store
	if err := model.DB.First(&stroe, "id = ?", handler.PartyCreate.ID).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}
	}

	if stroe.Id != "" {
		return errors.New(handler.PartyCreate.ID + "部门已存在")
	}

	party, err := handler.getInfo(l, handler.PartyCreate.ID)
	if err != nil {
		return err
	}

	// 名称不含有"店"，则返回
	if !strings.Contains(party.Name, "店") {
		log.Printf("部门名称不含有'店': %v\n", party.Name)
		return nil
	}

	store := model.Store{
		Name:     party.Name,
		ParentId: fmt.Sprint(party.ParentID),
		Order:    party.Order,
	}
	store.Id = fmt.Sprint(party.ID)

	if len(party.DepartmentLeaders) > 0 {
		var superiors []model.Account
		if err := model.DB.
			Where("username IN (?)", party.DepartmentLeaders).
			Where(&model.Account{
				Platform: enums.PlatformTypeWxWork,
			}).
			Preload("Staff").
			Find(&superiors).Error; err != nil {
			return err
		}
		if len(superiors) > 0 {
			for _, superior := range superiors {
				if superior.Staff != nil {
					store.Superiors = append(store.Superiors, *superior.Staff)
				}
			}
		}
	}

	if err := model.DB.Create(&store).Error; err != nil {
		return err
	}

	return nil
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

	var party model.Store
	if err := model.DB.First(&party, "id = ?", msg.PartyDelete.ID).Error; err != nil {
		return err
	}

	if err := model.DB.Delete(&party).Error; err != nil {
		return err
	}

	return nil
}

func (h *PartyHandle) getInfo(l *EventChangeContactEvent, Id string) (*models.Department, error) {
	// 转换部门ID
	id, err := strconv.Atoi(Id)
	if err != nil {
		log.Printf("部门ID转换失败: %v\n", err)
		return nil, err
	}

	// 获取部门信息
	res, err := l.Handle.Wechat.JdyWork.Department.Get(l.Handle.Ctx, id)
	if err != nil || res.ErrCode != 0 {
		log.Printf("获取部门失败: %+v\n", res)
		return nil, err
	}

	return res.Department, nil
}
