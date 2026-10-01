package parens

import "testing"

func TestIsValid(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"()", true},
		{"()[]{}", true},
		{"([])", true},
		{"{[()]}", true},
		{"", true},
		{"(]", false},
		{"([)]", false},
		{"(", false},
		{")", false},
		{"(((", false},
		{"())", false},
	}
	for _, tc := range cases {
		if got := IsValid(tc.in); got != tc.want {
			t.Errorf("IsValid(%q) = %v, ожидали %v", tc.in, got, tc.want)
		}
	}
}
