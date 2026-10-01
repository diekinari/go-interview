package mergech

import (
	"sort"
	"testing"
	"time"
)

func feed(values ...int) <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for _, v := range values {
			ch <- v
		}
	}()
	return ch
}

func TestMergeAllValues(t *testing.T) {
	out := Merge(feed(1, 2, 3), feed(4, 5), feed(6))

	var got []int
	done := make(chan struct{})
	go func() {
		defer close(done)
		for v := range out {
			got = append(got, v)
		}
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("range по результату не завершился — выходной канал не закрыт")
	}

	sort.Ints(got)
	want := []int{1, 2, 3, 4, 5, 6}
	if len(got) != len(want) {
		t.Fatalf("ожидали %d значений, получили %d: %v", len(want), len(got), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("ожидали %v, получили %v", want, got)
		}
	}
}

func TestMergeSingle(t *testing.T) {
	out := Merge(feed(42))
	v, ok := <-out
	if !ok || v != 42 {
		t.Fatalf("ожидали 42, получили (%v, %v)", v, ok)
	}
	if _, ok := <-out; ok {
		t.Fatal("после единственного значения канал должен быть закрыт")
	}
}

func TestMergeEmpty(t *testing.T) {
	out := Merge()
	select {
	case _, ok := <-out:
		if ok {
			t.Fatal("из Merge() без каналов не должно приходить значений")
		}
	case <-time.After(time.Second):
		t.Fatal("Merge() без каналов должен вернуть закрытый канал, а не блокировать")
	}
}

func TestMergeSlowProducers(t *testing.T) {
	slow := make(chan int)
	go func() {
		defer close(slow)
		time.Sleep(50 * time.Millisecond)
		slow <- 100
	}()

	out := Merge(feed(1), slow)
	count := 0
	for range out {
		count++
	}
	if count != 2 {
		t.Fatalf("ожидали 2 значения, получили %d — Merge не дождался медленного канала", count)
	}
}
