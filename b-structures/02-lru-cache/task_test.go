package lru

import "testing"

func TestBasicSetGet(t *testing.T) {
	c := NewLRU(2)
	c.Set("a", 1)

	if v, ok := c.Get("a"); !ok || v != 1 {
		t.Fatalf("ожидали (1, true), получили (%v, %v)", v, ok)
	}
	if _, ok := c.Get("b"); ok {
		t.Fatal("несуществующий ключ должен дать ok=false")
	}
}

func TestEvictsOldest(t *testing.T) {
	c := NewLRU(2)
	c.Set("a", 1)
	c.Set("b", 2)
	c.Set("c", 3) // вытесняет "a" — самый старый

	if _, ok := c.Get("a"); ok {
		t.Fatal("\"a\" должен был быть вытеснен")
	}
	if v, ok := c.Get("b"); !ok || v != 2 {
		t.Fatal("\"b\" должен был остаться")
	}
	if v, ok := c.Get("c"); !ok || v != 3 {
		t.Fatal("\"c\" должен был остаться")
	}
	if c.Len() != 2 {
		t.Fatalf("Len() = %d, ожидали 2", c.Len())
	}
}

func TestGetRefreshesRecency(t *testing.T) {
	c := NewLRU(2)
	c.Set("a", 1)
	c.Set("b", 2)

	c.Get("a")    // "a" теперь свежее "b"
	c.Set("c", 3) // вытесняется "b", а не "a"

	if _, ok := c.Get("b"); ok {
		t.Fatal("\"b\" должен был быть вытеснен: Get(\"a\") освежил \"a\"")
	}
	if _, ok := c.Get("a"); !ok {
		t.Fatal("\"a\" должен был остаться после освежения")
	}
}

func TestSetRefreshesRecency(t *testing.T) {
	c := NewLRU(2)
	c.Set("a", 1)
	c.Set("b", 2)

	c.Set("a", 10) // обновление тоже освежает
	c.Set("c", 3)  // вытесняется "b"

	if _, ok := c.Get("b"); ok {
		t.Fatal("\"b\" должен был быть вытеснен")
	}
	if v, ok := c.Get("a"); !ok || v != 10 {
		t.Fatalf("\"a\" должен был остаться со значением 10, получили (%v, %v)", v, ok)
	}
}

func TestUpdateDoesNotGrow(t *testing.T) {
	c := NewLRU(2)
	c.Set("a", 1)
	c.Set("a", 2)
	c.Set("a", 3)

	if c.Len() != 1 {
		t.Fatalf("обновления одного ключа не должны растить кэш: Len() = %d", c.Len())
	}
}

// Сценарий из LeetCode 146 целиком.
func TestLeetcodeScenario(t *testing.T) {
	c := NewLRU(2)
	c.Set(1, 1)
	c.Set(2, 2)
	if v, _ := c.Get(1); v != 1 {
		t.Fatal("Get(1) -> 1")
	}
	c.Set(3, 3) // вытесняет 2
	if _, ok := c.Get(2); ok {
		t.Fatal("Get(2) должен промахнуться")
	}
	c.Set(4, 4) // вытесняет 1
	if _, ok := c.Get(1); ok {
		t.Fatal("Get(1) должен промахнуться")
	}
	if v, _ := c.Get(3); v != 3 {
		t.Fatal("Get(3) -> 3")
	}
	if v, _ := c.Get(4); v != 4 {
		t.Fatal("Get(4) -> 4")
	}
}
