package search

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/hanzki/moviebox-server/core"
	. "github.com/smartystreets/goconvey/convey"
)

func TestParseTorznab(t *testing.T) {
	Convey("Given simple torznab response", t, func() {
		torznabResponse := loadFile("testdata/torznab_sample.txt")
		categories := simpleCategories()

		Convey("Should parse without errors", func() {
			So(func() { parseTorznab(torznabResponse, categories) }, ShouldNotPanic)
		})

		Convey("Should return parsed search result", func() {
			pubDate, err := time.Parse(time.RFC1123Z, "Sun, 01 Jun 2008 00:00:00 +0000")
			check(err)
			expected := []*core.SearchResult{
				{
					Title:       "Big Buck Bunny (1280x720 msmp4)",
					IndexerGUID: "http://www.legittorrents.info/index.php?page=torrent-details&id=d42644b1a0eb635a4ffaf6ce17042d71428e6a8e",
					Indexer:     "legittorrents",
					PubDate:     pubDate,
					Link:        "http://localhost:9117/dl/legittorrents/?jackett_apikey=gpowobdo7ztigmxjoeokamhjjh7bz8us&path=Q2ZESjhPcWNEVnJwNmdWQmhqNDh0dGcxM2VtcXRJM0c3WGhuOUZCZ3lSeVRyWjkxWGhheThGd3VMWUc3RnM0d0xvWkhHMk5meDRRY3VCRk02TlExSjF5Tl81X0thTWY5ZzlWSE9xU1I5SlBIdjk2MEJIbXNHSElKdXJRMFdFVF81M1hETGNkaFM2MmFRaTU3aTczT0hsUm8zRE1zd2hRY1AtcHd5a1prSTNzdDhRcWp4bm5lTm1CVEREdEdSb0MyOEE0OEY1WnZZVTYwOW4zWWw4Zkh4QXpJa3BpbGtGT1VDa2FqcnNBWVlpRGY3N3IzLWdQQ1JKc0dkNld2TGlIbzNQLUZtcUdHcDYyRkdyVVNVZmplV1dWc3IyeDNhLXA1MjhVNXRmQXZZbk5OeERHdg&file=Big+Buck+Bunny+(1280x720+msmp4)",
					Categories:  []string{"Movies"},
					Seeds:       59,
					Peers:       59,
				},
			}
			result := parseTorznab(torznabResponse, categories)
			So(result, ShouldResemble, expected)
		})
	})
}

func loadFile(filename string) []byte {
	contents, err := ioutil.ReadFile(filename)
	check(err)
	return contents
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func simpleCategories() map[string]*JackettCategory {
	return map[string]*JackettCategory{"2000": {"2000", "Movies", nil}}
}

func TestParseCaps(t *testing.T) {

	Convey("Given caps response", t, func() {
		capsResponse := loadFile("testdata/torznab_caps_sample.txt")
		Convey("Should parse without errors", func() {
			So(func() { parseCaps(capsResponse) }, ShouldNotPanic)
		})

		Convey("Should return parsed categories", func() {
			expected := &JackettCategory{
				ID:   "5000",
				Name: "TV",
				SubCategories: []*JackettCategory{
					{"5010", "TV/WEB-DL", nil},
					{"5020", "TV/FOREIGN", nil},
					{"5030", "TV/SD", nil},
					{"5040", "TV/HD", nil},
					{"5045", "TV/UHD", nil},
					{"5050", "TV/OTHER", nil},
					{"5060", "TV/Sport", nil},
					{"5070", "TV/Anime", nil},
					{"5080", "TV/Documentary", nil},
				},
			}
			results := parseCaps(capsResponse)
			So(expected, ShouldBeIn, results)
		})
	})
}
