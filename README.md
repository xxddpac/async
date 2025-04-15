# Async

## Usage

```
import "github.com/xxddpac/async"

//create a default async pool(maxWorkers=100, maxQueue=100)
p := async.New()

// or create a async pool with options
// p := async.New(async.WithMaxWorkers(), async.WithMaxQueue(),async.WithOnError(func(err error) {}))

// close the pool when you don't need it
defer p.Close()

// add tasks to the pool
for i := 0; i < 10000; i++ {
    p.Wg.Add(1)
    p.Add(func(args ...interface{}) error {
	defer p.Wg.Done()
	// do something
	return nil
	})
	
// wait for all tasks to finish
p.Wg.Wait()
```