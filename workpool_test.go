package async

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func BenchmarkThroughput1ms(b *testing.B) {
	pool := NewPoolWithFunc()
	defer pool.Close()

	var executed int64
	task := func(args ...interface{}) error {
		atomic.AddInt64(&executed, 1)
		time.Sleep(1 * time.Millisecond)
		return nil
	}

	b.ResetTimer()
	start := time.Now()
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			pool.Add(task, i)
		}(i)
	}
	wg.Wait()
	elapsed := time.Since(start)
	b.StopTimer()
	b.ReportMetric(float64(executed)/elapsed.Seconds(), "tasks/s")
}

func BenchmarkHighConcurrency(b *testing.B) {
	pool := NewPoolWithFunc()
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
	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(start int) {
			defer wg.Done()
			for j := 0; j < tasksPerGoroutine; j++ {
				pool.Add(task, start+j)
			}
		}(i * tasksPerGoroutine)
	}
	wg.Wait()
	b.StopTimer()

	b.ReportMetric(float64(b.N)/b.Elapsed().Seconds(), "tasks/s")
}

func BenchmarkMemoryUsage(b *testing.B) {
	pool := NewPoolWithFunc()
	defer pool.Close()

	task := func(args ...interface{}) error {
		time.Sleep(1 * time.Millisecond)
		return nil
	}

	b.ResetTimer()
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			pool.Add(task, i)
		}(i)
	}
	wg.Wait()
	b.StopTimer()
}
