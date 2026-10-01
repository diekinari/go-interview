package lfu

import "testing"

func TestBasic(t *testing.T) {
	c := NewLFU(2)
	c.Set("a", 1)

	if v, ok := c.Get("a"); !ok || v != 1 {
		t.Fatalf("ожидали (1, true), получили (%v, %v)", v, ok)
	}
}

func TestEvictsLeastFrequent(t *testing.T) {
	c := NewLFU(2)
	c.Set("a", 1)
	c.Set("b", 2)

	c.Get("a") // freq(a)=2, freq(b)=1
	c.Set("c", 3) // вытесняется "b" — у него минимальная частота

	if _, ok := c.Get("b"); ok {
		t.Fatal("\"b\" с минимальной частотой должен был быть вытеснен")
	}
	if _, ok := c.Get("a"); !ok {
		t.Fatal("\"a\" с большей частотой должен был остаться")
	}
}

func TestLRUAmongEqualFreq(t *testing.T) {
	c := NewLFU(2)
	c.Set("a", 1)
	c.Set("b", 2)
	// freq обоих = 1; "a" старше по использованию
	c.Set("c", 3) // вытесняется "a" — LRU среди минимальной частоты

	if _, ok := c.Get("a"); ok {
		t.Fatal("при равной частоте должен вытесняться наименее недавно использованный (\"a\")")
	}
	if _, ok := c.Get("b"); !ok {
		t.Fatal("\"b\" должен был остаться")
	}
}

func TestUpdateKeepsFreq(t *testing.T) {
	c := NewLFU(2)
	c.Set("a", 1)
	c.Set("b", 2)

	c.Get("b")     // freq(b)=2
	c.Set("a", 10) // обновление НЕ меняет частоту: freq(a) всё ещё 1
	c.Set("c", 3)  // вытесняется "a"

	if _, ok := c.Get("a"); ok {
		t.Fatal("Set по существующему ключу не должен поднимать частоту — \"a\" должен был быть вытеснен")
	}
	if v, ok := c.Get("b"); !ok || v != 2 {
		t.Fatal("\"b\" должен был остаться")
	}
}

func TestMinFreqAdvances(t *testing.T) {
	c := NewLFU(2)
	c.Set("a", 1)
	c.Set("b", 2)
	c.Get("a")
	c.Get("b") // freq обоих = 2, minFreq должен подняться

	c.Set("c", 3) // вытесняется "a" (LRU среди freq=2)

	if _, ok := c.Get("a"); ok {
		t.Fatal("\"a\" должен был быть вытеснен")
	}
	// "c" с freq=1 теперь минимальный
	c.Set("d", 4)
	if _, ok := c.Get("c"); ok {
		t.Fatal("\"c\" с freq=1 должен был быть вытеснен следующим")
	}
	if _, ok := c.Get("b"); !ok {
		t.Fatal("\"b\" должен был остаться")
	}
}
