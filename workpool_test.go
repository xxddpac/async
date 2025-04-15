package async

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func BenchmarkThroughput1ms(b *testing.B) {
	pool := New()
	defer pool.Close()

	var executed int64
	task := func(args ...interface{}) error {
		atomic.AddInt64(&executed, 1)
		time.Sleep(1 * time.Millisecond)
		return nil
	}

	b.ResetTimer()
	start := time.Now()
	for i := 0; i < b.N; i++ {
		pool.Wg.Add(1)
		go func(i int) {
			defer pool.Wg.Done()
			pool.Add(task, i)
		}(i)
	}
	pool.Wg.Wait()
	elapsed := time.Since(start)
	b.StopTimer()
	b.ReportMetric(float64(executed)/elapsed.Seconds(), "tasks/s")
}

func BenchmarkHighConcurrency(b *testing.B) {
	pool := New()
	defer pool.Close()

	var executed int64
	task := func(args ...interface{}) error {
		atomic.AddInt64(&executed, 1)
		time.Sleep(1 * time.Millisecond)
		return nil
	}

	concurrency := 100
	tasksPerGoroutine := b.N / concurrency

	b.ResetTimer()
	for i := 0; i < concurrency; i++ {
		pool.Wg.Add(1)
		go func(start int) {
			defer pool.Wg.Done()
			for j := 0; j < tasksPerGoroutine; j++ {
				pool.Add(task, start+j)
			}
		}(i * tasksPerGoroutine)
	}
	pool.Wg.Wait()
	b.StopTimer()

	b.ReportMetric(float64(b.N)/b.Elapsed().Seconds(), "tasks/s")
}

func BenchmarkMemoryUsage(b *testing.B) {
	pool := New()
	defer pool.Close()

	task := func(args ...interface{}) error {
		time.Sleep(1 * time.Millisecond)
		return nil
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pool.Wg.Add(1)
		go func(i int) {
			defer pool.Wg.Done()
			pool.Add(task, i)
		}(i)
	}
	pool.Wg.Wait()
	b.StopTimer()
}

func TestPoolWithOptions(t *testing.T) {
	var (
		count   = 50000
		total   int
		mu      sync.Mutex
		errChan = make(chan error, count)
		pool    = New(WithMaxWorkers(count), WithMaxQueue(count), WithOnError(func(err error) {
			mu.Lock()
			total++
			mu.Unlock()
		}))
	)
	go func() {
		for err := range errChan {
			_ = err
		}
	}()
	defer func() {
		close(errChan)
		pool.Close()
	}()
	pool.Wg.Add(count)
	for i := 0; i < count; i++ {
		pool.Add(func(args ...interface{}) error {
			defer pool.Wg.Done()
			var err error
			err = fmt.Errorf("%s %d", args[0], args[1])
			select {
			case errChan <- err:
			default:
			}
			return err
		}, "test error", i)
	}
	pool.Wg.Wait()
	if total != count {
		t.Errorf("expected %d errors, got %d", count, total)
	}
}
