package download

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/hanzki/moviebox-server/core"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		log.Fatal(err)
	}
	greeting, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", greeting)
}

func TestNewTransmissionClient(t *testing.T) {
	Convey("Should create new TransmissionClient", t, func() {
		config := validConfiguration()

		tc := NewTransmissionClient(config)
		So(tc, ShouldNotBeNil)
	})
}

type RPCRequest struct {
	Method    string                 `json:"method"`
	Tag       uint64                 `json:"tag"`
	Arguments map[string]interface{} `json:"arguments"`
}

type RPCResponse struct {
	Result    string                 `json:"result"`
	Tag       uint64                 `json:"tag"`
	Arguments map[string]interface{} `json:"arguments"`
}

func TestStartDownload(t *testing.T) {
	Convey("Given configured TransmissionClient", t, func() {
		tc, ts := transmissionClientWithTestServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			var result RPCRequest
			err := json.NewDecoder(r.Body).Decode(&result)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Test server called:\n%v\n", result)
			log.Printf("Request tag %v", result.Tag)
			log.Println(r.Header)
			response := RPCResponse{
				Result: "success",
				Tag:    result.Tag,
				Arguments: map[string]interface{}{
					"torrent-added": map[string]interface{}{
						"hashString": "2625bb07c6658bd906080ff3d7353c960c0d257a",
						"id":         1,
						"name":       "big_buck_bunny_480p_stereo.avi",
					},
				},
			}
			json.NewEncoder(w).Encode(response)
		}))
		defer ts.Close()

		Convey("Should start a new download", func() {
			searchResult := &core.SearchResult{}
			download, err := tc.StartDownload(searchResult)
			So(err, ShouldBeNil)
			So(download, ShouldNotBeNil)
			So(download.Status, ShouldEqual, core.Downloading)
			So(download.Hash, ShouldEqual, "2625bb07c6658bd906080ff3d7353c960c0d257a")
			log.Println(download)
		})

	})
}

func TestProgress(t *testing.T) {
	Convey("Given configured TransmissionClient", t, func() {
		tc, ts := transmissionClientWithTestServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			var result RPCRequest
			err := json.NewDecoder(r.Body).Decode(&result)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Test server called:\n%v\n", result)
			log.Printf("Request tag %v", result.Tag)
			log.Println(r.Header)
			response := RPCResponse{
				Result: "success",
				Tag:    result.Tag,
				Arguments: map[string]interface{}{
					"torrents": []map[string]interface{}{
						{
							"hashString":   "2625bb07c6658bd906080ff3d7353c960c0d257a",
							"id":           1,
							"name":         "big_buck_bunny_480p_stereo.avi",
							"totalSize":    893096902,
							"haveValid":    893096902,
							"sizeWhenDone": 893096902,
							"percentDone":  1,
							"isFinished":   false,
						},
					},
				},
			}
			json.NewEncoder(w).Encode(response)
		}))
		defer ts.Close()

		Convey("Return up to date progress of the download", func() {
			download := &core.Download{
				Hash:     "2625bb07c6658bd906080ff3d7353c960c0d257a",
				Status:   core.Downloading,
				Progress: 0,
			}
			download, err := tc.Progress(download)
			So(err, ShouldBeNil)
			So(download, ShouldNotBeNil)
			So(download.Progress, ShouldEqual, 1)
			So(download.Status, ShouldEqual, core.Complete)
			log.Println(download)
		})

	})
}

func validConfiguration() *Config {
	return &Config{
		RPCHost:     "localhost",
		RPCUser:     "testUser",
		RPCPassword: "testPass",
	}
}

func transmissionClientWithTestServer(handler http.Handler) (*TransmissionClient, *httptest.Server) {
	ts := httptest.NewServer(handler)
	url, err := url.Parse(ts.URL)
	if err != nil {
		log.Fatal(err)
	}
	config := &Config{
		RPCHost:     url.Hostname(),
		RPCPort:     url.Port(),
		RPCUser:     "testUser",
		RPCPassword: "testPass",
	}
	tc := NewTransmissionClient(config)
	return tc, ts
}
