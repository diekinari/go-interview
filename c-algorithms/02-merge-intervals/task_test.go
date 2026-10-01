package intervals

import (
	"reflect"
	"testing"
)

func TestMerge(t *testing.T) {
	cases := []struct {
		name string
		in   [][]int
		want [][]int
	}{
		{"классика", [][]int{{1, 3}, {2, 6}, {8, 10}, {15, 18}}, [][]int{{1, 6}, {8, 10}, {15, 18}}},
		{"касание", [][]int{{1, 4}, {4, 5}}, [][]int{{1, 5}}},
		{"не отсортировано", [][]int{{8, 10}, {1, 3}, {2, 6}}, [][]int{{1, 6}, {8, 10}}},
		{"вложенный", [][]int{{1, 10}, {2, 3}, {4, 5}}, [][]int{{1, 10}}},
		{"один", [][]int{{1, 2}}, [][]int{{1, 2}}},
		{"пусто", [][]int{}, [][]int{}},
		{"всё в один", [][]int{{1, 2}, {2, 3}, {3, 4}}, [][]int{{1, 4}}},
	}
	for _, tc := range cases {
		got := Merge(tc.in)
		if len(got) == 0 && len(tc.want) == 0 {
			continue
		}
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("%s: Merge(%v) = %v, ожидали %v", tc.name, tc.in, got, tc.want)
		}
	}
}
