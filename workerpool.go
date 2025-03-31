package async

type Job func(args ...interface{}) error

type JobWithArgs struct {
	Fn  Job
	Arg []interface{}
}

var (
	maxWorkers = 100
	maxQueue   = 100
)

type WorkerPool struct {
	MaxWorkers  int
	MaxQueue    int
	workers     []*worker
	jobQueue    chan JobWithArgs
	workerQueue chan chan JobWithArgs
	quit        chan struct{}
	Logger      Logger
}

func NewPoolWithFunc(opts ...Option) *WorkerPool {
	wp := &WorkerPool{
		MaxWorkers: maxWorkers,
		MaxQueue:   maxQueue,
		Logger:     Printf,
	}
	for _, opt := range opts {
		opt(wp)
	}
	wp.jobQueue = make(chan JobWithArgs, wp.MaxQueue)
	wp.workerQueue = make(chan chan JobWithArgs, wp.MaxWorkers)
	wp.quit = make(chan struct{})
	wp.Run()
	return wp
}

func (w *WorkerPool) Run() *WorkerPool {
	w.Logger.Printf("WorkerPool init... maxWorkers: %d, maxQueue: %d", w.MaxWorkers, w.MaxQueue)
	for i := 0; i < w.MaxWorkers; i++ {
		nw := NewWorker(i, w)
		nw.Run(w.workerQueue)
		w.workers = append(w.workers, nw)
	}
	go func() {
		for {
			select {
			case job := <-w.jobQueue: // get job
				wq := <-w.workerQueue // get worker
				wq <- job             // send job to worker
			case <-w.quit:
				w.Logger.Printf("WorkerPool quited")
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
	return w.MaxWorkers
}

func (w *WorkerPool) Add(fn Job, args ...interface{}) {
	w.jobQueue <- JobWithArgs{Fn: fn, Arg: args}
}

func (w *WorkerPool) Close() {
	go func() {
		w.quit <- struct{}{}
	}()
}
