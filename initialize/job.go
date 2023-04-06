package initialize

import (
	"github.com/robfig/cron/v3"
	"time"
)

func Job() *cron.Cron {
	return cron.New(cron.WithLocation(time.Local))
}
