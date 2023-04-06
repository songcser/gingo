package initialize

import (
	"github.com/gin-gonic/gin"
	"github.com/songcser/gingo/pkg/admin"
)

func InitAdmin(r *gin.Engine) {
	admin.Init(r, nil)
}
