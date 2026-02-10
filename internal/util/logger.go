package util

import (
	"context"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	gormLogger "gorm.io/gorm/logger"
)

const (
	LogPath = "logs/api.log"
)

// Embeds *zap.Logger so it can be used wherever *zap.Logger is expected
type Logger struct {
	*zap.Logger
}

// NewLogger creates a new logger instance
func NewLogger() (*Logger, error) {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{LogPath}
	config.EncoderConfig = zap.NewProductionEncoderConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return &Logger{Logger: logger}, nil
}

// NewLoggerWithConfig creates a logger with custom configuration
func NewLoggerWithConfig(outputPaths []string, level zapcore.Level) (*Logger, error) {
	config := zap.NewProductionConfig()
	config.OutputPaths = outputPaths
	config.Level = zap.NewAtomicLevelAt(level)
	config.EncoderConfig = zap.NewProductionEncoderConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return &Logger{Logger: logger}, nil
}

// EnsureLogDirectory creates the log directory if it doesn't exist
func EnsureLogDirectory() error {
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		if err := os.Mkdir("logs", 0755); err != nil {
			return err
		}
	}

	// Create the log file if it doesn't exist
	if _, err := os.Stat(LogPath); os.IsNotExist(err) {
		file, err := os.Create(LogPath)
		if err != nil {
			return err
		}
		file.Close()
	}

	return nil
}

type GormLogger struct {
	logger *zap.Logger
}

func NewGormLogger(logger *zap.Logger) *GormLogger {
	return &GormLogger{logger: logger}
}

func (l *GormLogger) Printf(format string, args ...interface{}) {
	l.logger.Info(format, zap.Any("args", args))
}

func (l *GormLogger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	return l
}

func (l *GormLogger) Info(ctx context.Context, msg string, args ...interface{}) {
	l.logger.Info(msg, zap.Any("args", args))
}

func (l *GormLogger) Warn(ctx context.Context, msg string, args ...interface{}) {
	l.logger.Warn(msg, zap.Any("args", args))
}

func (l *GormLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	l.logger.Error(msg, zap.Any("args", args))
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	l.logger.Debug("SQL Trace",
		zap.Duration("elapsed", elapsed),
		zap.String("sql", sql),
		zap.Int64("rows", rows),
		zap.Error(err),
	)
}
