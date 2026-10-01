package mergelists

import (
	"reflect"
	"testing"
)

// build собирает список из значений, начиная с хвоста.
func build(vals ...int) *ListNode {
	var head *ListNode
	for i := len(vals) - 1; i >= 0; i-- {
		head = &ListNode{Val: vals[i], Next: head}
	}
	return head
}

func toSlice(l *ListNode) []int {
	out := []int{}
	for ; l != nil; l = l.Next {
		out = append(out, l.Val)
	}
	return out
}

func nodes(l *ListNode) map[*ListNode]bool {
	set := map[*ListNode]bool{}
	for ; l != nil; l = l.Next {
		set[l] = true
	}
	return set
}

func TestMerge(t *testing.T) {
	cases := []struct {
		name   string
		l1, l2 *ListNode
		want   []int
	}{
		{"классика", build(1, 2, 4), build(1, 3, 4), []int{1, 1, 2, 3, 4, 4}},
		{"оба пустые", nil, nil, []int{}},
		{"первый пустой", nil, build(0), []int{0}},
		{"второй пустой", build(1, 2), nil, []int{1, 2}},
		{"без пересечений", build(1, 2, 3), build(10, 20), []int{1, 2, 3, 10, 20}},
		{"один длиннее", build(5), build(1, 2, 3, 4), []int{1, 2, 3, 4, 5}},
		{"дубликаты", build(2, 2), build(2, 2), []int{2, 2, 2, 2}},
	}
	for _, tc := range cases {
		got := toSlice(Merge(tc.l1, tc.l2))
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("%s: получили %v, ожидали %v", tc.name, got, tc.want)
		}
	}
}

func TestReusesNodes(t *testing.T) {
	l1, l2 := build(1, 2, 4), build(1, 3, 4)
	input := nodes(l1)
	for n := range nodes(l2) {
		input[n] = true
	}

	count := 0
	for n := Merge(l1, l2); n != nil; n = n.Next {
		if !input[n] {
			t.Fatalf("в результате узел %v, которого не было во входных списках — узлы нужно переиспользовать", n.Val)
		}
		count++
	}
	if count != len(input) {
		t.Fatalf("в результате %d узлов, во входных было %d", count, len(input))
	}
}
