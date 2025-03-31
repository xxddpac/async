package async

import (
	"fmt"
	"testing"
	"time"
)

func TestNewWorkerDefault(t *testing.T) {
	pool := NewPoolWithFunc()
	for i := 0; i < 100; i++ {
		pool.Add(func(args ...interface{}) error {
			fmt.Println("do")
			return nil
		})
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
	pool := NewPoolWithFunc(WithMaxWorkers(200), WithMaxQueue(200), WithLogger(Log{}))
	for i := 0; i < 100; i++ {
		pool.Add(func(args ...interface{}) error {
			var id int
			if len(args) > 0 {
				id = args[0].(int)
			}
			fmt.Println("do", id)
			return nil
		}, 100)
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
