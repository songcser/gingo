# Gingo

## Introduce

Gingo是基于 gin 框架为核心的脚手架，使用本项目可以快速完成业务逻辑开发。

## Feature

* gin框架，简单，高效，轻量
* gorm数据库ORM框架，封装mapper，使用简单
* viper配置管理
* zap日志框架，输出日志更灵活
* api接口封装，快速实现CURD操作
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



### Api




