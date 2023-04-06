package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/songcser/gingo/pkg/auth"
)

func Init(r *gin.Engine, a Admin) {

	if a == nil {
		a = BaseAdmin{
			User: auth.BaseUser{},
		}
	}
	r.HTMLRender = a.Render()
	r.GET("/admin/login", a.LoginView)
	r.GET("/admin/register", a.RegisterView)
	r.POST("/admin/login", a.Login)
	r.POST("/admin/register", a.Register)

	g := r.Group("admin")
	g.Use(a.Auth())
	g.GET("/", a.Home)
	g.GET(":model", a.List)
	g.GET(":model/form", a.AddView)
	g.POST(":model/add", a.AddItem)
	g.GET(":model/:id", a.ViewItem)
	g.POST(":model/:id", a.EditItem)
	g.GET(":model/delete/:id", a.DeleteItem)
}
