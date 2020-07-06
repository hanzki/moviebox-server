package search

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/hanzki/moviebox-server/utils/mocks"

	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	Client = &mocks.MockClient{}
}

func TestNewJackettClient(t *testing.T) {
	Convey("Given valid Jackett configuration", t, func() {
		host := "localhost"
		port := "9117"
		apikey := "secret"

		Convey("Should fetch categories from Jackett API", func() {
			mocks.GetDoFunc = mocks.SuccessFromFile("testdata/torznab_caps_sample.txt")

			jc, err := NewJackettClient(host, port, apikey)
			So(jc, ShouldNotBeNil)
			So(err, ShouldBeNil)
		})
	})
}

func TestSearchMovies(t *testing.T) {
	Convey("Given initialised JackettClient", t, func() {
		host := "localhost"
		port := "9117"
		apikey := "secret"

		mocks.GetDoFunc = mocks.SuccessFromFile("testdata/torznab_caps_sample.txt")

		jc, err := NewJackettClient(host, port, apikey)
		So(err, ShouldBeNil)

		Convey("Should search with category 2000 and subcategories", func() {
			mocks.GetDoFunc = func(req *http.Request) (*http.Response, error) {
				q := req.URL.Query()
				So(q.Get("cat"), ShouldEqual, "2000,2010,2020,2030,2040,2045,2050,2060,2070,2080")
				return &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(bytes.NewReader(nil)),
				}, nil
			}
			results := jc.SearchMovies("Big Buck Bunny")
			So(results, ShouldBeEmpty)
		})

		Convey("Should return parsed search results", func() {
			mocks.GetDoFunc = mocks.SuccessFromFile("testdata/torznab_sample.txt")

			results := jc.SearchMovies("Big Buck Bunny")
			So(len(results), ShouldEqual, 1)
			So(results[0].Title, ShouldEqual, "Big Buck Bunny (1280x720 msmp4)")
		})
	})
}

func TestSearchTVSeries(t *testing.T) {
	Convey("Given initialised JackettClient", t, func() {
		host := "localhost"
		port := "9117"
		apikey := "secret"

		mocks.GetDoFunc = mocks.SuccessFromFile("testdata/torznab_caps_sample.txt")

		jc, err := NewJackettClient(host, port, apikey)
		So(err, ShouldBeNil)

		Convey("Should search with category 5000 and subcategories", func() {
			mocks.GetDoFunc = func(req *http.Request) (*http.Response, error) {
				q := req.URL.Query()
				So(q.Get("cat"), ShouldEqual, "5000,5010,5020,5030,5040,5045,5050,5060,5070,5080")
				return &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(bytes.NewReader(nil)),
				}, nil
			}
			results := jc.SearchTVSeries("Big Buck Bunny")
			So(results, ShouldBeEmpty)
		})

		Convey("Should return parsed search results", func() {
			mocks.GetDoFunc = mocks.SuccessFromFile("testdata/torznab_sample.txt")

			results := jc.SearchTVSeries("Big Buck Bunny")
			So(len(results), ShouldEqual, 1)
			So(results[0].Title, ShouldEqual, "Big Buck Bunny (1280x720 msmp4)")
		})
	})
}
