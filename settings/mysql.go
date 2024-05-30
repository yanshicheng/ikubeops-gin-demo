package settings

import (
	"fmt"
	"github.com/yanshicheng/ikubeops-gin-demo/global"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func LoadMysql() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s", global.C.Mysql.User, global.C.Mysql.Password, global.C.Mysql.Host, global.C.Mysql.Port, global.C.Mysql.DbName, global.C.Mysql.Config)
	mysqlConfig := mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}

	// SetMaxOpenConns 设置连接池最大打开连接数
	if db, err := gorm.Open(mysql.New(mysqlConfig), gormConfig()); err != nil {
		global.L.Error("MySQL启动异常", zap.Any("err", err))
		return nil, err
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(global.C.Mysql.MaxIdleConns)
		sqlDB.SetMaxOpenConns(global.C.Mysql.MaxOpenConns)
		return db, nil
	}
}

// @author: SliverHorn
// @function: gormConfig
// @description: 根据配置决定是否开启日志
// @param: mod bool
// @return: *gorm.Config

func gormConfig() *gorm.Config {
	config := &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true}
	switch global.C.Mysql.LogMode {
	case "silent", "Silent":
		config.Logger = Default.LogMode(logger.Silent)
	case "error", "Error":
		config.Logger = Default.LogMode(logger.Error)
	case "warn", "Warn":
		config.Logger = Default.LogMode(logger.Warn)
	case "info", "Info":
		config.Logger = Default.LogMode(logger.Info)
	default:
		config.Logger = Default.LogMode(logger.Info)
	}
	return config
}

func GetDb() *gorm.DB {
	return global.DB
}
func RegisterModels(models ...interface{}) {
	global.M = append(global.M, models...)
}
