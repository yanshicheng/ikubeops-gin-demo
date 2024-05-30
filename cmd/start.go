package cmd

import (
	"fmt"
	"github.com/yanshicheng/ikubeops-gin-demo/apps"
	_ "github.com/yanshicheng/ikubeops-gin-demo/apps/all"
	"github.com/yanshicheng/ikubeops-gin-demo/global"
	"github.com/yanshicheng/ikubeops-gin-demo/settings"
	"github.com/yanshicheng/ikubeops-gin-demo/version"
	//"github.com/yanshicheng/ikubeops-gin-demo/global"
	"github.com/spf13/cobra"
)

// 注册所有服务

// startCmd represents the start command
var serviceCmd = &cobra.Command{
	Use:   "start",
	Short: fmt.Sprintf("%s API服务", version.IkubeopsProjectName),
	Long:  fmt.Sprintf("%s API服务", version.IkubeopsProjectName),
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		// 初始化全局变量
		if err := settings.LoadGlobalConfig(confType, confFile); err != nil {
			fmt.Printf("init validator failed, err:%v\n", err)
			return err
		}

		// 初始化全局日志配置
		if global.L, err = settings.LoadGlobalLogger(); err != nil {
			fmt.Printf("init logger failed, err:%v\n", err)
			return err
		}
		global.LSys = global.L.Named("system")
		//if global.L, err = settings.LoadLogger(); err != nil {
		//	return err
		//}
		if err := settings.InitTrans(); err != nil {
			global.LSys.Error("init InitTrans failed, err:%v\n", err)
			return err
		}

		// 初始化 mysql
		if global.C.Mysql.Enable {
			if global.DB, err = settings.LoadMysql(); err != nil {
				global.LSys.Error("初始化 mysql 出错: %s", err)
				return err
			}
		}
		// 初始化全局app
		apps.InitImpl()

		// 加载路由
		//router := router.RouterGrouputers()
		// 注册路由
		router := apps.InitGin()
		server := settings.NewIkubeopsServerManager(router)
		// 启动服务
		server.Run()
		// 优雅退出
		server.GracefulExit()
		return nil
	},
}

// config 为全局变量, 只需要load 即可全局可用户

// func newService(cnf *conf.Config) (*service, error) {
// 	handler := protocol.NewHTTPService()
// 	grpc := protocol.NewGRPCService()
// 	svr := &service{
// 		handler: handler,
// 		grpc: grpc,
// 		log:  zap.L().Named("CLI"),
// 	}

// 	return svr, nil
// }

// type service struct {
// 	handler *protocol.HTTPService
// 	grpc *protocol.GRPCService

// 	log logger.Logger
// }

// func (s *service) start() error {
// 	s.log.Infof("loaded grpc app: %s", app.LoadedGrpcApp())
// 	s.log.Infof("loaded handler app: %s", app.LoadedRESTfulApp())
// 	s.log.Infof("loaded internal app: %s", app.LoadedInternalApp())

// 	go s.grpc.Start()
// 	return s.handler.Start()
// }

// config 为全局变量, 只需要load 即可全局可用户
//func loadGlobalConfig(configType, configFile string) error {
//	// 配置加载
//	switch configType {
//	case "file":
//		err := global.LoadConfig(global.ConfigModeFile, configFile)
//		if err != nil {
//			return err
//		}
//	case "env":
//		err := global.LoadConfig(global.ConfigModeEnv, configFile)
//		if err != nil {
//			return err
//		}
//	default:
//		return errors.New("unknown config type")
//	}
//
//	return nil
//}

// // log 为全局变量, 只需要load 即可全局可用户, 依赖全局配置先初始化
// func loadGlobalLogger() error {
// 	var (
// 		logInitMsg string
// 		level      zap.Level
// 	)
// 	lc := conf.C().Log
// 	lv, err := zap.NewLevel(lc.Level)
// 	if err != nil {
// 		logInitMsg = fmt.Sprintf("%s, use default level INFO", err)
// 		level = zap.InfoLevel
// 	} else {
// 		level = lv
// 		logInitMsg = fmt.Sprintf("log level: %s", lv)
// 	}
// 	zapConfig := zap.DefaultConfig()
// 	zapConfig.Level = level
// 	switch lc.To {
// 	case conf.ToStdout:
// 		zapConfig.ToStderr = true
// 		zapConfig.ToFiles = false
// 	case conf.ToFile:
// 		zapConfig.Files.Name = "api.log"
// 		zapConfig.Files.Path = lc.PathDir
// 	}
// 	switch lc.Format {
// 	case conf.JSONFormat:
// 		zapConfig.JSON = true
// 	}
// 	if err := zap.Configure(zapConfig); err != nil {
// 		return err
// 	}
// 	zap.L().Named("INIT").Info(logInitMsg)
// 	return nil
// }

// func (s *service) waitSign(sign chan os.Signal) {
// 	for sg := range sign {
// 		switch v := sg.(type) {
// 		default:
// 			s.log.Infof("receive signal '%v', start graceful shutdown", v.String())

// 			if err := s.grpc.Stop(); err != nil {
// 				s.log.Errorf("grpc graceful shutdown err: %s, force exit", err)
// 			} else {
// 				s.log.Info("grpc service stop complete")
// 			}

// 			if err := s.handler.Stop(); err != nil {
// 				s.log.Errorf("handler graceful shutdown err: %s, force exit", err)
// 			} else {
// 				s.log.Infof("handler service stop complete")
// 			}
// 			return
// 		}
// 	}
// }

func init() {
	rootCommand.AddCommand(serviceCmd)
}
