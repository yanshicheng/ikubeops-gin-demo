package apps

import (
	"fmt"
	"github.com/gin-gonic/gin"
	routerservice "github.com/yanshicheng/ikubeops-gin-demo/common/router-service"
	"github.com/yanshicheng/ikubeops-gin-demo/router"
)

// IOC 容器层 管理所有服务实力对象

var (
	// 维护当前所有服务
	logicApps = map[string]LogicService{}
	// gin
	ginApps = map[string]routerservice.GinService{}
)

// Get 一个Impl服务的实例：logicApps
// 返回一个对象, 任何类型都可以, 使用时, 由使用方进行断言
func GetLogic(name string) interface{} {
	for k, v := range logicApps {
		if k == name {
			return v
		}
	}

	return nil
}
func GetGinApp(name string) routerservice.GinService {
	app, ok := ginApps[name]
	if !ok {
		panic(fmt.Sprintf("handler app %s not registed", name))
	}

	return app
}

// 通过断言自动注册
func RegistryLogic(svc LogicService) {
	// 服务实例注册到svcs map当中
	if _, ok := logicApps[svc.Name()]; ok {
		panic(fmt.Sprintf("service %s has registried", svc.Name()))
	}

	logicApps[svc.Name()] = svc
}

// LoadedGinApp 查询加载成功的服务
func LoadedGinApp() (apps []string) {
	for k := range ginApps {
		apps = append(apps, k)
	}
	return
}

type LogicService interface {
	Config()
	Name() string
}

// 用于初始化注册到 IOC容器中的所有服务
func InitImpl() {
	for _, v := range logicApps {
		v.Config()
	}
}

// 通过断言自动注册
func RegistryGin(svc routerservice.GinService) {
	// 服务实例注册到svcs map当中
	if _, ok := ginApps[svc.Name()]; ok {
		panic(fmt.Sprintf("gin service %s has registried", svc.Name()))
	}

	ginApps[svc.Name()] = svc
}

// 用于初始化注册到 IOC容器中的所有服务
func InitGin() *gin.Engine {
	for _, v := range ginApps {
		v.Config()
	}
	// 自动注册路由
	return router.Registry(ginApps)
}
