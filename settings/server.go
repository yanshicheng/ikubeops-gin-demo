package settings

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yanshicheng/ikubeops-gin-demo/global"
	"github.com/yanshicheng/ikubeops-gin-demo/version"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type server interface {
	ListenAndServe() error
}

func RunIkubeopsServer(router *gin.Engine) {
	address := fmt.Sprintf("%s:%d", global.C.App.HttpAddr, global.C.App.HttpPort)
	s := initServer(address, router)

	//defer func() {
	//
	//}()

	time.Sleep(60 * time.Microsecond)
	global.L.Info("server run success on ", zap.String("address", address))

	fmt.Printf(`
欢迎使用: %s
当前版本: %s
配置文件: 
演示地址: www.ikubeops.com
代码地址: %s
运行地址: %s

`, version.IkubeopsProjectName, version.ShortTagVersion(), version.IkubeopsUrl, fmt.Sprintf("http://%s", address))
	if err := s.ListenAndServe(); err != nil {
		fmt.Println(err)
		//L.Error("server run error", zap.Error(err))
	}
}

func initServer(address string, router *gin.Engine) server {
	s := endless.NewServer(address, router)
	s.ReadHeaderTimeout = 10 * time.Second
	s.WriteTimeout = 60 * time.Second
	s.MaxHeaderBytes = 1 << 20

	return s
}

// IkubeopsServerManager 结构体管理HTTP服务器
type IkubeopsServerManager struct {
	server *http.Server
}

// NewIkubeopsServerManager 创建一个新的IkubeopsServerManager实例
func NewIkubeopsServerManager(router *gin.Engine) *IkubeopsServerManager {
	server := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", global.C.App.HttpAddr, global.C.App.HttpPort),
		Handler:           router,
		ReadTimeout:       time.Duration(global.C.App.ReadTimeout) * time.Second,
		WriteTimeout:      time.Duration(global.C.App.WriteTimeout) * time.Second,
		ReadHeaderTimeout: time.Duration(global.C.App.ReadHeaderTimeout) * time.Second,
		MaxHeaderBytes:    global.C.App.MaxHeaderSize * 1024 * 1024, // 1 MB
	}
	return &IkubeopsServerManager{server: server}
}

// Run 启动服务器
func (sm *IkubeopsServerManager) Run() {
	if global.C.App.Tls {
		go func() {
			err := sm.server.ListenAndServeTLS(global.C.App.CertFile, global.C.App.KeyFile)
			if err != nil {
				global.LSys.Error("TLS Server : %s", err.Error())
				return
			}
		}()
	} else {
		go func() {
			if err := sm.server.ListenAndServe(); err != nil {
				global.LSys.Error("Server : %s", err.Error())
				return
			}
		}()
	}
	fmt.Printf(`
欢迎使用: %s
当前版本: %s
配置文件: %s
演示地址: www.ikubeops.com
代码地址: %s
运行地址: %s

`, version.IkubeopsProjectName, version.ShortTagVersion(), version.GetConfig(), version.IkubeopsUrl, version.GetWebUrl())
}

// GracefulExit 优雅退出服务器
func (sm *IkubeopsServerManager) GracefulExit() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	sig := <-ch

	global.LSys.Info("接收到退出信号: ", sig)
	now := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(global.C.App.ShutdownTimeout)*time.Second)
	defer cancel()
	err := sm.server.Shutdown(ctx)
	if err != nil {
		global.LSys.Error("server shutdown error", zap.Error(err))
	}
	// 看看实际退出所耗费的时间
	global.LSys.Info("退出耗时: %s", time.Since(now))
}
