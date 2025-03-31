# Async

## Usage

```
import "github.com/xxddpac/async"

var (
	wg   = &sync.WaitGroup{}
	pool = NewPoolWithFunc() // create a new pool with default(100) workers
)

// close the pool when you don't need it anymore
defer pool.Close()

// diaptch tasks to the pool
for i := 0; i < 1000; i++ {
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Add(-1)
		pool.Add(func(args ...interface{}) error {
			time.Sleep(time.Millisecond * 100)
			return nil
		})
	}(wg)
}

// wait for all tasks to finish
wg.Wait()
```