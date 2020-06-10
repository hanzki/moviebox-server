package main

import (
	"fmt"
	"time"
)

type SearchResult struct {
	title       string
	indexerGUID string
	indexer     string
	size        int
	pubDate     time.Time
	link        string
	categories  []string
	seeds       int
	peers       int
}

type SearchResponse struct {
	id         string
	title      string
	size       int
	pubDate    time.Time
	categories []string
	seeds      int
	peers      int
}

type SearchRecord struct {
	id          string
	createdAt   time.Time
	title       string
	indexerGUID string
	indexer     string
	size        int
	pubDate     time.Time
	link        string
	categories  []string
	seeds       int
	peers       int
}

type SearchStorage interface {
	Save(r *SearchRecord) (*SearchRecord, error)
	Update(r *SearchRecord) error
	Load(id string) (*SearchRecord, error)
	Delete(id string) error
}

type SearchClient interface {
	SearchMovies(query string) []*SearchResult
	SearchTVSeries(query string) []*SearchResult
}

type SearchController struct {
	storage SearchStorage
	client  SearchClient
}

func (sc *SearchController) searchMovies(query string) []*SearchResponse {
	results := sc.client.SearchMovies(query)

	responses := make([]*SearchResponse, 0, len(results))
	for i, result := range results {
		record := searchResultToRecord(result)
		if record, err := sc.storage.Save(record); err != nil {
			fmt.Println("Save failed :( $v", record)
			panic(err)
		}
		responses[i] = searchRecordToResponse(record)
	}

	return responses
}

func searchResultToRecord(result *SearchResult) *SearchRecord {
	return &SearchRecord{
		title:       result.title,
		indexerGUID: result.indexerGUID,
		indexer:     result.indexer,
		size:        result.size,
		pubDate:     result.pubDate,
		link:        result.link,
		categories:  result.categories,
		seeds:       result.seeds,
		peers:       result.peers,
	}
}

func searchRecordToResponse(record *SearchRecord) *SearchResponse {
	return &SearchResponse{
		id:         record.id,
		title:      record.title,
		size:       record.size,
		pubDate:    record.pubDate,
		categories: record.categories,
		seeds:      record.seeds,
		peers:      record.peers,
	}
}
