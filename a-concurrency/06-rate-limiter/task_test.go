package ratelimit

import (
	"sort"
	"sync"
	"testing"
	"time"
)

// Время фиксируется уже после возврата из wait(), поэтому соседние
// отметки могут сблизиться на задержку планировщика.
const tolerance = 15 * time.Millisecond

// checkWindows проверяет требование «в любом интервале длиной 1 секунда —
// не больше rps возвратов»: любые rps+1 подряд идущих отметок должны
// занимать не меньше секунды.
func checkWindows(t *testing.T, times []time.Time, rps int) {
	t.Helper()
	sort.Slice(times, func(i, j int) bool { return times[i].Before(times[j]) })
	for i := 0; i+rps < len(times); i++ {
		if d := times[i+rps].Sub(times[i]); d < time.Second-tolerance {
			t.Fatalf("возвраты %d..%d (%d шт.) уложились в %v — больше %d за секунду",
				i, i+rps, rps+1, d, rps)
		}
	}
}

func TestLimitsRate(t *testing.T) {
	const rps, calls = 20, 30
	wait, stop := NewLimiter(rps)
	defer stop()

	start := time.Now()
	times := make([]time.Time, 0, calls)
	for i := 0; i < calls; i++ {
		wait()
		times = append(times, time.Now())
	}

	checkWindows(t, times, rps)
	// Лимит допускает 30 вызовов примерно за 1-1.5s.
	if elapsed := time.Since(start); elapsed > 3*time.Second {
		t.Fatalf("слишком медленно: %d вызовов при %d rps заняли %v", calls, rps, elapsed)
	}
}

// Пачка вызовов после простоя, начинающаяся в середине секунды:
// лимит должен держаться и на стыке секунд, а не только внутри каждой.
func TestBurstAfterIdle(t *testing.T) {
	const rps, calls = 20, 30
	wait, stop := NewLimiter(rps)
	defer stop()

	time.Sleep(525 * time.Millisecond)

	times := make([]time.Time, 0, calls)
	for i := 0; i < calls; i++ {
		wait()
		times = append(times, time.Now())
	}

	checkWindows(t, times, rps)
}

func TestConcurrentWait(t *testing.T) {
	const rps, calls = 20, 30
	wait, stop := NewLimiter(rps)
	defer stop()

	var (
		mu    sync.Mutex
		times []time.Time
		wg    sync.WaitGroup
	)
	start := time.Now()
	for i := 0; i < calls; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			wait()
			now := time.Now()
			mu.Lock()
			times = append(times, now)
			mu.Unlock()
		}()
	}
	wg.Wait()

	checkWindows(t, times, rps)
	if elapsed := time.Since(start); elapsed > 3*time.Second {
		t.Fatalf("слишком медленно: %d конкурентных вызовов при %d rps заняли %v", calls, rps, elapsed)
	}
}
