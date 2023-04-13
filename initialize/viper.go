package initialize

import (
	"flag"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/songcser/gingo/config"
	"github.com/spf13/viper"
	"os"
)

// Viper //
// 优先级: 命令行 > 环境变量 > 默认值
func Viper(path ...string) *viper.Viper {
	var cfg string

	if len(path) == 0 {
		flag.StringVar(&cfg, "c", "", "choose config file.")
		flag.Parse()
		if cfg == "" {
			/*
			   判断 internal.ConfigEnv 常量存储的环境变量是否为空
			   比如我们启动项目的时候，执行：GVA_CONFIG=config.yaml go run main.go
			   这时候 os.Getenv(internal.ConfigEnv) 得到的就是 config.yaml
			   当然，也可以通过 os.Setenv(internal.ConfigEnv, "config.yaml") 在初始化之前设置
			*/
			if configEnv := os.Getenv(ConfigEnv); configEnv == "" {
				switch gin.Mode() {
				case gin.DebugMode:
					cfg = ConfigDefaultFile
					fmt.Printf("正在使用gin模式的%s环境名称,config的路径为%s\n", gin.EnvGinMode, ConfigDefaultFile)
				case gin.ReleaseMode:
					cfg = ConfigReleaseFile
					fmt.Printf("正在使用gin模式的%s环境名称,config的路径为%s\n", gin.EnvGinMode, ConfigReleaseFile)
				case gin.TestMode:
					cfg = ConfigTestFile
					fmt.Printf("正在使用gin模式的%s环境名称,config的路径为%s\n", gin.EnvGinMode, ConfigTestFile)
				}
			} else {
				// internal.ConfigEnv 常量存储的环境变量不为空 将值赋值于config
				cfg = configEnv
				fmt.Printf("正在使用%s环境变量,config的路径为%s\n", ConfigEnv, cfg)
			}
		} else {
			cfg = "config/config.yaml"
			fmt.Printf("正在使用命令行的-c参数传递的值,config的路径为%s\n", cfg)
		}
	} else {
		cfg = path[0]
	}

	v := viper.New()
	v.SetConfigFile(cfg)
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		if err = v.Unmarshal(&config.GVA_CONFIG); err != nil {
			fmt.Println(err)
		}
	})
	if err = v.Unmarshal(&config.GVA_CONFIG); err != nil {
		fmt.Println(err)
	}

	return v
}
