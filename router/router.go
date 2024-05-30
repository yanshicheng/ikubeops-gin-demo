package router

import (
	"github.com/gin-gonic/gin"
	"github.com/yanshicheng/ikubeops-gin-demo/common/response"
	routerservice "github.com/yanshicheng/ikubeops-gin-demo/common/router-service"
	"github.com/yanshicheng/ikubeops-gin-demo/global"
	"github.com/yanshicheng/ikubeops-gin-demo/middleware"
	"github.com/yanshicheng/ikubeops-gin-demo/settings"
)

func routers() *gin.Engine {
	router := gin.New()
	router.Use(middleware.GinRecovery(true), settings.GinLogger(), middleware.Cors())
	// 处理跨域请求
	//Router.Use(middleware.Cors())
	// 获取路由组实例
	return router
}

func Registry(ginApps map[string]routerservice.GinService) *gin.Engine {
	r := routers()
	RouterPublicGroup := r.Group("")
	{
		// 健康监测
		RouterPublicGroup.GET("/healthz", healthStatusRoutes)
	}

	RouterAuthGroup := r.Group("api")
	RouterAuthGroup.Use(middleware.DemoMiddleware())
	{
	}

	for _, v := range ginApps {
		v.PublicRegistry(RouterPublicGroup)
		v.AuthRegistry(RouterAuthGroup)
	}
	return r
}

func healthStatusRoutes(c *gin.Context) {
	// 检查数据库是否正常
	if global.C.Mysql.Enable {
		db, err := settings.GetDb().DB()
		if err != nil {
			global.L.Named("healthz").Errorf("获取数据库异常: %s", err.Error())
			response.FailServerErr(c, "获取数据库异常")
			return
		}
		if db.Ping() != nil {
			global.L.Named("healthz").Errorf("数据库异常: %v", err)
			response.FailServerErr(c, "数据库异常")
			return
		}
	}
	response.Success(c, "ok")
}
