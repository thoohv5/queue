package log

import "log"

type ILog interface {
	Debug(msg string, values ...interface{})
	Info(msg string, values ...interface{})
	Warn(msg string, values ...interface{})
	Error(msg string, values ...interface{})
}

var l ILog

func GetLogger() ILog {
	if l == nil {
		l = NewDefaultLogger()
	}
	return l
}

func SetLogger(log ILog) {
	l = log
}

type defaultLogger struct {
}

func NewDefaultLogger() ILog {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	return &defaultLogger{}
}

func (dl *defaultLogger) Debug(msg string, values ...interface{}) {
	log.Printf(msg, values...)
}

func (dl *defaultLogger) Info(msg string, values ...interface{}) {
	log.Printf(msg, values...)
}

func (dl *defaultLogger) Warn(msg string, values ...interface{}) {
	log.Printf(msg, values...)
}

func (dl *defaultLogger) Error(msg string, values ...interface{}) {
	log.Printf(msg, values...)
}
