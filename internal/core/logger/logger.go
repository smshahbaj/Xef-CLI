// Package logger provides structured logging capabilities for XefCLI.
// It supports multiple output formats, log levels, and rotation-aware logging.
package logger

import (
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LogLevel represents the severity of a log entry.
type LogLevel string

// Standard log levels used by the logger.
const (
	// DebugLevel enables verbose debugging output.
	DebugLevel LogLevel = "debug"
	// InfoLevel is the default informational logging level.
	InfoLevel LogLevel = "info"
	// WarnLevel indicates non-fatal warnings.
	WarnLevel LogLevel = "warn"
	// ErrorLevel indicates errors requiring attention.
	ErrorLevel LogLevel = "error"
	// FatalLevel indicates unrecoverable errors that will exit.
	FatalLevel LogLevel = "fatal"
)

// Config holds logger configuration.
type Config struct {
	Level       LogLevel
	Format      string // json or pretty
	Output      string // stdout, stderr, or file path
	Development bool
}

// Logger is the application logger interface.
type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
	With(fields ...Field) Logger
	Sync() error
}

// Field represents a structured log field.
type Field struct {
	Value interface{}
	Key   string
}

// String creates a string field.
func String(key, value string) Field {
	return Field{Key: key, Value: value}
}

// Int creates an int field.
func Int(key string, value int) Field {
	return Field{Key: key, Value: value}
}

// Int64 creates an int64 field.
func Int64(key string, value int64) Field {
	return Field{Key: key, Value: value}
}

// Float64 creates a float64 field.
func Float64(key string, value float64) Field {
	return Field{Key: key, Value: value}
}

// Bool creates a bool field.
func Bool(key string, value bool) Field {
	return Field{Key: key, Value: value}
}

// Error creates an error field.
func Error(err error) Field {
	return Field{Key: "error", Value: err}
}

// zapLogger implements Logger using zap.
type zapLogger struct {
	*zap.SugaredLogger
}

// New creates a new Logger from Config.
func New(cfg Config) (Logger, error) {
	level, err := parseLevel(cfg.Level)
	if err != nil {
		return nil, fmt.Errorf("invalid log level %q: %w", cfg.Level, err)
	}

	var encoder zapcore.Encoder
	if strings.EqualFold(cfg.Format, "json") {
		encoder = zapcore.NewJSONEncoder(encoderConfig(cfg.Development))
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig(cfg.Development))
	}

	var writeSyncer zapcore.WriteSyncer
	switch strings.ToLower(cfg.Output) {
	case "stderr":
		writeSyncer = zapcore.AddSync(os.Stderr)
	case "stdout", "":
		writeSyncer = zapcore.AddSync(os.Stdout)
	default:
		f, err := os.OpenFile(cfg.Output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}
		writeSyncer = zapcore.AddSync(f)
	}

	core := zapcore.NewCore(encoder, writeSyncer, level)
	options := []zap.Option{zap.AddCallerSkip(1)}
	if cfg.Development {
		options = append(options, zap.Development())
	}

	z := zap.New(core, options...)
	return &zapLogger{SugaredLogger: z.Sugar()}, nil
}

func parseLevel(level LogLevel) (zapcore.Level, error) {
	switch strings.ToLower(string(level)) {
	case "debug":
		return zapcore.DebugLevel, nil
	case "info":
		return zapcore.InfoLevel, nil
	case "warn", "warning":
		return zapcore.WarnLevel, nil
	case "error":
		return zapcore.ErrorLevel, nil
	case "fatal":
		return zapcore.FatalLevel, nil
	default:
		return zapcore.InfoLevel, fmt.Errorf("unknown level: %s", level)
	}
}

func encoderConfig(dev bool) zapcore.EncoderConfig {
	cfg := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	if dev {
		cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	return cfg
}

func (l *zapLogger) Debug(msg string, fields ...Field) {
	l.Debugw(msg, toArgs(fields...)...)
}

func (l *zapLogger) Info(msg string, fields ...Field) {
	l.Infow(msg, toArgs(fields...)...)
}

func (l *zapLogger) Warn(msg string, fields ...Field) {
	l.Warnw(msg, toArgs(fields...)...)
}

func (l *zapLogger) Error(msg string, fields ...Field) {
	l.Errorw(msg, toArgs(fields...)...)
}

func (l *zapLogger) Fatal(msg string, fields ...Field) {
	l.Fatalw(msg, toArgs(fields...)...)
}

func (l *zapLogger) With(fields ...Field) Logger {
	return &zapLogger{SugaredLogger: l.SugaredLogger.With(toArgs(fields...)...)}
}

func toArgs(fields ...Field) []interface{} {
	args := make([]interface{}, 0, len(fields)*2)
	for _, f := range fields {
		args = append(args, f.Key, f.Value)
	}
	return args
}

// Nop returns a no-op logger for testing.
func Nop() Logger {
	return &zapLogger{SugaredLogger: zap.NewNop().Sugar()}
}
