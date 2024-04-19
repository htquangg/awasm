package logger

type Fields map[string]interface{}

var DefaultLogger = NewZapLogger()

type Logger interface {
	Debug(v ...any)
	Debugf(format string, v ...any)
	Debugw(msg string, fields Fields)
	Info(v ...any)
	Infof(format string, v ...any)
	Infow(msg string, fields Fields)
	Warn(v ...any)
	Warnf(format string, v ...any)
	Warnw(msg string, fields Fields)
	Error(v ...any)
	Errorf(format string, v ...any)
	Errorw(msg string, fields Fields)
}
