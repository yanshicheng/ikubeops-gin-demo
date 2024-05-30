package settings

import (
	"fmt"
	"github.com/yanshicheng/ikubeops-gin-demo/global"
	logger "github.com/yanshicheng/ikubeops-gin-demo/settings/ikube-logger"
	"github.com/yanshicheng/ikubeops-gin-demo/utils"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// 初始化Logger

// InitLogger 初始化Logger
func LoadGlobalLogger() (*logger.Logger, error) {
	config := logger.Encoder(global.C.Logger)
	err := logger.NewLogger(&config)
	return logger.L(), err
}
func LoadLogger() (logger *zap.Logger, err error) {
	if ok, _ := utils.PathExists(global.C.Logger.FilePath); !ok { // 判断是否有Director文件夹
		fmt.Printf("create %v directory\n", global.C.Logger.FilePath)
		_ = os.Mkdir(global.C.Logger.FilePath, os.ModePerm)
	}

	// 调试级别
	debugPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev == zap.DebugLevel
	})
	// 日志级别
	infoPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev == zap.InfoLevel
	})
	// 警告级别
	warnPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev == zap.WarnLevel
	})
	// 错误级别
	errorPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev >= zap.ErrorLevel
	})
	//encoder := getEncoder()
	cores := [...]zapcore.Core{
		getEncoderCore(fmt.Sprintf("./%s/server_debug.log", global.C.Logger.FilePath), debugPriority),
		getEncoderCore(fmt.Sprintf("./%s/server_info.log", global.C.Logger.FilePath), infoPriority),
		getEncoderCore(fmt.Sprintf("./%s/server_warn.log", global.C.Logger.FilePath), warnPriority),
		getEncoderCore(fmt.Sprintf("./%s/server_error.log", global.C.Logger.FilePath), errorPriority),
	}
	var core zapcore.Core
	if global.C.Logger.Output == "console" {
		consoleEnbcoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		core = zapcore.NewTee(
			zapcore.NewCore(consoleEnbcoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel),
			zapcore.NewTee(cores[:]...),
		)
	} else {
		core = zapcore.NewTee(
			zapcore.NewTee(cores[:]...),
		)
	}
	logger = zap.New(core, zap.AddCaller())
	// 显示行
	logger = logger.WithOptions(zap.AddCaller())
	zap.ReplaceGlobals(logger) // 替换zap包中全局的logger实例，后续在其他包中只需使用zap.L()调用即可
	return logger, nil
}

// getEncoderCore 获取Encoder的zapcore.Core
func getEncoderCore(fileName string, level zapcore.LevelEnabler) (core zapcore.Core) {
	writer := GetWriteSyncer(fileName) // 使用lumberjack进行日志分割
	return zapcore.NewCore(getEncoder(), writer, level)
}

// 配置日志切割
func GetWriteSyncer(file string) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   file,                       //日志文件的位置
		MaxSize:    global.C.Logger.MaxSize,    //在进行切割之前，日志文件的最大大小（以MB为单位）
		MaxBackups: global.C.Logger.MaxBackups, //保留旧文件的最大个数
		MaxAge:     global.C.Logger.MaxAge,     //保留旧文件的最大天数
		Compress:   true,                       //是否压缩/归档旧文件
	}

	return zapcore.AddSync(lumberJackLogger)
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

// GinLogger 接收gin框架默认的日志
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()
		//如果路径是 /healthz，则不记录日志
		if path == "/healthz" {
			return
		}
		cost := time.Since(start)
		global.L.Named("gin-web").Info(
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
		)
	}
}
