package log

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// global variable Logger to avoid log injections everywhere
var (
	logger       *zap.Logger
	levelChanger zap.AtomicLevel
)

const (
	CorrelationId = "X-Request-Id"
)

func init() {
	atom := zap.NewAtomicLevelAt(zap.DebugLevel)
	l := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()),
		zapcore.Lock(os.Stdout),
		atom,
	))
	SetupLogger(l, atom)
}

func SetupLogger(l *zap.Logger, _levelChanger zap.AtomicLevel) {
	logger = l.WithOptions(zap.AddCallerSkip(1))
	levelChanger = _levelChanger
}

func GetLogger() *zap.Logger {
	return logger
}

func SetLogLevel(level zapcore.Level) {
	levelChanger.SetLevel(level)
}

func GetLogLevel() zapcore.Level {
	return levelChanger.Level()
}

func Named(s string) *zap.Logger {
	return logger.Named(s)
}

func WithOptions(opts ...zap.Option) *zap.Logger {
	return logger.WithOptions(opts...)
}

func With(fields ...zap.Field) *zap.Logger {
	return logger.With(fields...)
}

func WithContext(ctx context.Context, fields ...zap.Field) *zap.Logger {
	log := logger.With(fields...)

	// add CorrelationId
	corId, ok := ctx.Value(CorrelationId).(int)
	if ok {
		log = log.With(zap.Field{Key: "CorId", Type: zapcore.Int64Type, Integer: int64(corId)})
	}
	return log
}

func Check(lvl zapcore.Level, msg string) *zapcore.CheckedEntry {
	return logger.Check(lvl, msg)
}

func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

func DPanic(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

func Panic(msg string, fields ...zap.Field) {
	logger.Panic(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	logger.Fatal(msg, fields...)
}

func Sync() error {
	return logger.Sync()
}

func Core() zapcore.Core {
	return logger.Core()
}
