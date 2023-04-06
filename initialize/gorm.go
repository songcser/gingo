package initialize

import (
	"github.com/songcser/gingo/config"
	"github.com/songcser/gingo/pkg/auth"
	"gorm.io/gorm"
	"os"
)

func Gorm() *gorm.DB {
	switch config.GVA_CONFIG.DbType {
	case "mysql":
		return GormMysql()
	default:
		return GormMysql()
	}
}

// RegisterTables 注册数据库表专用
// Author SliverHorn
func RegisterTables(db *gorm.DB) {
	err := db.AutoMigrate(
		// 系统模块表
		auth.BaseUser{},
	)
	if err != nil {
		os.Exit(0)
	}
}
