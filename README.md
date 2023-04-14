# Gingo

## Introduce

Gingo是基于 gin 框架为核心的脚手架，使用本项目可以快速完成业务逻辑开发。

## Feature

* gin框架，简单，高效，轻量
* gorm数据库ORM框架，封装mapper，使用简单
* viper配置管理
* zap日志框架，输出日志更灵活
* api接口封装，快速实现CURD操作，提供Restful风格接口
* admin后台管理，实现了简单的后台管理，方便进行数据管理
* 使用范型，go版本不能低于1.18

## Catalogue

```
.
|——.gitignore
|——go.mod
|——go.sum
|——cmd          
   └──migrate
      └──main.go            // 注册数据库表
   └──main.go               // 项目入口main
|——README.md
|——config                   // 配置文件目录
|  └──autoload              // 配置文件的结构体定义包
|     └──admin.go           // admin配置
|     └──db.go
|     └──jwt.go             // jwt配置
|     └──mysql.go           // mysql配置
|     └──zap.go             // zap日志配置
|  └──config.yaml           // .yaml配置示例文件
|  └──config.go             // 配置初始化文件
|——initialize               // 数据初始化目录
|  └──admin.go              // admin初始化
|  └──constants.go          // 常量数据
|  └──gorm.go               // 数据库初始化
|  └──mysql.go              // mysql初始化
|  └──router.go             // gin初始化
|  └──swagger.go            
|  └──viper.go              // viper配置初始化
|  └──zap.go                // zap日志初始化
|——internal                 // 该服务所有不对外暴露的代码，通常的业务逻辑都在这下面，使用internal避免错误引用
|──middleware               // 中间件目录
|  └──logger.go             // 日志中间件，打印请求数据
|  └──recovery.go           // 自定义recovery, 输出错误格式化
|──pkg                      // 内部服务包    
|  └──admin                 // admin实现逻辑
|     └──admin.go
|     └──init.go
|     └──model.go
|     └──service.go
|  └──api                   // API接口封装
|     └──api.go
|  └──auth                  // 登陆授权接口封装
|     └──model.go           
|     └──user.go            
|  └──model                 // 底层模型封装
|     └──mapper.go           
|     └──model.go           
|     └──page.go           
|     └──wrapper.go           
|  └──response              // 响应数据模型封装
|     └──page.go         
|     └──response.go         
|  └──router                // 路由模块封装
|     └──router.go          
|  └──service               // 服务模块封装
|     └──service.go
|──templates                // admin模版页面
|  └──add.html              // 新建页面
|  └──edit.html             // 编辑页面
|  └──embed.go          
|  └──header.html           // 头部页面
|  └──home.html             // 首页
|  └──index.html            // 主页面
|  └──login.html            // 登陆页面
|  └──register.html         // 注册页面页面
|  └──sidebar.html          // 左边栏页面
|──utils                    // 一些工具方法
|  └──cache.go              // 缓存
|  └──error.go              // error检查
|  └──hash.go               // hash加密解密
|  └──http.go               // http客户端请求
|  └──json.go               // 
|  └──jwt.go                // JWT
|  └──path.go               // 文件路径
|  └──time.go               // time相关方法
|  └──translator.go         // 中英文翻译
```

## Usage

internal目录是不对外暴露的代码，在做go get时，此目录不会被下载，所以通常业务逻辑放在这个下面。
我们将在这个目录下面加一些业务代码，说明脚手架的使用。

新增 app 包目录

### Model

```go
// Package app model.go
package app

import "github.com/songcser/gingo/pkg/model"

type App struct {
	model.BaseModel
	Name        string `json:"name" gorm:"column:name;type:varchar(255);not null"`
	Description string `json:"description" gorm:"column:description;type:varchar(4096);not null"`
	Level       string `json:"level" gorm:"column:level;type:varchar(8);not null"`
	Type        string `json:"type" gorm:"column:type;type:varchar(16);not null"`
}
```
App模型有4个自定义字段，gorm标签会对应到数据库的字段。

```go
package model

type Model interface {
	Get() int64
}

type BaseModel struct {
	ID        int64          `json:"id" gorm:"primarykey" admin:"disable"`                // 主键ID
	CreatedAt utils.JsonTime `json:"createdAt" gorm:"index;comment:创建时间" admin:"disable"` // 创建时间
	UpdatedAt utils.JsonTime `json:"updatedAt" gorm:"index;comment:更新时间" admin:"disable"` // 更新时间
}

func (m BaseModel) Get() int64 {
	return m.ID
}
```
BaseModel 是基础模型，有一些公共字段, 并且实现了 Model interface, 所有引用 BaseModel 的模型都实现了 Model interface。

创建数据库表
```go
# Package initialize gorm.go
package initialize
// RegisterTables 注册数据库表专用
func RegisterTables(db *gorm.DB) {
	err := db.Set("gorm:table_options", "CHARSET=utf8mb4").AutoMigrate(
		// 系统模块表
		auth.BaseUser{},
		app.App{}, // app表注册
	)
	if err != nil {
		os.Exit(0)
	}
}

```
执行 migrate 的 main 方法会在数据库创建对应的表。


### Api

```go
// Package app api.go
package app

import (
	"github.com/songcser/gingo/pkg/api"
	"github.com/songcser/gingo/pkg/service"
)

type Api struct {
	api.Api
}

func NewApi() Api {
	var app App
	baseApi := api.NewApi[App](service.NewBaseService(app))
	return Api{baseApi}
}
```
api.Api接口
```go
// Package api api.go
package api

import "github.com/gin-gonic/gin"

type Api interface {
	Query(c *gin.Context)
	Get(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

```

api.Api接口定义了CURD方法，并且方法都是gin.HandlerFunc类型，可以直接绑定到gin Router上。
BaseApi实现了CURD的基本方法，app.Api类型组合了BaseApi的方法。

### Router

```go
// Package app router.go
package app

import (
	"github.com/gin-gonic/gin"
	"github.com/songcser/gingo/pkg/router"
)

func InitRouter(g *gin.RouterGroup) {
	r := router.NewRouter(g.Group("app"))
	a := NewApi()
	r.BindApi("", a)
}
```

router 是对gin.RouterGroup做了简单封装，方便和Api类型做绑定。
BindApi方法将Api的 CURD 方法和router进行了绑定。

启动服务之后，执行脚本或者使用 postman 请求服务

> 创建数据

```shell
curl --location 'http://localhost:8080/api/v1/app' \
--header 'Content-Type: application/json' \
--data '{
    "name": "测试应用",
    "description": "测试应用服务",
    "level": "S3",
    "type": "container"
}'
```
返回内容
```json
{
    "code": 0,
    "data": true,
    "message": "success"
}
```
成功创建数据

> 查询数据

```shell
curl --location 'http://localhost:8080/api/v1/app'
```
返回内容

```json
{
    "code": 0,
    "data": {
        "total": 1,
        "size": 10,
        "current": 1,
        "results": [
            {
                "id": 1,
                "createdAt": "2023-04-13 16:35:59",
                "updatedAt": "2023-04-13 16:35:59",
                "name": "测试应用",
                "description": "测试应用服务",
                "level": "S3",
                "type": "container"
            }
        ]
    },
    "message": "success"
}
```
> 查询单个数据
```shell
curl --location 'http://localhost:8080/api/v1/app/1'
```
返回内容
```json
{
    "code": 0,
    "data": {
        "id": 1,
        "createdAt": "2023-04-13 16:56:09",
        "updatedAt": "2023-04-13 16:58:29",
        "name": "测试应用",
        "description": "测试应用服务",
        "level": "S3",
        "type": "container"
    },
    "message": "success"
}
```

> 更新数据

```shell
curl --location --request PUT 'http://localhost:8080/api/v1/app/1' \
--header 'Content-Type: application/json' \
--data '{
    "name": "测试应用",
    "description": "测试应用服务",
    "level": "S1",
    "type": "container"
}'
```
返回内容
```json
{
    "code": 0,
    "data": true,
    "message": "success"
}
```
> 删除数据

```shell
curl --location --request DELETE 'http://localhost:8080/api/v1/app/1'
```
返回内容
```json
{
    "code": 0,
    "data": true,
    "message": "success"
}
```
 ***自定义方法***
 
api添加新的方法
```go
// Package app api.go
package app

func (a Api) Hello(c *gin.Context) {
	response.OkWithData("Hello World", c)
}
```
router进行绑定
```go
    r.BindGet("hello", a.Hello)
```
接口请求

```shell
curl --location 'http://localhost:8080/api/v1/app/hello'
```
返回内容
```json
{
    "code": 0,
    "data": "Hello World",
    "message": "success"
}
```

### Service

在 NewApi 时使用的是BaseService，可以实现自定义的Service，重写一些方法。

```go
package app

import (
	"github.com/gin-gonic/gin"
	"github.com/songcser/gingo/pkg/model"
	"github.com/songcser/gingo/pkg/service"
	"github.com/songcser/gingo/utils"
)

type Service struct {
	service.Service[App]
}

func NewService(a App) Service {
	return Service{service.NewBaseService[App](a)}
}

func (s Service) MakeMapper(c *gin.Context) model.Mapper[App] {
	var r Request
	err := c.ShouldBindQuery(&r)
	utils.CheckError(err)
	w := model.NewWrapper()
	w.Like("name", r.Name)
	w.Eq("level", r.Level)
	m := model.NewMapper[App](App{}, w)
	return m
}

func (s Service) MakeResponse(val model.Model) any {
	a := val.(App)
	res := Response{
		Name:        a.Name,
		Description: fmt.Sprintf("名称：%s, 等级: %s, 类型: %s", a.Name, a.Level, a.Type),
		Level:       a.Level,
		Type:        a.Type,
	}
	return res
}

```
* MakeMapper方法重写新的Mapper，可以自定义一些查询条件。

* MakeResponse方法可以重写返回的内容

在Api中替换新的Service

```go
// Package app api.go

func NewApi() Api {
	var app App
	s := NewService(app)  //使用新的Service
	baseApi := api.NewApi[App](s)
	return Api{Api: baseApi}
}
```
> 查询数据

```shell
curl --location 'http://localhost:8080/api/v1/app?name=测试&level=S3'
```
返回内容
```json
{
    "code": 0,
    "data": {
        "total": 1,
        "size": 10,
        "current": 1,
        "results": [
            {
                "name": "测试应用",
                "description": "名称：测试应用, 等级: S3, 类型: container",
                "level": "S3",
                "type": "container"
            }
        ]
    },
    "message": "success"
}
```

如果要在Service增加新的方法，需要在Api模型中重写Service

```go
// Package app service.go
func (s Service) Hello() string {
	return "Hello World"
}
```
Api实现

```go
// Package app api.go
package app

type Api struct {
	api.Api
	Service Service  // 重写Service
}

func NewApi() Api {
	var app App
	s := NewService(app)
	baseApi := api.NewApi[App](s)
	return Api{Api: baseApi, Service: s}
}

func (a Api) Hello(c *gin.Context) {
	str := a.Service.Hello() // 调用Service方法
	response.OkWithData(str, c)
}

```

## Admin

Admin 提供了简单的后台管理服务，可以很方便的对数据进行管理。

首先需要在配置文件中开启 Admin

```yaml
admin:
  enable: true
  auth: true
```

auth 会开启认证授权，需要登陆

接入 Admin 也比较简单
```go
// Package app admin.go
package app

import (
	"github.com/songcser/gingo/pkg/admin"
)

func Admin() {
	var a App
	admin.New(a, "app", "应用")
}
```

在 initialize 中引入 admin

```go
// Package initialize admin.go
package initialize

func Admin(r *gin.Engine) {
	if !config.GVA_CONFIG.Admin.Enable {
		return
	}
	admin.Init(r, nil)
	app.Admin()
}
```
在管理页面可以进行简单的查询，创建，修改，删除操作。

Model 字段中可以配置 admin 标签，管理页面可以根据类型展示。

```go
// Package app model.go
package app

import "github.com/songcser/gingo/pkg/model"

type App struct {
	model.BaseModel
	Name        string `json:"name" form:"name" gorm:"column:name;type:varchar(255);not null" admin:"type:input;name:name;label:应用名"`
	Description string `json:"description" form:"description" gorm:"column:description;type:varchar(4096);not null" admin:"type:textarea;name:description;label:描述"`
	Level       string `json:"level" form:"level" gorm:"column:level;type:varchar(8);not null" admin:"type:radio;enum:S1,S2,S3,S4,S5;label:级别"`
	Type        string `json:"type" form:"type" gorm:"column:type;type:varchar(16);not null" admin:"type:select;enum:container=容器应用,web=前端应用,mini=小程序应用;label:应用类型"`
}
```

* type: 字段类型，取值 input, textarea, select, radio。
* name: 字段名称
* label: 字段展示名称
* enum: 枚举数据，在select，radio类型时展示更加友好

**_需要添加 form 标签，由于是使用 html 的 form 提交数据。_**

访问链接

```
http://localhost:8080/admin/
```
> 首页

![首页](https://raw.githubusercontent.com/songcser/gingo/master/docs/img/img.png)

> 点击应用，查看应用列表

![列表](https://raw.githubusercontent.com/songcser/gingo/master/docs/img/list.png)

> 点击 创建应用 按钮

![新建](https://raw.githubusercontent.com/songcser/gingo/master/docs/img/add.png)

