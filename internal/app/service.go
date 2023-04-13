package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/songcser/gingo/pkg/model"
	"github.com/songcser/gingo/pkg/service"
	"github.com/songcser/gingo/utils"
)

type Service struct {
	service.Service[App]
}

func NewService(a App) Service {
	return Service{service.NewBaseService[App](a)}
}

func (s Service) MakeMapper(c *gin.Context) model.Mapper[App] {
	var r Request
	err := c.ShouldBindQuery(&r)
	utils.CheckError(err)
	w := model.NewWrapper()
	w.Like("name", r.Name)
	w.Eq("level", r.Level)
	m := model.NewMapper[App](App{}, w)
	return m
}

func (s Service) MakeResponse(val model.Model) any {
	a := val.(App)
	res := Response{
		Name:        a.Name,
		Description: fmt.Sprintf("名称：%s, 等级: %s, 类型: %s", a.Name, a.Level, a.Type),
		Level:       a.Level,
		Type:        a.Type,
	}
	return res
}

func (s Service) Hello() string {
	return "Hello World"
}
