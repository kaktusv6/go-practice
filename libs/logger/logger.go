package logger

import (
	"net/http"
)

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var globalLogger *zap.Logger

func Init(config Config) {
	globalLogger = New(config)
}

type Config struct {
	Level string
	Env   string
}

func isProd(env string) bool {
	return env == "production"
}

func getLevelByKey(key string) zapcore.Level {
	var result zapcore.Level

	switch key {
	case "debug":
		result = zap.DebugLevel
	case "info":
		result = zap.InfoLevel
	case "warn":
		result = zap.WarnLevel
	case "error":
		result = zap.ErrorLevel
	case "fatal":
		result = zap.FatalLevel
	}

	return result
}

func New(config Config) *zap.Logger {
	var logger *zap.Logger
	var err error

	if isProd(config.Env) {
		cfg := zap.NewProductionConfig()
		cfg.DisableCaller = true
		cfg.DisableStacktrace = true
		cfg.Level = zap.NewAtomicLevelAt(getLevelByKey(config.Level))
		logger, err = cfg.Build()
	} else {
		logger, err = zap.NewDevelopment()
	}

	if err != nil {
		panic(err)
	}

	return logger
}

func Middleware(logger *zap.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Debug(
			"incoming http request",
			zap.String("path", r.URL.Path),
			zap.String("query", r.URL.RawQuery),
		)

		next.ServeHTTP(w, r)
	})
}

func Debug(msg string, fields ...zap.Field) {
	globalLogger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	globalLogger.Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	globalLogger.Error(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	globalLogger.Warn(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	globalLogger.Fatal(msg, fields...)
}
