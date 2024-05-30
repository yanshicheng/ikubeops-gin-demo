package cmd

import (
	"database/sql"
	"fmt"
	_ "github.com/yanshicheng/ikubeops-gin-demo/apps/all"
	"github.com/yanshicheng/ikubeops-gin-demo/global"
	"github.com/yanshicheng/ikubeops-gin-demo/settings"
	"github.com/spf13/cobra"
)

var (
	db      string
	migrate bool
)

var dbCommand = &cobra.Command{
	Use:   "db",
	Short: "db console",
	Long:  "db console",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		// 初始化全局变量
		if err := settings.LoadGlobalConfig(confType, confFile); err != nil {
			return err
		}

		// 初始化全局日志配置
		if global.L, err = settings.LoadGlobalLogger(); err != nil {
			fmt.Printf("init logger failed, err:%v\n", err)
			return err
		}
		global.LSys = global.L.Named("system")
		// 初始化 mysql
		//fmt.Println(global.M)
		if !global.C.Mysql.Enable {
			global.LSys.Error(fmt.Sprintf("数据迁移失败，未启用mysql配置，如需迁移数据库，请先启用mysql配置"))
			return nil
		}
		if global.DB, err = settings.LoadMysql(); err != nil {
			global.LSys.Error(fmt.Sprintf("初始化 mysql 出错: %s", err))
			return nil
		}
		global.LSys.Info("开始数据迁移...")
		if global.DB != nil && len(global.M) > 0 {
			db, _ := global.DB.DB()
			if err := global.DB.AutoMigrate(global.M...); err != nil {
				global.LSys.Error("数据库迁移失败: %s\n", err.Error())
			} else {
				global.LSys.Info("数据库迁移成功")
				return nil
			}
			defer func(db *sql.DB) {
				err := db.Close()
				if err != nil {
					global.LSys.Error("数据库关闭失败: ", err)
				}
			}(db)
		} else {
			global.LSys.Info("未检测到 model")
		}
		return nil
	},
}

func init() {
	rootCommand.AddCommand(dbCommand)
	dbCommand.Flags().StringVarP(&db, "database", "d", "default", "database")
	dbCommand.Flags().BoolVarP(&migrate, "migrate", "m", false, "force syncdb")
}
