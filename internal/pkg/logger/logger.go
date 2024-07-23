package logger

import (
	"context"
	"io"
	"os"

	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

type contextKey string

const ContextKeyLogger contextKey = "logger"

func GetLoggerFromCtx(ctx context.Context) *logrus.Entry {
	return ctx.Value(ContextKeyLogger).(*logrus.Entry)
}

func NewLogger(level, format string) (*logrus.Logger, error) {
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		os.Mkdir("logs", 0o755)
	}
	f, err := os.OpenFile("./logs/logs.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
	if err != nil {
		return nil, err
	}

	logger := &logrus.Logger{
		Out:          io.MultiWriter(os.Stdout, f),
		ReportCaller: true,
	}

	switch level {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	case "fatal":
		logger.SetLevel(logrus.FatalLevel)
	case "panic":
		logger.SetLevel(logrus.PanicLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}

	if format == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	} else {
		logger.Formatter = &prefixed.TextFormatter{
			DisableColors:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
			ForceFormatting: true,
		}
	}
	// logrus.SetReportCaller(true)

	return logger, nil
}
