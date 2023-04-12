package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/gphper/multitemplate"
	"github.com/songcser/gingo/config"
	"github.com/songcser/gingo/pkg/auth"
	"github.com/songcser/gingo/pkg/model"
	"github.com/songcser/gingo/pkg/service"
	"github.com/songcser/gingo/templates"
	"github.com/songcser/gingo/utils"
	"net/http"
)

type Admin interface {
	LoginView(c *gin.Context)
	RegisterView(c *gin.Context)
	Login(c *gin.Context)
	Register(c *gin.Context)
	Auth() gin.HandlerFunc
	Home(c *gin.Context)
	List(c *gin.Context)
	AddView(c *gin.Context)
	AddItem(c *gin.Context)
	ViewItem(c *gin.Context)
	EditItem(c *gin.Context)
	DeleteItem(c *gin.Context)
	Render() multitemplate.Renderer
	GetModels() []ModelAdmin
	GetModel(name string) ModelAdmin
}

func New[T model.Model](m T, name string, alias string) {
	admin := BaseModelAdmin[T]{Name: name, Alias: alias, Service: service.NewBaseService(m), model: m}
	factory.Add(admin)
}

func NewAdmin(a ModelAdmin) {
	factory.Add(a)
}

type BaseAdmin struct {
	User auth.User
}

func (BaseAdmin) LoginView(c *gin.Context) {
	c.HTML(200, "login", gin.H{})
}

func (BaseAdmin) RegisterView(c *gin.Context) {
	c.HTML(200, "register", gin.H{})
}

func (b BaseAdmin) Login(c *gin.Context) {
	err := b.User.Login(c)
	utils.CheckError(err)
	c.Redirect(http.StatusFound, HomePath)
}

func (b BaseAdmin) Register(c *gin.Context) {
	err := b.User.Register(c)
	utils.CheckError(err)
	c.Redirect(http.StatusFound, HomePath)
}

func (b BaseAdmin) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !config.GVA_CONFIG.Admin.Auth {
			c.Next()
			return
		}
		err := b.User.Auth(c)
		if err != nil {
			c.Redirect(http.StatusFound, LoginPath)
		}
		c.Next()
	}
}

func (b BaseAdmin) Home(c *gin.Context) {
	admins := b.GetModels()
	user, _ := c.Get("user")
	c.HTML(200, "home", gin.H{
		"admins": admins,
		"user":   user,
	})
}

func (b BaseAdmin) List(c *gin.Context) {
	obj := c.Param("model")
	a := b.GetModel(obj)
	admins := b.GetModels()
	data := a.Query(c)
	result := a.FormatData(a.Header(), data.GetResults())
	current := data.GetCurrent()
	totalPage := int(data.GetTotal())/data.GetSize() + 1
	nextPage := current + 1
	if nextPage > totalPage {
		nextPage = totalPage
	}
	prePage := current - 1
	if prePage <= 0 {
		prePage = 1
	}
	pages := make([]int, 0, 10)
	minPage := current - 4
	if minPage <= 0 {
		minPage = 1
	}
	maxPage := minPage + 9
	if maxPage > totalPage {
		maxPage = totalPage
	}
	for i := minPage; i <= maxPage; i++ {
		pages = append(pages, i)
	}
	user, _ := c.Get("user")
	h := gin.H{
		"user":      user,
		"admins":    admins,
		"header":    a.Header(),
		"results":   result,
		"current":   current,
		"size":      data.GetSize(),
		"total":     data.GetTotal(),
		"totalPage": totalPage,
		"nextPage":  nextPage,
		"prePage":   prePage,
		"admin":     a,
		"name":      a.GetName(),
		"pages":     pages,
	}
	c.HTML(200, "index", h)
}

func (b BaseAdmin) AddView(c *gin.Context) {
	obj := c.Param("model")
	a := b.GetModel(obj)
	forms := a.Form()
	admins := b.GetModels()
	user, _ := c.Get("user")
	h := gin.H{
		"user":   user,
		"name":   a.GetName(),
		"admins": admins,
		"form":   forms,
	}
	c.HTML(200, "add", h)
}

func (b BaseAdmin) AddItem(c *gin.Context) {
	obj := c.Param("model")
	a := b.GetModel(obj)
	err := a.Add(c)
	if err == nil {
		c.Redirect(http.StatusFound, HomePath+obj)
	}
}

func (b BaseAdmin) ViewItem(c *gin.Context) {
	obj := c.Param("model")
	a := b.GetModel(obj)
	data, _ := a.Get(c)
	forms := a.FormValue(data)
	admins := b.GetModels()
	user, _ := c.Get("user")
	h := gin.H{
		"user":   user,
		"id":     data.Get(),
		"name":   a.GetName(),
		"admins": admins,
		"form":   forms,
	}
	c.HTML(200, "edit", h)
}

func (b BaseAdmin) EditItem(c *gin.Context) {
	obj := c.Param("model")
	a := b.GetModel(obj)
	err := a.Edit(c)
	if err == nil {
		c.Redirect(http.StatusFound, HomePath+obj)
	}
}

func (b BaseAdmin) DeleteItem(c *gin.Context) {
	obj := c.Param("model")
	a := b.GetModel(obj)
	err := a.Delete(c)
	if err == nil {
		c.Redirect(http.StatusFound, HomePath+obj)
	}
}

func (BaseAdmin) Render() multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	r.AddFromFs("home", templates.Staticfiles, "home.html", "sidebar.html", "header.html")
	r.AddFromFs("index", templates.Staticfiles, "index.html", "sidebar.html", "header.html")
	r.AddFromFs("add", templates.Staticfiles, "add.html", "sidebar.html", "header.html")
	r.AddFromFs("edit", templates.Staticfiles, "edit.html", "sidebar.html", "header.html")
	r.AddFromFs("login", templates.Staticfiles, "login.html")
	r.AddFromFs("register", templates.Staticfiles, "register.html")
	return r
}

func (BaseAdmin) GetModels() []ModelAdmin {
	return factory.GetAll()
}

func (BaseAdmin) GetModel(name string) ModelAdmin {
	return factory.Get(name)
}
