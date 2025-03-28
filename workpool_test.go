package async

import (
	"fmt"
	"testing"
	"time"
)

func Do() error {
	fmt.Println("do")
	return nil
}

func TestNewWorkerDefault(t *testing.T) {
	pool := NewPoolWithFunc()
	for i := 0; i < 100; i++ {
		pool.Add(Do)
	}
	select {
	case <-time.After(time.Second):
	}
	pool.Close()
}

type Log struct {
}

func (l Log) Printf(format string, v ...interface{}) {
}

func TestNewWorkerWithOptions(t *testing.T) {
	pool := NewPoolWithFunc(WithMaxWorkers(100), WithMaxQueue(100), WithLogger(Log{}))
	for i := 0; i < 100; i++ {
		pool.Add(Do)
	}
	fmt.Println(pool.WorkerCount())
	for _, wk := range pool.workers {
		fmt.Println(wk.id, wk.status)
	}
	select {
	case <-time.After(time.Second):
	}
	pool.Close()
}
