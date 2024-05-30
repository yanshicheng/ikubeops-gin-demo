package logger

import (
	"os"
	"path/filepath"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	TimeLayout = "2006-01-02 15:04:05"
	LogName    = "app.log"
)

type Encoder struct {
	Output     string `json:"output" yaml:"output" mapstructure:"output"`
	Format     string `json:"format" yaml:"format" mapstructure:"format"`
	Level      string `json:"level" yaml:"level" mapstructure:"level"`
	Dev        bool   `json:"dev" yaml:"dev" mapstructure:"dev"`
	FilePath   string `json:"file_path" yaml:"file_path" mapstructure:"file_path"`
	MaxSize    int    `json:"max_size" yaml:"max_size" mapstructure:"max_size"`
	MaxAge     int    `json:"max_age" yaml:"max_age" mapstructure:"max_age"`
	MaxBackups int    `json:"max_backups" yaml:"max_backups" mapstructure:"max_backups"`
}

func (e *Encoder) EncoderConfig() (*zap.Logger, *zap.Logger) {
	return e.logger(), e.webLogger()
}

func (e *Encoder) logger() *zap.Logger {
	var encoder zapcore.Encoder
	if e.Format == "json" {
		encoder = zapcore.NewJSONEncoder(e.getEncoder())
	} else {
		encoder = zapcore.NewConsoleEncoder(e.getEncoder())
	}

	var fileCore zapcore.Core = zapcore.NewCore(
		encoder,
		e.getWriteSyncer(filepath.Join(getDirPath(e.FilePath), LogName)),
		getLogLevel(e.Level),
	)

	var consoleCore zapcore.Core = zapcore.NewCore(
		encoder, // 使用相同的编码器
		zapcore.Lock(os.Stdout),
		getLogLevel(e.Level),
	)

	var teeCore zapcore.Core
	if e.Output == "file" {
		teeCore = zapcore.NewTee(fileCore)
	} else {
		// 输出到文件和控制台
		teeCore = zapcore.NewTee(fileCore, consoleCore)
	}

	var logger *zap.Logger
	if e.Dev {
		logger = zap.New(teeCore, zap.AddCaller(), zap.AddCallerSkip(1))
	} else {
		logger = zap.New(teeCore)
	}
	return logger
}

func (e *Encoder) webLogger() *zap.Logger {
	var cores zapcore.Core
	var teeCore zapcore.Core
	if e.Format == "json" {
		zap.NewDevelopmentEncoderConfig()
		cores = zapcore.NewCore(
			zapcore.NewJSONEncoder(e.getEncoder()),
			e.getWriteSyncer(filepath.Join(getDirPath(e.FilePath), "web.log")),
			getLogLevel(e.Level),
		)
	} else {
		cores = zapcore.NewCore(
			zapcore.NewConsoleEncoder(e.getEncoder()),
			e.getWriteSyncer(filepath.Join(getDirPath(e.FilePath), "web.log")),
			getLogLevel(e.Level),
		)
	}
	if e.Output == "file" {
		teeCore = zapcore.NewTee(
			cores,
		)
	} else {
		teeCore = zapcore.NewTee(
			zapcore.NewCore(zapcore.NewConsoleEncoder(e.getEncoder()), zapcore.Lock(os.Stdout), getLogLevel(e.Level)),
			cores,
		)
	}
	var logger *zap.Logger
	if e.Dev {
		logger = zap.New(teeCore, zap.AddCaller(), zap.AddCallerSkip(1))

	} else {
		logger = zap.New(teeCore)

	}
	return logger
}

func (e *Encoder) getEncoder() zapcore.EncoderConfig {
	if !e.Dev {
		return zapcore.EncoderConfig{
			TimeKey:       "timestamp",
			LevelKey:      "level",
			NameKey:       "logger",
			CallerKey:     "caller",
			FunctionKey:   zapcore.OmitKey,
			MessageKey:    "message",
			StacktraceKey: "stacktrace",
			LineEnding:    zapcore.DefaultLineEnding,
			EncodeLevel:   zapcore.LowercaseLevelEncoder,
			EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
				enc.AppendString(t.Format(TimeLayout))
			},
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}
	} else {
		return zapcore.EncoderConfig{
			TimeKey:       "T",
			LevelKey:      "L",
			NameKey:       "N",
			CallerKey:     "C",
			FunctionKey:   zapcore.OmitKey,
			MessageKey:    "M",
			StacktraceKey: "S",
			LineEnding:    zapcore.DefaultLineEnding,
			EncodeLevel:   zapcore.CapitalLevelEncoder,
			EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
				enc.AppendString(t.Format(TimeLayout))
			},
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}
	}
}

func (e *Encoder) getWriteSyncer(file string) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   file,         //日志文件的位置
		MaxSize:    e.MaxSize,    //在进行切割之前，日志文件的最大大小（以MB为单位）
		MaxBackups: e.MaxBackups, //保留旧文件的最大个数
		MaxAge:     e.MaxAge,     //保留旧文件的最大天数
		Compress:   false,        //是否压缩/归档旧文件
	}
	// defer lumberJackLogger.Close()

	return zapcore.AddSync(lumberJackLogger)
}

func getLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	default:
		return zap.InfoLevel
	}
}

func getDirPath(filePath string) string {
	return filepath.Dir(filePath)
}
