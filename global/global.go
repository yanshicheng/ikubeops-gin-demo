package global

import (
	ut "github.com/go-playground/universal-translator"
	logger "github.com/yanshicheng/ikubeops-gin-demo/settings/ikube-logger"
	"gorm.io/gorm"
)

const (
	AppName                   = ""
	EnvPrefix                 = "IKUBEOPS_"
	ConfigModeFile ConfigMode = "file"
	ConfigModeEnv  ConfigMode = "env"
)

type ConfigMode string

var (
	// Config 全局配置
	IkubeopsTrans ut.Translator
	C             *Config = NewDefaultConfig()
	L             *logger.Logger
	LSys          *logger.Logger
	DB            *gorm.DB
	M             []interface{}
	// Log 全局日志
)
