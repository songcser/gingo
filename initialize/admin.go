package initialize

import (
	"github.com/gin-gonic/gin"
	"github.com/songcser/gingo/config"
	"github.com/songcser/gingo/internal/app"
	"github.com/songcser/gingo/pkg/admin"
)

func Admin(r *gin.Engine) {
	if !config.GVA_CONFIG.Admin.Enable {
		return
	}
	admin.Init(r, nil)
	app.Admin()
}
