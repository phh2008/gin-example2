package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Config struct {
	Server Server `json:"server" yaml:"server"`
	Db     Db     `json:"db" yaml:"db"`
	Jwt    Jwt    `json:"jwt" yaml:"jwt"`
	Log    Log    `json:"log" yaml:"log"`
	Cors   Cors   `json:"cors" yaml:"cors"`
}

type Server struct {
	Port       string `yaml:"port" json:"port"`
	SignToken  string `yaml:"signToken" json:"signToken"`
	ExpireTime int64  `yaml:"expireTime" json:"expireTime"` // 签名有效期
}

type Db struct {
	Url string `yaml:"url" json:"url"`
}

type Jwt struct {
	Key string `yaml:"key" json:"key"`
}

type Log struct {
	Level      string `yaml:"level" json:"level"`
	Filename   string `yaml:"filename" json:"filename"`
	MaxSize    int    `yaml:"maxSize" json:"maxSize"`
	MaxBackups int    `yaml:"maxBackups" json:"maxBackups"`
	MaxAge     int    `yaml:"maxAge" json:"maxAge"`
	Compress   bool   `yaml:"compress" json:"compress"`
	LocalTime  bool   `yaml:"localTime" json:"localTime"`
}

type Cors struct {
	AllowedOriginPatterns []string `yaml:"allowedOriginPatterns" json:"allowedOriginPatterns"`
	AllowedMethods        string   `yaml:"allowedMethods" json:"allowedMethods"`
	AllowedHeaders        string   `yaml:"allowedHeaders" json:"allowedHeaders"`
	ExposeHeaders         string   `yaml:"exposeHeaders" json:"exposeHeaders"`
	MaxAge                int64    `yaml:"maxAge" json:"maxAge"`
	AllowCredentials      bool     `yaml:"allowCredentials" json:"allowCredentials"`
}

var Path string

// Active 当前环境，比如：dev，测会加载 config-dev.yml的配置
var Active string

func NewConfig(path string) *Config {
	var env string
	if Active != "" {
		env = "-" + Active
	}
	Path = path
	vp := viper.New()
	vp.SetConfigName("config" + env)
	vp.SetConfigType("yml")
	vp.AddConfigPath(Path)
	err := vp.ReadInConfig()
	if err != nil {
		zap.S().Errorf("加载配置错误,error:%s", err.Error())
		panic(err)
	}
	var conf Config
	if err = vp.Unmarshal(&conf); err != nil {
		zap.S().Errorf("绑定配置出错,error:%s", err.Error())
		panic(err)
	}
	vp.WatchConfig()
	vp.OnConfigChange(func(e fsnotify.Event) {
		zap.S().Infof("config file changed:%s", e.Name)
		if err = vp.Unmarshal(&conf); err != nil {
			zap.S().Errorf("更新配置出错,error:%s", err.Error())
		}
	})
	return &conf
}
