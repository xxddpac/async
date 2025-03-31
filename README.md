# Async

## Usage

```
import "github.com/xxddpac/async"

// create a pool with default(100) workers
pool := NewPoolWithFunc()

//dispatch tasks
for i := 0; i < 100; i++ {
	pool.Add(func(args ...interface{}) error {
		fmt.Println("do")
		return nil
	    })
    }
    
// block until you are ready to shut down
select {
case <-time.After(time.Second):
}

// close pool
pool.Close()
```