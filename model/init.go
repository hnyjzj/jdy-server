package model

import (
	"fmt"
	"jdy/config"
	"reflect"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB  *gorm.DB
	err error
)

// 初始化模型
func Init() {
	var (
		drive = config.Config.Database.Drive
		dsn   = config.Config.Database.Dsn()
	)

	switch drive {
	case "mysql":
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Error),
		})
	default:
		panic("不支持的数据库驱动程序")
	}

	if err != nil {
		panic("数据库初始化失败")
	}

	// 执行迁移
	migrator()
}

// 需要迁移的表
var Tables []interface{} = []interface{}{}

// 注册数据表
func RegisterModels(models ...interface{}) {
	Tables = append(Tables, models...)
}

// 需要重建的表
var Refresh []interface{} = []interface{}{}

// 注册需要重新迁移的数据表
func RegisterRefreshModels(models ...interface{}) {
	Refresh = append(Refresh, models...)
}

// 迁移数据表
func migrator() {
	for _, table := range Tables {
		// 重建模型
		if config.Config.Database.Refresh {
			for _, againTable := range Refresh {
				if reflect.DeepEqual(table, againTable) {
					fmt.Println("删除了", reflect.TypeOf(table).Elem().Name())
					DB.Migrator().DropTable(table)
					break
				}
			}
		}
		// 迁移模型
		err := DB.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci").Migrator().AutoMigrate(table)
		if err != nil {
			panic(fmt.Sprintf("迁移表 %s 失败: %s", strings.Split(reflect.TypeOf(table).Elem().Name(), ""), err.Error()))
		}
		// fmt.Println("迁移了", reflect.TypeOf(table).Elem().Name())
	}
}
