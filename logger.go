package async

import "log"

// Logger is the interface that wraps the Printf method.
var _ Logger = LoggerFunc(nil)

type Logger interface {
	Printf(string, ...interface{})
}

type LoggerFunc func(string, ...interface{})

func (f LoggerFunc) Printf(msg string, args ...interface{}) { f(msg, args...) }

// Printf is a simple wrapper around log.Printf.
var Printf = LoggerFunc(log.Printf)
