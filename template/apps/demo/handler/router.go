package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/yanshicheng/ikubeops-gin-demo/apps"
	"github.com/yanshicheng/ikubeops-gin-demo/template/apps/demo"
	"github.com/yanshicheng/ikubeops-gin-demo/template/apps/demo/logic"
)

// 通过一个实体类，把内部的接口通过 HTTP 协议暴露出去
// 依赖内部接口实现

var handler = &Handler{}

type Handler struct {
	svc *logic.Service
}

// Config 配置函数，在这里注入依赖，并且初始化实例，供其他函数使用。
func (h *Handler) Config() {
	h.svc = logic.NewLogicService()
}

func (h *Handler) PublicRegistry(r gin.IRouter) {

}
func (h *Handler) AuthRegistry(r gin.IRouter) {
	// 分组路由
	group := r.Group("v1/book")
	{
		// group.GET("/list", h.List)
		group.POST("/add", h.bookCreate)
		group.GET("/get", h.bookGet)
		group.DELETE("/delete/:id", h.bookDelete)
	}
}
func (h *Handler) Name() string {
	return demo.AppName
}

func init() {
	apps.RegistryGin(handler)
}
