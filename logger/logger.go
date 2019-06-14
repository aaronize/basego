package logger

import (
	"encoding/json"
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

type Options struct {
	ErrLog	 		string
	InfoLog 		string
	MaxSize 		int
	MaxBackups 		int
	MaxAge 			int
	Level 			string
	LocalTime 		bool
}

func (lc *Options) InitLogger() error {
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

func (lc *Options) newLogger(path string) *zap.Logger {
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

func Debug(brief string, detail string, mps ...map[string]interface{}) {
	var fields = make([]zap.Field, 0)
	fields = append(fields, zap.String("detail", detail))

	for _, mp := range mps {
		for k, v := range mp {
			fields = appendFields(fields, k, v)
		}
	}

	zapInfoLogger.Debug(brief, fields...)
}

func Info(brief string, detail string, mps ...map[string]interface{}) {
	var fields = make([]zap.Field, 0)
	fields = append(fields, zap.String("detail", detail))

	for _, mp := range mps {
		for k, v := range mp {
			fields = appendFields(fields, k, v)
		}
	}

	zapInfoLogger.Info(brief, fields...)
}

func Warn(brief string, detail string, mps ...map[string]interface{}) {
	var fields = make([]zap.Field, 0)
	fields = append(fields, zap.String("detail", detail))

	for _, mp := range mps {
		for k, v := range mp {
			fields = appendFields(fields, k, v)
		}
	}

	zapInfoLogger.Warn(brief, fields...)
}

func Error(brief string, detail string, mps ...map[string]interface{}) {
	var fields = make([]zap.Field, 0)
	fields = append(fields, zap.String("detail", detail))

	for _, mp := range mps {
		for k, v := range mp {
			fields = appendFields(fields, k, v)
		}
	}

	zapLogger.Error(brief, fields...)
}

func Fatal(brief string, detail string, mps ...map[string]interface{}) {
	var fields = make([]zap.Field, 0)
	fields = append(fields, zap.String("detail", detail))

	for _, mp := range mps {
		for k, v := range mp {
			fields = appendFields(fields, k, v)
		}
	}

	zapLogger.Fatal(brief, fields...)
}

func appendFields(fields []zap.Field, k string, v interface{}) []zap.Field {

	switch v.(type) {
	case string:
		fields = append(fields, zap.String(k, v.(string)))
	case int:
		fields = append(fields, zap.Int(k, v.(int)))
	case int8:
		fields = append(fields, zap.Int8(k, v.(int8)))
	case int16:
		fields = append(fields, zap.Int16(k, v.(int16)))
	case int32:
		fields = append(fields, zap.Int32(k, v.(int32)))
	case int64:
		fields = append(fields, zap.Int64(k, v.(int64)))
	case uint:
		fields = append(fields, zap.Uint(k, v.(uint)))
	case uint8:
		fields = append(fields, zap.Uint8(k, v.(uint8)))
	case uint16:
		fields = append(fields, zap.Uint16(k, v.(uint16)))
	case uint32:
		fields = append(fields, zap.Uint32(k, v.(uint32)))
	case uint64:
		fields = append(fields, zap.Uint64(k, v.(uint64)))
	case float32:
		fields = append(fields, zap.Float32(k, v.(float32)))
	case float64:
		fields = append(fields, zap.Float64(k, v.(float64)))
	case bool:
		fields = append(fields, zap.Bool(k, v.(bool)))
	case []byte:
		fields = append(fields, zap.Binary(k, v.([]byte)))
	default:
		bt, _ := json.Marshal(v)
		fields = append(fields, zap.String(k, string(bt)))
	}

	return fields
}
