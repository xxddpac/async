package async

import (
	"runtime/debug"
)

type Status int

const (
	Idle Status = iota + 1
	Busy
	Closed
)

type worker struct {
	id         int
	jobChannel chan *JobWithArgs
	quit       chan struct{}
	status     Status
	pool       *WorkerPool
}

func NewWorker(id int, pool *WorkerPool) *worker {
	return &worker{
		id:         id,
		jobChannel: make(chan *JobWithArgs),
		quit:       make(chan struct{}),
		status:     Idle,
		pool:       pool,
	}
}

func (w *worker) Run(wq chan<- chan *JobWithArgs) {
	go func() {
		defer func() {
			panicErr := recover()
			if panicErr != nil {
				debug.PrintStack()
				w.pool.Logger.Printf("【worker-%d】run panic, err info: %v", w.id, panicErr)
				w.pool.Logger.Printf("【worker-%d】panic error stack ==> %s", w.id, debug.Stack())
				w.pool.Logger.Printf("【worker-%d】 recover worker...", w.id)
				w.Run(wq)
			}
		}()
		for {
			wq <- w.jobChannel
			w.status = Idle
			select {
			case job := <-w.jobChannel:
				func() { defer w.pool.sp.Put(job) }()
				w.status = Busy
				if err := job.fn(job.args...); err != nil {
					w.pool.Logger.Printf("【worker-%d】run error, err info: %v", w.id, err)
					if w.pool.onError != nil {
						w.pool.onError(err)
					}
				}
			case <-w.quit:
				w.status = Closed
				return
			}
		}
	}()
}

func (w *worker) Status() Status {
	return w.status
}

func (w *worker) Close() {
	go func() {
		w.quit <- struct{}{}
	}()
}
