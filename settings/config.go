package settings

import (
	"errors"
	"fmt"
	"github.com/yanshicheng/ikubeops-gin-demo/global"
	"github.com/yanshicheng/ikubeops-gin-demo/version"
	"os"

	"github.com/caarlos0/env/v8"
	"github.com/spf13/viper"
)

func LoadFileConfig(configFile string) error {
	// 检查文件是否存在
	fmt.Println(configFile)
	if _, err := os.Stat(configFile); err != nil {
		return fmt.Errorf("config file does not exist: %s", configFile)
	}
	viper.SetConfigFile(configFile)
	err := viper.ReadInConfig() // 读取配置信息
	if err != nil {             // 读取配置信息失败
		return fmt.Errorf("Fatal error config file: %s \n", err)
	}
	// 将读取的配置信息保存至全局变量Conf
	if err := viper.Unmarshal(&global.C); err != nil {
		return fmt.Errorf("unmarshal conf failed, err:%s \n", err)
	}
	if err := loadConfigFromEnv(); err != nil {
		return fmt.Errorf("load config from env failed, err:%s \n", err)
	}
	//// 监控配置文件变化
	//viper.WatchConfig()
	//// 注意！！！配置文件发生变化后要同步到全局变量Conf
	//viper.OnConfigChange(func(in fsnotify.Event) {
	//	fmt.Println("配置文件被修改啦,正在重载...")
	//	if err := viper.Unmarshal(&global.C); err != nil {
	//		panic(fmt.Errorf("unmarshal conf failed, err:%s \n", err))
	//	}
	//})
	return nil
}

func loadConfigFromEnv() error {
	opts := env.Options{
		Prefix: global.EnvPrefix, // 只加载以 T_ 开头的环境变量。
	}
	err := env.ParseWithOptions(global.C, opts)
	if err != nil {
		return fmt.Errorf("failed to load config from environment variables: %w", err)
	}
	return nil
}

func LoadGlobalConfig(configType, configFile string) error {
	// 配置加载
	version.IkubeopsConfigFile = configFile
	version.IkubeopsConfigType = configType
	switch configType {
	case "file":
		return LoadFileConfig(configFile)
	default:
		return errors.New("unknown config type")
	}
}
