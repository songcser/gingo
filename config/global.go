package config

import (
	"context"
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Context struct {
	Ctx    context.Context
	Cancel context.CancelFunc
}

var (
	GVA_CONFIG Configuration
	GVA_DB     *gorm.DB
	GVA_LOG    *zap.Logger
	GVA_VP     *viper.Viper

	GVA_CTX *Context
	GVA_JOB *cron.Cron
)

func NewContext() *Context {
	ctx, cancel := context.WithCancel(context.Background())
	GVA_CTX = &Context{
		Ctx:    ctx,
		Cancel: cancel,
	}
	return GVA_CTX
}
