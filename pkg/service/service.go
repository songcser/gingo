package service

import (
	"github.com/gin-gonic/gin"
	"github.com/songcser/gingo/pkg/model"
	"github.com/songcser/gingo/utils"
)

type Service[T model.Model] interface {
	MakeMapper(c *gin.Context) model.Mapper[T]
	MakeResponse(val model.Model) any
	Query(c *gin.Context, mapper model.Mapper[T]) model.Page
	Get(id int64) (T, error)
	Create(data any) error
	Update(id int64, data T) error
	Delete(id int64) error
}

type BaseService[T model.Model] struct {
	Model T
}

func NewBaseService[T model.Model](t T) BaseService[T] {
	return BaseService[T]{Model: t}
}

type pageReq struct {
	Size    int `form:"size"`
	Current int `form:"current"`
}

func (bs BaseService[T]) MakeResponse(val model.Model) any {
	return val
}

func (bs BaseService[T]) MakeMapper(c *gin.Context) model.Mapper[T] {
	m := model.NewMapper[T](bs.Model, nil)
	return m
}

func (bs BaseService[T]) Query(c *gin.Context, mapper model.Mapper[T]) model.Page {
	var req pageReq
	err := c.ShouldBindQuery(&req)
	utils.CheckError(err)
	if req.Size == 0 {
		req.Size = 10
	}
	if req.Current == 0 {
		req.Current = 1
	}
	mapper.OrderBy("created_at desc")
	page, _ := mapper.SelectPage(req.Size, req.Current)
	return page
}

func (bs BaseService[T]) Get(id int64) (T, error) {
	mapper := model.NewMapper(bs.Model, nil)
	val, err := mapper.GetById(id)
	return val, err
}

func (bs BaseService[T]) Create(data any) error {
	mapper := model.NewMapper(bs.Model, nil)
	return mapper.Insert(data)
}

func (bs BaseService[T]) Update(id int64, data T) error {
	mapper := model.NewMapper(bs.Model, nil)
	return mapper.UpdatesById(id, data)
}

func (bs BaseService[T]) Delete(id int64) error {
	mapper := model.NewMapper(bs.Model, nil)
	return mapper.DeleteById(id)
}
