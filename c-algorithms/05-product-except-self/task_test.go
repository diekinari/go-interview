package product

import (
	"reflect"
	"testing"
)

func TestProductExceptSelf(t *testing.T) {
	cases := []struct {
		in, want []int
	}{
		{[]int{1, 2, 3, 4}, []int{24, 12, 8, 6}},
		{[]int{2, 3}, []int{3, 2}},
		{[]int{1, 0, 3}, []int{0, 3, 0}},
		{[]int{0, 0, 1}, []int{0, 0, 0}},
		{[]int{5}, []int{1}},
		{[]int{-1, 2, -3}, []int{-6, 3, -2}},
	}
	for _, tc := range cases {
		if got := ProductExceptSelf(tc.in); !reflect.DeepEqual(got, tc.want) {
			t.Errorf("ProductExceptSelf(%v) = %v, ожидали %v", tc.in, got, tc.want)
		}
	}
}
