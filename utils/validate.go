package utils

import (
	"github.com/go-playground/validator/v10"
)

func Validate(val interface{}) {
	validate := validator.New()
	err := validate.Struct(val)
	if err != nil {
		panic(err)
	}
}
