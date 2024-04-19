package logger

import (
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type ZapLogger struct {
	cfg  *ZapLoggerConfig
	log  *zap.Logger
	slog *zap.SugaredLogger
}

type ZapLoggerConfig struct {
	Level      string
	Filename   string
	MaxBackups int
	MaxSize    int
	MaxAge     int
	LocalTime  bool
	Compress   bool
}

type ZapLoggerOptions func(*ZapLogger)

func WithZapLevel(level string) ZapLoggerOptions {
	return func(l *ZapLogger) {
		l.cfg.Level = level
	}
}

func WithZapFilename(filename string) ZapLoggerOptions {
	return func(l *ZapLogger) {
		l.cfg.Filename = filename
	}
}

func WithZapMaxBackups(maxBackups int) ZapLoggerOptions {
	return func(l *ZapLogger) {
		l.cfg.MaxBackups = maxBackups
	}
}

func WithZapMaxSize(maxSize int) ZapLoggerOptions {
	return func(l *ZapLogger) {
		l.cfg.MaxSize = maxSize
	}
}

func WithZapMaxAge(maxAge int) ZapLoggerOptions {
	return func(l *ZapLogger) {
		l.cfg.MaxAge = maxAge
	}
}

func WithZapLocalTime(localTime bool) ZapLoggerOptions {
	return func(l *ZapLogger) {
		l.cfg.LocalTime = localTime
	}
}

func WithZapCompress(compress bool) ZapLoggerOptions {
	return func(l *ZapLogger) {
		l.cfg.Compress = compress
	}
}

var _ Logger = (*ZapLogger)(nil)

func NewZapLogger(options ...ZapLoggerOptions) *ZapLogger {
	l := &ZapLogger{
		cfg: &ZapLoggerConfig{
			Level:      "info",
			Filename:   "",
			MaxBackups: 10,
			MaxSize:    100,
			MaxAge:     30,
			LocalTime:  false,
			Compress:   false,
		},
	}

	for _, option := range options {
		option(l)
	}

	var level zap.AtomicLevel
	var syncOutput zapcore.WriteSyncer

	switch strings.ToLower(l.cfg.Level) {
	case "", "info":
		level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "debug":
		level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "warn":
		level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	case "panic":
		level = zap.NewAtomicLevelAt(zap.PanicLevel)
	case "fatal":
		level = zap.NewAtomicLevelAt(zap.FatalLevel)
	default:
		fmt.Printf("Invalid log level supplied. Defaulting to info: %s", l.cfg.Level)
		level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	if l.cfg.Filename != "" {
		syncOutput = zapcore.AddSync(&lumberjack.Logger{
			Filename:   l.cfg.Filename,
			MaxSize:    l.cfg.MaxSize,
			MaxBackups: l.cfg.MaxBackups,
			MaxAge:     l.cfg.MaxAge,
			LocalTime:  l.cfg.LocalTime,
			Compress:   l.cfg.Compress,
		})
	} else {
		syncOutput = zapcore.Lock(os.Stdout)
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		syncOutput,
		level,
	)

	l.log = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	l.slog = l.log.Sugar()
	zap.ReplaceGlobals(l.log)

	return l
}

func (l *ZapLogger) Debug(v ...any) {
	l.slog.Debug(v...)
}

func (l *ZapLogger) Debugf(format string, v ...any) {
	l.slog.Debugf(format, v...)
}

func (l *ZapLogger) Debugw(msg string, fields Fields) {
	zapFields := mapToZapFields(fields)
	l.log.Debug(msg, zapFields...)
}

func (l *ZapLogger) Info(v ...any) {
	l.slog.Info(v...)
}

func (l *ZapLogger) Infof(format string, v ...any) {
	l.slog.Infof(format, v...)
}

func (l *ZapLogger) Infow(msg string, fields Fields) {
	zapFields := mapToZapFields(fields)
	l.log.Info(msg, zapFields...)
}

func (l *ZapLogger) Warn(v ...any) {
	l.slog.Warn(v...)
}

func (l *ZapLogger) Warnf(format string, v ...any) {
	l.slog.Warnf(format, v...)
}

func (l *ZapLogger) Warnw(msg string, fields Fields) {
	zapFields := mapToZapFields(fields)
	l.log.Warn(msg, zapFields...)
}

func (l *ZapLogger) Error(v ...any) {
	l.slog.Error(v...)
}

func (l *ZapLogger) Errorf(format string, v ...any) {
	l.slog.Errorf(format, v...)
}

func (l *ZapLogger) Errorw(msg string, fields Fields) {
	zapFields := mapToZapFields(fields)
	l.log.Error(msg, zapFields...)
}

func mapToZapFields(fields map[string]interface{}) []zap.Field {
	var zapFields []zap.Field
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}

	return zapFields
}
