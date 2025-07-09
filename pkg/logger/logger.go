package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

// Init 初始化日志
func Init(mode string) {
	log = logrus.New()

	// 设置输出
	log.SetOutput(os.Stdout)

	// 根据模式设置日志级别和格式
	switch mode {
	case "debug":
		log.SetLevel(logrus.DebugLevel)
		log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	case "release":
		log.SetLevel(logrus.InfoLevel)
		log.SetFormatter(&logrus.JSONFormatter{})
	default:
		log.SetLevel(logrus.InfoLevel)
		log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}
}

// Debug 调试日志
func Debug(format string, args ...interface{}) {
	if log != nil {
		log.Debugf(format, args...)
	}
}

// Info 信息日志
func Info(format string, args ...interface{}) {
	if log != nil {
		log.Infof(format, args...)
	}
}

// Warn 警告日志
func Warn(format string, args ...interface{}) {
	if log != nil {
		log.Warnf(format, args...)
	}
}

// Error 错误日志
func Error(format string, args ...interface{}) {
	if log != nil {
		log.Errorf(format, args...)
	}
}

// Fatal 致命错误日志
func Fatal(format string, args ...interface{}) {
	if log != nil {
		log.Fatalf(format, args...)
	}
}

// WithFields 带字段的日志
func WithFields(fields map[string]interface{}) *logrus.Entry {
	if log != nil {
		return log.WithFields(fields)
	}
	return nil
}
