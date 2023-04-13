package router

import (
	"github.com/gin-gonic/gin"
	"github.com/songcser/gingo/pkg/api"
	"net/http"
)

type DefaultRouter struct {
	Router *gin.RouterGroup
}

func NewRouter(r *gin.RouterGroup) DefaultRouter {
	return DefaultRouter{Router: r}
}

func (d DefaultRouter) BindApi(path string, a api.Api) {
	d.BindGet(path, a.Query)
	d.BindGet(path+"/:id", a.Get)
	d.BindPost(path, a.Create)
	d.BindPut(path+"/:id", a.Update)
	d.BindDelete(path+"/:id", a.Delete)
}

func (d DefaultRouter) BindGet(path string, handle gin.HandlerFunc) {
	d.Bind(http.MethodGet, path, handle)
}

func (d DefaultRouter) BindPost(path string, handle gin.HandlerFunc) {
	d.Bind(http.MethodPost, path, handle)
}

func (d DefaultRouter) BindPut(path string, handle gin.HandlerFunc) {
	d.Bind(http.MethodPut, path, handle)
}

func (d DefaultRouter) BindDelete(path string, handle gin.HandlerFunc) {
	d.Bind(http.MethodDelete, path, handle)
}

func (d DefaultRouter) Bind(method string, path string, handle gin.HandlerFunc) {
	d.Router.Handle(method, path, handle)
}
