package search

import (
	"encoding/xml"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/hanzki/moviebox-server/core"
)

type XMLJackettIndexer struct {
	ID string `xml:"id,attr"`
}

type XMLTorznabAttr struct {
	Key   string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type XMLItem struct {
	Title          string            `xml:"title"`
	GUID           string            `xml:"guid"`
	JackettIndexer XMLJackettIndexer `xml:"jackettindexer"`
	PubDate        string            `xml:"pubDate"`
	Categories     []string          `xml:"category"`
	Link           string            `xml:"link"`
	Size           int               `xml:"size"`
	TorznabAttrs   []XMLTorznabAttr  `xml:"http://torznab.com/schemas/2015/feed attr"`
}

type XMLChannel struct {
	Items []XMLItem `xml:"channel>item"`
}

func parseTorznab(xmlBytes []byte, categories map[string]*JackettCategory) []*core.SearchResult {
	result := XMLChannel{}
	xml.Unmarshal(xmlBytes, &result)
	results := make([]*core.SearchResult, 0, 10)
	for _, item := range result.Items {
		sr := core.SearchResult{}
		sr.Title = item.Title
		sr.IndexerGUID = item.GUID
		sr.Indexer = item.JackettIndexer.ID
		sr.Link = item.Link
		sr.Size = item.Size

		pubDate, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			panic(err)
		}
		sr.PubDate = pubDate

		if seedsStr, ok := findAttr(item.TorznabAttrs, "seeders"); ok {
			if seeds, err := strconv.Atoi(seedsStr); err != nil {
				fmt.Fprintf(os.Stderr, "parseTorznab: Unable to parse seeders value %s, %v", seedsStr, err)
			} else {
				sr.Seeds = seeds
			}
		}

		if peersStr, ok := findAttr(item.TorznabAttrs, "peers"); ok {
			if peers, err := strconv.Atoi(peersStr); err != nil {
				fmt.Fprintf(os.Stderr, "parseTorznab: Unable to parse peers value %s, %v", peersStr, err)
			} else {
				sr.Peers = peers
			}
		}

		cats := []string{}
		for _, v := range item.Categories {
			if jackettCategory, ok := categories[v]; ok {
				cats = append(cats, jackettCategory.Name)
			}
		}
		sr.Categories = cats

		results = append(results, &sr)
	}
	return results
}

func findAttr(attrs []XMLTorznabAttr, name string) (string, bool) {
	for _, v := range attrs {
		if v.Key == name {
			return v.Value, true
		}
	}
	return "", false
}

type JackettCategory struct {
	ID            string             `xml:"id,attr"`
	Name          string             `xml:"name,attr"`
	SubCategories []*JackettCategory `xml:"subcat"`
}

type XMLCategories struct {
	Categories []*JackettCategory `xml:"categories>category"`
}

func parseCaps(xmlBytes []byte) []*JackettCategory {
	result := XMLCategories{}
	xml.Unmarshal(xmlBytes, &result)
	return result.Categories
}
