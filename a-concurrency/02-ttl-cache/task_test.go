package ttlcache

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestGetBeforeExpiry(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c := NewTTLCache(ctx, 200*time.Millisecond)

	c.Set("a", 1)
	if v, ok := c.Get("a"); !ok || v != 1 {
		t.Fatalf("свежая запись должна читаться: получили (%v, %v)", v, ok)
	}
}

func TestExpiredNotReturned(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c := NewTTLCache(ctx, 50*time.Millisecond)

	c.Set("a", 1)
	time.Sleep(100 * time.Millisecond)

	if _, ok := c.Get("a"); ok {
		t.Fatal("протухшая запись не должна отдаваться из Get")
	}
}

func TestSetRefreshesTTL(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c := NewTTLCache(ctx, 120*time.Millisecond)

	c.Set("a", 1)
	time.Sleep(70 * time.Millisecond)
	c.Set("a", 2) // срок должен обновиться
	time.Sleep(70 * time.Millisecond)

	if v, ok := c.Get("a"); !ok || v != 2 {
		t.Fatalf("повторный Set обновляет TTL: ожидали (2, true), получили (%v, %v)", v, ok)
	}
}

func TestCleanupWithoutReads(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	const ttl = 40 * time.Millisecond
	c := NewTTLCache(ctx, ttl)

	for i := 0; i < 10; i++ {
		c.Set(fmt.Sprintf("key-%d", i), i)
	}
	// Никто не зовёт Get — протухшие записи всё равно должны исчезнуть
	// не позже чем через ttl после истечения; ждём с запасом.
	time.Sleep(3 * ttl)

	if n := c.Len(); n != 0 {
		t.Fatalf("протухшие записи без обращений не исчезли, Len()=%d", n)
	}
}

func TestConcurrent(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c := NewTTLCache(ctx, 50*time.Millisecond)

	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			for j := 0; j < 200; j++ {
				key := fmt.Sprintf("key-%d", j%10)
				if n%2 == 0 {
					c.Set(key, j)
				} else {
					c.Get(key)
				}
			}
		}(i)
	}
	wg.Wait()
}
