package services

import (
	"github.com/yanshicheng/ikubeops-gin-demo/common/entities"
	"github.com/yanshicheng/ikubeops-gin-demo/common/pagination"
	"github.com/yanshicheng/ikubeops-gin-demo/template/apps/demo/models"
)

type Service interface {
	// 创建主机
	CreateBook(*models.Book) error
	QueryList(*pagination.Pagination) (*string, error)
	DeleteBook(*entities.Person) (*models.Book, error)
	// 更新主机
	//Update(context.Context, *UpdateHostRequest) (*Hosts, error)
	// 删除主机
	//Delete(context.Context, *DeleteHostRequest) (*Hosts, error)
	// 主机详情
	//Describe(context.Context, *QueryHostRequest) (*Hosts, error)
}
