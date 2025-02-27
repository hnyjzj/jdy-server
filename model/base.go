package model

import (
	"time"

	"github.com/bwmarrin/snowflake"
	"gorm.io/gorm"
)

// 基础模型
type BaseModel struct {
	Id string `json:"id" gorm:"uniqueIndex;size:255;primaryKey;AUTO_INCREMENT:false;comment:ID"`
}

// 通用模型
type Model struct {
	BaseModel
	CreatedAt *time.Time `json:"created_at" gorm:"type:datetime;not null;autoCreateTime;index;comment:创建时间"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"type:datetime;not null;autoUpdateTime;index;comment:更新时间"`
}

// 软删除模型
type SoftDelete struct {
	Model
	DeletedAt gorm.DeletedAt `json:"-" sql:"index" gorm:"type:datetime;comment:删除时间"`
}

// 创建雪花节点
var Node *snowflake.Node

// 初始化
func init() {
	var err error
	// 创建一个雪花节点
	Node, err = snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}
}

// BeforeCreate GORM回调，在创建记录之前自动设置雪花ID
func (pm *BaseModel) BeforeCreate(tx *gorm.DB) error {
	if pm.Id == "" {
		pm.Id = Node.Generate().String()
	}

	return nil
}
