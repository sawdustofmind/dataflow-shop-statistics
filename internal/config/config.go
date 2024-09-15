package config

import (
	"strings"
	"time"

	"github.com/iamolegga/enviper"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	UserServer *ServerConfig `mapstructure:"user_server"`
	Logger     *LoggerConf   `mapstructure:"logger"`
	Storage    *StorageConf  `mapstructure:"storage"`
}

type ServerConfig struct {
	GinMode        string        `mapstructure:"gin_mode"`
	Address        string        `mapstructure:"address"`
	ReadTimeout    time.Duration `mapstructure:"read_timeout"`
	WriteTimeout   time.Duration `mapstructure:"write_timeout"`
	MaxHeaderBytes int           `mapstructure:"max_header_bytes"`
}

type LoggerConf struct {
	Mode string `mapstructure:"mode"`
}

type StorageConf struct {
	Type            string `mapstructure:"type"`
	InitialCapacity int    `mapstructure:"initial_capacity"`
}

var Flags = &pflag.FlagSet{}

func init() {
	Flags.String("conf", "config", "Config file")
	Flags.String("confpath", "./configs", "Config path")
}

func Init() (*Config, error) {
	cfgName, err := Flags.GetString("conf")
	if err != nil {
		cfgName = "config"
	}

	cfgPath, err := Flags.GetString("confpath")
	if err != nil {
		cfgPath = "./configs"
	}

	return InitConfig(cfgName, cfgPath)
}

func InitConfig(cfgName, cfgPath string) (*Config, error) {
	vpr := viper.New()

	vpr.SetConfigName(cfgName)
	vpr.AddConfigPath(cfgPath)

	replacer := strings.NewReplacer(".", "_")
	vpr.SetEnvKeyReplacer(replacer)
	vpr.AutomaticEnv()

	if err := vpr.ReadInConfig(); err != nil {
		return nil, err
	}

	conf := &Config{}
	envAwareViper := enviper.New(vpr)
	if err := envAwareViper.Unmarshal(conf); err != nil {
		return nil, err
	}

	return conf, nil
}
