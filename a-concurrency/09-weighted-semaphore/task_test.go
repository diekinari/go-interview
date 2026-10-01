package wsema

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestTryLock(t *testing.T) {
	s := NewSema(5)

	if !s.TryLock(3) {
		t.Fatal("TryLock(3) при 5 свободных должен пройти")
	}
	if s.Free() != 2 {
		t.Fatalf("после захвата 3 из 5 свободно должно быть 2, а не %d", s.Free())
	}
	if s.TryLock(3) {
		t.Fatal("TryLock(3) при 2 свободных должен вернуть false")
	}
	if !s.TryLock(2) {
		t.Fatal("TryLock(2) при 2 свободных должен пройти")
	}

	s.Release(5)
	if s.Free() != 5 {
		t.Fatalf("после Release всё должно освободиться, свободно %d", s.Free())
	}
}

func TestLockBlocksUntilRelease(t *testing.T) {
	s := NewSema(4)
	s.Lock(3)

	var acquired atomic.Bool
	done := make(chan struct{})
	go func() {
		s.Lock(2) // должен ждать: свободен только 1
		acquired.Store(true)
		close(done)
	}()

	time.Sleep(50 * time.Millisecond)
	if acquired.Load() {
		t.Fatal("Lock(2) при 1 свободном не должен был пройти")
	}

	s.Release(3)
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("после Release ожидающий Lock должен был проснуться")
	}
}

// Две горутины одновременно хотят по 6 слотов из 10: одна должна
// получить их, вторая — дождаться её Release. Повторяем много раундов,
// чтобы неудачное чередование почти наверняка случилось хотя бы раз.
func TestTwoBigLocks(t *testing.T) {
	done := make(chan struct{})
	go func() {
		defer close(done)
		for round := 0; round < 50; round++ {
			s := NewSema(10)
			start := make(chan struct{})
			var wg sync.WaitGroup
			for i := 0; i < 2; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					<-start
					s.Lock(6)
					time.Sleep(time.Millisecond)
					s.Release(6)
				}()
			}
			close(start)
			wg.Wait()
		}
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("две горутины по Lock(6) при cap=10 зависли — взаимная блокировка")
	}
}

func TestConcurrentInvariant(t *testing.T) {
	const cap = 10
	s := NewSema(cap)

	var current, max int64
	var wg sync.WaitGroup
	for i := 0; i < 30; i++ {
		wg.Add(1)
		go func(weight int) {
			defer wg.Done()
			s.Lock(weight)
			c := atomic.AddInt64(&current, int64(weight))
			for {
				m := atomic.LoadInt64(&max)
				if c <= m || atomic.CompareAndSwapInt64(&max, m, c) {
					break
				}
			}
			time.Sleep(5 * time.Millisecond)
			atomic.AddInt64(&current, -int64(weight))
			s.Release(weight)
		}(i%3 + 1)
	}
	wg.Wait()

	if m := atomic.LoadInt64(&max); m > cap {
		t.Fatalf("инвариант нарушен: суммарный захваченный вес доходил до %d при cap=%d", m, cap)
	}
	if s.Free() != cap {
		t.Fatalf("после завершения всех горутин свободно %d из %d", s.Free(), cap)
	}
}
