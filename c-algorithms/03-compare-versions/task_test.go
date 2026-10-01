package versions

import "testing"

func TestCompare(t *testing.T) {
	cases := []struct {
		v1, v2 string
		want   int
	}{
		{"1.2", "1.10", -1},
		{"1.01", "1.001", 0},
		{"1.0", "1.0.0.0", 0},
		{"2.1", "1.9", 1},
		{"1.0.1", "1", 1},
		{"7.5.2.4", "7.5.3", -1},
		{"1", "1", 0},
		{"0.1", "1.1", -1},
	}
	for _, tc := range cases {
		if got := Compare(tc.v1, tc.v2); got != tc.want {
			t.Errorf("Compare(%q, %q) = %d, ожидали %d", tc.v1, tc.v2, got, tc.want)
		}
	}
}
