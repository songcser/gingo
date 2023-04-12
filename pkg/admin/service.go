package admin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/songcser/gingo/pkg/model"
	"github.com/songcser/gingo/pkg/service"
	"github.com/songcser/gingo/utils"
	"reflect"
	"strconv"
	"strings"
)

type ModelAdmin interface {
	Query(c *gin.Context) model.Page
	Add(c *gin.Context) error
	Edit(c *gin.Context) error
	Get(c *gin.Context) (model.Model, error)
	Delete(c *gin.Context) error
	GetName() string
	Header() *[]Header
	Form() *[]Form
	FormValue(data model.Model) *[]Form
	FormatData(header *[]Header, data *[]model.Model) [][]string
}

type ModelAdminFactory struct {
	models []ModelAdmin
}

func (f *ModelAdminFactory) Add(a ModelAdmin) {
	f.models = append(f.models, a)
}

func (f *ModelAdminFactory) Get(name string) ModelAdmin {
	for _, a := range f.models {
		if a.GetName() == name {
			return a
		}
	}
	return nil
}

func (f *ModelAdminFactory) GetAll() []ModelAdmin {
	return f.models
}

var factory = &ModelAdminFactory{models: make([]ModelAdmin, 0)}

type BaseModelAdmin[T model.Model] struct {
	model   T
	Name    string
	Alias   string
	Service service.Service[T]
}

func (b BaseModelAdmin[T]) Add(c *gin.Context) error {
	var data T
	err := c.ShouldBind(&data)
	utils.CheckError(err)
	err = b.Service.Create(&data)
	return err
}

func (b BaseModelAdmin[T]) Edit(c *gin.Context) error {
	i := c.Param("id")
	id, _ := strconv.ParseInt(i, 10, 64)
	var data T
	err := c.ShouldBind(&data)
	utils.CheckError(err)
	err = b.Service.Update(id, data)
	return err
}

func (b BaseModelAdmin[T]) GetName() string {
	return b.Name
}

func (b BaseModelAdmin[T]) Header() *[]Header {
	m := b.model
	var header []Header
	t := reflect.TypeOf(m)
	v := reflect.ValueOf(m)
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).CanInterface() { //判断是否为可导出字段
			val := v.Field(i)
			if val.Kind() == reflect.Struct {
				typ := reflect.TypeOf(val.Interface())
				for j := 0; j < typ.NumField(); j++ {
					b.addHeaderByTag(typ.Field(j).Tag, &header)
				}
			} else {
				b.addHeaderByTag(t.Field(i).Tag, &header)
			}
		}
	}
	return &header
}

func (b BaseModelAdmin[T]) addHeaderByTag(tags reflect.StructTag, header *[]Header) {
	tag, ok := b.parseTag(tags)
	if ok {
		h := Header{Name: tag.Name, Label: tag.Label}
		*header = append(*header, h)
	}
}

func (b BaseModelAdmin[T]) Query(c *gin.Context) model.Page {
	mapper := model.NewMapper[T](b.model, model.NewWrapper())
	result := b.Service.Query(c, mapper)
	return result
}

func (b BaseModelAdmin[T]) Form() *[]Form {
	a := b.model
	t := reflect.TypeOf(a)
	v := reflect.ValueOf(a)
	forms := make([]Form, 0, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).CanInterface() {
			val := v.Field(i)
			if val.Kind() == reflect.Struct {
				typ := reflect.TypeOf(val.Interface())
				for j := 0; j < typ.NumField(); j++ {
					b.addFormByTag(typ.Field(j).Tag, &forms)
				}
			} else {
				b.addFormByTag(t.Field(i).Tag, &forms)
			}
		}
	}
	return &forms
}

func (b BaseModelAdmin[T]) FormValue(data model.Model) *[]Form {
	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)
	forms := make([]Form, 0, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).CanInterface() {
			val := v.Field(i)
			if val.Kind() == reflect.Struct {
				typ := reflect.TypeOf(val.Interface())
				vv := reflect.ValueOf(val.Interface())
				for j := 0; j < typ.NumField(); j++ {
					b.addFormValueByTag(typ.Field(j).Tag, &forms, vv.Field(j).Interface())
				}
			} else {
				b.addFormValueByTag(t.Field(i).Tag, &forms, val.Interface())
			}
		}
	}
	return &forms
}

func (b BaseModelAdmin[T]) parseTag(tag reflect.StructTag) (Tag, bool) {
	tags := Tag{Admin: false}

	gormTag := tag.Get("gorm")
	if gormTag != "" && gormTag != "-" {
		labels := strings.Split(gormTag, ";")
		for _, label := range labels {
			lab := strings.Split(label, ":")
			if len(lab) < 2 {
				continue
			}
			if lab[0] == "comment" && tags.Label == "" {
				tags.Label = lab[1]
			} else if lab[0] == "type" && tags.Type == "" {
				if strings.HasPrefix(lab[1], "varchar") || strings.HasPrefix(lab[1], "char") {
					tags.Type = "text"
				} else if strings.HasPrefix(lab[1], "int") || strings.HasPrefix(lab[1], "bigint") {
					tags.Type = "number"
				} else if strings.HasPrefix(lab[1], "tinyint") {
					tags.Type = "number"
				}
			}
		}
	}

	jsonTag := tag.Get("json")
	if jsonTag != "" && jsonTag != "-" {
		if tags.Name == "" {
			tags.Name = jsonTag
		}
		if tags.Label == "" {
			tags.Label = jsonTag
		}
		if tags.Type == "" {
			tags.Type = "input"
		}
	}

	adminTag := tag.Get("admin")
	if adminTag != "" && adminTag != "-" {
		tags.Admin = true
		labels := strings.Split(adminTag, ";")
		for _, label := range labels {
			lab := strings.Split(label, ":")
			if len(lab) == 1 {
				if lab[0] == "disable" {
					tags.Disable = true
				}
				continue
			} else if len(lab) < 2 {
				continue
			}
			labelName := lab[0]
			labelValue := lab[1]
			if labelName == "name" {
				if labelValue == "" || labelValue == "-" {
					return Tag{}, false
				}
				tags.Name = labelValue
			} else if labelName == "type" {
				tags.Type = labelValue
			} else if labelName == "label" {
				tags.Label = labelValue
			} else if labelName == "enum" {
				values := strings.Split(labelValue, ",")
				enums := make([]Enum, 0, len(values))
				for _, enum := range values {
					e := strings.Split(enum, "=")
					if len(e) > 1 {
						enums = append(enums, Enum{Key: e[0], Value: e[1]})
					} else {
						enums = append(enums, Enum{Key: e[0], Value: e[0]})
					}
				}
				tags.Enum = enums
			}
		}
	}

	return tags, true
}

func (b BaseModelAdmin[T]) addFormByTag(tag reflect.StructTag, forms *[]Form) {
	tags, ok := b.parseTag(tag)
	if ok && tags.Name != "id" && tags.Name != "createdAt" && tags.Name != "updatedAt" {
		f := Form{
			Label:   tags.Label,
			Type:    tags.Type,
			Name:    tags.Name,
			Enum:    tags.Enum,
			Disable: tags.Disable,
		}
		*forms = append(*forms, f)
	}
}

func (b BaseModelAdmin[T]) addFormValueByTag(tag reflect.StructTag, forms *[]Form, value any) {
	tags, ok := b.parseTag(tag)
	if ok {
		if val, o := value.(utils.JsonTime); o {
			value = utils.JsonTimeFormat(val)
		}
		f := Form{
			Label:   tags.Label,
			Type:    tags.Type,
			Name:    tags.Name,
			Value:   value,
			Enum:    tags.Enum,
			Disable: tags.Disable,
		}
		*forms = append(*forms, f)
	}
}

func (b BaseModelAdmin[T]) formToMap(forms *[]Form) map[string]Form {
	res := make(map[string]Form, len(*forms))
	for _, f := range *forms {
		res[f.Name] = f
	}
	return res
}

func (b BaseModelAdmin[T]) FormatData(header *[]Header, data *[]model.Model) [][]string {
	results := make([][]string, 0, len(*data))
	for _, m := range *data {
		form := b.FormValue(m)
		formMap := b.formToMap(form)
		res := make([]string, 0, len(*form))
		for _, h := range *header {
			if f, ok := formMap[h.Name]; ok {
				value := fmt.Sprintf("%v", f.Value)
				for _, e := range f.Enum {
					if e.Key == value {
						value = e.Value
						break
					}
				}
				res = append(res, value)
			} else {
				res = append(res, "")
			}
		}
		results = append(results, res)
	}
	return results
}

func (b BaseModelAdmin[T]) Get(c *gin.Context) (model.Model, error) {
	i := c.Param("id")
	id, _ := strconv.ParseInt(i, 10, 64)
	return b.Service.Get(id)
}

func (b BaseModelAdmin[T]) Delete(c *gin.Context) error {
	i := c.Param("id")
	id, _ := strconv.ParseInt(i, 10, 64)
	return b.Service.Delete(id)
}
