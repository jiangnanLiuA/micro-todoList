package dao

import (
	"context"
	"fmt"
	"micro_todoList/config"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var _db *gorm.DB

func InitDB() {
	host := config.DbHost
	port := config.DbPort
	database := config.DbName
	username := config.DbUser
	password := config.DbPassWord
	charset := config.Charset
	//dsn := strings.Join([]string{username, ":", password, "@tcp(", host, ":", port, ")/", database, "?charset=" + charset + "&parseTime=true"}, "")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true", username, password, host, port, database, charset)
	err := Database(dsn)
	if err != nil {
		fmt.Println(err)
	}
}

// 数据库连接

func Database(connString string) error {
	var ormLogger logger.Interface //打印日志在控制台
	if gin.Mode() == "debug" {
		ormLogger = logger.Default.LogMode(logger.Info)
	} else {
		ormLogger = logger.Default
	}
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       connString, // 连接数据库 dsn -> data source name
		DefaultStringSize:         256,        // string 类型字段的默认长度
		DisableDatetimePrecision:  true,       // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,       // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,       // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,      // 根据版本自动配置
	}), &gorm.Config{
		Logger: ormLogger, //打印日志
		NamingStrategy: schema.NamingStrategy{ //命名策略
			SingularTable: true, //表名不加s   默认会加s ->   user -> users
		},
	})
	if err != nil {
		panic(err)
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(20)  //设置连接池
	sqlDB.SetMaxOpenConns(100) //最大打开数
	sqlDB.SetConnMaxLifetime(time.Second * 30)
	_db = db
	migration()
	return err
}

// 让 context 贯穿整条链路

func NewDBClient(ctx context.Context) *gorm.DB {
	db := _db
	return db.WithContext(ctx)
}
