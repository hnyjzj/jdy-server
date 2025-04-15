package model

import (
	"fmt"
	"jdy/config"
	"log"
	"os"
	"reflect"
	"strings"
	"time"

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
			Logger: logger.New(
				// 保存到指定位置的日志文件
				log.New(os.Stdout, "\r\n", log.LstdFlags),
				logger.Config{
					SlowThreshold:             time.Second, // Slow SQL threshold
					LogLevel:                  logger.Warn, // Log level
					IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
					ParameterizedQueries:      false,       // Don't include params in the SQL log
					Colorful:                  true,        // Disable color
				},
			),
			DisableForeignKeyConstraintWhenMigrating: true, // 禁用外键约束
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
var Tables []any = []any{}

// 注册数据表
func RegisterModels(models ...any) {
	Tables = append(Tables, models...)
}

// 需要重建的表
var Refresh []any = []any{}

// 注册需要重新迁移的数据表
func RegisterRefreshModels(models ...any) {
	Refresh = append(Refresh, models...)
}

// 迁移数据表
func migrator() {
	for _, table := range Tables {
		// 重建模型
		if config.Config.Database.Refresh {
			for _, againTable := range Refresh {
				if reflect.DeepEqual(table, againTable) {
					log.Println("删除了", reflect.TypeOf(table).Elem().Name())
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
		// log.Println("迁移了", reflect.TypeOf(table).Elem().Name())

		// 如果是门店表，则插入一条默认数据
		if reflect.TypeOf(table).Elem().Name() == "Store" {
			var store Store
			if err := DB.Where(&Store{Name: "公司总店"}).FirstOrCreate(&store).Error; err != nil {
				panic(fmt.Sprintf("插入默认门店失败: %s", err.Error()))
			}
		}
	}
}
