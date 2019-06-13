package logger

import (
	"errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"strings"
)

var (
	zapLogger *zap.Logger
	zapInfoLogger *zap.Logger
)

type LoggerConfig struct {
	ErrLog	 		string
	InfoLog 		string
	MaxSize 		int
	MaxBackups 		int
	MaxAge 			int
	Level 			string
	LocalTime 		bool
}

func (lc *LoggerConfig) InitLogger() error {
	if strings.Trim(lc.ErrLog, " ") == "" || strings.Trim(lc.InfoLog, " ") == "" {
		return errors.New("plz specify logger path")
	}

	zapLogger = lc.newLogger(lc.ErrLog)
	if lc.ErrLog == lc.InfoLog {
		zapInfoLogger = zapLogger
		return nil
	}
	zapInfoLogger = lc.newLogger(lc.InfoLog)
	return nil
}

func (lc *LoggerConfig) newLogger(path string) *zap.Logger {
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename: path,
		MaxSize: lc.MaxSize,
		MaxBackups: lc.MaxBackups,
		MaxAge: lc.MaxAge,
		LocalTime: lc.LocalTime,
	})

	encoderConf := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	level := zap.DebugLevel
	levelLower := strings.ToLower(lc.Level)
	switch levelLower {
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	case "fatal":
		level = zap.FatalLevel
	}

	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderConf), w, level)
	return zap.New(core)
}

func Debug(brief string, detail string, mps ...map[string]string) {
	var fields = make([]zap.Field, 0)
	fields = append(fields, zap.String("detail", detail))

	for _, mp := range mps {
		for k, v := range mp {
			fields = append(fields, zap.String(k, v))
		}
	}

	zapInfoLogger.Debug(brief, fields...)
}

func Info(brief string, detail string, mps ...map[string]string) {
	var fields = make([]zap.Field, 0)
	fields = append(fields, zap.String("detail", detail))

	for _, mp := range mps {
		for k, v := range mp {
			fields = append(fields, zap.String(k, v))
		}
	}

	zapInfoLogger.Info(brief, fields...)
}

func Warn(brief string, detail string, mps ...map[string]string) {
	var fields = make([]zap.Field, 0)
	fields = append(fields, zap.String("detail", detail))

	for _, mp := range mps {
		for k, v := range mp {
			fields = append(fields, zap.String(k, v))
		}
	}

	zapInfoLogger.Warn(brief, fields...)
}

func Error(brief string, detail string, mps ...map[string]string) {
	var fields = make([]zap.Field, 0)
	fields = append(fields, zap.String("detail", detail))

	for _, mp := range mps {
		for k, v := range mp {
			fields = append(fields, zap.String(k, v))
		}
	}

	zapLogger.Error(brief, fields...)
}

func Fatal(brief string, detail string, mps ...map[string]string) {
	var fields = make([]zap.Field, 0)
	fields = append(fields, zap.String("detail", detail))

	for _, mp := range mps {
		for k, v := range mp {
			fields = append(fields, zap.String(k, v))
		}
	}

	zapLogger.Fatal(brief, fields...)
}
