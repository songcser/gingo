package app

import (
	"github.com/gin-gonic/gin"
	"github.com/songcser/gingo/pkg/api"
	"github.com/songcser/gingo/pkg/response"
)

type Api struct {
	api.Api
	Service Service
}

func NewApi() Api {
	var app App
	s := NewService(app)
	baseApi := api.NewApi[App](s)
	return Api{Api: baseApi, Service: s}
}

func (a Api) Hello(c *gin.Context) {
	str := a.Service.Hello()
	response.OkWithData(str, c)
}
