package api

import (
	"github.com/gin-gonic/gin"
	"github.com/songcser/gingo/pkg/model"
	"github.com/songcser/gingo/pkg/response"
	"github.com/songcser/gingo/pkg/service"
	"github.com/songcser/gingo/utils"
	"strconv"
)

type Api interface {
	Query(c *gin.Context)
	Get(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type BaseApi[T model.Model] struct {
	Service service.Service[T]
}

func NewApi[T model.Model](s service.Service[T]) BaseApi[T] {
	return BaseApi[T]{Service: s}
}

func (b BaseApi[T]) parseId(c *gin.Context) int64 {
	i := c.Param("id")
	if i == "" {
		return 0
	}
	id, err := strconv.ParseInt(i, 10, 64)
	utils.CheckError(err)
	return id
}

func (b BaseApi[T]) Query(c *gin.Context) {
	mapper := b.Service.MakeMapper(c)
	result := b.Service.Query(c, mapper)
	page := response.NewPage(result)
	page.Results = utils.AsyncMap(result.GetResults(), func(t1 model.Model) any {
		return b.Service.MakeResponse(t1)
	})
	response.OkWithData(page, c)
}

func (b BaseApi[T]) Get(c *gin.Context) {
	id := b.parseId(c)
	obj, err := b.Service.Get(id)
	val := b.Service.MakeResponse(obj)
	utils.CheckError(err)
	response.OkWithData(val, c)
}

func (b BaseApi[T]) Create(c *gin.Context) {
	var data T
	err := c.ShouldBind(&data)
	utils.CheckError(err)
	err = b.Service.Create(&data)
	utils.CheckError(err)
	response.OkWithData(true, c)
}

func (b BaseApi[T]) Update(c *gin.Context) {
	id := b.parseId(c)
	var data T
	err := c.ShouldBind(&data)
	utils.CheckError(err)
	err = b.Service.Update(id, data)
	utils.CheckError(err)
	response.OkWithData(true, c)
}

func (b BaseApi[T]) Delete(c *gin.Context) {
	id := b.parseId(c)
	err := b.Service.Delete(id)
	utils.CheckError(err)
	response.OkWithData(true, c)
}
