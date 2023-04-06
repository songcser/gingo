package config

type Configuration struct {
	Domain string `mapstructure:"domain" json:"domain" yaml:"domain"`
	DbType string `mapstructure:"dbType" json:"dbType" yaml:"dbType"`
	Mysql  Mysql  `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Zap    Zap    `mapstructure:"zap" json:"zap" yaml:"zap"`
	JWT    JWT    `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
}
