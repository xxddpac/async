package async

type Option func(*WorkerPool)

func WithMaxWorkers(count int) Option {
	return func(c *WorkerPool) {
		c.maxWorkers = count
	}
}

func WithMaxQueue(count int) Option {
	return func(c *WorkerPool) {
		c.maxQueue = count
	}
}

func WithLogger(l Logger) Option {
	return func(c *WorkerPool) {
		c.Logger = l
	}
}
