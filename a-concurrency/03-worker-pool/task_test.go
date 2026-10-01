package workerpool

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestAllJobsProcessed(t *testing.T) {
	p := NewPool(5, func(j int) int { return j * 2 })

	const jobs = 100
	go func() {
		for i := 0; i < jobs; i++ {
			p.Submit(i)
		}
		p.Close()
	}()

	got := map[int]bool{}
	count := 0
	for r := range p.Results() {
		got[r] = true
		count++
	}

	if count != jobs {
		t.Fatalf("ожидали %d результатов, получили %d", jobs, count)
	}
	for i := 0; i < jobs; i++ {
		if !got[i*2] {
			t.Fatalf("потерян результат для задачи %d", i)
		}
	}
}

func TestWorkersLimited(t *testing.T) {
	var current, max int64

	p := NewPool(3, func(j int) int {
		c := atomic.AddInt64(&current, 1)
		for {
			m := atomic.LoadInt64(&max)
			if c <= m || atomic.CompareAndSwapInt64(&max, m, c) {
				break
			}
		}
		time.Sleep(20 * time.Millisecond)
		atomic.AddInt64(&current, -1)
		return j
	})

	go func() {
		for i := 0; i < 12; i++ {
			p.Submit(i)
		}
		p.Close()
	}()
	for range p.Results() {
	}

	if m := atomic.LoadInt64(&max); m > 3 {
		t.Fatalf("одновременно работало %d воркеров при лимите 3", m)
	}
	if m := atomic.LoadInt64(&max); m < 2 {
		t.Fatalf("параллельности не видно: максимум одновременных воркеров %d", m)
	}
}

func TestParallelSpeedup(t *testing.T) {
	p := NewPool(10, func(j int) int {
		time.Sleep(10 * time.Millisecond)
		return j
	})

	start := time.Now()
	go func() {
		for i := 0; i < 50; i++ {
			p.Submit(i)
		}
		p.Close()
	}()
	for range p.Results() {
	}
	elapsed := time.Since(start)

	// Последовательно было бы ~500ms; с 10 воркерами — ~50-100ms.
	if elapsed > 300*time.Millisecond {
		t.Fatalf("слишком медленно (%v): похоже, задачи выполняются последовательно", elapsed)
	}
}

func TestConcurrentSubmit(t *testing.T) {
	p := NewPool(4, func(j int) int { return j })

	var wg sync.WaitGroup
	for g := 0; g < 5; g++ {
		wg.Add(1)
		go func(base int) {
			defer wg.Done()
			for i := 0; i < 20; i++ {
				p.Submit(base*100 + i)
			}
		}(g)
	}
	go func() {
		wg.Wait()
		p.Close()
	}()

	count := 0
	for range p.Results() {
		count++
	}
	if count != 100 {
		t.Fatalf("ожидали 100 результатов, получили %d", count)
	}
}
