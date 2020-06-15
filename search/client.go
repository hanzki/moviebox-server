// Package search allows searching for movies and TV series
package search

import (
	"fmt"
	"io/ioutil"
	"net/http"
	urlUtils "net/url"
	"strings"

	"github.com/hanzki/moviebox-server/core"
)

// Movies queries the search service and returns slice of results
func Movies(query string) []core.SearchResult {
	return nil
}

type JackettClient struct {
	host       string
	port       string
	apikey     string
	categories map[string]*JackettCategory
}

func NewJackettClient(host, port, apikey string) (*JackettClient, error) {
	categories, err := getCategories(host, port, apikey)
	if err != nil {
		return nil, err
	}
	return &JackettClient{host, port, apikey, categories}, nil
}

func getCategories(host, port, apikey string) (map[string]*JackettCategory, error) {
	url := fmt.Sprintf("http://%s:%s/api/v2.0/indexers/all/results/torznab", host, port)
	queryString := fmt.Sprintf("?apikey=%s&t=caps", apikey)
	fullURI := url + queryString

	xml, err := getXML(fullURI)
	_ = err //TODO: Check for errors

	jackettCategories := parseCaps(xml)

	if len(jackettCategories) < 1 {
		return nil, fmt.Errorf("getCategories: Couldn't get category information")
	}

	categories := make(map[string]*JackettCategory)
	for _, v := range jackettCategories {
		categories[v.ID] = v
	}

	return categories, nil
}

func (jc *JackettClient) SearchMovies(query string) []*core.SearchResult {
	url := fmt.Sprintf("http://%s:%s/api/v2.0/indexers/all/results/torznab", jc.host, jc.port)
	queryString := fmt.Sprintf("?apikey=%s&t=search&q=%s&cat=%s", jc.apikey, urlUtils.QueryEscape(query), catString(jc.categories["2000"]))

	fullURI := url + queryString
	fmt.Println(fullURI)

	xml, err := getXML(fullURI)
	_ = err //TODO: Check for errors

	results := parseTorznab(xml, jc.categories)

	return results
}

func (jc *JackettClient) SearchTVSeries(query string) []*core.SearchResult {
	url := fmt.Sprintf("http://%s:%s/api/v2.0/indexers/all/results/torznab", jc.host, jc.port)
	queryString := fmt.Sprintf("?apikey=%s&t=search&q=%s&cat=%s", jc.apikey, urlUtils.QueryEscape(query), catString(jc.categories["5000"]))

	fullURI := url + queryString
	fmt.Println(fullURI)

	xml, err := getXML(fullURI)
	_ = err //TODO: Check for errors

	results := parseTorznab(xml, jc.categories)

	return results
}

func catString(cat *JackettCategory) string {
	catIDs := []string{cat.ID}
	for _, v := range cat.SubCategories {
		catIDs = append(catIDs, v.ID)
	}
	return strings.Join(catIDs, ",")
}

func getXML(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, fmt.Errorf("GET error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("Status error: %v", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("Read body: %v", err)
	}

	return data, nil
}
