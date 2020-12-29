package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/hanzki/moviebox-server/core"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "hello world"}`))
}

type JsonSearchResult struct {
	ID         core.SearchID `json:"id"`
	Title      string        `json:"title"`
	Indexer    string        `json:"indexer"`
	Size       int           `json:"size"`
	PubDate    time.Time     `json:"published"`
	Categories []string      `json:"categories"`
	Seeds      int           `json:"seeds"`
	Peers      int           `json:"peers"`
}

func searchResultToJson(sr *core.SearchResult) *JsonSearchResult {
	return &JsonSearchResult{
		sr.ID,
		sr.Title,
		sr.Indexer,
		sr.Size,
		sr.PubDate,
		sr.Categories,
		sr.Seeds,
		sr.Peers,
	}
}

func (s *server) searchTorrents(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	log.Printf("search query=%s", query)

	results := s.searchCtrl.SearchMovies(query)

	jsonResults := make([]*JsonSearchResult, len(results))
	for i, sr := range results {
		jsonResults[i] = searchResultToJson(sr)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(jsonResults)
}

func (s *server) getSearchResult(w http.ResponseWriter, r *http.Request) {
	searchID := core.SearchID(mux.Vars(r)["searchID"])
	log.Printf("search-get searchID=%s", searchID)

	searchResult, err := s.searchCtrl.Storage.Load(searchID)
	if err != nil {
		panic(err) //TODO: Handle error
	}

	jsonResult := searchResultToJson(searchResult)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(jsonResult)
}

type JsonDownload struct {
	ID       core.DownloadID     `json:"id"`
	SearchID core.SearchID       `json:"searchId"`
	Status   core.DownloadStatus `json:"status"`
	Progress float64             `json:"progress"`
}

func downloadToJsonDownload(d *core.Download) *JsonDownload {
	return &JsonDownload{
		d.ID,
		d.SearchID,
		d.Status,
		d.Progress,
	}
}

func (s *server) downloadTorrent(w http.ResponseWriter, r *http.Request) {
	searchID := core.SearchID(mux.Vars(r)["searchID"])
	log.Printf("download searchID=%s", searchID)

	searchResult, err := s.searchCtrl.Storage.Load(searchID)
	if err != nil {
		panic(err) //TODO: Handle error
	}

	download, err := s.downloadCtrl.StartNewDownload(searchResult)
	if err != nil {
		panic(err) //TODO: Handle error
	}
	jsonDownload := downloadToJsonDownload(download)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(jsonDownload)
}

func (s *server) getDownload(w http.ResponseWriter, r *http.Request) {
	downloadID := core.DownloadID(mux.Vars(r)["downloadID"])
	log.Printf("download-get downloadID=%s", downloadID)

	download, err := s.downloadCtrl.GetProgress(downloadID)
	if err != nil {
		panic(err) //TODO: Handle error
	}

	jsonDownload := downloadToJsonDownload(download)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(jsonDownload)
}

type server struct {
	searchCtrl   *core.SearchController
	downloadCtrl *core.DownloadController
}

func startServer(searchCtrl *core.SearchController, downloadCtrl *core.DownloadController) {
	s := &server{searchCtrl, downloadCtrl}

	r := mux.NewRouter()
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/search", s.searchTorrents).Methods(http.MethodGet)
	api.HandleFunc("/search/{searchID}", s.getSearchResult).Methods(http.MethodGet)
	api.HandleFunc("/search/{searchID}/download", s.downloadTorrent).Methods(http.MethodPost)
	api.HandleFunc("/download/{downloadID}", s.getDownload).Methods(http.MethodGet)

	log.Fatal(http.ListenAndServe(":8080", r))
}
