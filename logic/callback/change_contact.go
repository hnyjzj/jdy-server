package callback

import (
	"errors"
	"fmt"
	"jdy/enums"
	"jdy/model"
	"log"
	"strconv"
	"strings"

	models1 "github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/models"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/work/server/handlers/models"
	"gorm.io/gorm"
)

// 分配
func (l *EventChangeContactEvent) Distribute() error {
	switch l.Message.ChangeType {
	case models.CALLBACK_EVENT_CHANGE_TYPE_CREATE_USER: // 新增成员
		return l.CreateUser(&l.Message)
	case models.CALLBACK_EVENT_CHANGE_TYPE_UPDATE_USER: // 更新成员
		return l.UpdateUser(&l.Message)
	case models.CALLBACK_EVENT_CHANGE_TYPE_DELETE_USER: // 删除成员
		return l.DeleteUser(&l.Message)
	case models.CALLBACK_EVENT_CHANGE_TYPE_CREATE_PARTY: // 新增部门
		return l.CreateParty(&l.Message)
	case models.CALLBACK_EVENT_CHANGE_TYPE_DELETE_PARTY: // 删除部门
		return l.DeleteParty(&l.Message)
	default:
		err := fmt.Errorf("不支持更改类型(%v)", l.Message.ChangeType)
		log.Printf(err.Error()+": %+v", l.Message)
		return err
	}
}

// 员工变更事件
type EventChangeContactEvent struct {
	Handle  *WxWork                       // 处理器
	Message models1.CallbackMessageHeader // 消息体

	UserCreate models.EventUserCreate // 新增成员
	UserUpdate models.EventUserUpdate // 更新成员
	UserDelete models.EventUserDelete // 删除成员

	PartyCreate models.EventPartyCreate // 新增部门
	PartyDelete models.EventPartyDelete // 删除部门
}

// 创建用户
func (l *EventChangeContactEvent) CreateUser(message *models1.CallbackMessageHeader) error {
	// 解析消息体
	if err := l.Handle.Event.ReadMessage(&l.UserCreate); err != nil {
		return err
	}

	if l.UserCreate.UserID == "" {
		return nil
	}

	var mobile *string
	if l.UserCreate.Mobile == "" {
		mobile = nil
		log.Printf("%v,手机号为空", l.UserCreate.UserID)
	} else {
		mobile = &l.UserCreate.Mobile
	}

	var account model.Account
	if err := model.DB.Where(model.Account{
		Username: &l.UserCreate.UserID,
		Platform: enums.PlatformTypeWxWork,
	}).Attrs(model.Account{
		Phone:    mobile,
		Nickname: &l.UserCreate.Name,
		Avatar:   &l.UserCreate.Avatar,
		Email:    &l.UserCreate.Email,
		Gender:   enums.GenderUnknown.Convert(l.UserCreate.Gender),
		Info:     &l.UserCreate,
	}).FirstOrCreate(&account).Error; err != nil {
		return err
	}

	return nil
}

// 更新用户
func (l *EventChangeContactEvent) UpdateUser(message *models1.CallbackMessageHeader) error {
	// 解析消息体
	if err := l.Handle.Event.ReadMessage(&l.UserUpdate); err != nil {
		return err
	}

	if l.UserUpdate.UserID == "" {
		return nil
	}

	var account model.Account
	if err := model.DB.Where(model.Account{
		Username: &l.UserUpdate.UserID,
		Platform: enums.PlatformTypeWxWork,
	}).First(&account).Error; err != nil {
		return err
	}

	uid := l.UserUpdate.UserID
	if l.UserUpdate.NewUserID != "" {
		uid = l.UserUpdate.NewUserID
	}

	if err := model.DB.Model(&account).Updates(model.Account{
		Username: &uid,
		Info:     &l.UserUpdate,
	}).Error; err != nil {
		return err
	}

	return nil
}

// 删除用户
func (l *EventChangeContactEvent) DeleteUser(message *models1.CallbackMessageHeader) error {
	// 解析消息体
	if err := l.Handle.Event.ReadMessage(&l.UserDelete); err != nil {
		return err
	}

	if l.UserDelete.UserID == "" {
		return nil
	}

	var account model.Account
	if err := model.DB.Where(model.Account{
		Username: &l.UserDelete.UserID,
		Platform: enums.PlatformTypeWxWork,
	}).Preload("Staff").First(&account).Error; err != nil {
		return errors.New(l.UserDelete.UserID + "用户不存在")
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		if account.Staff != nil {
			if err := tx.Delete(&account.Staff).Error; err != nil {
				return errors.New(l.UserDelete.UserID + "删除员工失败")
			}
			if err := tx.Where(model.Account{
				StaffId: account.StaffId,
			}).Delete(&model.Account{}).Error; err != nil {
				return errors.New(l.UserDelete.UserID + "删除账号失败")
			}
		} else {
			if err := tx.Delete(&account).Error; err != nil {
				return errors.New(l.UserDelete.UserID + "删除空账号失败")
			}
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

// 创建部门
func (l *EventChangeContactEvent) CreateParty(message *models1.CallbackMessageHeader) error {
	// 解析消息体
	if err := l.Handle.Event.ReadMessage(&l.PartyCreate); err != nil {
		return err
	}

	if l.PartyCreate.ID == "" {
		return nil
	}

	var party model.Store
	if err := model.DB.First(&party, "id = ?", l.PartyCreate.ID).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}
	}

	if party.Id != "" {
		return errors.New(l.PartyCreate.ID + "部门已存在")
	}

	id, err := strconv.Atoi(l.PartyCreate.ID)
	if err != nil {
		fmt.Printf("转换失败: %v\n", err)
		return err
	}

	res, err := l.Handle.Wechat.JdyWork.Department.Get(l.Handle.Ctx, id)
	if err != nil || res.ErrCode != 0 {
		log.Printf("获取部门失败: %+v\n", res)
		return err
	}

	// 名称不含有"店"，则返回
	if !strings.Contains(res.Department.Name, "店") {
		log.Printf("部门名称不含有'店': %v\n", res.Department.Name)
		return nil
	}

	store := model.Store{
		Name:     res.Department.Name,
		ParentId: fmt.Sprint(res.Department.ParentID),
		Order:    res.Department.Order,
	}
	store.Id = fmt.Sprint(res.Department.ID)

	if len(res.Department.DepartmentLeaders) > 0 {
		var superiors []model.Account
		if err := model.DB.
			Where("username IN (?)", res.Department.DepartmentLeaders).
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
func (l *EventChangeContactEvent) DeleteParty(message *models1.CallbackMessageHeader) error {
	// 解析消息体
	if err := l.Handle.Event.ReadMessage(&l.PartyDelete); err != nil {
		return err
	}

	if l.PartyDelete.ID == "" {
		return nil
	}

	var party model.Store
	if err := model.DB.First(&party, "id = ?", l.PartyDelete.ID).Error; err != nil {
		return err
	}

	if err := model.DB.Delete(&party).Error; err != nil {
		return err
	}

	return nil
}

// 员工变更事件处理
func (Handle *WxWork) ChangeContactEvent() any {
	var (
		l = EventChangeContactEvent{
			Handle: Handle,
		}
	)

	// 解析消息体
	if err := l.Handle.Event.ReadMessage(&l.Message); err != nil {
		log.Printf("TemplateCardEvent.ReadMessage.Error(): %v\n", err.Error())
		return "error"
	}

	// 处理事件
	if err := l.Distribute(); err != nil {
		log.Printf("TemplateCardEvent.GetEventKey.Error(): %v\n", err.Error())
		return "error"
	}

	return nil
}
