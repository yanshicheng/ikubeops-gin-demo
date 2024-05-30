package logic

import (
	"github.com/yanshicheng/ikubeops-gin-demo/apps"
	"github.com/yanshicheng/ikubeops-gin-demo/global"
	logger "github.com/yanshicheng/ikubeops-gin-demo/settings/ikube-logger"
	"github.com/yanshicheng/ikubeops-gin-demo/template/apps/demo"
	"github.com/yanshicheng/ikubeops-gin-demo/template/apps/demo/services"
	"gorm.io/gorm"
)

// 类型检查 Service 是不是等于 LogicService
var _ services.Service = (*Service)(nil)

var impl = &Service{}

type Service struct {
	services.Service
	L  *logger.Logger
	db *gorm.DB
}

// 只需要保证 全局对象Config和全局Logger已经加载完成
func (i *Service) Config() {
	// Host service 服务的子Loggger
	// 封装的Zap让其满足 Logger接口
	// 为
	//什么要封装:
	// 		1. Logger全局实例
	// 		2. Logger Level的动态调整, Logrus不支持Level共同调整
	// 		3. 加入日志轮转功能的集合
	i.L = global.L.Named(demo.AppName)
	i.db = global.DB
}

func (i *Service) Name() string {
	return demo.AppName
}

func init() {
	// 注册
	apps.RegistryLogic(impl)
}

func NewLogicService() *Service {
	return &Service{
		L:  global.L.Named(demo.AppName),
		db: global.DB,
	}
}
