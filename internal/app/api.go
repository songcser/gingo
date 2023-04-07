package app

import (
	"github.com/gin-gonic/gin"
	"github.com/songcser/gingo/pkg/api"
	"github.com/songcser/gingo/pkg/service"
	"net/http"
)

type Api struct {
	api.Api
}

func NewApi() Api {
	var app App
	baseApi := api.NewApi[App](service.NewBaseService(app))
	return Api{baseApi}
}

func (a Api) Recipients(c *gin.Context) {
	users := []string{"jiyi.song", "yanlong.liang", "xuewen.zhang01"}
	c.JSON(http.StatusOK, struct {
		Errcode int      `json:"errcode"`
		Errmsg  string   `json:"errmsg"`
		Data    []string `json:"data"`
	}{Errcode: 0, Errmsg: "sucess", Data: users})
}
