package dlist

import (
	"reflect"
	"testing"
)

func TestPushFrontOrder(t *testing.T) {
	l := NewList()
	l.PushFront(1)
	l.PushFront(2)
	l.PushFront(3)

	want := []any{3, 2, 1}
	if got := l.Values(); !reflect.DeepEqual(got, want) {
		t.Fatalf("ожидали %v, получили %v", want, got)
	}
}

func TestBack(t *testing.T) {
	l := NewList()
	if l.Back() != nil {
		t.Fatal("Back() пустого списка должен вернуть nil")
	}

	l.PushFront(1)
	l.PushFront(2)
	if b := l.Back(); b == nil || b.Value != 1 {
		t.Fatalf("Back() должен вернуть узел со значением 1, получили %v", b)
	}
}

func TestRemove(t *testing.T) {
	l := NewList()
	l.PushFront(1)
	mid := l.PushFront(2)
	l.PushFront(3)

	l.Remove(mid)
	want := []any{3, 1}
	if got := l.Values(); !reflect.DeepEqual(got, want) {
		t.Fatalf("после удаления середины ожидали %v, получили %v", want, got)
	}

	l.Remove(l.Back())
	want = []any{3}
	if got := l.Values(); !reflect.DeepEqual(got, want) {
		t.Fatalf("после удаления хвоста ожидали %v, получили %v", want, got)
	}
}

func TestMoveToFront(t *testing.T) {
	l := NewList()
	a := l.PushFront("a")
	l.PushFront("b")
	l.PushFront("c") // c b a

	l.MoveToFront(a)
	want := []any{"a", "c", "b"}
	if got := l.Values(); !reflect.DeepEqual(got, want) {
		t.Fatalf("ожидали %v, получили %v", want, got)
	}

	// Перемещение уже первого узла ничего не ломает.
	l.MoveToFront(a)
	if got := l.Values(); !reflect.DeepEqual(got, want) {
		t.Fatalf("повторный MoveToFront первого узла: ожидали %v, получили %v", want, got)
	}
}

func TestEmptyAfterRemovals(t *testing.T) {
	l := NewList()
	n := l.PushFront(1)
	l.Remove(n)

	if l.Back() != nil {
		t.Fatal("после удаления единственного узла список должен быть пуст")
	}
	if got := l.Values(); len(got) != 0 {
		t.Fatalf("Values() пустого списка: ожидали пусто, получили %v", got)
	}

	// Список остаётся рабочим после опустошения.
	l.PushFront(7)
	if b := l.Back(); b == nil || b.Value != 7 {
		t.Fatal("список сломался после полного опустошения")
	}
}
