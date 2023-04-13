package app

import "github.com/songcser/gingo/pkg/model"

type App struct {
	model.BaseModel
	Name        string `json:"name" form:"name" gorm:"column:name;type:varchar(255);not null" admin:"type:input;name:name;label:应用名"`
	Description string `json:"description" form:"description" gorm:"column:description;type:varchar(4096);not null" admin:"type:textarea;name:description;label:描述"`
	Level       string `json:"level" form:"level" gorm:"column:level;type:varchar(8);not null" admin:"type:radio;enum:S1,S2,S3,S4,S5;label:级别"`
	Type        string `json:"type" form:"type" gorm:"column:type;type:varchar(16);not null" admin:"type:select;enum:container=容器应用,web=前端应用,mini=小程序应用;label:应用类型"`
}
