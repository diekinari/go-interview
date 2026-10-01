package crawler

import (
	"fmt"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// Детерминированная заглушка: код зависит от URL.
func stubFetch(url string) int {
	time.Sleep(15 * time.Millisecond)
	if strings.HasSuffix(url, "7") {
		return 500
	}
	return 200
}

func makeURLs(n int) []string {
	urls := make([]string, n)
	for i := range urls {
		urls[i] = fmt.Sprintf("https://site.ru/page/%d", i)
	}
	return urls
}

func TestOrderAndValues(t *testing.T) {
	urls := makeURLs(20)
	got := Crawl(urls, 5, stubFetch)

	if len(got) != len(urls) {
		t.Fatalf("ожидали %d кодов, получили %d", len(urls), len(got))
	}
	for i, url := range urls {
		want := 200
		if strings.HasSuffix(url, "7") {
			want = 500
		}
		if got[i] != want {
			t.Fatalf("позиция %d (%s): ожидали %d, получили %d — порядок результатов нарушен", i, url, want, got[i])
		}
	}
}

func TestConcurrencyLimit(t *testing.T) {
	var current, max int64
	instrumented := func(url string) int {
		c := atomic.AddInt64(&current, 1)
		for {
			m := atomic.LoadInt64(&max)
			if c <= m || atomic.CompareAndSwapInt64(&max, m, c) {
				break
			}
		}
		time.Sleep(20 * time.Millisecond)
		atomic.AddInt64(&current, -1)
		return 200
	}

	Crawl(makeURLs(15), 3, instrumented)

	if m := atomic.LoadInt64(&max); m > 3 {
		t.Fatalf("одновременно выполнялось %d запросов при лимите 3", m)
	}
	if m := atomic.LoadInt64(&max); m < 2 {
		t.Fatalf("параллельности не видно: максимум одновременных запросов %d", m)
	}
}

func TestActuallyParallel(t *testing.T) {
	start := time.Now()
	Crawl(makeURLs(12), 4, stubFetch) // 12 задач по 15ms при k=4 -> ~45-60ms
	elapsed := time.Since(start)

	if elapsed > 120*time.Millisecond {
		t.Fatalf("слишком медленно (%v): похоже, обход последовательный", elapsed)
	}
}

func TestEmptyInput(t *testing.T) {
	got := Crawl(nil, 3, stubFetch)
	if len(got) != 0 {
		t.Fatalf("для пустого входа ожидали пустой результат, получили %v", got)
	}
}
