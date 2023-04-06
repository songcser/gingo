package main

import (
	"github.com/songcser/gingo/config"
	"github.com/songcser/gingo/initialize"
)

func main() {
	config.GVA_VP = initialize.Viper()
	config.GVA_DB = initialize.Gorm() // gorm连接数据库
	if config.GVA_DB != nil {
		initialize.RegisterTables(config.GVA_DB) // 初始化表
		// 程序结束前关闭数据库链接
		db, _ := config.GVA_DB.DB()
		defer db.Close()
	}
}
