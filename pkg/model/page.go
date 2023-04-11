package model

import "github.com/songcser/gingo/utils"

type Page interface {
	SetTotal(total int64)
	GetTotal() int64
	GetSize() int
	GetCurrent() int
	GetResults() *[]Model
}

type BasePage[T Model] struct {
	Total   int64 `json:"total"`
	Size    int   `json:"size"`
	Current int   `json:"current"`
	Results *[]T  `json:"results"`
}

func (p BasePage[T]) SetTotal(total int64) {
	p.Total = total
}

func (p BasePage[T]) GetTotal() int64 {
	return p.Total
}

func (p BasePage[T]) GetSize() int {
	return p.Size
}

func (p BasePage[T]) GetCurrent() int {
	return p.Current
}

func (p BasePage[T]) GetResults() *[]Model {
	return utils.Map(p.Results, func(t T) Model {
		return t
	})
}
