package model

import (
	"github.com/songcser/gingo/config"
	"github.com/songcser/gingo/utils"
	"gorm.io/gorm"
)

type Mapper interface {
	OrderBy(o string)
	Insert(data any) error
	Update(column string, value any) error
	Updates(value Model) error
	UpdatesById(id int64, data any) error
	DeleteById(id int64) error
	SelectById(id int64) error
	Select() (*[]Model, error)
	SelectOne() (Model, error)
	QueryOne() (Model, error)
	GetById(id int64) (Model, error)
	Count() int64
	SelectPage(size, current int) (Page, error)
	QueryPage(size, current int) (Page, error)
	Query() (*[]Model, error)
	CheckWrapper() bool
}

type BaseMapper[T Model] struct {
	t     T
	w     *wrapper
	order string
}

func (m *BaseMapper[T]) CheckWrapper() bool {
	return m.w != nil
}

func (m *BaseMapper[T]) OrderBy(o string) {
	m.order = o
}

func (m *BaseMapper[T]) Insert(data any) error {
	err := config.GVA_DB.Model(m.t).Create(data).Error
	return err
}

func (m *BaseMapper[T]) Update(column string, value any) error {
	return config.GVA_DB.Model(m.t).Update(column, value).Error
}

func (m *BaseMapper[T]) Updates(value Model) error {
	return config.GVA_DB.Model(m.t).Updates(value).Error
}

func (m *BaseMapper[T]) UpdatesById(id int64, value any) error {
	return config.GVA_DB.Model(m.t).Select("*").
		Omit("created_at", "id").
		Where("id = ?", id).
		Updates(value).Error
}

func (m *BaseMapper[T]) DeleteById(id int64) error {
	err := config.GVA_DB.Model(m.t).Delete("id", id).Error
	return err
}

func (m *BaseMapper[T]) SelectById(id int64) error {
	err := config.GVA_DB.Model(m.t).First(m.t, id).Error
	return err
}

func (m *BaseMapper[T]) GetById(id int64) (Model, error) {
	var val T
	err := config.GVA_DB.Model(m.t).First(&val, id).Error
	return val, err
}

func (m *BaseMapper[T]) Count() int64 {
	w := m.w
	var count int64
	db := m.sql(0, 0)
	if len(w.distinct) > 0 {
		db.Distinct(w.distinct[0])
	}
	err := db.Count(&count).Error
	if err != nil {
		panic(err)
	}
	return count
}

func (m *BaseMapper[T]) QueryCount() int64 {
	var count int64
	db := config.GVA_DB.Model(m.t).Where(m.t)
	err := db.Count(&count).Error
	if err != nil {
		panic(err)
	}
	return count
}

func (m *BaseMapper[T]) sql(limit int, offset int) *gorm.DB {
	w := m.w
	db := config.GVA_DB.Model(m.t)
	if len(w.s) > 0 {
		db.Select(w.s)
	}
	if limit > 0 {
		db.Limit(limit)
	}
	if offset > 0 {
		db.Offset(offset)
	}
	if len(w.joins) > 0 {
		for _, join := range w.joins {
			db.Joins(join)
		}
	}
	if m.order != "" {
		db.Order(m.order)
	}
	where, params := w.Where()
	db.Where(where, params...)
	return db
}

func (m *BaseMapper[T]) Select() (*[]Model, error) {
	w := m.w
	db := m.sql(0, 0)
	if len(w.distinct) > 0 {
		db.Distinct(w.distinct)
	}
	var val []T
	err := db.Find(&val).Error
	res := utils.AsyncMap(&val, func(t1 T) Model {
		return t1
	})
	return res, err
}

func (m *BaseMapper[T]) SelectOne() (Model, error) {
	w := m.w
	db := m.sql(0, 0)
	if len(w.distinct) > 0 {
		db.Distinct(w.distinct)
	}
	var val T
	err := db.First(&val).Error
	return val, err
}

func (m *BaseMapper[T]) QueryOne() (Model, error) {
	var val T
	err := config.GVA_DB.Model(m.t).Where(m.t).First(&val).Error
	return val, err
}

func (m *BaseMapper[T]) Query() (*[]Model, error) {
	var val []T
	err := config.GVA_DB.Model(m.t).Where(m.t).Find(&val).Error
	if err != nil {
		return nil, err
	}
	res := utils.AsyncMap(&val, func(t1 T) Model {
		return t1
	})
	return res, err
}

func (m *BaseMapper[T]) QueryPage(size, current int) (Page, error) {
	res := make([]T, 0, size)
	page := BasePage[T]{Size: size, Current: current, Results: &res}
	count := m.QueryCount()
	page.Total = count
	limit := page.Size
	offset := 0
	if page.Current > 0 {
		offset = (page.Current - 1) * page.Size
	}
	db := config.GVA_DB.Model(m.t).Where(m.t)
	if limit > 0 {
		db.Limit(limit)
	}
	if offset > 0 {
		db.Offset(offset)
	}
	if m.order != "" {
		db.Order(m.order)
	}
	err := db.Find(page.Results).Error
	return page, err
}

func (m *BaseMapper[T]) SelectPage(size, current int) (Page, error) {
	res := make([]T, 0, size)
	page := BasePage[T]{Size: size, Current: current, Results: &res}

	w := m.w
	count := m.Count()
	page.Total = count
	limit := page.Size
	offset := 0
	if page.Current > 0 {
		offset = (page.Current - 1) * page.Size
	}
	db := m.sql(limit, offset)
	if len(w.distinct) > 0 {
		db.Distinct(w.distinct)
	}
	err := db.Find(page.Results).Error
	return page, err
}

func NewMapper[T Model](t T, w *wrapper) Mapper {
	return &BaseMapper[T]{t: t, w: w}
}
