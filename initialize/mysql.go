package initialize

import (
	"fmt"
	"github.com/songcser/gingo/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

type writer struct {
	logger.Writer
}

// NewWriter writer 构造函数
func NewWriter(w logger.Writer) *writer {
	return &writer{Writer: w}
}

// Printf 格式化打印日志
func (w *writer) Printf(message string, data ...interface{}) {
	var logZap bool
	switch config.GVA_CONFIG.DbType {
	case "mysql":
		logZap = config.GVA_CONFIG.Mysql.LogZap
	}
	if logZap {
		config.GVA_LOG.Info(fmt.Sprintf(message+"\n", data...))
	} else {
		w.Writer.Printf(message, data...)
	}
}

type DBBASE interface {
	GetLogMode() string
}

var orm = new(_gorm)

type _gorm struct{}

// Config gorm 自定义配置
func (g *_gorm) Config(prefix string, singular bool) *gorm.Config {
	cfg := &gorm.Config{
		// 命名策略
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   prefix,   // 表前缀，在表名前添加前缀，如添加用户模块的表前缀 user_
			SingularTable: singular, // 是否使用单数形式的表名，如果设置为 true，那么 User 模型会对应 users 表
		},

		DisableForeignKeyConstraintWhenMigrating: true,
	}
	_default := logger.New(NewWriter(log.New(os.Stdout, "\r\n", log.LstdFlags)), logger.Config{
		SlowThreshold: 200 * time.Millisecond,
		LogLevel:      logger.Warn,
		Colorful:      true,
	})

	var logMode DBBASE
	switch config.GVA_CONFIG.DbType {
	case "mysql":
		logMode = &config.GVA_CONFIG.Mysql
	default:
		logMode = &config.GVA_CONFIG.Mysql
	}

	switch logMode.GetLogMode() {
	case "silent", "Silent":
		cfg.Logger = _default.LogMode(logger.Silent)
	case "error", "Message":
		cfg.Logger = _default.LogMode(logger.Error)
	case "warn", "Warn":
		cfg.Logger = _default.LogMode(logger.Warn)
	case "info", "Info":
		cfg.Logger = _default.LogMode(logger.Info)
	default:
		cfg.Logger = _default.LogMode(logger.Info)
	}
	return cfg

}

func GormMysql() *gorm.DB {
	m := config.GVA_CONFIG.Mysql
	if m.Dbname == "" {
		return nil
	}
	fmt.Println(m.Dsn())
	mysqlConfig := mysql.Config{
		DSN:                       m.Dsn(), // DSN data source name
		DefaultStringSize:         191,     // string 类型字段的默认长度
		SkipInitializeWithVersion: false,   // 根据版本自动配置
	}
	if db, err := gorm.Open(mysql.New(mysqlConfig), orm.Config(m.Prefix, m.Singular)); err != nil {
		return nil
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(m.MaxIdleConns)
		sqlDB.SetMaxOpenConns(m.MaxOpenConns)
		return db
	}
}
