package logger

import (
	"fmt"
	"sync/atomic"
	"unsafe"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	_log unsafe.Pointer // Pointer to a coreLogger. Access via atomic.LoadPointer.
)

type LogOption = zap.Option
type coreLogger struct {
	logger       *Logger     // Logger that is the basis for all logp.Loggers.
	rootLogger   *zap.Logger // Root logger without any options configured.
	webLogger    *Logger
	globalLogger *zap.Logger
	atom         zap.AtomicLevel // 动态Level设置
}

type Logger struct {
	logger *zap.Logger
	sugar  *zap.SugaredLogger
}

func storeLogger(l *coreLogger) {
	if old := loadLogger(); old != nil {
		old.rootLogger.Sync()
	}
	atomic.StorePointer(&_log, unsafe.Pointer(l))
}
func newLogger(rootLogger *zap.Logger, selector string, options ...LogOption) *Logger {
	log := rootLogger.
		WithOptions().
		WithOptions(options...).
		Named(selector)
	return &Logger{log, log.Sugar()}
}

func newGinLogger(rootLogger *zap.Logger, selector string, options ...LogOption) *Logger {
	log := rootLogger.
		WithOptions().
		WithOptions(options...).
		Named(selector)
	return &Logger{log, log.Sugar()}
}

func NewLogger(e *Encoder) error {
	atom := zap.NewAtomicLevel()
	logger, webLogger := e.EncoderConfig()
	storeLogger(&coreLogger{
		rootLogger:   logger,
		logger:       newLogger(logger, ""),
		globalLogger: logger.WithOptions(),
		webLogger:    newGinLogger(webLogger, ""),
		atom:         atom,
	})
	return nil
}

// Named returns a logger that adds a new path segment to the logger's name.
func (l *Logger) Named(name string) *Logger {
	logger := l.logger.Named(name)
	return &Logger{logger, logger.Sugar()}
}

// SetLevel dynamically sets the logging level.
func (l *Logger) SetLevel(level string) {
	var zapLevel zap.AtomicLevel
	zapLevel.UnmarshalText([]byte(level))
	l.logger.Core().Enabled(zapLevel.Level())
}

type Option func(*zap.Config)

// WithCaller enables the caller field in the log output.
func WithCaller(caller bool) Option {
	return func(config *zap.Config) {
		config.Development = !caller
		config.DisableCaller = !caller
	}
}

// Print uses fmt.Sprint to construct and log a message.
func (l *Logger) Print(args ...interface{}) {
	l.sugar.Debug(args...)
}

// Println todo
func (l *Logger) Println(args ...interface{}) {
	l.sugar.Debug(args...)
}

// Debug uses fmt.Sprint to construct and log a message.
func (l *Logger) Debug(args ...interface{}) {
	l.sugar.Debug(args...)
}

// Info uses fmt.Sprint to construct and log a message.
func (l *Logger) Info(args ...interface{}) {
	l.sugar.Info(args...)
}

// Warn uses fmt.Sprint to construct and log a message.
func (l *Logger) Warn(args ...interface{}) {
	l.sugar.Warn(args...)
}

// Error uses fmt.Sprint to construct and log a message.
func (l *Logger) Error(args ...interface{}) {
	l.sugar.Error(args...)
}

// Fatal uses fmt.Sprint to construct and log a message, then calls os.Exit(1).
func (l *Logger) Fatal(args ...interface{}) {
	l.sugar.Fatal(args...)
}

// Panic uses fmt.Sprint to construct and log a message, then panics.
func (l *Logger) Panic(args ...interface{}) {
	l.sugar.Panic(args...)
}

// DPanic uses fmt.Sprint to construct and log a message. In development, the
// logger then panics.
func (l *Logger) DPanic(args ...interface{}) {
	l.sugar.DPanic(args...)
}

// IsDebug checks to see if the given logger is Debug enabled.
func (l *Logger) IsDebug() bool {
	return l.logger.Check(zapcore.DebugLevel, "") != nil
}

// Printf todo
func (l *Logger) Printf(format string, args ...interface{}) {
	l.sugar.Debugf(format, args...)
}

// Debugf uses fmt.Sprintf to construct and log a message.
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.sugar.Debugf(format, args...)
}

// Infof uses fmt.Sprintf to log a templated message.
func (l *Logger) Infof(format string, args ...interface{}) {
	l.sugar.Infof(format, args...)
}

// Warnf uses fmt.Sprintf to log a templated message.
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.sugar.Warnf(format, args...)
}

// Errorf uses fmt.Sprintf to log a templated message.
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.sugar.Errorf(format, args...)
}

// Fatalf uses fmt.Sprintf to log a templated message, then calls os.Exit(1).
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.sugar.Fatalf(format, args...)
}

// Panicf uses fmt.Sprintf to log a templated message, then panics.
func (l *Logger) Panicf(format string, args ...interface{}) {
	l.sugar.Panicf(format, args...)
}

// DPanicf uses fmt.Sprintf to log a templated message. In development, the
// logger then panics.
func (l *Logger) DPanicf(format string, args ...interface{}) {
	l.sugar.DPanicf(format, args...)
}

// With context (reflection based)

// Debugw logs a message with some additional context. The additional context
// is added in the form of key-value pairs. The optimal way to write the value
// to the log message will be inferred by the value's type. To explicitly
// specify a type you can pass a Field such as logp.Stringer.
func (l *Logger) Debugw(msg string, fields ...Field) {
	l.sugar.Debugw(msg, transfer(fields)...)
}

// Infow logs a message with some additional context. The additional context
// is added in the form of key-value pairs. The optimal way to write the value
// to the log message will be inferred by the value's type. To explicitly
// specify a type you can pass a Field such as logp.Stringer.
func (l *Logger) Infow(msg string, fields ...Field) {
	l.sugar.Infow(msg, transfer(fields)...)
}

// Warnw logs a message with some additional context. The additional context
// is added in the form of key-value pairs. The optimal way to write the value
// to the log message will be inferred by the value's type. To explicitly
// specify a type you can pass a Field such as logp.Stringer.
func (l *Logger) Warnw(msg string, fields ...Field) {
	l.sugar.Warnw(msg, transfer(fields)...)
}

// Errorw logs a message with some additional context. The additional context
// is added in the form of key-value pairs. The optimal way to write the value
// to the log message will be inferred by the value's type. To explicitly
// specify a type you can pass a Field such as logp.Stringer.
func (l *Logger) Errorw(msg string, fields ...Field) {
	l.sugar.Errorw(msg, transfer(fields)...)
}

// Fatalw logs a message with some additional context, then calls os.Exit(1).
// The additional context is added in the form of key-value pairs. The optimal
// way to write the value to the log message will be inferred by the value's
// type. To explicitly specify a type you can pass a Field such as
// logp.Stringer.
func (l *Logger) Fatalw(msg string, fields ...Field) {
	l.sugar.Fatalw(msg, transfer(fields)...)
}

// Panicw logs a message with some additional context, then panics. The
// additional context is added in the form of key-value pairs. The optimal way
// to write the value to the log message will be inferred by the value's type.
// To explicitly specify a type you can pass a Field such as logp.Stringer.
func (l *Logger) Panicw(msg string, fields ...Field) {
	l.sugar.Panicw(msg, transfer(fields)...)
}

// DPanicw logs a message with some additional context. The logger panics only
// in Development mode.  The additional context is added in the form of
// key-value pairs. The optimal way to write the value to the log message will
// be inferred by the value's type. To explicitly specify a type you can pass a
// Field such as logp.Stringer.
func (l *Logger) DPanicw(msg string, fields ...Field) {
	l.sugar.DPanicw(msg, transfer(fields)...)
}

type Field struct {
	Key   string
	Value interface{}
}

// transfer 用于方便zap记录
func transfer(m []Field) (ma []interface{}) {
	for i := range m {
		ma = append(ma, zap.Any(m[i].Key, m[i].Value))
	}

	return
}
func globalLogger() *zap.Logger {
	return loadLogger().globalLogger
}

// Recover stops a panicking goroutine and logs an Error.
func (l *Logger) Recover(msg string) {
	if r := recover(); r != nil {
		msg := fmt.Sprintf("%s. Recovering, but please report this.", msg)
		globalLogger().WithOptions().
			Error(msg, zap.Any("panic", r), zap.Stack("stack"))
	}
}
func loadLogger() *coreLogger {
	p := atomic.LoadPointer(&_log)
	return (*coreLogger)(p)
}

// SetLevel 动态设置日志等级
func SetLevel(lv Level) {
	loadLogger().atom.SetLevel(lv.zapLevel())
}

func L() *Logger {
	return loadLogger().logger
}

func W() *Logger {
	return loadLogger().webLogger
}
