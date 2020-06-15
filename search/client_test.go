package search

import "testing"

func TestMovies(t *testing.T) {
	jc := &JackettClient{
		"localhost",
		"9117",
		"gpowobdo7ztigmxjoeokamhjjh7bz8us",
		map[string]*JackettCategory{"2000": {"2000", "Movies", nil}},
	}

	jc.SearchMovies("big buck bunny")
}
