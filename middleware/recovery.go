package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/songcser/gingo/config"
	"github.com/songcser/gingo/pkg/response"
	"github.com/songcser/gingo/utils"
	"net/http"
)

func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if errs, ok := recovered.(validator.ValidationErrors); ok {
			response.FailWithDetailed(errs.Translate(utils.Trans), "入参校验失败", c)
			return
		}
		if err, ok := recovered.(error); ok {
			config.GVA_LOG.Error(err.Error())
			response.FailWithError(err, c)
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	})
}
