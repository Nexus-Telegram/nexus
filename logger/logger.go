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

func NewLogger(isProduction bool) (*zap.Logger, error) {
	var logger *zap.Logger
	var err error

	// Define encoder configuration
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = customTimeEncoder // Use custom time encoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	// Set log level based on environment
	var level zapcore.Level
	if isProduction {
		level = zapcore.InfoLevel
	} else {
		level = zapcore.DebugLevel
	}

	// Create core with JSON output for production, console for development
	var core zapcore.Core
	if isProduction {
		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(os.Stdout),
			level,
		)
	} else {
		core = zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.AddSync(os.Stdout),
			level,
		)
	}

	logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.WarnLevel))
	return logger, err
}
