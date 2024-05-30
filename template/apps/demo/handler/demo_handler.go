package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yanshicheng/ikubeops-gin-demo/common/entities"
	"github.com/yanshicheng/ikubeops-gin-demo/common/response"
	"github.com/yanshicheng/ikubeops-gin-demo/template/apps/demo/models"
	"time"
)

func (h *Handler) bookCreate(c *gin.Context) {
	var ins *models.Book
	// 解析用户传递的参数
	if err := c.ShouldBindJSON(&ins); err != nil {
		fmt.Println(err.Error())
		response.FailedParam(c, err)
		return
	}
	// 解析成功后调用 logic 层创建主机
	if err := h.svc.CreateBook(ins); err != nil {
		response.Failed(c, err.Error())
		return
	}
	response.Success(c, ins)
	time.Sleep(20 * time.Second)
	return
}

func (h *Handler) bookGet(context *gin.Context) {
	time.Sleep(20 * time.Second)
	response.Success(context, "ok")
	return
}

func (h *Handler) bookDelete(c *gin.Context) {
	var person entities.Person
	pid := c.ShouldBindUri(&person)
	if pid != nil {
		response.Failed(c, "id is null")
		return
	}
	if ins, err := h.svc.DeleteBook(&person); err != nil {
		response.Failed(c, err.Error())
		return
	} else {
		response.Success(c, ins)
		return
	}
}
