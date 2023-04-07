package app

import (
	"github.com/songcser/gingo/pkg/admin"
)

func Admin() {
	var a App
	admin.New(a, "app", "应用")
}
