package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var globalLogger *zap.Logger

// InitLogger initializes the global logger with the specified level and log file
func InitLogger(level string, logFile string) error {
	// Parse log level
	var zapLevel zapcore.Level
	switch level {
	case "debug":
		zapLevel = zapcore.DebugLevel
	case "info":
		zapLevel = zapcore.InfoLevel
	case "warn":
		zapLevel = zapcore.WarnLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	default:
		zapLevel = zapcore.InfoLevel
	}

	// Create encoder config
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Create console encoder (with colors)
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	// Create file encoder (without colors)
	fileEncoderConfig := encoderConfig
	fileEncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	fileEncoder := zapcore.NewJSONEncoder(fileEncoderConfig)

	// Create console writer
	consoleWriter := zapcore.AddSync(os.Stdout)

	// Create cores
	var cores []zapcore.Core
	cores = append(cores, zapcore.NewCore(consoleEncoder, consoleWriter, zapLevel))

	// Add file writer if log file is specified
	if logFile != "" {
		file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("failed to open log file: %w", err)
		}
		fileWriter := zapcore.AddSync(file)
		cores = append(cores, zapcore.NewCore(fileEncoder, fileWriter, zapLevel))
	}

	// Create logger
	core := zapcore.NewTee(cores...)
	globalLogger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return nil
}

// GetLogger returns the global logger instance
func GetLogger() *zap.Logger {
	if globalLogger == nil {
		// Fallback to a default logger if not initialized
		globalLogger, _ = zap.NewProduction()
	}
	return globalLogger
}

// Sync flushes any buffered log entries
func Sync() {
	if globalLogger != nil {
		_ = globalLogger.Sync()
	}
}

// Helper functions for common logging patterns
func Debug(msg string, fields ...zap.Field) {
	GetLogger().Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	GetLogger().Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	GetLogger().Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	GetLogger().Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	GetLogger().Fatal(msg, fields...)
}

/*
// Initialize logger
	if err := logger.InitLogger("info", "logs/app.log"); err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("Shutting down Agent Server...")

	logger.Info("Agent Server started successfully",
		zap.String("grpc", fmt.Sprintf(":%d", cfg.GRPCPort)),
		zap.String("tcp", fmt.Sprintf(":%d", cfg.TCPPort)),
		zap.String("http", fmt.Sprintf(":%d", cfg.HTTPPort)),
	)

	logger.Info("Starting Agent Server",
		zap.String("name", cfg.ServerName),
		zap.Int("grpc_port", cfg.GRPCPort),
		zap.Int("tcp_port", cfg.TCPPort),
		zap.Int("http_port", cfg.HTTPPort),
	)

	logger.Fatal("Failed", zap.Error(err))
	logger.Error("failed", zap.Error(err))
*/
