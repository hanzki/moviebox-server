package search

import "testing"

func TestMovies(t *testing.T) {
	cases := []struct {
		in   string
		want int
	}{
		{"", 0},
	}
	for _, c := range cases {
		got := len(Movies(c.in))
		if got != c.want {
			t.Errorf("len(Movies(%q)) == %q, want %q", c.in, got, c.want)
		}
	}
}
