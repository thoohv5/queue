package log

import "log"

type ILog interface {
	Debug(msg string, values ...interface{})
	Info(msg string, values ...interface{})
	Warn(msg string, values ...interface{})
	Error(msg string, values ...interface{})
}

func GetLogger(log ILog, flag bool) ILog {
	if !flag {
		return NewEmptyLogger()
	}
	if log == nil {
		log = NewDefaultLogger()
	}
	return log
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

type emptyLogger struct {
}

func NewEmptyLogger() ILog {
	return &emptyLogger{}
}

func (dl *emptyLogger) Debug(msg string, values ...interface{}) {
}

func (dl *emptyLogger) Info(msg string, values ...interface{}) {
}

func (dl *emptyLogger) Warn(msg string, values ...interface{}) {
}

func (dl *emptyLogger) Error(msg string, values ...interface{}) {
}
