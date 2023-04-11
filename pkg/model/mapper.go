package model

import (
	"github.com/songcser/gingo/config"
	"github.com/songcser/gingo/utils"
	"gorm.io/gorm"
)

type Mapper[T Model] interface {
	OrderBy(o string)
	Insert(data any) error
	Update(column string, value any) error
	Updates(value T) error
	UpdatesById(id int64, data any) error
	DeleteById(id int64) error
	SelectById(id int64) error
	Select() (*[]T, error)
	SelectOne() (T, error)
	GetById(id int64) (T, error)
	Count() int64
	SelectPage(size, current int) (Page, error)
}

type BaseMapper[T Model] struct {
	t     T
	w     *wrapper
	order string
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

func (m *BaseMapper[T]) Updates(value T) error {
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

func (m *BaseMapper[T]) GetById(id int64) (T, error) {
	var val T
	err := config.GVA_DB.Model(m.t).First(&val, id).Error
	return val, err
}

func (m *BaseMapper[T]) Count() int64 {
	var count int64
	db := m.sql(0, 0)
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

	var db *gorm.DB
	if m.w != nil {
		w := m.w
		db = config.GVA_DB.Model(m.t)
		if len(w.s) > 0 {
			db.Select(w.s)
		}
		if len(w.joins) > 0 {
			for _, join := range w.joins {
				db.Joins(join)
			}
		}
		if len(w.distinct) > 0 {
			db.Distinct(w.distinct)
		}
		where, params := w.Where()
		db.Where(where, params...)
	} else {
		db = config.GVA_DB.Model(m.t).Where(m.t)
	}

	if limit > 0 {
		db.Limit(limit)
	}
	if offset > 0 {
		db.Offset(offset)
	}

	if m.order != "" {
		db.Order(m.order)
	}

	return db
}

func (m *BaseMapper[T]) Select() (*[]T, error) {
	db := m.sql(0, 0)
	var val []T
	err := db.Find(&val).Error
	res := utils.AsyncMap(&val, func(t1 T) T {
		return t1
	})
	return res, err
}

func (m *BaseMapper[T]) SelectOne() (T, error) {
	db := m.sql(0, 0)
	var val T
	err := db.First(&val).Error
	return val, err
}

func (m *BaseMapper[T]) SelectPage(size, current int) (Page, error) {
	res := make([]T, 0, size)
	page := BasePage[T]{Size: size, Current: current, Results: &res}

	count := m.Count()
	page.Total = count
	limit := page.Size
	offset := 0
	if page.Current > 0 {
		offset = (page.Current - 1) * page.Size
	}
	db := m.sql(limit, offset)
	err := db.Find(page.Results).Error
	return page, err
}

func NewMapper[T Model](t T, w *wrapper) Mapper[T] {
	return &BaseMapper[T]{t: t, w: w}
}
