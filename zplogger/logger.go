package zplogger

import (
	"os"
	"path/filepath"

	"go-echo-api/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Config holds logger options.
type LogConfig struct {
	LogDir     string
	LogFile    string
	ErrorFile  string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
	Level      string
	Console    bool
	DevMode    bool
}

func defaultConfig(cfg *config.Config) LogConfig {
	return LogConfig{
		LogDir:     filepath.Join(cfg.Cwd, "logs"),
		LogFile:    cfg.ServiceName + ".log",
		ErrorFile:  cfg.ServiceName + "-error.log",
		MaxSize:    50,
		MaxBackups: 7,
		MaxAge:     28,
		Compress:   false,
		Level:      cfg.LogLevel,
		Console:    true,
	}
}

// Logger wraps *zap.Logger
type Logger struct {
	*zap.Logger
}

// NewLogger creates a new zap logger with lumberjack rotation and custom encoders.
func NewLogger(cfg *config.Config) *Logger {
	opts := defaultConfig(cfg)

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "@timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stack_trace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	consoleEncoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stack_trace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var level zapcore.Level
	switch opts.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	var cores []zapcore.Core

	logDir := opts.LogDir
	_ = os.MkdirAll(logDir, 0755)

	// Main log file with lumberjack
	if opts.LogFile != "" {
		w := &lumberjack.Logger{
			Filename:   filepath.Join(logDir, opts.LogFile),
			MaxSize:    opts.MaxSize,
			MaxBackups: opts.MaxBackups,
			MaxAge:     opts.MaxAge,
			Compress:   opts.Compress,
		}
		cores = append(cores, zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(w),
			level,
		))
	}

	// Error log file with lumberjack
	if opts.ErrorFile != "" {
		w := &lumberjack.Logger{
			Filename:   filepath.Join(logDir, opts.ErrorFile),
			MaxSize:    opts.MaxSize,
			MaxBackups: opts.MaxBackups,
			MaxAge:     opts.MaxAge,
			Compress:   opts.Compress,
		}
		cores = append(cores, zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(w),
			zapcore.ErrorLevel,
		))
	}

	// Console output
	if opts.Console {
		cores = append(cores, zapcore.NewCore(
			zapcore.NewConsoleEncoder(consoleEncoderConfig),
			zapcore.AddSync(os.Stdout),
			level,
		))
	}

	if len(cores) == 0 {
		cores = append(cores, zapcore.NewCore(
			zapcore.NewConsoleEncoder(consoleEncoderConfig),
			zapcore.AddSync(os.Stderr),
			level,
		))
	}

	core := zapcore.NewTee(cores...)
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	return &Logger{logger}
}

// Close syncs the underlying zap logger.
func (l *Logger) Close() error {
	if l.Logger != nil {
		_ = l.Logger.Sync()
	}
	return nil
}
