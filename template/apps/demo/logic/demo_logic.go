package logic

import (
	"fmt"
	"github.com/yanshicheng/ikubeops-gin-demo/common/entities"
	"github.com/yanshicheng/ikubeops-gin-demo/common/pagination"
	"github.com/yanshicheng/ikubeops-gin-demo/template/apps/demo/models"
)

func (i *Service) CreateBook(book *models.Book) error {
	// 调用 db 创建数据
	err := i.db.Create(book).Error
	return err
}

func (i *Service) QueryList(*pagination.Pagination) (*string, error) {
	fmt.Println("CreateHost ")
	msg := "CreateHost ins"
	return &msg, nil
}

func (i *Service) DeleteBook(p *entities.Person) (*models.Book, error) {
	ins := models.Book{}
	// 查询数据
	err := i.db.First(&ins, p.Id)
	return &ins, err.Error
}
