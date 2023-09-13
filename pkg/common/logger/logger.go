package logger

import (
	"log"

	config "backend-service/config/marketplace"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var singleton *zap.Logger

// Init initializes a thread-safe singleton logger
// This would be called from a main method when the application starts up
// This function would ideally, take zap configuration, but is left out
// in favor of simplicity using the example logger.
func Init(cfg *config.Config) *zap.Logger {
	if cfg.Log.Prod || cfg.Env == config.EnvProd {
		singleton, _ = zap.NewProduction()

		return singleton
	}

	// once ensures the singleton is initialized only once
	zaplv, err := zapcore.ParseLevel(cfg.Log.Level)
	if err != nil {
		log.Fatalf("Cannot init logger")
	}

	zapcfg := zap.Config{
		Encoding:    "console",
		Level:       zap.NewAtomicLevelAt(zaplv),
		OutputPaths: []string{"stderr"},

		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:    "message",
			StacktraceKey: "stacktrace",
			TimeKey:       "time",
			LevelKey:      "level",
			CallerKey:     "caller",
			FunctionKey:   zapcore.OmitKey,
			EncodeLevel:   zapcore.CapitalLevelEncoder,
			EncodeCaller:  zapcore.FullCallerEncoder,
			EncodeTime:    zapcore.RFC3339TimeEncoder,
		},
	}

	logger, err := zapcfg.Build()
	if err != nil {
		log.Fatalf("Cannot init logger")
	}

	singleton = logger

	return singleton
}

// Debug logs a debug message with the given fields.
func Debug(message string, fields ...zap.Field) {
	singleton.Debug(message, fields...)
}

// Info logs a debug message with the given fields.
func Info(message string, fields ...zap.Field) {
	singleton.Info(message, fields...)
}

// Info logs a debug message with the given fields.
func Infow(message string, keysAndValues ...any) {
	singleton.Sugar().Infow(message, keysAndValues...)
}

// Warn logs a debug message with the given fields.
func Warn(message string, fields ...zap.Field) {
	singleton.Warn(message, fields...)
}

// Warnw logs a debug message with the given fields.
func Warnw(message string, keysAndValues ...any) {
	singleton.Sugar().Warnw(message, keysAndValues...)
}

// Error logs a debug message with the given fields.
func Error(message string, fields ...zap.Field) {
	singleton.Error(message, fields...)
}

// Errorw logs a debug message with the given key-value pairs.
func Errorw(message string, keysAndValues ...any) {
	singleton.Sugar().Errorw(message, keysAndValues...)
}

// Fatal logs a message.
func Fatal(message string, keysAndValues ...any) {
	singleton.Sugar().Fatalw(message, keysAndValues...)
}

func ErrWrap(err error) zap.Field {
	return zap.Error(err)
}
