package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func Init() {
	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:      false,
		Encoding:         "json",
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:    "timestamp",
			LevelKey:   "level",
			MessageKey: "message",

			NameKey:       "",
			CallerKey:     "",
			FunctionKey:   "",
			StacktraceKey: "",

			EncodeTime:  zapcore.ISO8601TimeEncoder,
			EncodeLevel: zapcore.CapitalLevelEncoder,
		},
	}

	var err error
	Log, err = config.Build()
	if err != nil {
		panic(err)
	}
}
