package main

import (
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/songcser/gingo/config"
	"github.com/songcser/gingo/initialize"
	"go.uber.org/zap"
	"time"
)

type server interface {
	ListenAndServe() error
}

func RunWindowServer() {

	Router := initialize.Routers()
	address := ":8080"
	s := initServer(address, Router)

	if err := s.ListenAndServe(); err != nil {
	}
}

func initServer(address string, router *gin.Engine) server {
	s := endless.NewServer(address, router)
	s.ReadHeaderTimeout = 20 * time.Second
	s.WriteTimeout = 20 * time.Second
	s.MaxHeaderBytes = 1 << 20
	return s
}

func main() {
	// 初始化配置
	config.GVA_VP = initialize.Viper()

	// 初始化日志
	config.GVA_LOG = initialize.Zap()
	zap.ReplaceGlobals(config.GVA_LOG)

	config.GVA_DB = initialize.Gorm() // gorm连接数据库
	if config.GVA_DB != nil {
		initialize.RegisterTables(config.GVA_DB) // 初始化表
		// 程序结束前关闭数据库链接
		db, _ := config.GVA_DB.DB()
		defer db.Close()
	}

	if config.GVA_JOB != nil {
		defer config.GVA_JOB.Stop()
	}
	RunWindowServer()
}
