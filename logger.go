package async

import "log"

var _ Logger = LoggerFunc(nil)

type Logger interface {
	Printf(string, ...interface{})
}

type LoggerFunc func(string, ...interface{})

func (f LoggerFunc) Printf(msg string, args ...interface{}) { f(msg, args...) }

var Printf = LoggerFunc(log.Printf)
