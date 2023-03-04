package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// GetZapLogger returns a zap.Logger
func GetZapLogger(Debug bool) *zap.Logger {

	return getZapLogger(Debug, true, "log")
}

// getZapLogger returns a zap.Logger
func getZapLogger(Debug bool, console bool, label string) *zap.Logger {

	// Zap Logger
	var logger *zap.Logger
	var core zapcore.Core

	// override time format
	zapConfig := zap.NewProductionEncoderConfig()

	zapConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	var encoder zapcore.Encoder
	if console {
		encoder = zapcore.NewConsoleEncoder(zapConfig)
	} else {
		encoder = zapcore.NewJSONEncoder(zapConfig)
	}

	// First, define our level-handling logic.
	errorPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	infoPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl > zapcore.DebugLevel && lvl < zapcore.ErrorLevel
	})
	debugPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.DebugLevel
	})
	// default writer for logger
	// consoleStdout := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)

	if Debug {
		// set log level to writer
		core = zapcore.NewTee(
			zapcore.NewCore(encoder, consoleErrors, debugPriority),
			zapcore.NewCore(encoder, consoleErrors, infoPriority),
			zapcore.NewCore(encoder, consoleErrors, errorPriority),
		)
	} else {

		// set log level to writer
		core = zapcore.NewTee(
			zapcore.NewCore(encoder, consoleErrors, infoPriority),
			zapcore.NewCore(encoder, consoleErrors, errorPriority),
		)

	}

	// add function caller and stack trace on error
	logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return logger

}
