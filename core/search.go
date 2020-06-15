package core

import (
	"fmt"
	"sort"
	"time"
)

type SearchID string

// SearchResult represents a single item in search results
type SearchResult struct {
	ID          SearchID
	CreatedAt   time.Time
	Title       string
	IndexerGUID string
	Indexer     string
	Size        int
	PubDate     time.Time
	Link        string
	Categories  []string
	Seeds       int
	Peers       int
}

// SearchStorage allows storing and retrieving SearchResults from storage
type SearchStorage interface {
	Save(r *SearchResult) (*SearchResult, error)
	Update(r *SearchResult) error
	Load(id SearchID) (*SearchResult, error)
	Delete(id SearchID) error
}

// SearchClient provides methods for querying the search provider for Movies and TVSeries
type SearchClient interface {
	SearchMovies(query string) []*SearchResult
	SearchTVSeries(query string) []*SearchResult
}

// SearchController orchestrates the searching of movies and TV series
type SearchController struct {
	Storage SearchStorage
	Client  SearchClient
}

func (sc *SearchController) SearchMovies(query string) []*SearchResult {
	results := sc.Client.SearchMovies(query)
	results = pruneSearchResults(results, 5)
	results = storeResults(results, sc.Storage)
	return results
}

func (sc *SearchController) SearchTVSeries(query string) []*SearchResult {
	results := sc.Client.SearchTVSeries(query)
	results = pruneSearchResults(results, 5)
	results = storeResults(results, sc.Storage)
	return results
}

func storeResults(results []*SearchResult, storage SearchStorage) []*SearchResult {
	for i, result := range results {
		if result, err := storage.Save(result); err != nil {
			fmt.Println("Save failed :( $v", result)
			panic(err)
		}
		results[i] = result
	}
	return results
}

func pruneSearchResults(results []*SearchResult, n int) []*SearchResult {
	sort.Sort(sort.Reverse(bySeeds(results)))
	return results[:min(len(results), n)]
}

func min(a int, b int) int {
	if a <= b {
		return a
	}
	return b
}

type bySeeds []*SearchResult

func (s bySeeds) Len() int {
	return len(s)
}
func (s bySeeds) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s bySeeds) Less(i, j int) bool {
	return s[i].Seeds < s[j].Seeds
}
