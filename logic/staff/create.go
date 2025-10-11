package staff

import (
	"jdy/config"
	"jdy/enums"
	"jdy/errors"
	"jdy/message"
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

// 创建员工
func (l *StaffLogic) StaffCreate(ctx *gin.Context, req *types.StaffCreateReq) error {
	logic := &StaffCreateLogic{
		Ctx:      ctx,
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

		// 创建账号
		if err := logic.createStaff(); err != nil {
			return err
		}

		// 创建企业微信
		if err := logic.createWechat(); err != nil {
			return err
		}

		// 添加记录
		if err := logic.addlogs(); err != nil {
			return err
		}

		// 发送消息
		if err := logic.sendMessage(); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

type StaffCreateLogic struct {
	Ctx      *gin.Context
	Req      *types.StaffCreateReq
	Db       *gorm.DB
	Operator *model.Staff

	Staff  *model.Staff
	Leader *model.Staff
	Tags   *model.StaffTag

	Stores  map[string]model.Store
	Regions map[string]model.Region
}

// 查询员工是否存在
func (l *StaffCreateLogic) getStaff() error {
	var (
		staff model.Staff
	)

	// 根据手机号或用户名查询账号
	db := l.Db.Unscoped()
	db = db.Where(&model.Staff{
		Phone: l.Req.Phone,
	})
	db = db.Or(&model.Staff{
		Username: l.Req.Username,
	})
	if err := db.First(&staff).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return errors.New("查询账号失败")
		}
	}

	if staff.Id != "" {
		return errors.New("账号已存在")
	}

	return nil
}

// 查询上级
func (l *StaffCreateLogic) getLeader() error {
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
func (l *StaffCreateLogic) getTag() error {
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
func (l *StaffCreateLogic) buildStaff() error {
	// 创建账号
	l.Staff = &model.Staff{
		Username:   l.Req.Username,
		Phone:      l.Req.Phone,
		Nickname:   l.Req.Nickname,
		Avatar:     l.Req.Avatar,
		Email:      l.Req.Email,
		Gender:     l.Req.Gender,
		Identity:   l.Req.Identity,
		LeaderName: l.Req.LeaderName,
		TagId:      l.Tags.Id,
		RoleId:     l.Req.RoleId,
	}

	password, err := l.Staff.HashPassword(&l.Req.Password)
	if err != nil {
		return errors.New("密码加密失败")
	}
	l.Staff.Password = password

	return nil
}

// 查询门店
func (l *StaffCreateLogic) getStore() error {
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
func (l *StaffCreateLogic) getRegion() error {
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

// 创建账号
func (l *StaffCreateLogic) createStaff() error {
	if err := l.Db.Create(&l.Staff).Error; err != nil {
		return errors.New("创建账号失败")
	}

	for _, id := range l.Req.StoreIds {
		store, ok := l.Stores[id]
		if !ok {
			return errors.New("所属门店不存在")
		}
		if err := l.Db.Model(&l.Staff).Association("Stores").Append(&store); err != nil {
			return errors.New("关联所属门店失败")
		}
	}
	for _, id := range l.Req.StoreSuperiorIds {
		store, ok := l.Stores[id]
		if !ok {
			return errors.New("负责门店不存在")
		}
		if err := l.Db.Model(&l.Staff).Association("StoreSuperiors").Append(&store); err != nil {
			return errors.New("关联负责门店失败")
		}
	}
	for _, id := range l.Req.StoreAdminIds {
		store, ok := l.Stores[id]
		if !ok {
			return errors.New("管理门店不存在")
		}
		if err := l.Db.Model(&l.Staff).Association("StoreAdmins").Append(&store); err != nil {
			return errors.New("关联管理门店失败")
		}
	}

	for _, id := range l.Req.RegionIds {
		region, ok := l.Regions[id]
		if !ok {
			return errors.New("所属区域不存在")
		}
		if err := l.Db.Model(&l.Staff).Association("Regions").Append(&region); err != nil {
			return errors.New("关联所属区域失败")
		}
	}
	for _, id := range l.Req.RegionSuperiorIds {
		region, ok := l.Regions[id]
		if !ok {
			return errors.New("负责区域不存在")
		}
		if err := l.Db.Model(&l.Staff).Association("RegionSuperiors").Append(&region); err != nil {
			return errors.New("关联负责区域失败")
		}
	}
	for _, id := range l.Req.RegionAdminIds {
		region, ok := l.Regions[id]
		if !ok {
			return errors.New("管理区域不存在")
		}
		if err := l.Db.Model(&l.Staff).Association("RegionAdmins").Append(&region); err != nil {
			return errors.New("关联管理区域失败")
		}
	}

	return nil
}

// 企业微信创建员工
func (l *StaffCreateLogic) createWechat() error {
	wxwork := config.NewWechatService().ContactsWork
	jdy := config.NewWechatService().JdyWork

	user := &request.RequestUserDetail{
		Userid:           l.Staff.Username,
		Name:             l.Staff.Nickname,
		Mobile:           l.Staff.Phone,
		Department:       make([]int, 0), // 所属部门
		Position:         l.Staff.Identity.String(),
		Gender:           uint32(l.Staff.Gender),
		Email:            l.Staff.Email,
		IsLeaderInDept:   make([]int, 0),    // 部门负责人
		DirectLeader:     make([]string, 0), // 直属上级
		Enable:           1,
		ToInvite:         true,
		ExternalPosition: l.Staff.Identity.StringExternal(),
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

	cres, err := wxwork.User.Create(l.Ctx, user)
	if err != nil || (cres != nil && cres.ErrCode != 0) {
		log.Printf("企业微信创建员工失败: %+v, %+v", err, cres)
		switch cres.ErrCode {
		case 60104:
			return errors.New("手机号已存在")
		default:
			return errors.New("企业微信创建员工失败")
		}
	}

	tres, err := jdy.UserTag.TagUsers(l.Ctx, l.Tags.TagId, []string{user.Userid})
	if err != nil || (tres != nil && tres.ErrCode != 0) {
		log.Printf("更新标签失败: %+v, %+v", err, tres)
		return errors.New("更新标签失败")
	}

	return nil
}

// 添加记录
func (l *StaffCreateLogic) addlogs() error {
	logs := model.StaffLog{
		Type:       enums.StaffLogTypeCreate,
		StaffId:    l.Staff.Id,
		NewValue:   *l.Staff,
		OperatorId: l.Operator.Id,
	}
	if err := l.Db.Create(&logs).Error; err != nil {
		return errors.New("添加记录失败")
	}

	return nil
}

// 发送消息
func (l *StaffCreateLogic) sendMessage() error {

	go func() {
		m := message.NewMessage(l.Ctx)
		m.SendRegisterMessage(&message.RegisterMessageContent{
			Username: l.Req.Username,
			Nickname: l.Req.Nickname,
			Phone:    l.Req.Phone,
			Password: l.Req.Password,
		})
	}()

	return nil
}
