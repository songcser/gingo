package app

import "github.com/songcser/gingo/pkg/model"

type App struct {
	model.BaseModel
	Name        string `json:"name" form:"name" admin:"type:input;name:name;label:应用名"`
	Description string `json:"description" form:"description" admin:"type:textarea;name:description;label:应用名"`
}
