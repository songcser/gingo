package initialize

import (
	"github.com/gin-gonic/gin"
	"github.com/songcser/gingo/pkg/doc"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Swagger(r *gin.Engine) {
	//docs.SwaggerInfo.BasePath = "/api/v1"
	doc.InitSwagger(doc.SwaggerConfig{Title: "service api", Description: "service api"})
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
