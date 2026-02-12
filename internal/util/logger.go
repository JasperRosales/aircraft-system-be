package util

import (
	"context"
	"fmt"
	"os"
	"time"

	gormLogger "gorm.io/gorm/logger"
)

type Logger struct{}

func NewLogger() *Logger {
	return &Logger{}
}

func formatTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func (l *Logger) Info(message string, keysAndValues ...interface{}) {
	logWithFormat("INFO", message, keysAndValues...)
}

func (l *Logger) Warn(message string, keysAndValues ...interface{}) {
	logWithFormat("WARN", message, keysAndValues...)
}

func (l *Logger) Error(message string, keysAndValues ...interface{}) {
	logWithFormat("ERROR", message, keysAndValues...)
}

func (l *Logger) Debug(message string, keysAndValues ...interface{}) {
	logWithFormat("DEBUG", message, keysAndValues...)
}

func (l *Logger) Fatal(message string, keysAndValues ...interface{}) {
	logWithFormat("FATAL", message, keysAndValues...)
	os.Exit(1)
}

func logWithFormat(level, message string, keysAndValues ...interface{}) {
	logLine := fmt.Sprintf("[%s] %s: %s", formatTime(), level, message)

	if len(keysAndValues) > 0 {
		logLine += " |"
		for i := 0; i < len(keysAndValues); i += 2 {
			if i+1 < len(keysAndValues) {
				key, ok := keysAndValues[i].(string)
				if ok {
					logLine += fmt.Sprintf(" %s=%v", key, keysAndValues[i+1])
				}
			}
		}
	}

	fmt.Println(logLine)
}

type GormLogger struct{}

func NewGormLogger() *GormLogger {
	return &GormLogger{}
}

func (l *GormLogger) Printf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	logWithFormat("SQL", msg)
}

func (l *GormLogger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	return l
}

func (l *GormLogger) Info(ctx context.Context, msg string, args ...interface{}) {
	logWithFormat("SQL-INFO", msg, args...)
}

func (l *GormLogger) Warn(ctx context.Context, msg string, args ...interface{}) {
	logWithFormat("SQL-WARN", msg, args...)
}

func (l *GormLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	logWithFormat("SQL-ERROR", msg, args...)
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	logWithFormat("SQL-TRACE", "SQL executed",
		"elapsed", elapsed,
		"sql", sql,
		"rows", rows,
		"error", err,
	)
}

func EnsureLogDirectory() error {
	return nil
}
