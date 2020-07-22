package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/hanzki/moviebox-server/download"
	"github.com/hanzki/moviebox-server/storage"

	"github.com/hanzki/moviebox-server/core"

	"github.com/davecgh/go-spew/spew"
	"github.com/hanzki/moviebox-server/search"
)

var (
	apikeyPtr *string
	queryPtr  *string
	hostPtr   *string
	portPtr   *string
	typePtr   *string
)

func main() {
	fmt.Println("Moviebox server v0.1.0")
	apikeyPtr = flag.String("apikey", "", "Jackett apikey")
	queryPtr = flag.String("query", "", "Search query")
	hostPtr = flag.String("host", "localhost", "Jackett host")
	portPtr = flag.String("port", "9117", "Jackett port")
	typePtr = flag.String("type", "movies", "Search type (movies/tv)")
	flag.Parse()

	args := flag.Args()

	var command string
	if len(args) > 0 {
		command = args[0]
		args = args[1:]
	} else {
		command = "search"
	}

	switch command {
	case "search":
		searchCmd()
	case "download":
		downloadCmd()
	default:
		log.Fatalf("moviebox: Unknown command %s", command)
	}

}

func searchCmd() {
	if *apikeyPtr == "" || *queryPtr == "" {
		fmt.Print("Missing required parameters.\n\n")
		flag.Usage()
		os.Exit(1)
	}

	jackettClient, err := search.NewJackettClient(
		*hostPtr,
		*portPtr,
		*apikeyPtr,
	)
	if err != nil {
		panic(err)
	}

	dummyStorage := storage.NewMock()

	searchController := &core.SearchController{
		Storage: dummyStorage,
		Client:  jackettClient,
	}

	var results []*core.SearchResult
	switch strings.ToLower(*typePtr) {
	case "movies", "movie":
		results = searchController.SearchMovies(*queryPtr)
	case "tv", "tvseries", "tv-series":
		results = searchController.SearchTVSeries(*queryPtr)
	default:
		fmt.Printf("Unknown search type %s.\n", *typePtr)
		flag.Usage()
		os.Exit(1)
	}
	spew.Dump(results)
}

func downloadCmd() {
	transmissionClient := download.NewTransmissionClient(&download.Config{
		RPCHost:     "localhost",
		RPCUser:     "transmission",
		RPCPassword: "transmission",
	})

	dummyStorage := storage.NewDownloadStorageMock()

	downloadController := &core.DownloadController{
		Storage: dummyStorage,
		Client:  transmissionClient,
	}

	sr := &core.SearchResult{
		ID:          "d8e2f06b-9fdd-4a7a-9712-1d01b8622417",
		CreatedAt:   time.Now(),
		Title:       "Big Buck Bunny (1920x1080 h.264)",
		IndexerGUID: "http://www.legittorrents.info/index.php?page=torrent-details&id=7f34612e0fac5e7b051b78bdf1060113350ebfe0",
		Indexer:     "legittorrents",
		Size:        0,
		PubDate:     time.Date(2008, time.June, 1, 0, 0, 0, 0, time.UTC),
		Link:        "http://localhost:9117/dl/legittorrents/?jackett_apikey=gpowobdo7ztigmxjoeokamhjjh7bz8us&path=Q2ZESjhPcWNEVnJwNmdWQmhqNDh0dGcxM2VrbmpmWjUtVGpaUHF6VVRmTGdTREs3N2ZMaXVDVmxSNkVPcWN6cThHS1o2ck50YUxhSXBjb0gyUk9jTXVWNHFQM1RCdHRhck92dU9sOUtTTDd5LVU4X3VheXJyaVB4WUxsWVlMNEZPN2R2VG1EQWVkaEEwS1V1M0EyS25UTGdGVGt4cXAyLUZ6bmhMMXRxWm9iM2hnSHIzemRkdldfX3h0bWQ1TWt1M0d5LXBqRjc2Y05QV0ExQ2ZXcVZ4TGxvM3JsYkwyVUJrYkVnTzdxaE90QmU1UGIzeVVKNXlWcjU3eDhXbm5PMFJIUGcwM1FWd1hPMks1bmNqY0tSak9QeUhvbXplby1ZYU1mTHFvOUNhVkoweEpEUg&file=Big+Buck+Bunny+(1920x1080+h.264)",
	}

	download, err := downloadController.StartNewDownload(sr)
	if err != nil {
		log.Fatal(err)
	}

	spew.Dump(download)
}
