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

func WithOnError(fn func(err error)) Option {
	return func(c *WorkerPool) {
		c.onError = fn
	}
}
