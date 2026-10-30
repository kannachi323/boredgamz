package identity

import "testing"

func TestSameAccount(t *testing.T) {
	cases := []struct {
		a, b string
		want bool
	}{
		{"u1:tabA", "u1:tabB", true},
		{"u1:tabA", "u2:tabA", false},
		{"one", "two", false},
		{"one", "one", true},
		{"u1:tabA:spectator", "u1:tabA", true},
	}
	for _, c := range cases {
		if got := SameAccount(c.a, c.b); got != c.want {
			t.Errorf("SameAccount(%q, %q) = %v, want %v", c.a, c.b, got, c.want)
		}
	}
}
