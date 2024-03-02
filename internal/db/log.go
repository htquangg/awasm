package db

import (
	"os"
	"sync/atomic"

	"github.com/htquangg/a-wasm/internal/constants"

	"github.com/segmentfault/pacman/log"
	xormlog "xorm.io/xorm/log"
)

type XORMLogBridge struct {
	showSQL  atomic.Bool
	logLevel string
}

func NewXORMLogger(showSQL bool) xormlog.Logger {
	logLevel := os.Getenv(constants.LogLevel)
	l := &XORMLogBridge{
		logLevel: logLevel,
	}
	l.showSQL.Store(showSQL)
	return l
}

func (l *XORMLogBridge) Debug(v ...interface{}) {
	log.Debug(v...)
}

func (l *XORMLogBridge) Debugf(format string, v ...interface{}) {
	log.Debugf(format, v...)
}

func (l *XORMLogBridge) Error(v ...interface{}) {
	log.Error(v...)
}

func (l *XORMLogBridge) Errorf(format string, v ...interface{}) {
	log.Errorf(format, v...)
}

func (l *XORMLogBridge) Info(v ...interface{}) {
	log.Info(v...)
}

func (l *XORMLogBridge) Infof(format string, v ...interface{}) {
	log.Infof(format, v...)
}

func (l *XORMLogBridge) Warn(v ...interface{}) {
	log.Warn(v...)
}

func (l *XORMLogBridge) Warnf(format string, v ...interface{}) {
	log.Warnf(format, v...)
}

func (l *XORMLogBridge) Level() xormlog.LogLevel {
	switch l.logLevel {
	case "debug":
		return xormlog.LOG_DEBUG
	case "info":
		return xormlog.LOG_INFO
	case "warn":
		return xormlog.LOG_WARNING
	case "error", "panic", "fatal":
		return xormlog.LOG_ERR
	}
	return xormlog.LOG_UNKNOWN
}

func (*XORMLogBridge) SetLevel(xormlog.LogLevel) {
}

func (l *XORMLogBridge) IsShowSQL() bool {
	return l.showSQL.Load()
}

func (l *XORMLogBridge) ShowSQL(show ...bool) {
	if len(show) == 0 {
		show = []bool{true}
	}
	l.showSQL.Store(show[0])
}
