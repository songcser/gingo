package autoload

type Admin struct {
	Enable bool `mapstructure:"enable" json:"enable" yaml:"enable"`
	Auth   bool `mapstructure:"auth" json:"auth" yaml:"auth"`
}
