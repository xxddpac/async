package async

import (
	"sync"
)

type Job func(args ...interface{}) error

type JobWithArgs struct {
	fn   Job
	args []interface{}
}

var (
	maxWorkers = 100
	maxQueue   = 100
)

type WorkerPool struct {
	maxWorkers  int
	maxQueue    int
	workers     []*worker
	jobQueue    chan *JobWithArgs
	workerQueue chan chan *JobWithArgs
	quit        chan struct{}
	Logger      Logger
	onError     func(err error)
	sp          sync.Pool
	Wg          sync.WaitGroup
}

func New(opts ...Option) *WorkerPool {
	wp := &WorkerPool{
		maxWorkers: maxWorkers,
		maxQueue:   maxQueue,
		Logger:     Printf,
	}
	for _, opt := range opts {
		opt(wp)
	}
	wp.jobQueue = make(chan *JobWithArgs, wp.maxQueue)
	wp.workerQueue = make(chan chan *JobWithArgs, wp.maxWorkers)
	wp.quit = make(chan struct{})
	wp.sp.New = func() interface{} { return &JobWithArgs{} }
	wp.Run()
	return wp
}

func (w *WorkerPool) Run() *WorkerPool {
	w.Logger.Printf("WorkerPool init... maxWorkers: %d, maxQueue: %d", w.maxWorkers, w.maxQueue)
	for i := 0; i < w.maxWorkers; i++ {
		nw := NewWorker(i, w)
		nw.Run(w.workerQueue)
		w.workers = append(w.workers, nw)
	}
	go func() {
		for {
			select {
			case job := <-w.jobQueue:
				wq := <-w.workerQueue
				wq <- job
			case <-w.quit:
				w.Logger.Printf("WorkerPool close...")
				for _, wk := range w.workers {
					wk.Close()
				}
				return
			}
		}
	}()
	return w
}

func (w *WorkerPool) WorkerCount() int {
	return w.maxWorkers
}

func (w *WorkerPool) Add(fn Job, args ...interface{}) {
	job := w.sp.Get().(*JobWithArgs)
	job.fn = fn
	job.args = args
	w.jobQueue <- job
}

func (w *WorkerPool) Close() {
	go func() {
		w.quit <- struct{}{}
	}()
}
