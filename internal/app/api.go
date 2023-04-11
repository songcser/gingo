package app

import (
	"github.com/songcser/gingo/pkg/api"
	"github.com/songcser/gingo/pkg/service"
)

type Api struct {
	api.Api
}

func NewApi() Api {
	var app App
	baseApi := api.NewApi[App](service.NewBaseService(app))
	return Api{baseApi}
}
