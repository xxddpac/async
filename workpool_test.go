package async

import (
	"log"
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

// TestPoolWithOptions tests the worker pool with options

type Log struct {
}

func (l Log) Printf(format string, args ...interface{}) {
	//zap.S().Infof(msg, args) // Uncomment this line if you want to use zap logger
	log.Printf(format, args...)
}

func TestPoolWithOptions(t *testing.T) {
	var (
		count = 10000
		pool  = New(
			WithMaxWorkers(count),
			WithMaxQueue(count),
			WithOnError(func(err error) {
				// callback function for error handling
			}),
			WithLogger(Log{}))
	)
	defer pool.Close()
	pool.Wg.Add(count)
	for i := 0; i < count; i++ {
		pool.Add(func(args ...interface{}) error {
			defer pool.Wg.Done()
			time.Sleep(10 * time.Millisecond)
			return nil
		})
	}
	pool.Wg.Wait()
	pool.Logger.Printf("all tasks completed")
}
