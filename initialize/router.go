package initialize

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/songcser/gingo/pkg/doc"
	"github.com/songcser/gingo/pkg/response"
	"net/http"
)

var trans ut.Translator

func InitTrans(locale string) (err error) {
	//修改gin框架中Validator引擎属性，实现自定制
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {

		zhT := zh.New() //中文翻译器
		enT := en.New() //英文翻译器
		// 第一个参数是备用(fallback)语言环境
		// 后面参数是应该支持语言环境(可支持多个)
		uni := ut.New(enT, zhT, enT)
		// locale通常取决于http请求\'Accept-language\'
		var ok bool
		trans, ok = uni.GetTranslator(locale)
		if !ok {
			return fmt.Errorf("uni.GetTranslator(%s) failed", locale)
		}
		// 注册翻译器
		switch locale {
		case "en":
			err = enTranslations.RegisterDefaultTranslations(v, trans)
		case "zh":
			err = zhTranslations.RegisterDefaultTranslations(v, trans)
		default:
			err = enTranslations.RegisterDefaultTranslations(v, trans)
		}
		return
	}
	return
}

func HealthCheck(g *gin.Context) {
	g.JSON(http.StatusOK, "ok")
}

func Routers() *gin.Engine {

	if err := InitTrans("zh"); err != nil {
		return nil
	}

	Router := gin.Default()
	//gin.SetMode(gin.DebugMode)

	Router.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if errs, ok := recovered.(validator.ValidationErrors); ok {
			response.FailWithDetailed(errs.Translate(trans), "入参校验失败", c)
			return
		}
		if err, ok := recovered.(error); ok {
			response.FailWithError(err, c)
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	}))
	HealthGroup := Router.Group("")
	{
		// 健康监测
		HealthGroup.GET("/healthcheck", HealthCheck)
	}
	InitSwagger(Router)
	InitAdmin(Router)

	//ApiGroup := Router.Group("api/v1")

	//PrivateGroup := Router.Group("api/v1/private")

	doc.GenerateSwagger()
	return Router
}
