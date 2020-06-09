// Package search allows searching for movies and TV series
package search

import (
	"time"
)

// Result represents a single item in search results
type Result struct {
	title    string
	guid     string
	size     int
	pubDate  time.Time
	link     string
	category []string
	seeds    int
	peers    int
}

// Movies queries the search service and returns slice of results
func Movies(query string) []Result {
	return nil
}
