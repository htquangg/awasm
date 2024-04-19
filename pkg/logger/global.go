package logger

var global Logger

func init() {
	global = DefaultLogger
}

func SetLogger(log Logger) {
	global = log
}

func GetLogger() Logger {
	return global
}

func Debug(v ...any) {
	global.Debug(v...)
}

func Debugf(format string, v ...any) {
	global.Debugf(format, v...)
}

func Debugw(msg string, fields Fields) {
	global.Debugw(msg, fields)
}

func Info(v ...any) {
	global.Info(v...)
}

func Infof(format string, v ...any) {
	global.Infof(format, v...)
}

func Infow(msg string, fields Fields) {
	global.Infow(msg, fields)
}

func Warn(v ...any) {
	global.Warn(v...)
}

func Warnf(format string, v ...any) {
	global.Warnf(format, v...)
}

func Warnw(msg string, fields Fields) {
	global.Warnw(msg, fields)
}

func Error(v ...any) {
	global.Error(v...)
}

func Errorf(format string, v ...any) {
	global.Errorf(format, v...)
}

func Errorw(msg string, fields Fields) {
	global.Errorw(msg, fields)
}
