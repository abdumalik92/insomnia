package v1

import (
	"time"

	"github.com/alifcapital/keycloak_module/conf"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func NewZapLogger(cfg *conf.Config) *zap.Logger {
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   cfg.LogFilePath,
		MaxSize:    500,
		MaxAge:     30,
		MaxBackups: 0,
		LocalTime:  false,
		Compress:   false,
	})

	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(config),
		w,
		zap.InfoLevel,
	)
	logger := zap.New(core)

	return logger
}
