package async

type Option func(*WorkerPool)

func WithMaxWorkers(count int) Option {
	return func(c *WorkerPool) {
		c.MaxWorkers = count
	}
}

func WithMaxQueue(count int) Option {
	return func(c *WorkerPool) {
		c.MaxQueue = count
	}
}

func WithLogger(l Logger) Option {
	return func(c *WorkerPool) {
		c.Logger = l
	}
}
