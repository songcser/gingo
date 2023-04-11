package model

import (
	"reflect"
	"strings"
)

const (
	EQ   = iota // ==
	NE          // <>
	GT          // >
	GE          // >=
	LT          // <
	LE          // <=
	LIKE        // like
	IN          // in
)

type item struct {
	Name   string
	Flag   int
	Params interface{}
}

type wrapper struct {
	where    []item
	or       []item
	joins    []string
	s        []string
	distinct []string
}

func (w *wrapper) isZeroOfUnderlyingType(x interface{}) bool {
	return reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}

func (w *wrapper) append(name string, params interface{}, flag int, force bool) {
	if !force && w.isZeroOfUnderlyingType(params) {
		return
	}
	i := item{Name: name, Params: params, Flag: flag}
	w.where = append(w.where, i)
}

func (w *wrapper) Select(s ...string) *wrapper {
	w.s = s
	return w
}

func (w *wrapper) Joins(where string) *wrapper {
	w.joins = append(w.joins, where)
	return w
}

func (w *wrapper) Eq(name string, params interface{}) *wrapper {
	w.append(name, params, EQ, false)
	return w
}

func (w *wrapper) Ne(name string, params interface{}) *wrapper {
	w.append(name, params, NE, false)
	return w
}

func (w *wrapper) NeF(name string, params interface{}) *wrapper {
	w.append(name, params, NE, true)
	return w
}

func (w *wrapper) Gt(name string, params interface{}) *wrapper {
	w.append(name, params, GT, false)
	return w
}

func (w *wrapper) Ge(name string, params interface{}) *wrapper {
	w.append(name, params, GE, false)
	return w
}

func (w *wrapper) Lt(name string, params interface{}) *wrapper {
	w.append(name, params, LT, false)
	return w
}

func (w *wrapper) Le(name string, params interface{}) *wrapper {
	w.append(name, params, LE, false)
	return w
}

func (w *wrapper) Like(name string, params interface{}) *wrapper {
	w.append(name, params, LIKE, false)
	return w
}

func (w *wrapper) In(name string, params interface{}) *wrapper {
	w.append(name, params, IN, false)
	return w
}

func (w *wrapper) Distinct(args []string) *wrapper {
	w.distinct = args
	return w
}

func (w *wrapper) Or(name string, params interface{}) *wrapper {
	w.or = append(w.or, item{Name: name, Params: params, Flag: EQ})
	return w
}

func (w *wrapper) Where() (string, []interface{}) {
	var params []interface{}
	var builder strings.Builder
	for i, item := range w.where {
		if i > 0 {
			builder.WriteString(" AND ")
		}
		builder.WriteString(item.Name)
		switch item.Flag {
		case EQ:
			builder.WriteString(" = ?")
			params = append(params, item.Params)
		case NE:
			builder.WriteString(" <> ?")
			params = append(params, item.Params)
		case GT:
			builder.WriteString(" > ?")
			params = append(params, item.Params)
		case GE:
			builder.WriteString(" >= ?")
			params = append(params, item.Params)
		case LT:
			builder.WriteString(" < ?")
			params = append(params, item.Params)
		case LE:
			builder.WriteString(" <= ?")
			params = append(params, item.Params)
		case LIKE:
			builder.WriteString(" like ?")
			params = append(params, "%"+item.Params.(string)+"%")
		case IN:
			builder.WriteString(" in ?")
			params = append(params, item.Params)
		}
	}
	if len(w.or) > 0 {
		builder.WriteString("AND ( ")
		for i, item := range w.or {
			if i > 0 {
				builder.WriteString(" OR ")
			}
			builder.WriteString(item.Name)
			builder.WriteString(" = ?")
			params = append(params, item.Params)
		}
		builder.WriteString(" )")
	}
	return builder.String(), params
}

func NewWrapper() *wrapper {
	return &wrapper{}
}
