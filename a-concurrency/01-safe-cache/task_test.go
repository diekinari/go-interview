package safecache

import (
	"fmt"
	"sync"
	"testing"
)

func TestBasic(t *testing.T) {
	c := NewCache()

	if _, ok := c.Get("missing"); ok {
		t.Fatal("Get по несуществующему ключу должен вернуть ok=false")
	}

	c.Set("a", 1)
	v, ok := c.Get("a")
	if !ok || v != 1 {
		t.Fatalf("ожидали (1, true), получили (%v, %v)", v, ok)
	}

	c.Set("a", 2)
	if v, _ := c.Get("a"); v != 2 {
		t.Fatalf("после перезаписи ожидали 2, получили %v", v)
	}

	c.Delete("a")
	if _, ok := c.Get("a"); ok {
		t.Fatal("после Delete ключ не должен находиться")
	}
}

func TestConcurrent(t *testing.T) {
	c := NewCache()
	var wg sync.WaitGroup

	// 10 писателей, 10 читателей, 5 удаляющих — гоняем одновременно.
	// Смысл теста раскрывается флагом -race.
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				c.Set(fmt.Sprintf("key-%d", j%20), n*1000+j)
			}
		}(i)
	}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				c.Get(fmt.Sprintf("key-%d", j%20))
			}
		}()
	}
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 50; j++ {
				c.Delete(fmt.Sprintf("key-%d", j%20))
			}
		}()
	}
	wg.Wait()

	// После гонки кэш должен остаться рабочим.
	c.Set("final", 42)
	if v, ok := c.Get("final"); !ok || v != 42 {
		t.Fatal("кэш сломался после конкурентной нагрузки")
	}
}
