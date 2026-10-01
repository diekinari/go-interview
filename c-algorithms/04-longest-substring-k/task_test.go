package ksubstr

import "testing"

func TestLongestKDistinct(t *testing.T) {
	cases := []struct {
		s    string
		k    int
		want int
	}{
		{"abbc", 2, 3},
		{"aaa", 1, 3},
		{"eceba", 2, 3},
		{"aa", 1, 2},
		{"abc", 3, 3},
		{"abc", 10, 3},
		{"abcabcbb", 2, 4}, // "bcbb"
		{"", 2, 0},
		{"abc", 0, 0},
	}
	for _, tc := range cases {
		if got := LongestKDistinct(tc.s, tc.k); got != tc.want {
			t.Errorf("LongestKDistinct(%q, %d) = %d, ожидали %d", tc.s, tc.k, got, tc.want)
		}
	}
}
