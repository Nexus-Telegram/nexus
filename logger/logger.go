package logger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("02/01/2006 03:04 PM"))
}

func New() *zap.Logger {
	// Define log level (info by default)
	logLevel := zapcore.InfoLevel

	// Enable debug mode based on environment (optional)
	if os.Getenv("DEBUG") == "true" {
		logLevel = zapcore.DebugLevel
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeTime:    customTimeEncoder,
		EncodeLevel:   zapcore.CapitalLevelEncoder,
		EncodeCaller:  zapcore.ShortCallerEncoder,
	}

	// Create console and file writers
	consoleWriter := zapcore.Lock(os.Stdout)

	// Always log to a file (e.g., app.log)
	fileWriter, err := os.Create("app.log")
	if err != nil {
		panic("failed to create log file: " + err.Error())
	}

	// Create core with multiple outputs (console + file)
	core := zapcore.NewTee(
		// Console output (info level or debug if DEBUG=true)
		zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), consoleWriter, logLevel),

		// File output (always log everything to the file in JSON format)
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(fileWriter), zapcore.DebugLevel),
	)

	// Add options like caller info and stack traces
	options := []zap.Option{
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel), // Stack trace for errors and above
	}

	// Build and return logger
	return zap.New(core, options...)
}
