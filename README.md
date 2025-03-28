# Async

## Usage

```
import "github.com/xxddpac/async"

func Do() error {
	// do something
	return nil
}

pool := async.NewPoolWithFunc()

for i := 0; i < 100; i++ {
	pool.Add(Do)
}

pool.Close()
```