package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/hanzki/moviebox-server/storage"

	"github.com/hanzki/moviebox-server/core"

	"github.com/davecgh/go-spew/spew"
	"github.com/hanzki/moviebox-server/search"
)

func main() {
	apikeyPtr := flag.String("apikey", "", "Jackett apikey")
	queryPtr := flag.String("query", "", "Search query")
	hostPtr := flag.String("host", "localhost", "Jackett host")
	portPtr := flag.String("port", "9117", "Jackett port")
	typePtr := flag.String("type", "movies", "Search type (movies/tv)")
	flag.Parse()

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
