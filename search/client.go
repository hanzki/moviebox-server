// Package search allows searching for movies and TV series
package search

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/hanzki/moviebox-server/core"
)

var (
	Client HTTPClient
)

func init() {
	Client = &http.Client{}
}

type JackettClient struct {
	host       string
	port       string
	apikey     string
	categories map[string]*JackettCategory
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func NewJackettClient(host, port, apikey string) (*JackettClient, error) {
	categories, err := getCategories(host, port, apikey)
	if err != nil {
		return nil, err
	}
	return &JackettClient{host, port, apikey, categories}, nil
}

func torznabAPIurl(host, port string) string {
	return fmt.Sprintf("http://%s:%s/api/v2.0/indexers/all/results/torznab", host, port)
}

func getCategories(host, port, apikey string) (map[string]*JackettCategory, error) {
	req := buildRequest(torznabAPIurl(host, port), apikey, "caps", "", nil)
	xml, err := getXML(req)
	if err != nil {
		panic(err)
	}

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
	req := buildRequest(torznabAPIurl(jc.host, jc.port), jc.apikey, "search", query, jc.categories["2000"])
	xml, err := getXML(req)
	if err != nil {
		panic(err)
	}

	results := parseTorznab(xml, jc.categories)

	return results
}

func (jc *JackettClient) SearchTVSeries(query string) []*core.SearchResult {
	req := buildRequest(torznabAPIurl(jc.host, jc.port), jc.apikey, "search", query, jc.categories["5000"])
	xml, err := getXML(req)
	if err != nil {
		panic(err)
	}

	results := parseTorznab(xml, jc.categories)

	return results
}

func buildRequest(url, apikey, queryType, query string, category *JackettCategory) *http.Request {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	q := req.URL.Query()
	q.Add("apikey", apikey)
	q.Add("t", queryType)
	if queryType == "search" {
		q.Add("q", query)
		if category != nil {
			q.Add("cat", catString(category))
		}
	}
	req.URL.RawQuery = q.Encode()
	return req
}

func catString(cat *JackettCategory) string {
	catIDs := []string{cat.ID}
	for _, v := range cat.SubCategories {
		catIDs = append(catIDs, v.ID)
	}
	return strings.Join(catIDs, ",")
}

func getXML(req *http.Request) ([]byte, error) {
	resp, err := Client.Do(req)
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
