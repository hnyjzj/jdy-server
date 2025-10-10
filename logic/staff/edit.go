package staff

import (
	"jdy/config"
	"jdy/enums"
	"jdy/errors"
	"jdy/model"
	"jdy/types"
	"jdy/utils"
	"log"
	"slices"
	"strconv"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/work/user/request"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 修改员工
func (l *StaffLogic) StaffEdit(req *types.StaffEditReq) error {
	logic := &StaffEditLogic{
		Ctx:      l.Ctx,
		Req:      req,
		Operator: l.Staff,
	}
	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		logic.Db = tx

		// 查询员工是否存在
		if err := logic.getStaff(); err != nil {
			return err
		}

		// 查询上级
		if err := logic.getLeader(); err != nil {
			return err
		}

		// 查询标签
		if err := logic.getTag(); err != nil {
			return err
		}

		// 构建员工信息
		if err := logic.buildStaff(); err != nil {
			return err
		}

		// 查询门店
		if err := logic.getStore(); err != nil {
			return err
		}

		// 查询区域
		if err := logic.getRegion(); err != nil {
			return err
		}

		// 更新员工信息
		if err := logic.updateStaff(); err != nil {
			return err
		}

		// 更新企业微信
		if err := logic.updateWechat(); err != nil {
			return err
		}

		// 添加记录
		if err := logic.addlogs(); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

type StaffEditLogic struct {
	Ctx      *gin.Context        // 上下文
	Req      *types.StaffEditReq // 请求参数
	Db       *gorm.DB            // 数据库连接
	Operator *model.Staff        // 操作人

	Staff  *model.Staff    // 员工
	Data   *model.Staff    // 员工更新数据
	Leader *model.Staff    // 上级
	Tags   *model.StaffTag // 标签

	Stores  map[string]model.Store  // 门店
	Regions map[string]model.Region // 区域
}

// 查询员工是否存在
func (l *StaffEditLogic) getStaff() error {
	var (
		staff model.Staff
	)

	// 查询员工
	db := l.Db.Where("id = ?", l.Req.Id).Preload("Tag")
	if err := db.First(&staff).Error; err != nil {
		return errors.New("查询员工失败")
	}
	l.Staff = &staff

	return nil
}

// 查询上级
func (l *StaffEditLogic) getLeader() error {
	var (
		leader model.Staff
	)

	if l.Req.LeaderName == "" {
		return nil
	}

	if err := l.Db.Unscoped().Where(&model.Staff{
		Username: l.Req.LeaderName,
	}).First(&leader).Error; err != nil {
		return errors.New("查询上级失败")
	}

	l.Leader = &leader
	return nil
}

// 查询标签
func (l *StaffEditLogic) getTag() error {
	var (
		tags model.StaffTag
		app  = config.NewWechatService().JdyWork
		name = l.Req.Identity.String()
	)

	if err := l.Db.Where(&model.StaffTag{
		Name: name,
	}).First(&tags).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return errors.New("查询标签失败")
		}
	}

	if tags.Id == "" {
		cres, err := app.UserTag.Create(l.Ctx, name, 0)
		if err != nil || (cres != nil && cres.ErrCode != 0) {
			return err
		}
		tags = model.StaffTag{
			TagId: cres.TagID,
			Name:  name,
		}
		if err := l.Db.Create(&tags).Error; err != nil {
			return errors.New("创建标签失败")
		}
	}

	l.Tags = &tags
	return nil
}

// 构建员工信息
func (l *StaffEditLogic) buildStaff() error {
	// 构建员工信息
	data := model.Staff{
		Phone:      l.Req.Phone,
		Username:   l.Req.Username,
		Nickname:   l.Req.Nickname,
		Avatar:     l.Req.Avatar,
		Email:      l.Req.Email,
		Gender:     l.Req.Gender,
		IsDisabled: l.Req.IsDisabled,
		Identity:   l.Req.Identity,
		RoleId:     l.Req.RoleId,
		TagId:      l.Tags.Id,
	}

	if l.Leader != nil {
		data.LeaderName = l.Leader.Username
	}

	if l.Req.Password != "" {
		password, err := data.HashPassword(&l.Req.Password)
		if err != nil {
			return errors.New("密码加密失败")
		}
		data.Password = password
	}

	l.Data = &data

	return nil
}

// 查询门店
func (l *StaffEditLogic) getStore() error {
	allids := make([]string, 0)
	allids = append(allids, l.Req.StoreIds...)
	allids = append(allids, l.Req.StoreSuperiorIds...)
	allids = append(allids, l.Req.StoreAdminIds...)

	ids := utils.ArrayUnique(allids, func(item string) string { return item })
	if len(ids) == 0 {
		return nil
	}

	var stores []model.Store
	if err := l.Db.Where("id in (?)", ids).Find(&stores).Error; err != nil {
		return errors.New("查询门店失败")
	}

	l.Stores = make(map[string]model.Store)
	for _, store := range stores {
		l.Stores[store.Id] = store
	}

	return nil
}

// 查询区域
func (l *StaffEditLogic) getRegion() error {
	allids := make([]string, 0)
	allids = append(allids, l.Req.RegionIds...)
	allids = append(allids, l.Req.RegionSuperiorIds...)
	allids = append(allids, l.Req.RegionAdminIds...)

	ids := utils.ArrayUnique(allids, func(item string) string { return item })
	if len(ids) == 0 {
		return nil
	}

	var regions []model.Region
	if err := l.Db.Where("id in (?)", ids).Find(&regions).Error; err != nil {
		return errors.New("查询区域失败")
	}

	l.Regions = make(map[string]model.Region)
	for _, region := range regions {
		l.Regions[region.Id] = region
	}

	return nil
}

// 更新员工信息
func (l *StaffEditLogic) updateStaff() error {
	// 更新员工信息
	if err := l.Db.Model(&model.Staff{}).Where("id = ?", l.Staff.Id).Updates(l.Data).Error; err != nil {
		return err
	}

	// 关联门店
	if len(l.Req.StoreIds) == 0 {
		if err := l.Db.Model(&l.Staff).Association("Stores").Clear(); err != nil {
			return err
		}
	} else {
		var stores []model.Store
		if err := l.Db.Where("id in (?)", l.Req.StoreIds).Find(&stores).Error; err != nil {
			return err
		}
		if err := l.Db.Model(&l.Staff).Association("Stores").Replace(stores); err != nil {
			return err
		}
	}

	// 关联负责的门店
	if len(l.Req.StoreSuperiorIds) == 0 {
		if err := l.Db.Model(&l.Staff).Association("StoreSuperiors").Clear(); err != nil {
			return err
		}
	} else {
		var store_superiors []model.Store
		if err := l.Db.Where("id in (?)", l.Req.StoreSuperiorIds).Find(&store_superiors).Error; err != nil {
			return err
		}
		if err := l.Db.Model(&l.Staff).Association("StoreSuperiors").Replace(store_superiors); err != nil {
			return err
		}
	}

	// 关联管理的门店
	if len(l.Req.StoreAdminIds) == 0 {
		if err := l.Db.Model(&l.Staff).Association("StoreAdmins").Clear(); err != nil {
			return err
		}
	} else {
		var store_admins []model.Store
		if err := l.Db.Where("id in (?)", l.Req.StoreAdminIds).Find(&store_admins).Error; err != nil {
			return err
		}
		if err := l.Db.Model(&l.Staff).Association("StoreAdmins").Replace(store_admins); err != nil {
			return err
		}
	}

	// 关联区域
	if len(l.Req.RegionIds) == 0 {
		if err := l.Db.Model(&l.Staff).Association("Regions").Clear(); err != nil {
			return err
		}
	} else {
		var regions []model.Region
		if err := l.Db.Where("id in (?)", l.Req.RegionIds).Find(&regions).Error; err != nil {
			return err
		}
		if err := l.Db.Model(&l.Staff).Association("Regions").Replace(regions); err != nil {
			return err
		}
	}

	// 关联负责区域
	if len(l.Req.RegionSuperiorIds) == 0 {
		if err := l.Db.Model(&l.Staff).Association("RegionSuperiors").Clear(); err != nil {
			return err
		}
	} else {
		var region_superiors []model.Region
		if err := l.Db.Where("id in (?)", l.Req.RegionSuperiorIds).Find(&region_superiors).Error; err != nil {
			return err
		}
		if err := l.Db.Model(&l.Staff).Association("RegionSuperiors").Replace(region_superiors); err != nil {
			return err
		}
	}

	// 关联管理的区域
	if len(l.Req.RegionAdminIds) == 0 {
		if err := l.Db.Model(&l.Staff).Association("RegionAdmins").Clear(); err != nil {
			return err
		}
	} else {
		var region_admins []model.Region
		if err := l.Db.Where("id in (?)", l.Req.RegionAdminIds).Find(&region_admins).Error; err != nil {
			return err
		}
		if err := l.Db.Model(&l.Staff).Association("RegionAdmins").Replace(region_admins); err != nil {
			return err
		}
	}

	// 查询更新完的员工信息
	if err := l.Db.First(&l.Data, "id = ?", l.Staff.Id).Error; err != nil {
		return err
	}

	return nil
}

// 更新企业微信
func (l *StaffEditLogic) updateWechat() error {
	wxwork := config.NewWechatService().ContactsWork
	jdy := config.NewWechatService().JdyWork

	user := &request.RequestUserDetail{
		Userid:           l.Data.Username,
		Name:             l.Data.Nickname,
		Mobile:           l.Data.Phone,
		Department:       make([]int, 0), // 所属部门
		Position:         l.Data.Identity.String(),
		Gender:           uint32(l.Data.Gender),
		Email:            l.Data.Email,
		IsLeaderInDept:   make([]int, 0),    // 部门负责人
		DirectLeader:     make([]string, 0), // 直属上级
		Enable:           1,
		ToInvite:         true,
		ExternalPosition: l.Data.Identity.StringExternal(),
	}

	if l.Req.IsDisabled {
		user.Enable = 0
	}

	for _, rid := range l.Req.StoreIds {
		store, ok := l.Stores[rid]
		if !ok {
			return errors.New("所属门店不存在")
		}
		wid, err := strconv.Atoi(store.IdWx)
		if err != nil {
			return errors.New("门店ID转换失败")
		}
		isLeaderInDept := 0
		if slices.Contains(l.Req.StoreSuperiorIds, rid) {
			isLeaderInDept = 1
		}
		user.Department = append(user.Department, wid)
		user.IsLeaderInDept = append(user.IsLeaderInDept, isLeaderInDept)

	}
	for _, id := range l.Req.RegionIds {
		region, ok := l.Regions[id]
		if !ok {
			return errors.New("所属区域不存在")
		}
		wid, err := strconv.Atoi(region.IdWx)
		if err != nil {
			return errors.New("区域ID转换失败")
		}

		isLeaderInDept := 0
		if slices.Contains(l.Req.RegionSuperiorIds, id) {
			isLeaderInDept = 1
		}
		user.Department = append(user.Department, wid)
		user.IsLeaderInDept = append(user.IsLeaderInDept, isLeaderInDept)
	}

	if l.Leader != nil {
		user.DirectLeader = append(user.DirectLeader, l.Leader.Username)
	}

	cres, err := wxwork.User.Update(l.Ctx, user)
	if err != nil || (cres != nil && cres.ErrCode != 0) {
		log.Printf("企业微信更新员工失败: %+v, %+v", err, cres)
		return errors.New("企业微信更新员工失败")
	}

	dres, err := jdy.UserTag.TagDelUsers(l.Ctx, l.Staff.Tag.TagId, []string{l.Staff.Username}, nil)
	if err != nil || (dres != nil && dres.ErrCode != 0) {
		log.Printf("删除标签失败: %+v, %+v", err, dres)
		return errors.New("删除标签失败")
	}

	tres, err := jdy.UserTag.TagUsers(l.Ctx, l.Tags.TagId, []string{user.Userid})
	if err != nil || (tres != nil && tres.ErrCode != 0) {
		log.Printf("更新标签失败: %+v, %+v", err, tres)
		return errors.New("更新标签失败")
	}

	return nil
}

// 添加记录
func (l *StaffEditLogic) addlogs() error {
	logs := model.StaffLog{
		Type:       enums.StaffLogTypeUpdate,
		StaffId:    l.Staff.Id,
		OldValue:   *l.Staff,
		NewValue:   *l.Data,
		OperatorId: l.Operator.Id,
	}
	if err := l.Db.Create(&logs).Error; err != nil {
		return errors.New("添加记录失败")
	}

	return nil
}
